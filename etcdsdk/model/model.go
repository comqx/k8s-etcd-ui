package model

// 这里保存sdk使用到的所有结构图

type ListNodes struct {
	IsDir   bool   `json:"is_dir"`
	Path    string `json:"path"`
	Name    string `json:"name"`
	Version int64  `json:"version"`
	Lease   int64  `json:"lease"`
}

// Node 一个key 目录或文件
type Node struct {
	IsDir   bool   `json:"is_dir"`
	Path    string `json:"path"`
	Name    string `json:"name"`
	Value   string `json:"value"`
	Version int64  `json:"version"`
	Lease   int64  `json:"lease"`
}

const (
	ROLE_LEADER   = "leader"
	ROLE_FOLLOWER = "follower"

	STATUS_HEALTHY   = "healthy"
	STATUS_UNHEALTHY = "unhealthy"
)

// Member 节点信息
type Member struct {
	// *etcdserverpb.Member
	ID         string   `json:"ID"`
	Name       string   `json:"name"`
	PeerURLs   []string `json:"peerURLs"`
	ClientURLs []string `json:"clientURLs"`
	Role       string   `json:"role"`
	Status     string   `json:"status"`
	DbSize     int64    `json:"db_size"`
}
