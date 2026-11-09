package models

// Response adalah bentuk standar body JSON yang dikirim controller.
type Response struct {
	Status   string `json:"status"`
	Message  string `json:"message"`
	Token    string `json:"token,omitempty"`
	Username string `json:"username,omitempty"`
}
