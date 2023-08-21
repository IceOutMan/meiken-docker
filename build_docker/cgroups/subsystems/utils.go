package subsystems

import (
	"bufio"
	"fmt"
	"os"
	"path"
	"strings"
)

// 找出某个 subsystem 的 hierarchy cgroup 根节点所在的目录
func FindCgroupMountpoint(subsystem string) string {
	/*
     * 它包含有关当前进程的挂载信息
	 * /proc/self/mountinfo 文件内容的格式
	 * 45 32 0:40 / /sys/fs/cgroup/memory rw,nosuid,nodev,noexec,relatime shared:23 - cgroup cgroup rw,memory
	 * 如： subsystem = memory -> 匹配到末尾的 memory -> 就返回固定的位置的信息是对应的path是： /sys/fs/cgroup/memory 
	*/
	f, err := os.Open("/proc/self/mountinfo")
	if err != nil {
		return ""
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		txt := scanner.Text()
		fields := strings.Split(txt, " ")
		for _, opt := range strings.Split(fields[len(fields)-1], ",") {
			if opt == subsystem {
				return fields[4]
			}
		}
	}
	if err := scanner.Err(); err != nil {
		return ""
	}

	return ""
}

// 得到 cgroup 在文件系统中的绝对路径, 比如 subsystem=memory, cgroupPath = docker_1
// 这里的 cgroup 是建立在系统的 hierarchy 中, 得到就是 /sys/fs/cgroup/memory/docker_1
func GetCgroupPath(subsystem string, cgroupPath string, autoCreate bool) (string, error) {
	cgroupRoot := FindCgroupMountpoint(subsystem)
	if _, err := os.Stat(path.Join(cgroupRoot, cgroupPath)); err == nil || (autoCreate && os.IsNotExist(err)) {
		if os.IsNotExist(err) {
			if err := os.Mkdir(path.Join(cgroupRoot, cgroupPath), 0755); err == nil {
			} else {
				return "", fmt.Errorf("error create cgroup %v", err)
			}
		}
		return path.Join(cgroupRoot, cgroupPath), nil
	} else {
		return "", fmt.Errorf("cgroup path error %v", err)
	}
}