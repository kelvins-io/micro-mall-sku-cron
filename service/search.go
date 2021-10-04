package service

import (
	"context"
	"gitee.com/cristiane/micro-mall-sku-cron/model/args"
	"gitee.com/cristiane/micro-mall-sku-cron/model/mysql"
	"gitee.com/cristiane/micro-mall-sku-cron/repository"
	"gitee.com/cristiane/micro-mall-sku-cron/vars"
	"gitee.com/kelvins-io/common/json"
	"gitee.com/kelvins-io/kelvins"
	"github.com/google/uuid"
)

var (
	pageSize = 50
	pageNum  = 1
)

func SkuInventorySearchSync() {
	if vars.SkuInventorySearchSyncTaskSetting != nil {
		if vars.SkuInventorySearchSyncTaskSetting.SingleSyncNum > 0 {
			pageSize = vars.SkuInventorySearchSyncTaskSetting.SingleSyncNum
		}
	}
	count := 0
	for {
		if count > 5 {
			break
		}
		skuInventorySearchSyncOne(pageSize, pageNum)
		count++
		pageNum++
	}
}

const (
	sqlSelectSkuInventorySearch = "sku_code,shop_id,price"
)

func skuInventorySearchSyncOne(pageSize, pageNum int) {
	ctx := context.TODO()
	where := map[string]interface{}{}
	list, err := repository.FindSkuInventoryList(sqlSelectSkuInventorySearch, where, pageSize, pageNum)
	if err != nil {
		kelvins.ErrLogger.Errorf(ctx, "FindSkuInventoryList err %v", err)
		return
	}
	if len(list) == 0 {
		return
	}
	skuCodeList := make([]string, len(list))
	for i := 0; i < len(list); i++ {
		skuCodeList[i] = list[i].SkuCode
	}
	skuPropertyList, err := repository.FindSkuProperty(skuCodeList)
	if err != nil {
		kelvins.ErrLogger.Errorf(ctx, "FindSkuProperty err %v", err)
		return
	}
	if len(skuPropertyList) == 0 {
		return
	}
	skuCodeToProperty := map[string]mysql.SkuProperty{}
	for i := 0; i < len(skuPropertyList); i++ {
		skuCodeToProperty[skuPropertyList[i].Code] = skuPropertyList[i]
	}
	for i := 0; i < len(list); i++ {
		propertyInfo, ok := skuCodeToProperty[list[i].SkuCode]
		if ok {
			info := &args.SkuInventoryInfo{
				ShopId:        list[i].ShopId,
				SkuCode:       list[i].SkuCode,
				Name:          propertyInfo.Name,
				Price:         list[i].Price,
				Title:         propertyInfo.Title,
				SubTitle:      propertyInfo.SubTitle,
				Desc:          propertyInfo.Desc,
				Production:    propertyInfo.Production,
				Supplier:      propertyInfo.Supplier,
				Category:      int32(propertyInfo.Category),
				Color:         propertyInfo.Color,
				ColorCode:     int32(propertyInfo.ColorCode),
				Specification: propertyInfo.Specification,
				DescLink:      propertyInfo.DescLink,
			}
			_ = skuInventorySearchNotice(info)
		}
	}
}

func skuInventorySearchNotice(info *args.SkuInventoryInfo) error {
	kelvins.GPool.SendJob(func() {
		var ctx = context.TODO()
		msg := &args.CommonBusinessMsg{
			Type:    args.SkuInventorySearchNotice,
			Tag:     args.GetMsg(args.SkuInventorySearchNotice),
			UUID:    uuid.New().String(),
			Content: json.MarshalToStringNoError(info),
		}
		vars.SkuInventorySearchNoticePusher.PushMessage(ctx, msg)
	})
	return nil
}
