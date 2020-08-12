package model

import "errors"

// etcd 错误定义
var (
	ERR_CONFIG_ISNIL        = errors.New("Config is nil")
	ERR_TLS_CONFIG_ISNIL    = errors.New("TLSConfig is nil")
	ERR_ETCD_ADDRESS_EMPTY  = errors.New("Etcd connection address cannot be empty")
	ERR_UNSUPPORTED_VERSION = errors.New("Unsupported etcd version")

	ERR_ADD_KEY       = errors.New("Add key error")
	ERR_KEY_NOT_FOUND = errors.New("Key does not exist")
	ERR_KEY_NOT_DIR   = errors.New("Key is not a directory")
)
