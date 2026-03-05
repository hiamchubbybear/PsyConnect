package redis

import (
	"context"
)


type ScanIterator struct {
	cmd *ScanCmd
	pos int
}


func (it *ScanIterator) Err() error {
	return it.cmd.Err()
}


func (it *ScanIterator) Next(ctx context.Context) bool {
	
	if it.cmd.Err() != nil {
		return false
	}

	
	if it.pos < len(it.cmd.page) {
		it.pos++
		return true
	}

	for {
		
		if it.cmd.cursor == 0 {
			return false
		}

		
		switch it.cmd.args[0] {
		case "scan", "qscan":
			it.cmd.args[1] = it.cmd.cursor
		default:
			it.cmd.args[2] = it.cmd.cursor
		}

		err := it.cmd.process(ctx, it.cmd)
		if err != nil {
			return false
		}

		it.pos = 1

		
		if len(it.cmd.page) > 0 {
			return true
		}
	}
}


func (it *ScanIterator) Val() string {
	var v string
	if it.cmd.Err() == nil && it.pos > 0 && it.pos <= len(it.cmd.page) {
		v = it.cmd.page[it.pos-1]
	}
	return v
}
