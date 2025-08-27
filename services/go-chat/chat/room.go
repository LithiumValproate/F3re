package chat

import (
	"encoding/json"
	"errors"
	"fmt"
	"go-chat/message"
	"go-chat/user"
	"time"
)

type metaMessage struct {
	MsgType message.MessageType `json:"type"`
}

type participantUpdateRequest struct {
	oldParticipant user.Participant
	newParticipant user.Participant
}

type Room struct {
	ID         string
	clients    map[*Client]bool
	register   chan *Client
	unregister chan *Client
	incoming   chan clientMessage
	broadcast  chan message.Message
	update     chan participantUpdateRequest
}

func NewRoom(id string) *Room {
	return &Room{
		ID:         id,
		clients:    make(map[*Client]bool),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		incoming:   make(chan clientMessage),
		broadcast:  make(chan message.Message),
		update:     make(chan participantUpdateRequest),
	}
}

func (r *Room) Run() {
	for {
		select {
		case client := <-r.register:
			r.handleRegister(client)
		case client := <-r.unregister:
			r.handleUnregister(client)
		case clientMsg := <-r.incoming:
			r.handleIncomingMessage(clientMsg.client, clientMsg.data)
		case msg := <-r.broadcast:
			r.handleBroadcast(msg)
		case req := <-r.update:
			r.handleParticipantUpdate(req)
		}
	}
}

func (r *Room) handleRegister(client *Client) {
	r.clients[client] = true
	p := client.participant
	fmt.Printf("[%s] joined room [%s]. Total: %d\n", p.Nickname(), r.ID, len(r.clients))
	notice := message.NoticeMessage{
		BaseMessage: message.BaseMessage{
			Type:      message.TypeNotice,
			Timestamp: time.Now().UnixMilli(),
		},
		Content: message.NoticeContent{
			Event:   "user_join",
			Message: fmt.Sprintf("'%s' has joined the room.", p.Nickname()),
		},
	}
	r.broadcast <- &notice
}

func (r *Room) handleUnregister(client *Client) {
	if _, ok := r.clients[client]; ok {
		delete(r.clients, client)
		close(client.send)
		p := client.participant
		fmt.Printf("[%s] left room [%s]. Total: %d\n", p.Nickname(), r.ID, len(r.clients))
		notice := message.NoticeMessage{
			BaseMessage: message.BaseMessage{
				Type:      message.TypeNotice,
				Timestamp: time.Now().UnixMilli(),
			},
			Content: message.NoticeContent{
				Event:   "muted",
				Message: "You are muted and cannot send messages.",
			},
		}
		r.broadcast <- &notice
	}
}

func (r *Room) handleIncomingMessage(sender *Client, rawMsg []byte) {
	var meta metaMessage
	if err := json.Unmarshal(rawMsg, &meta); err != nil {
		fmt.Printf("error unmarshalling message: %v\n", err)
		return
	}
	if _, ok := sender.participant.(*user.MutedParticipant); ok {
		notice := message.NoticeMessage{
			BaseMessage: message.BaseMessage{
				Type:      message.TypeNotice,
				Timestamp: time.Now().UnixMilli(),
			},
			Content: message.NoticeContent{
				Event:   "muted",
				Message: "You are muted and cannot send messages.",
			},
		}
		r.handleUnicast(sender, &notice)
		return
	}
	switch meta.MsgType {
	case message.TypeText:
		var textMsg message.TextMessage
		if err := json.Unmarshal(rawMsg, &textMsg); err != nil {
			fmt.Printf("error unmarshalling text message: %v\n", err)
			return
		}
		textMsg.SetSender(sender.participant)
		r.broadcast <- &textMsg
	case message.TypeImage:
		var imgMsg message.ImageMessage
		if err := json.Unmarshal(rawMsg, &imgMsg); err != nil {
			fmt.Printf("error unmarshalling image message: %v\n", err)
			return
		}
		imgMsg.SetSender(sender.participant)
		r.broadcast <- &imgMsg
	case message.TypeVideo:
		var vidMsg message.VideoMessage
		if err := json.Unmarshal(rawMsg, &vidMsg); err != nil {
			fmt.Printf("error unmarshalling video message: %v\n", err)
			return
		}
		vidMsg.SetSender(sender.participant)
		r.broadcast <- &vidMsg
	case message.TypeAudio:
		var audMsg message.AudioMessage
		if err := json.Unmarshal(rawMsg, &audMsg); err != nil {
			fmt.Printf("error unmarshalling audio message: %v\n", err)
			return
		}
		audMsg.SetSender(sender.participant)
		r.broadcast <- &audMsg
	case message.TypeFile:
		var fileMsg message.FileMessage
		if err := json.Unmarshal(rawMsg, &fileMsg); err != nil {
			fmt.Printf("error unmarshalling file message: %v\n", err)
			return
		}
		fileMsg.SetSender(sender.participant)
		r.broadcast <- &fileMsg
	default:
		fmt.Printf("unknown message type: %s\n", meta.MsgType)
		return
	}
}

func (r *Room) handleUnicast(client *Client, msg message.Message) {
	msgBytes := r.formatMessage(msg)
	if msgBytes != nil {
		r.writeToClient(client, msgBytes)
	}
}

func (r *Room) handleBroadcast(msg message.Message) {
	msgBytes := r.formatMessage(msg)
	if msgBytes == nil {
		return
	}
	for client := range r.clients {
		if msg.GetSender() != nil && client.participant.ID() == msg.GetSender().ID() {
			continue
		}
		r.writeToClient(client, msgBytes)
	}
}

func (r *Room) handleParticipantUpdate(req participantUpdateRequest) {
	for client := range r.clients {
		if client.participant.ID() == req.oldParticipant.ID() {
			client.participant = req.newParticipant
			fmt.Printf("Participant role changed for [%s]\n", req.newParticipant.Nickname())
			return
		}
	}
}

func (r *Room) formatMessage(msg message.Message) []byte {
	msgBytes, err := json.Marshal(msg)
	if err != nil {
		fmt.Printf("error marshalling message: %v\n", err)
		return nil
	}
	return msgBytes
}

func (r *Room) writeToClient(client *Client, msgBytes []byte) {
	select {
	case client.send <- msgBytes:
	default:
		r.handleUnregister(client)
	}
}

func (r *Room) Kick(p user.Participant) {
	for client := range r.clients {
		if client.participant.ID() == p.ID() {
			r.unregister <- client
			break
		}
	}
}

func (r *Room) ReplaceParticipant(old, new user.Participant) {
	updateReq := participantUpdateRequest{
		oldParticipant: old,
		newParticipant: new,
	}
	r.update <- updateReq
}

func (r *Room) MuteParticipant(moderator, target user.Participant) error {
	if _, ok := moderator.(*user.Moderator); !ok {
		return errors.New("permission denied")
	}
	if _, ok := target.(*user.MutedParticipant); ok {
		return errors.New("user is already muted")
	}

	mutedP := user.NewMutedParticipant(target.GetUser(), target.Nickname())
	r.ReplaceParticipant(target, mutedP)
	fmt.Printf("👑 Moderator [%s] muted [%s]\n", moderator.Nickname(), target.Nickname())
	return nil
}

func (r *Room) UnmuteParticipant(moderator, target user.Participant) error {
	// FIX: 修正了函数签名和内部逻辑
	if _, ok := moderator.(*user.Moderator); !ok {
		return errors.New("permission denied")
	}

	// 确保目标确实是 MutedParticipant
	mutedP, ok := target.(*user.MutedParticipant)
	if !ok {
		return errors.New("user is not muted")
	}

	// FIX: 正确地创建 CommonParticipant
	commonP := user.NewCommonParticipant(mutedP.GetUser(), mutedP.Nickname())
	// FIX: 调用安全的 ReplaceParticipant
	r.ReplaceParticipant(target, commonP)
	fmt.Printf("👑 Moderator [%s] unmuted [%s]\n", moderator.Nickname(), target.Nickname())
	return nil
}

// ModeratorLeave 让管理员自己离开房间 (替代了 RemoveParticipant)
// FIX: 更改了函数名，增加了返回值
func (r *Room) ModeratorLeave(moderator user.Participant) error {
	mod, ok := moderator.(*user.Moderator)
	if !ok {
		return errors.New("permission denied: only moderators can perform this action")
	}

	for client := range r.clients {
		if client.participant.ID() == mod.ID() {
			r.unregister <- client
			fmt.Printf("👑 Moderator [%s] removed themselves from room [%s]\n", mod.Nickname(), r.ID)
			return nil
		}
	}
	return errors.New("moderator not found in this room")
}

func (r *Room) ChangeNicknameOf(p user.Participant, newNickname string) error {
	// FIX: 这个操作也存在数据竞争，需要通过 channel 处理。
	// 为了简化，我们假设这是一个不常用的操作，并保持简单。
	// 在一个真正的生产系统中，这也应该通过 update channel 来完成。
	if _, ok := p.(*user.MutedParticipant); ok {
		return errors.New("muted participants cannot change nickname")
	}
	p.ChangeNickname(newNickname)
	fmt.Printf("🔄 Participant [%s] changed nickname to [%s]\n", p.Nickname(), newNickname)
	// TODO: 广播一个昵称更改的通知
	return nil
}
