package validation

import (
	"strings"
	"unicode"
)

// NormalizeAnswer normalizes an answer for comparison
func NormalizeAnswer(answer string) string {
	// Convert to lowercase
	answer = strings.ToLower(answer)

	// Remove common prefixes
	prefixes := []string{"the ", "a ", "an "}
	for _, prefix := range prefixes {
		if strings.HasPrefix(answer, prefix) {
			answer = strings.TrimPrefix(answer, prefix)
		}
	}

	// Remove punctuation and extra spaces
	var result strings.Builder
	for _, r := range answer {
		if !unicode.IsPunct(r) {
			result.WriteRune(r)
		}
	}

	// Trim spaces and normalize internal spaces
	return strings.Join(strings.Fields(result.String()), " ")
}

// IsSimilarAnswer checks if two answers are similar enough to be considered the same
func IsSimilarAnswer(answer1, answer2 string) bool {
	normalized1 := NormalizeAnswer(answer1)
	normalized2 := NormalizeAnswer(answer2)

	// Exact match after normalization
	if normalized1 == normalized2 {
		return true
	}

	// Check if one is contained within the other
	if strings.Contains(normalized1, normalized2) || strings.Contains(normalized2, normalized1) {
		return true
	}

	// Calculate Levenshtein distance for close matches
	distance := levenshteinDistance(normalized1, normalized2)
	maxLen := max(len(normalized1), len(normalized2))

	// If the distance is less than 20% of the longer string length, consider it similar
	return float64(distance)/float64(maxLen) < 0.2
}

// levenshteinDistance calculates the Levenshtein distance between two strings
func levenshteinDistance(s1, s2 string) int {
	if len(s1) == 0 {
		return len(s2)
	}
	if len(s2) == 0 {
		return len(s1)
	}

	matrix := make([][]int, len(s1)+1)
	for i := range matrix {
		matrix[i] = make([]int, len(s2)+1)
		matrix[i][0] = i
	}
	for j := range matrix[0] {
		matrix[0][j] = j
	}

	for i := 1; i <= len(s1); i++ {
		for j := 1; j <= len(s2); j++ {
			if s1[i-1] == s2[j-1] {
				matrix[i][j] = matrix[i-1][j-1]
			} else {
				matrix[i][j] = min(
					matrix[i-1][j]+1,   // deletion
					matrix[i][j-1]+1,   // insertion
					matrix[i-1][j-1]+1, // substitution
				)
			}
		}
	}

	return matrix[len(s1)][len(s2)]
}

func min(a, b, c int) int {
	if a < b {
		if a < c {
			return a
		}
		return c
	}
	if b < c {
		return b
	}
	return c
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
