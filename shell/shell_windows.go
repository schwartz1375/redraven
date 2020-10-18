//+build windows !linux !darwin

package shell

import (
	"os/exec"
	"syscall"
)

//SetHide funcation to hide the shell
func SetHide(cmd *exec.Cmd) {
	//cmd.SysProcAttr = &windows.SysProcAttr{HideWindow: true} //"golang.org/x/sys/windows"
	cmd.SysProcAttr = &syscall.SysProcAttr{CreationFlags: 0x08000000} // CREATE_NO_WINDOW
}
