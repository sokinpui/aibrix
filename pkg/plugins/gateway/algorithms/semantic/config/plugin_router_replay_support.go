package config

const (
	defaultRouterReplayMaxRecords   = 10000
	defaultRouterReplayMaxBodyBytes = 4096
)

func DefaultRouterReplayPluginConfig() RouterReplayPluginConfig {
	return RouterReplayPluginConfig{
		Enabled:             true,
		MaxRecords:          defaultRouterReplayMaxRecords,
		CaptureRequestBody:  true,
		CaptureResponseBody: true,
		MaxBodyBytes:        defaultRouterReplayMaxBodyBytes,
	}
}

// EffectiveRouterReplayConfigForDecision returns the replay configuration that
// should apply to a decision after layering global enablement and any
// per-decision router_replay plugin overrides.
func (c *RouterConfig) EffectiveRouterReplayConfigForDecision(decisionName string) *RouterReplayPluginConfig {
	base := DefaultRouterReplayPluginConfig()
	if c == nil {
		return &base
	}
	base.Enabled = c.RouterReplay.Enabled

	decision := c.GetDecisionByName(decisionName)
	if decision == nil {
		if base.Enabled {
			return &base
		}
		return nil
	}

	plugin := decision.GetPlugin(DecisionPluginRouterReplay)
	if plugin == nil {
		if base.Enabled {
			return &base
		}
		return nil
	}
	if plugin.Configuration == nil {
		if base.Enabled {
			return &base
		}
		return nil
	}

	if err := UnmarshalPluginConfig(plugin.Configuration, &base); err != nil {
		return nil
	}
	if !base.Enabled {
		return nil
	}
	return &base
}
