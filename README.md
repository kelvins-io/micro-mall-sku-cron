# micro-mall-sku-cron

#### 介绍
商品库cron

#### 软件架构
cron

#### 框架，库依赖
kelvins框架支持（gRPC，cron，queue，web支持）：https://gitee.com/kelvins-io/kelvins   
g2cache缓存库支持（两级缓存）：https://gitee.com/kelvins-io/g2cache   

#### 安装教程

1.仅构建  sh build.sh   
2 运行  sh build-run.sh   
3 停止 sh stop.sh

#### 使用说明
配置参考
```toml
[kelvins-server]
Environment = "dev"

[kelvins-logger]
RootPath = "./logs"
Level = "debug"

[kelvins-mysql]
Host = "127.0.0.1:3306"
UserName = "root"
Password = "fasdfasdf"
DBName = "micro_mall_sku"
Charset = "utf8mb4"
PoolNum =  10
MaxIdleConns = 5
ConnMaxLifeSecond = 3600
MultiStatements = true
ParseTime = true

[email-config]
Enable = false
User = "fasdfa@qq.com"
Password = "fasdfa"
Host = "smtp.qq.com"
Port = "465"

[order-failed-inventory-restore-task]
Cron = "0 */4 * * * *"
```

#### 参与贡献

1.  Fork 本仓库
2.  新建 Feat_xxx 分支
3.  提交代码
4.  新建 Pull Request

