




package murmur3

type bmixer interface {
	bmix(p []byte) (tail []byte)
	Size() (n int)
	reset()
}

type digest struct {
	clen int      
	tail []byte   
	buf  [16]byte 
	seed uint32   
	bmixer
}

func (d *digest) BlockSize() int { return 1 }

func (d *digest) Write(p []byte) (n int, err error) {
	n = len(p)
	d.clen += n

	if len(d.tail) > 0 {
		
		nfree := d.Size() - len(d.tail) 
		if nfree < len(p) {
			
			block := append(d.tail, p[:nfree]...)
			p = p[nfree:]
			_ = d.bmix(block) 
		} else {
			
			p = append(d.tail, p...)
		}
	}

	d.tail = d.bmix(p)

	
	nn := copy(d.buf[:], d.tail)
	d.tail = d.buf[:nn]

	return n, nil
}

func (d *digest) Reset() {
	d.clen = 0
	d.tail = nil
	d.bmixer.reset()
}
