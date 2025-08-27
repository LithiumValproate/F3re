package message

import (
	"go-chat/participant"
)

type MessageType string

const (
	TypeText   MessageType = "text"
	TypeImage  MessageType = "image"
	TypeVideo  MessageType = "video"
	TypeAudio  MessageType = "audio"
	TypeFile   MessageType = "file"
	TypeNotice MessageType = "notice"
)

type Message interface {
	GetSender() participant.Participant
	SetSender(sender participant.Participant)
	GetType() MessageType
	GetTimestamp() int64
}

type BaseMessage struct {
	Sender    participant.Participant `json:"sender"`
	Type      MessageType             `json:"type"`
	Timestamp int64                   `json:"timestamp"`
}

func (m *BaseMessage) GetSender() participant.Participant {
	return m.Sender
}

func (m *BaseMessage) SetSender(s participant.Participant) {
	m.Sender = s
}

func (m *BaseMessage) GetType() MessageType {
	return m.Type
}

func (m *BaseMessage) GetTimestamp() int64 {
	return m.Timestamp
}

type TextContent struct {
	Text string `json:"text"`
}

type TextMessage struct {
	BaseMessage
	Content TextContent `json:"content"`
}

type Image struct {
	URL          string `json:"url"`
	ThumbnailURL string `json:"thumbnail_url"`
	Width        int    `json:"width"`
	Height       int    `json:"height"`
	ImageSize    int64  `json:"image_size"`
	Format       string `json:"format"`
	FileName     string `json:"file_name"`
}

type ImageMessage struct {
	BaseMessage
	Content Image `json:"content"`
}

type Video struct {
	URL          string  `json:"url"`
	ThumbnailURL string  `json:"thumbnail_url"`
	Width        int     `json:"width"`
	Height       int     `json:"height"`
	Duration     float64 `json:"duration"`
	VideoSize    int64   `json:"video_size"`
	Format       string  `json:"format"`
	FileName     string  `json:"file_name"`
}

type VideoMessage struct {
	BaseMessage
	Content Video `json:"content"`
}

type Audio struct {
	URL       string  `json:"url"`
	Duration  float64 `json:"duration"`
	AudioSize int64   `json:"audio_size"`
	Format    string  `json:"format"`
	FileName  string  `json:"file_name"`
}

type AudioMessage struct {
	BaseMessage
	Content Audio `json:"content"`
}

type File struct {
	URL      string `json:"url"`
	FileSize int64  `json:"file_size"`
	Format   string `json:"format"`
	FileName string `json:"file_name"`
}

type FileMessage struct {
	BaseMessage
	Content File `json:"content"`
}

type NoticeContent struct {
	Event   string `json:"event"`
	Message string `json:"message"`
}

type NoticeMessage struct {
	BaseMessage
	Content NoticeContent `json:"content"`
}
