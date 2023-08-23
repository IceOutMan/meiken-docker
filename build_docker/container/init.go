package container

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"syscall"

	"github.com/sirupsen/logrus"
)

func RunContainerInitProcess() error {

	// 读到管道传递进来的命令
	cmdArray := readUserCommand()
	if len(cmdArray) == 0 {
		return fmt.Errorf(" Run container get user command error, cmdArray is nil")
	}

	// 初始化内容：挂载 proc
	defaultMountFlags := syscall.MS_NOEXEC | syscall.MS_NOSUID | syscall.MS_NODEV
	syscall.Mount("proc", "/proc", "proc", uintptr(defaultMountFlags), "")

	path, err := exec.LookPath(cmdArray[0])
	if err != nil {
		logrus.Errorf("Exec loop path error %v", err)
		return err
	}
	logrus.Infof("Find Path %s", path)
	/*
	* 要执行的用户命令,execve 系统调用会在当前进程中执行命令，使用当前进程的资源
	* 相当于是要执行的任务使用了当前进程的外壳，夺舍了当前进程
	 */
	if err := syscall.Exec(path, cmdArray[0:], os.Environ()); err != nil {
		logrus.Errorf(err.Error())
	}
	return nil
}

// 从 pipe 中读出来 命令
func readUserCommand() []string {
	pipe := os.NewFile(uintptr(3), "pipe")
	msg, err := ioutil.ReadAll(pipe)
	if err != nil {
		logrus.Errorf("init read pipe error %v", err)
		return nil
	}

	msgStr := string(msg)
	return strings.Split(msgStr, " ")
}

// // Init 挂载点
// func setUpMount(){
// 	pwd, err := os.Getwd()
// 	if err != nil{
// 		logrus.Errorf("Get current location error %v", err)
// 		return
// 	}
// 	logrus.Infof("Current location is %s", pwd)
// 	privotRoot(pwd)
// }

// func privotRoot(root string) error {
// 	/**
// 	 为了当前root 的老root和新root 不在同一个文件系统下，我们把root重新mount了一次
// 	 bind mount 是把相同的内容换了一个挂载点的挂载方法
// 	*/

// 	if err := syscall.Mount(root, root, "bind", syscall.MS_BIND|syscall.MS_REC, ""); err != nil{
// 		return fmt.Errorf("Moutn rootfs to itself error : %v", err)
// 	}

// 	// 创建 rootfs/.pivot_root 存储 old_root
// 	pivotDir := filepath.Join(root, ".pivot_root")
// 	if err := os.Mkdir(pivotDir, 0777); err != nil{
// 		return err
// 	}
// 	// pivot_root 到新的rootfs，现在老的 old_root 是挂载在 rootfs/.pivot_root
// 	// 挂载点现在依然可以在 mount 命令中看到
// 	if err := syscall.PivotRoot(root, pivotDir); err != nil{
// 		return fmt.Errorf("pivot_root %v", err)
// 	}
// 	// 修改当前的工作目录到根目录
// 	if err := syscall.Chdir("/"); err != nil{
// 		return fmt.Errorf("chdir / %v", err)
// 	}

// 	pivotDir = filepath.Join("/", ".pivot_root")
// 	// umount rootfs/.pivot_root
// 	if err := syscall.Unmount(pivotDir, syscall.MNT_DETACH); err != nil{
// 		return fmt.Errorf("unmount pivot_root dir %v", err)
// 	}

// 	// 删除临时文件
// 	return os.Remove(pivotDir)
// }
