package container

import (
	"os"
	"os/exec"
	"syscall"
)

func NewParentProcess(tty bool, command string) *exec.Cmd {
	args := []string{"init", command}
	// 创建了一个新的进程，新进程将会执行与当前程序相同的可执行文件, 也就是说 CMD 此时也是一个 cli 的客户端
	// 这个新创建的进程（也可以称为子进程）将会执行与当前程序相同的代码，并且带有传递给它的命令行参数。
	// 比我我们执行的是 mydocker run -ti /bin/bash, 此时 args 就是 init /bin/bash, 会触发 initCommand
	cmd := exec.Command("/proc/self/exe", args...)
    cmd.SysProcAttr = &syscall.SysProcAttr{
        Cloneflags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID | syscall.CLONE_NEWNS |
		syscall.CLONE_NEWNET | syscall.CLONE_NEWIPC,
    }
	if tty {
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	}
	return cmd
}