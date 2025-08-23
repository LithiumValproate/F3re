module go-chat

go 1.24.4

require (
	github.com/gorilla/websocket v1.5.3
	go-identity v0.0.0
)

require golang.org/x/crypto v0.41.0 // indirect

replace go-identity => ./../go-identity
