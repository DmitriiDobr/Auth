package types

type Status int

const (
	Success Status = iota + 1
	Warning
	Error
)

type Message struct {
	UserID int    `json:"user_id"`
	Status Status `json:"status"`
	Header string `json:"header"`
	Body   string `json:"body"`
}
