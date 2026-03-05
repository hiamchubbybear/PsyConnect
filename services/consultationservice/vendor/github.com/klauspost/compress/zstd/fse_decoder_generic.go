//go:build !amd64 || appengine || !gc || noasm
// +build !amd64 appengine !gc noasm

package zstd

import (
	"errors"
	"fmt"
)


func (s *fseDecoder) buildDtable() error {
	tableSize := uint32(1 << s.actualTableLog)
	highThreshold := tableSize - 1
	symbolNext := s.stateTable[:256]

	
	{
		for i, v := range s.norm[:s.symbolLen] {
			if v == -1 {
				s.dt[highThreshold].setAddBits(uint8(i))
				highThreshold--
				symbolNext[i] = 1
			} else {
				symbolNext[i] = uint16(v)
			}
		}
	}

	
	{
		tableMask := tableSize - 1
		step := tableStep(tableSize)
		position := uint32(0)
		for ss, v := range s.norm[:s.symbolLen] {
			for i := 0; i < int(v); i++ {
				s.dt[position].setAddBits(uint8(ss))
				position = (position + step) & tableMask
				for position > highThreshold {
					
					position = (position + step) & tableMask
				}
			}
		}
		if position != 0 {
			
			return errors.New("corrupted input (position != 0)")
		}
	}

	
	{
		tableSize := uint16(1 << s.actualTableLog)
		for u, v := range s.dt[:tableSize] {
			symbol := v.addBits()
			nextState := symbolNext[symbol]
			symbolNext[symbol] = nextState + 1
			nBits := s.actualTableLog - byte(highBits(uint32(nextState)))
			s.dt[u&maxTableMask].setNBits(nBits)
			newState := (nextState << nBits) - tableSize
			if newState > tableSize {
				return fmt.Errorf("newState (%d) outside table size (%d)", newState, tableSize)
			}
			if newState == uint16(u) && nBits == 0 {
				
				return fmt.Errorf("newState (%d) == oldState (%d) and no bits", newState, u)
			}
			s.dt[u&maxTableMask].setNewState(newState)
		}
	}
	return nil
}
