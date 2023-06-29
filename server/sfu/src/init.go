package src

import (
	"log"
	"net/http"
	_ "net/http/pprof"
	"server/pkg/etcd"
	"server/pkg/proto"
	"server/pkg/util"
	"server/server/sfu/conf"
	"server/server/sfu/rtc"
	"strings"
	"time"

	nprotoo "github.com/zhuanxin-sz/nats-protoo"
)

const (
	statCycle = time.Second * 10
)

var (
	node   *etcd.ServiceNode
	nats   *nprotoo.NatsProtoo
	caster *nprotoo.Broadcaster
)

// Start 启动服务
func Start() {
	// 服务注册
	node = etcd.NewServiceNode(conf.Etcd.Addrs, conf.Global.Ndc, conf.Global.Nid, conf.Global.Name)
	node.RegisterNode()
	// 消息注册
	nats = nprotoo.NewNatsProtoo(conf.Nats.URL)
	nats.OnRequest(node.GetRPCChannel(), handleRpcMsg)
	// 消息广播
	caster = nats.NewBroadcaster(node.GetEventChannel())
	// 启动RTC
	rtc.InitRTC()
	// 启动调试
	if conf.Global.Pprof != "" {
		go debug()
	}
	// 启动其他
	go CheckRTC()
	go UpdatePayload()
}

// Stop 关闭连接
func Stop() {
	rtc.FreeRTC()
	if nats != nil {
		nats.Close()
	}
	if node != nil {
		node.Close()
	}
}

// CheckRTC 通知信令服务器流被移除
func CheckRTC() {
	for id := range rtc.CleanRouter {
		str := strings.Split(id, "/")
		rid := str[3]
		uid := str[5]
		mid := str[7]
		caster.Say(proto.SfuToBizOnStreamRemove, util.Map("rid", rid, "uid", uid, "mid", mid))
	}
}

// UpdatePayload 更新sfu服务器负载
func UpdatePayload() {
	t := time.NewTicker(statCycle)
	defer t.Stop()
	for range t.C {
		node.UpdateNodePayload(rtc.GetRouters())
	}
}

func debug() {
	log.Printf("Start sfu pprof on %s", conf.Global.Pprof)
	http.ListenAndServe(conf.Global.Pprof, nil)
}
