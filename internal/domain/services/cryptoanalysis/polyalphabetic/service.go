package polyalphabetic

import (
	"cryptoapi/internal/domain/dto/cryptoanalysis"
	"cryptoapi/pkg/alphabets"
	"fmt"
	"math"
	"sort"
	"strings"
	"unicode"
)

type Service struct{}

func NewService() *Service {
	return &Service{}
}

type Result struct {
	KeyLen int
	Score  float64
}

func cleanText(text, lang string) []rune {
	text = strings.ToUpper(text)
	var targetAlphabet []rune
	if lang == "ru" {
		targetAlphabet = alphabets.Dictionary["rusUp"]
	} else {
		targetAlphabet = alphabets.Dictionary["enUp"]
	}

	validRunes := make(map[rune]bool)
	for _, r := range targetAlphabet {
		validRunes[r] = true
	}

	var cleaned []rune
	for _, r := range text {
		if validRunes[r] {
			cleaned = append(cleaned, r)
		}
	}
	return cleaned
}

func calculateOverallIC(runes []rune) float64 {
	n := float64(len(runes))
	if n <= 1 {
		return 0.0
	}
	counts := make(map[rune]float64)
	for _, r := range runes {
		counts[r]++
	}
	sum := 0.0
	for _, count := range counts {
		sum += count * (count - 1)
	}
	return sum / (n * (n - 1))
}

func getMaxKeyLen(textLen int) int {
	maxLen := 30
	if textLen/10 < maxLen {
		maxLen = textLen / 10
	}
	if maxLen < 2 {
		maxLen = 2
	}
	return maxLen
}

func indexOfCoincidence(runes []rune) []Result {
	if len(runes) == 0 {
		return nil
	}

	var results []Result
	maxKeyLen := getMaxKeyLen(len(runes))

	for k := 2; k <= maxKeyLen; k++ {
		totalIC := 0.0
		for col := 0; col < k; col++ {
			counts := make(map[rune]float64)
			colLen := 0.0
			for i := col; i < len(runes); i += k {
				counts[runes[i]]++
				colLen++
			}

			ic := 0.0
			if colLen > 1 {
				for _, count := range counts {
					ic += count * (count - 1)
				}
				ic /= colLen * (colLen - 1)
			}
			totalIC += ic
		}

		avgIC := totalIC / float64(k)
		results = append(results, Result{KeyLen: k, Score: avgIC})
	}

	sort.Slice(results, func(i, j int) bool {
		if math.Abs(results[i].Score-results[j].Score) < 0.002 {
			return results[i].KeyLen < results[j].KeyLen
		}
		return results[i].Score > results[j].Score
	})

	return results
}

func autocorrelationDetailed(runes []rune) (map[string]int, int) {
	n := len(runes)
	shiftScores := make(map[string]int)
	if n == 0 {
		return shiftScores, 0
	}

	maxKeyLen := getMaxKeyLen(n)

	type keyRate struct {
		k    int
		rate float64
	}
	var rates []keyRate

	for k := 2; k <= maxKeyLen; k++ {
		matches := 0
		for i := 0; i < n-k; i++ {
			if runes[i] == runes[i+k] {
				matches++
			}
		}
		shiftScores[fmt.Sprintf("%d", k)] = matches

		rate := 0.0
		if n-k > 0 {
			rate = float64(matches) / float64(n-k)
		}
		rates = append(rates, keyRate{k: k, rate: rate})
	}

	bestKeyLen := 0
	maxRate := -1.0
	for _, r := range rates {
		if r.rate > maxRate {
			maxRate = r.rate
			bestKeyLen = r.k
		}
	}

	for _, r := range rates {
		if bestKeyLen%r.k == 0 && r.k < bestKeyLen {
			if r.rate >= maxRate*0.85 {
				bestKeyLen = r.k
				break
			}
		}
	}

	return shiftScores, bestKeyLen
}

func kasiskiDetailed(runes []rune) ([]cryptoanalysis.KasiskiResult, int) {
	if len(runes) == 0 {
		return nil, 0
	}

	positions := make(map[string][]int)
	l := 3
	for i := 0; i <= len(runes)-l; i++ {
		gram := string(runes[i : i+l])
		positions[gram] = append(positions[gram], i)
	}

	var allResults []cryptoanalysis.KasiskiResult
	gcdCounts := make(map[int]int)

	for gram, posList := range positions {
		if len(posList) > 1 {
			var distances []int
			for i := 1; i < len(posList); i++ {
				dist := posList[i] - posList[i-1]
				distances = append(distances, dist)
			}

			possibleLen := distances[0]
			for _, d := range distances {
				possibleLen = gcd(possibleLen, d)
			}

			if possibleLen > 1 && possibleLen <= 30 {
				gcdCounts[possibleLen]++
			}

			allResults = append(allResults, cryptoanalysis.KasiskiResult{
				Substring:      gram,
				Positions:      posList,
				Distances:      distances,
				PossibleLength: possibleLen,
			})
		}
	}

	sort.Slice(allResults, func(i, j int) bool {
		return len(allResults[i].Positions) > len(allResults[j].Positions)
	})

	var results []cryptoanalysis.KasiskiResult
	if len(allResults) > 10 {
		results = allResults[:10]
	} else {
		results = allResults
	}

	bestKeyLen := 0
	maxCount := -1
	for k, count := range gcdCounts {
		if count > maxCount {
			maxCount = count
			bestKeyLen = k
		} else if count == maxCount && k > bestKeyLen {
			bestKeyLen = k
		}
	}

	return results, bestKeyLen
}

func gcd(a, b int) int {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

type shiftCandidate struct {
	char  rune
	score float64
}

type keyCandidate struct {
	runes []rune
	score float64
}

func GuessKeys(runes []rune, lang string, keyLen, topN int) []string {
	if len(runes) == 0 || keyLen <= 0 {
		return nil
	}

	var alphabet []rune
	var refFreq map[rune]float64

	if lang == "ru" {
		alphabet = alphabets.Dictionary["rusUp"]
		refFreq = alphabets.RussianReference
	} else {
		alphabet = alphabets.Dictionary["enUp"]
		refFreq = alphabets.EnglishReference
	}

	m := len(alphabet)
	runeToIndex := make(map[rune]int)
	for i, r := range alphabet {
		runeToIndex[r] = i
	}

	upRefFreq := make(map[rune]float64)
	for r, prob := range refFreq {
		upRefFreq[unicode.ToUpper(r)] = prob
	}

	beamWidth := 10
	candidates := []keyCandidate{{runes: make([]rune, 0, keyLen), score: 0.0}}

	for col := 0; col < keyLen; col++ {
		counts := make(map[rune]int)
		colLen := 0

		for i := col; i < len(runes); i += keyLen {
			if _, ok := runeToIndex[runes[i]]; ok {
				counts[runes[i]]++
				colLen++
			}
		}

		var colShifts []shiftCandidate
		if colLen == 0 {
			colShifts = append(colShifts, shiftCandidate{char: alphabet[0], score: 0})
		} else {
			for s := 0; s < m; s++ {
				score := 0.0
				for cRune, count := range counts {
					cIdx := runeToIndex[cRune]
					pIdx := (cIdx - s + m) % m
					pRune := alphabet[pIdx]
					observedFreq := float64(count) / float64(colLen)
					score += observedFreq * upRefFreq[pRune]
				}
				colShifts = append(colShifts, shiftCandidate{char: alphabet[s], score: score})
			}
			sort.Slice(colShifts, func(i, j int) bool {
				return colShifts[i].score > colShifts[j].score
			})
		}

		var newCandidates []keyCandidate
		for _, cand := range candidates {
			limit := 3
			if len(colShifts) < limit {
				limit = len(colShifts)
			}
			for _, shift := range colShifts[:limit] {
				newRunes := make([]rune, len(cand.runes))
				copy(newRunes, cand.runes)
				newRunes = append(newRunes, shift.char)
				newCandidates = append(newCandidates, keyCandidate{
					runes: newRunes,
					score: cand.score + shift.score,
				})
			}
		}

		sort.Slice(newCandidates, func(i, j int) bool {
			return newCandidates[i].score > newCandidates[j].score
		})
		if len(newCandidates) > beamWidth {
			newCandidates = newCandidates[:beamWidth]
		}
		candidates = newCandidates
	}

	var results []string
	limit := topN
	if len(candidates) < limit {
		limit = len(candidates)
	}
	for i := 0; i < limit; i++ {
		results = append(results, string(candidates[i].runes))
	}
	return results
}

func (s *Service) Analyze(text, lang string) cryptoanalysis.PolyAnalysisResponse {
	runes := cleanText(text, lang)

	overallIC := calculateOverallIC(runes)
	icResults := indexOfCoincidence(runes)
	suggestedICPeriod := 0
	if len(icResults) > 0 {
		suggestedICPeriod = icResults[0].KeyLen
	}

	kasiskiRes, suggestedKasiskiPeriod := kasiskiDetailed(runes)
	shiftScores, suggestedAutocorrPeriod := autocorrelationDetailed(runes)

	finalKeyLength := suggestedAutocorrPeriod
	if suggestedKasiskiPeriod > 1 && finalKeyLength%suggestedKasiskiPeriod == 0 {
		finalKeyLength = suggestedKasiskiPeriod
	}

	if finalKeyLength == 0 {
		finalKeyLength = suggestedKasiskiPeriod
	}
	if finalKeyLength == 0 {
		finalKeyLength = suggestedICPeriod
	}
	if finalKeyLength == 0 {
		finalKeyLength = 3
	}

	guessedKeys := GuessKeys(runes, lang, finalKeyLength, 3)

	return cryptoanalysis.PolyAnalysisResponse{
		KeyLength: finalKeyLength,
		Keys:      guessedKeys,
		ICAnalysis: cryptoanalysis.ICAnalysisResult{
			OverallIC:       overallIC,
			SuggestedPeriod: suggestedICPeriod,
		},
		Kasiski: kasiskiRes,
		Autocorr: cryptoanalysis.AutocorrResult{
			SuggestedPeriod: suggestedAutocorrPeriod,
			ShiftScores:     shiftScores,
		},
	}
}
