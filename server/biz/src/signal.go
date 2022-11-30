package src

import (
	"net/http"
	"server/pkg/util"
	"server/server/biz/ws"

	"github.com/zhuanxin-sz/go-protoo/logger"
	"github.com/zhuanxin-sz/go-protoo/transport"
)

// InitSignalServer 初始化ws
func InitSignalServer(host string, port int, cert, key string) {
	config := ws.DefaultConfig()
	config.Host = host
	config.Port = port
	config.CertFile = cert
	config.KeyFile = key
	wsServer := ws.NewWebSocketServer(handler)
	go wsServer.Bind(config)
}

func handler(transport *transport.WebSocketTransport, request *http.Request) {
	logger.Debugf("handler = %v", request.URL.Query())
	vars := request.URL.Query()
	peerID := vars["peer"]
	if peerID == nil || len(peerID) < 1 {
		return
	}
	id := peerID[0]
	peer := ws.NewPeer(id, transport)

	handleRequest := func(request map[string]interface{}, accept ws.AcceptFunc, reject ws.RejectFunc) {
		defer util.Recover("biz handleRequest")

		method := util.Val(request, "method")
		if method == "" {
			reject(-1, ws.ErrInvalidMethod)
			return
		}

		data := request["data"]
		if data == nil {
			reject(-1, ws.ErrInvalidData)
			return
		}

		msg := data.(map[string]interface{})
		handlerWebsocket(method, peer, msg, accept, reject)
	}

	handleNotification := func(notification map[string]interface{}) {
		defer util.Recover("biz handleNotification")

		method := util.Val(notification, "method")
		if method == "" {
			ws.DefaultReject(-1, ws.ErrInvalidMethod)
			return
		}

		data := notification["data"]
		if data == nil {
			ws.DefaultReject(-1, ws.ErrInvalidData)
			return
		}

		msg := data.(map[string]interface{})
		handlerWebsocket(method, peer, msg, ws.DefaultAccept, ws.DefaultReject)
	}

	handleClose := func(code int, err string) {
		peer.Close()
	}

	peer.On("request", handleRequest)
	peer.On("notification", handleNotification)
	peer.On("close", handleClose)
	peer.On("error", handleClose)
}
