package routingalgorithms

import (
	"fmt"
	"math/rand"

	"github.com/vllm-project/aibrix/pkg/plugins/gateway/algorithms/semantic"
	"github.com/vllm-project/aibrix/pkg/types"
	"github.com/vllm-project/aibrix/pkg/utils"
	v1 "k8s.io/api/core/v1"
	"k8s.io/klog/v2"
)

const (
	RouterSemantic     types.RoutingAlgorithm = "semantic"
	SemanticRouteLabel string                 = "aibrix.ai/semantic-route"
)

func init() {
	Register(RouterSemantic, NewSemanticRouter)
}

type semanticRouter struct {
	selector semantic.Selector
}

func NewSemanticRouter() (types.Router, error) {
	// In a real scenario, these would be loaded from configuration
	return &semanticRouter{
		selector: semantic.NewLabelSelector(SemanticRouteLabel),
	}, nil
}

func (r *semanticRouter) Route(ctx *types.RoutingContext, readyPodList types.PodList) (string, error) {
	pods := readyPodList.All()
	if len(pods) == 0 {
		return "", fmt.Errorf("no ready pods available")
	}

	reqPrompt := ctx.Message
	reqModel := ctx.Model
	// reqBody := ctx.ReqBody
	// reqHeaders := ctx.ReqHeaders
	// reqUser := ctx.User

	klog.V(4).InfoS("Analyzing request for semantic routing",
		"request_id", ctx.RequestID,
		"model", reqModel,
		"prompt_length", len(reqPrompt))

	targetPod, err := r.selector.Select("", pods)
	if err != nil {
		return r.fallback(ctx, pods)
	}

	ctx.SetTargetPod(targetPod)
	return ctx.TargetAddress(), nil
}

func (r *semanticRouter) fallback(ctx *types.RoutingContext, pods []*v1.Pod) (string, error) {
	targetPod, err := utils.SelectRandomPod(pods, rand.Intn)
	if err != nil {
		return "", fmt.Errorf("random fallback selection failed: %w", err)
	}

	ctx.SetTargetPod(targetPod)
	return ctx.TargetAddress(), nil
}
