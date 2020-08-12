# k8s Etcd manage Server


## 功能介绍
etcd-manage-server 是一个用go编写的etcd可视化管理工具，具有友好的界面，管理key就像管理本地文件一样方便。支持简单权限管理区分只读和读写权限。

**备注**

1. 本项目适配k8s-etcd的乱码问题
2. 将sql文件导入到mysql数据库(sql文件路径 ./data/etcd-manage.sql)，默认用户 admin/111111 
3. 当前只实现了etcd v3 api管理
4. 在使用时可直接修改默认的两个etcd连接地址为真实可用地址即可开始体验。

## 使用
```
# 1. sql数据导入到mysql数据库

# 2. 打包、运行
go mod tidy
go mod download
go build .

./etcd-manage-server
```
## 效果演示

etcd管理页面
![](https://github.com/qixiang-liu/k8s-etcd-ui/blob/master/tupian/etcd%E7%AE%A1%E7%90%86%E9%A1%B5%E9%9D%A2.png)

多个etcd集群健康图
![](https://github.com/qixiang-liu/k8s-etcd-ui/blob/master/tupian/%E5%A4%9A%E4%B8%AAetcd%E9%9B%86%E7%BE%A4%E5%81%A5%E5%BA%B7%E5%9B%BE.png)

key树以及value展示图
![](https://github.com/qixiang-liu/k8s-etcd-ui/blob/master/tupian/etcd%20%20key%E6%A0%91%E4%BB%A5%E5%8F%8Avalue%E5%B1%95%E7%A4%BA%E5%9B%BE.png)

添加etcd界面
![](https://github.com/qixiang-liu/k8s-etcd-ui/blob/master/tupian/%E6%B7%BB%E5%8A%A0etcd%E7%95%8C%E9%9D%A2.png)

etcd-key树图
![](https://github.com/qixiang-liu/k8s-etcd-ui/blob/master/tupian/etcd-key%E6%A0%91%E5%9B%BE.png)


