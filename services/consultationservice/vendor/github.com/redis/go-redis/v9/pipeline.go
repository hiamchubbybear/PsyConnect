package redis

import (
	"context"
	"errors"
)

type pipelineExecer func(context.Context, []Cmder) error














type Pipeliner interface {
	StatefulCmdable

	
	Len() int

	
	
	Do(ctx context.Context, args ...interface{}) *Cmd

	
	Process(ctx context.Context, cmd Cmder) error

	
	Discard()

	
	Exec(ctx context.Context) ([]Cmder, error)
}

var _ Pipeliner = (*Pipeline)(nil)




type Pipeline struct {
	cmdable
	statefulCmdable

	exec pipelineExecer
	cmds []Cmder
}

func (c *Pipeline) init() {
	c.cmdable = c.Process
	c.statefulCmdable = c.Process
}


func (c *Pipeline) Len() int {
	return len(c.cmds)
}


func (c *Pipeline) Do(ctx context.Context, args ...interface{}) *Cmd {
	cmd := NewCmd(ctx, args...)
	if len(args) == 0 {
		cmd.SetErr(errors.New("redis: please enter the command to be executed"))
		return cmd
	}
	_ = c.Process(ctx, cmd)
	return cmd
}


func (c *Pipeline) Process(ctx context.Context, cmd Cmder) error {
	c.cmds = append(c.cmds, cmd)
	return nil
}


func (c *Pipeline) Discard() {
	c.cmds = c.cmds[:0]
}






func (c *Pipeline) Exec(ctx context.Context) ([]Cmder, error) {
	if len(c.cmds) == 0 {
		return nil, nil
	}

	cmds := c.cmds
	c.cmds = nil

	return cmds, c.exec(ctx, cmds)
}

func (c *Pipeline) Pipelined(ctx context.Context, fn func(Pipeliner) error) ([]Cmder, error) {
	if err := fn(c); err != nil {
		return nil, err
	}
	return c.Exec(ctx)
}

func (c *Pipeline) Pipeline() Pipeliner {
	return c
}

func (c *Pipeline) TxPipelined(ctx context.Context, fn func(Pipeliner) error) ([]Cmder, error) {
	return c.Pipelined(ctx, fn)
}

func (c *Pipeline) TxPipeline() Pipeliner {
	return c
}
