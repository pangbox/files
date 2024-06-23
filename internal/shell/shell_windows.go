package shell

import "golang.org/x/sys/windows"

func OpenFolder(folder string) error {
	return windows.ShellExecute(windows.Handle(0), windows.StringToUTF16Ptr("explore"), windows.StringToUTF16Ptr(folder), nil, nil, windows.SW_SHOWNORMAL)
}
