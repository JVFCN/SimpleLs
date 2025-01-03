package TraverseDIR

import (
	"os"
	"syscall"
)

func isWindowsHidden(info os.DirEntry) bool {
	const FileAttributeHidden = 0x2

	// 调用Windows API获取文件/文件夹属性
	namePtr, err := syscall.UTF16PtrFromString(info.Name())
	if err != nil {
		return false
	}

	winAttrs, err := syscall.GetFileAttributes(namePtr)
	if err != nil {
		return false
	}

	return winAttrs&FileAttributeHidden != 0
}
func isHidden(info os.DirEntry) bool {
	// 检查操作系统是否认为文件/文件夹是隐藏的
	if info.Name()[0] == '.' || isWindowsHidden(info) {
		return true
	}
	return false
}
