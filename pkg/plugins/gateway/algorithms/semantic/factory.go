/*
Copyright 2025 vLLM Semantic Router.

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

package semantic

import (
	"k8s.io/klog/v2"
)

// DefaultModelSelectionConfig returns the default configuration
func DefaultModelSelectionConfig() *ModelSelectionConfig {
	return &ModelSelectionConfig{
		Method: "static",
	}
}

// Factory creates and initializes selectors based on configuration
type Factory struct {
	cfg           *ModelSelectionConfig
	modelConfig   map[string]ModelParams
	categories    []Category
	embeddingFunc func(string) ([]float32, error)
}

// NewFactory creates a new selector factory
func NewFactory(cfg *ModelSelectionConfig) *Factory {
	if cfg == nil {
		cfg = DefaultModelSelectionConfig()
	}
	return &Factory{
		cfg: cfg,
	}
}

// WithModelConfig sets the model configuration
func (f *Factory) WithModelConfig(modelConfig map[string]ModelParams) *Factory {
	f.modelConfig = modelConfig
	return f
}

// WithCategories sets the category configuration
func (f *Factory) WithCategories(categories []Category) *Factory {
	f.categories = categories
	return f
}

// WithEmbeddingFunc sets the embedding function for RouterDC
func (f *Factory) WithEmbeddingFunc(fn func(string) ([]float32, error)) *Factory {
	f.embeddingFunc = fn
	return f
}

// Create creates and initializes a selector based on the configured method
func (f *Factory) Create() Selector {
	method := SelectionMethod(f.cfg.Method)

	var selector Selector

	switch method {
	case MethodRouterDC:
		routerDCSelector := NewRouterDCSelector(toRouterDCConfig(f.cfg.RouterDC))
		if f.embeddingFunc != nil {
			routerDCSelector.SetEmbeddingFunc(f.embeddingFunc)
		}
		// Initialize model embeddings from descriptions in model config
		if f.modelConfig != nil {
			if err := routerDCSelector.InitializeFromConfig(f.modelConfig); err != nil {
				klog.Errorf("[SelectionFactory] RouterDC initialization failed: %v", err)
			}
		}
		selector = routerDCSelector

	default:
		// Default to static selector
		staticSelector := NewStaticSelector(DefaultStaticConfig())
		if f.categories != nil {
			staticSelector.InitializeFromConfig(f.categories)
		}
		selector = staticSelector
	}

	klog.Infof("[SelectionFactory] Created selector: method=%s", method)
	return selector
}

// CreateAll creates all available selectors and registers them
func (f *Factory) CreateAll() *Registry {
	registry := NewRegistry()

	// Always create static selector
	staticSelector := NewStaticSelector(DefaultStaticConfig())
	if f.categories != nil {
		staticSelector.InitializeFromConfig(f.categories)
	}
	registry.Register(MethodStatic, staticSelector)

	// Create RouterDC selector
	routerDCSelector := NewRouterDCSelector(toRouterDCConfig(f.cfg.RouterDC))
	if f.embeddingFunc != nil {
		routerDCSelector.SetEmbeddingFunc(f.embeddingFunc)
	}
	// Initialize model embeddings from descriptions in model config
	if f.modelConfig != nil {
		if err := routerDCSelector.InitializeFromConfig(f.modelConfig); err != nil {
			klog.Errorf("[SelectionFactory] RouterDC initialization failed: %v", err)
		}
	}
	registry.Register(MethodRouterDC, routerDCSelector)

	LogRegisteredAlgorithms(registry)
	klog.InfoS("selection_factory_initialized", "component", "selection", "details", map[string]interface{}{
		"selector_count": len(registry.selectors),
	})
	return registry
}

// LogRegisteredAlgorithms logs the tier and dependencies of each registered algorithm
func LogRegisteredAlgorithms(registry *Registry) {
	registry.mu.RLock()
	defer registry.mu.RUnlock()

	for method, selector := range registry.selectors {
		klog.Infof("[Selection] Registered algorithm: %s (tier=%s)", method, selector.Tier())
	}
}

// Initialize sets up the global registry with all selectors
func Initialize(cfg *ModelSelectionConfig, modelConfig map[string]ModelParams, categories []Category, embeddingFunc func(string) ([]float32, error)) {
	factory := NewFactory(cfg).
		WithModelConfig(modelConfig).
		WithCategories(categories).
		WithEmbeddingFunc(embeddingFunc)

	// Create all selectors and register globally
	GlobalRegistry = factory.CreateAll()

	klog.InfoS("selection_registry_initialized", "component", "selection", "details", map[string]interface{}{
		"selector_count": len(GlobalRegistry.selectors),
	})
}

// GetSelector returns a selector for the specified method from global registry
func GetSelector(method SelectionMethod) Selector {
	selector, ok := GlobalRegistry.Get(method)
	if !ok {
		// Fallback to static
		selector, _ = GlobalRegistry.Get(MethodStatic)
	}
	return selector
}

func toRouterDCConfig(cfg RouterDCSelectionConfig) *RouterDCConfig {
	return &RouterDCConfig{
		Temperature:         cfg.Temperature,
		MinSimilarity:       cfg.MinSimilarity,
		UseCapabilities:     cfg.UseCapabilities,
		RequireDescriptions: cfg.RequireDescriptions,
		DimensionSize:       768, // Default
	}
}
