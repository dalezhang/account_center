## 环境变量
#### LOG_DIR
日志保存目录，默认在当前项目的tmp/log下
#### ENV
环境，默认development，关系到log的保存文件名
#### MYSQL_DNS
mysql的连接语句，为空的时候使用默认
```
root:root@tcp(mysql:33061)/arc-warden_test?parseTime=True&loc=Local
```
#### DEBUG
开启后打印http包下请求返回信息

#### PG_DNS
postgres的连接语句，为空的时候使用默认
```
host=pg port=54320 user=postgres password=postgres dbname=account_center_test sslmode=disable
```


## 生产swagger的json文件
```
swagger generate spec -o ./swagger.json
```

