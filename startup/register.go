package startup

import (
	"gitee.com/cristiane/micro-mall-sku-cron/service"
	"gitee.com/kelvins-io/kelvins"
)

const (
	CronHandleOrderFailedInventoryRestore  = "30 */10 * * * *"
	CronHandleOrderSuccessInventoryRestore = "30 */2 * * * *"
)

func GenCronJobs() []*kelvins.CronJob {
	tasks := make([]*kelvins.CronJob, 0)
	tasks = append(tasks, &kelvins.CronJob{
		Name: "库存恢复-订单创建失败",
		Spec: CronHandleOrderFailedInventoryRestore,
		Job:  service.HandleOrderFailedSkuInventoryRestore,
	})
	tasks = append(tasks, &kelvins.CronJob{
		Name: "库存恢复-订单创建成功",
		Spec: CronHandleOrderSuccessInventoryRestore,
		Job:  service.HandleOrderSuccessSkuInventoryRestore,
	})
	return tasks
}
