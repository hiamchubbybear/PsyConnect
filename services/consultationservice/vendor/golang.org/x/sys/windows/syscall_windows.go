





package windows

import (
	errorspkg "errors"
	"fmt"
	"runtime"
	"sync"
	"syscall"
	"time"
	"unicode/utf16"
	"unsafe"
)

type (
	Handle uintptr
	HWND   uintptr
)

const (
	InvalidHandle = ^Handle(0)
	InvalidHWND   = ^HWND(0)

	
	DDD_EXACT_MATCH_ON_REMOVE = 0x00000004
	DDD_NO_BROADCAST_SYSTEM   = 0x00000008
	DDD_RAW_TARGET_PATH       = 0x00000001
	DDD_REMOVE_DEFINITION     = 0x00000002

	
	DRIVE_UNKNOWN     = 0
	DRIVE_NO_ROOT_DIR = 1
	DRIVE_REMOVABLE   = 2
	DRIVE_FIXED       = 3
	DRIVE_REMOTE      = 4
	DRIVE_CDROM       = 5
	DRIVE_RAMDISK     = 6

	
	FILE_CASE_SENSITIVE_SEARCH        = 0x00000001
	FILE_CASE_PRESERVED_NAMES         = 0x00000002
	FILE_FILE_COMPRESSION             = 0x00000010
	FILE_DAX_VOLUME                   = 0x20000000
	FILE_NAMED_STREAMS                = 0x00040000
	FILE_PERSISTENT_ACLS              = 0x00000008
	FILE_READ_ONLY_VOLUME             = 0x00080000
	FILE_SEQUENTIAL_WRITE_ONCE        = 0x00100000
	FILE_SUPPORTS_ENCRYPTION          = 0x00020000
	FILE_SUPPORTS_EXTENDED_ATTRIBUTES = 0x00800000
	FILE_SUPPORTS_HARD_LINKS          = 0x00400000
	FILE_SUPPORTS_OBJECT_IDS          = 0x00010000
	FILE_SUPPORTS_OPEN_BY_FILE_ID     = 0x01000000
	FILE_SUPPORTS_REPARSE_POINTS      = 0x00000080
	FILE_SUPPORTS_SPARSE_FILES        = 0x00000040
	FILE_SUPPORTS_TRANSACTIONS        = 0x00200000
	FILE_SUPPORTS_USN_JOURNAL         = 0x02000000
	FILE_UNICODE_ON_DISK              = 0x00000004
	FILE_VOLUME_IS_COMPRESSED         = 0x00008000
	FILE_VOLUME_QUOTAS                = 0x00000020

	
	LOCKFILE_FAIL_IMMEDIATELY = 0x00000001
	LOCKFILE_EXCLUSIVE_LOCK   = 0x00000002

	
	WAIT_IO_COMPLETION = 0x000000C0
)




func StringToUTF16(s string) []uint16 {
	a, err := UTF16FromString(s)
	if err != nil {
		panic("windows: string with NUL passed to StringToUTF16")
	}
	return a
}




func UTF16FromString(s string) ([]uint16, error) {
	return syscall.UTF16FromString(s)
}



func UTF16ToString(s []uint16) string {
	return syscall.UTF16ToString(s)
}




func StringToUTF16Ptr(s string) *uint16 { return &StringToUTF16(s)[0] }




func UTF16PtrFromString(s string) (*uint16, error) {
	a, err := UTF16FromString(s)
	if err != nil {
		return nil, err
	}
	return &a[0], nil
}




func UTF16PtrToString(p *uint16) string {
	if p == nil {
		return ""
	}
	if *p == 0 {
		return ""
	}

	
	n := 0
	for ptr := unsafe.Pointer(p); *(*uint16)(ptr) != 0; n++ {
		ptr = unsafe.Pointer(uintptr(ptr) + unsafe.Sizeof(*p))
	}
	return UTF16ToString(unsafe.Slice(p, n))
}

func Getpagesize() int { return 4096 }




func NewCallback(fn interface{}) uintptr {
	return syscall.NewCallback(fn)
}




func NewCallbackCDecl(fn interface{}) uintptr {
	return syscall.NewCallbackCDecl(fn)
}












































































































































































































































































































































func GetCurrentProcess() (Handle, error) {
	return CurrentProcess(), nil
}



func CurrentProcess() Handle { return Handle(^uintptr(1 - 1)) }







func GetCurrentThread() (Handle, error) {
	return CurrentThread(), nil
}



func CurrentThread() Handle { return Handle(^uintptr(2 - 1)) }



func GetProcAddressByOrdinal(module Handle, ordinal uintptr) (proc uintptr, err error) {
	r0, _, e1 := syscall.Syscall(procGetProcAddress.Addr(), 2, uintptr(module), ordinal, 0)
	proc = uintptr(r0)
	if proc == 0 {
		err = errnoErr(e1)
	}
	return
}

func Exit(code int) { ExitProcess(uint32(code)) }

func makeInheritSa() *SecurityAttributes {
	var sa SecurityAttributes
	sa.Length = uint32(unsafe.Sizeof(sa))
	sa.InheritHandle = 1
	return &sa
}

func Open(path string, mode int, perm uint32) (fd Handle, err error) {
	if len(path) == 0 {
		return InvalidHandle, ERROR_FILE_NOT_FOUND
	}
	pathp, err := UTF16PtrFromString(path)
	if err != nil {
		return InvalidHandle, err
	}
	var access uint32
	switch mode & (O_RDONLY | O_WRONLY | O_RDWR) {
	case O_RDONLY:
		access = GENERIC_READ
	case O_WRONLY:
		access = GENERIC_WRITE
	case O_RDWR:
		access = GENERIC_READ | GENERIC_WRITE
	}
	if mode&O_CREAT != 0 {
		access |= GENERIC_WRITE
	}
	if mode&O_APPEND != 0 {
		access &^= GENERIC_WRITE
		access |= FILE_APPEND_DATA
	}
	sharemode := uint32(FILE_SHARE_READ | FILE_SHARE_WRITE)
	var sa *SecurityAttributes
	if mode&O_CLOEXEC == 0 {
		sa = makeInheritSa()
	}
	var createmode uint32
	switch {
	case mode&(O_CREAT|O_EXCL) == (O_CREAT | O_EXCL):
		createmode = CREATE_NEW
	case mode&(O_CREAT|O_TRUNC) == (O_CREAT | O_TRUNC):
		createmode = CREATE_ALWAYS
	case mode&O_CREAT == O_CREAT:
		createmode = OPEN_ALWAYS
	case mode&O_TRUNC == O_TRUNC:
		createmode = TRUNCATE_EXISTING
	default:
		createmode = OPEN_EXISTING
	}
	var attrs uint32 = FILE_ATTRIBUTE_NORMAL
	if perm&S_IWRITE == 0 {
		attrs = FILE_ATTRIBUTE_READONLY
	}
	h, e := CreateFile(pathp, access, sharemode, sa, createmode, attrs, 0)
	return h, e
}

func Read(fd Handle, p []byte) (n int, err error) {
	var done uint32
	e := ReadFile(fd, p, &done, nil)
	if e != nil {
		if e == ERROR_BROKEN_PIPE {
			
			return 0, nil
		}
		return 0, e
	}
	return int(done), nil
}

func Write(fd Handle, p []byte) (n int, err error) {
	if raceenabled {
		raceReleaseMerge(unsafe.Pointer(&ioSync))
	}
	var done uint32
	e := WriteFile(fd, p, &done, nil)
	if e != nil {
		return 0, e
	}
	return int(done), nil
}

func ReadFile(fd Handle, p []byte, done *uint32, overlapped *Overlapped) error {
	err := readFile(fd, p, done, overlapped)
	if raceenabled {
		if *done > 0 {
			raceWriteRange(unsafe.Pointer(&p[0]), int(*done))
		}
		raceAcquire(unsafe.Pointer(&ioSync))
	}
	return err
}

func WriteFile(fd Handle, p []byte, done *uint32, overlapped *Overlapped) error {
	if raceenabled {
		raceReleaseMerge(unsafe.Pointer(&ioSync))
	}
	err := writeFile(fd, p, done, overlapped)
	if raceenabled && *done > 0 {
		raceReadRange(unsafe.Pointer(&p[0]), int(*done))
	}
	return err
}

var ioSync int64

func Seek(fd Handle, offset int64, whence int) (newoffset int64, err error) {
	var w uint32
	switch whence {
	case 0:
		w = FILE_BEGIN
	case 1:
		w = FILE_CURRENT
	case 2:
		w = FILE_END
	}
	hi := int32(offset >> 32)
	lo := int32(offset)
	
	ft, _ := GetFileType(fd)
	if ft == FILE_TYPE_PIPE {
		return 0, syscall.EPIPE
	}
	rlo, e := SetFilePointer(fd, lo, &hi, w)
	if e != nil {
		return 0, e
	}
	return int64(hi)<<32 + int64(rlo), nil
}

func Close(fd Handle) (err error) {
	return CloseHandle(fd)
}

var (
	Stdin  = getStdHandle(STD_INPUT_HANDLE)
	Stdout = getStdHandle(STD_OUTPUT_HANDLE)
	Stderr = getStdHandle(STD_ERROR_HANDLE)
)

func getStdHandle(stdhandle uint32) (fd Handle) {
	r, _ := GetStdHandle(stdhandle)
	return r
}

const ImplementsGetwd = true

func Getwd() (wd string, err error) {
	b := make([]uint16, 300)
	n, e := GetCurrentDirectory(uint32(len(b)), &b[0])
	if e != nil {
		return "", e
	}
	return string(utf16.Decode(b[0:n])), nil
}

func Chdir(path string) (err error) {
	pathp, err := UTF16PtrFromString(path)
	if err != nil {
		return err
	}
	return SetCurrentDirectory(pathp)
}

func Mkdir(path string, mode uint32) (err error) {
	pathp, err := UTF16PtrFromString(path)
	if err != nil {
		return err
	}
	return CreateDirectory(pathp, nil)
}

func Rmdir(path string) (err error) {
	pathp, err := UTF16PtrFromString(path)
	if err != nil {
		return err
	}
	return RemoveDirectory(pathp)
}

func Unlink(path string) (err error) {
	pathp, err := UTF16PtrFromString(path)
	if err != nil {
		return err
	}
	return DeleteFile(pathp)
}

func Rename(oldpath, newpath string) (err error) {
	from, err := UTF16PtrFromString(oldpath)
	if err != nil {
		return err
	}
	to, err := UTF16PtrFromString(newpath)
	if err != nil {
		return err
	}
	return MoveFileEx(from, to, MOVEFILE_REPLACE_EXISTING)
}

func ComputerName() (name string, err error) {
	var n uint32 = MAX_COMPUTERNAME_LENGTH + 1
	b := make([]uint16, n)
	e := GetComputerName(&b[0], &n)
	if e != nil {
		return "", e
	}
	return string(utf16.Decode(b[0:n])), nil
}

func DurationSinceBoot() time.Duration {
	return time.Duration(getTickCount64()) * time.Millisecond
}

func Ftruncate(fd Handle, length int64) (err error) {
	type _FILE_END_OF_FILE_INFO struct {
		EndOfFile int64
	}
	var info _FILE_END_OF_FILE_INFO
	info.EndOfFile = length
	return SetFileInformationByHandle(fd, FileEndOfFileInfo, (*byte)(unsafe.Pointer(&info)), uint32(unsafe.Sizeof(info)))
}

func Gettimeofday(tv *Timeval) (err error) {
	var ft Filetime
	GetSystemTimeAsFileTime(&ft)
	*tv = NsecToTimeval(ft.Nanoseconds())
	return nil
}

func Pipe(p []Handle) (err error) {
	if len(p) != 2 {
		return syscall.EINVAL
	}
	var r, w Handle
	e := CreatePipe(&r, &w, makeInheritSa(), 0)
	if e != nil {
		return e
	}
	p[0] = r
	p[1] = w
	return nil
}

func Utimes(path string, tv []Timeval) (err error) {
	if len(tv) != 2 {
		return syscall.EINVAL
	}
	pathp, e := UTF16PtrFromString(path)
	if e != nil {
		return e
	}
	h, e := CreateFile(pathp,
		FILE_WRITE_ATTRIBUTES, FILE_SHARE_WRITE, nil,
		OPEN_EXISTING, FILE_FLAG_BACKUP_SEMANTICS, 0)
	if e != nil {
		return e
	}
	defer CloseHandle(h)
	a := NsecToFiletime(tv[0].Nanoseconds())
	w := NsecToFiletime(tv[1].Nanoseconds())
	return SetFileTime(h, nil, &a, &w)
}

func UtimesNano(path string, ts []Timespec) (err error) {
	if len(ts) != 2 {
		return syscall.EINVAL
	}
	pathp, e := UTF16PtrFromString(path)
	if e != nil {
		return e
	}
	h, e := CreateFile(pathp,
		FILE_WRITE_ATTRIBUTES, FILE_SHARE_WRITE, nil,
		OPEN_EXISTING, FILE_FLAG_BACKUP_SEMANTICS, 0)
	if e != nil {
		return e
	}
	defer CloseHandle(h)
	a := NsecToFiletime(TimespecToNsec(ts[0]))
	w := NsecToFiletime(TimespecToNsec(ts[1]))
	return SetFileTime(h, nil, &a, &w)
}

func Fsync(fd Handle) (err error) {
	return FlushFileBuffers(fd)
}

func Chmod(path string, mode uint32) (err error) {
	p, e := UTF16PtrFromString(path)
	if e != nil {
		return e
	}
	attrs, e := GetFileAttributes(p)
	if e != nil {
		return e
	}
	if mode&S_IWRITE != 0 {
		attrs &^= FILE_ATTRIBUTE_READONLY
	} else {
		attrs |= FILE_ATTRIBUTE_READONLY
	}
	return SetFileAttributes(p, attrs)
}

func LoadGetSystemTimePreciseAsFileTime() error {
	return procGetSystemTimePreciseAsFileTime.Find()
}

func LoadCancelIoEx() error {
	return procCancelIoEx.Find()
}

func LoadSetFileCompletionNotificationModes() error {
	return procSetFileCompletionNotificationModes.Find()
}

func WaitForMultipleObjects(handles []Handle, waitAll bool, waitMilliseconds uint32) (event uint32, err error) {
	
	
	

	var handlePtr *Handle
	if len(handles) > 0 {
		handlePtr = &handles[0]
	}
	return waitForMultipleObjects(uint32(len(handles)), uintptr(unsafe.Pointer(handlePtr)), waitAll, waitMilliseconds)
}



const socket_error = uintptr(^uint32(0))




















































var SocketDisableIPv6 bool

type RawSockaddrInet4 struct {
	Family uint16
	Port   uint16
	Addr   [4]byte 
	Zero   [8]uint8
}

type RawSockaddrInet6 struct {
	Family   uint16
	Port     uint16
	Flowinfo uint32
	Addr     [16]byte 
	Scope_id uint32
}

type RawSockaddr struct {
	Family uint16
	Data   [14]int8
}

type RawSockaddrAny struct {
	Addr RawSockaddr
	Pad  [100]int8
}

type Sockaddr interface {
	sockaddr() (ptr unsafe.Pointer, len int32, err error) 
}

type SockaddrInet4 struct {
	Port int
	Addr [4]byte
	raw  RawSockaddrInet4
}

func (sa *SockaddrInet4) sockaddr() (unsafe.Pointer, int32, error) {
	if sa.Port < 0 || sa.Port > 0xFFFF {
		return nil, 0, syscall.EINVAL
	}
	sa.raw.Family = AF_INET
	p := (*[2]byte)(unsafe.Pointer(&sa.raw.Port))
	p[0] = byte(sa.Port >> 8)
	p[1] = byte(sa.Port)
	sa.raw.Addr = sa.Addr
	return unsafe.Pointer(&sa.raw), int32(unsafe.Sizeof(sa.raw)), nil
}

type SockaddrInet6 struct {
	Port   int
	ZoneId uint32
	Addr   [16]byte
	raw    RawSockaddrInet6
}

func (sa *SockaddrInet6) sockaddr() (unsafe.Pointer, int32, error) {
	if sa.Port < 0 || sa.Port > 0xFFFF {
		return nil, 0, syscall.EINVAL
	}
	sa.raw.Family = AF_INET6
	p := (*[2]byte)(unsafe.Pointer(&sa.raw.Port))
	p[0] = byte(sa.Port >> 8)
	p[1] = byte(sa.Port)
	sa.raw.Scope_id = sa.ZoneId
	sa.raw.Addr = sa.Addr
	return unsafe.Pointer(&sa.raw), int32(unsafe.Sizeof(sa.raw)), nil
}

type RawSockaddrUnix struct {
	Family uint16
	Path   [UNIX_PATH_MAX]int8
}

type SockaddrUnix struct {
	Name string
	raw  RawSockaddrUnix
}

func (sa *SockaddrUnix) sockaddr() (unsafe.Pointer, int32, error) {
	name := sa.Name
	n := len(name)
	if n > len(sa.raw.Path) {
		return nil, 0, syscall.EINVAL
	}
	if n == len(sa.raw.Path) && name[0] != '@' {
		return nil, 0, syscall.EINVAL
	}
	sa.raw.Family = AF_UNIX
	for i := 0; i < n; i++ {
		sa.raw.Path[i] = int8(name[i])
	}
	
	sl := int32(2)
	if n > 0 {
		sl += int32(n) + 1
	}
	if sa.raw.Path[0] == '@' || (sa.raw.Path[0] == 0 && sl > 3) {
		
		sa.raw.Path[0] = 0
		
		sl--
	}

	return unsafe.Pointer(&sa.raw), sl, nil
}

type RawSockaddrBth struct {
	AddressFamily  [2]byte
	BtAddr         [8]byte
	ServiceClassId [16]byte
	Port           [4]byte
}

type SockaddrBth struct {
	BtAddr         uint64
	ServiceClassId GUID
	Port           uint32

	raw RawSockaddrBth
}

func (sa *SockaddrBth) sockaddr() (unsafe.Pointer, int32, error) {
	family := AF_BTH
	sa.raw = RawSockaddrBth{
		AddressFamily:  *(*[2]byte)(unsafe.Pointer(&family)),
		BtAddr:         *(*[8]byte)(unsafe.Pointer(&sa.BtAddr)),
		Port:           *(*[4]byte)(unsafe.Pointer(&sa.Port)),
		ServiceClassId: *(*[16]byte)(unsafe.Pointer(&sa.ServiceClassId)),
	}
	return unsafe.Pointer(&sa.raw), int32(unsafe.Sizeof(sa.raw)), nil
}

func (rsa *RawSockaddrAny) Sockaddr() (Sockaddr, error) {
	switch rsa.Addr.Family {
	case AF_UNIX:
		pp := (*RawSockaddrUnix)(unsafe.Pointer(rsa))
		sa := new(SockaddrUnix)
		if pp.Path[0] == 0 {
			
			
			
			
			
			pp.Path[0] = '@'
		}

		
		
		
		
		
		n := 0
		for n < len(pp.Path) && pp.Path[n] != 0 {
			n++
		}
		sa.Name = string(unsafe.Slice((*byte)(unsafe.Pointer(&pp.Path[0])), n))
		return sa, nil

	case AF_INET:
		pp := (*RawSockaddrInet4)(unsafe.Pointer(rsa))
		sa := new(SockaddrInet4)
		p := (*[2]byte)(unsafe.Pointer(&pp.Port))
		sa.Port = int(p[0])<<8 + int(p[1])
		sa.Addr = pp.Addr
		return sa, nil

	case AF_INET6:
		pp := (*RawSockaddrInet6)(unsafe.Pointer(rsa))
		sa := new(SockaddrInet6)
		p := (*[2]byte)(unsafe.Pointer(&pp.Port))
		sa.Port = int(p[0])<<8 + int(p[1])
		sa.ZoneId = pp.Scope_id
		sa.Addr = pp.Addr
		return sa, nil
	}
	return nil, syscall.EAFNOSUPPORT
}

func Socket(domain, typ, proto int) (fd Handle, err error) {
	if domain == AF_INET6 && SocketDisableIPv6 {
		return InvalidHandle, syscall.EAFNOSUPPORT
	}
	return socket(int32(domain), int32(typ), int32(proto))
}

func SetsockoptInt(fd Handle, level, opt int, value int) (err error) {
	v := int32(value)
	return Setsockopt(fd, int32(level), int32(opt), (*byte)(unsafe.Pointer(&v)), int32(unsafe.Sizeof(v)))
}

func Bind(fd Handle, sa Sockaddr) (err error) {
	ptr, n, err := sa.sockaddr()
	if err != nil {
		return err
	}
	return bind(fd, ptr, n)
}

func Connect(fd Handle, sa Sockaddr) (err error) {
	ptr, n, err := sa.sockaddr()
	if err != nil {
		return err
	}
	return connect(fd, ptr, n)
}

func GetBestInterfaceEx(sa Sockaddr, pdwBestIfIndex *uint32) (err error) {
	ptr, _, err := sa.sockaddr()
	if err != nil {
		return err
	}
	return getBestInterfaceEx(ptr, pdwBestIfIndex)
}

func Getsockname(fd Handle) (sa Sockaddr, err error) {
	var rsa RawSockaddrAny
	l := int32(unsafe.Sizeof(rsa))
	if err = getsockname(fd, &rsa, &l); err != nil {
		return
	}
	return rsa.Sockaddr()
}

func Getpeername(fd Handle) (sa Sockaddr, err error) {
	var rsa RawSockaddrAny
	l := int32(unsafe.Sizeof(rsa))
	if err = getpeername(fd, &rsa, &l); err != nil {
		return
	}
	return rsa.Sockaddr()
}

func Listen(s Handle, n int) (err error) {
	return listen(s, int32(n))
}

func Shutdown(fd Handle, how int) (err error) {
	return shutdown(fd, int32(how))
}

func WSASendto(s Handle, bufs *WSABuf, bufcnt uint32, sent *uint32, flags uint32, to Sockaddr, overlapped *Overlapped, croutine *byte) (err error) {
	var rsa unsafe.Pointer
	var l int32
	if to != nil {
		rsa, l, err = to.sockaddr()
		if err != nil {
			return err
		}
	}
	return WSASendTo(s, bufs, bufcnt, sent, flags, (*RawSockaddrAny)(unsafe.Pointer(rsa)), l, overlapped, croutine)
}

func LoadGetAddrInfo() error {
	return procGetAddrInfoW.Find()
}

var connectExFunc struct {
	once sync.Once
	addr uintptr
	err  error
}

func LoadConnectEx() error {
	connectExFunc.once.Do(func() {
		var s Handle
		s, connectExFunc.err = Socket(AF_INET, SOCK_STREAM, IPPROTO_TCP)
		if connectExFunc.err != nil {
			return
		}
		defer CloseHandle(s)
		var n uint32
		connectExFunc.err = WSAIoctl(s,
			SIO_GET_EXTENSION_FUNCTION_POINTER,
			(*byte)(unsafe.Pointer(&WSAID_CONNECTEX)),
			uint32(unsafe.Sizeof(WSAID_CONNECTEX)),
			(*byte)(unsafe.Pointer(&connectExFunc.addr)),
			uint32(unsafe.Sizeof(connectExFunc.addr)),
			&n, nil, 0)
	})
	return connectExFunc.err
}

func connectEx(s Handle, name unsafe.Pointer, namelen int32, sendBuf *byte, sendDataLen uint32, bytesSent *uint32, overlapped *Overlapped) (err error) {
	r1, _, e1 := syscall.Syscall9(connectExFunc.addr, 7, uintptr(s), uintptr(name), uintptr(namelen), uintptr(unsafe.Pointer(sendBuf)), uintptr(sendDataLen), uintptr(unsafe.Pointer(bytesSent)), uintptr(unsafe.Pointer(overlapped)), 0, 0)
	if r1 == 0 {
		if e1 != 0 {
			err = error(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return
}

func ConnectEx(fd Handle, sa Sockaddr, sendBuf *byte, sendDataLen uint32, bytesSent *uint32, overlapped *Overlapped) error {
	err := LoadConnectEx()
	if err != nil {
		return errorspkg.New("failed to find ConnectEx: " + err.Error())
	}
	ptr, n, err := sa.sockaddr()
	if err != nil {
		return err
	}
	return connectEx(fd, ptr, n, sendBuf, sendDataLen, bytesSent, overlapped)
}

var sendRecvMsgFunc struct {
	once     sync.Once
	sendAddr uintptr
	recvAddr uintptr
	err      error
}

func loadWSASendRecvMsg() error {
	sendRecvMsgFunc.once.Do(func() {
		var s Handle
		s, sendRecvMsgFunc.err = Socket(AF_INET, SOCK_DGRAM, IPPROTO_UDP)
		if sendRecvMsgFunc.err != nil {
			return
		}
		defer CloseHandle(s)
		var n uint32
		sendRecvMsgFunc.err = WSAIoctl(s,
			SIO_GET_EXTENSION_FUNCTION_POINTER,
			(*byte)(unsafe.Pointer(&WSAID_WSARECVMSG)),
			uint32(unsafe.Sizeof(WSAID_WSARECVMSG)),
			(*byte)(unsafe.Pointer(&sendRecvMsgFunc.recvAddr)),
			uint32(unsafe.Sizeof(sendRecvMsgFunc.recvAddr)),
			&n, nil, 0)
		if sendRecvMsgFunc.err != nil {
			return
		}
		sendRecvMsgFunc.err = WSAIoctl(s,
			SIO_GET_EXTENSION_FUNCTION_POINTER,
			(*byte)(unsafe.Pointer(&WSAID_WSASENDMSG)),
			uint32(unsafe.Sizeof(WSAID_WSASENDMSG)),
			(*byte)(unsafe.Pointer(&sendRecvMsgFunc.sendAddr)),
			uint32(unsafe.Sizeof(sendRecvMsgFunc.sendAddr)),
			&n, nil, 0)
	})
	return sendRecvMsgFunc.err
}

func WSASendMsg(fd Handle, msg *WSAMsg, flags uint32, bytesSent *uint32, overlapped *Overlapped, croutine *byte) error {
	err := loadWSASendRecvMsg()
	if err != nil {
		return err
	}
	r1, _, e1 := syscall.Syscall6(sendRecvMsgFunc.sendAddr, 6, uintptr(fd), uintptr(unsafe.Pointer(msg)), uintptr(flags), uintptr(unsafe.Pointer(bytesSent)), uintptr(unsafe.Pointer(overlapped)), uintptr(unsafe.Pointer(croutine)))
	if r1 == socket_error {
		err = errnoErr(e1)
	}
	return err
}

func WSARecvMsg(fd Handle, msg *WSAMsg, bytesReceived *uint32, overlapped *Overlapped, croutine *byte) error {
	err := loadWSASendRecvMsg()
	if err != nil {
		return err
	}
	r1, _, e1 := syscall.Syscall6(sendRecvMsgFunc.recvAddr, 5, uintptr(fd), uintptr(unsafe.Pointer(msg)), uintptr(unsafe.Pointer(bytesReceived)), uintptr(unsafe.Pointer(overlapped)), uintptr(unsafe.Pointer(croutine)), 0)
	if r1 == socket_error {
		err = errnoErr(e1)
	}
	return err
}


type Rusage struct {
	CreationTime Filetime
	ExitTime     Filetime
	KernelTime   Filetime
	UserTime     Filetime
}

type WaitStatus struct {
	ExitCode uint32
}

func (w WaitStatus) Exited() bool { return true }

func (w WaitStatus) ExitStatus() int { return int(w.ExitCode) }

func (w WaitStatus) Signal() Signal { return -1 }

func (w WaitStatus) CoreDump() bool { return false }

func (w WaitStatus) Stopped() bool { return false }

func (w WaitStatus) Continued() bool { return false }

func (w WaitStatus) StopSignal() Signal { return -1 }

func (w WaitStatus) Signaled() bool { return false }

func (w WaitStatus) TrapCause() int { return -1 }



type Timespec struct {
	Sec  int64
	Nsec int64
}

func TimespecToNsec(ts Timespec) int64 { return int64(ts.Sec)*1e9 + int64(ts.Nsec) }

func NsecToTimespec(nsec int64) (ts Timespec) {
	ts.Sec = nsec / 1e9
	ts.Nsec = nsec % 1e9
	return
}



func Accept(fd Handle) (nfd Handle, sa Sockaddr, err error) { return 0, nil, syscall.EWINDOWS }

func Recvfrom(fd Handle, p []byte, flags int) (n int, from Sockaddr, err error) {
	var rsa RawSockaddrAny
	l := int32(unsafe.Sizeof(rsa))
	n32, err := recvfrom(fd, p, int32(flags), &rsa, &l)
	n = int(n32)
	if err != nil {
		return
	}
	from, err = rsa.Sockaddr()
	return
}

func Sendto(fd Handle, p []byte, flags int, to Sockaddr) (err error) {
	ptr, l, err := to.sockaddr()
	if err != nil {
		return err
	}
	return sendto(fd, p, int32(flags), ptr, l)
}

func SetsockoptTimeval(fd Handle, level, opt int, tv *Timeval) (err error) { return syscall.EWINDOWS }








type Linger struct {
	Onoff  int32
	Linger int32
}

type sysLinger struct {
	Onoff  uint16
	Linger uint16
}

type IPMreq struct {
	Multiaddr [4]byte 
	Interface [4]byte 
}

type IPv6Mreq struct {
	Multiaddr [16]byte 
	Interface uint32
}

func GetsockoptInt(fd Handle, level, opt int) (int, error) {
	v := int32(0)
	l := int32(unsafe.Sizeof(v))
	err := Getsockopt(fd, int32(level), int32(opt), (*byte)(unsafe.Pointer(&v)), &l)
	return int(v), err
}

func SetsockoptLinger(fd Handle, level, opt int, l *Linger) (err error) {
	sys := sysLinger{Onoff: uint16(l.Onoff), Linger: uint16(l.Linger)}
	return Setsockopt(fd, int32(level), int32(opt), (*byte)(unsafe.Pointer(&sys)), int32(unsafe.Sizeof(sys)))
}

func SetsockoptInet4Addr(fd Handle, level, opt int, value [4]byte) (err error) {
	return Setsockopt(fd, int32(level), int32(opt), (*byte)(unsafe.Pointer(&value[0])), 4)
}

func SetsockoptIPMreq(fd Handle, level, opt int, mreq *IPMreq) (err error) {
	return Setsockopt(fd, int32(level), int32(opt), (*byte)(unsafe.Pointer(mreq)), int32(unsafe.Sizeof(*mreq)))
}

func SetsockoptIPv6Mreq(fd Handle, level, opt int, mreq *IPv6Mreq) (err error) {
	return syscall.EWINDOWS
}

func EnumProcesses(processIds []uint32, bytesReturned *uint32) error {
	
	
	var p *uint32
	if len(processIds) > 0 {
		p = &processIds[0]
	}
	size := uint32(len(processIds) * 4)
	return enumProcesses(p, size, bytesReturned)
}

func Getpid() (pid int) { return int(GetCurrentProcessId()) }

func FindFirstFile(name *uint16, data *Win32finddata) (handle Handle, err error) {
	
	
	
	
	
	
	
	
	var data1 win32finddata1
	handle, err = findFirstFile1(name, &data1)
	if err == nil {
		copyFindData(data, &data1)
	}
	return
}

func FindNextFile(handle Handle, data *Win32finddata) (err error) {
	var data1 win32finddata1
	err = findNextFile1(handle, &data1)
	if err == nil {
		copyFindData(data, &data1)
	}
	return
}

func getProcessEntry(pid int) (*ProcessEntry32, error) {
	snapshot, err := CreateToolhelp32Snapshot(TH32CS_SNAPPROCESS, 0)
	if err != nil {
		return nil, err
	}
	defer CloseHandle(snapshot)
	var procEntry ProcessEntry32
	procEntry.Size = uint32(unsafe.Sizeof(procEntry))
	if err = Process32First(snapshot, &procEntry); err != nil {
		return nil, err
	}
	for {
		if procEntry.ProcessID == uint32(pid) {
			return &procEntry, nil
		}
		err = Process32Next(snapshot, &procEntry)
		if err != nil {
			return nil, err
		}
	}
}

func Getppid() (ppid int) {
	pe, err := getProcessEntry(Getpid())
	if err != nil {
		return -1
	}
	return int(pe.ParentProcessID)
}


func Fchdir(fd Handle) (err error)             { return syscall.EWINDOWS }
func Link(oldpath, newpath string) (err error) { return syscall.EWINDOWS }
func Symlink(path, link string) (err error)    { return syscall.EWINDOWS }

func Fchmod(fd Handle, mode uint32) (err error)        { return syscall.EWINDOWS }
func Chown(path string, uid int, gid int) (err error)  { return syscall.EWINDOWS }
func Lchown(path string, uid int, gid int) (err error) { return syscall.EWINDOWS }
func Fchown(fd Handle, uid int, gid int) (err error)   { return syscall.EWINDOWS }

func Getuid() (uid int)                  { return -1 }
func Geteuid() (euid int)                { return -1 }
func Getgid() (gid int)                  { return -1 }
func Getegid() (egid int)                { return -1 }
func Getgroups() (gids []int, err error) { return nil, syscall.EWINDOWS }

type Signal int

func (s Signal) Signal() {}

func (s Signal) String() string {
	if 0 <= s && int(s) < len(signals) {
		str := signals[s]
		if str != "" {
			return str
		}
	}
	return "signal " + itoa(int(s))
}

func LoadCreateSymbolicLink() error {
	return procCreateSymbolicLinkW.Find()
}


func Readlink(path string, buf []byte) (n int, err error) {
	fd, err := CreateFile(StringToUTF16Ptr(path), GENERIC_READ, 0, nil, OPEN_EXISTING,
		FILE_FLAG_OPEN_REPARSE_POINT|FILE_FLAG_BACKUP_SEMANTICS, 0)
	if err != nil {
		return -1, err
	}
	defer CloseHandle(fd)

	rdbbuf := make([]byte, MAXIMUM_REPARSE_DATA_BUFFER_SIZE)
	var bytesReturned uint32
	err = DeviceIoControl(fd, FSCTL_GET_REPARSE_POINT, nil, 0, &rdbbuf[0], uint32(len(rdbbuf)), &bytesReturned, nil)
	if err != nil {
		return -1, err
	}

	rdb := (*reparseDataBuffer)(unsafe.Pointer(&rdbbuf[0]))
	var s string
	switch rdb.ReparseTag {
	case IO_REPARSE_TAG_SYMLINK:
		data := (*symbolicLinkReparseBuffer)(unsafe.Pointer(&rdb.reparseBuffer))
		p := (*[0xffff]uint16)(unsafe.Pointer(&data.PathBuffer[0]))
		s = UTF16ToString(p[data.PrintNameOffset/2 : (data.PrintNameLength-data.PrintNameOffset)/2])
	case IO_REPARSE_TAG_MOUNT_POINT:
		data := (*mountPointReparseBuffer)(unsafe.Pointer(&rdb.reparseBuffer))
		p := (*[0xffff]uint16)(unsafe.Pointer(&data.PathBuffer[0]))
		s = UTF16ToString(p[data.PrintNameOffset/2 : (data.PrintNameLength-data.PrintNameOffset)/2])
	default:
		
		
		return -1, syscall.ENOENT
	}
	n = copy(buf, []byte(s))

	return n, nil
}



func GUIDFromString(str string) (GUID, error) {
	guid := GUID{}
	str16, err := syscall.UTF16PtrFromString(str)
	if err != nil {
		return guid, err
	}
	err = clsidFromString(str16, &guid)
	if err != nil {
		return guid, err
	}
	return guid, nil
}


func GenerateGUID() (GUID, error) {
	guid := GUID{}
	err := coCreateGuid(&guid)
	if err != nil {
		return guid, err
	}
	return guid, nil
}



func (guid GUID) String() string {
	var str [100]uint16
	chars := stringFromGUID2(&guid, &str[0], int32(len(str)))
	if chars <= 1 {
		return ""
	}
	return string(utf16.Decode(str[:chars-1]))
}



func KnownFolderPath(folderID *KNOWNFOLDERID, flags uint32) (string, error) {
	return Token(0).KnownFolderPath(folderID, flags)
}



func (t Token) KnownFolderPath(folderID *KNOWNFOLDERID, flags uint32) (string, error) {
	var p *uint16
	err := shGetKnownFolderPath(folderID, flags, t, &p)
	if err != nil {
		return "", err
	}
	defer CoTaskMemFree(unsafe.Pointer(p))
	return UTF16PtrToString(p), nil
}



func RtlGetVersion() *OsVersionInfoEx {
	info := &OsVersionInfoEx{}
	info.osVersionInfoSize = uint32(unsafe.Sizeof(*info))
	
	
	
	
	_ = rtlGetVersion(info)
	return info
}



func RtlGetNtVersionNumbers() (majorVersion, minorVersion, buildNumber uint32) {
	rtlGetNtVersionNumbers(&majorVersion, &minorVersion, &buildNumber)
	buildNumber &= 0xffff
	return
}


func GetProcessPreferredUILanguages(flags uint32) ([]string, error) {
	return getUILanguages(flags, getProcessPreferredUILanguages)
}


func GetThreadPreferredUILanguages(flags uint32) ([]string, error) {
	return getUILanguages(flags, getThreadPreferredUILanguages)
}


func GetUserPreferredUILanguages(flags uint32) ([]string, error) {
	return getUILanguages(flags, getUserPreferredUILanguages)
}


func GetSystemPreferredUILanguages(flags uint32) ([]string, error) {
	return getUILanguages(flags, getSystemPreferredUILanguages)
}

func getUILanguages(flags uint32, f func(flags uint32, numLanguages *uint32, buf *uint16, bufSize *uint32) error) ([]string, error) {
	size := uint32(128)
	for {
		var numLanguages uint32
		buf := make([]uint16, size)
		err := f(flags, &numLanguages, &buf[0], &size)
		if err == ERROR_INSUFFICIENT_BUFFER {
			continue
		}
		if err != nil {
			return nil, err
		}
		buf = buf[:size]
		if numLanguages == 0 || len(buf) == 0 { 
			return []string{}, nil
		}
		if buf[len(buf)-1] == 0 {
			buf = buf[:len(buf)-1] 
		}
		languages := make([]string, 0, numLanguages)
		from := 0
		for i, c := range buf {
			if c == 0 {
				languages = append(languages, string(utf16.Decode(buf[from:i])))
				from = i + 1
			}
		}
		return languages, nil
	}
}

func SetConsoleCursorPosition(console Handle, position Coord) error {
	return setConsoleCursorPosition(console, *((*uint32)(unsafe.Pointer(&position))))
}

func GetStartupInfo(startupInfo *StartupInfo) error {
	getStartupInfo(startupInfo)
	return nil
}

func (s NTStatus) Errno() syscall.Errno {
	return rtlNtStatusToDosErrorNoTeb(s)
}

func langID(pri, sub uint16) uint32 { return uint32(sub)<<10 | uint32(pri) }

func (s NTStatus) Error() string {
	b := make([]uint16, 300)
	n, err := FormatMessage(FORMAT_MESSAGE_FROM_SYSTEM|FORMAT_MESSAGE_FROM_HMODULE|FORMAT_MESSAGE_ARGUMENT_ARRAY, modntdll.Handle(), uint32(s), langID(LANG_ENGLISH, SUBLANG_ENGLISH_US), b, nil)
	if err != nil {
		return fmt.Sprintf("NTSTATUS 0x%08x", uint32(s))
	}
	
	for ; n > 0 && (b[n-1] == '\n' || b[n-1] == '\r'); n-- {
	}
	return string(utf16.Decode(b[:n]))
}





func NewNTUnicodeString(s string) (*NTUnicodeString, error) {
	s16, err := UTF16FromString(s)
	if err != nil {
		return nil, err
	}
	n := uint16(len(s16) * 2)
	return &NTUnicodeString{
		Length:        n - 2, 
		MaximumLength: n,
		Buffer:        &s16[0],
	}, nil
}


func (s *NTUnicodeString) Slice() []uint16 {
	slice := unsafe.Slice(s.Buffer, s.MaximumLength)
	return slice[:s.Length]
}

func (s *NTUnicodeString) String() string {
	return UTF16ToString(s.Slice())
}





func NewNTString(s string) (*NTString, error) {
	var nts NTString
	s8, err := BytePtrFromString(s)
	if err != nil {
		return nil, err
	}
	RtlInitString(&nts, s8)
	return &nts, nil
}


func (s *NTString) Slice() []byte {
	slice := unsafe.Slice(s.Buffer, s.MaximumLength)
	return slice[:s.Length]
}

func (s *NTString) String() string {
	return ByteSliceToString(s.Slice())
}


func FindResource(module Handle, name, resType ResourceIDOrString) (Handle, error) {
	var namePtr, resTypePtr uintptr
	var name16, resType16 *uint16
	var err error
	resolvePtr := func(i interface{}, keep **uint16) (uintptr, error) {
		switch v := i.(type) {
		case string:
			*keep, err = UTF16PtrFromString(v)
			if err != nil {
				return 0, err
			}
			return uintptr(unsafe.Pointer(*keep)), nil
		case ResourceID:
			return uintptr(v), nil
		}
		return 0, errorspkg.New("parameter must be a ResourceID or a string")
	}
	namePtr, err = resolvePtr(name, &name16)
	if err != nil {
		return 0, err
	}
	resTypePtr, err = resolvePtr(resType, &resType16)
	if err != nil {
		return 0, err
	}
	resInfo, err := findResource(module, namePtr, resTypePtr)
	runtime.KeepAlive(name16)
	runtime.KeepAlive(resType16)
	return resInfo, err
}

func LoadResourceData(module, resInfo Handle) (data []byte, err error) {
	size, err := SizeofResource(module, resInfo)
	if err != nil {
		return
	}
	resData, err := LoadResource(module, resInfo)
	if err != nil {
		return
	}
	ptr, err := LockResource(resData)
	if err != nil {
		return
	}
	data = unsafe.Slice((*byte)(unsafe.Pointer(ptr)), size)
	return
}


type PSAPI_WORKING_SET_EX_BLOCK uint64



func (b PSAPI_WORKING_SET_EX_BLOCK) Valid() bool {
	return (b & 1) == 1
}


func (b PSAPI_WORKING_SET_EX_BLOCK) ShareCount() uint64 {
	return b.intField(1, 3)
}



func (b PSAPI_WORKING_SET_EX_BLOCK) Win32Protection() uint64 {
	return b.intField(4, 11)
}



func (b PSAPI_WORKING_SET_EX_BLOCK) Shared() bool {
	return (b & (1 << 15)) == 1
}


func (b PSAPI_WORKING_SET_EX_BLOCK) Node() uint64 {
	return b.intField(16, 6)
}



func (b PSAPI_WORKING_SET_EX_BLOCK) Locked() bool {
	return (b & (1 << 22)) == 1
}



func (b PSAPI_WORKING_SET_EX_BLOCK) LargePage() bool {
	return (b & (1 << 23)) == 1
}



func (b PSAPI_WORKING_SET_EX_BLOCK) Bad() bool {
	return (b & (1 << 31)) == 1
}


func (b PSAPI_WORKING_SET_EX_BLOCK) intField(start, length int) uint64 {
	var mask PSAPI_WORKING_SET_EX_BLOCK
	for pos := start; pos < start+length; pos++ {
		mask |= (1 << pos)
	}

	masked := b & mask
	return uint64(masked >> start)
}


type PSAPI_WORKING_SET_EX_INFORMATION struct {
	
	VirtualAddress Pointer
	
	VirtualAttributes PSAPI_WORKING_SET_EX_BLOCK
}


func CreatePseudoConsole(size Coord, in Handle, out Handle, flags uint32, pconsole *Handle) error {
	
	
	return createPseudoConsole(*((*uint32)(unsafe.Pointer(&size))), in, out, flags, pconsole)
}


func ResizePseudoConsole(pconsole Handle, size Coord) error {
	
	
	return resizePseudoConsole(pconsole, *((*uint32)(unsafe.Pointer(&size))))
}


const (
	CBR_110    = 110
	CBR_300    = 300
	CBR_600    = 600
	CBR_1200   = 1200
	CBR_2400   = 2400
	CBR_4800   = 4800
	CBR_9600   = 9600
	CBR_14400  = 14400
	CBR_19200  = 19200
	CBR_38400  = 38400
	CBR_57600  = 57600
	CBR_115200 = 115200
	CBR_128000 = 128000
	CBR_256000 = 256000

	DTR_CONTROL_DISABLE   = 0x00000000
	DTR_CONTROL_ENABLE    = 0x00000010
	DTR_CONTROL_HANDSHAKE = 0x00000020

	RTS_CONTROL_DISABLE   = 0x00000000
	RTS_CONTROL_ENABLE    = 0x00001000
	RTS_CONTROL_HANDSHAKE = 0x00002000
	RTS_CONTROL_TOGGLE    = 0x00003000

	NOPARITY    = 0
	ODDPARITY   = 1
	EVENPARITY  = 2
	MARKPARITY  = 3
	SPACEPARITY = 4

	ONESTOPBIT   = 0
	ONE5STOPBITS = 1
	TWOSTOPBITS  = 2
)


const (
	SETXOFF  = 1
	SETXON   = 2
	SETRTS   = 3
	CLRRTS   = 4
	SETDTR   = 5
	CLRDTR   = 6
	SETBREAK = 8
	CLRBREAK = 9
)


const (
	PURGE_TXABORT = 0x0001
	PURGE_RXABORT = 0x0002
	PURGE_TXCLEAR = 0x0004
	PURGE_RXCLEAR = 0x0008
)


const (
	EV_RXCHAR  = 0x0001
	EV_RXFLAG  = 0x0002
	EV_TXEMPTY = 0x0004
	EV_CTS     = 0x0008
	EV_DSR     = 0x0010
	EV_RLSD    = 0x0020
	EV_BREAK   = 0x0040
	EV_ERR     = 0x0080
	EV_RING    = 0x0100
)
