package win32

import (
	"os"
	"runtime"
	"syscall"
	"unsafe"
)

func CreateJobObject(sa *syscall.SecurityAttributes, name *uint16) (syscall.Handle, error) {
	r1, _, e1 := procCreateJobObjectW.Call(
		uintptr(unsafe.Pointer(sa)),
		uintptr(unsafe.Pointer(name)))
	runtime.KeepAlive(sa)
	runtime.KeepAlive(name)
	if int(r1) == 0 {
		return syscall.InvalidHandle, os.NewSyscallError("CreateJobObject", e1)
	}
	return syscall.Handle(r1), nil
}

func QueryInformationJobObject(job syscall.Handle, infoclass uint32, info unsafe.Pointer, length uint32) (uint32, error) {
	var nLength uint32
	r1, _, e1 := procQueryInformationJobObject.Call(
		uintptr(job),
		uintptr(infoclass),
		uintptr(info),
		uintptr(length),
		uintptr(unsafe.Pointer(&nLength)))
	runtime.KeepAlive(&nLength)
	if int(r1) == 0 {
		return nLength, os.NewSyscallError("QueryInformationJobObject", e1)
	}
	return nLength, nil
}

func SetInformationJobObject(job syscall.Handle, infoclass uint32, info unsafe.Pointer, length uint32) error {
	r1, _, e1 := procSetInformationJobObject.Call(
		uintptr(job),
		uintptr(infoclass),
		uintptr(info),
		uintptr(length))
	if int(r1) == 0 {
		return os.NewSyscallError("SetInformationJobObject", e1)
	}
	return nil
}

type JobObjectBasicAccountingInformation struct {
	TotalUserTime             uint64
	TotalKernelTime           uint64
	ThisPeriodTotalUserTime   uint64
	ThisPeriodTotalKernelTime uint64
	TotalPageFaultCount       uint32
	TotalProcesses            uint32
	ActiveProcesses           uint32
	TotalTerminatedProcesses  uint32
}

type JobObjectBasicUiRestrictions struct {
	UIRestrictionClass uint32
}

const (
	JOB_OBJECT_UILIMIT_DESKTOP          = 0x40
	JOB_OBJECT_UILIMIT_DISPLAYSETTINGS  = 0x10
	JOB_OBJECT_UILIMIT_EXITWINDOWS      = 0x80
	JOB_OBJECT_UILIMIT_GLOBALATOMS      = 0x20
	JOB_OBJECT_UILIMIT_HANDLES          = 1
	JOB_OBJECT_UILIMIT_READCLIPBOARD    = 2
	JOB_OBJECT_UILIMIT_SYSTEMPARAMETERS = 8
	JOB_OBJECT_UILIMIT_WRITECLIPBOARD   = 4
)

func GetJobObjectBasicAccountingInformation(job syscall.Handle) (*JobObjectBasicAccountingInformation, error) {
	var jinfo JobObjectBasicAccountingInformation
	_, err := QueryInformationJobObject(job, 1, unsafe.Pointer(&jinfo), uint32(unsafe.Sizeof(jinfo)))
	if err != nil {
		return nil, err
	}
	return &jinfo, nil
}

type JobObjectBasicLimitInformation struct {
	PerProcessUserTimeLimit uint64  // LARGE_INTEGER
	PerJobUserTimeLimit     uint64  // LARGE_INTEGER
	LimitFlags              uint32  // DWORD
	MinimumWorkingSetSize   uintptr // SIZE_T
	MaximumWorkingSetSize   uintptr // SIZE_T
	ActiveProcessLimit      uint32  // DWORD
	Affinity                uintptr // ULONG_PTR
	PriorityClass           uint32  // DWORD
	SchedulingClass         uint32  // DWORD
}

const (
	JOB_OBJECT_LIMIT_KILL_ON_JOB_CLOSE          = 0x2000
	JOB_OBJECT_LIMIT_DIE_ON_UNHANDLED_EXCEPTION = 0x400
	JOB_OBJECT_LIMIT_ACTIVE_PROCESS             = 8
	JOB_OBJECT_LIMIT_JOB_MEMORY                 = 0x200
	JOB_OBJECT_LIMIT_JOB_TIME                   = 4
	JOB_OBJECT_LIMIT_PROCESS_MEMORY             = 0x100
	JOB_OBJECT_LIMIT_PROCESS_TIME               = 2
	JOB_OBJECT_LIMIT_WORKINGSET                 = 1
	JOB_OBJECT_LIMIT_AFFINITY                   = 0x00000010
)

type IoCounters struct {
	ReadOperationCount  uint64 // ULONGLONG
	WriteOperationCount uint64 // ULONGLONG
	OtherOperationCount uint64 // ULONGLONG
	ReadTransferCount   uint64 // ULONGLONG
	WriteTransferCount  uint64 // ULONGLONG
	OtherTransferCount  uint64 // ULONGLONG
}

func GetJobObjectExtendedLimitInformation(job syscall.Handle) (*JobObjectExtendedLimitInformation, error) {
	var jinfo JobObjectExtendedLimitInformation
	_, err := QueryInformationJobObject(job, 9, unsafe.Pointer(&jinfo), uint32(unsafe.Sizeof(jinfo)))
	if err != nil {
		return nil, err
	}
	return &jinfo, nil
}

func SetJobObjectBasicUiRestrictions(job syscall.Handle, info *JobObjectBasicUiRestrictions) error {
	return SetInformationJobObject(job, 4, unsafe.Pointer(info), uint32(unsafe.Sizeof(*info)))
}

func SetJobObjectExtendedLimitInformation(job syscall.Handle, info *JobObjectExtendedLimitInformation) error {
	return SetInformationJobObject(job, 9, unsafe.Pointer(info), uint32(unsafe.Sizeof(*info)))
}

func AssignProcessToJobObject(job syscall.Handle, process syscall.Handle) error {
	r1, _, e1 := procAssignProcessToJobObject.Call(
		uintptr(job),
		uintptr(process))
	if int(r1) == 0 {
		return os.NewSyscallError("AssignProcessToJobObject", e1)
	}
	return nil
}

const (
	MEM_COMMIT     = 0x00001000
	PAGE_READWRITE = 0x04
)
