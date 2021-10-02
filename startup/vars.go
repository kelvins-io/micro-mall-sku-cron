package startup

import (
	"gitee.com/cristiane/micro-mall-sku-cron/model/args"
	"gitee.com/cristiane/micro-mall-sku-cron/vars"
	"gitee.com/kelvins-io/kelvins"
	"gitee.com/kelvins-io/kelvins/setup"
	"gitee.com/kelvins-io/kelvins/util/queue_helper"
)

// SetupVars 加载变量
func SetupVars() error {
	var err error
	err = setupQueueSkuInventorySearchNotice()

	return err
}

func setupQueueSkuInventorySearchNotice() error {
	var err error
	if vars.SkuInventorySearchNoticeSetting != nil {
		vars.SkuInventorySearchNoticeServer, err = setup.NewAMQPQueue(vars.SkuInventorySearchNoticeSetting, nil)
		if err != nil {
			return err
		}
		vars.SkuInventorySearchNoticePusher, err = queue_helper.NewPublishService(
			vars.SkuInventorySearchNoticeServer, &queue_helper.PushMsgTag{
				DeliveryTag:    args.SkuInventorySearchNoticeTag,
				DeliveryErrTag: args.SkuInventorySearchNoticeTagErr,
				RetryCount:     vars.SkuInventorySearchNoticeSetting.TaskRetryCount,
				RetryTimeout:   vars.SkuInventorySearchNoticeSetting.TaskRetryTimeout,
			}, kelvins.BusinessLogger)
		if err != nil {
			return err
		}
	}

	return err
}
