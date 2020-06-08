package model

//Message struct
type Message struct {
	Sender  string `json:"sender,omitempty"`
	Reciver string `json:"receiver,omitempty"`
	Content string `json:"content,omitempty"`
}
