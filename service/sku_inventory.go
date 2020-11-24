package service

import (
	"context"
	"gitee.com/cristiane/micro-mall-sku-cron/model/args"
	"gitee.com/cristiane/micro-mall-sku-cron/model/mysql"
	"gitee.com/cristiane/micro-mall-sku-cron/pkg/util"
	"gitee.com/cristiane/micro-mall-sku-cron/proto/micro_mall_order_proto/order_business"
	"gitee.com/cristiane/micro-mall-sku-cron/repository"
	"gitee.com/kelvins-io/kelvins"
	"github.com/google/uuid"
	"time"
)

const (
	sqlSelectSkuInventoryRestore = "op_tx_id"
	sqlSelectSkuInventoryRecord  = "id,shop_id,sku_code,amount,update_time,op_tx_id"
)

func HandleOrderFailedSkuInventoryRestore() {
	ctx := context.Background()
	where := map[string]interface{}{
		"verify":  0, // 记录未经验证
		"op_type": 1, // 出库
	}
	recordList, err := repository.FindSkuInventoryRecordList(sqlSelectSkuInventoryRestore, where, 300, 1)
	if err != nil {
		kelvins.ErrLogger.Errorf(ctx, "FindSkuInventoryRecordList,err: %v, req: %+v", err, where)
		return
	}
	if len(recordList) == 0 {
		return
	}
	opTxIds := make([]string, 0)
	opTxIdsSet := map[string]struct{}{}
	for i := 0; i < len(recordList); i++ {
		if _, ok := opTxIdsSet[recordList[i].OpTxId]; !ok {
			opTxIdsSet[recordList[i].OpTxId] = struct{}{}
			opTxIds = append(opTxIds, recordList[i].OpTxId)
		}
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
	skuInventoryFailedOrderSet := map[string]struct{}{}
	for i := 0; i < len(rsp.List); i++ {
		if !rsp.List[i].IsExist {
			if _, ok := skuInventoryFailedOrderSet[rsp.List[i].OrderCode]; !ok {
				skuInventoryFailedOrderSet[rsp.List[i].OrderCode] = struct{}{}
				skuInventoryFailedOrder = append(skuInventoryFailedOrder, rsp.List[i].OrderCode)
			}
		}
	}
	if len(skuInventoryFailedOrder) == 0 {
		return
	}
	skuInventoryRecordWhere := map[string]interface{}{
		"op_tx_id": skuInventoryFailedOrder,
		"verify":   0, // 记录未经验证
		"op_type":  1, // 出库
	}
	skuInventoryRecordList, err := repository.FindSkuInventoryRecordList(sqlSelectSkuInventoryRecord, skuInventoryRecordWhere, 300, 1)
	if err != nil {
		kelvins.ErrLogger.Errorf(ctx, "FindSkuInventoryRecordList err: %v, where: %+v", err, skuInventoryRecordWhere)
		return
	}
	if len(skuInventoryRecordList) == 0 {
		return
	}
	for i := 0; i < len(skuInventoryRecordList); i++ {
		tx := kelvins.XORM_DBEngine.NewSession()
		err = tx.Begin()
		if err != nil {
			kelvins.ErrLogger.Errorf(ctx, "HandleSkuInventoryRestore Begin err: %v, ", err)
			return
		}
		row := skuInventoryRecordList[i]
		getSkuInventoryWhere := map[string]interface{}{
			"shop_id":  row.ShopId,
			"sku_code": row.SkuCode,
		}
		skuInventory, err := repository.GetSkuInventory(tx, "shop_id,sku_code,amount,last_tx_id", getSkuInventoryWhere)
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
		opTxId := uuid.New().String()
		// 更新库存记录
		inventoryRecordWhere := map[string]interface{}{
			"id":       row.Id,
			"op_tx_id": row.OpTxId,
			"op_type":  1, // 出库
			"verify":   0, // 未验证
		}
		inventoryRecordMaps := map[string]interface{}{
			"verify":      1,
			"op_tx_id":    opTxId,
			"update_time": time.Now(),
		}
		rowAffected, err := repository.UpdateSkuInventoryRecordByTx(tx, inventoryRecordWhere, inventoryRecordMaps)
		if err != nil {
			err = tx.Rollback()
			if err != nil {
				kelvins.ErrLogger.Errorf(ctx, "HandleSkuInventoryRestore  Rollback err: %v, ", err)
			}
			kelvins.ErrLogger.Errorf(ctx, "UpdateSkuInventoryRecord err: %v, where: %+v,maps: %+v ", err, inventoryRecordWhere, inventoryRecordMaps)
			return
		}
		// 库存记录可能是同一个订单扣减的
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
			OpTxId:       opTxId, // 恢复库存
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
			"shop_id":    skuInventory.ShopId,
			"sku_code":   skuInventory.SkuCode,
			"amount":     skuInventory.Amount,
			"last_tx_id": skuInventory.LastTxId,
		}
		updateSkuInventoryMaps := map[string]interface{}{
			"amount":      skuInventory.Amount + row.Amount,
			"last_tx_id":  opTxId,
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
		err = tx.Commit()
		if err != nil {
			kelvins.ErrLogger.Errorf(ctx, "HandleSkuInventoryRestore Commit err: %v, ", err)
		}
	}

	return
}

func HandleOrderSuccessSkuInventoryRestore() {
	ctx := context.Background()
	where := map[string]interface{}{
		"verify":  0, // 记录未经验证
		"op_type": 1, // 出库
	}
	recordList, err := repository.FindSkuInventoryRecordList(sqlSelectSkuInventoryRestore, where, 300, 1)
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
	// 支付成功的订单
	skuInventorySuccessOrder := make([]string, 0)
	skuInventorySuccessOrderSet := map[string]struct{}{}
	for i := 0; i < len(rsp.List); i++ {
		if rsp.List[i].IsExist && rsp.List[i].PayState == order_business.OrderPayStateType_PAY_SUCCESS {
			if _, ok := skuInventorySuccessOrderSet[rsp.List[i].OrderCode]; !ok {
				skuInventorySuccessOrderSet[rsp.List[i].OrderCode] = struct{}{}
				skuInventorySuccessOrder = append(skuInventorySuccessOrder, rsp.List[i].OrderCode)
			}
		}
	}
	if len(skuInventorySuccessOrder) == 0 {
		return
	}
	updateWhere := map[string]interface{}{
		"op_tx_id": skuInventorySuccessOrder,
		"verify":   0, // 记录未经验证
		"op_type":  1, // 出库
	}
	updateMaps := map[string]interface{}{
		"verify":      1,
		"op_tx_id":    uuid.New().String(),
		"update_time": time.Now(),
	}
	rowAffected, err := repository.UpdateSkuInventoryRecord(updateWhere, updateMaps)
	if err != nil {
		kelvins.ErrLogger.Errorf(ctx, "UpdateSkuInventoryRecord err: %v, where: %+v", err, updateWhere)
		return
	}
	_ = rowAffected
}
