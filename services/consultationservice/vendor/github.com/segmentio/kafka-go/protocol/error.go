package protocol

import (
	"fmt"
)


type Error string

func (e Error) Error() string { return string(e) }

func Errorf(msg string, args ...interface{}) Error {
	return Error(fmt.Sprintf(msg, args...))
}

const (
	
	ErrNoTopic Error = "topic not found"

	
	
	ErrNoPartition Error = "topic partition not found"

	
	
	
	ErrNoLeader Error = "topic partition has no leader"

	
	
	
	
	
	
	ErrNoRecord Error = "record set contains no records"

	
	
	ErrNoReset Error = "record sequence does not support reset"
)

type TopicError struct {
	Topic string
	Err   error
}

func NewTopicError(topic string, err error) *TopicError {
	return &TopicError{Topic: topic, Err: err}
}

func NewErrNoTopic(topic string) *TopicError {
	return NewTopicError(topic, ErrNoTopic)
}

func (e *TopicError) Error() string {
	return fmt.Sprintf("%v (topic=%q)", e.Err, e.Topic)
}

func (e *TopicError) Unwrap() error {
	return e.Err
}

type TopicPartitionError struct {
	Topic     string
	Partition int32
	Err       error
}

func NewTopicPartitionError(topic string, partition int32, err error) *TopicPartitionError {
	return &TopicPartitionError{
		Topic:     topic,
		Partition: partition,
		Err:       err,
	}
}

func NewErrNoPartition(topic string, partition int32) *TopicPartitionError {
	return NewTopicPartitionError(topic, partition, ErrNoPartition)
}

func NewErrNoLeader(topic string, partition int32) *TopicPartitionError {
	return NewTopicPartitionError(topic, partition, ErrNoLeader)
}

func (e *TopicPartitionError) Error() string {
	return fmt.Sprintf("%v (topic=%q partition=%d)", e.Err, e.Topic, e.Partition)
}

func (e *TopicPartitionError) Unwrap() error {
	return e.Err
}
