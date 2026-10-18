package frequency

import (
	"cryptoapi/internal/domain/dto/cryptoanalysis"
	"cryptoapi/pkg/alphabets"
	"math"
	"strings"
)

type Service struct{}

func NewService() *Service {
	return &Service{}
}

func (s *Service) Analyze(text, lang string) (cryptoanalysis.AnalysisResult, error) {
	var reference map[rune]float64
	var alphabet []rune

	if strings.ToLower(lang) == "ru" {
		reference = alphabets.RussianReference
		alphabet = alphabets.Dictionary["rusUp"]
	} else {
		reference = alphabets.EnglishReference
		alphabet = alphabets.Dictionary["enUp"]
	}

	hist, total := calculateHistogram(text, alphabet)

	if total == 0 {
		return cryptoanalysis.AnalysisResult{
			ProbableKey:   0,
			TextHistogram: hist,
		}, nil
	}

	bestShift := 0
	minScore := math.MaxFloat64
	alphabetLen := len(alphabet)

	for shift := 0; shift < alphabetLen; shift++ {
		score := 0.0
		for i, char := range alphabet {
			shiftedIdx := (i + shift) % alphabetLen
			cipherChar := string(alphabet[shiftedIdx])

			var actualFreq float64
			if stat, ok := hist[cipherChar]; ok {
				actualFreq = stat.Percentage / 100.0
			}

			expectedFreq := reference[char]
			score += math.Pow(actualFreq-expectedFreq, 2)
		}

		if score < minScore {
			minScore = score
			bestShift = shift
		}
	}
	return cryptoanalysis.AnalysisResult{
		ProbableKey:   bestShift,
		TextHistogram: hist,
	}, nil
}

func calculateHistogram(text string, alphabet []rune) (cryptoanalysis.Histogram, int) {
	hist := make(cryptoanalysis.Histogram)
	total := 0

	alphaMap := make(map[rune]bool)
	for _, r := range alphabet {
		alphaMap[r] = true
	}
	for _, r := range strings.ToUpper(text) {
		if alphaMap[r] {
			key := string(r)
			if _, ok := hist[key]; !ok {
				hist[key] = &cryptoanalysis.CharStat{}
			}
			hist[key].Count++
			total++
		}
	}
	for _, stat := range hist {
		if total > 0 {
			stat.Percentage = (float64(stat.Count) / float64(total)) * 100.0
		}
	}
	return hist, total
}
