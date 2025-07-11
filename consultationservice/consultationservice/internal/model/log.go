package model

import (
	"consultationservice/internal/enum"
	"time"
)

type LogModel struct {
	Service   string
	Action    string
	LogLevel  enum.LogLevel
	UserId    string
	Timestamp int64
	Metadata  map[string]string
	Message   string
}
type LogBuilder struct {
	log LogModel
}

func NewLogBuilder() *LogBuilder {
	return &LogBuilder{
		log: LogModel{
			Timestamp: time.Now().Unix(),
			Metadata:  make(map[string]string),
		},
	}
}
func (b *LogBuilder) Service(service string) *LogBuilder {
	b.log.Message = service
	return b
}
func (b *LogBuilder) Action(action string) *LogBuilder {
	b.log.Action = action
	return b
}

func (b *LogBuilder) Level(level enum.LogLevel) *LogBuilder {
	b.log.LogLevel = level
	return b
}

func (b *LogBuilder) User(userId string) *LogBuilder {
	b.log.UserId = userId
	return b
}

func (b *LogBuilder) Message(msg string) *LogBuilder {
	b.log.Message = msg
	return b
}

func (b *LogBuilder) AddMetadata(key, value string) *LogBuilder {
	b.log.Metadata[key] = value
	return b
}
func (b *LogBuilder) Build() LogModel {
	return b.log
}
