package main

import (
	"os"

	"com.meiken/meiken-docker/build_docker/container"
	log "github.com/sirupsen/logrus"
)


func Run(tty bool, command string) {
	parent := container.NewParentProcess(tty, command)
	if err := parent.Start(); err != nil {
		log.Error(err)
	}
	parent.Wait()
	os.Exit(-1)
}