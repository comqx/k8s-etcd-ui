
# mysql相关命令

## mysql数据备份语句
```
 mysqldump -u root -p123456 -h127.0.0.1 -B  --add-drop-database --column-statistics=0  etcd_servers > etcd-manage.sql 
```
## mysql导入语句
```
mysql -uroot -p123456 -h127.0.0.1 < etcd-manage.sql
```

