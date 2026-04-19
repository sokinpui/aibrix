package classification

import (
	"strings"
)

// applyAuthzFailOpenOnClassifyError clears Classify error when authz.fail_open is true and user ID
// is empty (e.g. Envoy stripped identity headers before ext_proc). Returns anonymous authz result.
func applyAuthzFailOpenOnClassifyError(failOpen bool, userID string, result *AuthzResult, err error) (*AuthzResult, error) {
	if err == nil {
		return result, nil
	}
	if failOpen && strings.TrimSpace(userID) == "" {
		return &AuthzResult{}, nil
	}
	return result, err
}
