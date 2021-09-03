echo 当前分支
git branch

echo 拉取依赖
go mod vendor

echo 开始构建
go build -o micro-mall-sku-cron main.go

cp -n ./etc/app.ini.example ./etc/app.ini

echo 开始运行micro-mall-sku-cron
nohup ./micro-mall-sku-cron -s start >nohup.out  2>&1  &