package kafka

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"io"
	"net"
	"strconv"
	"strings"
	"time"

	"github.com/segmentio/kafka-go/sasl"
)



type Dialer struct {
	
	ClientID string

	
	
	
	
	
	DialFunc func(ctx context.Context, network string, address string) (net.Conn, error)

	
	
	
	
	
	
	
	
	
	
	Timeout time.Duration

	
	
	
	
	Deadline time.Time

	
	
	
	LocalAddr net.Addr

	
	
	
	
	DualStack bool

	
	
	
	FallbackDelay time.Duration

	
	
	
	
	KeepAlive time.Duration

	
	
	
	
	
	
	
	Resolver Resolver

	
	
	TLS *tls.Config

	
	
	SASLMechanism sasl.Mechanism

	
	
	
	
	TransactionalID string
}


func (d *Dialer) Dial(network string, address string) (*Conn, error) {
	return d.DialContext(context.Background(), network, address)
}














func (d *Dialer) DialContext(ctx context.Context, network string, address string) (*Conn, error) {
	return d.connect(
		ctx,
		network,
		address,
		ConnConfig{
			ClientID:        d.ClientID,
			TransactionalID: d.TransactionalID,
		},
	)
}









func (d *Dialer) DialLeader(ctx context.Context, network string, address string, topic string, partition int) (*Conn, error) {
	p, err := d.LookupPartition(ctx, network, address, topic, partition)
	if err != nil {
		return nil, err
	}
	return d.DialPartition(ctx, network, address, p)
}




func (d *Dialer) DialPartition(ctx context.Context, network string, address string, partition Partition) (*Conn, error) {
	return d.connect(ctx, network, net.JoinHostPort(partition.Leader.Host, strconv.Itoa(partition.Leader.Port)), ConnConfig{
		ClientID:        d.ClientID,
		Topic:           partition.Topic,
		Partition:       partition.ID,
		Broker:          partition.Leader.ID,
		Rack:            partition.Leader.Rack,
		TransactionalID: d.TransactionalID,
	})
}



func (d *Dialer) LookupLeader(ctx context.Context, network string, address string, topic string, partition int) (Broker, error) {
	p, err := d.LookupPartition(ctx, network, address, topic, partition)
	return p.Leader, err
}


func (d *Dialer) LookupPartition(ctx context.Context, network string, address string, topic string, partition int) (Partition, error) {
	c, err := d.DialContext(ctx, network, address)
	if err != nil {
		return Partition{}, err
	}
	defer c.Close()

	brkch := make(chan Partition, 1)
	errch := make(chan error, 1)

	go func() {
		for attempt := 0; true; attempt++ {
			if attempt != 0 {
				if !sleep(ctx, backoff(attempt, 100*time.Millisecond, 10*time.Second)) {
					errch <- ctx.Err()
					return
				}
			}

			partitions, err := c.ReadPartitions(topic)
			if err != nil {
				if isTemporary(err) {
					continue
				}
				errch <- err
				return
			}

			for _, p := range partitions {
				if p.ID == partition {
					brkch <- p
					return
				}
			}
		}

		errch <- UnknownTopicOrPartition
	}()

	var prt Partition
	select {
	case prt = <-brkch:
	case err = <-errch:
	case <-ctx.Done():
		err = ctx.Err()
	}
	return prt, err
}


func (d *Dialer) LookupPartitions(ctx context.Context, network string, address string, topic string) ([]Partition, error) {
	conn, err := d.DialContext(ctx, network, address)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	prtch := make(chan []Partition, 1)
	errch := make(chan error, 1)

	go func() {
		if prt, err := conn.ReadPartitions(topic); err != nil {
			errch <- err
		} else {
			prtch <- prt
		}
	}()

	var prt []Partition
	select {
	case prt = <-prtch:
	case err = <-errch:
	case <-ctx.Done():
		err = ctx.Err()
	}
	return prt, err
}


func (d *Dialer) connectTLS(ctx context.Context, conn net.Conn, config *tls.Config) (tlsConn *tls.Conn, err error) {
	tlsConn = tls.Client(conn, config)
	errch := make(chan error)

	go func() {
		defer close(errch)
		errch <- tlsConn.Handshake()
	}()

	select {
	case <-ctx.Done():
		conn.Close()
		tlsConn.Close()
		<-errch 
		err = ctx.Err()

	case err = <-errch:
	}

	return
}



func (d *Dialer) connect(ctx context.Context, network, address string, connCfg ConnConfig) (*Conn, error) {
	if d.Timeout != 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, d.Timeout)
		defer cancel()
	}

	if !d.Deadline.IsZero() {
		var cancel context.CancelFunc
		ctx, cancel = context.WithDeadline(ctx, d.Deadline)
		defer cancel()
	}

	c, err := d.dialContext(ctx, network, address)
	if err != nil {
		return nil, fmt.Errorf("failed to dial: %w", err)
	}

	conn := NewConnWith(c, connCfg)

	if d.SASLMechanism != nil {
		host, port, err := splitHostPortNumber(address)
		if err != nil {
			return nil, fmt.Errorf("could not determine host/port for SASL authentication: %w", err)
		}
		metadata := &sasl.Metadata{
			Host: host,
			Port: port,
		}
		if err := d.authenticateSASL(sasl.WithMetadata(ctx, metadata), conn); err != nil {
			_ = conn.Close()
			return nil, fmt.Errorf("could not successfully authenticate to %s:%d with SASL: %w", host, port, err)
		}
	}

	return conn, nil
}







func (d *Dialer) authenticateSASL(ctx context.Context, conn *Conn) error {
	if err := conn.saslHandshake(d.SASLMechanism.Name()); err != nil {
		return fmt.Errorf("SASL handshake failed: %w", err)
	}

	sess, state, err := d.SASLMechanism.Start(ctx)
	if err != nil {
		return fmt.Errorf("SASL authentication process could not be started: %w", err)
	}

	for completed := false; !completed; {
		challenge, err := conn.saslAuthenticate(state)
		switch {
		case err == nil:
		case errors.Is(err, io.EOF):
			
			
			
			return SASLAuthenticationFailed
		default:
			return err
		}

		completed, state, err = sess.Next(ctx, challenge)
		if err != nil {
			return fmt.Errorf("SASL authentication process has failed: %w", err)
		}
	}

	return nil
}

func (d *Dialer) dialContext(ctx context.Context, network string, addr string) (net.Conn, error) {
	address, err := lookupHost(ctx, addr, d.Resolver)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve host: %w", err)
	}

	dial := d.DialFunc
	if dial == nil {
		dial = (&net.Dialer{
			LocalAddr:     d.LocalAddr,
			DualStack:     d.DualStack,
			FallbackDelay: d.FallbackDelay,
			KeepAlive:     d.KeepAlive,
		}).DialContext
	}

	conn, err := dial(ctx, network, address)
	if err != nil {
		return nil, fmt.Errorf("failed to open connection to %s: %w", address, err)
	}

	if d.TLS != nil {
		c := d.TLS
		
		
		if c.ServerName == "" {
			c = d.TLS.Clone()
			
			colonPos := strings.LastIndex(address, ":")
			if colonPos == -1 {
				colonPos = len(address)
			}
			hostname := address[:colonPos]
			c.ServerName = hostname
		}
		return d.connectTLS(ctx, conn, c)
	}

	return conn, nil
}


var DefaultDialer = &Dialer{
	Timeout:   10 * time.Second,
	DualStack: true,
}


func Dial(network string, address string) (*Conn, error) {
	return DefaultDialer.Dial(network, address)
}


func DialContext(ctx context.Context, network string, address string) (*Conn, error) {
	return DefaultDialer.DialContext(ctx, network, address)
}


func DialLeader(ctx context.Context, network string, address string, topic string, partition int) (*Conn, error) {
	return DefaultDialer.DialLeader(ctx, network, address, topic, partition)
}


func DialPartition(ctx context.Context, network string, address string, partition Partition) (*Conn, error) {
	return DefaultDialer.DialPartition(ctx, network, address, partition)
}


func LookupPartition(ctx context.Context, network string, address string, topic string, partition int) (Partition, error) {
	return DefaultDialer.LookupPartition(ctx, network, address, topic, partition)
}


func LookupPartitions(ctx context.Context, network string, address string, topic string) ([]Partition, error) {
	return DefaultDialer.LookupPartitions(ctx, network, address, topic)
}

func sleep(ctx context.Context, duration time.Duration) bool {
	if duration == 0 {
		select {
		default:
			return true
		case <-ctx.Done():
			return false
		}
	}
	timer := time.NewTimer(duration)
	defer timer.Stop()
	select {
	case <-timer.C:
		return true
	case <-ctx.Done():
		return false
	}
}

func backoff(attempt int, min time.Duration, max time.Duration) time.Duration {
	d := time.Duration(attempt*attempt) * min
	if d > max {
		d = max
	}
	return d
}

func canonicalAddress(s string) string {
	return net.JoinHostPort(splitHostPort(s))
}

func splitHostPort(s string) (host string, port string) {
	host, port, _ = net.SplitHostPort(s)
	if len(host) == 0 && len(port) == 0 {
		host = s
		port = "9092"
	}
	return
}

func splitHostPortNumber(s string) (host string, portNumber int, err error) {
	host, port := splitHostPort(s)
	portNumber, err = strconv.Atoi(port)
	if err != nil {
		return host, 0, fmt.Errorf("%s: %w", s, err)
	}
	return host, portNumber, nil
}

func lookupHost(ctx context.Context, address string, resolver Resolver) (string, error) {
	host, port := splitHostPort(address)

	if resolver != nil {
		resolved, err := resolver.LookupHost(ctx, host)
		if err != nil {
			return "", fmt.Errorf("failed to resolve host %s: %w", host, err)
		}

		
		
		if len(resolved) > 0 {
			resolvedHost, resolvedPort := splitHostPort(resolved[0])

			
			host = resolvedHost

			
			
			if port == "" {
				port = resolvedPort
			}
		}
	}

	return net.JoinHostPort(host, port), nil
}
