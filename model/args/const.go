package args

type MerchantsMaterialInfo struct {
	Uid          int64
	MaterialId   int64
	RegisterAddr string
	HealthCardNo string
	Identity     int32
	State        int32
	TaxCardNo    string
}

type SkuInventoryInfo struct {
	ShopId        int64  `json:"shop_id"`
	SkuCode       string `json:"sku_code"`
	Name          string `json:"name"`
	Price         string `json:"price"`
	Title         string `json:"title"`
	SubTitle      string `json:"sub_title"`
	Desc          string `json:"desc"`
	Production    string `json:"production"`
	Supplier      string `json:"supplier"`
	Category      int32  `json:"category"`
	Color         string `json:"color"`
	ColorCode     int32  `json:"color_code"`
	Specification string `json:"specification"`
	DescLink      string `json:"desc_link"`
}

const (
	RpcServiceMicroMallUsers = "micro-mall-users"
	RpcServiceMicroMallShop  = "micro-mall-shop"
	RpcServiceMicroMallOrder = "micro-mall-order"
)

const (
	SkuInventorySearchNotice       = 1000
	SkuInventorySearchNoticeTag    = "sku_inventory_search_notice"
	SkuInventorySearchNoticeTagErr = "sku_inventory_search_notice_err"
)

type CommonBusinessMsg struct {
	Type    int    `json:"type"`
	Tag     string `json:"tag"`
	UUID    string `json:"uuid"`
	Content string `json:"content"`
}

type TradeOrderDetail struct {
	ShopId    int64  `json:"shop_id"`
	OrderCode string `json:"order_code"`
}

type TradeOrderNotice struct {
	Uid  int64  `json:"uid"`
	Time string `json:"time"`
	// 9-19修改为 直接通知交易号, 放弃通知[]TradeOrderDetail
	TxCode string `json:"tx_code"`
}

const (
	Unknown                   = 0
	TradeOrderEventTypeCreate = 10014
	TradeOrderEventTypeExpire = 10015
)

var MsgFlags = map[int]string{
	Unknown:                   "未知",
	TradeOrderEventTypeCreate: "交易订单创建",
	TradeOrderEventTypeExpire: "交易订单过期",
}

func GetMsg(code int) string {
	msg, ok := MsgFlags[code]
	if ok {
		return msg
	}
	return MsgFlags[Unknown]
}
