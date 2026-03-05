









package aws

import (
	"io"
)













func ReadSeekCloser(r io.Reader) ReaderSeekerCloser {
	return ReaderSeekerCloser{r}
}



type ReaderSeekerCloser struct {
	r io.Reader
}




func IsReaderSeekable(r io.Reader) bool {
	switch v := r.(type) {
	case ReaderSeekerCloser:
		return v.IsSeeker()
	case *ReaderSeekerCloser:
		return v.IsSeeker()
	case io.ReadSeeker:
		return true
	default:
		return false
	}
}








func (r ReaderSeekerCloser) Read(p []byte) (int, error) {
	switch t := r.r.(type) {
	case io.Reader:
		return t.Read(p)
	}
	return 0, nil
}







func (r ReaderSeekerCloser) Seek(offset int64, whence int) (int64, error) {
	switch t := r.r.(type) {
	case io.Seeker:
		return t.Seek(offset, whence)
	}
	return int64(0), nil
}


func (r ReaderSeekerCloser) IsSeeker() bool {
	_, ok := r.r.(io.Seeker)
	return ok
}



func (r ReaderSeekerCloser) HasLen() (int, bool) {
	type lenner interface {
		Len() int
	}

	if lr, ok := r.r.(lenner); ok {
		return lr.Len(), true
	}

	return 0, false
}






func (r ReaderSeekerCloser) GetLen() (int64, error) {
	if l, ok := r.HasLen(); ok {
		return int64(l), nil
	}

	if s, ok := r.r.(io.Seeker); ok {
		return seekerLen(s)
	}

	return -1, nil
}



func SeekerLen(s io.Seeker) (int64, error) {
	
	
	switch v := s.(type) {
	case ReaderSeekerCloser:
		return v.GetLen()
	case *ReaderSeekerCloser:
		return v.GetLen()
	}

	return seekerLen(s)
}

func seekerLen(s io.Seeker) (int64, error) {
	curOffset, err := s.Seek(0, io.SeekCurrent)
	if err != nil {
		return 0, err
	}

	endOffset, err := s.Seek(0, io.SeekEnd)
	if err != nil {
		return 0, err
	}

	_, err = s.Seek(curOffset, io.SeekStart)
	if err != nil {
		return 0, err
	}

	return endOffset - curOffset, nil
}
