package semantic

import (
	"errors"
	"math/rand"

	v1 "k8s.io/api/core/v1"
)

var (
	ErrNoPodsAvailable = errors.New("no pods available")
	ErrNoPodsMatched   = errors.New("no pods found matching route label")
)

type Selector interface {
	Select(decision *Decision, pods []*v1.Pod) (*v1.Pod, error)
}

type labelSelector struct {
	labelKey string
}

func NewLabelSelector(labelKey string) Selector {
	return &labelSelector{
		labelKey: labelKey,
	}
}

func (s *labelSelector) Select(decision *Decision, pods []*v1.Pod) (*v1.Pod, error) {
	if len(pods) == 0 {
		return nil, ErrNoPodsAvailable
	}

	if decision == nil || decision.RouteLabel == "" {
		return pods[rand.Intn(len(pods))], nil
	}

	var matched []*v1.Pod
	for _, pod := range pods {
		if pod.Labels != nil && pod.Labels[s.labelKey] == decision.RouteLabel {
			matched = append(matched, pod)
		}
	}

	if len(matched) == 0 {
		return nil, ErrNoPodsMatched
	}

	return matched[rand.Intn(len(matched))], nil
}
