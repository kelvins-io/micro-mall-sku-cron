package vars

type EmailConfigSettingS struct {
	Enable   bool   `json:"enable"`
	User     string `json:"user"`
	Password string `json:"password"`
	Host     string `json:"host"`
	Port     string `json:"port"`
}

type OrderFailedInventoryRestoreTaskSettingS struct {
	Cron string `json:"cron"`
}

type SkuInventorySearchSyncTaskSettingS struct {
	Cron string `json:"cron"`
}
