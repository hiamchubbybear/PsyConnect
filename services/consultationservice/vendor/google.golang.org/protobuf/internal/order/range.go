




package order

import (
	"sort"
	"sync"

	"google.golang.org/protobuf/reflect/protoreflect"
)

type messageField struct {
	fd protoreflect.FieldDescriptor
	v  protoreflect.Value
}

var messageFieldPool = sync.Pool{
	New: func() any { return new([]messageField) },
}

type (
	
	
	FieldRanger interface{ Range(VisitField) }
	
	VisitField = func(protoreflect.FieldDescriptor, protoreflect.Value) bool
)


func RangeFields(fs FieldRanger, less FieldOrder, fn VisitField) {
	if less == nil {
		fs.Range(fn)
		return
	}

	
	p := messageFieldPool.Get().(*[]messageField)
	fields := (*p)[:0]
	defer func() {
		if cap(fields) < 1024 {
			*p = fields
			messageFieldPool.Put(p)
		}
	}()

	
	fs.Range(func(fd protoreflect.FieldDescriptor, v protoreflect.Value) bool {
		fields = append(fields, messageField{fd, v})
		return true
	})
	sort.Slice(fields, func(i, j int) bool {
		return less(fields[i].fd, fields[j].fd)
	})

	
	for _, f := range fields {
		if !fn(f.fd, f.v) {
			return
		}
	}
}

type mapEntry struct {
	k protoreflect.MapKey
	v protoreflect.Value
}

var mapEntryPool = sync.Pool{
	New: func() any { return new([]mapEntry) },
}

type (
	
	
	EntryRanger interface{ Range(VisitEntry) }
	
	VisitEntry = func(protoreflect.MapKey, protoreflect.Value) bool
)


func RangeEntries(es EntryRanger, less KeyOrder, fn VisitEntry) {
	if less == nil {
		es.Range(fn)
		return
	}

	
	p := mapEntryPool.Get().(*[]mapEntry)
	entries := (*p)[:0]
	defer func() {
		if cap(entries) < 1024 {
			*p = entries
			mapEntryPool.Put(p)
		}
	}()

	
	es.Range(func(k protoreflect.MapKey, v protoreflect.Value) bool {
		entries = append(entries, mapEntry{k, v})
		return true
	})
	sort.Slice(entries, func(i, j int) bool {
		return less(entries[i].k, entries[j].k)
	})

	
	for _, e := range entries {
		if !fn(e.k, e.v) {
			return
		}
	}
}
