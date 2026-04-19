package semantic

type Plugin interface {
	Name() string
	Execute(prompt string, features Features, decision *Decision) (string, error)
}

type PluginRegistry struct {
	plugins
