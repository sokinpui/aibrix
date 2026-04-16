package routingalgorithms

import (
	"context"
	"testing"

	"github.com/vllm-project/aibrix/pkg/types"
	"github.com/vllm-project/aibrix/pkg/utils"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestSemanticRouter(t *testing.T) {
	podGeneral := &v1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name: "pod-general",
			Labels: map[string]string{
				SemanticRouteLabel:     "general",
				"model.aibrix.ai/port": "8000",
			},
		},
		Status: v1.PodStatus{
			PodIP:      "10.0.0.1",
			Conditions: []v1.PodCondition{{Type: v1.PodReady, Status: v1.ConditionTrue}},
		},
	}
	podCoding := &v1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name: "pod-coding",
			Labels: map[string]string{
				SemanticRouteLabel:     "coding",
				"model.aibrix.ai/port": "8000",
			},
		},
		Status: v1.PodStatus{
			PodIP:      "10.0.0.2",
			Conditions: []v1.PodCondition{{Type: v1.PodReady, Status: v1.ConditionTrue}},
		},
	}
	podList := &utils.PodArray{Pods: []*v1.Pod{podGeneral, podCoding}}

	t.Run("Successful Routing to Coding", func(t *testing.T) {
		router, _ := NewSemanticRouter()
		ctx := types.NewRoutingContext(context.Background(), "id-1", "model", "python code", "/v1/completions", "user")

		// 模拟上游 Semantic Router ext_proc 注入了 Header
		ctx.ReqHeaders = map[string]string{
			HeaderSemanticDecision: "coding",
		}

		_, err := router.Route(ctx, podList)
		if err != nil {
			t.Fatalf("Route failed: %v", err)
		}
		if ctx.TargetPod() == nil || ctx.TargetPod().Name != "pod-coding" {
			t.Errorf("Expected pod-coding, got %v", ctx.TargetPod())
		}
	})

	t.Run("Fallback when Decision matches no pods", func(t *testing.T) {
		router, _ := NewSemanticRouter()
		ctx := types.NewRoutingContext(context.Background(), "id-2", "model", "1+1", "/v1/completions", "user")

		ctx.ReqHeaders = map[string]string{
			HeaderSemanticDecision: "math",
		}

		_, err := router.Route(ctx, podList)
		if err != nil {
			t.Fatalf("Route failed during fallback: %v", err)
		}
		if ctx.TargetPod() == nil {
			t.Errorf("Expected a fallback pod, got nil")
		}
	})

	t.Run("Fallback when Header is missing", func(t *testing.T) {
		router, _ := NewSemanticRouter()
		ctx := types.NewRoutingContext(context.Background(), "id-3", "model", "broken", "/v1/completions", "user")

		// Context headers 为空
		_, err := router.Route(ctx, podList)
		if err != nil {
			t.Fatalf("Route failed during header missing fallback: %v", err)
		}
		if ctx.TargetPod() == nil {
			t.Errorf("Expected a fallback pod on missing header, got nil")
		}
	})
}
