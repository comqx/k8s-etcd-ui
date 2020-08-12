package publicFunc

import (
	"strings"

	"github.com/etcd-manage/etcd-manage-server/program/common"

	"github.com/etcd-manage/etcd-manage-server/etcdsdk"
	"github.com/etcd-manage/etcd-manage-server/etcdsdk/model"
	"github.com/etcd-manage/etcd-manage-server/program/logger"
	"github.com/etcd-manage/etcd-manage-server/program/models"
)

type PublicF struct {
}

func (p PublicF) GetEtcdClient(etcdId int) (client etcdsdk.EtcdV3, err error) {
	// 查询etcd服务信息
	etcdOne := new(models.EtcdServersModel)
	etcdOne, err = etcdOne.FirstById(int32(etcdId))
	if err != nil {
		logger.Log.Errorw("获取etcd服务信息错误", "EtcdID", etcdId, "err", err)
	}
	if etcdOne.TlsEnable == "true" {
		if !strings.HasPrefix(etcdOne.Address, "https") {
			etcdOne.Address = strings.ReplaceAll(etcdOne.Address, "http", "https")
		}
	} else {
		if strings.HasPrefix(etcdOne.Address, "https") {
			etcdOne.Address = strings.ReplaceAll(etcdOne.Address, "https", "http")
		}
	}
	// 连接etcd
	cfg := model.Config{
		EtcdId:    int32(etcdId),
		Version:   etcdOne.Version,
		Address:   strings.Split(etcdOne.Address, ","),
		TlsEnable: etcdOne.TlsEnable,
		CertFile:  common.Base64Decode(etcdOne.CertFile),
		KeyFile:   common.Base64Decode(etcdOne.KeyFile),
		CaFile:    common.Base64Decode(etcdOne.CaFile),
	}
	if client, err = etcdsdk.NewClient(cfg); err != nil {
		logger.Log.Errorf("初始化etcd client 失败：%s", err.Error())
	}
	return
}
