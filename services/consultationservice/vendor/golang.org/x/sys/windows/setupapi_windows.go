



package windows

import (
	"encoding/binary"
	"errors"
	"fmt"
	"runtime"
	"strings"
	"syscall"
	"unsafe"
)







const (
	ERROR_EXPECTED_SECTION_NAME                  Errno = 0x20000000 | 0xC0000000 | 0
	ERROR_BAD_SECTION_NAME_LINE                  Errno = 0x20000000 | 0xC0000000 | 1
	ERROR_SECTION_NAME_TOO_LONG                  Errno = 0x20000000 | 0xC0000000 | 2
	ERROR_GENERAL_SYNTAX                         Errno = 0x20000000 | 0xC0000000 | 3
	ERROR_WRONG_INF_STYLE                        Errno = 0x20000000 | 0xC0000000 | 0x100
	ERROR_SECTION_NOT_FOUND                      Errno = 0x20000000 | 0xC0000000 | 0x101
	ERROR_LINE_NOT_FOUND                         Errno = 0x20000000 | 0xC0000000 | 0x102
	ERROR_NO_BACKUP                              Errno = 0x20000000 | 0xC0000000 | 0x103
	ERROR_NO_ASSOCIATED_CLASS                    Errno = 0x20000000 | 0xC0000000 | 0x200
	ERROR_CLASS_MISMATCH                         Errno = 0x20000000 | 0xC0000000 | 0x201
	ERROR_DUPLICATE_FOUND                        Errno = 0x20000000 | 0xC0000000 | 0x202
	ERROR_NO_DRIVER_SELECTED                     Errno = 0x20000000 | 0xC0000000 | 0x203
	ERROR_KEY_DOES_NOT_EXIST                     Errno = 0x20000000 | 0xC0000000 | 0x204
	ERROR_INVALID_DEVINST_NAME                   Errno = 0x20000000 | 0xC0000000 | 0x205
	ERROR_INVALID_CLASS                          Errno = 0x20000000 | 0xC0000000 | 0x206
	ERROR_DEVINST_ALREADY_EXISTS                 Errno = 0x20000000 | 0xC0000000 | 0x207
	ERROR_DEVINFO_NOT_REGISTERED                 Errno = 0x20000000 | 0xC0000000 | 0x208
	ERROR_INVALID_REG_PROPERTY                   Errno = 0x20000000 | 0xC0000000 | 0x209
	ERROR_NO_INF                                 Errno = 0x20000000 | 0xC0000000 | 0x20A
	ERROR_NO_SUCH_DEVINST                        Errno = 0x20000000 | 0xC0000000 | 0x20B
	ERROR_CANT_LOAD_CLASS_ICON                   Errno = 0x20000000 | 0xC0000000 | 0x20C
	ERROR_INVALID_CLASS_INSTALLER                Errno = 0x20000000 | 0xC0000000 | 0x20D
	ERROR_DI_DO_DEFAULT                          Errno = 0x20000000 | 0xC0000000 | 0x20E
	ERROR_DI_NOFILECOPY                          Errno = 0x20000000 | 0xC0000000 | 0x20F
	ERROR_INVALID_HWPROFILE                      Errno = 0x20000000 | 0xC0000000 | 0x210
	ERROR_NO_DEVICE_SELECTED                     Errno = 0x20000000 | 0xC0000000 | 0x211
	ERROR_DEVINFO_LIST_LOCKED                    Errno = 0x20000000 | 0xC0000000 | 0x212
	ERROR_DEVINFO_DATA_LOCKED                    Errno = 0x20000000 | 0xC0000000 | 0x213
	ERROR_DI_BAD_PATH                            Errno = 0x20000000 | 0xC0000000 | 0x214
	ERROR_NO_CLASSINSTALL_PARAMS                 Errno = 0x20000000 | 0xC0000000 | 0x215
	ERROR_FILEQUEUE_LOCKED                       Errno = 0x20000000 | 0xC0000000 | 0x216
	ERROR_BAD_SERVICE_INSTALLSECT                Errno = 0x20000000 | 0xC0000000 | 0x217
	ERROR_NO_CLASS_DRIVER_LIST                   Errno = 0x20000000 | 0xC0000000 | 0x218
	ERROR_NO_ASSOCIATED_SERVICE                  Errno = 0x20000000 | 0xC0000000 | 0x219
	ERROR_NO_DEFAULT_DEVICE_INTERFACE            Errno = 0x20000000 | 0xC0000000 | 0x21A
	ERROR_DEVICE_INTERFACE_ACTIVE                Errno = 0x20000000 | 0xC0000000 | 0x21B
	ERROR_DEVICE_INTERFACE_REMOVED               Errno = 0x20000000 | 0xC0000000 | 0x21C
	ERROR_BAD_INTERFACE_INSTALLSECT              Errno = 0x20000000 | 0xC0000000 | 0x21D
	ERROR_NO_SUCH_INTERFACE_CLASS                Errno = 0x20000000 | 0xC0000000 | 0x21E
	ERROR_INVALID_REFERENCE_STRING               Errno = 0x20000000 | 0xC0000000 | 0x21F
	ERROR_INVALID_MACHINENAME                    Errno = 0x20000000 | 0xC0000000 | 0x220
	ERROR_REMOTE_COMM_FAILURE                    Errno = 0x20000000 | 0xC0000000 | 0x221
	ERROR_MACHINE_UNAVAILABLE                    Errno = 0x20000000 | 0xC0000000 | 0x222
	ERROR_NO_CONFIGMGR_SERVICES                  Errno = 0x20000000 | 0xC0000000 | 0x223
	ERROR_INVALID_PROPPAGE_PROVIDER              Errno = 0x20000000 | 0xC0000000 | 0x224
	ERROR_NO_SUCH_DEVICE_INTERFACE               Errno = 0x20000000 | 0xC0000000 | 0x225
	ERROR_DI_POSTPROCESSING_REQUIRED             Errno = 0x20000000 | 0xC0000000 | 0x226
	ERROR_INVALID_COINSTALLER                    Errno = 0x20000000 | 0xC0000000 | 0x227
	ERROR_NO_COMPAT_DRIVERS                      Errno = 0x20000000 | 0xC0000000 | 0x228
	ERROR_NO_DEVICE_ICON                         Errno = 0x20000000 | 0xC0000000 | 0x229
	ERROR_INVALID_INF_LOGCONFIG                  Errno = 0x20000000 | 0xC0000000 | 0x22A
	ERROR_DI_DONT_INSTALL                        Errno = 0x20000000 | 0xC0000000 | 0x22B
	ERROR_INVALID_FILTER_DRIVER                  Errno = 0x20000000 | 0xC0000000 | 0x22C
	ERROR_NON_WINDOWS_NT_DRIVER                  Errno = 0x20000000 | 0xC0000000 | 0x22D
	ERROR_NON_WINDOWS_DRIVER                     Errno = 0x20000000 | 0xC0000000 | 0x22E
	ERROR_NO_CATALOG_FOR_OEM_INF                 Errno = 0x20000000 | 0xC0000000 | 0x22F
	ERROR_DEVINSTALL_QUEUE_NONNATIVE             Errno = 0x20000000 | 0xC0000000 | 0x230
	ERROR_NOT_DISABLEABLE                        Errno = 0x20000000 | 0xC0000000 | 0x231
	ERROR_CANT_REMOVE_DEVINST                    Errno = 0x20000000 | 0xC0000000 | 0x232
	ERROR_INVALID_TARGET                         Errno = 0x20000000 | 0xC0000000 | 0x233
	ERROR_DRIVER_NONNATIVE                       Errno = 0x20000000 | 0xC0000000 | 0x234
	ERROR_IN_WOW64                               Errno = 0x20000000 | 0xC0000000 | 0x235
	ERROR_SET_SYSTEM_RESTORE_POINT               Errno = 0x20000000 | 0xC0000000 | 0x236
	ERROR_SCE_DISABLED                           Errno = 0x20000000 | 0xC0000000 | 0x238
	ERROR_UNKNOWN_EXCEPTION                      Errno = 0x20000000 | 0xC0000000 | 0x239
	ERROR_PNP_REGISTRY_ERROR                     Errno = 0x20000000 | 0xC0000000 | 0x23A
	ERROR_REMOTE_REQUEST_UNSUPPORTED             Errno = 0x20000000 | 0xC0000000 | 0x23B
	ERROR_NOT_AN_INSTALLED_OEM_INF               Errno = 0x20000000 | 0xC0000000 | 0x23C
	ERROR_INF_IN_USE_BY_DEVICES                  Errno = 0x20000000 | 0xC0000000 | 0x23D
	ERROR_DI_FUNCTION_OBSOLETE                   Errno = 0x20000000 | 0xC0000000 | 0x23E
	ERROR_NO_AUTHENTICODE_CATALOG                Errno = 0x20000000 | 0xC0000000 | 0x23F
	ERROR_AUTHENTICODE_DISALLOWED                Errno = 0x20000000 | 0xC0000000 | 0x240
	ERROR_AUTHENTICODE_TRUSTED_PUBLISHER         Errno = 0x20000000 | 0xC0000000 | 0x241
	ERROR_AUTHENTICODE_TRUST_NOT_ESTABLISHED     Errno = 0x20000000 | 0xC0000000 | 0x242
	ERROR_AUTHENTICODE_PUBLISHER_NOT_TRUSTED     Errno = 0x20000000 | 0xC0000000 | 0x243
	ERROR_SIGNATURE_OSATTRIBUTE_MISMATCH         Errno = 0x20000000 | 0xC0000000 | 0x244
	ERROR_ONLY_VALIDATE_VIA_AUTHENTICODE         Errno = 0x20000000 | 0xC0000000 | 0x245
	ERROR_DEVICE_INSTALLER_NOT_READY             Errno = 0x20000000 | 0xC0000000 | 0x246
	ERROR_DRIVER_STORE_ADD_FAILED                Errno = 0x20000000 | 0xC0000000 | 0x247
	ERROR_DEVICE_INSTALL_BLOCKED                 Errno = 0x20000000 | 0xC0000000 | 0x248
	ERROR_DRIVER_INSTALL_BLOCKED                 Errno = 0x20000000 | 0xC0000000 | 0x249
	ERROR_WRONG_INF_TYPE                         Errno = 0x20000000 | 0xC0000000 | 0x24A
	ERROR_FILE_HASH_NOT_IN_CATALOG               Errno = 0x20000000 | 0xC0000000 | 0x24B
	ERROR_DRIVER_STORE_DELETE_FAILED             Errno = 0x20000000 | 0xC0000000 | 0x24C
	ERROR_UNRECOVERABLE_STACK_OVERFLOW           Errno = 0x20000000 | 0xC0000000 | 0x300
	EXCEPTION_SPAPI_UNRECOVERABLE_STACK_OVERFLOW Errno = ERROR_UNRECOVERABLE_STACK_OVERFLOW
	ERROR_NO_DEFAULT_INTERFACE_DEVICE            Errno = ERROR_NO_DEFAULT_DEVICE_INTERFACE
	ERROR_INTERFACE_DEVICE_ACTIVE                Errno = ERROR_DEVICE_INTERFACE_ACTIVE
	ERROR_INTERFACE_DEVICE_REMOVED               Errno = ERROR_DEVICE_INTERFACE_REMOVED
	ERROR_NO_SUCH_INTERFACE_DEVICE               Errno = ERROR_NO_SUCH_DEVICE_INTERFACE
)

const (
	MAX_DEVICE_ID_LEN   = 200
	MAX_DEVNODE_ID_LEN  = MAX_DEVICE_ID_LEN
	MAX_GUID_STRING_LEN = 39 
	MAX_CLASS_NAME_LEN  = 32
	MAX_PROFILE_LEN     = 80
	MAX_CONFIG_VALUE    = 9999
	MAX_INSTANCE_VALUE  = 9999
	CONFIGMG_VERSION    = 0x0400
)


const (
	LINE_LEN                    = 256  
	MAX_INF_STRING_LENGTH       = 4096 
	MAX_INF_SECTION_NAME_LENGTH = 255  
	MAX_TITLE_LEN               = 60
	MAX_INSTRUCTION_LEN         = 256
	MAX_LABEL_LEN               = 30
	MAX_SERVICE_NAME_LEN        = 256
	MAX_SUBTITLE_LEN            = 256
)

const (
	
	SP_MAX_MACHINENAME_LENGTH = MAX_PATH + 3
)


type HSPFILEQ uintptr


type DevInfo Handle


type DEVINST uint32


type DevInfoData struct {
	size      uint32
	ClassGUID GUID
	DevInst   DEVINST
	_         uintptr
}


type DevInfoListDetailData struct {
	size                uint32 
	ClassGUID           GUID
	RemoteMachineHandle Handle
	remoteMachineName   [SP_MAX_MACHINENAME_LENGTH]uint16
}

func (*DevInfoListDetailData) unsafeSizeOf() uint32 {
	if unsafe.Sizeof(uintptr(0)) == 4 {
		
		return uint32(unsafe.Offsetof(DevInfoListDetailData{}.remoteMachineName) + unsafe.Sizeof(DevInfoListDetailData{}.remoteMachineName))
	}
	return uint32(unsafe.Sizeof(DevInfoListDetailData{}))
}

func (data *DevInfoListDetailData) RemoteMachineName() string {
	return UTF16ToString(data.remoteMachineName[:])
}

func (data *DevInfoListDetailData) SetRemoteMachineName(remoteMachineName string) error {
	str, err := UTF16FromString(remoteMachineName)
	if err != nil {
		return err
	}
	copy(data.remoteMachineName[:], str)
	return nil
}


type DI_FUNCTION uint32

const (
	DIF_SELECTDEVICE                   DI_FUNCTION = 0x00000001
	DIF_INSTALLDEVICE                  DI_FUNCTION = 0x00000002
	DIF_ASSIGNRESOURCES                DI_FUNCTION = 0x00000003
	DIF_PROPERTIES                     DI_FUNCTION = 0x00000004
	DIF_REMOVE                         DI_FUNCTION = 0x00000005
	DIF_FIRSTTIMESETUP                 DI_FUNCTION = 0x00000006
	DIF_FOUNDDEVICE                    DI_FUNCTION = 0x00000007
	DIF_SELECTCLASSDRIVERS             DI_FUNCTION = 0x00000008
	DIF_VALIDATECLASSDRIVERS           DI_FUNCTION = 0x00000009
	DIF_INSTALLCLASSDRIVERS            DI_FUNCTION = 0x0000000A
	DIF_CALCDISKSPACE                  DI_FUNCTION = 0x0000000B
	DIF_DESTROYPRIVATEDATA             DI_FUNCTION = 0x0000000C
	DIF_VALIDATEDRIVER                 DI_FUNCTION = 0x0000000D
	DIF_DETECT                         DI_FUNCTION = 0x0000000F
	DIF_INSTALLWIZARD                  DI_FUNCTION = 0x00000010
	DIF_DESTROYWIZARDDATA              DI_FUNCTION = 0x00000011
	DIF_PROPERTYCHANGE                 DI_FUNCTION = 0x00000012
	DIF_ENABLECLASS                    DI_FUNCTION = 0x00000013
	DIF_DETECTVERIFY                   DI_FUNCTION = 0x00000014
	DIF_INSTALLDEVICEFILES             DI_FUNCTION = 0x00000015
	DIF_UNREMOVE                       DI_FUNCTION = 0x00000016
	DIF_SELECTBESTCOMPATDRV            DI_FUNCTION = 0x00000017
	DIF_ALLOW_INSTALL                  DI_FUNCTION = 0x00000018
	DIF_REGISTERDEVICE                 DI_FUNCTION = 0x00000019
	DIF_NEWDEVICEWIZARD_PRESELECT      DI_FUNCTION = 0x0000001A
	DIF_NEWDEVICEWIZARD_SELECT         DI_FUNCTION = 0x0000001B
	DIF_NEWDEVICEWIZARD_PREANALYZE     DI_FUNCTION = 0x0000001C
	DIF_NEWDEVICEWIZARD_POSTANALYZE    DI_FUNCTION = 0x0000001D
	DIF_NEWDEVICEWIZARD_FINISHINSTALL  DI_FUNCTION = 0x0000001E
	DIF_INSTALLINTERFACES              DI_FUNCTION = 0x00000020
	DIF_DETECTCANCEL                   DI_FUNCTION = 0x00000021
	DIF_REGISTER_COINSTALLERS          DI_FUNCTION = 0x00000022
	DIF_ADDPROPERTYPAGE_ADVANCED       DI_FUNCTION = 0x00000023
	DIF_ADDPROPERTYPAGE_BASIC          DI_FUNCTION = 0x00000024
	DIF_TROUBLESHOOTER                 DI_FUNCTION = 0x00000026
	DIF_POWERMESSAGEWAKE               DI_FUNCTION = 0x00000027
	DIF_ADDREMOTEPROPERTYPAGE_ADVANCED DI_FUNCTION = 0x00000028
	DIF_UPDATEDRIVER_UI                DI_FUNCTION = 0x00000029
	DIF_FINISHINSTALL_ACTION           DI_FUNCTION = 0x0000002A
)


type DevInstallParams struct {
	size                     uint32
	Flags                    DI_FLAGS
	FlagsEx                  DI_FLAGSEX
	hwndParent               uintptr
	InstallMsgHandler        uintptr
	InstallMsgHandlerContext uintptr
	FileQueue                HSPFILEQ
	_                        uintptr
	_                        uint32
	driverPath               [MAX_PATH]uint16
}

func (params *DevInstallParams) DriverPath() string {
	return UTF16ToString(params.driverPath[:])
}

func (params *DevInstallParams) SetDriverPath(driverPath string) error {
	str, err := UTF16FromString(driverPath)
	if err != nil {
		return err
	}
	copy(params.driverPath[:], str)
	return nil
}


type DI_FLAGS uint32

const (
	
	DI_SHOWOEM       DI_FLAGS = 0x00000001 
	DI_SHOWCOMPAT    DI_FLAGS = 0x00000002 
	DI_SHOWCLASS     DI_FLAGS = 0x00000004 
	DI_SHOWALL       DI_FLAGS = 0x00000007 
	DI_NOVCP         DI_FLAGS = 0x00000008 
	DI_DIDCOMPAT     DI_FLAGS = 0x00000010 
	DI_DIDCLASS      DI_FLAGS = 0x00000020 
	DI_AUTOASSIGNRES DI_FLAGS = 0x00000040 

	
	DI_NEEDRESTART DI_FLAGS = 0x00000080 
	DI_NEEDREBOOT  DI_FLAGS = 0x00000100 

	
	DI_NOBROWSE DI_FLAGS = 0x00000200 

	
	DI_MULTMFGS DI_FLAGS = 0x00000400 

	
	DI_DISABLED DI_FLAGS = 0x00000800 

	
	DI_GENERALPAGE_ADDED  DI_FLAGS = 0x00001000
	DI_RESOURCEPAGE_ADDED DI_FLAGS = 0x00002000

	
	DI_PROPERTIES_CHANGE DI_FLAGS = 0x00004000

	
	DI_INF_IS_SORTED DI_FLAGS = 0x00008000

	
	DI_ENUMSINGLEINF DI_FLAGS = 0x00010000

	
	
	DI_DONOTCALLCONFIGMG DI_FLAGS = 0x00020000

	
	DI_INSTALLDISABLED DI_FLAGS = 0x00040000

	
	
	DI_COMPAT_FROM_CLASS DI_FLAGS = 0x00080000

	
	DI_CLASSINSTALLPARAMS DI_FLAGS = 0x00100000

	
	DI_NODI_DEFAULTACTION DI_FLAGS = 0x00200000

	
	DI_QUIETINSTALL        DI_FLAGS = 0x00800000 
	DI_NOFILECOPY          DI_FLAGS = 0x01000000 
	DI_FORCECOPY           DI_FLAGS = 0x02000000 
	DI_DRIVERPAGE_ADDED    DI_FLAGS = 0x04000000 
	DI_USECI_SELECTSTRINGS DI_FLAGS = 0x08000000 
	DI_OVERRIDE_INFFLAGS   DI_FLAGS = 0x10000000 
	DI_PROPS_NOCHANGEUSAGE DI_FLAGS = 0x20000000 

	DI_NOSELECTICONS DI_FLAGS = 0x40000000 

	DI_NOWRITE_IDS DI_FLAGS = 0x80000000 
)


type DI_FLAGSEX uint32

const (
	DI_FLAGSEX_CI_FAILED                DI_FLAGSEX = 0x00000004 
	DI_FLAGSEX_FINISHINSTALL_ACTION     DI_FLAGSEX = 0x00000008 
	DI_FLAGSEX_DIDINFOLIST              DI_FLAGSEX = 0x00000010 
	DI_FLAGSEX_DIDCOMPATINFO            DI_FLAGSEX = 0x00000020 
	DI_FLAGSEX_FILTERCLASSES            DI_FLAGSEX = 0x00000040
	DI_FLAGSEX_SETFAILEDINSTALL         DI_FLAGSEX = 0x00000080
	DI_FLAGSEX_DEVICECHANGE             DI_FLAGSEX = 0x00000100
	DI_FLAGSEX_ALWAYSWRITEIDS           DI_FLAGSEX = 0x00000200
	DI_FLAGSEX_PROPCHANGE_PENDING       DI_FLAGSEX = 0x00000400 
	DI_FLAGSEX_ALLOWEXCLUDEDDRVS        DI_FLAGSEX = 0x00000800
	DI_FLAGSEX_NOUIONQUERYREMOVE        DI_FLAGSEX = 0x00001000
	DI_FLAGSEX_USECLASSFORCOMPAT        DI_FLAGSEX = 0x00002000 
	DI_FLAGSEX_NO_DRVREG_MODIFY         DI_FLAGSEX = 0x00008000 
	DI_FLAGSEX_IN_SYSTEM_SETUP          DI_FLAGSEX = 0x00010000 
	DI_FLAGSEX_INET_DRIVER              DI_FLAGSEX = 0x00020000 
	DI_FLAGSEX_APPENDDRIVERLIST         DI_FLAGSEX = 0x00040000 
	DI_FLAGSEX_PREINSTALLBACKUP         DI_FLAGSEX = 0x00080000 
	DI_FLAGSEX_BACKUPONREPLACE          DI_FLAGSEX = 0x00100000 
	DI_FLAGSEX_DRIVERLIST_FROM_URL      DI_FLAGSEX = 0x00200000 
	DI_FLAGSEX_EXCLUDE_OLD_INET_DRIVERS DI_FLAGSEX = 0x00800000 
	DI_FLAGSEX_POWERPAGE_ADDED          DI_FLAGSEX = 0x01000000 
	DI_FLAGSEX_FILTERSIMILARDRIVERS     DI_FLAGSEX = 0x02000000 
	DI_FLAGSEX_INSTALLEDDRIVER          DI_FLAGSEX = 0x04000000 
	DI_FLAGSEX_NO_CLASSLIST_NODE_MERGE  DI_FLAGSEX = 0x08000000 
	DI_FLAGSEX_ALTPLATFORM_DRVSEARCH    DI_FLAGSEX = 0x10000000 
	DI_FLAGSEX_RESTART_DEVICE_ONLY      DI_FLAGSEX = 0x20000000 
	DI_FLAGSEX_RECURSIVESEARCH          DI_FLAGSEX = 0x40000000 
	DI_FLAGSEX_SEARCH_PUBLISHED_INFS    DI_FLAGSEX = 0x80000000 
)


type ClassInstallHeader struct {
	size            uint32
	InstallFunction DI_FUNCTION
}

func MakeClassInstallHeader(installFunction DI_FUNCTION) *ClassInstallHeader {
	hdr := &ClassInstallHeader{InstallFunction: installFunction}
	hdr.size = uint32(unsafe.Sizeof(*hdr))
	return hdr
}


type DICS_STATE uint32

const (
	DICS_ENABLE     DICS_STATE = 0x00000001 
	DICS_DISABLE    DICS_STATE = 0x00000002 
	DICS_PROPCHANGE DICS_STATE = 0x00000003 
	DICS_START      DICS_STATE = 0x00000004 
	DICS_STOP       DICS_STATE = 0x00000005 
)


type DICS_FLAG uint32

const (
	DICS_FLAG_GLOBAL         DICS_FLAG = 0x00000001 
	DICS_FLAG_CONFIGSPECIFIC DICS_FLAG = 0x00000002 
	DICS_FLAG_CONFIGGENERAL  DICS_FLAG = 0x00000004 
)


type PropChangeParams struct {
	ClassInstallHeader ClassInstallHeader
	StateChange        DICS_STATE
	Scope              DICS_FLAG
	HwProfile          uint32
}


type DI_REMOVEDEVICE uint32

const (
	DI_REMOVEDEVICE_GLOBAL         DI_REMOVEDEVICE = 0x00000001 
	DI_REMOVEDEVICE_CONFIGSPECIFIC DI_REMOVEDEVICE = 0x00000002 
)


type RemoveDeviceParams struct {
	ClassInstallHeader ClassInstallHeader
	Scope              DI_REMOVEDEVICE
	HwProfile          uint32
}


type DrvInfoData struct {
	size          uint32
	DriverType    uint32
	_             uintptr
	description   [LINE_LEN]uint16
	mfgName       [LINE_LEN]uint16
	providerName  [LINE_LEN]uint16
	DriverDate    Filetime
	DriverVersion uint64
}

func (data *DrvInfoData) Description() string {
	return UTF16ToString(data.description[:])
}

func (data *DrvInfoData) SetDescription(description string) error {
	str, err := UTF16FromString(description)
	if err != nil {
		return err
	}
	copy(data.description[:], str)
	return nil
}

func (data *DrvInfoData) MfgName() string {
	return UTF16ToString(data.mfgName[:])
}

func (data *DrvInfoData) SetMfgName(mfgName string) error {
	str, err := UTF16FromString(mfgName)
	if err != nil {
		return err
	}
	copy(data.mfgName[:], str)
	return nil
}

func (data *DrvInfoData) ProviderName() string {
	return UTF16ToString(data.providerName[:])
}

func (data *DrvInfoData) SetProviderName(providerName string) error {
	str, err := UTF16FromString(providerName)
	if err != nil {
		return err
	}
	copy(data.providerName[:], str)
	return nil
}


func (data *DrvInfoData) IsNewer(driverDate Filetime, driverVersion uint64) bool {
	if data.DriverDate.HighDateTime > driverDate.HighDateTime {
		return true
	}
	if data.DriverDate.HighDateTime < driverDate.HighDateTime {
		return false
	}

	if data.DriverDate.LowDateTime > driverDate.LowDateTime {
		return true
	}
	if data.DriverDate.LowDateTime < driverDate.LowDateTime {
		return false
	}

	if data.DriverVersion > driverVersion {
		return true
	}
	if data.DriverVersion < driverVersion {
		return false
	}

	return false
}


type DrvInfoDetailData struct {
	size            uint32 
	InfDate         Filetime
	compatIDsOffset uint32
	compatIDsLength uint32
	_               uintptr
	sectionName     [LINE_LEN]uint16
	infFileName     [MAX_PATH]uint16
	drvDescription  [LINE_LEN]uint16
	hardwareID      [1]uint16
}

func (*DrvInfoDetailData) unsafeSizeOf() uint32 {
	if unsafe.Sizeof(uintptr(0)) == 4 {
		
		return uint32(unsafe.Offsetof(DrvInfoDetailData{}.hardwareID) + unsafe.Sizeof(DrvInfoDetailData{}.hardwareID))
	}
	return uint32(unsafe.Sizeof(DrvInfoDetailData{}))
}

func (data *DrvInfoDetailData) SectionName() string {
	return UTF16ToString(data.sectionName[:])
}

func (data *DrvInfoDetailData) InfFileName() string {
	return UTF16ToString(data.infFileName[:])
}

func (data *DrvInfoDetailData) DrvDescription() string {
	return UTF16ToString(data.drvDescription[:])
}

func (data *DrvInfoDetailData) HardwareID() string {
	if data.compatIDsOffset > 1 {
		bufW := data.getBuf()
		return UTF16ToString(bufW[:wcslen(bufW)])
	}

	return ""
}

func (data *DrvInfoDetailData) CompatIDs() []string {
	a := make([]string, 0)

	if data.compatIDsLength > 0 {
		bufW := data.getBuf()
		bufW = bufW[data.compatIDsOffset : data.compatIDsOffset+data.compatIDsLength]
		for i := 0; i < len(bufW); {
			j := i + wcslen(bufW[i:])
			if i < j {
				a = append(a, UTF16ToString(bufW[i:j]))
			}
			i = j + 1
		}
	}

	return a
}

func (data *DrvInfoDetailData) getBuf() []uint16 {
	len := (data.size - uint32(unsafe.Offsetof(data.hardwareID))) / 2
	sl := struct {
		addr *uint16
		len  int
		cap  int
	}{&data.hardwareID[0], int(len), int(len)}
	return *(*[]uint16)(unsafe.Pointer(&sl))
}


func (data *DrvInfoDetailData) IsCompatible(hwid string) bool {
	hwidLC := strings.ToLower(hwid)
	if strings.ToLower(data.HardwareID()) == hwidLC {
		return true
	}
	a := data.CompatIDs()
	for i := range a {
		if strings.ToLower(a[i]) == hwidLC {
			return true
		}
	}

	return false
}


type DICD uint32

const (
	DICD_GENERATE_ID       DICD = 0x00000001
	DICD_INHERIT_CLASSDRVS DICD = 0x00000002
)


type SUOI uint32

const (
	SUOI_FORCEDELETE SUOI = 0x0001
)




type SPDIT uint32

const (
	SPDIT_NODRIVER     SPDIT = 0x00000000
	SPDIT_CLASSDRIVER  SPDIT = 0x00000001
	SPDIT_COMPATDRIVER SPDIT = 0x00000002
)


type DIGCF uint32

const (
	DIGCF_DEFAULT         DIGCF = 0x00000001 
	DIGCF_PRESENT         DIGCF = 0x00000002
	DIGCF_ALLCLASSES      DIGCF = 0x00000004
	DIGCF_PROFILE         DIGCF = 0x00000008
	DIGCF_DEVICEINTERFACE DIGCF = 0x00000010
)


type DIREG uint32

const (
	DIREG_DEV  DIREG = 0x00000001 
	DIREG_DRV  DIREG = 0x00000002 
	DIREG_BOTH DIREG = 0x00000004 
)









type SPDRP uint32

const (
	SPDRP_DEVICEDESC                  SPDRP = 0x00000000 
	SPDRP_HARDWAREID                  SPDRP = 0x00000001 
	SPDRP_COMPATIBLEIDS               SPDRP = 0x00000002 
	SPDRP_SERVICE                     SPDRP = 0x00000004 
	SPDRP_CLASS                       SPDRP = 0x00000007 
	SPDRP_CLASSGUID                   SPDRP = 0x00000008 
	SPDRP_DRIVER                      SPDRP = 0x00000009 
	SPDRP_CONFIGFLAGS                 SPDRP = 0x0000000A 
	SPDRP_MFG                         SPDRP = 0x0000000B 
	SPDRP_FRIENDLYNAME                SPDRP = 0x0000000C 
	SPDRP_LOCATION_INFORMATION        SPDRP = 0x0000000D 
	SPDRP_PHYSICAL_DEVICE_OBJECT_NAME SPDRP = 0x0000000E 
	SPDRP_CAPABILITIES                SPDRP = 0x0000000F 
	SPDRP_UI_NUMBER                   SPDRP = 0x00000010 
	SPDRP_UPPERFILTERS                SPDRP = 0x00000011 
	SPDRP_LOWERFILTERS                SPDRP = 0x00000012 
	SPDRP_BUSTYPEGUID                 SPDRP = 0x00000013 
	SPDRP_LEGACYBUSTYPE               SPDRP = 0x00000014 
	SPDRP_BUSNUMBER                   SPDRP = 0x00000015 
	SPDRP_ENUMERATOR_NAME             SPDRP = 0x00000016 
	SPDRP_SECURITY                    SPDRP = 0x00000017 
	SPDRP_SECURITY_SDS                SPDRP = 0x00000018 
	SPDRP_DEVTYPE                     SPDRP = 0x00000019 
	SPDRP_EXCLUSIVE                   SPDRP = 0x0000001A 
	SPDRP_CHARACTERISTICS             SPDRP = 0x0000001B 
	SPDRP_ADDRESS                     SPDRP = 0x0000001C 
	SPDRP_UI_NUMBER_DESC_FORMAT       SPDRP = 0x0000001D 
	SPDRP_DEVICE_POWER_DATA           SPDRP = 0x0000001E 
	SPDRP_REMOVAL_POLICY              SPDRP = 0x0000001F 
	SPDRP_REMOVAL_POLICY_HW_DEFAULT   SPDRP = 0x00000020 
	SPDRP_REMOVAL_POLICY_OVERRIDE     SPDRP = 0x00000021 
	SPDRP_INSTALL_STATE               SPDRP = 0x00000022 
	SPDRP_LOCATION_PATHS              SPDRP = 0x00000023 
	SPDRP_BASE_CONTAINERID            SPDRP = 0x00000024 

	SPDRP_MAXIMUM_PROPERTY SPDRP = 0x00000025 
)



type DEVPROPTYPE uint32

const (
	DEVPROP_TYPEMOD_ARRAY DEVPROPTYPE = 0x00001000
	DEVPROP_TYPEMOD_LIST  DEVPROPTYPE = 0x00002000

	DEVPROP_TYPE_EMPTY                      DEVPROPTYPE = 0x00000000
	DEVPROP_TYPE_NULL                       DEVPROPTYPE = 0x00000001
	DEVPROP_TYPE_SBYTE                      DEVPROPTYPE = 0x00000002
	DEVPROP_TYPE_BYTE                       DEVPROPTYPE = 0x00000003
	DEVPROP_TYPE_INT16                      DEVPROPTYPE = 0x00000004
	DEVPROP_TYPE_UINT16                     DEVPROPTYPE = 0x00000005
	DEVPROP_TYPE_INT32                      DEVPROPTYPE = 0x00000006
	DEVPROP_TYPE_UINT32                     DEVPROPTYPE = 0x00000007
	DEVPROP_TYPE_INT64                      DEVPROPTYPE = 0x00000008
	DEVPROP_TYPE_UINT64                     DEVPROPTYPE = 0x00000009
	DEVPROP_TYPE_FLOAT                      DEVPROPTYPE = 0x0000000A
	DEVPROP_TYPE_DOUBLE                     DEVPROPTYPE = 0x0000000B
	DEVPROP_TYPE_DECIMAL                    DEVPROPTYPE = 0x0000000C
	DEVPROP_TYPE_GUID                       DEVPROPTYPE = 0x0000000D
	DEVPROP_TYPE_CURRENCY                   DEVPROPTYPE = 0x0000000E
	DEVPROP_TYPE_DATE                       DEVPROPTYPE = 0x0000000F
	DEVPROP_TYPE_FILETIME                   DEVPROPTYPE = 0x00000010
	DEVPROP_TYPE_BOOLEAN                    DEVPROPTYPE = 0x00000011
	DEVPROP_TYPE_STRING                     DEVPROPTYPE = 0x00000012
	DEVPROP_TYPE_STRING_LIST                DEVPROPTYPE = DEVPROP_TYPE_STRING | DEVPROP_TYPEMOD_LIST
	DEVPROP_TYPE_SECURITY_DESCRIPTOR        DEVPROPTYPE = 0x00000013
	DEVPROP_TYPE_SECURITY_DESCRIPTOR_STRING DEVPROPTYPE = 0x00000014
	DEVPROP_TYPE_DEVPROPKEY                 DEVPROPTYPE = 0x00000015
	DEVPROP_TYPE_DEVPROPTYPE                DEVPROPTYPE = 0x00000016
	DEVPROP_TYPE_BINARY                     DEVPROPTYPE = DEVPROP_TYPE_BYTE | DEVPROP_TYPEMOD_ARRAY
	DEVPROP_TYPE_ERROR                      DEVPROPTYPE = 0x00000017
	DEVPROP_TYPE_NTSTATUS                   DEVPROPTYPE = 0x00000018
	DEVPROP_TYPE_STRING_INDIRECT            DEVPROPTYPE = 0x00000019

	MAX_DEVPROP_TYPE    DEVPROPTYPE = 0x00000019
	MAX_DEVPROP_TYPEMOD DEVPROPTYPE = 0x00002000

	DEVPROP_MASK_TYPE    DEVPROPTYPE = 0x00000FFF
	DEVPROP_MASK_TYPEMOD DEVPROPTYPE = 0x0000F000
)


type DEVPROPGUID GUID


type DEVPROPID uint32

const DEVPROPID_FIRST_USABLE DEVPROPID = 2



type DEVPROPKEY struct {
	FmtID DEVPROPGUID
	PID   DEVPROPID
}


type CONFIGRET uint32

func (ret CONFIGRET) Error() string {
	if win32Error, ok := ret.Unwrap().(Errno); ok {
		return fmt.Sprintf("%s (CfgMgr error: 0x%08x)", win32Error.Error(), uint32(ret))
	}
	return fmt.Sprintf("CfgMgr error: 0x%08x", uint32(ret))
}

func (ret CONFIGRET) Win32Error(defaultError Errno) Errno {
	return cm_MapCrToWin32Err(ret, defaultError)
}

func (ret CONFIGRET) Unwrap() error {
	const noMatch = Errno(^uintptr(0))
	win32Error := ret.Win32Error(noMatch)
	if win32Error == noMatch {
		return nil
	}
	return win32Error
}

const (
	CR_SUCCESS                  CONFIGRET = 0x00000000
	CR_DEFAULT                  CONFIGRET = 0x00000001
	CR_OUT_OF_MEMORY            CONFIGRET = 0x00000002
	CR_INVALID_POINTER          CONFIGRET = 0x00000003
	CR_INVALID_FLAG             CONFIGRET = 0x00000004
	CR_INVALID_DEVNODE          CONFIGRET = 0x00000005
	CR_INVALID_DEVINST                    = CR_INVALID_DEVNODE
	CR_INVALID_RES_DES          CONFIGRET = 0x00000006
	CR_INVALID_LOG_CONF         CONFIGRET = 0x00000007
	CR_INVALID_ARBITRATOR       CONFIGRET = 0x00000008
	CR_INVALID_NODELIST         CONFIGRET = 0x00000009
	CR_DEVNODE_HAS_REQS         CONFIGRET = 0x0000000A
	CR_DEVINST_HAS_REQS                   = CR_DEVNODE_HAS_REQS
	CR_INVALID_RESOURCEID       CONFIGRET = 0x0000000B
	CR_DLVXD_NOT_FOUND          CONFIGRET = 0x0000000C
	CR_NO_SUCH_DEVNODE          CONFIGRET = 0x0000000D
	CR_NO_SUCH_DEVINST                    = CR_NO_SUCH_DEVNODE
	CR_NO_MORE_LOG_CONF         CONFIGRET = 0x0000000E
	CR_NO_MORE_RES_DES          CONFIGRET = 0x0000000F
	CR_ALREADY_SUCH_DEVNODE     CONFIGRET = 0x00000010
	CR_ALREADY_SUCH_DEVINST               = CR_ALREADY_SUCH_DEVNODE
	CR_INVALID_RANGE_LIST       CONFIGRET = 0x00000011
	CR_INVALID_RANGE            CONFIGRET = 0x00000012
	CR_FAILURE                  CONFIGRET = 0x00000013
	CR_NO_SUCH_LOGICAL_DEV      CONFIGRET = 0x00000014
	CR_CREATE_BLOCKED           CONFIGRET = 0x00000015
	CR_NOT_SYSTEM_VM            CONFIGRET = 0x00000016
	CR_REMOVE_VETOED            CONFIGRET = 0x00000017
	CR_APM_VETOED               CONFIGRET = 0x00000018
	CR_INVALID_LOAD_TYPE        CONFIGRET = 0x00000019
	CR_BUFFER_SMALL             CONFIGRET = 0x0000001A
	CR_NO_ARBITRATOR            CONFIGRET = 0x0000001B
	CR_NO_REGISTRY_HANDLE       CONFIGRET = 0x0000001C
	CR_REGISTRY_ERROR           CONFIGRET = 0x0000001D
	CR_INVALID_DEVICE_ID        CONFIGRET = 0x0000001E
	CR_INVALID_DATA             CONFIGRET = 0x0000001F
	CR_INVALID_API              CONFIGRET = 0x00000020
	CR_DEVLOADER_NOT_READY      CONFIGRET = 0x00000021
	CR_NEED_RESTART             CONFIGRET = 0x00000022
	CR_NO_MORE_HW_PROFILES      CONFIGRET = 0x00000023
	CR_DEVICE_NOT_THERE         CONFIGRET = 0x00000024
	CR_NO_SUCH_VALUE            CONFIGRET = 0x00000025
	CR_WRONG_TYPE               CONFIGRET = 0x00000026
	CR_INVALID_PRIORITY         CONFIGRET = 0x00000027
	CR_NOT_DISABLEABLE          CONFIGRET = 0x00000028
	CR_FREE_RESOURCES           CONFIGRET = 0x00000029
	CR_QUERY_VETOED             CONFIGRET = 0x0000002A
	CR_CANT_SHARE_IRQ           CONFIGRET = 0x0000002B
	CR_NO_DEPENDENT             CONFIGRET = 0x0000002C
	CR_SAME_RESOURCES           CONFIGRET = 0x0000002D
	CR_NO_SUCH_REGISTRY_KEY     CONFIGRET = 0x0000002E
	CR_INVALID_MACHINENAME      CONFIGRET = 0x0000002F
	CR_REMOTE_COMM_FAILURE      CONFIGRET = 0x00000030
	CR_MACHINE_UNAVAILABLE      CONFIGRET = 0x00000031
	CR_NO_CM_SERVICES           CONFIGRET = 0x00000032
	CR_ACCESS_DENIED            CONFIGRET = 0x00000033
	CR_CALL_NOT_IMPLEMENTED     CONFIGRET = 0x00000034
	CR_INVALID_PROPERTY         CONFIGRET = 0x00000035
	CR_DEVICE_INTERFACE_ACTIVE  CONFIGRET = 0x00000036
	CR_NO_SUCH_DEVICE_INTERFACE CONFIGRET = 0x00000037
	CR_INVALID_REFERENCE_STRING CONFIGRET = 0x00000038
	CR_INVALID_CONFLICT_LIST    CONFIGRET = 0x00000039
	CR_INVALID_INDEX            CONFIGRET = 0x0000003A
	CR_INVALID_STRUCTURE_SIZE   CONFIGRET = 0x0000003B
	NUM_CR_RESULTS              CONFIGRET = 0x0000003C
)

const (
	CM_GET_DEVICE_INTERFACE_LIST_PRESENT     = 0 
	CM_GET_DEVICE_INTERFACE_LIST_ALL_DEVICES = 1 
)

const (
	DN_ROOT_ENUMERATED       = 0x00000001        
	DN_DRIVER_LOADED         = 0x00000002        
	DN_ENUM_LOADED           = 0x00000004        
	DN_STARTED               = 0x00000008        
	DN_MANUAL                = 0x00000010        
	DN_NEED_TO_ENUM          = 0x00000020        
	DN_NOT_FIRST_TIME        = 0x00000040        
	DN_HARDWARE_ENUM         = 0x00000080        
	DN_LIAR                  = 0x00000100        
	DN_HAS_MARK              = 0x00000200        
	DN_HAS_PROBLEM           = 0x00000400        
	DN_FILTERED              = 0x00000800        
	DN_MOVED                 = 0x00001000        
	DN_DISABLEABLE           = 0x00002000        
	DN_REMOVABLE             = 0x00004000        
	DN_PRIVATE_PROBLEM       = 0x00008000        
	DN_MF_PARENT             = 0x00010000        
	DN_MF_CHILD              = 0x00020000        
	DN_WILL_BE_REMOVED       = 0x00040000        
	DN_NOT_FIRST_TIMEE       = 0x00080000        
	DN_STOP_FREE_RES         = 0x00100000        
	DN_REBAL_CANDIDATE       = 0x00200000        
	DN_BAD_PARTIAL           = 0x00400000        
	DN_NT_ENUMERATOR         = 0x00800000        
	DN_NT_DRIVER             = 0x01000000        
	DN_NEEDS_LOCKING         = 0x02000000        
	DN_ARM_WAKEUP            = 0x04000000        
	DN_APM_ENUMERATOR        = 0x08000000        
	DN_APM_DRIVER            = 0x10000000        
	DN_SILENT_INSTALL        = 0x20000000        
	DN_NO_SHOW_IN_DM         = 0x40000000        
	DN_BOOT_LOG_PROB         = 0x80000000        
	DN_NEED_RESTART          = DN_LIAR           
	DN_DRIVER_BLOCKED        = DN_NOT_FIRST_TIME 
	DN_LEGACY_DRIVER         = DN_MOVED          
	DN_CHILD_WITH_INVALID_ID = DN_HAS_MARK       
	DN_DEVICE_DISCONNECTED   = DN_NEEDS_LOCKING  
	DN_QUERY_REMOVE_PENDING  = DN_MF_PARENT      
	DN_QUERY_REMOVE_ACTIVE   = DN_MF_CHILD       
	DN_CHANGEABLE_FLAGS      = DN_NOT_FIRST_TIME | DN_HARDWARE_ENUM | DN_HAS_MARK | DN_DISABLEABLE | DN_REMOVABLE | DN_MF_CHILD | DN_MF_PARENT | DN_NOT_FIRST_TIMEE | DN_STOP_FREE_RES | DN_REBAL_CANDIDATE | DN_NT_ENUMERATOR | DN_NT_DRIVER | DN_SILENT_INSTALL | DN_NO_SHOW_IN_DM
)




func SetupDiCreateDeviceInfoListEx(classGUID *GUID, hwndParent uintptr, machineName string) (deviceInfoSet DevInfo, err error) {
	var machineNameUTF16 *uint16
	if machineName != "" {
		machineNameUTF16, err = UTF16PtrFromString(machineName)
		if err != nil {
			return
		}
	}
	return setupDiCreateDeviceInfoListEx(classGUID, hwndParent, machineNameUTF16, 0)
}




func SetupDiGetDeviceInfoListDetail(deviceInfoSet DevInfo) (deviceInfoSetDetailData *DevInfoListDetailData, err error) {
	data := &DevInfoListDetailData{}
	data.size = data.unsafeSizeOf()

	return data, setupDiGetDeviceInfoListDetail(deviceInfoSet, data)
}


func (deviceInfoSet DevInfo) DeviceInfoListDetail() (*DevInfoListDetailData, error) {
	return SetupDiGetDeviceInfoListDetail(deviceInfoSet)
}




func SetupDiCreateDeviceInfo(deviceInfoSet DevInfo, deviceName string, classGUID *GUID, deviceDescription string, hwndParent uintptr, creationFlags DICD) (deviceInfoData *DevInfoData, err error) {
	deviceNameUTF16, err := UTF16PtrFromString(deviceName)
	if err != nil {
		return
	}

	var deviceDescriptionUTF16 *uint16
	if deviceDescription != "" {
		deviceDescriptionUTF16, err = UTF16PtrFromString(deviceDescription)
		if err != nil {
			return
		}
	}

	data := &DevInfoData{}
	data.size = uint32(unsafe.Sizeof(*data))

	return data, setupDiCreateDeviceInfo(deviceInfoSet, deviceNameUTF16, classGUID, deviceDescriptionUTF16, hwndParent, creationFlags, data)
}


func (deviceInfoSet DevInfo) CreateDeviceInfo(deviceName string, classGUID *GUID, deviceDescription string, hwndParent uintptr, creationFlags DICD) (*DevInfoData, error) {
	return SetupDiCreateDeviceInfo(deviceInfoSet, deviceName, classGUID, deviceDescription, hwndParent, creationFlags)
}




func SetupDiEnumDeviceInfo(deviceInfoSet DevInfo, memberIndex int) (*DevInfoData, error) {
	data := &DevInfoData{}
	data.size = uint32(unsafe.Sizeof(*data))

	return data, setupDiEnumDeviceInfo(deviceInfoSet, uint32(memberIndex), data)
}


func (deviceInfoSet DevInfo) EnumDeviceInfo(memberIndex int) (*DevInfoData, error) {
	return SetupDiEnumDeviceInfo(deviceInfoSet, memberIndex)
}





func (deviceInfoSet DevInfo) Close() error {
	return SetupDiDestroyDeviceInfoList(deviceInfoSet)
}




func (deviceInfoSet DevInfo) BuildDriverInfoList(deviceInfoData *DevInfoData, driverType SPDIT) error {
	return SetupDiBuildDriverInfoList(deviceInfoSet, deviceInfoData, driverType)
}




func (deviceInfoSet DevInfo) CancelDriverInfoSearch() error {
	return SetupDiCancelDriverInfoSearch(deviceInfoSet)
}




func SetupDiEnumDriverInfo(deviceInfoSet DevInfo, deviceInfoData *DevInfoData, driverType SPDIT, memberIndex int) (*DrvInfoData, error) {
	data := &DrvInfoData{}
	data.size = uint32(unsafe.Sizeof(*data))

	return data, setupDiEnumDriverInfo(deviceInfoSet, deviceInfoData, driverType, uint32(memberIndex), data)
}


func (deviceInfoSet DevInfo) EnumDriverInfo(deviceInfoData *DevInfoData, driverType SPDIT, memberIndex int) (*DrvInfoData, error) {
	return SetupDiEnumDriverInfo(deviceInfoSet, deviceInfoData, driverType, memberIndex)
}




func SetupDiGetSelectedDriver(deviceInfoSet DevInfo, deviceInfoData *DevInfoData) (*DrvInfoData, error) {
	data := &DrvInfoData{}
	data.size = uint32(unsafe.Sizeof(*data))

	return data, setupDiGetSelectedDriver(deviceInfoSet, deviceInfoData, data)
}


func (deviceInfoSet DevInfo) SelectedDriver(deviceInfoData *DevInfoData) (*DrvInfoData, error) {
	return SetupDiGetSelectedDriver(deviceInfoSet, deviceInfoData)
}




func (deviceInfoSet DevInfo) SetSelectedDriver(deviceInfoData *DevInfoData, driverInfoData *DrvInfoData) error {
	return SetupDiSetSelectedDriver(deviceInfoSet, deviceInfoData, driverInfoData)
}




func SetupDiGetDriverInfoDetail(deviceInfoSet DevInfo, deviceInfoData *DevInfoData, driverInfoData *DrvInfoData) (*DrvInfoDetailData, error) {
	reqSize := uint32(2048)
	for {
		buf := make([]byte, reqSize)
		data := (*DrvInfoDetailData)(unsafe.Pointer(&buf[0]))
		data.size = data.unsafeSizeOf()
		err := setupDiGetDriverInfoDetail(deviceInfoSet, deviceInfoData, driverInfoData, data, uint32(len(buf)), &reqSize)
		if err == ERROR_INSUFFICIENT_BUFFER {
			continue
		}
		if err != nil {
			return nil, err
		}
		data.size = reqSize
		return data, nil
	}
}


func (deviceInfoSet DevInfo) DriverInfoDetail(deviceInfoData *DevInfoData, driverInfoData *DrvInfoData) (*DrvInfoDetailData, error) {
	return SetupDiGetDriverInfoDetail(deviceInfoSet, deviceInfoData, driverInfoData)
}




func (deviceInfoSet DevInfo) DestroyDriverInfoList(deviceInfoData *DevInfoData, driverType SPDIT) error {
	return SetupDiDestroyDriverInfoList(deviceInfoSet, deviceInfoData, driverType)
}




func SetupDiGetClassDevsEx(classGUID *GUID, enumerator string, hwndParent uintptr, flags DIGCF, deviceInfoSet DevInfo, machineName string) (handle DevInfo, err error) {
	var enumeratorUTF16 *uint16
	if enumerator != "" {
		enumeratorUTF16, err = UTF16PtrFromString(enumerator)
		if err != nil {
			return
		}
	}
	var machineNameUTF16 *uint16
	if machineName != "" {
		machineNameUTF16, err = UTF16PtrFromString(machineName)
		if err != nil {
			return
		}
	}
	return setupDiGetClassDevsEx(classGUID, enumeratorUTF16, hwndParent, flags, deviceInfoSet, machineNameUTF16, 0)
}





func (deviceInfoSet DevInfo) CallClassInstaller(installFunction DI_FUNCTION, deviceInfoData *DevInfoData) error {
	return SetupDiCallClassInstaller(installFunction, deviceInfoSet, deviceInfoData)
}





func (deviceInfoSet DevInfo) OpenDevRegKey(DeviceInfoData *DevInfoData, Scope DICS_FLAG, HwProfile uint32, KeyType DIREG, samDesired uint32) (Handle, error) {
	return SetupDiOpenDevRegKey(deviceInfoSet, DeviceInfoData, Scope, HwProfile, KeyType, samDesired)
}




func SetupDiGetDeviceProperty(deviceInfoSet DevInfo, deviceInfoData *DevInfoData, propertyKey *DEVPROPKEY) (value interface{}, err error) {
	reqSize := uint32(256)
	for {
		var dataType DEVPROPTYPE
		buf := make([]byte, reqSize)
		err = setupDiGetDeviceProperty(deviceInfoSet, deviceInfoData, propertyKey, &dataType, &buf[0], uint32(len(buf)), &reqSize, 0)
		if err == ERROR_INSUFFICIENT_BUFFER {
			continue
		}
		if err != nil {
			return
		}
		switch dataType {
		case DEVPROP_TYPE_STRING:
			ret := UTF16ToString(bufToUTF16(buf))
			runtime.KeepAlive(buf)
			return ret, nil
		}
		return nil, errors.New("unimplemented property type")
	}
}




func SetupDiGetDeviceRegistryProperty(deviceInfoSet DevInfo, deviceInfoData *DevInfoData, property SPDRP) (value interface{}, err error) {
	reqSize := uint32(256)
	for {
		var dataType uint32
		buf := make([]byte, reqSize)
		err = setupDiGetDeviceRegistryProperty(deviceInfoSet, deviceInfoData, property, &dataType, &buf[0], uint32(len(buf)), &reqSize)
		if err == ERROR_INSUFFICIENT_BUFFER {
			continue
		}
		if err != nil {
			return
		}
		return getRegistryValue(buf[:reqSize], dataType)
	}
}

func getRegistryValue(buf []byte, dataType uint32) (interface{}, error) {
	switch dataType {
	case REG_SZ:
		ret := UTF16ToString(bufToUTF16(buf))
		runtime.KeepAlive(buf)
		return ret, nil
	case REG_EXPAND_SZ:
		value := UTF16ToString(bufToUTF16(buf))
		if value == "" {
			return "", nil
		}
		p, err := syscall.UTF16PtrFromString(value)
		if err != nil {
			return "", err
		}
		ret := make([]uint16, 100)
		for {
			n, err := ExpandEnvironmentStrings(p, &ret[0], uint32(len(ret)))
			if err != nil {
				return "", err
			}
			if n <= uint32(len(ret)) {
				return UTF16ToString(ret[:n]), nil
			}
			ret = make([]uint16, n)
		}
	case REG_BINARY:
		return buf, nil
	case REG_DWORD_LITTLE_ENDIAN:
		return binary.LittleEndian.Uint32(buf), nil
	case REG_DWORD_BIG_ENDIAN:
		return binary.BigEndian.Uint32(buf), nil
	case REG_MULTI_SZ:
		bufW := bufToUTF16(buf)
		a := []string{}
		for i := 0; i < len(bufW); {
			j := i + wcslen(bufW[i:])
			if i < j {
				a = append(a, UTF16ToString(bufW[i:j]))
			}
			i = j + 1
		}
		runtime.KeepAlive(buf)
		return a, nil
	case REG_QWORD_LITTLE_ENDIAN:
		return binary.LittleEndian.Uint64(buf), nil
	default:
		return nil, fmt.Errorf("Unsupported registry value type: %v", dataType)
	}
}


func bufToUTF16(buf []byte) []uint16 {
	sl := struct {
		addr *uint16
		len  int
		cap  int
	}{(*uint16)(unsafe.Pointer(&buf[0])), len(buf) / 2, cap(buf) / 2}
	return *(*[]uint16)(unsafe.Pointer(&sl))
}


func utf16ToBuf(buf []uint16) []byte {
	sl := struct {
		addr *byte
		len  int
		cap  int
	}{(*byte)(unsafe.Pointer(&buf[0])), len(buf) * 2, cap(buf) * 2}
	return *(*[]byte)(unsafe.Pointer(&sl))
}

func wcslen(str []uint16) int {
	for i := 0; i < len(str); i++ {
		if str[i] == 0 {
			return i
		}
	}
	return len(str)
}


func (deviceInfoSet DevInfo) DeviceRegistryProperty(deviceInfoData *DevInfoData, property SPDRP) (interface{}, error) {
	return SetupDiGetDeviceRegistryProperty(deviceInfoSet, deviceInfoData, property)
}




func SetupDiSetDeviceRegistryProperty(deviceInfoSet DevInfo, deviceInfoData *DevInfoData, property SPDRP, propertyBuffers []byte) error {
	return setupDiSetDeviceRegistryProperty(deviceInfoSet, deviceInfoData, property, &propertyBuffers[0], uint32(len(propertyBuffers)))
}


func (deviceInfoSet DevInfo) SetDeviceRegistryProperty(deviceInfoData *DevInfoData, property SPDRP, propertyBuffers []byte) error {
	return SetupDiSetDeviceRegistryProperty(deviceInfoSet, deviceInfoData, property, propertyBuffers)
}


func (deviceInfoSet DevInfo) SetDeviceRegistryPropertyString(deviceInfoData *DevInfoData, property SPDRP, str string) error {
	str16, err := UTF16FromString(str)
	if err != nil {
		return err
	}
	err = SetupDiSetDeviceRegistryProperty(deviceInfoSet, deviceInfoData, property, utf16ToBuf(append(str16, 0)))
	runtime.KeepAlive(str16)
	return err
}




func SetupDiGetDeviceInstallParams(deviceInfoSet DevInfo, deviceInfoData *DevInfoData) (*DevInstallParams, error) {
	params := &DevInstallParams{}
	params.size = uint32(unsafe.Sizeof(*params))

	return params, setupDiGetDeviceInstallParams(deviceInfoSet, deviceInfoData, params)
}


func (deviceInfoSet DevInfo) DeviceInstallParams(deviceInfoData *DevInfoData) (*DevInstallParams, error) {
	return SetupDiGetDeviceInstallParams(deviceInfoSet, deviceInfoData)
}




func SetupDiGetDeviceInstanceId(deviceInfoSet DevInfo, deviceInfoData *DevInfoData) (string, error) {
	reqSize := uint32(1024)
	for {
		buf := make([]uint16, reqSize)
		err := setupDiGetDeviceInstanceId(deviceInfoSet, deviceInfoData, &buf[0], uint32(len(buf)), &reqSize)
		if err == ERROR_INSUFFICIENT_BUFFER {
			continue
		}
		if err != nil {
			return "", err
		}
		return UTF16ToString(buf), nil
	}
}


func (deviceInfoSet DevInfo) DeviceInstanceID(deviceInfoData *DevInfoData) (string, error) {
	return SetupDiGetDeviceInstanceId(deviceInfoSet, deviceInfoData)
}





func (deviceInfoSet DevInfo) ClassInstallParams(deviceInfoData *DevInfoData, classInstallParams *ClassInstallHeader, classInstallParamsSize uint32, requiredSize *uint32) error {
	return SetupDiGetClassInstallParams(deviceInfoSet, deviceInfoData, classInstallParams, classInstallParamsSize, requiredSize)
}




func (deviceInfoSet DevInfo) SetDeviceInstallParams(deviceInfoData *DevInfoData, deviceInstallParams *DevInstallParams) error {
	return SetupDiSetDeviceInstallParams(deviceInfoSet, deviceInfoData, deviceInstallParams)
}





func (deviceInfoSet DevInfo) SetClassInstallParams(deviceInfoData *DevInfoData, classInstallParams *ClassInstallHeader, classInstallParamsSize uint32) error {
	return SetupDiSetClassInstallParams(deviceInfoSet, deviceInfoData, classInstallParams, classInstallParamsSize)
}




func SetupDiClassNameFromGuidEx(classGUID *GUID, machineName string) (className string, err error) {
	var classNameUTF16 [MAX_CLASS_NAME_LEN]uint16

	var machineNameUTF16 *uint16
	if machineName != "" {
		machineNameUTF16, err = UTF16PtrFromString(machineName)
		if err != nil {
			return
		}
	}

	err = setupDiClassNameFromGuidEx(classGUID, &classNameUTF16[0], MAX_CLASS_NAME_LEN, nil, machineNameUTF16, 0)
	if err != nil {
		return
	}

	className = UTF16ToString(classNameUTF16[:])
	return
}




func SetupDiClassGuidsFromNameEx(className string, machineName string) ([]GUID, error) {
	classNameUTF16, err := UTF16PtrFromString(className)
	if err != nil {
		return nil, err
	}

	var machineNameUTF16 *uint16
	if machineName != "" {
		machineNameUTF16, err = UTF16PtrFromString(machineName)
		if err != nil {
			return nil, err
		}
	}

	reqSize := uint32(4)
	for {
		buf := make([]GUID, reqSize)
		err = setupDiClassGuidsFromNameEx(classNameUTF16, &buf[0], uint32(len(buf)), &reqSize, machineNameUTF16, 0)
		if err == ERROR_INSUFFICIENT_BUFFER {
			continue
		}
		if err != nil {
			return nil, err
		}
		return buf[:reqSize], nil
	}
}




func SetupDiGetSelectedDevice(deviceInfoSet DevInfo) (*DevInfoData, error) {
	data := &DevInfoData{}
	data.size = uint32(unsafe.Sizeof(*data))

	return data, setupDiGetSelectedDevice(deviceInfoSet, data)
}


func (deviceInfoSet DevInfo) SelectedDevice() (*DevInfoData, error) {
	return SetupDiGetSelectedDevice(deviceInfoSet)
}





func (deviceInfoSet DevInfo) SetSelectedDevice(deviceInfoData *DevInfoData) error {
	return SetupDiSetSelectedDevice(deviceInfoSet, deviceInfoData)
}




func SetupUninstallOEMInf(infFileName string, flags SUOI) error {
	infFileName16, err := UTF16PtrFromString(infFileName)
	if err != nil {
		return err
	}
	return setupUninstallOEMInf(infFileName16, flags, 0)
}






func CM_Get_Device_Interface_List(deviceID string, interfaceClass *GUID, flags uint32) ([]string, error) {
	deviceID16, err := UTF16PtrFromString(deviceID)
	if err != nil {
		return nil, err
	}
	var buf []uint16
	var buflen uint32
	for {
		if ret := cm_Get_Device_Interface_List_Size(&buflen, interfaceClass, deviceID16, flags); ret != CR_SUCCESS {
			return nil, ret
		}
		buf = make([]uint16, buflen)
		if ret := cm_Get_Device_Interface_List(interfaceClass, deviceID16, &buf[0], buflen, flags); ret == CR_SUCCESS {
			break
		} else if ret != CR_BUFFER_SMALL {
			return nil, ret
		}
	}
	var interfaces []string
	for i := 0; i < len(buf); {
		j := i + wcslen(buf[i:])
		if i < j {
			interfaces = append(interfaces, UTF16ToString(buf[i:j]))
		}
		i = j + 1
	}
	if interfaces == nil {
		return nil, ERROR_NO_SUCH_DEVICE_INTERFACE
	}
	return interfaces, nil
}



func CM_Get_DevNode_Status(status *uint32, problemNumber *uint32, devInst DEVINST, flags uint32) error {
	ret := cm_Get_DevNode_Status(status, problemNumber, devInst, flags)
	if ret == CR_SUCCESS {
		return nil
	}
	return ret
}
