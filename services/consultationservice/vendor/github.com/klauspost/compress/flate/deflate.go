




package flate

import (
	"encoding/binary"
	"fmt"
	"io"
	"math"
)

const (
	NoCompression      = 0
	BestSpeed          = 1
	BestCompression    = 9
	DefaultCompression = -1

	
	
	
	
	
	
	
	
	
	HuffmanOnly         = -2
	ConstantCompression = HuffmanOnly 

	logWindowSize    = 15
	windowSize       = 1 << logWindowSize
	windowMask       = windowSize - 1
	logMaxOffsetSize = 15  
	minMatchLength   = 4   
	maxMatchLength   = 258 
	minOffsetSize    = 1   

	
	
	
	
	maxFlateBlockTokens = 1 << 15
	maxStoreBlockSize   = 65535
	hashBits            = 17 
	hashSize            = 1 << hashBits
	hashMask            = (1 << hashBits) - 1
	hashShift           = (hashBits + minMatchLength - 1) / minMatchLength
	maxHashOffset       = 1 << 28

	skipNever = math.MaxInt32

	debugDeflate = false
)

type compressionLevel struct {
	good, lazy, nice, chain, fastSkipHashing, level int
}




var levels = []compressionLevel{
	{}, 
	
	{0, 0, 0, 0, 0, 1},
	{0, 0, 0, 0, 0, 2},
	{0, 0, 0, 0, 0, 3},
	{0, 0, 0, 0, 0, 4},
	{0, 0, 0, 0, 0, 5},
	{0, 0, 0, 0, 0, 6},
	
	
	{8, 12, 16, 24, skipNever, 7},
	{16, 30, 40, 64, skipNever, 8},
	{32, 258, 258, 1024, skipNever, 9},
}


type advancedState struct {
	
	length         int
	offset         int
	maxInsertIndex int
	chainHead      int
	hashOffset     int

	ii uint16 

	
	index     int
	hashMatch [maxMatchLength + minMatchLength]uint32

	
	
	
	
	
	hashHead [hashSize]uint32
	hashPrev [windowSize]uint32
}

type compressor struct {
	compressionLevel

	h *huffmanEncoder
	w *huffmanBitWriter

	
	fill func(*compressor, []byte) int 
	step func(*compressor)             

	window     []byte
	windowEnd  int
	blockStart int 
	err        error

	
	tokens tokens
	fast   fastEnc
	state  *advancedState

	sync          bool 
	byteAvailable bool 
}

func (d *compressor) fillDeflate(b []byte) int {
	s := d.state
	if s.index >= 2*windowSize-(minMatchLength+maxMatchLength) {
		
		
		*(*[windowSize]byte)(d.window) = *(*[windowSize]byte)(d.window[windowSize:])
		s.index -= windowSize
		d.windowEnd -= windowSize
		if d.blockStart >= windowSize {
			d.blockStart -= windowSize
		} else {
			d.blockStart = math.MaxInt32
		}
		s.hashOffset += windowSize
		if s.hashOffset > maxHashOffset {
			delta := s.hashOffset - 1
			s.hashOffset -= delta
			s.chainHead -= delta
			
			
			for i, v := range s.hashPrev[:] {
				if int(v) > delta {
					s.hashPrev[i] = uint32(int(v) - delta)
				} else {
					s.hashPrev[i] = 0
				}
			}
			for i, v := range s.hashHead[:] {
				if int(v) > delta {
					s.hashHead[i] = uint32(int(v) - delta)
				} else {
					s.hashHead[i] = 0
				}
			}
		}
	}
	n := copy(d.window[d.windowEnd:], b)
	d.windowEnd += n
	return n
}

func (d *compressor) writeBlock(tok *tokens, index int, eof bool) error {
	if index > 0 || eof {
		var window []byte
		if d.blockStart <= index {
			window = d.window[d.blockStart:index]
		}
		d.blockStart = index
		
		d.w.writeBlockDynamic(tok, eof, window, d.sync)
		return d.w.err
	}
	return nil
}




func (d *compressor) writeBlockSkip(tok *tokens, index int, eof bool) error {
	if index > 0 || eof {
		if d.blockStart <= index {
			window := d.window[d.blockStart:index]
			
			
			if int(tok.n) > len(window)-int(tok.n>>6) {
				d.w.writeBlockHuff(eof, window, d.sync)
			} else {
				
				d.w.writeBlockDynamic(tok, eof, window, d.sync)
			}
		} else {
			d.w.writeBlock(tok, eof, nil)
		}
		d.blockStart = index
		return d.w.err
	}
	return nil
}





func (d *compressor) fillWindow(b []byte) {
	
	if d.level <= 0 {
		return
	}
	if d.fast != nil {
		
		if len(b) > maxMatchOffset {
			b = b[len(b)-maxMatchOffset:]
		}
		d.fast.Encode(&d.tokens, b)
		d.tokens.Reset()
		return
	}
	s := d.state
	
	if len(b) > windowSize {
		b = b[len(b)-windowSize:]
	}
	
	n := copy(d.window[d.windowEnd:], b)

	
	loops := (n + 256 - minMatchLength) / 256
	for j := 0; j < loops; j++ {
		startindex := j * 256
		end := startindex + 256 + minMatchLength - 1
		if end > n {
			end = n
		}
		tocheck := d.window[startindex:end]
		dstSize := len(tocheck) - minMatchLength + 1

		if dstSize <= 0 {
			continue
		}

		dst := s.hashMatch[:dstSize]
		bulkHash4(tocheck, dst)
		var newH uint32
		for i, val := range dst {
			di := i + startindex
			newH = val & hashMask
			
			
			s.hashPrev[di&windowMask] = s.hashHead[newH]
			
			s.hashHead[newH] = uint32(di + s.hashOffset)
		}
	}
	
	d.windowEnd += n
	s.index = n
}




func (d *compressor) findMatch(pos int, prevHead int, lookahead int) (length, offset int, ok bool) {
	minMatchLook := maxMatchLength
	if lookahead < minMatchLook {
		minMatchLook = lookahead
	}

	win := d.window[0 : pos+minMatchLook]

	
	nice := len(win) - pos
	if d.nice < nice {
		nice = d.nice
	}

	
	tries := d.chain
	length = minMatchLength - 1

	wEnd := win[pos+length]
	wPos := win[pos:]
	minIndex := pos - windowSize
	if minIndex < 0 {
		minIndex = 0
	}
	offset = 0

	if d.chain < 100 {
		for i := prevHead; tries > 0; tries-- {
			if wEnd == win[i+length] {
				n := matchLen(win[i:i+minMatchLook], wPos)
				if n > length {
					length = n
					offset = pos - i
					ok = true
					if n >= nice {
						
						break
					}
					wEnd = win[pos+n]
				}
			}
			if i <= minIndex {
				
				break
			}
			i = int(d.state.hashPrev[i&windowMask]) - d.state.hashOffset
			if i < minIndex {
				break
			}
		}
		return
	}

	
	cGain := 4

	
	const baseCost = 3
	
	

	for i := prevHead; tries > 0; tries-- {
		if wEnd == win[i+length] {
			n := matchLen(win[i:i+minMatchLook], wPos)
			if n > length {
				
				newGain := d.h.bitLengthRaw(wPos[:n]) - int(offsetExtraBits[offsetCode(uint32(pos-i))]) - baseCost - int(lengthExtraBits[lengthCodes[(n-3)&255]])

				
				if newGain > cGain {
					length = n
					offset = pos - i
					cGain = newGain
					ok = true
					if n >= nice {
						
						break
					}
					wEnd = win[pos+n]
				}
			}
		}
		if i <= minIndex {
			
			break
		}
		i = int(d.state.hashPrev[i&windowMask]) - d.state.hashOffset
		if i < minIndex {
			break
		}
	}
	return
}

func (d *compressor) writeStoredBlock(buf []byte) error {
	if d.w.writeStoredHeader(len(buf), false); d.w.err != nil {
		return d.w.err
	}
	d.w.writeBytes(buf)
	return d.w.err
}




func hash4(b []byte) uint32 {
	return hash4u(binary.LittleEndian.Uint32(b), hashBits)
}



func hash4u(u uint32, h uint8) uint32 {
	return (u * prime4bytes) >> (32 - h)
}



func bulkHash4(b []byte, dst []uint32) {
	if len(b) < 4 {
		return
	}
	hb := binary.LittleEndian.Uint32(b)

	dst[0] = hash4u(hb, hashBits)
	end := len(b) - 4 + 1
	for i := 1; i < end; i++ {
		hb = (hb >> 8) | uint32(b[i+3])<<24
		dst[i] = hash4u(hb, hashBits)
	}
}

func (d *compressor) initDeflate() {
	d.window = make([]byte, 2*windowSize)
	d.byteAvailable = false
	d.err = nil
	if d.state == nil {
		return
	}
	s := d.state
	s.index = 0
	s.hashOffset = 1
	s.length = minMatchLength - 1
	s.offset = 0
	s.chainHead = -1
}



func (d *compressor) deflateLazy() {
	s := d.state
	
	
	
	const sanity = debugDeflate

	if d.windowEnd-s.index < minMatchLength+maxMatchLength && !d.sync {
		return
	}
	if d.windowEnd != s.index && d.chain > 100 {
		
		if d.h == nil {
			d.h = newHuffmanEncoder(maxFlateBlockTokens)
		}
		var tmp [256]uint16
		for _, v := range d.window[s.index:d.windowEnd] {
			tmp[v]++
		}
		d.h.generate(tmp[:], 15)
	}

	s.maxInsertIndex = d.windowEnd - (minMatchLength - 1)

	for {
		if sanity && s.index > d.windowEnd {
			panic("index > windowEnd")
		}
		lookahead := d.windowEnd - s.index
		if lookahead < minMatchLength+maxMatchLength {
			if !d.sync {
				return
			}
			if sanity && s.index > d.windowEnd {
				panic("index > windowEnd")
			}
			if lookahead == 0 {
				
				if d.byteAvailable {
					
					d.tokens.AddLiteral(d.window[s.index-1])
					d.byteAvailable = false
				}
				if d.tokens.n > 0 {
					if d.err = d.writeBlock(&d.tokens, s.index, false); d.err != nil {
						return
					}
					d.tokens.Reset()
				}
				return
			}
		}
		if s.index < s.maxInsertIndex {
			
			hash := hash4(d.window[s.index:])
			ch := s.hashHead[hash]
			s.chainHead = int(ch)
			s.hashPrev[s.index&windowMask] = ch
			s.hashHead[hash] = uint32(s.index + s.hashOffset)
		}
		prevLength := s.length
		prevOffset := s.offset
		s.length = minMatchLength - 1
		s.offset = 0
		minIndex := s.index - windowSize
		if minIndex < 0 {
			minIndex = 0
		}

		if s.chainHead-s.hashOffset >= minIndex && lookahead > prevLength && prevLength < d.lazy {
			if newLength, newOffset, ok := d.findMatch(s.index, s.chainHead-s.hashOffset, lookahead); ok {
				s.length = newLength
				s.offset = newOffset
			}
		}

		if prevLength >= minMatchLength && s.length <= prevLength {
			
			
			
			
			const checkOff = 2

			
			if prevLength < maxMatchLength-checkOff {
				prevIndex := s.index - 1
				if prevIndex+prevLength < s.maxInsertIndex {
					end := lookahead
					if lookahead > maxMatchLength+checkOff {
						end = maxMatchLength + checkOff
					}
					end += prevIndex

					
					h := hash4(d.window[prevIndex+prevLength:])
					ch2 := int(s.hashHead[h]) - s.hashOffset - prevLength
					if prevIndex-ch2 != prevOffset && ch2 > minIndex+checkOff {
						length := matchLen(d.window[prevIndex+checkOff:end], d.window[ch2+checkOff:])
						
						if length > prevLength {
							prevLength = length
							prevOffset = prevIndex - ch2

							
							for i := checkOff - 1; i >= 0; i-- {
								if prevLength >= maxMatchLength || d.window[prevIndex+i] != d.window[ch2+i] {
									
									for j := 0; j <= i; j++ {
										d.tokens.AddLiteral(d.window[prevIndex+j])
										if d.tokens.n == maxFlateBlockTokens {
											
											if d.err = d.writeBlock(&d.tokens, s.index, false); d.err != nil {
												return
											}
											d.tokens.Reset()
										}
										s.index++
										if s.index < s.maxInsertIndex {
											h := hash4(d.window[s.index:])
											ch := s.hashHead[h]
											s.chainHead = int(ch)
											s.hashPrev[s.index&windowMask] = ch
											s.hashHead[h] = uint32(s.index + s.hashOffset)
										}
									}
									break
								} else {
									prevLength++
								}
							}
						} else if false {
							
							
							prevIndex++
							h := hash4(d.window[prevIndex+prevLength:])
							ch2 := int(s.hashHead[h]) - s.hashOffset - prevLength
							if prevIndex-ch2 != prevOffset && ch2 > minIndex+checkOff {
								length := matchLen(d.window[prevIndex+checkOff:end], d.window[ch2+checkOff:])
								
								if length > prevLength+checkOff {
									prevLength = length
									prevOffset = prevIndex - ch2
									prevIndex--

									
									for i := checkOff; i >= 0; i-- {
										if prevLength >= maxMatchLength || d.window[prevIndex+i] != d.window[ch2+i-1] {
											
											for j := 0; j <= i; j++ {
												d.tokens.AddLiteral(d.window[prevIndex+j])
												if d.tokens.n == maxFlateBlockTokens {
													
													if d.err = d.writeBlock(&d.tokens, s.index, false); d.err != nil {
														return
													}
													d.tokens.Reset()
												}
												s.index++
												if s.index < s.maxInsertIndex {
													h := hash4(d.window[s.index:])
													ch := s.hashHead[h]
													s.chainHead = int(ch)
													s.hashPrev[s.index&windowMask] = ch
													s.hashHead[h] = uint32(s.index + s.hashOffset)
												}
											}
											break
										} else {
											prevLength++
										}
									}
								}
							}
						}
					}
				}
			}
			
			
			d.tokens.AddMatch(uint32(prevLength-3), uint32(prevOffset-minOffsetSize))

			
			
			
			
			newIndex := s.index + prevLength - 1
			
			end := newIndex
			if end > s.maxInsertIndex {
				end = s.maxInsertIndex
			}
			end += minMatchLength - 1
			startindex := s.index + 1
			if startindex > s.maxInsertIndex {
				startindex = s.maxInsertIndex
			}
			tocheck := d.window[startindex:end]
			dstSize := len(tocheck) - minMatchLength + 1
			if dstSize > 0 {
				dst := s.hashMatch[:dstSize]
				bulkHash4(tocheck, dst)
				var newH uint32
				for i, val := range dst {
					di := i + startindex
					newH = val & hashMask
					
					
					s.hashPrev[di&windowMask] = s.hashHead[newH]
					
					s.hashHead[newH] = uint32(di + s.hashOffset)
				}
			}

			s.index = newIndex
			d.byteAvailable = false
			s.length = minMatchLength - 1
			if d.tokens.n == maxFlateBlockTokens {
				
				if d.err = d.writeBlock(&d.tokens, s.index, false); d.err != nil {
					return
				}
				d.tokens.Reset()
			}
			s.ii = 0
		} else {
			
			if s.length >= minMatchLength {
				s.ii = 0
			}
			
			if d.byteAvailable {
				s.ii++
				d.tokens.AddLiteral(d.window[s.index-1])
				if d.tokens.n == maxFlateBlockTokens {
					if d.err = d.writeBlock(&d.tokens, s.index, false); d.err != nil {
						return
					}
					d.tokens.Reset()
				}
				s.index++

				
				
				if n := int(s.ii) - d.chain; n > 0 {
					n = 1 + int(n>>6)
					for j := 0; j < n; j++ {
						if s.index >= d.windowEnd-1 {
							break
						}
						d.tokens.AddLiteral(d.window[s.index-1])
						if d.tokens.n == maxFlateBlockTokens {
							if d.err = d.writeBlock(&d.tokens, s.index, false); d.err != nil {
								return
							}
							d.tokens.Reset()
						}
						
						if s.index < s.maxInsertIndex {
							h := hash4(d.window[s.index:])
							ch := s.hashHead[h]
							s.chainHead = int(ch)
							s.hashPrev[s.index&windowMask] = ch
							s.hashHead[h] = uint32(s.index + s.hashOffset)
						}
						s.index++
					}
					
					d.tokens.AddLiteral(d.window[s.index-1])
					d.byteAvailable = false
					
					if d.tokens.n == maxFlateBlockTokens {
						if d.err = d.writeBlock(&d.tokens, s.index, false); d.err != nil {
							return
						}
						d.tokens.Reset()
					}
				}
			} else {
				s.index++
				d.byteAvailable = true
			}
		}
	}
}

func (d *compressor) store() {
	if d.windowEnd > 0 && (d.windowEnd == maxStoreBlockSize || d.sync) {
		d.err = d.writeStoredBlock(d.window[:d.windowEnd])
		d.windowEnd = 0
	}
}



func (d *compressor) fillBlock(b []byte) int {
	n := copy(d.window[d.windowEnd:], b)
	d.windowEnd += n
	return n
}




func (d *compressor) storeHuff() {
	if d.windowEnd < len(d.window) && !d.sync || d.windowEnd == 0 {
		return
	}
	d.w.writeBlockHuff(false, d.window[:d.windowEnd], d.sync)
	d.err = d.w.err
	d.windowEnd = 0
}




func (d *compressor) storeFast() {
	
	if d.windowEnd < len(d.window) {
		if !d.sync {
			return
		}
		
		if d.windowEnd < 128 {
			if d.windowEnd == 0 {
				return
			}
			if d.windowEnd <= 32 {
				d.err = d.writeStoredBlock(d.window[:d.windowEnd])
			} else {
				d.w.writeBlockHuff(false, d.window[:d.windowEnd], true)
				d.err = d.w.err
			}
			d.tokens.Reset()
			d.windowEnd = 0
			d.fast.Reset()
			return
		}
	}

	d.fast.Encode(&d.tokens, d.window[:d.windowEnd])
	
	if d.tokens.n == 0 {
		d.err = d.writeStoredBlock(d.window[:d.windowEnd])
		
	} else if int(d.tokens.n) > d.windowEnd-(d.windowEnd>>4) {
		d.w.writeBlockHuff(false, d.window[:d.windowEnd], d.sync)
		d.err = d.w.err
	} else {
		d.w.writeBlockDynamic(&d.tokens, false, d.window[:d.windowEnd], d.sync)
		d.err = d.w.err
	}
	d.tokens.Reset()
	d.windowEnd = 0
}



func (d *compressor) write(b []byte) (n int, err error) {
	if d.err != nil {
		return 0, d.err
	}
	n = len(b)
	for len(b) > 0 {
		if d.windowEnd == len(d.window) || d.sync {
			d.step(d)
		}
		b = b[d.fill(d, b):]
		if d.err != nil {
			return 0, d.err
		}
	}
	return n, d.err
}

func (d *compressor) syncFlush() error {
	d.sync = true
	if d.err != nil {
		return d.err
	}
	d.step(d)
	if d.err == nil {
		d.w.writeStoredHeader(0, false)
		d.w.flush()
		d.err = d.w.err
	}
	d.sync = false
	return d.err
}

func (d *compressor) init(w io.Writer, level int) (err error) {
	d.w = newHuffmanBitWriter(w)

	switch {
	case level == NoCompression:
		d.window = make([]byte, maxStoreBlockSize)
		d.fill = (*compressor).fillBlock
		d.step = (*compressor).store
	case level == ConstantCompression:
		d.w.logNewTablePenalty = 10
		d.window = make([]byte, 32<<10)
		d.fill = (*compressor).fillBlock
		d.step = (*compressor).storeHuff
	case level == DefaultCompression:
		level = 5
		fallthrough
	case level >= 1 && level <= 6:
		d.w.logNewTablePenalty = 7
		d.fast = newFastEnc(level)
		d.window = make([]byte, maxStoreBlockSize)
		d.fill = (*compressor).fillBlock
		d.step = (*compressor).storeFast
	case 7 <= level && level <= 9:
		d.w.logNewTablePenalty = 8
		d.state = &advancedState{}
		d.compressionLevel = levels[level]
		d.initDeflate()
		d.fill = (*compressor).fillDeflate
		d.step = (*compressor).deflateLazy
	default:
		return fmt.Errorf("flate: invalid compression level %d: want value in range [-2, 9]", level)
	}
	d.level = level
	return nil
}


func (d *compressor) reset(w io.Writer) {
	d.w.reset(w)
	d.sync = false
	d.err = nil
	
	if d.fast != nil {
		d.fast.Reset()
		d.windowEnd = 0
		d.tokens.Reset()
		return
	}
	switch d.compressionLevel.chain {
	case 0:
		
		d.windowEnd = 0
	default:
		s := d.state
		s.chainHead = -1
		for i := range s.hashHead {
			s.hashHead[i] = 0
		}
		for i := range s.hashPrev {
			s.hashPrev[i] = 0
		}
		s.hashOffset = 1
		s.index, d.windowEnd = 0, 0
		d.blockStart, d.byteAvailable = 0, false
		d.tokens.Reset()
		s.length = minMatchLength - 1
		s.offset = 0
		s.ii = 0
		s.maxInsertIndex = 0
	}
}

func (d *compressor) close() error {
	if d.err != nil {
		return d.err
	}
	d.sync = true
	d.step(d)
	if d.err != nil {
		return d.err
	}
	if d.w.writeStoredHeader(0, true); d.w.err != nil {
		return d.w.err
	}
	d.w.flush()
	d.w.reset(nil)
	return d.w.err
}













func NewWriter(w io.Writer, level int) (*Writer, error) {
	var dw Writer
	if err := dw.d.init(w, level); err != nil {
		return nil, err
	}
	return &dw, nil
}







func NewWriterDict(w io.Writer, level int, dict []byte) (*Writer, error) {
	zw, err := NewWriter(w, level)
	if err != nil {
		return nil, err
	}
	zw.d.fillWindow(dict)
	zw.dict = append(zw.dict, dict...) 
	return zw, err
}



type Writer struct {
	d    compressor
	dict []byte
}



func (w *Writer) Write(data []byte) (n int, err error) {
	return w.d.write(data)
}










func (w *Writer) Flush() error {
	
	
	return w.d.syncFlush()
}


func (w *Writer) Close() error {
	return w.d.close()
}




func (w *Writer) Reset(dst io.Writer) {
	if len(w.dict) > 0 {
		
		w.d.reset(dst)
		if dst != nil {
			w.d.fillWindow(w.dict)
		}
	} else {
		
		w.d.reset(dst)
	}
}




func (w *Writer) ResetDict(dst io.Writer, dict []byte) {
	w.dict = dict
	w.d.reset(dst)
	w.d.fillWindow(w.dict)
}
