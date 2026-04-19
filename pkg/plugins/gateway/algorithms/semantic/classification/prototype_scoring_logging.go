package classification

import (
	"fmt"
	"strings"
)

const prototypeMedoidPreviewLimit = 96

func logPrototypeBankSummary(family string, owner string, bank *prototypeBank) {
}

func truncatePrototypePreview(text string) string {
	text = strings.TrimSpace(text)
	if len(text) <= prototypeMedoidPreviewLimit {
		return text
	}
	return text[:prototypeMedoidPreviewLimit-3] + "..."
}
