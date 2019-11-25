package getter

import "github.com/ruoklive/proxypool/pkg/models"

// IPGetter IP获取接口
type IPGetter interface {
	// Execute 执行获取IP操作
	Execute() []*models.IP
}