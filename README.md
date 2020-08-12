# k8s Etcd manage Server


## 功能介绍

etcd-manage-server 是一个用go编写的etcd可视化管理工具，具有友好的界面，管理key就像管理本地文件一样方便。支持简单权限管理区分只读和读写权限。

**备注**

1. 本项目适配k8s-etcd的乱码问题
2. 将sql文件导入到mysql数据库(sql文件路径 ./data/etcd-manage.sql)，默认用户 admin/111111 
5. 当前只实现了etcd v3 api管理
6. 在使用时可直接修改默认的两个etcd连接地址为真实可用地址即可开始体验。


## 效果演示

etcd服务列表管理

key 管理

key 编辑

key 查看

用户管理

