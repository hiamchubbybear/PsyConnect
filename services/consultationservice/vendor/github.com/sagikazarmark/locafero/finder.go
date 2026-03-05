
package locafero

import (
	"errors"
	"io/fs"
	"path/filepath"
	"strings"

	"github.com/sourcegraph/conc/iter"
	"github.com/spf13/afero"
)


type Finder struct {
	
	
	
	
	
	
	
	Paths []string

	
	
	
	
	
	
	
	
	
	
	
	Names []string

	
	
	
	Type FileType
}


func (f Finder) Find(fsys afero.Fs) ([]string, error) {
	
	

	type searchItem struct {
		path string
		name string
	}

	var searchItems []searchItem

	for _, searchPath := range f.Paths {
		searchPath := searchPath

		for _, searchName := range f.Names {
			searchName := searchName

			searchItems = append(searchItems, searchItem{searchPath, searchName})

			
			
			
			
			
			
			
			
		}
	}

	
	
	
	

	allResults, err := iter.MapErr(searchItems, func(item *searchItem) ([]string, error) {
		
		if strings.ContainsAny(item.name, globMatch) {
			return globWalkSearch(fsys, item.path, item.name, f.Type)
		}

		return statSearch(fsys, item.path, item.name, f.Type)
	})
	if err != nil {
		return nil, err
	}

	var results []string

	for _, r := range allResults {
		results = append(results, r...)
	}

	
	

	return results, nil
}

func globWalkSearch(fsys afero.Fs, searchPath string, searchName string, searchType FileType) ([]string, error) {
	var results []string

	err := afero.Walk(fsys, searchPath, func(p string, fileInfo fs.FileInfo, err error) error {
		if err != nil {
			return err
		}

		
		if p == searchPath {
			return nil
		}

		var result error

		
		
		if fileInfo.IsDir() && filepath.Dir(p) == searchPath {
			result = fs.SkipDir
		}

		
		if !searchType.matchFileInfo(fileInfo) {
			return result
		}

		match, err := filepath.Match(searchName, fileInfo.Name())
		if err != nil {
			return err
		}

		if match {
			results = append(results, p)
		}

		return result
	})
	if err != nil {
		return results, err
	}

	return results, nil
}

func statSearch(fsys afero.Fs, searchPath string, searchName string, searchType FileType) ([]string, error) {
	filePath := filepath.Join(searchPath, searchName)

	fileInfo, err := fsys.Stat(filePath)
	if errors.Is(err, fs.ErrNotExist) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	
	if !searchType.matchFileInfo(fileInfo) {
		return nil, nil
	}

	return []string{filePath}, nil
}
