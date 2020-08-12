package model

// EtcdSdk 统一操作etcd接口
type EtcdSdk interface {
	List(path string) (list []*Node, err error) // 显示当前path下所有key
	Val(path string) (data *Node, err error)    // 获取path的值
	Add(path string, data []byte) (err error)   // 添加key
	Put(path string, data []byte) (err error)   // 修改key
	Del(path string) (err error)                // 删除key
	Members() (members []*Member, err error)    // 获取节点列表
	Close() (err error)                         // 关闭连接
}
