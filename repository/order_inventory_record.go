package repository

import (
	"gitee.com/cristiane/micro-mall-sku-cron/model/mysql"
	"gitee.com/kelvins-io/kelvins"
	"xorm.io/xorm"
)

func FindSkuInventoryRecordList(sqlSelect string, where interface{}) ([]mysql.SkuInventoryRecord, error) {
	var result = make([]mysql.SkuInventoryRecord, 0)
	err := kelvins.XORM_DBEngine.Table(mysql.TableSkuInventoryRecord).Select(sqlSelect).Where(where).Find(&result)
	return result, err
}

func UpdateSkuInventoryRecord(tx *xorm.Session, where, maps interface{}) (int64, error) {
	return tx.Table(mysql.TableSkuInventoryRecord).Where(where).Update(maps)
}

func CreateSkuInventoryRecord(tx *xorm.Session, model *mysql.SkuInventoryRecord) error {
	_, err := tx.Table(mysql.TableSkuInventoryRecord).Insert(model)
	return err
}
