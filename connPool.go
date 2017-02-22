// Copyright 2017 Lynn Deng (ldeng7)

package golua

import (
	"bufio"
	"errors"
	"fmt"
	"net"
	"sync"
	"time"
)

type tcpConn struct {
	*net.TCPConn
	bufReader *bufio.Reader
	poolKey   string
	reuseCnt  int
}

type connPool struct {
	sync.Pool
}

func (ctx *vmCtx) getConnFromPool(host string, port int, timeout int, key string) (*tcpConn, error) {
	ctx.connPoolsLock.Lock()
	pool, ok := ctx.connPools[key]
	if !ok {
		pool = &connPool{}
		pool.New = func() interface{} {
			conn, _ := net.DialTimeout("tcp", fmt.Sprintf("%s:%d", host, port),
				time.Duration(timeout)*time.Millisecond)
			goTcpConn, _ := conn.(*net.TCPConn)
			return &tcpConn{goTcpConn, bufio.NewReader(goTcpConn), key, 0}
		}
		ctx.connPools[key] = pool
	}
	ctx.connPoolsLock.Unlock()

	o := pool.Get()
	if nil == o {
		return nil, errors.New("failed to get connection from pool")
	}
	tc, _ := o.(*tcpConn)
	ctx.goObjRefs.connsLock.Lock()
	ctx.goObjRefs.conns[tc] = true
	ctx.goObjRefs.connsLock.Unlock()
	return tc, nil
}

func (ctx *vmCtx) returnConnToPool(conn *tcpConn) error {
	pool, ok := ctx.connPools[conn.poolKey]
	if !ok {
		return errors.New("invalid pool key")
	}
	conn.SetKeepAlive(true)
	conn.reuseCnt = conn.reuseCnt + 1
	pool.Put(conn)
	return nil
}

func (ctx *vmCtx) closeConn(conn *tcpConn) error {
	ctx.goObjRefs.connsLock.Lock()
	delete(ctx.goObjRefs.conns, conn)
	ctx.goObjRefs.connsLock.Unlock()
	return conn.Close()
}
