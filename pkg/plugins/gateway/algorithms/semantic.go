package routingalgorithms

import (
	"fmt"
	"math/rand"

	"github.com/vllm-project/aibrix/pkg/types"
	"github.com/vllm-project/aibrix/pkg/utils"
	v1 "k8s.io/api/core/v1"
	"k8s.io/klog/v2"
)

const (
	RouterSemantic         types.RoutingAlgorithm = "semantic"
	SemanticRouteLabel     string                 = "aibrix.ai/semantic-route"
	HeaderSemanticDecision                        = "x-semantic-decision"
)

func init() {
	Register(RouterSemantic, NewSemanticRouter)
}

type semanticRouter struct{}

func NewSemanticRouter() (types.Router, error) {
	return &semanticRouter{}, nil
}

func (r *semanticRouter) Route(ctx *types.RoutingContext, readyPodList types.PodList) (string, error) {
	pods := readyPodList.All()
	if len(pods) == 0 {
		return "", fmt.Errorf("no ready pods available for semantic routing")
	}

	var decision string
	if ctx.ReqHeaders != nil {
		decision = ctx.ReqHeaders[HeaderSemanticDecision]
	}

	if decision == "" {
		klog.V(2).InfoS("semantic decision header not found, falling back to random", "request_id", ctx.RequestID)
		return r.fallback(ctx, pods)
	}

	targetPod := r.matchPodByDecision(pods, decision)
	if targetPod == nil {
		klog.V(2).InfoS("no pod matched semantic decision, falling back to random", "request_id", ctx.RequestID, "decision", decision)
		return r.fallback(ctx, pods)
	}

	klog.V(4).InfoS("semantic routing successful", "request_id", ctx.RequestID, "decision", decision, "target_pod", targetPod.Name)

	ctx.SetTargetPod(targetPod)
	return ctx.TargetAddress(), nil
}

func (r *semanticRouter) matchPodByDecision(pods []*v1.Pod, decision string) *v1.Pod {
	var candidates []*v1.Pod
	for _, pod := range pods {
		if val, ok := pod.Labels[SemanticRouteLabel]; ok && val == decision {
			candidates = append(candidates, pod)
		}
	}

	if len(candidates) == 0 {
		return nil
	}

	// TODO: change later
	return candidates[rand.Intn(len(candidates))]
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
