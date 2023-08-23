package main

import (
	"os"
	"strings"

	"com.meiken/meiken-docker/build_docker/container"
	"github.com/sirupsen/logrus"
)

func Run(tty bool, comArry []string) {
	parent, writePipe := container.NewParentProcess(tty)
	if parent == nil {
		logrus.Errorf("New parent process error")
		return
	}
	if err := parent.Start(); err != nil {
		logrus.Error(err)
	}
	sendInitCommand(comArry, writePipe)
	parent.Wait()
	os.Exit(0)
}

// 命令数据合并成一个字符串，然后通过 writePipe 写入给 init 进程的管道
func sendInitCommand(comArray []string, writePipe *os.File) {
	command := strings.Join(comArray, " ")
	logrus.Infof("command all is %s", command)

	writePipe.WriteString(command)
	writePipe.Close()
}
