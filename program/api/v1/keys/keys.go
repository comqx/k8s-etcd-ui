package keys

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/etcd-manage/etcd-manage-server/etcdsdk/model"

	"github.com/etcd-manage/etcd-manage-server/etcdsdk"
	"github.com/etcd-manage/etcd-manage-server/program/publicFunc"

	"github.com/gin-gonic/gin"
)

// KeyController key控制器
type KeyController struct {
	publicFunc.PublicF
	client    etcdsdk.EtcdV3
	etcdIdNum int
}

// List 获取目录下key列表
// /v1/keys?path=
func (api *KeyController) List(c *gin.Context) {
	var (
		path string
		err  error
		list []*model.Node
	)
	defer func() {
		api.client.Close()
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"msg": err.Error(),
			})
		}
	}()

	path = c.DefaultQuery("path", "/")
	if api.etcdIdNum, err = strconv.Atoi(c.GetHeader("EtcdID")); err != nil {
		fmt.Printf("获取etcdid失败：%s\n", err.Error())
	}

	if api.client, err = api.GetEtcdClient(api.etcdIdNum); err != nil {
		fmt.Println("获取client失败，err", err.Error())
		return
	}

	if list, err = api.client.List(path); err != nil {
		fmt.Println("获取etcd 数据失败，err", err.Error())
		return
	}
	c.JSON(http.StatusOK, list)
}

// Val 获取一个key的值
func (api *KeyController) Val(c *gin.Context) {
	var (
		err  error
		list *model.Node
	)
	defer func() {
		api.client.Close()
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"msg": err.Error(),
			})
		}
	}()
	path := c.Query("path")
	if api.etcdIdNum, err = strconv.Atoi(c.GetHeader("EtcdID")); err != nil {
		fmt.Printf("获取etcdId 失败：%s\n", err.Error())
	}

	if api.client, err = api.GetEtcdClient(api.etcdIdNum); err != nil {
		fmt.Println("获取etcd 数据失败，err", err.Error())
		return
	}

	if list, err = api.client.Val(path); err != nil {
		fmt.Println("获取etcd 数据失败，err", err.Error())
		return
	}
	c.JSON(http.StatusOK, list)
}

// ReqKeyBody 添加和修改key请求body
type ReqKeyBody struct {
	Path  string `json:"path"`
	Value string `json:"value"`
}

// Add 添加key
func (api *KeyController) Add(c *gin.Context) {
	var (
		req ReqKeyBody
		err error
	)
	defer func() {
		api.client.Close()
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"msg": err.Error(),
			})
		}
	}()
	if err = c.Bind(req); err != nil {
		return
	}
	if api.etcdIdNum, err = strconv.Atoi(c.GetHeader("EtcdID")); err != nil {
		fmt.Printf("获取etcdId 失败：%s\n", err.Error())
	}

	if api.client, err = api.GetEtcdClient(api.etcdIdNum); err != nil {
		fmt.Println("获取etcd 数据失败，err", err.Error())
		return
	}

	if err = api.client.Add(req.Path, []byte(req.Value)); err != nil {
		return
	}
	c.JSON(http.StatusOK, "ok")
}

// Put 修改key
func (api *KeyController) Put(c *gin.Context) {
	var (
		req ReqKeyBody
		err error
	)
	defer func() {
		api.client.Close()
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"msg": err.Error(),
			})
		}
	}()

	if err = c.Bind(req); err != nil {
		return
	}
	if api.etcdIdNum, err = strconv.Atoi(c.GetHeader("EtcdID")); err != nil {
		fmt.Printf("获取etcdId 失败：%s\n", err.Error())
	}

	if api.client, err = api.GetEtcdClient(api.etcdIdNum); err != nil {
		fmt.Println("获取etcd 数据失败，err", err.Error())
		return
	}

	if err = api.client.Put(req.Path, []byte(req.Value)); err != nil {
		return
	}
	c.JSON(http.StatusOK, "ok")
}

// Del 删除key
func (api *KeyController) Del(c *gin.Context) {
	var (
		err  error
		path string
	)
	defer func() {
		api.client.Close()
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"msg": err.Error(),
			})
		}
	}()
	if path = c.Query("path"); path == "" {
		err = errors.New("Path cannot be empty")
		return
	}
	if api.etcdIdNum, err = strconv.Atoi(c.GetHeader("EtcdID")); err != nil {
		fmt.Printf("获取etcdId 失败：%s\n", err.Error())
		return
	}

	if api.client, err = api.GetEtcdClient(api.etcdIdNum); err != nil {
		fmt.Println("获取etcd 数据失败，err", err.Error())
		return
	}
	if err = api.client.Del(path); err != nil {
		return
	}
	c.JSON(http.StatusOK, "ok")
}

// Members 获取etcd服务节点
func (api *KeyController) Members(c *gin.Context) {
	var (
		err     error
		members []*model.Member
	)
	defer func() {
		api.client.Close()
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"msg": err.Error(),
			})
		}
	}()

	if api.etcdIdNum, err = strconv.Atoi(c.GetHeader("EtcdID")); err != nil {
		fmt.Printf("获取etcdId 失败：%s\n", err.Error())
		return
	}

	if api.client, err = api.GetEtcdClient(api.etcdIdNum); err != nil {
		fmt.Println("获取etcd 数据失败，err", err.Error())
		return
	}

	if members, err = api.client.Members(); err != nil {
		return
	}
	c.JSON(http.StatusOK, members)
}
