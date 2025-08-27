package participant

import (
	"encoding/json"
	"go-identity/user"
)

type ParticipantType string

const (
	TypeCommon    ParticipantType = "common"
	TypeModerator ParticipantType = "moderator"
	TypeMuted     ParticipantType = "muted"
	TypeBot       ParticipantType = "bot"
)

type Participant interface {
	ID() string
	Nickname() string
	GetUser() user.User
	ChangeNickname(nickname string)
	Type() ParticipantType
}

type baseParticipant struct {
	user     user.User
	nickname string
}

func newBaseParticipant(u user.User, nick string) *baseParticipant {
	if nick == "" {
		nick = u.Name()
	}
	return &baseParticipant{u, nick}
}

func (p *baseParticipant) ID() string {
	return p.user.ID()
}

func (p *baseParticipant) Nickname() string {
	if p.nickname == "" {
		return p.user.Name()
	}
	return p.nickname
}

func (p *baseParticipant) GetUser() user.User {
	return p.user
}

func (p *baseParticipant) ChangeNickname(nick string) {
	if nick == "" {
		nick = p.user.Name()
	}
	p.nickname = nick
}

func (p *baseParticipant) marshalJSON(pType ParticipantType) ([]byte, error) {
	return json.Marshal(struct {
		ID       string          `json:"id"`
		Name     string          `json:"name"`
		Nickname string          `json:"nickname"`
		Type     ParticipantType `json:"type"`
	}{
		ID:       p.ID(),
		Name:     p.user.Name(),
		Nickname: p.Nickname(),
		Type:     pType,
	})
}

type CommonParticipant struct {
	*baseParticipant
}

func NewCommonParticipant(u user.User, nick string) *CommonParticipant {
	return &CommonParticipant{
		newBaseParticipant(u, nick),
	}
}

func (p *CommonParticipant) Type() ParticipantType {
	return TypeCommon
}

func (p *CommonParticipant) MarshalJSON() ([]byte, error) {
	return p.marshalJSON(p.Type())
}

type Moderator struct {
	*baseParticipant
}

func NewModerator(u user.User, nick string) *Moderator {
	return &Moderator{
		newBaseParticipant(u, nick),
	}
}

func (m *Moderator) Type() ParticipantType {
	return TypeModerator
}

func (m *Moderator) MarshalJSON() ([]byte, error) {
	return m.marshalJSON(m.Type())
}

type MutedParticipant struct {
	*baseParticipant
}

func NewMutedParticipant(u user.User, nick string) *MutedParticipant {
	return &MutedParticipant{
		newBaseParticipant(u, nick),
	}
}

func (mp *MutedParticipant) ChangeNickname(nick string) {
	// 禁止更改昵称
}

func (mp *MutedParticipant) Type() ParticipantType {
	return TypeMuted
}

func (mp *MutedParticipant) MarshalJSON() ([]byte, error) {
	return mp.marshalJSON(mp.Type())
}

type Bot struct {
	*baseParticipant
}

func NewBot(u user.User, nick string) *Bot {
	return &Bot{
		newBaseParticipant(u, nick),
	}
}

func (b *Bot) Type() ParticipantType {
	return TypeBot
}

func (b *Bot) MarshalJSON() ([]byte, error) {
	return b.marshalJSON(b.Type())
}
