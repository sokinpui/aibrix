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
	RouterSemantic     types.RoutingAlgorithm = "semantic"
	SemanticRouteLabel string                 = "aibrix.ai/semantic-route"
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
		return "", fmt.Errorf("no ready pods available")
	}

	klog.V(4).InfoS("Executing semantic routing logic", "request_id", ctx.RequestID)

	return r.fallback(ctx, pods)
}

func (r *semanticRouter) fallback(ctx *types.RoutingContext, pods []*v1.Pod) (string, error) {
	targetPod, err := utils.SelectRandomPod(pods, rand.Intn)
	if err != nil {
		return "", fmt.Errorf("random fallback selection failed: %w", err)
	}

	ctx.SetTargetPod(targetPod)
	return ctx.TargetAddress(), nil
}
