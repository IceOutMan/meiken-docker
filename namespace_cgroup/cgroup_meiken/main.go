package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"strconv"
	"syscall"
)

// 挂载了memory subsystem 的 hierarchy 的目录位置
const cgroupMemoryHierarchyMount = "/sys/fs/cgroup/memory"

func main() {
	if os.Args[0] == "/proc/self/exe" {
				// fork出的容器进程 - 进行压测
		fmt.Printf("current pid %d", syscall.Getpid())
		fmt.Println()
		cmd := exec.Command("sh", "-c", `stress --vm-bytes 90m --vm-keep -m 1`)
		cmd.SysProcAttr = &syscall.SysProcAttr{}
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		if err := cmd.Run(); err != nil {
			fmt.Println(err)
			fmt.Println("fork stress error")
			fmt.Println(err)
			os.Exit(-1)
		}
	}
		// 执行 - 启动当前进程的可执行文件 - go run main.go 会运行到此处
	cmd := exec.Command("/proc/self/exe")
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID | syscall.CLONE_NEWNS,
	}
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Start(); err != nil {
		fmt.Println("ERROR", err)
		os.Exit(-1)
	} else {
		// 得到 fork 出来进程映射在外部命名空间的pid
		fmt.Printf("%v\n", cmd.Process.Pid)

		// 在系统默认创建挂载 memory subsystem 的 Hierarchy上创建cgroup
		os.Mkdir(path.Join(cgroupMemoryHierarchyMount, "test-memory-limit"), 0755)
		// 将容器进程加入到这个 cgroup
		ioutil.WriteFile(path.Join(cgroupMemoryHierarchyMount, "test-memory-limit", "tasks"), []byte(strconv.Itoa(cmd.Process.Pid)), 0644)
		// 限制cgroup进程使用
		ioutil.WriteFile(path.Join(cgroupMemoryHierarchyMount, "test-memory-limit", "memory.limit_in_bytes"), []byte("100m"), 0644)
	}
	cmd.Process.Wait()
}
