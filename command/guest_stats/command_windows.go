/*
guest-stats - returns disk and memory stats from guest

Example:
        { "execute": "guest-stats", "arguments": {}}
*/
package guest_stats

import "github.com/vtolstov/qemu-ga/qga"

func init() {
	qga.RegisterCommand(&qga.Command{
		Name:    "guest-stats",
		Func:    fnGuestStats,
		Enabled: true,
		Returns: true,
	})
}

func fnGuestStats(req *qga.Request) *qga.Response {
	res := &qga.Response{Id: req.Id}
	/*
		resData := struct {
			MemoryTotal uint64
			MemoryFree  uint64
			SwapTotal   uint64
			SwapFree    uint64
			BlkTotal    uint64
			BlkFree     uint64
			InodeTotal  uint64
			InodeFree   uint64
		}{}

		kernel32, err := syscall.LoadLibrary("Kernel32.dll")
		if err != nil {
			res.Error = &qga.Error{Code: -1, Desc: err.Error()}
			return res
		}
		defer syscall.FreeLibrary(kernel32)
		GetDiskFreeSpaceEx, err := syscall.GetProcAddress(syscall.Handle(kernel32), "GetDiskFreeSpaceExW")
		if err != nil {
			res.Error = &qga.Error{Code: -1, Desc: err.Error()}
			return res
		}

		lpFreeBytesAvailable := int64(0)
		lpTotalNumberOfBytes := int64(0)
		lpTotalNumberOfFreeBytes := int64(0)

		r, a, b := syscall.Syscall6(uintptr(GetDiskFreeSpaceEx), 4,
			uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr("C:"))),
			uintptr(unsafe.Pointer(&lpFreeBytesAvailable)),
			uintptr(unsafe.Pointer(&lpTotalNumberOfBytes)),
			uintptr(unsafe.Pointer(&lpTotalNumberOfFreeBytes)), 0, 0)

		resData.BlkTotal = lpTotalNumberOfBytes
		resData.BlkFree = lpTotalNumberOfFreeBytes

		f, err := windows.Open("C:")
		if err != nil {
			res.Error = &qga.Error{Code: -1, Desc: err.Error()}
			return res
		}
		defer f.Close()

		var d syscall.ByHandleFileInformation
		err := syscall.GetFileInformationByHandle(syscall.Handle(f.Fd()), &d)
		if err != nil {
			res.Error = &qga.Error{Code: -1, Desc: err.Error()}
			return res
		}

		resData.InodeTotal = d.FileIndexHigh
		resData.InodeFree = d.FileIndexLow

		m := struct {
		}{}

		GlobalMemoryStatusEx, err := syscall.GetProcAddress(syscall.Handle(kernel32), "GlobalMemoryStatusExW")
		if err != nil {
			res.Error = &qga.Error{Code: -1, Desc: err.Error()}
			return res
		}



		res.Return = resData
	*/
	return res
}
