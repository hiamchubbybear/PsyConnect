package viper

import (
	"errors"

	"github.com/spf13/afero"
)


func WithFinder(f Finder) Option {
	return optionFunc(func(v *Viper) {
		if f == nil {
			return
		}

		v.finder = f
	})
}


type Finder interface {
	Find(fsys afero.Fs) ([]string, error)
}


func Finders(finders ...Finder) Finder {
	return &combinedFinder{finders: finders}
}


type combinedFinder struct {
	finders []Finder
}


func (c *combinedFinder) Find(fsys afero.Fs) ([]string, error) {
	var results []string
	var errs []error

	for _, finder := range c.finders {
		if finder == nil {
			continue
		}

		r, err := finder.Find(fsys)
		if err != nil {
			errs = append(errs, err)
			continue
		}

		results = append(results, r...)
	}

	return results, errors.Join(errs...)
}
