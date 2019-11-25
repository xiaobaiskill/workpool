package all

import (
	//_ "github.com/ruoklive/proxypool/internal/ipgetter/data5u"
	_ "github.com/ruoklive/proxypool/internal/ipgetter/feiyi"
	//_ "github.com/ruoklive/proxypool/internal/ipgetter/goubanjia" // 因为网站限制，无法正常下载数据
	//_ "github.com/ruoklive/proxypool/internal/ipgetter/ip181" // 已经无法使用
	_ "github.com/ruoklive/proxypool/internal/ipgetter/ip66"
	_ "github.com/ruoklive/proxypool/internal/ipgetter/ip89"
	_ "github.com/ruoklive/proxypool/internal/ipgetter/kuaidl"
	//_ "github.com/ruoklive/proxypool/internal/ipgetter/plp" // 网址无法访问
	//_ "github.com/ruoklive/proxypool/internal/ipgetter/xdaili" // 网址无法访问
	//_ "github.com/ruoklive/proxypool/internal/ipgetter/xicidl" // 网址无法访问
	//_ "github.com/ruoklive/proxypool/internal/ipgetter/youdl" //失效的采集脚本，用作系统容错实验
)
