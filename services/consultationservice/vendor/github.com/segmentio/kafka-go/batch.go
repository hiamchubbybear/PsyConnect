package kafka

import (
	"bufio"
	"errors"
	"io"
	"sync"
	"time"
)










type Batch struct {
	mutex         sync.Mutex
	conn          *Conn
	lock          *sync.Mutex
	msgs          *messageSetReader
	deadline      time.Time
	throttle      time.Duration
	topic         string
	partition     int
	offset        int64
	highWaterMark int64
	err           error
	
	
	
	
	
	
	
	
	lastOffset int64
}



func (batch *Batch) Throttle() time.Duration {
	return batch.throttle
}


func (batch *Batch) HighWaterMark() int64 {
	return batch.highWaterMark
}


func (batch *Batch) Partition() int {
	return batch.partition
}


func (batch *Batch) Offset() int64 {
	batch.mutex.Lock()
	offset := batch.offset
	batch.mutex.Unlock()
	return offset
}



func (batch *Batch) Close() error {
	batch.mutex.Lock()
	err := batch.close()
	batch.mutex.Unlock()
	return err
}

func (batch *Batch) close() (err error) {
	conn := batch.conn
	lock := batch.lock

	batch.conn = nil
	batch.lock = nil

	if batch.msgs != nil {
		batch.msgs.discard()
	}

	if batch.msgs != nil && batch.msgs.decompressed != nil {
		releaseBuffer(batch.msgs.decompressed)
		batch.msgs.decompressed = nil
	}

	if err = batch.err; errors.Is(batch.err, io.EOF) {
		err = nil
	}

	if conn != nil {
		conn.rdeadline.unsetConnReadDeadline()
		conn.mutex.Lock()
		conn.offset = batch.offset
		conn.mutex.Unlock()

		if err != nil {
			var kafkaError Error
			if !errors.As(err, &kafkaError) && !errors.Is(err, io.ErrShortBuffer) {
				conn.Close()
			}
		}
	}

	if lock != nil {
		lock.Unlock()
	}

	return
}












func (batch *Batch) Err() error { return batch.err }











func (batch *Batch) Read(b []byte) (int, error) {
	n := 0

	batch.mutex.Lock()
	offset := batch.offset

	_, _, _, err := batch.readMessage(
		func(r *bufio.Reader, size int, nbytes int) (int, error) {
			if nbytes < 0 {
				return size, nil
			}
			return discardN(r, size, nbytes)
		},
		func(r *bufio.Reader, size int, nbytes int) (int, error) {
			if nbytes < 0 {
				return size, nil
			}
			
			
			if nbytes > size {
				return size, errShortRead
			}
			n = nbytes 
			if nbytes > cap(b) {
				nbytes = cap(b)
			}
			if nbytes > len(b) {
				b = b[:nbytes]
			}
			nbytes, err := io.ReadFull(r, b[:nbytes])
			if err != nil {
				return size - nbytes, err
			}
			return discardN(r, size-nbytes, n-nbytes)
		},
	)

	if err == nil && n > len(b) {
		n, err = len(b), io.ErrShortBuffer
		batch.err = io.ErrShortBuffer
		batch.offset = offset 
	}

	batch.mutex.Unlock()
	return n, err
}






func (batch *Batch) ReadMessage() (Message, error) {
	msg := Message{}
	batch.mutex.Lock()

	var offset, timestamp int64
	var headers []Header
	var err error

	offset, timestamp, headers, err = batch.readMessage(
		func(r *bufio.Reader, size int, nbytes int) (remain int, err error) {
			msg.Key, remain, err = readNewBytes(r, size, nbytes)
			return
		},
		func(r *bufio.Reader, size int, nbytes int) (remain int, err error) {
			msg.Value, remain, err = readNewBytes(r, size, nbytes)
			return
		},
	)
	
	
	for batch.conn != nil && offset < batch.conn.offset {
		if err != nil {
			break
		}
		offset, timestamp, headers, err = batch.readMessage(
			func(r *bufio.Reader, size int, nbytes int) (remain int, err error) {
				msg.Key, remain, err = readNewBytes(r, size, nbytes)
				return
			},
			func(r *bufio.Reader, size int, nbytes int) (remain int, err error) {
				msg.Value, remain, err = readNewBytes(r, size, nbytes)
				return
			},
		)
	}

	batch.mutex.Unlock()
	msg.Topic = batch.topic
	msg.Partition = batch.partition
	msg.Offset = offset
	msg.HighWaterMark = batch.highWaterMark
	msg.Time = makeTime(timestamp)
	msg.Headers = headers

	return msg, err
}

func (batch *Batch) readMessage(
	key func(*bufio.Reader, int, int) (int, error),
	val func(*bufio.Reader, int, int) (int, error),
) (offset int64, timestamp int64, headers []Header, err error) {
	if err = batch.err; err != nil {
		return
	}

	var lastOffset int64
	offset, lastOffset, timestamp, headers, err = batch.msgs.readMessage(batch.offset, key, val)
	switch {
	case err == nil:
		batch.offset = offset + 1
		batch.lastOffset = lastOffset
	case errors.Is(err, errShortRead):
		
		
		
		err = batch.msgs.discard()
		switch {
		case err != nil:
			
			
			
			
			
			batch.err = dontExpectEOF(err)
		case batch.msgs.remaining() == 0:
			
			
			
			
			
			
			
			err = checkTimeoutErr(batch.deadline)
			batch.err = err

			
			
			
			
			
			
			if errors.Is(batch.err, io.EOF) && batch.msgs.lengthRemain == 0 && batch.lastOffset != -1 {
				
				
				
				
				
				
				
				batch.offset = batch.lastOffset + 1
			}
		}
	default:
		
		
		
		
		
		batch.err = dontExpectEOF(err)
	}

	return
}

func checkTimeoutErr(deadline time.Time) (err error) {
	if !deadline.IsZero() && time.Now().After(deadline) {
		err = RequestTimedOut
	} else {
		err = io.EOF
	}
	return
}
