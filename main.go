package main

import (
	"gitee.com/cristiane/micro-mall-sku-cron/startup"
	"gitee.com/kelvins-io/kelvins"
	"gitee.com/kelvins-io/kelvins/app"
)

const APP_NAME = "micro-mall-sku-cron"

func main() {
	application := &kelvins.CronApplication{
		Application: &kelvins.Application{
			LoadConfig: startup.LoadConfig,
			SetupVars:  startup.SetupVars,
			Name:       APP_NAME,
		},
		GenCronJobs: startup.GenCronJobs,
	}
	app.RunCronApplication(application)
}
