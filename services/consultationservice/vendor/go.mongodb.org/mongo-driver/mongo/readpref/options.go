





package readpref

import (
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/tag"
)


var ErrInvalidTagSet = errors.New("an even number of tags must be specified")


type Option func(*ReadPref) error



func WithMaxStaleness(ms time.Duration) Option {
	return func(rp *ReadPref) error {
		rp.maxStaleness = ms
		rp.maxStalenessSet = true
		return nil
	}
}









func WithTags(tags ...string) Option {
	return func(rp *ReadPref) error {
		length := len(tags)
		if length < 2 || length%2 != 0 {
			return ErrInvalidTagSet
		}

		tagset := make(tag.Set, 0, length/2)

		for i := 1; i < length; i += 2 {
			tagset = append(tagset, tag.Tag{Name: tags[i-1], Value: tags[i]})
		}

		return WithTagSets(tagset)(rp)
	}
}











func WithTagSets(tagSets ...tag.Set) Option {
	return func(rp *ReadPref) error {
		rp.tagSets = tagSets
		return nil
	}
}





func WithHedgeEnabled(hedgeEnabled bool) Option {
	return func(rp *ReadPref) error {
		rp.hedgeEnabled = &hedgeEnabled
		return nil
	}
}
