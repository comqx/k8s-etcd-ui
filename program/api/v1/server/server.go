package server

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/etcd-manage/etcd-manage-server/etcdsdk"
	"github.com/etcd-manage/etcd-manage-server/etcdsdk/model"
	"github.com/etcd-manage/etcd-manage-server/program/common"
	"github.com/etcd-manage/etcd-manage-server/program/logger"
	"github.com/etcd-manage/etcd-manage-server/program/models"
	"github.com/gin-gonic/gin"

	//"strings"
	"time"
)

// ServerController etcd服务列表相关操作
type ServerController struct {
}

// List 获取etcd服务列表，全部
func (api *ServerController) List(c *gin.Context) {
	name := c.Query("name")
	// 查询当前角色权限
	userinfoObj, exist := c.Get("userinfo")
	if exist == false {
		logger.Log.Warnw("用户登录信息不存在")
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	userinfo := userinfoObj.(*models.UsersModel)

	list, err := new(models.EtcdServersModel).All(name, userinfo.RoleId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": err.Error(),
		})
	}
	c.JSON(http.StatusOK, list)
}

// Add 添加服务
func (api *ServerController) Add(c *gin.Context) {
	var err error
	defer func() {
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"msg": err.Error(),
			})
		}
	}()
	// 添加
	req := new(models.EtcdServersModel)
	err = c.Bind(req)
	if err != nil {
		return
	}
	now := models.JSONTime(time.Now())
	req.CreatedAt = now
	err = req.Insert()
	if err != nil {
		return
	}
	// 添加超级管理员权限
	re := &models.RoleEtcdServersModel{
		EtcdServerId: req.ID,
		Type:         1,
		RoleId:       1,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	if err = re.Save(); err != nil {
		return
	}
	c.JSON(http.StatusOK, "ok")
}

// Update 修改服务
func (api *ServerController) Update(c *gin.Context) {
	var err error
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": err.Error(),
		})
	}
	// 添加
	req := new(models.EtcdServersModel)
	err = c.Bind(req)
	if err != nil {
		return
	}
	req.UpdatedAt = models.JSONTime(time.Now())
	err = req.Update()
	if err != nil {
		return
	}

	c.JSON(http.StatusOK, "ok")
}

// Restore 修复v1版本或e3w对目录的标记
// /v1/server/restore?etcd_id=6"
func (api *ServerController) Restore(c *gin.Context) {
	etcdId := c.Query("etcd_id")
	var err error
	defer func() {
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"msg": err.Error(),
			})
		}
	}()

	etcdIdNum, _ := strconv.Atoi(etcdId)
	etcdOne := new(models.EtcdServersModel)
	etcdOne, err = etcdOne.FirstById(int32(etcdIdNum))
	if err != nil {
		return
	}
	if etcdOne.Version != model.ETCD_VERSION_V3 {
		err = errors.New("Only V3 version is allowed to be repaired")
		return
	}
	// 连接etcd
	cfg := model.Config{
		Version:   etcdOne.Version,
		Address:   strings.Split(etcdOne.Address, ","),
		TlsEnable: etcdOne.TlsEnable,
		CertFile:  common.Base64Decode(etcdOne.CaFile),
		KeyFile:   common.Base64Decode(etcdOne.KeyFile),
		CaFile:    common.Base64Decode(etcdOne.CaFile),
		Username:  etcdOne.Username,
		Password:  etcdOne.Password,
	}
	//cfg := &model.Config{
	//	Version:   "v3",
	//	Address:   []string{"https://127.0.0.1:2379"},
	//	TlsEnable: "true",
	//	CertFile:  "/Users/liuqixiang/project/goStudyProject/etcd-ui/etcd-manage-server/0_cert.pem",
	//	KeyFile:   "/Users/liuqixiang/project/goStudyProject/etcd-ui/etcd-manage-server/0_key.pem",
	//	CaFile:    "/Users/liuqixiang/project/goStudyProject/etcd-ui/etcd-manage-server/0_ca.pem",
	//}

	if _, err = etcdsdk.NewClient(cfg); err != nil {
		fmt.Println("etcdv3.NewClient(cfg)>>>", err)
		return
	}
	//clientV3, ok := client.(*etcdsdk.EtcdV3)
	//if ok == false {
	//	err = errors.New("Connecting etcd V3 service error")
	//	return
	//}
	//err = client.Restore()
	//if err != nil {
	//	fmt.Println("clientV3.Restore()>>>", err)
	//	return
	//}
	c.JSON(http.StatusOK, "ok")
}

// SetRoles 设置etcd服务角色
func (api *ServerController) SetRoles(c *gin.Context) {
	req := make([]*models.AllByEtcdIdData, 0)
	err := c.Bind(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": err.Error(),
		})
		return
	}
	if len(req) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "未设置任何权限",
		})
		return
	}
	m := new(models.RoleEtcdServersModel)
	err = m.UpByEtcdId(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, "ok")
}

// GetRoles 获取etcd服务权限列表
func (api *ServerController) GetRoles(c *gin.Context) {
	etcdId := common.GetHttpToInt(c, "etcd_id")
	if etcdId <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "参数错误",
		})
		return
	}
	list, err := new(models.RoleEtcdServersModel).AllByEtcdId(int32(etcdId))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, list)
}

// Del 删除
func (s *ServerController) Del(c *gin.Context) {
	id := c.Query("id")
	idNum, _ := strconv.Atoi(id)
	if idNum == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "参数错误",
		})
		return
	}
	err := new(models.EtcdServersModel).Del(int32(idNum))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, nil)
}
