package semantic

// This file centralizes Tier() and ExternalDependencies() implementations
// for all selection algorithms. Keeping them in one place avoids adding lines
// to algorithm files that are already at or over the 800-line hard limit.

// --- Supported-tier algorithms ---

// Tier returns the production readiness tier
func (s *StaticSelector) Tier() AlgorithmTier {
	return TierSupported
}

// Tier returns the production readiness tier
func (r *RouterDCSelector) Tier() AlgorithmTier {
	return TierSupported
}
