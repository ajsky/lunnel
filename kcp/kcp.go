package kcp

import (
	"net"

	"github.com/pkg/errors"
	kcp "github.com/xtaci/kcp-go"
)

var (
	noDelay        = 0
	interval       = 40
	resend         = 0
	noCongestion   = 0
	SockBuf        = 4194304
	dataShard      = 10
	parityShard    = 3
	udpSegmentSize = 1472
)

func Dial(addr string) (net.Conn, error) {
	block, _ := kcp.NewNoneBlockCrypt([]byte{12})
	kcpconn, err := kcp.DialWithOptions(addr, block, dataShard, parityShard)
	if err != nil {
		return nil, errors.Wrap(err, "create kcpConn")
	}
	kcpconn.SetStreamMode(true)
	kcpconn.SetNoDelay(noDelay, interval, resend, noCongestion)
	kcpconn.SetWindowSize(512, 512)
	kcpconn.SetMtu(udpSegmentSize)
	kcpconn.SetACKNoDelay(true)
	kcpconn.SetKeepAlive(10)

	if err := kcpconn.SetDSCP(0); err != nil {
		return nil, errors.Wrap(err, "kcpConn SetDSCP")
	}

	if err := kcpconn.SetReadBuffer(SockBuf); err != nil {
		return nil, errors.Wrap(err, "kcpConn SetReadBuffer")
	}
	if err := kcpconn.SetWriteBuffer(SockBuf); err != nil {
		return nil, errors.Wrap(err, "kcpConn SetWriteBuffer")
	}
	return kcpconn, nil
}

type Listener struct {
	lis *kcp.Listener
}

func Listen(addr string) (*Listener, error) {
	block, _ := kcp.NewNoneBlockCrypt([]byte{12})
	lis, err := kcp.ListenWithOptions(addr, block, dataShard, parityShard)
	if err != nil {
		return nil, errors.Wrap(err, "kcp ListenWithOptions")
	}

	if err := lis.SetDSCP(0); err != nil {
		return nil, errors.Wrap(err, "kcp SetDSCP")
	}
	if err := lis.SetReadBuffer(SockBuf); err != nil {
		return nil, errors.Wrap(err, "kcp SetReadBuffer")
	}
	if err := lis.SetWriteBuffer(SockBuf); err != nil {
		return nil, errors.Wrap(err, "kcp SetWriteBuffer")
	}
	return &Listener{lis: lis}, nil
}

func (l *Listener) Close() error {
	return l.lis.Close()
}
func (l *Listener) Addr() net.Addr {
	return l.lis.Addr()
}
func (l *Listener) Accept() (net.Conn, error) {
	conn, err := l.lis.AcceptKCP()
	if err != nil {
		return nil, errors.Wrap(err, "kcp AcceptKcp")
	}
	conn.SetStreamMode(true)
	conn.SetNoDelay(noDelay, interval, resend, noCongestion)
	conn.SetMtu(udpSegmentSize)
	conn.SetWindowSize(512, 512)
	conn.SetACKNoDelay(true)
	conn.SetKeepAlive(10)
	return conn, nil
}
