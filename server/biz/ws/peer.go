package ws

import (
	"github.com/zhuanxin-sz/go-protoo/logger"
	"github.com/zhuanxin-sz/go-protoo/peer"
	"github.com/zhuanxin-sz/go-protoo/transport"
)

// Peer peer对象
type Peer struct {
	peer.Peer
}

// NewPeer 新建Peer对象
func NewPeer(uid string, t *transport.WebSocketTransport) *Peer {
	return &Peer{
		Peer: *peer.NewPeer(uid, t),
	}
}

// On 事件处理
func (peer *Peer) On(event, listener interface{}) {
	peer.Peer.On(event, listener)
}

// Close Peer关闭
func (peer *Peer) Close() {
	logger.Debugf("Close Room Peer=%s", peer.ID())
	peer.Peer.Close()
}
