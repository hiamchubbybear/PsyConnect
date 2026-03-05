package protocol

import (
	"errors"
	"io"
	"time"
)






type RecordReader interface {
	
	
	
	
	
	
	ReadRecord() (*Record, error)
}


func NewRecordReader(records ...Record) RecordReader {
	switch len(records) {
	case 0:
		return emptyRecordReader{}
	default:
		r := &recordReader{records: make([]Record, len(records))}
		copy(r.records, records)
		return r
	}
}


func MultiRecordReader(batches ...RecordReader) RecordReader {
	switch len(batches) {
	case 0:
		return emptyRecordReader{}
	case 1:
		return batches[0]
	default:
		m := &multiRecordReader{batches: make([]RecordReader, len(batches))}
		copy(m.batches, batches)
		return m
	}
}

func forEachRecord(r RecordReader, f func(int, *Record) error) error {
	for i := 0; ; i++ {
		rec, err := r.ReadRecord()

		if err != nil {
			if errors.Is(err, io.EOF) {
				err = nil
			}
			return err
		}

		if err := handleRecord(i, rec, f); err != nil {
			return err
		}
	}
}

func handleRecord(i int, r *Record, f func(int, *Record) error) error {
	if r.Key != nil {
		defer r.Key.Close()
	}
	if r.Value != nil {
		defer r.Value.Close()
	}
	return f(i, r)
}

type recordReader struct {
	records []Record
	index   int
}

func (r *recordReader) ReadRecord() (*Record, error) {
	if i := r.index; i >= 0 && i < len(r.records) {
		r.index++
		return &r.records[i], nil
	}
	return nil, io.EOF
}

type multiRecordReader struct {
	batches []RecordReader
	index   int
}

func (m *multiRecordReader) ReadRecord() (*Record, error) {
	for {
		if m.index == len(m.batches) {
			return nil, io.EOF
		}
		r, err := m.batches[m.index].ReadRecord()
		if err == nil {
			return r, nil
		}
		if !errors.Is(err, io.EOF) {
			return nil, err
		}
		m.index++
	}
}



type optimizedRecordReader struct {
	records []optimizedRecord
	index   int
	buffer  Record
	headers [][]Header
}

func (r *optimizedRecordReader) ReadRecord() (*Record, error) {
	if i := r.index; i >= 0 && i < len(r.records) {
		rec := &r.records[i]
		r.index++
		r.buffer = Record{
			Offset: rec.offset,
			Time:   rec.time(),
			Key:    rec.key(),
			Value:  rec.value(),
		}
		if i < len(r.headers) {
			r.buffer.Headers = r.headers[i]
		}
		return &r.buffer, nil
	}
	return nil, io.EOF
}

type optimizedRecord struct {
	offset    int64
	timestamp int64
	keyRef    *pageRef
	valueRef  *pageRef
}

func (r *optimizedRecord) time() time.Time {
	return makeTime(r.timestamp)
}

func (r *optimizedRecord) key() Bytes {
	return makeBytes(r.keyRef)
}

func (r *optimizedRecord) value() Bytes {
	return makeBytes(r.valueRef)
}

func makeBytes(ref *pageRef) Bytes {
	if ref == nil {
		return nil
	}
	return ref
}

type emptyRecordReader struct{}

func (emptyRecordReader) ReadRecord() (*Record, error) { return nil, io.EOF }


type ControlRecord struct {
	Offset  int64
	Time    time.Time
	Version int16
	Type    int16
	Data    []byte
	Headers []Header
}

func ReadControlRecord(r *Record) (*ControlRecord, error) {
	if r.Key != nil {
		defer r.Key.Close()
	}
	if r.Value != nil {
		defer r.Value.Close()
	}

	k, err := ReadAll(r.Key)
	if err != nil {
		return nil, err
	}
	if k == nil {
		return nil, Error("invalid control record with nil key")
	}
	if len(k) != 4 {
		return nil, Errorf("invalid control record with key of size %d", len(k))
	}

	v, err := ReadAll(r.Value)
	if err != nil {
		return nil, err
	}

	c := &ControlRecord{
		Offset:  r.Offset,
		Time:    r.Time,
		Version: readInt16(k[:2]),
		Type:    readInt16(k[2:]),
		Data:    v,
		Headers: r.Headers,
	}

	return c, nil
}

func (cr *ControlRecord) Key() Bytes {
	k := make([]byte, 4)
	writeInt16(k[:2], cr.Version)
	writeInt16(k[2:], cr.Type)
	return NewBytes(k)
}

func (cr *ControlRecord) Value() Bytes {
	return NewBytes(cr.Data)
}

func (cr *ControlRecord) Record() Record {
	return Record{
		Offset:  cr.Offset,
		Time:    cr.Time,
		Key:     cr.Key(),
		Value:   cr.Value(),
		Headers: cr.Headers,
	}
}



type ControlBatch struct {
	Attributes           Attributes
	PartitionLeaderEpoch int32
	BaseOffset           int64
	ProducerID           int64
	ProducerEpoch        int16
	BaseSequence         int32
	Records              RecordReader
}



func NewControlBatch(records ...ControlRecord) *ControlBatch {
	rawRecords := make([]Record, len(records))
	for i, cr := range records {
		rawRecords[i] = cr.Record()
	}
	return &ControlBatch{
		Records: NewRecordReader(rawRecords...),
	}
}

func (c *ControlBatch) ReadRecord() (*Record, error) {
	return c.Records.ReadRecord()
}

func (c *ControlBatch) ReadControlRecord() (*ControlRecord, error) {
	r, err := c.ReadRecord()
	if err != nil {
		return nil, err
	}
	if r.Key != nil {
		defer r.Key.Close()
	}
	if r.Value != nil {
		defer r.Value.Close()
	}
	return ReadControlRecord(r)
}

func (c *ControlBatch) Offset() int64 {
	return c.BaseOffset
}

func (c *ControlBatch) Version() int {
	return 2
}



type RecordBatch struct {
	Attributes           Attributes
	PartitionLeaderEpoch int32
	BaseOffset           int64
	ProducerID           int64
	ProducerEpoch        int16
	BaseSequence         int32
	Records              RecordReader
}

func (r *RecordBatch) ReadRecord() (*Record, error) {
	return r.Records.ReadRecord()
}

func (r *RecordBatch) Offset() int64 {
	return r.BaseOffset
}

func (r *RecordBatch) Version() int {
	return 2
}



type MessageSet struct {
	Attributes Attributes
	BaseOffset int64
	Records    RecordReader
}

func (m *MessageSet) ReadRecord() (*Record, error) {
	return m.Records.ReadRecord()
}

func (m *MessageSet) Offset() int64 {
	return m.BaseOffset
}

func (m *MessageSet) Version() int {
	return 1
}




type RecordStream struct {
	Records []RecordReader
	index   int
}

func (s *RecordStream) ReadRecord() (*Record, error) {
	for {
		if s.index < 0 || s.index >= len(s.Records) {
			return nil, io.EOF
		}

		if _, isControl := s.Records[s.index].(*ControlBatch); isControl {
			s.index++
			continue
		}

		r, err := s.Records[s.index].ReadRecord()
		if err != nil {
			if errors.Is(err, io.EOF) {
				s.index++
				continue
			}
		}

		return r, err
	}
}
