package repository

import (
	"gitee.com/cristiane/micro-mall-sku-cron/model/mysql"
	"gitee.com/kelvins-io/kelvins"
	"xorm.io/xorm"
)

func FindSkuInventoryList(sqlSelect string, where interface{}) ([]mysql.SkuInventory, error) {
	var result = make([]mysql.SkuInventory, 0)
	err := kelvins.XORM_DBEngine.Table(mysql.TableSkuInventory).Select(sqlSelect).Where(where).Find(&result)
	return result, err
}

func GetSkuInventory(tx *xorm.Session, sqlSelect string, where interface{}) (*mysql.SkuInventory, error) {
	var model mysql.SkuInventory
	_, err := tx.Table(mysql.TableSkuInventory).Select(sqlSelect).Where(where).Get(&model)
	return &model, err
}

func UpdateSkuInventory(tx *xorm.Session, where, maps interface{}) (int64, error) {
	return tx.Table(mysql.TableSkuInventory).Where(where).Update(maps)
}
