package classification

import (
	"fmt"
	"regexp"
	"strings"
	"unicode"

	"github.com/vllm-project/aibrix/pkg/plugins/gateway/algorithms/semantic/config"
	nlp_binding "github.com/vllm-project/semantic-router/nlp-binding"
)

// preppedKeywordRule stores preprocessed keywords for efficient regex matching.
type preppedKeywordRule struct {
	Name              string // Name is also used as category
	Operator          string
	CaseSensitive     bool
	OriginalKeywords  []string         // For logging/returning original case
	CompiledRegexpsCS []*regexp.Regexp // Compiled regex for case-sensitive
	CompiledRegexpsCI []*regexp.Regexp // Compiled regex for case-insensitive

	// Fuzzy matching fields
	FuzzyMatch        bool     // Enable approximate matching with Levenshtein distance
	FuzzyThreshold    int      // Maximum edit distance for fuzzy matching (default: 2)
	LowercaseKeywords []string // Pre-computed lowercase for fuzzy matching
}

// KeywordClassifier implements keyword-based classification logic.
// It supports three matching methods per rule: regex (default), bm25, and ngram.
// BM25 and N-gram rules are dispatched to Rust-backed classifiers via nlp-binding.
type KeywordClassifier struct {
	regexRules []preppedKeywordRule // Regex-based rules (original behavior)

	// Rust-backed classifiers via nlp-binding FFI
	bm25Classifier  *nlp_binding.BM25Classifier
	ngramClassifier *nlp_binding.NgramClassifier

	// Track which rules use which method for ordered evaluation
	ruleOrder []ruleRef
}

// ruleRef tracks the method and index for ordered rule evaluation.
type ruleRef struct {
	method string // "regex", "bm25", "ngram"
	name   string // rule name for logging
}

type ruleMatch struct {
	matched       bool
	ruleName      string
	keywords      []string
	matchCount    int
	totalKeywords int
}

// NewKeywordClassifier creates a new KeywordClassifier.
// Rules with method "bm25" or "ngram" are dispatched to Rust-backed classifiers;
// all others (including default/empty method) use the original regex engine.
func NewKeywordClassifier(cfgRules []config.KeywordRule) (*KeywordClassifier, error) {
	return &KeywordClassifier{}, nil
}

// Free releases Rust-side resources. Call when the classifier is no longer needed.
func (c *KeywordClassifier) Free() {
}

// prepRegexRule creates a preppedKeywordRule from a config rule (original regex logic).
func prepRegexRule(rule config.KeywordRule) (preppedKeywordRule, error) {
	return preppedKeywordRule{}, nil
}

func regexPatterns(keyword string, useExplicitRegex bool) (string, string) {
	return "", ""
}

// Classify performs keyword-based classification on the given text.
// Returns category, confidence (0.0-1.0), and error.
// For regex: confidence = 0.5 + (matchCount / totalKeywords * 0.5)
// For BM25/N-gram: confidence derived from match scores
func (c *KeywordClassifier) Classify(text string) (string, float64, error) {
	return "", 0.0, nil
}

// ClassifyWithKeywords performs keyword-based classification and returns the matched keywords.
func (c *KeywordClassifier) ClassifyWithKeywords(text string) (string, []string, error) {
	return "", nil, nil
}

// ClassifyWithKeywordsAndCount performs keyword-based classification and returns:
// - category: the matched rule name (or "" if no match)
// - matchedKeywords: slice of keywords that matched
// - matchCount: number of keywords that matched
// - totalKeywords: total number of keywords in the matched rule
// - error: any error that occurred
//
// Rules are evaluated in the order they were defined in the config (first-match semantics),
// regardless of method. Each rule is dispatched to its respective engine.
func (c *KeywordClassifier) ClassifyWithKeywordsAndCount(text string) (string, []string, int, int, error) {
	return "", nil, 0, 0, nil
}

func (c *KeywordClassifier) classifyRule(text string, ref ruleRef, regexIdx *int) (ruleMatch, error) {
	return ruleMatch{}, nil
}

func (c *KeywordClassifier) classifyRegexRule(text string, regexIdx *int) (ruleMatch, error) {
	return ruleMatch{}, nil
}

// matchesWithCount checks if the text matches the given keyword rule.
func (c *KeywordClassifier) matchesWithCount(text string, rule preppedKeywordRule) (bool, []string, int, error) {
	return false, nil, 0, nil
}

func regexpsForRule(rule preppedKeywordRule) []*regexp.Regexp {
	if rule.CaseSensitive {
		return rule.CompiledRegexpsCS
	}
	return rule.CompiledRegexpsCI
}

func lowerWordsForRule(rule preppedKeywordRule, text string) []string {
	if !rule.FuzzyMatch {
		return nil
	}
	return extractLowerWords(text)
}

func validateRegexp(ruleName string, idx int, re *regexp.Regexp) error {
	if re != nil {
		return nil
	}
	return fmt.Errorf("nil regular expression found in rule %q at index %d. This indicates a failed compilation during initialization", ruleName, idx)
}

func matchAND(text string, rule preppedKeywordRule, regexpsToUse []*regexp.Regexp, lowerTextWords []string) (bool, []string, int, error) {
	matchedKeywords := make([]string, 0, len(rule.OriginalKeywords))
	for i, re := range regexpsToUse {
		if err := validateRegexp(rule.Name, i, re); err != nil {
			return false, nil, 0, err
		}
		if re.MatchString(text) {
			matchedKeywords = append(matchedKeywords, rule.OriginalKeywords[i])
			continue
		}
		if hasFuzzyMatch(rule, i, lowerTextWords) {
			matchedKeywords = append(matchedKeywords, rule.OriginalKeywords[i]+" (fuzzy)")
			continue
		}
		return false, nil, 0, nil
	}
	return true, matchedKeywords, len(matchedKeywords), nil
}

func matchOR(text string, rule preppedKeywordRule, regexpsToUse []*regexp.Regexp, lowerTextWords []string) (bool, []string, int, error) {
	matchedKeywords := make([]string, 0, len(rule.OriginalKeywords))
	matchedSet := make(map[string]bool)
	for i, re := range regexpsToUse {
		if err := validateRegexp(rule.Name, i, re); err != nil {
			return false, nil, 0, err
		}
		keyword := rule.OriginalKeywords[i]
		if re.MatchString(text) {
			addMatchedKeyword(keyword, matchedSet, &matchedKeywords)
			continue
		}
		if hasFuzzyMatch(rule, i, lowerTextWords) {
			addMatchedKeyword(keyword+" (fuzzy)", matchedSet, &matchedKeywords)
		}
	}
	if len(matchedKeywords) == 0 {
		return false, nil, 0, nil
	}
	return true, matchedKeywords, len(matchedKeywords), nil
}

func matchNOR(text string, rule preppedKeywordRule, regexpsToUse []*regexp.Regexp, lowerTextWords []string) (bool, []string, int, error) {
	for i, re := range regexpsToUse {
		if err := validateRegexp(rule.Name, i, re); err != nil {
			return false, nil, 0, err
		}
		if re.MatchString(text) || hasFuzzyMatch(rule, i, lowerTextWords) {
			return false, nil, 0, nil
		}
	}
	return true, nil, 0, nil
}

func hasFuzzyMatch(rule preppedKeywordRule, idx int, lowerTextWords []string) bool {
	if !rule.FuzzyMatch || idx >= len(rule.LowercaseKeywords) {
		return false
	}
	return fuzzyMatch(rule.LowercaseKeywords[idx], lowerTextWords, rule.FuzzyThreshold)
}

func addMatchedKeyword(keyword string, matchedSet map[string]bool, matchedKeywords *[]string) {
	if matchedSet[keyword] {
		return
	}
	matchedSet[keyword] = true
	*matchedKeywords = append(*matchedKeywords, keyword)
}

// ----------- Fuzzy Matching -----------

// levenshteinDistance calculates the edit distance between two strings.
// Uses Wagner-Fischer dynamic programming approach with O(m*n) time complexity.
func levenshteinDistance(s1, s2 string) int {
	s1 = strings.ToLower(s1)
	s2 = strings.ToLower(s2)

	r1 := []rune(s1)
	r2 := []rune(s2)
	len1 := len(r1)
	len2 := len(r2)

	if len1 == 0 {
		return len2
	}
	if len2 == 0 {
		return len1
	}

	// Optimize space to O(min(m,n))
	if len1 > len2 {
		r1, r2 = r2, r1
		len1, len2 = len2, len1
	}

	prev := make([]int, len1+1)
	curr := make([]int, len1+1)

	for i := 0; i <= len1; i++ {
		prev[i] = i
	}

	for j := 1; j <= len2; j++ {
		curr[0] = j
		for i := 1; i <= len1; i++ {
			cost := 0
			if r1[i-1] != r2[j-1] {
				cost = 1
			}
			curr[i] = min(prev[i]+1, min(curr[i-1]+1, prev[i-1]+cost))
		}
		prev, curr = curr, prev
	}

	return prev[len1]
}

// fuzzyMatch checks if any word in text fuzzy-matches the keyword within threshold.
func fuzzyMatch(lowerKeyword string, lowerTextWords []string, threshold int) bool {
	for _, textWord := range lowerTextWords {
		if levenshteinDistance(textWord, lowerKeyword) <= threshold {
			return true
		}
	}
	return false
}

// extractLowerWords splits text into lowercase words for fuzzy matching.
func extractLowerWords(text string) []string {
	var words []string
	var currentWord strings.Builder

	for _, r := range text {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			currentWord.WriteRune(r)
		} else if currentWord.Len() > 0 {
			words = append(words, strings.ToLower(currentWord.String()))
			currentWord.Reset()
		}
	}

	if currentWord.Len() > 0 {
		words = append(words, strings.ToLower(currentWord.String()))
	}

	return words
}
