





package windows

import (
	errorspkg "errors"
	"unsafe"
)










func EscapeArg(s string) string {
	if len(s) == 0 {
		return `""`
	}
	n := len(s)
	hasSpace := false
	for i := 0; i < len(s); i++ {
		switch s[i] {
		case '"', '\\':
			n++
		case ' ', '\t':
			hasSpace = true
		}
	}
	if hasSpace {
		n += 2 
	}
	if n == len(s) {
		return s
	}

	qs := make([]byte, n)
	j := 0
	if hasSpace {
		qs[j] = '"'
		j++
	}
	slashes := 0
	for i := 0; i < len(s); i++ {
		switch s[i] {
		default:
			slashes = 0
			qs[j] = s[i]
		case '\\':
			slashes++
			qs[j] = s[i]
		case '"':
			for ; slashes > 0; slashes-- {
				qs[j] = '\\'
				j++
			}
			qs[j] = '\\'
			j++
			qs[j] = s[i]
		}
		j++
	}
	if hasSpace {
		for ; slashes > 0; slashes-- {
			qs[j] = '\\'
			j++
		}
		qs[j] = '"'
		j++
	}
	return string(qs[:j])
}




func ComposeCommandLine(args []string) string {
	if len(args) == 0 {
		return ""
	}

	
	
	
	
	
	
	prog := args[0]
	mustQuote := len(prog) == 0
	for i := 0; i < len(prog); i++ {
		c := prog[i]
		if c <= ' ' || (c == '"' && i == 0) {
			
			
			
			
			
			mustQuote = true
			break
		}
	}
	var commandLine []byte
	if mustQuote {
		commandLine = make([]byte, 0, len(prog)+2)
		commandLine = append(commandLine, '"')
		for i := 0; i < len(prog); i++ {
			c := prog[i]
			if c == '"' {
				
				
				
				continue
			}
			commandLine = append(commandLine, c)
		}
		commandLine = append(commandLine, '"')
	} else {
		if len(args) == 1 {
			
			
			return prog
		}
		commandLine = []byte(prog)
	}

	for _, arg := range args[1:] {
		commandLine = append(commandLine, ' ')
		
		
		
		commandLine = append(commandLine, EscapeArg(arg)...)
	}
	return string(commandLine)
}





func DecomposeCommandLine(commandLine string) ([]string, error) {
	if len(commandLine) == 0 {
		return []string{}, nil
	}
	utf16CommandLine, err := UTF16FromString(commandLine)
	if err != nil {
		return nil, errorspkg.New("string with NUL passed to DecomposeCommandLine")
	}
	var argc int32
	argv, err := commandLineToArgv(&utf16CommandLine[0], &argc)
	if err != nil {
		return nil, err
	}
	defer LocalFree(Handle(unsafe.Pointer(argv)))

	var args []string
	for _, p := range unsafe.Slice(argv, argc) {
		args = append(args, UTF16PtrToString(p))
	}
	return args, nil
}











func CommandLineToArgv(cmd *uint16, argc *int32) (argv *[8192]*[8192]uint16, err error) {
	argp, err := commandLineToArgv(cmd, argc)
	argv = (*[8192]*[8192]uint16)(unsafe.Pointer(argp))
	return argv, err
}

func CloseOnExec(fd Handle) {
	SetHandleInformation(Handle(fd), HANDLE_FLAG_INHERIT, 0)
}


func FullPath(name string) (path string, err error) {
	p, err := UTF16PtrFromString(name)
	if err != nil {
		return "", err
	}
	n := uint32(100)
	for {
		buf := make([]uint16, n)
		n, err = GetFullPathName(p, uint32(len(buf)), &buf[0], nil)
		if err != nil {
			return "", err
		}
		if n <= uint32(len(buf)) {
			return UTF16ToString(buf[:n]), nil
		}
	}
}


func NewProcThreadAttributeList(maxAttrCount uint32) (*ProcThreadAttributeListContainer, error) {
	var size uintptr
	err := initializeProcThreadAttributeList(nil, maxAttrCount, 0, &size)
	if err != ERROR_INSUFFICIENT_BUFFER {
		if err == nil {
			return nil, errorspkg.New("unable to query buffer size from InitializeProcThreadAttributeList")
		}
		return nil, err
	}
	alloc, err := LocalAlloc(LMEM_FIXED, uint32(size))
	if err != nil {
		return nil, err
	}
	
	al := &ProcThreadAttributeListContainer{data: (*ProcThreadAttributeList)(unsafe.Pointer(alloc))}
	err = initializeProcThreadAttributeList(al.data, maxAttrCount, 0, &size)
	if err != nil {
		return nil, err
	}
	return al, err
}


func (al *ProcThreadAttributeListContainer) Update(attribute uintptr, value unsafe.Pointer, size uintptr) error {
	al.pointers = append(al.pointers, value)
	return updateProcThreadAttribute(al.data, 0, attribute, value, size, nil, nil)
}


func (al *ProcThreadAttributeListContainer) Delete() {
	deleteProcThreadAttributeList(al.data)
	LocalFree(Handle(unsafe.Pointer(al.data)))
	al.data = nil
	al.pointers = nil
}


func (al *ProcThreadAttributeListContainer) List() *ProcThreadAttributeList {
	return al.data
}
