






package readpref 

import (
	"bytes"
	"errors"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/tag"
)

var (
	errInvalidReadPreference = errors.New("can not specify tags, max staleness, or hedge with mode primary")
)


func Primary() *ReadPref {
	return &ReadPref{mode: PrimaryMode}
}


func PrimaryPreferred(opts ...Option) *ReadPref {
	
	rp, _ := New(PrimaryPreferredMode, opts...)
	return rp
}


func SecondaryPreferred(opts ...Option) *ReadPref {
	
	rp, _ := New(SecondaryPreferredMode, opts...)
	return rp
}


func Secondary(opts ...Option) *ReadPref {
	
	rp, _ := New(SecondaryMode, opts...)
	return rp
}


func Nearest(opts ...Option) *ReadPref {
	
	rp, _ := New(NearestMode, opts...)
	return rp
}


func New(mode Mode, opts ...Option) (*ReadPref, error) {
	rp := &ReadPref{
		mode: mode,
	}

	if mode == PrimaryMode && len(opts) != 0 {
		return nil, errInvalidReadPreference
	}

	for _, opt := range opts {
		if opt == nil {
			continue
		}
		err := opt(rp)
		if err != nil {
			return nil, err
		}
	}

	return rp, nil
}


type ReadPref struct {
	maxStaleness    time.Duration
	maxStalenessSet bool
	mode            Mode
	tagSets         []tag.Set
	hedgeEnabled    *bool
}




func (r *ReadPref) MaxStaleness() (time.Duration, bool) {
	return r.maxStaleness, r.maxStalenessSet
}


func (r *ReadPref) Mode() Mode {
	return r.mode
}



func (r *ReadPref) TagSets() []tag.Set {
	return r.tagSets
}



func (r *ReadPref) HedgeEnabled() *bool {
	return r.hedgeEnabled
}


func (r *ReadPref) String() string {
	var b bytes.Buffer
	b.WriteString(r.mode.String())
	delim := "("
	if r.maxStalenessSet {
		fmt.Fprintf(&b, "%smaxStaleness=%v", delim, r.maxStaleness)
		delim = " "
	}
	for _, tagSet := range r.tagSets {
		fmt.Fprintf(&b, "%stagSet=%s", delim, tagSet.String())
		delim = " "
	}
	if r.hedgeEnabled != nil {
		fmt.Fprintf(&b, "%shedgeEnabled=%v", delim, *r.hedgeEnabled)
		delim = " "
	}
	if delim != "(" {
		b.WriteString(")")
	}
	return b.String()
}
