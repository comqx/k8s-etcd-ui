package etcdsdk

import (
	"bytes"
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"log"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/etcd-manage/etcd-manage-server/etcdsdk/model"
	"go.etcd.io/etcd/clientv3"
	jsonserializer "k8s.io/apimachinery/pkg/runtime/serializer/json"
	"k8s.io/kubectl/pkg/scheme"
)

// EtcdV3Sdk etcd v3版
type EtcdV3 struct {
	cli *clientv3.Client
}

var (
	// DefaultTimeout 默认查询超时
	DefaultTimeout = 5 * time.Second
	sm             = new(sync.Mutex)
)

func NewClient(cfg model.Config) (client EtcdV3, err error) {
	var (
		client3Config clientv3.Config
		cert          tls.Certificate
	)
	fmt.Println("进入NewClient")
	client3Config = clientv3.Config{
		Endpoints:   cfg.Address,
		DialTimeout: 5 * time.Second,
	}
	// 判断是否使用证书
	if cfg.TlsEnable == "true" {
		if cert, err = tls.X509KeyPair(cfg.CertFile, cfg.KeyFile); err != nil {
			return
		}
		pool := x509.NewCertPool()
		pool.AppendCertsFromPEM(cfg.CaFile)
		_tlsConfig := &tls.Config{
			Certificates: []tls.Certificate{cert},
			RootCAs:      pool,
		}
		client3Config.TLS = _tlsConfig
	}

	cliNew, err := clientv3.New(client3Config)
	if err != nil {
		fmt.Println("出错：", err.Error())
	}
	//EtcdV3 实现了etcdsdk的全方法
	client = EtcdV3{cli: cliNew}
	return
}
func (sdk *EtcdV3) List(path string) (list []*model.Node, err error) {
	var (
		resp *clientv3.GetResponse
	)
	if path == "/" || path == "" {
		resp, err = clientv3.NewKV(sdk.cli).Get(context.Background(), "/", clientv3.WithFromKey(), clientv3.WithKeysOnly())
	} else {
		resp, err = clientv3.NewKV(sdk.cli).Get(context.Background(), path, clientv3.WithPrefix())
	}
	if err != nil {
		fmt.Println("list 失败：", err.Error())
		return
	}
	/* 处理出当前目录层的key */
	if resp.Count == 0 {
		fmt.Println("查询为空")
		return
	}
	/* 处理出当前目录层的key */
	if resp.Count == 0 {
		return
	}
	list, err = ConvertToPath(path, resp)

	// etcd 排序无效，自己实现
	sort.Slice(list, func(i, j int) bool {
		return list[i].Path < list[j].Path
	})

	// 如果是值，则查询值内容
	for _, v := range list {
		rv, err := sdk.cli.Get(context.Background(), v.Path)
		if err != nil {
			log.Println("读取值错误")
			continue
		}
		if len(rv.Kvs) > 0 {
			v.Value = string(rv.Kvs[0].Value)
		}
	}
	return
}
func (sdk *EtcdV3) Val(path string) (data *model.Node, err error) {
	var (
		resp  *clientv3.GetResponse
		value string
	)
	if resp, err = clientv3.NewKV(sdk.cli).Get(context.Background(), path); err != nil {
		fmt.Println("获取数据失败")
		return
	}

	decoder := scheme.Codecs.UniversalDeserializer()
	encoder := jsonserializer.NewSerializer(jsonserializer.DefaultMetaFactory, scheme.Scheme, scheme.Scheme, true)
	obj, _, err := decoder.Decode(resp.Kvs[0].Value, nil, nil)
	if err != nil {
		fmt.Printf("WARN: unable to decode %s: %v\n", resp.Kvs[0].Key, err)
		value = string(resp.Kvs[0].Value)
	} else {
		buffer := new(bytes.Buffer)
		err = encoder.Encode(obj, buffer)
		value = buffer.String()
		if err != nil {
			fmt.Printf("WARN: unable to encode %s: %v\n", resp.Kvs[0].Key, err)
			value = string(resp.Kvs[0].Value)
		}
	}

	data = &model.Node{
		Path:    string(resp.Kvs[0].Key),
		Name:    string(resp.Kvs[0].Key),
		Value:   value,
		Version: resp.Kvs[0].Version,
		Lease:   resp.Kvs[0].Lease,
	}
	return
}
func (sdk *EtcdV3) Add(path string, data []byte) (err error) {
	// 使用事物，防止覆盖，添加就是添加，不可以覆盖
	ctx, cancel := context.WithTimeout(context.Background(), DefaultTimeout)
	defer cancel()
	txn := sdk.cli.Txn(ctx)
	txn.If(
		clientv3.Compare(
			clientv3.Version(path),
			"=",
			0,
		),
	).Then(
		clientv3.OpPut(path, string(data)),
	)

	txnResp, err := txn.Commit()
	if err != nil {
		return err
	}

	if !txnResp.Succeeded {
		return model.ERR_ADD_KEY
	}
	return
}
func (sdk *EtcdV3) Put(path string, data []byte) (err error) {
	ctx, cancel := context.WithTimeout(context.Background(), DefaultTimeout)
	defer cancel()
	_, err = sdk.cli.Put(ctx, path, string(data))
	if err != nil {
		return
	}
	return
}
func (sdk *EtcdV3) Del(path string) (err error) {
	ctx, cancel := context.WithTimeout(context.Background(), DefaultTimeout)
	defer cancel()
	_, err = sdk.cli.Delete(ctx, path)
	if err != nil {
		return
	}
	return
}
func (sdk *EtcdV3) Members() (members []*model.Member, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), DefaultTimeout)
	defer cancel()
	resp, err := sdk.cli.MemberList(ctx)
	if err != nil {
		fmt.Println("获取members失败：", err.Error())
		return nil, err
	}
	fmt.Println(resp)
	for _, member := range resp.Members {
		fmt.Println(member.Name)
		if len(member.ClientURLs) > 0 {
			m := &model.Member{
				ID:         fmt.Sprint(member.ID),
				Name:       member.Name,
				PeerURLs:   member.PeerURLs,
				ClientURLs: member.ClientURLs,
				Role:       model.ROLE_FOLLOWER,
				Status:     model.STATUS_UNHEALTHY,
			}
			ctx, cancel := context.WithTimeout(context.Background(), DefaultTimeout)
			defer cancel()
			// log.Println(m.ClientURLs[0])
			resp, err := sdk.cli.Status(ctx, m.ClientURLs[0])
			if err == nil {
				m.Status = model.STATUS_HEALTHY
				m.DbSize = resp.DbSize
				if resp.Leader == resp.Header.MemberId {
					m.Role = model.ROLE_LEADER
				}
			}
			members = append(members, m)
		}
	}
	return
}
func (sdk *EtcdV3) Close() error {
	return sdk.cli.Close()
}

//
//func (sdk *EtcdV3Sdk) getKey(client *clientv3.Client, key string) error {
//	resp, err := clientv3.NewKV(client).Get(context.Background(), key)
//	if err != nil {
//		return err
//	}
//
//	decoder := scheme.Codecs.UniversalDeserializer()
//	encoder := jsonserializer.NewSerializer(jsonserializer.DefaultMetaFactory, scheme.Scheme, scheme.Scheme, true)
//
//	for _, kv := range resp.Kvs {
//		obj, gvk, err := decoder.Decode(kv.Value, nil, nil)
//		if err != nil {
//			fmt.Fprintf(os.Stderr, "WARN: unable to decode %s: %v\n", kv.Key, err)
//			continue
//		}
//		fmt.Println(gvk)
//		err = encoder.Encode(obj, os.Stdout)
//		if err != nil {
//			fmt.Fprintf(os.Stderr, "WARN: unable to encode %s: %v\n", kv.Key, err)
//			continue
//		}
//	}
//
//	return nil
//}
//
//func (sdk *EtcdV3Sdk) getKeyAll(client *clientv3.Client, key string) error {
//	resp, err := clientv3.NewKV(client).Get(context.Background(), key, clientv3.WithPrefix())
//	if err != nil {
//		return err
//	}
//
//	decoder := scheme.Codecs.UniversalDeserializer()
//	encoder := jsonserializer.NewSerializer(jsonserializer.DefaultMetaFactory, scheme.Scheme, scheme.Scheme, true)
//
//	for _, kv := range resp.Kvs {
//		obj, gvk, err := decoder.Decode(kv.Value, nil, nil)
//		if err != nil {
//			fmt.Fprintf(os.Stderr, "WARN: unable to decode %s: %v\n", kv.Key, err)
//			continue
//		}
//		fmt.Println(gvk)
//		err = encoder.Encode(obj, os.Stdout)
//		if err != nil {
//			fmt.Fprintf(os.Stderr, "WARN: unable to encode %s: %v\n", kv.Key, err)
//			continue
//		}
//		//fmt.Println(string(kv.Key))
//	}
//
//	return nil
//}
//
//func (sdk *EtcdV3Sdk) dump(client *clientv3.Client) error {
//	response, err := clientv3.NewKV(client).Get(context.Background(), "/", clientv3.WithPrefix(), clientv3.WithSort(clientv3.SortByKey, clientv3.SortDescend))
//	if err != nil {
//		return err
//	}
//
//	kvData := []etcd3kv{}
//	decoder := scheme.Codecs.UniversalDeserializer()
//	encoder := jsonserializer.NewSerializer(jsonserializer.DefaultMetaFactory, scheme.Scheme, scheme.Scheme, false)
//	objJSON := &bytes.Buffer{}
//
//	for _, kv := range response.Kvs {
//		obj, _, err := decoder.Decode(kv.Value, nil, nil)
//		if err != nil {
//			fmt.Fprintf(os.Stderr, "WARN: error decoding value %q: %v\n", string(kv.Value), err)
//			continue
//		}
//		objJSON.Reset()
//		if err := encoder.Encode(obj, objJSON); err != nil {
//			fmt.Fprintf(os.Stderr, "WARN: error encoding object %#v as JSON: %v", obj, err)
//			continue
//		}
//		kvData = append(
//			kvData,
//			etcd3kv{
//				Key:            string(kv.Key),
//				Value:          string(objJSON.Bytes()),
//				CreateRevision: kv.CreateRevision,
//				ModRevision:    kv.ModRevision,
//				Version:        kv.Version,
//				Lease:          kv.Lease,
//			},
//		)
//	}
//
//	jsonData, err := json.MarshalIndent(kvData, "", "  ")
//	if err != nil {
//		return err
//	}
//
//	fmt.Println(string(jsonData))
//
//	return nil
//}

// ConvertToPath 处理etcd3 的key为目录形式 - path只能是/结尾或为空
func ConvertToPath(path string, resp *clientv3.GetResponse) (list []*model.Node, err error) {

	keyMapVal := make(map[string]*model.Node, 0)
	keyMapPath := make(map[string]*model.Node, 0)

	for _, val := range resp.Kvs {
		if ok := strings.HasPrefix(string(val.Key), path); ok == true {
			key := string(val.Key)[len(path):]
			// 判断是否是//开头，如果是，则本级目录为//
			if strings.HasPrefix(key, "//") == true {
				fullKey := path + "/"
				keyMapPath["/"] = &model.Node{
					IsDir:   true,
					Path:    fullKey,
					Name:    "/",
					Value:   string(val.Value),
					Version: 0,
				}
				continue
			}
			// 处理path为/情况
			if path == "" && strings.HasPrefix(key, "/") {
				keyMapPath["/"] = &model.Node{
					IsDir:   true,
					Path:    "/",
					Name:    "/",
					Value:   string(val.Value),
					Version: 0,
				}
				continue
			}
			// 查找下一个/位置
			i := strings.Index(key, "/")
			// 截取后第一个字符是/的情况
			if i == 0 {
				key = key[1:]
				i = strings.Index(key, "/")
			}
			// 未查找到，则为key，而不是目录
			if i == -1 { // 未查询到则证明是key，而不是目录 则取完整路径最后一个/之后的部分作为name
				lastIndex := strings.LastIndex(string(val.Key), "/")
				key = string(val.Key)[lastIndex+1:]
				keyMapVal[key] = &model.Node{
					IsDir:   false,
					Path:    string(val.Key),
					Name:    key,
					Value:   string(val.Value),
					Version: val.Version,
				}
			} else {
				key = key[:i]
				// 等于当前path，不返回
				if key == "" {
					continue
				}
				fullKey := path + "/" + key
				if path == "/" {
					fullKey = path + key
				} else if path == "" {
					fullKey = key
				}
				lastIndex := strings.LastIndex(fullKey, "/")
				key = fullKey[lastIndex+1:]
				// log.Println(path, " -- ", key)
				keyMapPath[key] = &model.Node{
					IsDir:   true,
					Path:    fullKey,
					Name:    key,
					Value:   string(val.Value),
					Version: 0,
				}
			}
		}
	}
	for _, val := range keyMapPath {
		list = append(list, val)
	}
	for _, val := range keyMapVal {
		list = append(list, val)
	}

	return
}
