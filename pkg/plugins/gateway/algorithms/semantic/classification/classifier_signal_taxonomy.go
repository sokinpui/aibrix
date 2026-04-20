package classification

import (
	"sync"

	"github.com/vllm-project/aibrix/pkg/plugins/gateway/algorithms/semantic/config"
)

// evaluateKBSignals runs all configured knowledge bases, records structured KB
// results, and maps routing.signals.kb bindings into normal matched signals.
func (c *Classifier) evaluateKBSignals(results *SignalResults, mu *sync.Mutex, text string) {
}

func kbSignalMatchConfidence(rule config.KBSignalRule, result *KBClassifyResult) (float64, bool) {
	return 0, false
}

func kbSignalMatchMode(rule config.KBSignalRule) string {
	return ""
}

func kbLabelMatchConfidence(labelName string, matchMode string, result *KBClassifyResult) (float64, bool) {
	return 0.0, false
}

func kbGroupMatchConfidence(groupName string, matchMode string, result *KBClassifyResult) (float64, bool) {
	return 0.0, false
}

func signalSetBestConfidence(confidences map[string]float64, signalType string, names []string) float64 {
	return 0.0
}
