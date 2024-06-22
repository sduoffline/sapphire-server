package dto

type NewMessage struct {
	Title      string
	Content    string
	ReceiverID []uint
	Type       int
}
