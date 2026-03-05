package redis

import (
	"context"
	"errors"
	"fmt"
	"net"
	"sync"
	"sync/atomic"
	"time"

	"github.com/redis/go-redis/v9/auth"
	"github.com/redis/go-redis/v9/internal"
	"github.com/redis/go-redis/v9/internal/hscan"
	"github.com/redis/go-redis/v9/internal/pool"
	"github.com/redis/go-redis/v9/internal/proto"
)


type Scanner = hscan.Scanner


const Nil = proto.Nil


func SetLogger(logger internal.Logging) {
	internal.Logger = logger
}



type Hook interface {
	DialHook(next DialHook) DialHook
	ProcessHook(next ProcessHook) ProcessHook
	ProcessPipelineHook(next ProcessPipelineHook) ProcessPipelineHook
}

type (
	DialHook            func(ctx context.Context, network, addr string) (net.Conn, error)
	ProcessHook         func(ctx context.Context, cmd Cmder) error
	ProcessPipelineHook func(ctx context.Context, cmds []Cmder) error
)

type hooksMixin struct {
	hooksMu *sync.RWMutex

	slice   []Hook
	initial hooks
	current hooks
}

func (hs *hooksMixin) initHooks(hooks hooks) {
	hs.hooksMu = new(sync.RWMutex)
	hs.initial = hooks
	hs.chain()
}

type hooks struct {
	dial       DialHook
	process    ProcessHook
	pipeline   ProcessPipelineHook
	txPipeline ProcessPipelineHook
}

func (h *hooks) setDefaults() {
	if h.dial == nil {
		h.dial = func(ctx context.Context, network, addr string) (net.Conn, error) { return nil, nil }
	}
	if h.process == nil {
		h.process = func(ctx context.Context, cmd Cmder) error { return nil }
	}
	if h.pipeline == nil {
		h.pipeline = func(ctx context.Context, cmds []Cmder) error { return nil }
	}
	if h.txPipeline == nil {
		h.txPipeline = func(ctx context.Context, cmds []Cmder) error { return nil }
	}
}





































func (hs *hooksMixin) AddHook(hook Hook) {
	hs.slice = append(hs.slice, hook)
	hs.chain()
}

func (hs *hooksMixin) chain() {
	hs.initial.setDefaults()

	hs.hooksMu.Lock()
	defer hs.hooksMu.Unlock()

	hs.current.dial = hs.initial.dial
	hs.current.process = hs.initial.process
	hs.current.pipeline = hs.initial.pipeline
	hs.current.txPipeline = hs.initial.txPipeline

	for i := len(hs.slice) - 1; i >= 0; i-- {
		if wrapped := hs.slice[i].DialHook(hs.current.dial); wrapped != nil {
			hs.current.dial = wrapped
		}
		if wrapped := hs.slice[i].ProcessHook(hs.current.process); wrapped != nil {
			hs.current.process = wrapped
		}
		if wrapped := hs.slice[i].ProcessPipelineHook(hs.current.pipeline); wrapped != nil {
			hs.current.pipeline = wrapped
		}
		if wrapped := hs.slice[i].ProcessPipelineHook(hs.current.txPipeline); wrapped != nil {
			hs.current.txPipeline = wrapped
		}
	}
}

func (hs *hooksMixin) clone() hooksMixin {
	hs.hooksMu.Lock()
	defer hs.hooksMu.Unlock()

	clone := *hs
	l := len(clone.slice)
	clone.slice = clone.slice[:l:l]
	clone.hooksMu = new(sync.RWMutex)
	return clone
}

func (hs *hooksMixin) withProcessHook(ctx context.Context, cmd Cmder, hook ProcessHook) error {
	for i := len(hs.slice) - 1; i >= 0; i-- {
		if wrapped := hs.slice[i].ProcessHook(hook); wrapped != nil {
			hook = wrapped
		}
	}
	return hook(ctx, cmd)
}

func (hs *hooksMixin) withProcessPipelineHook(
	ctx context.Context, cmds []Cmder, hook ProcessPipelineHook,
) error {
	for i := len(hs.slice) - 1; i >= 0; i-- {
		if wrapped := hs.slice[i].ProcessPipelineHook(hook); wrapped != nil {
			hook = wrapped
		}
	}
	return hook(ctx, cmds)
}

func (hs *hooksMixin) dialHook(ctx context.Context, network, addr string) (net.Conn, error) {
	
	
	
	hs.hooksMu.RLock()
	current := hs.current
	hs.hooksMu.RUnlock()

	return current.dial(ctx, network, addr)
}

func (hs *hooksMixin) processHook(ctx context.Context, cmd Cmder) error {
	return hs.current.process(ctx, cmd)
}

func (hs *hooksMixin) processPipelineHook(ctx context.Context, cmds []Cmder) error {
	return hs.current.pipeline(ctx, cmds)
}

func (hs *hooksMixin) processTxPipelineHook(ctx context.Context, cmds []Cmder) error {
	return hs.current.txPipeline(ctx, cmds)
}



type baseClient struct {
	opt      *Options
	connPool pool.Pooler
	hooksMixin

	onClose func() error 
}

func (c *baseClient) clone() *baseClient {
	clone := *c
	return &clone
}

func (c *baseClient) withTimeout(timeout time.Duration) *baseClient {
	opt := c.opt.clone()
	opt.ReadTimeout = timeout
	opt.WriteTimeout = timeout

	clone := c.clone()
	clone.opt = opt

	return clone
}

func (c *baseClient) String() string {
	return fmt.Sprintf("Redis<%s db:%d>", c.getAddr(), c.opt.DB)
}

func (c *baseClient) newConn(ctx context.Context) (*pool.Conn, error) {
	cn, err := c.connPool.NewConn(ctx)
	if err != nil {
		return nil, err
	}

	err = c.initConn(ctx, cn)
	if err != nil {
		_ = c.connPool.CloseConn(cn)
		return nil, err
	}

	return cn, nil
}

func (c *baseClient) getConn(ctx context.Context) (*pool.Conn, error) {
	if c.opt.Limiter != nil {
		err := c.opt.Limiter.Allow()
		if err != nil {
			return nil, err
		}
	}

	cn, err := c._getConn(ctx)
	if err != nil {
		if c.opt.Limiter != nil {
			c.opt.Limiter.ReportResult(err)
		}
		return nil, err
	}

	return cn, nil
}

func (c *baseClient) _getConn(ctx context.Context) (*pool.Conn, error) {
	cn, err := c.connPool.Get(ctx)
	if err != nil {
		return nil, err
	}

	if cn.Inited {
		return cn, nil
	}

	if err := c.initConn(ctx, cn); err != nil {
		c.connPool.Remove(ctx, cn, err)
		if err := errors.Unwrap(err); err != nil {
			return nil, err
		}
		return nil, err
	}

	return cn, nil
}

func (c *baseClient) newReAuthCredentialsListener(poolCn *pool.Conn) auth.CredentialsListener {
	return auth.NewReAuthCredentialsListener(
		c.reAuthConnection(poolCn),
		c.onAuthenticationErr(poolCn),
	)
}

func (c *baseClient) reAuthConnection(poolCn *pool.Conn) func(credentials auth.Credentials) error {
	return func(credentials auth.Credentials) error {
		var err error
		username, password := credentials.BasicAuth()
		ctx := context.Background()
		connPool := pool.NewSingleConnPool(c.connPool, poolCn)
		
		cn := newConn(c.opt, connPool, nil)

		if username != "" {
			err = cn.AuthACL(ctx, username, password).Err()
		} else {
			err = cn.Auth(ctx, password).Err()
		}
		return err
	}
}
func (c *baseClient) onAuthenticationErr(poolCn *pool.Conn) func(err error) {
	return func(err error) {
		if err != nil {
			if isBadConn(err, false, c.opt.Addr) {
				
				err := c.connPool.CloseConn(poolCn)
				if err != nil {
					internal.Logger.Printf(context.Background(), "redis: failed to close connection: %v", err)
					
					
					err := poolCn.Close()
					if err != nil {
						internal.Logger.Printf(context.Background(), "redis: failed to close network connection: %v", err)
					}
				}
			}
			internal.Logger.Printf(context.Background(), "redis: re-authentication failed: %v", err)
		}
	}
}

func (c *baseClient) wrappedOnClose(newOnClose func() error) func() error {
	onClose := c.onClose
	return func() error {
		var firstErr error
		err := newOnClose()
		
		
		
		if err != nil {
			firstErr = err
		}
		if onClose != nil {
			err = onClose()
			if err != nil && firstErr == nil {
				firstErr = err
			}
		}
		return firstErr
	}
}

func (c *baseClient) initConn(ctx context.Context, cn *pool.Conn) error {
	if cn.Inited {
		return nil
	}

	var err error
	cn.Inited = true
	connPool := pool.NewSingleConnPool(c.connPool, cn)
	conn := newConn(c.opt, connPool, &c.hooksMixin)

	username, password := "", ""
	if c.opt.StreamingCredentialsProvider != nil {
		credentials, unsubscribeFromCredentialsProvider, err := c.opt.StreamingCredentialsProvider.
			Subscribe(c.newReAuthCredentialsListener(cn))
		if err != nil {
			return fmt.Errorf("failed to subscribe to streaming credentials: %w", err)
		}
		c.onClose = c.wrappedOnClose(unsubscribeFromCredentialsProvider)
		cn.SetOnClose(unsubscribeFromCredentialsProvider)
		username, password = credentials.BasicAuth()
	} else if c.opt.CredentialsProviderContext != nil {
		username, password, err = c.opt.CredentialsProviderContext(ctx)
		if err != nil {
			return fmt.Errorf("failed to get credentials from context provider: %w", err)
		}
	} else if c.opt.CredentialsProvider != nil {
		username, password = c.opt.CredentialsProvider()
	} else if c.opt.Username != "" || c.opt.Password != "" {
		username, password = c.opt.Username, c.opt.Password
	}

	
	
	if err = conn.Hello(ctx, c.opt.Protocol, username, password, c.opt.ClientName).Err(); err == nil {
		
	} else if !isRedisError(err) {
		
		
		
		
		
		
		
		return err
	} else if password != "" {
		
		if username != "" {
			err = conn.AuthACL(ctx, username, password).Err()
		} else {
			err = conn.Auth(ctx, password).Err()
		}
		if err != nil {
			return fmt.Errorf("failed to authenticate: %w", err)
		}
	}

	_, err = conn.Pipelined(ctx, func(pipe Pipeliner) error {
		if c.opt.DB > 0 {
			pipe.Select(ctx, c.opt.DB)
		}

		if c.opt.readOnly {
			pipe.ReadOnly(ctx)
		}

		if c.opt.ClientName != "" {
			pipe.ClientSetName(ctx, c.opt.ClientName)
		}

		return nil
	})
	if err != nil {
		return fmt.Errorf("failed to initialize connection options: %w", err)
	}

	if !c.opt.DisableIdentity && !c.opt.DisableIndentity {
		libName := ""
		libVer := Version()
		if c.opt.IdentitySuffix != "" {
			libName = c.opt.IdentitySuffix
		}
		p := conn.Pipeline()
		p.ClientSetInfo(ctx, WithLibraryName(libName))
		p.ClientSetInfo(ctx, WithLibraryVersion(libVer))
		
		
		if _, err = p.Exec(ctx); err != nil && !isRedisError(err) {
			return err
		}
	}

	if c.opt.OnConnect != nil {
		return c.opt.OnConnect(ctx, conn)
	}

	return nil
}

func (c *baseClient) releaseConn(ctx context.Context, cn *pool.Conn, err error) {
	if c.opt.Limiter != nil {
		c.opt.Limiter.ReportResult(err)
	}

	if isBadConn(err, false, c.opt.Addr) {
		c.connPool.Remove(ctx, cn, err)
	} else {
		c.connPool.Put(ctx, cn)
	}
}

func (c *baseClient) withConn(
	ctx context.Context, fn func(context.Context, *pool.Conn) error,
) error {
	cn, err := c.getConn(ctx)
	if err != nil {
		return err
	}

	var fnErr error
	defer func() {
		c.releaseConn(ctx, cn, fnErr)
	}()

	fnErr = fn(ctx, cn)

	return fnErr
}

func (c *baseClient) dial(ctx context.Context, network, addr string) (net.Conn, error) {
	return c.opt.Dialer(ctx, network, addr)
}

func (c *baseClient) process(ctx context.Context, cmd Cmder) error {
	var lastErr error
	for attempt := 0; attempt <= c.opt.MaxRetries; attempt++ {
		attempt := attempt

		retry, err := c._process(ctx, cmd, attempt)
		if err == nil || !retry {
			return err
		}

		lastErr = err
	}
	return lastErr
}

func (c *baseClient) assertUnstableCommand(cmd Cmder) bool {
	switch cmd.(type) {
	case *AggregateCmd, *FTInfoCmd, *FTSpellCheckCmd, *FTSearchCmd, *FTSynDumpCmd:
		if c.opt.UnstableResp3 {
			return true
		} else {
			panic("RESP3 responses for this command are disabled because they may still change. Please set the flag UnstableResp3 .  See the [README](https://github.com/redis/go-redis/blob/master/README.md) and the release notes for guidance.")
		}
	default:
		return false
	}
}

func (c *baseClient) _process(ctx context.Context, cmd Cmder, attempt int) (bool, error) {
	if attempt > 0 {
		if err := internal.Sleep(ctx, c.retryBackoff(attempt)); err != nil {
			return false, err
		}
	}

	retryTimeout := uint32(0)
	if err := c.withConn(ctx, func(ctx context.Context, cn *pool.Conn) error {
		if err := cn.WithWriter(c.context(ctx), c.opt.WriteTimeout, func(wr *proto.Writer) error {
			return writeCmd(wr, cmd)
		}); err != nil {
			atomic.StoreUint32(&retryTimeout, 1)
			return err
		}
		readReplyFunc := cmd.readReply
		
		if c.opt.Protocol != 2 && c.assertUnstableCommand(cmd) {
			readReplyFunc = cmd.readRawReply
		}
		if err := cn.WithReader(c.context(ctx), c.cmdTimeout(cmd), readReplyFunc); err != nil {
			if cmd.readTimeout() == nil {
				atomic.StoreUint32(&retryTimeout, 1)
			} else {
				atomic.StoreUint32(&retryTimeout, 0)
			}
			return err
		}

		return nil
	}); err != nil {
		retry := shouldRetry(err, atomic.LoadUint32(&retryTimeout) == 1)
		return retry, err
	}

	return false, nil
}

func (c *baseClient) retryBackoff(attempt int) time.Duration {
	return internal.RetryBackoff(attempt, c.opt.MinRetryBackoff, c.opt.MaxRetryBackoff)
}

func (c *baseClient) cmdTimeout(cmd Cmder) time.Duration {
	if timeout := cmd.readTimeout(); timeout != nil {
		t := *timeout
		if t == 0 {
			return 0
		}
		return t + 10*time.Second
	}
	return c.opt.ReadTimeout
}




func (c *baseClient) context(ctx context.Context) context.Context {
	if c.opt.ContextTimeoutEnabled {
		return ctx
	}
	return context.Background()
}





func (c *baseClient) Close() error {
	var firstErr error
	if c.onClose != nil {
		if err := c.onClose(); err != nil {
			firstErr = err
		}
	}
	if err := c.connPool.Close(); err != nil && firstErr == nil {
		firstErr = err
	}
	return firstErr
}

func (c *baseClient) getAddr() string {
	return c.opt.Addr
}

func (c *baseClient) processPipeline(ctx context.Context, cmds []Cmder) error {
	if err := c.generalProcessPipeline(ctx, cmds, c.pipelineProcessCmds); err != nil {
		return err
	}
	return cmdsFirstErr(cmds)
}

func (c *baseClient) processTxPipeline(ctx context.Context, cmds []Cmder) error {
	if err := c.generalProcessPipeline(ctx, cmds, c.txPipelineProcessCmds); err != nil {
		return err
	}
	return cmdsFirstErr(cmds)
}

type pipelineProcessor func(context.Context, *pool.Conn, []Cmder) (bool, error)

func (c *baseClient) generalProcessPipeline(
	ctx context.Context, cmds []Cmder, p pipelineProcessor,
) error {
	var lastErr error
	for attempt := 0; attempt <= c.opt.MaxRetries; attempt++ {
		if attempt > 0 {
			if err := internal.Sleep(ctx, c.retryBackoff(attempt)); err != nil {
				setCmdsErr(cmds, err)
				return err
			}
		}

		
		canRetry := true
		lastErr = c.withConn(ctx, func(ctx context.Context, cn *pool.Conn) error {
			var err error
			canRetry, err = p(ctx, cn, cmds)
			return err
		})
		if lastErr == nil || !canRetry || !shouldRetry(lastErr, true) {
			return lastErr
		}
	}
	return lastErr
}

func (c *baseClient) pipelineProcessCmds(
	ctx context.Context, cn *pool.Conn, cmds []Cmder,
) (bool, error) {
	if err := cn.WithWriter(c.context(ctx), c.opt.WriteTimeout, func(wr *proto.Writer) error {
		return writeCmds(wr, cmds)
	}); err != nil {
		setCmdsErr(cmds, err)
		return true, err
	}

	if err := cn.WithReader(c.context(ctx), c.opt.ReadTimeout, func(rd *proto.Reader) error {
		return pipelineReadCmds(rd, cmds)
	}); err != nil {
		return true, err
	}

	return false, nil
}

func pipelineReadCmds(rd *proto.Reader, cmds []Cmder) error {
	for i, cmd := range cmds {
		err := cmd.readReply(rd)
		cmd.SetErr(err)
		if err != nil && !isRedisError(err) {
			setCmdsErr(cmds[i+1:], err)
			return err
		}
	}
	
	return cmds[0].Err()
}

func (c *baseClient) txPipelineProcessCmds(
	ctx context.Context, cn *pool.Conn, cmds []Cmder,
) (bool, error) {
	if err := cn.WithWriter(c.context(ctx), c.opt.WriteTimeout, func(wr *proto.Writer) error {
		return writeCmds(wr, cmds)
	}); err != nil {
		setCmdsErr(cmds, err)
		return true, err
	}

	if err := cn.WithReader(c.context(ctx), c.opt.ReadTimeout, func(rd *proto.Reader) error {
		statusCmd := cmds[0].(*StatusCmd)
		
		trimmedCmds := cmds[1 : len(cmds)-1]

		if err := txPipelineReadQueued(rd, statusCmd, trimmedCmds); err != nil {
			setCmdsErr(cmds, err)
			return err
		}

		return pipelineReadCmds(rd, trimmedCmds)
	}); err != nil {
		return false, err
	}

	return false, nil
}

func txPipelineReadQueued(rd *proto.Reader, statusCmd *StatusCmd, cmds []Cmder) error {
	
	if err := statusCmd.readReply(rd); err != nil {
		return err
	}

	
	for range cmds {
		if err := statusCmd.readReply(rd); err != nil && !isRedisError(err) {
			return err
		}
	}

	
	line, err := rd.ReadLine()
	if err != nil {
		if err == Nil {
			err = TxFailedErr
		}
		return err
	}

	if line[0] != proto.RespArray {
		return fmt.Errorf("redis: expected '*', but got line %q", line)
	}

	return nil
}








type Client struct {
	*baseClient
	cmdable
}


func NewClient(opt *Options) *Client {
	if opt == nil {
		panic("redis: NewClient nil options")
	}
	opt.init()

	c := Client{
		baseClient: &baseClient{
			opt: opt,
		},
	}
	c.init()
	c.connPool = newConnPool(opt, c.dialHook)

	return &c
}

func (c *Client) init() {
	c.cmdable = c.Process
	c.initHooks(hooks{
		dial:       c.baseClient.dial,
		process:    c.baseClient.process,
		pipeline:   c.baseClient.processPipeline,
		txPipeline: c.baseClient.processTxPipeline,
	})
}

func (c *Client) WithTimeout(timeout time.Duration) *Client {
	clone := *c
	clone.baseClient = c.baseClient.withTimeout(timeout)
	clone.init()
	return &clone
}

func (c *Client) Conn() *Conn {
	return newConn(c.opt, pool.NewStickyConnPool(c.connPool), &c.hooksMixin)
}

func (c *Client) Process(ctx context.Context, cmd Cmder) error {
	err := c.processHook(ctx, cmd)
	cmd.SetErr(err)
	return err
}


func (c *Client) Options() *Options {
	return c.opt
}

type PoolStats pool.Stats


func (c *Client) PoolStats() *PoolStats {
	stats := c.connPool.Stats()
	return (*PoolStats)(stats)
}

func (c *Client) Pipelined(ctx context.Context, fn func(Pipeliner) error) ([]Cmder, error) {
	return c.Pipeline().Pipelined(ctx, fn)
}

func (c *Client) Pipeline() Pipeliner {
	pipe := Pipeline{
		exec: pipelineExecer(c.processPipelineHook),
	}
	pipe.init()
	return &pipe
}

func (c *Client) TxPipelined(ctx context.Context, fn func(Pipeliner) error) ([]Cmder, error) {
	return c.TxPipeline().Pipelined(ctx, fn)
}


func (c *Client) TxPipeline() Pipeliner {
	pipe := Pipeline{
		exec: func(ctx context.Context, cmds []Cmder) error {
			cmds = wrapMultiExec(ctx, cmds)
			return c.processTxPipelineHook(ctx, cmds)
		},
	}
	pipe.init()
	return &pipe
}

func (c *Client) pubSub() *PubSub {
	pubsub := &PubSub{
		opt: c.opt,

		newConn: func(ctx context.Context, channels []string) (*pool.Conn, error) {
			return c.newConn(ctx)
		},
		closeConn: c.connPool.CloseConn,
	}
	pubsub.init()
	return pubsub
}



























func (c *Client) Subscribe(ctx context.Context, channels ...string) *PubSub {
	pubsub := c.pubSub()
	if len(channels) > 0 {
		_ = pubsub.Subscribe(ctx, channels...)
	}
	return pubsub
}



func (c *Client) PSubscribe(ctx context.Context, channels ...string) *PubSub {
	pubsub := c.pubSub()
	if len(channels) > 0 {
		_ = pubsub.PSubscribe(ctx, channels...)
	}
	return pubsub
}



func (c *Client) SSubscribe(ctx context.Context, channels ...string) *PubSub {
	pubsub := c.pubSub()
	if len(channels) > 0 {
		_ = pubsub.SSubscribe(ctx, channels...)
	}
	return pubsub
}






type Conn struct {
	baseClient
	cmdable
	statefulCmdable
}




func newConn(opt *Options, connPool pool.Pooler, parentHooks *hooksMixin) *Conn {
	c := Conn{
		baseClient: baseClient{
			opt:      opt,
			connPool: connPool,
		},
	}

	if parentHooks != nil {
		c.hooksMixin = parentHooks.clone()
	}

	c.cmdable = c.Process
	c.statefulCmdable = c.Process
	c.initHooks(hooks{
		dial:       c.baseClient.dial,
		process:    c.baseClient.process,
		pipeline:   c.baseClient.processPipeline,
		txPipeline: c.baseClient.processTxPipeline,
	})

	return &c
}

func (c *Conn) Process(ctx context.Context, cmd Cmder) error {
	err := c.processHook(ctx, cmd)
	cmd.SetErr(err)
	return err
}

func (c *Conn) Pipelined(ctx context.Context, fn func(Pipeliner) error) ([]Cmder, error) {
	return c.Pipeline().Pipelined(ctx, fn)
}

func (c *Conn) Pipeline() Pipeliner {
	pipe := Pipeline{
		exec: c.processPipelineHook,
	}
	pipe.init()
	return &pipe
}

func (c *Conn) TxPipelined(ctx context.Context, fn func(Pipeliner) error) ([]Cmder, error) {
	return c.TxPipeline().Pipelined(ctx, fn)
}


func (c *Conn) TxPipeline() Pipeliner {
	pipe := Pipeline{
		exec: func(ctx context.Context, cmds []Cmder) error {
			cmds = wrapMultiExec(ctx, cmds)
			return c.processTxPipelineHook(ctx, cmds)
		},
	}
	pipe.init()
	return &pipe
}
