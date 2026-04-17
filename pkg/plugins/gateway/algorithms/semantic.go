package routingalgorithms

import (
	"encoding/json"
	"fmt"
	"math/rand"

	"github.com/vllm-project/aibrix/pkg/constants"
	"github.com/vllm-project/aibrix/pkg/plugins/gateway/algorithms/semantic"
	"github.com/vllm-project/aibrix/pkg/types"
	"github.com/vllm-project/aibrix/pkg/utils"
	v1 "k8s.io/api/core/v1"
	"k8s.io/klog/v2"
)

const (
	RouterSemantic         types.RoutingAlgorithm = "semantic"
	SemanticRouteLabel     string                 = "aibrix.ai/semantic-route"
	HeaderSemanticDecision string                 = "x-semantic-decision"
)

func init() {
	Register(RouterSemantic, NewSemanticRouter)
}

type semanticRouter struct{}

func NewSemanticRouter() (types.Router, error) {
	cfg := semantic.DefaultModelSelectionConfig()
	semantic.Initialize(cfg, nil, nil, nil)
	return &semanticRouter{}, nil
}

func (r *semanticRouter) Route(ctx *types.RoutingContext, readyPodList types.PodList) (string, error) {
	pods := readyPodList.All()
	if len(pods) == 0 {
		return "", fmt.Errorf("no ready pods available for semantic routing")
	}

	candidateModels, modelToPods := r.groupPodsByModel(pods)
	if len(candidateModels) == 0 {
		klog.V(2).InfoS("no candidate models found in pods, falling back to random", "request_id", ctx.RequestID)
		return r.fallback(ctx, pods)
	}

	method := r.resolveSelectionMethod(ctx)
	selector := semantic.GetSelector(method)

	selCtx := &semantic.SelectionContext{
		Query:           ctx.Message,
		CandidateModels: candidateModels,
		DecisionName:    r.getDecisionName(ctx),
	}

	result, err := selector.Select(ctx.Context, selCtx)
	if err != nil {
		klog.ErrorS(err, "semantic selection failed, falling back to random", "request_id", ctx.RequestID)
		return r.fallback(ctx, pods)
	}

	klog.V(4).InfoS("semantic routing successful",
		"request_id", ctx.RequestID,
		"selected_model", result.SelectedModel,
		"method", result.Method,
		"score", result.Score,
		"confidence", result.Confidence)

	targetPods := modelToPods[result.SelectedModel]
	if len(targetPods) == 0 {
		klog.V(2).InfoS("selected model has no ready pods, falling back to random", "request_id", ctx.RequestID, "selected_model", result.SelectedModel)
		return r.fallback(ctx, pods)
	}

	targetPod, err := utils.SelectRandomPod(targetPods, rand.Intn)
	if err != nil {
		return "", fmt.Errorf("failed to select pod for model %s: %w", result.SelectedModel, err)
	}

	ctx.SetTargetPod(targetPod)
	return ctx.TargetAddress(), nil
}

func (r *semanticRouter) groupPodsByModel(pods []*v1.Pod) ([]semantic.ModelRef, map[string][]*v1.Pod) {
	var candidateModels []semantic.ModelRef
	modelToPods := make(map[string][]*v1.Pod)

	for _, pod := range pods {
		model := r.getModelNameFromPod(pod)
		if model == "" {
			continue
		}

		if _, exists := modelToPods[model]; !exists {
			candidateModels = append(candidateModels, semantic.ModelRef{Model: model})
		}
		modelToPods[model] = append(modelToPods[model], pod)
	}
	return candidateModels, modelToPods
}

func (r *semanticRouter) getModelNameFromPod(pod *v1.Pod) string {
	if model, ok := pod.Labels[constants.ModelLabelName]; ok && model != "" {
		return model
	}
	if model, ok := pod.Annotations[constants.ModelLabelName]; ok && model != "" {
		return model
	}
	return pod.Labels[SemanticRouteLabel]
}

func (r *semanticRouter) resolveSelectionMethod(ctx *types.RoutingContext) semantic.SelectionMethod {
	method := semantic.MethodStatic
	if ctx.ConfigProfile == nil || len(ctx.ConfigProfile.RoutingConfig) == 0 {
		return method
	}

	var cfg struct {
		Method string `json:"method"`
	}
	if err := json.Unmarshal(ctx.ConfigProfile.RoutingConfig, &cfg); err == nil && cfg.Method != "" {
		method = semantic.SelectionMethod(cfg.Method)
	}
	return method
}

func (r *semanticRouter) getDecisionName(ctx *types.RoutingContext) string {
	if ctx.ReqHeaders == nil {
		return ""
	}
	return ctx.ReqHeaders[HeaderSemanticDecision]
}

func (r *semanticRouter) getUserID(ctx *types.RoutingContext) string {
	if ctx.User == nil {
		return ""
	}
	return *ctx.User
}

func (r *semanticRouter) fallback(ctx *types.RoutingContext, pods []*v1.Pod) (string, error) {
	targetPod, err := utils.SelectRandomPod(pods, rand.Intn)
	if err != nil {
		return "", fmt.Errorf("random fallback selection failed: %w", err)
	}

	ctx.SetTargetPod(targetPod)
	return ctx.TargetAddress(), nil
}

func (r *semanticRouter) SubscribedMetrics() []string {
	return []string{}
}
