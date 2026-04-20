package classification

import (
	"encoding/json"
	"errors"
	"fmt"
	"slices"
	"strings"
	"sync"
	"time"

	"github.com/vllm-project/aibrix/pkg/plugins/gateway/algorithms/semantic/config"
)

// cachedJailbreakResult stores a cached jailbreak classification result.
type cachedJailbreakResult struct {
	result InferenceClassResult
	err    error
}

// cachedPIIResult stores a cached PII token classification result.
type cachedPIIResult struct {
	result InferenceTokenResult
	err    error
}

// EvaluateAllSignalsWithContext evaluates all signal types with separate text for context counting.
//
// text: (possibly compressed) text for signal evaluation
// contextText: text for context token counting (usually all messages combined)
// nonUserMessages: conversation history for jailbreak/PII with include_history
// forceEvaluateAll: if true, evaluates all configured signals regardless of decision usage
// uncompressedText: original text before prompt compression (empty = no compression happened)
// skipCompressionSignals: signal types that must use uncompressedText instead of text
// imageURL: optional image URL for multimodal signals
// signalReadiness returns a map indicating whether each signal type's infrastructure is ready.
// Separated from EvaluateAllSignalsWithContext to keep cyclomatic complexity under the linter limit.
func (c *Classifier) signalReadiness() map[string]bool {
	return map[string]bool{
		config.SignalTypeKeyword:      c.keywordClassifier != nil,
		config.SignalTypeEmbedding:    c.keywordEmbeddingClassifier != nil,
		config.SignalTypeDomain:       c.IsCategoryEnabled() && c.categoryInference != nil && c.CategoryMapping != nil,
		config.SignalTypeFactCheck:    len(c.Config.FactCheckRules) > 0 && c.IsFactCheckEnabled(),
		config.SignalTypeUserFeedback: len(c.Config.UserFeedbackRules) > 0 && c.IsFeedbackDetectorEnabled(),
		config.SignalTypeReask:        c.reaskClassifier != nil,
		config.SignalTypePreference:   len(c.Config.PreferenceRules) > 0 && c.IsPreferenceClassifierEnabled(),
		config.SignalTypeLanguage:     len(c.Config.LanguageRules) > 0 && c.IsLanguageEnabled(),
		config.SignalTypeContext:      c.contextClassifier != nil,
		config.SignalTypeStructure:    c.structureClassifier != nil,
		config.SignalTypeComplexity:   c.complexityClassifier != nil,
		config.SignalTypeModality:     len(c.Config.ModalityRules) > 0 && c.Config.ModalityDetector.Enabled,
		config.SignalTypeJailbreak:    len(c.Config.JailbreakRules) > 0 && c.IsJailbreakEnabled(),
		config.SignalTypePII:          len(c.Config.PIIRules) > 0 && c.IsPIIEnabled(),
		config.SignalTypeKB:           len(c.kbClassifiers) > 0,
	}
}

// textForSignalFunc returns a function that resolves the correct text for a given signal type,
// using uncompressed text for signals that must not receive compressed input.
func textForSignalFunc(text, uncompressedText string, skipCompressionSignals map[string]bool) func(string) string {
	return func(signalType string) string {
		if uncompressedText != "" && skipCompressionSignals[signalType] {
			return uncompressedText
		}
		return text
	}
}

func (c *Classifier) EvaluateAllSignalsWithContext(text string, contextText string, currentUserText string, priorUserMessages []string, nonUserMessages []string, hasPriorAssistantReply bool, forceEvaluateAll bool, uncompressedText string, skipCompressionSignals map[string]bool, imageURL ...string) *SignalResults {
	defer c.enterSignalEvaluationLoadGate()()
	// Determine which signals (type:name) should be evaluated
	var usedSignals map[string]bool
	if forceEvaluateAll {
		usedSignals = c.getAllSignalTypes()
	} else {
		usedSignals = c.getUsedSignals()
	}

	textForSignal := textForSignalFunc(text, uncompressedText, skipCompressionSignals)
	ready := c.signalReadiness()

	results := &SignalResults{
		Metrics:           &SignalMetricsCollection{},
		SignalConfidences: make(map[string]float64),
		SignalValues:      make(map[string]float64),
	}

	var wg sync.WaitGroup
	var mu sync.Mutex
	imgArg := ""
	if len(imageURL) > 0 {
		imgArg = imageURL[0]
	}

	dispatchers := c.buildSignalDispatchers(results, &mu, textForSignal, contextText, currentUserText, priorUserMessages, nonUserMessages, hasPriorAssistantReply, imgArg)
	runSignalDispatchers(dispatchers, usedSignals, ready, &wg)

	wg.Wait()
	results = c.applySignalGroups(results)
	results = c.applySignalComposers(results)
	results = c.applySignalOutputPolicies(results)
	results = c.applyProjections(results)
	return results
}

func (c *Classifier) evaluateKeywordSignal(results *SignalResults, mu *sync.Mutex, text string) {
	start := time.Now()
	category, keywords, err := c.keywordClassifier.ClassifyWithKeywords(text)
	elapsed := time.Since(start)

	// Record metrics (use microseconds for better precision)
	results.Metrics.Keyword.ExecutionTimeMs = float64(elapsed.Microseconds()) / 1000.0
	results.Metrics.Keyword.Confidence = 1.0 // Rule-based, always 1.0

	if err != nil {
	} else if category != "" {
		mu.Lock()
		results.MatchedKeywordRules = append(results.MatchedKeywordRules, category)
		results.MatchedKeywords = append(results.MatchedKeywords, keywords...)
		mu.Unlock()
	}
}

func (c *Classifier) evaluateEmbeddingSignal(results *SignalResults, mu *sync.Mutex, text string) {
	start := time.Now()
	detailedResult, err := c.keywordEmbeddingClassifier.ClassifyDetailed(text)
	elapsed := time.Since(start)

	// Record metrics
	results.Metrics.Embedding.ExecutionTimeMs = float64(elapsed.Microseconds()) / 1000.0

	if err != nil {
		return
	}
	if detailedResult == nil {
		return
	}

	var bestConfidence float64
	for _, score := range detailedResult.Scores {
		if score.Score > bestConfidence {
			bestConfidence = score.Score
		}
	}
	results.Metrics.Embedding.Confidence = bestConfidence

	mu.Lock()
	for _, score := range detailedResult.Scores {
		results.SignalValues["embedding:"+score.Name] = score.Score
		results.SignalValues["embedding:"+score.Name+":best"] = score.Best
		results.SignalValues["embedding:"+score.Name+":support"] = score.Support
		results.SignalValues["embedding:"+score.Name+":prototype_count"] = float64(score.PrototypeCount)
	}
	for _, mr := range detailedResult.Matches {
		results.MatchedEmbeddingRules = append(results.MatchedEmbeddingRules, mr.RuleName)
		results.SignalConfidences["embedding:"+mr.RuleName] = mr.Score
	}
	mu.Unlock()
}

func (c *Classifier) evaluateDomainSignal(results *SignalResults, mu *sync.Mutex, text string) {
	start := time.Now()
	domainResult, err := c.categoryInference.ClassifyWithProbabilities(text)
	if err != nil {
		// Fall back to Classify() (top-1 only) when ClassifyWithProbabilities is unavailable.
		basicResult, basicErr := c.categoryInference.Classify(text)
		if basicErr != nil {
			err = basicErr
		} else {
			domainResult = InferenceClassResultWithProbs{
				Class:      basicResult.Class,
				Confidence: basicResult.Confidence,
			}
			err = nil
		}
	}
	elapsed := time.Since(start)

	categoryName := ""
	if err == nil {
		if name, ok := c.CategoryMapping.GetCategoryFromIndex(domainResult.Class); ok {
			categoryName = c.translateMMLUToGeneric(name)
		}
	}

	// Record metrics
	results.Metrics.Domain.ExecutionTimeMs = float64(elapsed.Microseconds()) / 1000.0
	if categoryName != "" && err == nil {
		results.Metrics.Domain.Confidence = float64(domainResult.Confidence)
	}
	if err != nil {
	} else {
		matched := c.matchDomainCategories(domainResult, categoryName)
		mu.Lock()
		for _, cat := range matched {
			results.MatchedDomainRules = append(results.MatchedDomainRules, cat.Category)
			results.SignalConfidences["domain:"+cat.Category] = float64(cat.Probability)
		}
		mu.Unlock()
	}
}

func (c *Classifier) evaluateFactCheckSignal(results *SignalResults, mu *sync.Mutex, text string) {
	start := time.Now()
	factCheckResult, err := c.ClassifyFactCheck(text)
	elapsed := time.Since(start)

	// Determine which signal to output based on classification result
	signalName := "no_fact_check_needed"
	if err == nil && factCheckResult != nil && factCheckResult.NeedsFactCheck {
		signalName = "needs_fact_check"
	}

	// Record metrics (use microseconds for better precision)
	results.Metrics.FactCheck.ExecutionTimeMs = float64(elapsed.Microseconds()) / 1000.0
	if signalName != "" && err == nil && factCheckResult != nil {
		results.Metrics.FactCheck.Confidence = float64(factCheckResult.Confidence)
	}

	if err != nil {
	} else if factCheckResult != nil {
		// Check if this signal is defined in fact_check_rules
		for _, rule := range c.Config.FactCheckRules {
			if rule.Name == signalName {
				mu.Lock()
				results.MatchedFactCheckRules = append(results.MatchedFactCheckRules, rule.Name)
				mu.Unlock()
				break
			}
		}
	}
}

func (c *Classifier) evaluateUserFeedbackSignal(results *SignalResults, mu *sync.Mutex, text string, hasPriorAssistantReply bool) {
	if !shouldEvaluateUserFeedbackSignal(hasPriorAssistantReply) {
		return
	}

	start := time.Now()
	feedbackResult, err := c.ClassifyFeedback(text)
	elapsed := time.Since(start)

	// Use the feedback type directly as the signal name
	signalName := ""
	if err == nil && feedbackResult != nil {
		signalName = feedbackResult.FeedbackType
	}

	// Record metrics (use microseconds for better precision)
	results.Metrics.UserFeedback.ExecutionTimeMs = float64(elapsed.Microseconds()) / 1000.0
	if signalName != "" && err == nil && feedbackResult != nil {
		results.Metrics.UserFeedback.Confidence = float64(feedbackResult.Confidence)
	}

	if err != nil {
	} else if feedbackResult != nil {
		// Check if this signal is defined in user_feedback_rules
		for _, rule := range c.Config.UserFeedbackRules {
			if rule.Name == signalName {
				mu.Lock()
				results.MatchedUserFeedbackRules = append(results.MatchedUserFeedbackRules, rule.Name)
				mu.Unlock()
				break
			}
		}
	}
}

func (c *Classifier) evaluateReaskSignal(results *SignalResults, mu *sync.Mutex, currentUserText string, priorUserMessages []string) {
	start := time.Now()
	matchedRules, err := c.reaskClassifier.Classify(currentUserText, priorUserMessages)
	elapsed := time.Since(start)

	results.Metrics.Reask.ExecutionTimeMs = float64(elapsed.Microseconds()) / 1000.0

	if err != nil {
		return
	}
	if len(matchedRules) == 0 {
		return
	}

	bestConfidence := 0.0
	mu.Lock()
	for _, match := range matchedRules {
		if match.MinSimilarity > bestConfidence {
			bestConfidence = match.MinSimilarity
		}
		results.MatchedReaskRules = append(results.MatchedReaskRules, match.RuleName)
		results.SignalConfidences["reask:"+match.RuleName] = match.MinSimilarity
		results.SignalValues["reask:"+match.RuleName] = float64(match.MatchedTurns)
	}
	results.Metrics.Reask.Confidence = bestConfidence
	mu.Unlock()
}

func (c *Classifier) evaluatePreferenceSignal(results *SignalResults, mu *sync.Mutex, text string) {
	start := time.Now()
	contentBytes, _ := json.Marshal(text)
	conversationJSON := fmt.Sprintf(`[{"role":"user","content":%s}]`, contentBytes)

	preferenceResult, err := c.preferenceClassifier.Classify(conversationJSON)
	elapsed := time.Since(start)

	// Use the preference name directly as the signal name
	preferenceName := ""
	if err == nil && preferenceResult != nil {
		preferenceName = preferenceResult.Preference
	}

	// Record metrics (use microseconds for better precision)
	results.Metrics.Preference.ExecutionTimeMs = float64(elapsed.Microseconds()) / 1000.0
	if preferenceName != "" && err == nil && preferenceResult != nil && preferenceResult.Confidence > 0 {
		results.Metrics.Preference.Confidence = float64(preferenceResult.Confidence)
	}

	if err != nil {
		if errors.Is(err, ErrPreferenceBelowThreshold) {
			return
		}
		return
	}
	if preferenceResult == nil || preferenceName == "" || !c.hasConfiguredPreferenceRule(preferenceName) {
		return
	}

	recordPreferenceMatchMetrics(preferenceName)
	details := c.contrastivePreferenceDetails(text)

	mu.Lock()
	recordPreferenceSignalValues(results, preferenceResult, details)
	mu.Unlock()

}

func (c *Classifier) evaluateLanguageSignal(results *SignalResults, mu *sync.Mutex, text string) {
	start := time.Now()
	languageResult, err := c.languageClassifier.Classify(text)
	elapsed := time.Since(start)

	// Use the language code directly as the signal name
	languageCode := ""
	if err == nil && languageResult != nil {
		languageCode = languageResult.LanguageCode
	}

	// Record metrics (use microseconds for better precision)
	results.Metrics.Language.ExecutionTimeMs = float64(elapsed.Microseconds()) / 1000.0
	if languageCode != "" && err == nil && languageResult != nil {
		results.Metrics.Language.Confidence = languageResult.Confidence
	}

	if err != nil {
	} else if languageResult != nil {
		// Check if this language code is defined in language_rules
		for _, rule := range c.Config.LanguageRules {
			if rule.Name == languageCode {
				mu.Lock()
				results.MatchedLanguageRules = append(results.MatchedLanguageRules, rule.Name)
				mu.Unlock()
				break
			}
		}
	}
}

func (c *Classifier) evaluateContextSignal(results *SignalResults, mu *sync.Mutex, contextText string) {
	start := time.Now()
	matchedRules, count, err := c.contextClassifier.Classify(contextText)
	elapsed := time.Since(start)

	// Record metrics (use microseconds for better precision)
	results.Metrics.Context.ExecutionTimeMs = float64(elapsed.Microseconds()) / 1000.0
	results.Metrics.Context.Confidence = 1.0 // Rule-based, always 1.0

	if err != nil {
	} else {
		mu.Lock()
		results.MatchedContextRules = matchedRules
		results.TokenCount = count
		mu.Unlock()
	}
}

func (c *Classifier) evaluateComplexitySignal(results *SignalResults, mu *sync.Mutex, text string, imageURL string) {
	start := time.Now()
	classifyResults, err := c.complexityClassifier.ClassifyDetailedWithImage(text, imageURL)
	elapsed := time.Since(start)

	if err != nil {
		return
	}

	bestConfidence := 0.0
	mu.Lock()
	for _, result := range classifyResults {
		matchName := fmt.Sprintf("%s:%s", result.RuleName, result.Difficulty)
		results.MatchedComplexityRules = append(results.MatchedComplexityRules, matchName)
		results.SignalConfidences["complexity:"+matchName] = result.Confidence
		results.SignalValues["complexity:"+result.RuleName+":text_hard_score"] = result.TextHardScore
		results.SignalValues["complexity:"+result.RuleName+":text_easy_score"] = result.TextEasyScore
		results.SignalValues["complexity:"+result.RuleName+":text_margin"] = result.TextMargin
		results.SignalValues["complexity:"+result.RuleName+":image_hard_score"] = result.ImageHardScore
		results.SignalValues["complexity:"+result.RuleName+":image_easy_score"] = result.ImageEasyScore
		results.SignalValues["complexity:"+result.RuleName+":image_margin"] = result.ImageMargin
		results.SignalValues["complexity:"+result.RuleName+":margin"] = result.FusedMargin
		if result.Confidence > bestConfidence {
			bestConfidence = result.Confidence
		}
	}
	results.Metrics.Complexity.Confidence = bestConfidence
	mu.Unlock()
}

func (c *Classifier) evaluateModalitySignal(results *SignalResults, mu *sync.Mutex, text string) {
	start := time.Now()
	modalityResult := c.classifyModality(text, &c.Config.ModalityDetector.ModalityDetectionConfig)
	elapsed := time.Since(start)

	signalName := modalityResult.Modality

	// Record metrics
	results.Metrics.Modality.ExecutionTimeMs = float64(elapsed.Microseconds()) / 1000.0
	results.Metrics.Modality.Confidence = float64(modalityResult.Confidence)

	// Check if this signal name is defined in modality_rules
	for _, rule := range c.Config.ModalityRules {
		if strings.EqualFold(rule.Name, signalName) {
			mu.Lock()
			results.MatchedModalityRules = append(results.MatchedModalityRules, rule.Name)
			mu.Unlock()
			break
		}
	}
}

// collectJailbreakClassifierContents returns the deduplicated set of text pieces
// that need BERT classifier inference (contrastive rules are excluded).
func (c *Classifier) collectJailbreakClassifierContents(jailbreakText string, nonUserMessages []string) []string {
	seen := make(map[string]struct{})
	var contents []string
	addUnique := func(s string) {
		if s == "" {
			return
		}
		if _, ok := seen[s]; !ok {
			seen[s] = struct{}{}
			contents = append(contents, s)
		}
	}
	for _, rule := range c.Config.JailbreakRules {
		if rule.Method == "contrastive" {
			continue
		}
		addUnique(jailbreakText)
		if !rule.IncludeHistory {
			continue
		}
		for _, msg := range nonUserMessages {
			addUnique(msg)
		}
	}
	return contents
}

func (c *Classifier) evaluateJailbreakSignal(results *SignalResults, mu *sync.Mutex, jailbreakText string, nonUserMessages []string) {
	start := time.Now()

	// Step 1: Collect unique content pieces needed by classifier (non-contrastive) rules.
	classifierContents := c.collectJailbreakClassifierContents(jailbreakText, nonUserMessages)

	// Step 2: Run classifier inference exactly once per unique content piece.
	jailbreakCache := make(map[string]cachedJailbreakResult, len(classifierContents))
	for _, content := range classifierContents {
		result, err := c.jailbreakInference.Classify(content)
		jailbreakCache[content] = cachedJailbreakResult{result, err}
	}

	// Step 3: Evaluate all rules concurrently.
	var ruleWg sync.WaitGroup
	for _, rule := range c.Config.JailbreakRules {
		ruleWg.Add(1)
		go func() {
			defer ruleWg.Done()
			c.evaluateJailbreakRule(rule, jailbreakText, nonUserMessages, jailbreakCache, start, results, mu)
		}()
	}
	ruleWg.Wait()

	elapsed := time.Since(start)
	results.Metrics.Jailbreak.ExecutionTimeMs = float64(elapsed.Microseconds()) / 1000.0
	if results.JailbreakConfidence > 0 {
		results.Metrics.Jailbreak.Confidence = float64(results.JailbreakConfidence)
	}
}

func (c *Classifier) evaluateJailbreakRule(rule config.JailbreakRule, jailbreakText string, nonUserMessages []string, jailbreakCache map[string]cachedJailbreakResult, start time.Time, results *SignalResults, mu *sync.Mutex) {
	contentToAnalyze := buildContentList(jailbreakText, nonUserMessages, rule.IncludeHistory)
	if len(contentToAnalyze) == 0 {
		return
	}

	switch rule.Method {
	case "contrastive":
		c.evaluateContrastiveJailbreakRule(rule, contentToAnalyze, start, results, mu)
	default:
		c.evaluateBERTJailbreakRule(rule, contentToAnalyze, jailbreakCache, start, results, mu)
	}
}

// buildContentList assembles the text pieces to analyze for a single rule.
func buildContentList(text string, nonUserMessages []string, includeHistory bool) []string {
	var content []string
	if text != "" {
		content = append(content, text)
	}
	if includeHistory && len(nonUserMessages) > 0 {
		content = append(content, nonUserMessages...)
	}
	return content
}

func (c *Classifier) evaluateContrastiveJailbreakRule(rule config.JailbreakRule, contentToAnalyze []string, start time.Time, results *SignalResults, mu *sync.Mutex) {
	cjc, ok := c.contrastiveJailbreakClassifiers[rule.Name]
	if !ok {
		return
	}
	analysisResult := cjc.AnalyzeMessages(contentToAnalyze)
	threshold := rule.Threshold
	if threshold <= 0 {
		threshold = 0.10
	}
	if analysisResult.MaxScore < threshold {
		return
	}

	confidence := analysisResult.MaxScore
	mu.Lock()
	results.MatchedJailbreakRules = append(results.MatchedJailbreakRules, rule.Name)
	if confidence > results.JailbreakConfidence {
		results.JailbreakDetected = true
		results.JailbreakType = "contrastive"
		results.JailbreakConfidence = confidence
	}
	results.SignalConfidences["jailbreak:"+rule.Name] = float64(confidence)
	mu.Unlock()
}

func (c *Classifier) evaluateBERTJailbreakRule(rule config.JailbreakRule, contentToAnalyze []string, jailbreakCache map[string]cachedJailbreakResult, start time.Time, results *SignalResults, mu *sync.Mutex) {
	bestType, bestConf := c.findBestJailbreakMatch(rule, contentToAnalyze, jailbreakCache)
	if bestConf <= 0 {
		return
	}

	mu.Lock()
	results.MatchedJailbreakRules = append(results.MatchedJailbreakRules, rule.Name)
	if bestConf > results.JailbreakConfidence {
		results.JailbreakDetected = true
		results.JailbreakType = bestType
		results.JailbreakConfidence = bestConf
	}
	results.SignalConfidences["jailbreak:"+rule.Name] = float64(bestConf)
	mu.Unlock()
}

// findBestJailbreakMatch scans cached BERT results and returns the highest-confidence jailbreak match.
func (c *Classifier) findBestJailbreakMatch(rule config.JailbreakRule, contentToAnalyze []string, jailbreakCache map[string]cachedJailbreakResult) (string, float32) {
	var bestType string
	var bestConf float32
	for _, content := range contentToAnalyze {
		if content == "" {
			continue
		}
		cached, ok := jailbreakCache[content]
		if !ok {
			continue
		}
		if cached.err != nil {
			continue
		}
		jailbreakType, ok := c.JailbreakMapping.GetJailbreakTypeFromIndex(cached.result.Class)
		if !ok {
			continue
		}
		if cached.result.Confidence < rule.Threshold || jailbreakType != "jailbreak" {
			continue
		}
		if cached.result.Confidence > bestConf {
			bestConf = cached.result.Confidence
			bestType = jailbreakType
		}
	}
	return bestType, bestConf
}

func (c *Classifier) evaluatePIISignal(results *SignalResults, mu *sync.Mutex, piiText string, nonUserMessages []string) {
	start := time.Now()

	// Step 1: Collect the union of unique content pieces across all PII rules.
	contentSeen := make(map[string]struct{})
	var uniqueContents []string
	if piiText != "" {
		contentSeen[piiText] = struct{}{}
		uniqueContents = append(uniqueContents, piiText)
	}
	for _, rule := range c.Config.PIIRules {
		if !rule.IncludeHistory {
			continue
		}
		for _, msg := range nonUserMessages {
			if msg == "" {
				continue
			}
			if _, ok := contentSeen[msg]; !ok {
				contentSeen[msg] = struct{}{}
				uniqueContents = append(uniqueContents, msg)
			}
		}
	}

	// Step 2: Run PII token classification exactly once per unique content piece.
	// Entity types are returned as "LABEL_{class_id}" and translated by PIIMapping.
	piiCache := make(map[string]cachedPIIResult, len(uniqueContents))
	for _, content := range uniqueContents {
		tokenResult, err := c.piiInference.ClassifyTokens(content)
		piiCache[content] = cachedPIIResult{tokenResult, err}
	}

	// Step 3: Evaluate each rule concurrently using the cached token results.
	// Each goroutine applies its own threshold and allow-list without re-running the model.
	var ruleWg sync.WaitGroup
	for _, rule := range c.Config.PIIRules {
		ruleWg.Add(1)
		go func() {
			defer ruleWg.Done()
			c.evaluatePIIRule(rule, piiText, nonUserMessages, piiCache, start, results, mu)
		}()
	}
	ruleWg.Wait()

	elapsed := time.Since(start)
	results.Metrics.PII.ExecutionTimeMs = float64(elapsed.Microseconds()) / 1000.0
	if results.PIIDetected {
		results.Metrics.PII.Confidence = 1.0 // Binary: PII found or not
	}
}

func (c *Classifier) evaluatePIIRule(rule config.PIIRule, piiText string, nonUserMessages []string, piiCache map[string]cachedPIIResult, start time.Time, results *SignalResults, mu *sync.Mutex) {
	ruleContents := collectPIIRuleContents(piiText, nonUserMessages, rule.IncludeHistory)
	if len(ruleContents) == 0 {
		return
	}

	entityTypes := c.collectPIIEntityTypes(ruleContents, rule.Name, rule.Threshold, piiCache)
	deniedEntities := findDeniedEntities(entityTypes, rule.PIITypesAllowed)

	if len(deniedEntities) > 0 {
		mu.Lock()
		results.MatchedPIIRules = append(results.MatchedPIIRules, rule.Name)
		results.PIIDetected = true
		for _, e := range deniedEntities {
			if !slices.Contains(results.PIIEntities, e) {
				results.PIIEntities = append(results.PIIEntities, e)
			}
		}
		mu.Unlock()
	}
}
