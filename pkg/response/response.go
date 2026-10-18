package response

type Response struct {
	Text string `json:"text" binding:"required"`
	Lang string `json:"lang" binding:"required"`
	Key  string `json:"key,omitempty"`
}
