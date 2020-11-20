package service

import (
	"context"
	"gitee.com/cristiane/micro-mall-sku-cron/model/args"
	"gitee.com/cristiane/micro-mall-sku-cron/model/mysql"
	"gitee.com/cristiane/micro-mall-sku-cron/pkg/util"
	"gitee.com/cristiane/micro-mall-sku-cron/proto/micro_mall_order_proto/order_business"
	"gitee.com/cristiane/micro-mall-sku-cron/repository"
	"gitee.com/kelvins-io/kelvins"
	"time"
)

const (
	sqlSelectSkuInventoryRestore = "op_tx_id"
	sqlSelectSkuInventoryRecord  = "shop_id,sku_code,amount,update_time,op_tx_id"
)

func HandleSkuInventoryRestore() {
	ctx := context.Background()
	where := map[string]interface{}{
		"verify":  0, // 记录未经验证
		"op_type": 1, // 出库
	}
	recordList, err := repository.FindSkuInventoryRecordList(sqlSelectSkuInventoryRestore, where)
	if err != nil {
		kelvins.ErrLogger.Errorf(ctx, "FindSkuInventoryRecordList,err: %v, req: %+v", err, where)
		return
	}
	if len(recordList) == 0 {
		return
	}
	opTxIds := make([]string, len(recordList))
	for i := 0; i < len(recordList); i++ {
		opTxIds[i] = recordList[i].OpTxId
	}
	serverName := args.RpcServiceMicroMallOrder
	conn, err := util.GetGrpcClient(serverName)
	if err != nil {
		kelvins.ErrLogger.Errorf(ctx, "GetGrpcClient %v,err: %v", serverName, err)
		return
	}
	defer conn.Close()
	client := order_business.NewOrderBusinessServiceClient(conn)
	req := order_business.CheckOrderStateRequest{OrderCodes: opTxIds}
	rsp, err := client.CheckOrderState(ctx, &req)
	if err != nil {
		kelvins.ErrLogger.Errorf(ctx, "CheckOrderStateRequest %v,err: %v, req: %+v", serverName, err, req)
		return
	}
	if rsp.Common.Code != order_business.RetCode_SUCCESS {
		kelvins.ErrLogger.Errorf(ctx, "CheckOrderStateRequest %v,err: %v, req: %+v, rsp: %+v", serverName, err, req, rsp)
		return
	}
	if len(rsp.List) == 0 {
		return
	}
	skuInventoryFailedOrder := make([]string, 0)
	for i := 0; i < len(rsp.List); i++ {
		if !rsp.List[i].IsExist {
			skuInventoryFailedOrder = append(skuInventoryFailedOrder, rsp.List[i].OrderCode)
		}
	}
	if len(skuInventoryFailedOrder) == 0 {
		return
	}
	skuInventoryRecordWhere := map[string]interface{}{
		"op_tx_id": skuInventoryFailedOrder,
	}
	skuInventoryRecordList, err := repository.FindSkuInventoryRecordList(sqlSelectSkuInventoryRecord, skuInventoryRecordWhere)
	if err != nil {
		kelvins.ErrLogger.Errorf(ctx, "FindSkuInventoryRecordList err: %v, where: %+v", err, skuInventoryRecordWhere)
		return
	}
	if len(skuInventoryRecordList) == 0 {
		return
	}
	tx := kelvins.XORM_DBEngine.NewSession()
	err = tx.Begin()
	if err != nil {
		kelvins.ErrLogger.Errorf(ctx, "HandleSkuInventoryRestore Begin err: %v, ", err)
		return
	}
	for i := 0; i < len(skuInventoryRecordList); i++ {
		row := skuInventoryRecordList[i]
		getSkuInventoryWhere := map[string]interface{}{
			"shop_id":  row.ShopId,
			"sku_code": row.SkuCode,
		}
		skuInventory, err := repository.GetSkuInventory(tx, "*", getSkuInventoryWhere)
		if err != nil {
			err = tx.Rollback()
			if err != nil {
				kelvins.ErrLogger.Errorf(ctx, "GetSkuInventory  Rollback err: %v, ", err)
			}
			kelvins.ErrLogger.Errorf(ctx, "GetSkuInventory  Rollback err: %v, ", err)
		}
		if skuInventory.SkuCode == "" {
			continue
		}
		// 记录库存日志
		inventoryRecordWhere := map[string]interface{}{
			"op_tx_id":    row.OpTxId,
			"update_time": row.UpdateTime,
		}
		inventoryRecordMaps := map[string]interface{}{
			"verify":      1,
			"update_time": time.Now(),
		}
		rowAffected, err := repository.UpdateSkuInventoryRecord(tx, inventoryRecordWhere, inventoryRecordMaps)
		if err != nil {
			err = tx.Rollback()
			if err != nil {
				kelvins.ErrLogger.Errorf(ctx, "HandleSkuInventoryRestore  Rollback err: %v, ", err)
			}
			kelvins.ErrLogger.Errorf(ctx, "UpdateSkuInventoryRecord err: %v, where: %+v,maps: %+v ", err, inventoryRecordWhere, inventoryRecordMaps)
			return
		}
		if rowAffected != 1 {
			err = tx.Rollback()
			if err != nil {
				kelvins.ErrLogger.Errorf(ctx, "HandleSkuInventoryRestore rowAffected Rollback err: %v, ", err)
			}
			return
		}
		// 增加库存恢复记录
		inventoryRestoreRecord := &mysql.SkuInventoryRecord{
			ShopId:       row.ShopId,
			SkuCode:      row.SkuCode,
			OpType:       3, // 恢复库存
			OpUid:        0,
			OpIp:         "micro_mall_sku_cron",
			AmountBefore: skuInventory.Amount,
			Amount:       row.Amount,
			OpTxId:       row.OpTxId, // 恢复库存
			State:        0,
			Verify:       1,
			CreateTime:   time.Now(),
			UpdateTime:   time.Now(),
		}
		err = repository.CreateSkuInventoryRecord(tx, inventoryRestoreRecord)
		if err != nil {
			err = tx.Rollback()
			if err != nil {
				kelvins.ErrLogger.Errorf(ctx, "CreateSkuInventoryRecord  Rollback err: %v, ", err)
			}
			kelvins.ErrLogger.Errorf(ctx, "CreateSkuInventoryRecord err: %v, inventoryRestoreRecord: %+v", err, inventoryRestoreRecord)
			return
		}
		// 更新库存
		updateSkuInventoryWhere := map[string]interface{}{
			"shop_id":     skuInventory.ShopId,
			"sku_code":    skuInventory.SkuCode,
			"amount":      skuInventory.Amount,
			"update_time": skuInventory.UpdateTime,
		}
		updateSkuInventoryMaps := map[string]interface{}{
			"amount":      row.Amount,
			"update_time": time.Now(),
		}
		rowAffected, err = repository.UpdateSkuInventory(tx, updateSkuInventoryWhere, updateSkuInventoryMaps)
		if err != nil {
			err = tx.Rollback()
			if err != nil {
				kelvins.ErrLogger.Errorf(ctx, "UpdateSkuInventory  Rollback err: %v, ", err)
			}
			kelvins.ErrLogger.Errorf(ctx, "UpdateSkuInventory err: %v, where: %+v,maps: %+v ", err, updateSkuInventoryWhere, updateSkuInventoryMaps)
			return
		}
		if rowAffected != 1 {
			err = tx.Rollback()
			if err != nil {
				kelvins.ErrLogger.Errorf(ctx, "UpdateSkuInventory rowAffected Rollback err: %v, ", err)
			}
			return
		}
	}
	err = tx.Commit()
	if err != nil {
		kelvins.ErrLogger.Errorf(ctx, "HandleSkuInventoryRestore Commit err: %v, ", err)
	}
	return
}
