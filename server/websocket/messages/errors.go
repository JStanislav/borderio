package messages

type ErrorMessage struct {
	Type    string `json:"type"`
	Message string `json:"message"`
}
