package startup

import (
	"gitee.com/cristiane/micro-mall-sku-cron/vars"
	"gitee.com/kelvins-io/kelvins/config"
	"gitee.com/kelvins-io/kelvins/config/setting"
)

const (
	SectionEmailConfig                    = "email-config"
	OrderFailedInventoryRestoreTaskConfig = "order-failed-inventory-restore-task"
	SectionSkuInventorySearchNotice       = "sku-inventory-search-notice"
	SkuInventorySearchSyncTask            = "sku-inventory-search-notice-task"
)

// LoadConfig 加载配置对象映射
func LoadConfig() error {
	// 加载email数据源
	vars.EmailConfigSetting = new(vars.EmailConfigSettingS)
	config.MapConfig(SectionEmailConfig, vars.EmailConfigSetting)
	// 订单失败恢复库存
	vars.OrderFailedInventoryRestoreTaskSetting = new(vars.OrderFailedInventoryRestoreTaskSettingS)
	config.MapConfig(OrderFailedInventoryRestoreTaskConfig, vars.OrderFailedInventoryRestoreTaskSetting)
	// 商品库存搜素通知
	vars.SkuInventorySearchNoticeSetting = new(setting.QueueAMQPSettingS)
	config.MapConfig(SectionSkuInventorySearchNotice, vars.SkuInventorySearchNoticeSetting)
	// 商品库存搜索同步
	vars.SkuInventorySearchSyncTaskSetting = new(vars.SkuInventorySearchSyncTaskSettingS)
	config.MapConfig(SkuInventorySearchSyncTask, vars.SkuInventorySearchSyncTaskSetting)

	return nil
}
