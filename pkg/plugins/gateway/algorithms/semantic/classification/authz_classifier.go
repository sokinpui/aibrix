package classification

import (
	"github.com/vllm-project/aibrix/pkg/plugins/gateway/algorithms/semantic/config"
)

// AuthzResult represents the result of authz signal classification.
// It contains only the matched role names — the decision engine uses these
// to select models via modelRefs.
type AuthzResult struct {
	// MatchedRules contains the role names from all matched role bindings.
	// These are the Role field values, not the binding Name field.
	MatchedRules []string
}

// normalizedBinding is the internal representation with Kind/Name already normalized.
// This avoids repeated string normalization at request time.
type normalizedBinding struct {
	name     string // binding name (for logs)
	role     string // role name (emitted as signal)
	subjects []normalizedSubject
}

type normalizedSubject struct {
	kind string // "user" or "group" (already lowercased)
	name string // already trimmed
}

// AuthzClassifier evaluates user identity and group membership against RBAC role bindings.
// It follows the Kubernetes RoleBinding pattern:
//   - Subject  → user ID + groups (from auth backend: Authorino, Envoy Gateway JWT, etc.)
//   - Role     → RoleBinding.Role (emitted as the signal name)
//   - Permission → decision engine modelRefs (not this classifier's concern)
type AuthzClassifier struct {
	bindings []normalizedBinding
}

// NewAuthzClassifier creates a new AuthzClassifier from RBAC role bindings.
// All validation and normalization happens here at startup. If this function returns
// without error, Classify() is guaranteed to work correctly at request time.
//
// Validates at startup:
//   - Binding name must not be empty
//   - Binding name must be unique across all bindings
//   - Role must not be empty
//   - At least one subject must be specified
//   - Each subject must have kind "User" or "Group" (case-insensitive)
//   - Each subject must have a non-empty name (whitespace-only is rejected)
//
// Normalizes at startup:
//   - Subject.Kind is lowercased and trimmed
//   - Subject.Name is trimmed (preserving original case for exact matching)
func NewAuthzClassifier(bindings []config.RoleBinding) (*AuthzClassifier, error) {
	return &AuthzClassifier{}, nil
}

// Classify evaluates the RBAC role bindings against the user identity and groups.
//
// Match logic: a binding matches if ANY of its subjects match:
//   - kind: "user"  → matches if subject.name == userID
//   - kind: "group" → matches if subject.name is in userGroups
//
// When a binding matches, its role is emitted as the signal name.
// Multiple bindings can match. If multiple bindings grant the same role, it is deduplicated.
//
// Returns an error if userID is empty and role bindings are configured — this prevents
// silent bypass when ext_authz fails to inject the user identity header.
func (c *AuthzClassifier) Classify(userID string, userGroups []string) (*AuthzResult, error) {
	return &AuthzResult{}, nil
}

// ParseUserGroups parses a comma-separated groups header value into a slice of group names.
// Whitespace around group names is trimmed. Empty strings are excluded.
func ParseUserGroups(headerValue string) []string {
	return nil
}
