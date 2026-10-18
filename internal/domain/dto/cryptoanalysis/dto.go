package cryptoanalysis

type CharStat struct {
	Count      int     `json:"count"`
	Percentage float64 `json:"percentage"`
}

type Histogram map[string]*CharStat

type AnalysisResult struct {
	ProbableKey   int       `json:"probable_key"`
	TextHistogram Histogram `json:"histogram"`
}

type PolyAnalysisResponse struct {
	KeyLength  int              `json:"key_length"`
	Keys       []string         `json:"keys"`
	ICAnalysis ICAnalysisResult `json:"ic_analysis"`
	Kasiski    []KasiskiResult  `json:"kasiski"`
	Autocorr   AutocorrResult   `json:"autocorr"`
}

type ICAnalysisResult struct {
	OverallIC       float64 `json:"overall_ic"`
	SuggestedPeriod int     `json:"suggested_period"`
}

type KasiskiResult struct {
	Substring      string `json:"substring"`
	Positions      []int  `json:"positions"`
	Distances      []int  `json:"distances"`
	PossibleLength int    `json:"possible_length"`
}

type AutocorrResult struct {
	SuggestedPeriod int            `json:"suggested_period"`
	ShiftScores     map[string]int `json:"shift_scores"`
}
