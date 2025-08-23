package api

type Connection interface {
	Send(message []byte) error
}
