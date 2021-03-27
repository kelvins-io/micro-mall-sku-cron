package startup

import (
	"gitee.com/cristiane/micro-mall-sku-cron/vars"
	"gitee.com/kelvins-io/kelvins/config"
	"log"
)

const (
	SectionEmailConfig                    = "email-config"
	OrderFailedInventoryRestoreTaskConfig = "order-failed-inventory-restore-task"
)

// LoadConfig 加载配置对象映射
func LoadConfig() error {
	// 加载email数据源
	log.Printf("[info] Load default config %s", SectionEmailConfig)
	vars.EmailConfigSetting = new(vars.EmailConfigSettingS)
	config.MapConfig(SectionEmailConfig, vars.EmailConfigSetting)
	// 订单失败恢复库存
	log.Printf("[info] Load default config %s", OrderFailedInventoryRestoreTaskConfig)
	vars.OrderFailedInventoryRestoreTaskSetting = new(vars.OrderFailedInventoryRestoreTaskSettingS)
	config.MapConfig(OrderFailedInventoryRestoreTaskConfig, vars.OrderFailedInventoryRestoreTaskSetting)
	return nil
}
