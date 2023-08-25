package container

import (
	"os"
	"os/exec"
	"syscall"

	"github.com/sirupsen/logrus"
)

func NewParentProcess(tty bool) (*exec.Cmd, *os.File) {
	// 构建管道
	readPipe, writePipe, err := NewPipe()
	if err != nil {
		logrus.Errorf("New pipe error %v", err)
	}

	/*
		创建了一个新的进程，新进程将会执行与当前程序相同的可执行文件, 也就是说 CMD 此时也是一个 cli 的客户端
		这个新创建的进程（也可以称为子进程）将会执行与当前程序相同的代码，并且带有传递给它的命令行参数。
		比我我们执行的是 mydocker run -ti /bin/bash, 此时 args 就是 init /bin/bash, 会触发 initCommand
	*/
	cmd := exec.Command("/proc/self/exe", "init")
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID | syscall.CLONE_NEWNS |
			syscall.CLONE_NEWNET | syscall.CLONE_NEWIPC,
	}
	if tty {
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	}
	cmd.ExtraFiles = []*os.File{readPipe}
	cmd.Dir = "/root/busybox"
	// 新的init 进程 保留一个管道的读端
	// 管道的写端暴露出去，相当于是 init 进程用来通过管道接收外部的命令
	return cmd, writePipe
}

func NewPipe() (*os.File, *os.File, error) {
	read, write, err := os.Pipe()
	if err != nil {
		return nil, nil, err
	}
	return read, write, nil
}
