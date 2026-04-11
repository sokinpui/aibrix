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
	"bytes"
	"context"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"time"

	"github.com/bytedance/sonic"
	"github.com/vllm-project/aibrix/pkg/types"
	"github.com/vllm-project/aibrix/pkg/utils"
	v1 "k8s.io/api/core/v1"
	"k8s.io/klog/v2"
)

const (
	RouterSemantic     types.RoutingAlgorithm = "semantic"
	SemanticRouteLabel string                 = "aibrix.ai/semantic-route"
)

var (
	semanticRouterTimeout = utils.LoadEnvInt("AIBRIX_SEMANTIC_ROUTER_TIMEOUT", 2)
	semanticRouterBaseURL = utils.LoadEnv("AIBRIX_SEMANTIC_ROUTER_URL", "")
	// Task type: combined, intent, pii, security, etc.
	semanticRouterTask = utils.LoadEnv("AIBRIX_SEMANTIC_ROUTER_TASK", "combined")
)

func init() {
	Register(RouterSemantic, NewSemanticRouter)
}

type semanticRouter struct {
	httpClient *http.Client
	apiURL     string
}

// semanticRequest matches the classification API expected by vllm-semantic-router apiserver
type semanticRequest struct {
	Text  string `json:"text"`
	Model string `json:"model,omitempty"`
}

// semanticResponse captures the decision name from the classification result
type semanticResponse struct {
	Decision string `json:"decision"`
}

func NewSemanticRouter() (types.Router, error) {
	httpClient := &http.Client{
		Timeout: time.Duration(semanticRouterTimeout) * time.Second,
		Transport: &http.Transport{
			MaxIdleConns:        100,
			MaxIdleConnsPerHost: 20,
			IdleConnTimeout:     90 * time.Second,
		},
	}

	apiURL := fmt.Sprintf("%s/api/v1/classify/%s", semanticRouterBaseURL, semanticRouterTask)

	return &semanticRouter{
		httpClient: httpClient,
		apiURL:     apiURL,
	}, nil
}

func (r *semanticRouter) Route(ctx *types.RoutingContext, readyPodList types.PodList) (string, error) {
	pods := readyPodList.All()
	if len(pods) == 0 {
		return "", fmt.Errorf("no ready pods available for semantic routing")
	}

	decision, err := r.classifyRequest(ctx)
	if err != nil {
		klog.V(2).InfoS("semantic classification failed, falling back to random",
			"request_id", ctx.RequestID,
			"error", err)
		return r.fallback(ctx, pods)
	}

	targetPod := r.matchPodByDecision(pods, decision)
	if targetPod == nil {
		klog.V(2).InfoS("no pod matched semantic decision, falling back to random",
			"request_id", ctx.RequestID,
			"decision", decision)
		return r.fallback(ctx, pods)
	}

	klog.V(4).InfoS("semantic routing successful",
		"request_id", ctx.RequestID,
		"decision", decision,
		"target_pod", targetPod.Name)

	ctx.SetTargetPod(targetPod)
	return ctx.TargetAddress(), nil
}

func (r *semanticRouter) classifyRequest(ctx *types.RoutingContext) (string, error) {
	if semanticRouterBaseURL == "" {
		return "", fmt.Errorf("AIBRIX_SEMANTIC_ROUTER_URL is not configured")
	}

	payload, err := sonic.Marshal(semanticRequest{
		Text:  ctx.Message,
		Model: ctx.Model,
	})
	if err != nil {
		return "", err
	}

	httpCtx, cancel := context.WithTimeout(ctx.Context, time.Duration(semanticRouterTimeout)*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(httpCtx, "POST", r.apiURL, bytes.NewBuffer(payload))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Request-Id", ctx.RequestID)

	resp, err := r.httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("semantic engine error (status %d): %s", resp.StatusCode, string(body))
	}

	var result semanticResponse
	if err := sonic.ConfigDefault.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("failed to decode semantic response: %w", err)
	}

	return result.Decision, nil
}

func (r *semanticRouter) matchPodByDecision(pods []*v1.Pod, decision string) *v1.Pod {
	if decision == "" {
		return nil
	}

	var candidates []*v1.Pod
	for _, pod := range pods {
		if val, ok := pod.Labels[SemanticRouteLabel]; ok && val == decision {
			candidates = append(candidates, pod)
		}
	}

	if len(candidates) == 0 {
		return nil
	}

	// Shuffle to provide basic load balancing among pods serving the same decision category
	return candidates[rand.Intn(len(candidates))]
}

func (r *semanticRouter) fallback(ctx *types.RoutingContext, pods []*v1.Pod) (string, error) {
	targetPod, err := utils.SelectRandomPod(pods, rand.Intn)
	if err != nil {
		return "", fmt.Errorf("random fallback selection failed: %w", err)
	}

	if targetPod == nil {
		return "", fmt.Errorf("no pods available for fallback")
	}

	ctx.SetTargetPod(targetPod)
	return ctx.TargetAddress(), nil
}

func (r *semanticRouter) SubscribedMetrics() []string {
	return []string{}
}
