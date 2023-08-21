package subsystems

type ResourceConfig struct{
	MemoryLimit string
	CpuShare string
	CpuSet string
} 

// 这里将 cgroup 抽象成了path, 原因是 cgroup 在 hierarchy 的路径，便是虚拟文件系统中的虚拟路径
type Subsystem interface{
	// 返回 subsystem 的名字
	Name() string
	// 设置某个 cgroup 在 subsystem 中的资源限制
	Set(path string, res *ResourceConfig) error
	// 将进程添加到 cgroup 中
	Apply(path string, pid int ) error
	// 移除某个 cgroup
	Remove(path string) error
} 

var (
	// 通过不同的 subsystem 初始化实例创建资源限制处理链数组
	SubsystemsIns = []Subsystem{
		&CpusetSubSystem{},
		&MemorySubSystem{},
		&CpuSubSystem{},
	}
)
