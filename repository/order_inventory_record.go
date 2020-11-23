package repository

import (
	"gitee.com/cristiane/micro-mall-sku-cron/model/mysql"
	"gitee.com/kelvins-io/kelvins"
	"xorm.io/xorm"
)

func FindSkuInventoryRecordList(sqlSelect string, where interface{}, pageSize, pageNum int) ([]mysql.SkuInventoryRecord, error) {
	var result = make([]mysql.SkuInventoryRecord, 0)
	session := kelvins.XORM_DBEngine.Table(mysql.TableSkuInventoryRecord).Select(sqlSelect).Where(where)
	if pageSize > 0 && pageNum > 0 {
		session = session.Limit(pageSize, (pageNum-1)*pageSize)
	}
	err := session.Find(&result)
	return result, err
}

func UpdateSkuInventoryRecordByTx(tx *xorm.Session, where, maps interface{}) (int64, error) {
	return tx.Table(mysql.TableSkuInventoryRecord).Where(where).Update(maps)
}

func UpdateSkuInventoryRecord(where, maps interface{}) (int64, error) {
	return kelvins.XORM_DBEngine.Table(mysql.TableSkuInventoryRecord).Where(where).Update(maps)
}

func CreateSkuInventoryRecord(tx *xorm.Session, model *mysql.SkuInventoryRecord) error {
	_, err := tx.Table(mysql.TableSkuInventoryRecord).Insert(model)
	return err
}
