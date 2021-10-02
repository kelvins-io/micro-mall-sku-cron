package vars

import (
	"gitee.com/kelvins-io/common/queue"
	"gitee.com/kelvins-io/kelvins/config/setting"
	"gitee.com/kelvins-io/kelvins/util/queue_helper"
)

var (
	EmailConfigSetting                     *EmailConfigSettingS
	OrderFailedInventoryRestoreTaskSetting *OrderFailedInventoryRestoreTaskSettingS
	SkuInventorySearchSyncTaskSetting      *SkuInventorySearchSyncTaskSettingS
	SkuInventorySearchNoticeSetting        *setting.QueueAMQPSettingS
	SkuInventorySearchNoticeServer         *queue.MachineryQueue
	SkuInventorySearchNoticePusher         *queue_helper.PublishService
)
