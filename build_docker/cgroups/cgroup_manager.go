package cgroups

import (
	"com.meiken/meiken-docker/build_docker/cgroups/subsystems"
	"github.com/sirupsen/logrus"
)

// 需要将不同 subsystem 中的 cgroup 管理起来
type CgroupManager struct {
	// cgroup在 hierarchy 中的路径 相当于创建的 cgroup 目录相对于 root cgroup 目录的路径
	Path string
	// 资源配置
	Resource *subsystems.ResourceConfig
}

func NewCgroupManager(path string) *CgroupManager {
	return &CgroupManager{
		Path: path,
	}
}

// 将进程 pid 加入到每一个 cgroup 中
func (c *CgroupManager) Apply(pid int) error {
	for _, subSysIns := range(subsystems.SubsystemsIns) {
		subSysIns.Apply(c.Path, pid)
	}
	return nil
}

// 设置各个 subsystem 挂载中的 cgroup 资源限制
func (c *CgroupManager) Set(res *subsystems.ResourceConfig) error {
	for _, subSysIns := range(subsystems.SubsystemsIns) {
		subSysIns.Set(c.Path, res)
	}
	return nil
}

// 释放 cgroup, 释放各个 subsystem 挂载中的 cgroup
func (c *CgroupManager) Destroy() error {
	for _, subSysIns := range(subsystems.SubsystemsIns) {
		if err := subSysIns.Remove(c.Path); err != nil {
			logrus.Warnf("remove cgroup fail %v", err)
		}
	}
	return nil
}