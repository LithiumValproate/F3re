package chat

import "go-chat/participant"

type RoomManager interface {
	Kick(p participant.Participant)
	ReplaceParticipant(old, new participant.Participant)
}
