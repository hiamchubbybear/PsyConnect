



package zstd

import (
	"github.com/klauspost/compress/huff0"
)


type history struct {
	
	huffTree *huff0.Scratch

	
	decoders      sequenceDecs
	recentOffsets [3]int

	
	b []byte

	
	
	ignoreBuffer int

	windowSize       int
	allocFrameBuffer int 
	error            bool
	dict             *dict
}



func (h *history) reset() {
	h.b = h.b[:0]
	h.ignoreBuffer = 0
	h.error = false
	h.recentOffsets = [3]int{1, 4, 8}
	h.decoders.freeDecoders()
	h.decoders = sequenceDecs{br: h.decoders.br}
	h.freeHuffDecoder()
	h.huffTree = nil
	h.dict = nil
	
}

func (h *history) freeHuffDecoder() {
	if h.huffTree != nil {
		if h.dict == nil || h.dict.litEnc != h.huffTree {
			huffDecoderPool.Put(h.huffTree)
			h.huffTree = nil
		}
	}
}

func (h *history) setDict(dict *dict) {
	if dict == nil {
		return
	}
	h.dict = dict
	h.decoders.litLengths = dict.llDec
	h.decoders.offsets = dict.ofDec
	h.decoders.matchLengths = dict.mlDec
	h.decoders.dict = dict.content
	h.recentOffsets = dict.offsets
	h.huffTree = dict.litEnc
}




func (h *history) append(b []byte) {
	if len(b) >= h.windowSize {
		
		h.b = h.b[:h.windowSize]
		copy(h.b, b[len(b)-h.windowSize:])
		return
	}

	
	if len(b) < cap(h.b)-len(h.b) {
		h.b = append(h.b, b...)
		return
	}

	
	
	discard := len(b) + len(h.b) - h.windowSize
	copy(h.b, h.b[discard:])
	h.b = h.b[:h.windowSize]
	copy(h.b[h.windowSize-len(b):], b)
}


func (h *history) ensureBlock() {
	if cap(h.b) < h.allocFrameBuffer {
		h.b = make([]byte, 0, h.allocFrameBuffer)
		return
	}

	avail := cap(h.b) - len(h.b)
	if avail >= h.windowSize || avail > maxCompressedBlockSize {
		return
	}
	
	
	discard := len(h.b) - h.windowSize
	copy(h.b, h.b[discard:])
	h.b = h.b[:h.windowSize]
}


func (h *history) appendKeep(b []byte) {
	h.b = append(h.b, b...)
}
