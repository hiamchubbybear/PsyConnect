package redis

import (
	"context"

	"github.com/redis/go-redis/v9/internal/pool"
	"github.com/redis/go-redis/v9/internal/proto"
)


const TxFailedErr = proto.RedisError("redis: transaction failed")






type Tx struct {
	baseClient
	cmdable
	statefulCmdable
}

func (c *Client) newTx() *Tx {
	tx := Tx{
		baseClient: baseClient{
			opt:        c.opt,
			connPool:   pool.NewStickyConnPool(c.connPool),
			hooksMixin: c.hooksMixin.clone(),
		},
	}
	tx.init()
	return &tx
}

func (c *Tx) init() {
	c.cmdable = c.Process
	c.statefulCmdable = c.Process

	c.initHooks(hooks{
		dial:       c.baseClient.dial,
		process:    c.baseClient.process,
		pipeline:   c.baseClient.processPipeline,
		txPipeline: c.baseClient.processTxPipeline,
	})
}

func (c *Tx) Process(ctx context.Context, cmd Cmder) error {
	err := c.processHook(ctx, cmd)
	cmd.SetErr(err)
	return err
}





func (c *Client) Watch(ctx context.Context, fn func(*Tx) error, keys ...string) error {
	tx := c.newTx()
	defer tx.Close(ctx)
	if len(keys) > 0 {
		if err := tx.Watch(ctx, keys...).Err(); err != nil {
			return err
		}
	}
	return fn(tx)
}


func (c *Tx) Close(ctx context.Context) error {
	_ = c.Unwatch(ctx).Err()
	return c.baseClient.Close()
}



func (c *Tx) Watch(ctx context.Context, keys ...string) *StatusCmd {
	args := make([]interface{}, 1+len(keys))
	args[0] = "watch"
	for i, key := range keys {
		args[1+i] = key
	}
	cmd := NewStatusCmd(ctx, args...)
	_ = c.Process(ctx, cmd)
	return cmd
}


func (c *Tx) Unwatch(ctx context.Context, keys ...string) *StatusCmd {
	args := make([]interface{}, 1+len(keys))
	args[0] = "unwatch"
	for i, key := range keys {
		args[1+i] = key
	}
	cmd := NewStatusCmd(ctx, args...)
	_ = c.Process(ctx, cmd)
	return cmd
}


func (c *Tx) Pipeline() Pipeliner {
	pipe := Pipeline{
		exec: func(ctx context.Context, cmds []Cmder) error {
			return c.processPipelineHook(ctx, cmds)
		},
	}
	pipe.init()
	return &pipe
}



func (c *Tx) Pipelined(ctx context.Context, fn func(Pipeliner) error) ([]Cmder, error) {
	return c.Pipeline().Pipelined(ctx, fn)
}









func (c *Tx) TxPipelined(ctx context.Context, fn func(Pipeliner) error) ([]Cmder, error) {
	return c.TxPipeline().Pipelined(ctx, fn)
}


func (c *Tx) TxPipeline() Pipeliner {
	pipe := Pipeline{
		exec: func(ctx context.Context, cmds []Cmder) error {
			cmds = wrapMultiExec(ctx, cmds)
			return c.processTxPipelineHook(ctx, cmds)
		},
	}
	pipe.init()
	return &pipe
}

func wrapMultiExec(ctx context.Context, cmds []Cmder) []Cmder {
	if len(cmds) == 0 {
		panic("not reached")
	}
	cmdsCopy := make([]Cmder, len(cmds)+2)
	cmdsCopy[0] = NewStatusCmd(ctx, "multi")
	copy(cmdsCopy[1:], cmds)
	cmdsCopy[len(cmdsCopy)-1] = NewSliceCmd(ctx, "exec")
	return cmdsCopy
}
