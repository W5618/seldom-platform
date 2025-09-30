package services

var (
	// GlobalScheduler 全局调度服务实例
	GlobalScheduler *SchedulerService
)

// InitGlobalScheduler 初始化全局调度服务
func InitGlobalScheduler() error {
	GlobalScheduler = NewSchedulerService()
	return GlobalScheduler.Start()
}

// StopGlobalScheduler 停止全局调度服务
func StopGlobalScheduler() {
	if GlobalScheduler != nil {
		GlobalScheduler.Stop()
	}
}