



package unix

import "unsafe"




func IoctlRetInt(fd int, req uint) (int, error) {
	ret, _, err := Syscall(SYS_IOCTL, uintptr(fd), uintptr(req), 0)
	if err != 0 {
		return 0, err
	}
	return int(ret), nil
}

func IoctlGetUint32(fd int, req uint) (uint32, error) {
	var value uint32
	err := ioctlPtr(fd, req, unsafe.Pointer(&value))
	return value, err
}

func IoctlGetRTCTime(fd int) (*RTCTime, error) {
	var value RTCTime
	err := ioctlPtr(fd, RTC_RD_TIME, unsafe.Pointer(&value))
	return &value, err
}

func IoctlSetRTCTime(fd int, value *RTCTime) error {
	return ioctlPtr(fd, RTC_SET_TIME, unsafe.Pointer(value))
}

func IoctlGetRTCWkAlrm(fd int) (*RTCWkAlrm, error) {
	var value RTCWkAlrm
	err := ioctlPtr(fd, RTC_WKALM_RD, unsafe.Pointer(&value))
	return &value, err
}

func IoctlSetRTCWkAlrm(fd int, value *RTCWkAlrm) error {
	return ioctlPtr(fd, RTC_WKALM_SET, unsafe.Pointer(value))
}



func IoctlGetEthtoolDrvinfo(fd int, ifname string) (*EthtoolDrvinfo, error) {
	ifr, err := NewIfreq(ifname)
	if err != nil {
		return nil, err
	}

	value := EthtoolDrvinfo{Cmd: ETHTOOL_GDRVINFO}
	ifrd := ifr.withData(unsafe.Pointer(&value))

	err = ioctlIfreqData(fd, SIOCETHTOOL, &ifrd)
	return &value, err
}



func IoctlGetEthtoolTsInfo(fd int, ifname string) (*EthtoolTsInfo, error) {
	ifr, err := NewIfreq(ifname)
	if err != nil {
		return nil, err
	}

	value := EthtoolTsInfo{Cmd: ETHTOOL_GET_TS_INFO}
	ifrd := ifr.withData(unsafe.Pointer(&value))

	err = ioctlIfreqData(fd, SIOCETHTOOL, &ifrd)
	return &value, err
}



func IoctlGetHwTstamp(fd int, ifname string) (*HwTstampConfig, error) {
	ifr, err := NewIfreq(ifname)
	if err != nil {
		return nil, err
	}

	value := HwTstampConfig{}
	ifrd := ifr.withData(unsafe.Pointer(&value))

	err = ioctlIfreqData(fd, SIOCGHWTSTAMP, &ifrd)
	return &value, err
}



func IoctlSetHwTstamp(fd int, ifname string, cfg *HwTstampConfig) error {
	ifr, err := NewIfreq(ifname)
	if err != nil {
		return err
	}
	ifrd := ifr.withData(unsafe.Pointer(cfg))
	return ioctlIfreqData(fd, SIOCSHWTSTAMP, &ifrd)
}




func FdToClockID(fd int) int32 { return int32((int(^fd) << 3) | 3) }


func IoctlPtpClockGetcaps(fd int) (*PtpClockCaps, error) {
	var value PtpClockCaps
	err := ioctlPtr(fd, PTP_CLOCK_GETCAPS2, unsafe.Pointer(&value))
	return &value, err
}



func IoctlPtpSysOffsetPrecise(fd int) (*PtpSysOffsetPrecise, error) {
	var value PtpSysOffsetPrecise
	err := ioctlPtr(fd, PTP_SYS_OFFSET_PRECISE2, unsafe.Pointer(&value))
	return &value, err
}




func IoctlPtpSysOffsetExtended(fd int, samples uint) (*PtpSysOffsetExtended, error) {
	value := PtpSysOffsetExtended{Samples: uint32(samples)}
	err := ioctlPtr(fd, PTP_SYS_OFFSET_EXTENDED2, unsafe.Pointer(&value))
	return &value, err
}



func IoctlPtpPinGetfunc(fd int, index uint) (*PtpPinDesc, error) {
	value := PtpPinDesc{Index: uint32(index)}
	err := ioctlPtr(fd, PTP_PIN_GETFUNC2, unsafe.Pointer(&value))
	return &value, err
}



func IoctlPtpPinSetfunc(fd int, pd *PtpPinDesc) error {
	return ioctlPtr(fd, PTP_PIN_SETFUNC2, unsafe.Pointer(pd))
}



func IoctlPtpPeroutRequest(fd int, r *PtpPeroutRequest) error {
	return ioctlPtr(fd, PTP_PEROUT_REQUEST2, unsafe.Pointer(r))
}



func IoctlPtpExttsRequest(fd int, r *PtpExttsRequest) error {
	return ioctlPtr(fd, PTP_EXTTS_REQUEST2, unsafe.Pointer(r))
}




func IoctlGetWatchdogInfo(fd int) (*WatchdogInfo, error) {
	var value WatchdogInfo
	err := ioctlPtr(fd, WDIOC_GETSUPPORT, unsafe.Pointer(&value))
	return &value, err
}




func IoctlWatchdogKeepalive(fd int) error {
	
	return ioctl(fd, WDIOC_KEEPALIVE, 0)
}




func IoctlFileCloneRange(destFd int, value *FileCloneRange) error {
	return ioctlPtr(destFd, FICLONERANGE, unsafe.Pointer(value))
}




func IoctlFileClone(destFd, srcFd int) error {
	return ioctl(destFd, FICLONE, uintptr(srcFd))
}

type FileDedupeRange struct {
	Src_offset uint64
	Src_length uint64
	Reserved1  uint16
	Reserved2  uint32
	Info       []FileDedupeRangeInfo
}

type FileDedupeRangeInfo struct {
	Dest_fd       int64
	Dest_offset   uint64
	Bytes_deduped uint64
	Status        int32
	Reserved      uint32
}





func IoctlFileDedupeRange(srcFd int, value *FileDedupeRange) error {
	buf := make([]byte, SizeofRawFileDedupeRange+
		len(value.Info)*SizeofRawFileDedupeRangeInfo)
	rawrange := (*RawFileDedupeRange)(unsafe.Pointer(&buf[0]))
	rawrange.Src_offset = value.Src_offset
	rawrange.Src_length = value.Src_length
	rawrange.Dest_count = uint16(len(value.Info))
	rawrange.Reserved1 = value.Reserved1
	rawrange.Reserved2 = value.Reserved2

	for i := range value.Info {
		rawinfo := (*RawFileDedupeRangeInfo)(unsafe.Pointer(
			uintptr(unsafe.Pointer(&buf[0])) + uintptr(SizeofRawFileDedupeRange) +
				uintptr(i*SizeofRawFileDedupeRangeInfo)))
		rawinfo.Dest_fd = value.Info[i].Dest_fd
		rawinfo.Dest_offset = value.Info[i].Dest_offset
		rawinfo.Bytes_deduped = value.Info[i].Bytes_deduped
		rawinfo.Status = value.Info[i].Status
		rawinfo.Reserved = value.Info[i].Reserved
	}

	err := ioctlPtr(srcFd, FIDEDUPERANGE, unsafe.Pointer(&buf[0]))

	
	for i := range value.Info {
		rawinfo := (*RawFileDedupeRangeInfo)(unsafe.Pointer(
			uintptr(unsafe.Pointer(&buf[0])) + uintptr(SizeofRawFileDedupeRange) +
				uintptr(i*SizeofRawFileDedupeRangeInfo)))
		value.Info[i].Dest_fd = rawinfo.Dest_fd
		value.Info[i].Dest_offset = rawinfo.Dest_offset
		value.Info[i].Bytes_deduped = rawinfo.Bytes_deduped
		value.Info[i].Status = rawinfo.Status
		value.Info[i].Reserved = rawinfo.Reserved
	}

	return err
}

func IoctlHIDGetDesc(fd int, value *HIDRawReportDescriptor) error {
	return ioctlPtr(fd, HIDIOCGRDESC, unsafe.Pointer(value))
}

func IoctlHIDGetRawInfo(fd int) (*HIDRawDevInfo, error) {
	var value HIDRawDevInfo
	err := ioctlPtr(fd, HIDIOCGRAWINFO, unsafe.Pointer(&value))
	return &value, err
}

func IoctlHIDGetRawName(fd int) (string, error) {
	var value [_HIDIOCGRAWNAME_LEN]byte
	err := ioctlPtr(fd, _HIDIOCGRAWNAME, unsafe.Pointer(&value[0]))
	return ByteSliceToString(value[:]), err
}

func IoctlHIDGetRawPhys(fd int) (string, error) {
	var value [_HIDIOCGRAWPHYS_LEN]byte
	err := ioctlPtr(fd, _HIDIOCGRAWPHYS, unsafe.Pointer(&value[0]))
	return ByteSliceToString(value[:]), err
}

func IoctlHIDGetRawUniq(fd int) (string, error) {
	var value [_HIDIOCGRAWUNIQ_LEN]byte
	err := ioctlPtr(fd, _HIDIOCGRAWUNIQ, unsafe.Pointer(&value[0]))
	return ByteSliceToString(value[:]), err
}



func IoctlIfreq(fd int, req uint, value *Ifreq) error {
	
	
	return ioctlPtr(fd, req, unsafe.Pointer(&value.raw))
}





func ioctlIfreqData(fd int, req uint, value *ifreqData) error {
	
	
	return ioctlPtr(fd, req, unsafe.Pointer(value))
}




func IoctlKCMClone(fd int) (*KCMClone, error) {
	var info KCMClone
	if err := ioctlPtr(fd, SIOCKCMCLONE, unsafe.Pointer(&info)); err != nil {
		return nil, err
	}

	return &info, nil
}



func IoctlKCMAttach(fd int, info KCMAttach) error {
	return ioctlPtr(fd, SIOCKCMATTACH, unsafe.Pointer(&info))
}


func IoctlKCMUnattach(fd int, info KCMUnattach) error {
	return ioctlPtr(fd, SIOCKCMUNATTACH, unsafe.Pointer(&info))
}



func IoctlLoopGetStatus64(fd int) (*LoopInfo64, error) {
	var value LoopInfo64
	if err := ioctlPtr(fd, LOOP_GET_STATUS64, unsafe.Pointer(&value)); err != nil {
		return nil, err
	}
	return &value, nil
}



func IoctlLoopSetStatus64(fd int, value *LoopInfo64) error {
	return ioctlPtr(fd, LOOP_SET_STATUS64, unsafe.Pointer(value))
}


func IoctlLoopConfigure(fd int, value *LoopConfig) error {
	return ioctlPtr(fd, LOOP_CONFIGURE, unsafe.Pointer(value))
}
