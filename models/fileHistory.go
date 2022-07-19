package models

type FileHistory struct {
	Channel  int    `json:"channel"`
	Sender   string `json:"sender"`
	Receiver string `json:"receiver"`
	NameFile string `json:"nameFile"`
}
