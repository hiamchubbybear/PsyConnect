




package protolazy

import (
	"fmt"
	"sort"

	"google.golang.org/protobuf/encoding/protowire"
	piface "google.golang.org/protobuf/runtime/protoiface"
)



type IndexEntry struct {
	FieldNum uint32
	
	Start uint32
	
	
	End uint32
	
	MultipleContiguous bool
}




type XXX_lazyUnmarshalInfo struct {
	
	
	
	
	
	index *[]IndexEntry
	
	
	
	
	Protobuf []byte
	
	unmarshalFlags piface.UnmarshalInputFlags
}







func (lazy *XXX_lazyUnmarshalInfo) Buffer() []byte {
	return lazy.Protobuf
}




func (lazy *XXX_lazyUnmarshalInfo) SetBuffer(b []byte) {
	lazy.Protobuf = b
}



func (lazy *XXX_lazyUnmarshalInfo) SetUnmarshalFlags(f piface.UnmarshalInputFlags) {
	lazy.unmarshalFlags = f
}


func (lazy *XXX_lazyUnmarshalInfo) UnmarshalFlags() piface.UnmarshalInputFlags {
	return lazy.unmarshalFlags
}



func (lazy *XXX_lazyUnmarshalInfo) AllowedPartial() bool {
	return (lazy.unmarshalFlags & piface.UnmarshalCheckRequired) == 0
}

func protoFieldNumber(tag uint32) uint32 {
	return tag >> 3
}



func buildIndex(buf []byte) ([]IndexEntry, error) {
	index := make([]IndexEntry, 0, 16)
	var lastProtoFieldNum uint32
	var outOfOrder bool

	var r BufferReader = NewBufferReader(buf)

	for !r.Done() {
		var tag uint32
		var err error
		var curPos = r.Pos
		
		{
			i := r.Pos
			buf := r.Buf

			if i >= len(buf) {
				return nil, errOutOfBounds
			} else if buf[i] < 0x80 {
				r.Pos++
				tag = uint32(buf[i])
			} else if r.Remaining() < 5 {
				var v uint64
				v, err = r.DecodeVarintSlow()
				tag = uint32(v)
			} else {
				var v uint32
				
				tag = uint32(buf[i]) & 127
				i++

				v = uint32(buf[i])
				i++
				tag |= (v & 127) << 7
				if v < 128 {
					goto done
				}

				v = uint32(buf[i])
				i++
				tag |= (v & 127) << 14
				if v < 128 {
					goto done
				}

				v = uint32(buf[i])
				i++
				tag |= (v & 127) << 21
				if v < 128 {
					goto done
				}

				v = uint32(buf[i])
				i++
				tag |= (v & 127) << 28
				if v < 128 {
					goto done
				}

				return nil, errOutOfBounds

			done:
				r.Pos = i
			}
		}
		

		fieldNum := protoFieldNumber(tag)
		if fieldNum < lastProtoFieldNum {
			outOfOrder = true
		}

		
		
		wireType := tag & 0x7
		switch protowire.Type(wireType) {
		case protowire.VarintType:
			
			i := r.Pos

			if len(r.Buf)-i < 10 {
				
				
				_, err = r.DecodeVarintSlow()
				goto out2
			}
			if r.Buf[i] < 0x80 {
				goto out
			}
			i++

			if r.Buf[i] < 0x80 {
				goto out
			}
			i++

			if r.Buf[i] < 0x80 {
				goto out
			}
			i++

			if r.Buf[i] < 0x80 {
				goto out
			}
			i++

			if r.Buf[i] < 0x80 {
				goto out
			}
			i++

			if r.Buf[i] < 0x80 {
				goto out
			}
			i++

			if r.Buf[i] < 0x80 {
				goto out
			}
			i++

			if r.Buf[i] < 0x80 {
				goto out
			}
			i++

			if r.Buf[i] < 0x80 {
				goto out
			}
			i++

			if r.Buf[i] < 0x80 {
				goto out
			}
			return nil, errOverflow
		out:
			r.Pos = i + 1
			
		case protowire.Fixed64Type:
			err = r.SkipFixed64()
		case protowire.BytesType:
			var n uint32
			n, err = r.DecodeVarint32()
			if err == nil {
				err = r.Skip(int(n))
			}
		case protowire.StartGroupType:
			err = r.SkipGroup(tag)
		case protowire.Fixed32Type:
			err = r.SkipFixed32()
		default:
			err = fmt.Errorf("Unexpected wire type (%d)", wireType)
		}
		

	out2:
		if err != nil {
			return nil, err
		}
		if fieldNum != lastProtoFieldNum {
			index = append(index, IndexEntry{FieldNum: fieldNum,
				Start: uint32(curPos),
				End:   uint32(r.Pos)},
			)
		} else {
			index[len(index)-1].End = uint32(r.Pos)
			index[len(index)-1].MultipleContiguous = true
		}
		lastProtoFieldNum = fieldNum
	}
	if outOfOrder {
		sort.Slice(index, func(i, j int) bool {
			return index[i].FieldNum < index[j].FieldNum ||
				(index[i].FieldNum == index[j].FieldNum &&
					index[i].Start < index[j].Start)
		})
	}
	return index, nil
}

func (lazy *XXX_lazyUnmarshalInfo) SizeField(num uint32) (size int) {
	start, end, found, _, multipleEntries := lazy.FindFieldInProto(num)
	if multipleEntries != nil {
		for _, entry := range multipleEntries {
			size += int(entry.End - entry.Start)
		}
		return size
	}
	if !found {
		return 0
	}
	return int(end - start)
}

func (lazy *XXX_lazyUnmarshalInfo) AppendField(b []byte, num uint32) ([]byte, bool) {
	start, end, found, _, multipleEntries := lazy.FindFieldInProto(num)
	if multipleEntries != nil {
		for _, entry := range multipleEntries {
			b = append(b, lazy.Protobuf[entry.Start:entry.End]...)
		}
		return b, true
	}
	if !found {
		return nil, false
	}
	b = append(b, lazy.Protobuf[start:end]...)
	return b, true
}

func (lazy *XXX_lazyUnmarshalInfo) SetIndex(index []IndexEntry) {
	atomicStoreIndex(&lazy.index, &index)
}



func (lazy *XXX_lazyUnmarshalInfo) FindFieldInProto(fieldNum uint32) (start, end uint32, found, multipleContiguous bool, multipleEntries []IndexEntry) {
	if lazy.Protobuf == nil {
		
		return 0, 0, false, false, nil
	}
	index := atomicLoadIndex(&lazy.index)
	if index == nil {
		r, err := buildIndex(lazy.Protobuf)
		if err != nil {
			panic(fmt.Sprintf("findFieldInfo: error building index when looking for field %d: %v", fieldNum, err))
		}
		
		index = &r
		atomicStoreIndex(&lazy.index, index)
	}
	return lookupField(index, fieldNum)
}











func lookupField(indexp *[]IndexEntry, fieldNum uint32) (start, end uint32, found bool, multipleContiguous bool, multipleEntries []IndexEntry) {
	
	
	
	index := *indexp
	for i, entry := range index {
		if fieldNum == entry.FieldNum {
			if i < len(index)-1 && entry.FieldNum == index[i+1].FieldNum {
				
				
				
				multiple := make([]IndexEntry, 1, 2)
				multiple[0] = IndexEntry{fieldNum, entry.Start, entry.End, entry.MultipleContiguous}
				i++
				for i < len(index) && index[i].FieldNum == fieldNum {
					multiple = append(multiple, IndexEntry{fieldNum, index[i].Start, index[i].End, index[i].MultipleContiguous})
					i++
				}
				return 0, 0, false, false, multiple

			}
			return entry.Start, entry.End, true, entry.MultipleContiguous, nil
		}
		if fieldNum < entry.FieldNum {
			return 0, 0, false, false, nil
		}
	}
	return 0, 0, false, false, nil
}
