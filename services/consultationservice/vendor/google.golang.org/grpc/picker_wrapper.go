

package grpc

import (
	"context"
	"fmt"
	"io"
	"sync/atomic"

	"google.golang.org/grpc/balancer"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/internal/channelz"
	istatus "google.golang.org/grpc/internal/status"
	"google.golang.org/grpc/internal/transport"
	"google.golang.org/grpc/stats"
	"google.golang.org/grpc/status"
)



type pickerGeneration struct {
	
	
	picker balancer.Picker
	
	
	blockingCh chan struct{}
}



type pickerWrapper struct {
	
	pickerGen     atomic.Pointer[pickerGeneration]
	statsHandlers []stats.Handler 
}

func newPickerWrapper(statsHandlers []stats.Handler) *pickerWrapper {
	pw := &pickerWrapper{
		statsHandlers: statsHandlers,
	}
	pw.pickerGen.Store(&pickerGeneration{
		blockingCh: make(chan struct{}),
	})
	return pw
}



func (pw *pickerWrapper) updatePicker(p balancer.Picker) {
	old := pw.pickerGen.Swap(&pickerGeneration{
		picker:     p,
		blockingCh: make(chan struct{}),
	})
	close(old.blockingCh)
}






func doneChannelzWrapper(acbw *acBalancerWrapper, result *balancer.PickResult) {
	ac := acbw.ac
	ac.incrCallsStarted()
	done := result.Done
	result.Done = func(b balancer.DoneInfo) {
		if b.Err != nil && b.Err != io.EOF {
			ac.incrCallsFailed()
		} else {
			ac.incrCallsSucceeded()
		}
		if done != nil {
			done(b)
		}
	}
}








func (pw *pickerWrapper) pick(ctx context.Context, failfast bool, info balancer.PickInfo) (transport.ClientTransport, balancer.PickResult, error) {
	var ch chan struct{}

	var lastPickErr error

	for {
		pg := pw.pickerGen.Load()
		if pg == nil {
			return nil, balancer.PickResult{}, ErrClientConnClosing
		}
		if pg.picker == nil {
			ch = pg.blockingCh
		}
		if ch == pg.blockingCh {
			
			
			
			select {
			case <-ctx.Done():
				var errStr string
				if lastPickErr != nil {
					errStr = "latest balancer error: " + lastPickErr.Error()
				} else {
					errStr = fmt.Sprintf("%v while waiting for connections to become ready", ctx.Err())
				}
				switch ctx.Err() {
				case context.DeadlineExceeded:
					return nil, balancer.PickResult{}, status.Error(codes.DeadlineExceeded, errStr)
				case context.Canceled:
					return nil, balancer.PickResult{}, status.Error(codes.Canceled, errStr)
				}
			case <-ch:
			}
			continue
		}

		
		
		
		
		
		
		
		
		if ch != nil {
			for _, sh := range pw.statsHandlers {
				sh.HandleRPC(ctx, &stats.PickerUpdated{})
			}
		}

		ch = pg.blockingCh
		p := pg.picker

		pickResult, err := p.Pick(info)
		if err != nil {
			if err == balancer.ErrNoSubConnAvailable {
				continue
			}
			if st, ok := status.FromError(err); ok {
				
				
				if istatus.IsRestrictedControlPlaneCode(st) {
					err = status.Errorf(codes.Internal, "received picker error with illegal status: %v", err)
				}
				return nil, balancer.PickResult{}, dropError{error: err}
			}
			
			
			if !failfast {
				lastPickErr = err
				continue
			}
			return nil, balancer.PickResult{}, status.Error(codes.Unavailable, err.Error())
		}

		acbw, ok := pickResult.SubConn.(*acBalancerWrapper)
		if !ok {
			logger.Errorf("subconn returned from pick is type %T, not *acBalancerWrapper", pickResult.SubConn)
			continue
		}
		if t := acbw.ac.getReadyTransport(); t != nil {
			if channelz.IsOn() {
				doneChannelzWrapper(acbw, &pickResult)
				return t, pickResult, nil
			}
			return t, pickResult, nil
		}
		if pickResult.Done != nil {
			
			
			pickResult.Done(balancer.DoneInfo{})
		}
		logger.Infof("blockingPicker: the picked transport is not ready, loop back to repick")
		
		
		
		
	}
}

func (pw *pickerWrapper) close() {
	old := pw.pickerGen.Swap(nil)
	close(old.blockingCh)
}



func (pw *pickerWrapper) reset() {
	old := pw.pickerGen.Swap(&pickerGeneration{blockingCh: make(chan struct{})})
	close(old.blockingCh)
}



type dropError struct {
	error
}
