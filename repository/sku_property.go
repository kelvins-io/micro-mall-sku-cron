package repository

import (
	"gitee.com/cristiane/micro-mall-sku-cron/model/mysql"
	"gitee.com/kelvins-io/kelvins"
)

func FindSkuProperty(skuCodeList []string) ([]mysql.SkuProperty, error) {
	var result = make([]mysql.SkuProperty, 0)
	err := kelvins.XORM_DBEngine.Table(mysql.TableSkuProperty).In("code", skuCodeList).Find(&result)
	return result, err
}
