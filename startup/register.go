package startup

import (
	"gitee.com/cristiane/micro-mall-sku-cron/service"
	"gitee.com/cristiane/micro-mall-sku-cron/vars"
	"gitee.com/kelvins-io/kelvins"
)

func GenCronJobs() []*kelvins.CronJob {
	tasks := make([]*kelvins.CronJob, 0)
	if vars.OrderFailedInventoryRestoreTaskSetting != nil {
		if vars.OrderFailedInventoryRestoreTaskSetting.Cron != "" {
			tasks = append(tasks, &kelvins.CronJob{
				Name: "库存恢复-订单创建失败",
				Spec: vars.OrderFailedInventoryRestoreTaskSetting.Cron,
				Job:  service.HandleOrderFailedSkuInventoryRestore,
			})
		}
	}

	return tasks
}
