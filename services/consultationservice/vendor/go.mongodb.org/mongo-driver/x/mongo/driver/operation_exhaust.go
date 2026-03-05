





package driver

import (
	"context"
	"errors"
)



func (op Operation) ExecuteExhaust(ctx context.Context, conn StreamerConnection) error {
	if !conn.CurrentlyStreaming() {
		return errors.New("exhaust read must be done with a connection that is currently streaming")
	}

	res, err := op.readWireMessage(ctx, conn)
	if err != nil {
		return err
	}
	if op.ProcessResponseFn != nil {
		
		info := ResponseInfo{
			ServerResponse: res,
			Connection:     conn,
		}
		if err = op.ProcessResponseFn(info); err != nil {
			return err
		}
	}

	return nil
}
