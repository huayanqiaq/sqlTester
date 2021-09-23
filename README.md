# 使用说明

### 这是一个无需各种驱动来连接数据库测试的工具


```
Usage of ./xorm_test_darwin_amd64:
  -dbname string
    	Input Your dbname (default "mysql,orcl")
  -flagname int
    	Just for demo (default 1234)
  -host string
    	Input Your host (default "127.0.0.1")
  -password string
    	Input Your password (default "root")
  -port string
    	Input Your port (default "3306")
  -sql string
    	Input Your sql (default "select version()")
  -type string
    	Input Your databases type (default "mysql,mssql,oracle")
  -user string
    	Input Your user (default "root")
```



1.对于oracle

```
xorm_test_darwin_amd64 -host 192.168.1.4 -user system -type oracle -port 1521 -password "system" -sql "SELECT * FROM SCOTT.EMP" -dbname helowin
```

