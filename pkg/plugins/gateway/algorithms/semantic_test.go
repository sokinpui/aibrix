/*
Copyright 2024 The Aibrix Team.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package routingalgorithms

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
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
			PodIP: "10.0.0.1",
			Conditions: []v1.PodCondition{
				{
					Type:   v1.PodReady,
					Status: v1.ConditionTrue,
				},
			},
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
			PodIP: "10.0.0.2",
			Conditions: []v1.PodCondition{
				{
					Type:   v1.PodReady,
					Status: v1.ConditionTrue,
				},
			},
		},
	}
	podList := &utils.PodArray{Pods: []*v1.Pod{podGeneral, podCoding}}

	t.Run("Successful Routing to Coding", func(t *testing.T) {
		mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(semanticResponse{Decision: "coding"})
		}))
		defer mockServer.Close()

		semanticRouterBaseURL = mockServer.URL
		router, _ := NewSemanticRouter()
		ctx := types.NewRoutingContext(context.Background(), "id-1", "model", "python code", "/v1/completions", "user")

		_, err := router.Route(ctx, podList)
		if err != nil {
			t.Fatalf("Route failed: %v", err)
		}
		if ctx.TargetPod() == nil || ctx.TargetPod().Name != "pod-coding" {
			t.Errorf("Expected pod-coding, got %v", ctx.TargetPod())
		}
	})

	t.Run("Fallback when Decision matches no pods", func(t *testing.T) {
		mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			// "math" decision does not exist on any pods
			json.NewEncoder(w).Encode(semanticResponse{Decision: "math"})
		}))
		defer mockServer.Close()

		semanticRouterBaseURL = mockServer.URL
		router, _ := NewSemanticRouter()
		ctx := types.NewRoutingContext(context.Background(), "id-2", "model", "1+1", "/v1/completions", "user")

		// Should not error, should fallback to random selection among ready pods
		_, err := router.Route(ctx, podList)
		if err != nil {
			t.Fatalf("Route failed during fallback: %v", err)
		}
		if ctx.TargetPod() == nil {
			t.Errorf("Expected a fallback pod, got nil")
		}
	})

	t.Run("Fallback when Semantic Router API fails", func(t *testing.T) {
		mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
		}))
		defer mockServer.Close()

		semanticRouterBaseURL = mockServer.URL
		router, _ := NewSemanticRouter()
		ctx := types.NewRoutingContext(context.Background(), "id-3", "model", "broken", "/v1/completions", "user")

		// Should not error, should fallback to random selection
		_, err := router.Route(ctx, podList)
		if err != nil {
			t.Fatalf("Route failed during API error fallback: %v", err)
		}
		if ctx.TargetPod() == nil {
			t.Errorf("Expected a fallback pod on API error, got nil")
		}
	})
}
