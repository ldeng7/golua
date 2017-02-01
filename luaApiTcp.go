// Copyright 2017 Lynn Deng (ldeng7)

package golua

/*
#include <goluaCommon.h>
*/
import "C"
import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"time"
	"unsafe"
)

var (
	cStrLuaTableKey_tcp_conn     = C.CString("_conn")
	cStrLuaTableKey_tcp_ctimeout = C.CString("_ctimeout")
	cStrLuaTableKey_tcp_wtimeout = C.CString("_wtimeout")
	cStrLuaTableKey_tcp_rtimeout = C.CString("_rtimeout")
)

//export GoluaLuaApi_tcp
func GoluaLuaApi_tcp(l luaState) C.int {
	// conn = go.tcp()

	C.lua_createtable(l, 0, 4)
	C.lua_pushinteger(l, 10000)
	C.lua_setfield(l, 1, cStrLuaTableKey_tcp_ctimeout)
	C.lua_pushinteger(l, 10000)
	C.lua_setfield(l, 1, cStrLuaTableKey_tcp_wtimeout)
	C.lua_pushinteger(l, 10000)
	C.lua_setfield(l, 1, cStrLuaTableKey_tcp_rtimeout)
	return 1
}

//export GoluaLuaApi_tcp_close
func GoluaLuaApi_tcp_close(l luaState) C.int {
	// ok, err = go.tcp_close(conn)

	C.lua_getfield(l, C.LUA_REGISTRYINDEX, cStrLuaRegKey_vmCtx)
	ctx := (*vmCtx)(C.lua_touserdata(l, -1))
	C.lua_getfield(l, 1, cStrLuaTableKey_tcp_conn)
	conn := (*tcpConn)(C.lua_touserdata(l, -1))

	var ok C.int = 1
	err := ctx.closeConn(conn)
	if nil != err {
		ok = 0
	}
	C.lua_pushboolean(l, ok)
	luaApiPushErr(l, err)
	return 2
}

//export GoluaLuaApi_tcp_connect
func GoluaLuaApi_tcp_connect(l luaState) C.int {
	// ok, err = go.tcp_connect(conn, host, port, pool_key)

	C.lua_getfield(l, C.LUA_REGISTRYINDEX, cStrLuaRegKey_vmCtx)
	ctx := (*vmCtx)(C.lua_touserdata(l, -1))
	host := C.GoString(lua_tostring(l, 2))
	port := int(C.lua_tointeger(l, 3))
	C.lua_getfield(l, 1, cStrLuaTableKey_tcp_ctimeout)
	timeout := int(C.lua_tointeger(l, -1))
	var poolKey string
	if lua_isstring(l, 4) {
		poolKey = C.GoString(lua_tostring(l, 4))
	} else {
		poolKey = fmt.Sprintf("%s:%d", host, port)
	}

	conn, err := ctx.getConnFromPool(host, port, timeout, poolKey)
	if nil != conn {
		C.goluaLua_pushlightuserdata(l, C.ulong(uintptr(unsafe.Pointer(conn))))
		C.lua_setfield(l, 1, cStrLuaTableKey_tcp_conn)
		C.lua_pushboolean(l, 1)
		C.lua_pushnil(l)
	} else {
		C.lua_pushboolean(l, 0)
		luaApiPushErr(l, err)
	}
	return 2
}

//export GoluaLuaApi_tcp_getreusedtimes
func GoluaLuaApi_tcp_getreusedtimes(l luaState) C.int {
	// n, err = go.tcp_getreusedtimes(conn)

	C.lua_getfield(l, 1, cStrLuaTableKey_tcp_conn)
	conn := (*tcpConn)(C.lua_touserdata(l, -1))

	C.lua_pushinteger(l, C.lua_Integer(conn.reuseCnt))
	C.lua_pushnil(l)
	return 2
}

func luaApiTcpReceiveInner(conn *tcpConn, sz int, timeout int) ([]byte, error) {
	if 0 == timeout {
		conn.SetReadDeadline(time.Time{})
	} else {
		conn.SetReadDeadline(time.Now().Add(time.Duration(timeout) * time.Millisecond))
	}

	var bytes []byte
	var err error
	if sz == 0 {
		reader := bufio.NewReader(conn)
		bytes, err = reader.ReadBytes('\n')
		if len(bytes) > 0 && bytes[len(bytes)-1] == '\r' {
			bytes = bytes[:len(bytes)-1]
		}
	} else if sz < 0 {
		bytes, err = ioutil.ReadAll(conn)
	} else {
		sz := int(sz)
		reader := bufio.NewReader(conn)
		bytes = make([]byte, 0, sz)
		for i := 0; i < sz; i = i + 1 {
			var b byte
			b, err = reader.ReadByte()
			if nil != err {
				break
			}
			bytes[i] = b
		}
	}

	if io.EOF == err {
		err = nil
	}
	return bytes, err
}

//export GoluaLuaApi_tcp_receive
func GoluaLuaApi_tcp_receive(l luaState) C.int {
	// data, err, partial = go.tcp_receive(conn, size)

	C.lua_getfield(l, 1, cStrLuaTableKey_tcp_conn)
	conn := (*tcpConn)(C.lua_touserdata(l, -1))
	sz := int(C.lua_tonumber(l, 2))
	if 0 == sz {
		s := luaToGoString(l, 2)
		if "*a" == s {
			sz = -1
		}
	}
	C.lua_getfield(l, 1, cStrLuaTableKey_tcp_rtimeout)
	timeout := int(C.lua_tointeger(l, -1))

	buf, err := luaApiTcpReceiveInner(conn, sz, timeout)
	if nil == err {
		luaPushGoBytes(l, buf)
		C.lua_pushnil(l)
		C.lua_pushnil(l)
	} else {
		C.lua_pushnil(l)
		luaPushGoString(l, err.Error())
		luaPushGoBytes(l, buf)
	}
	return 3
}

//export GoluaLuaApi_tcp_send
func GoluaLuaApi_tcp_send(l luaState) C.int {
	// n_bytes, err = go.tcp_send(conn, data)

	C.lua_getfield(l, 1, cStrLuaTableKey_tcp_conn)
	conn := (*tcpConn)(C.lua_touserdata(l, -1))
	buf := luaToGoBytes(l, 2)
	C.lua_getfield(l, 1, cStrLuaTableKey_tcp_wtimeout)
	timeout := int(C.lua_tointeger(l, -1))

	if 0 == timeout {
		conn.SetWriteDeadline(time.Time{})
	} else {
		conn.SetWriteDeadline(time.Now().Add(time.Duration(timeout) * time.Millisecond))
	}
	n, err := conn.Write(buf)
	C.lua_pushinteger(l, C.lua_Integer(n))
	luaApiPushErr(l, err)
	return 2
}

//export GoluaLuaApi_tcp_setkeepalive
func GoluaLuaApi_tcp_setkeepalive(l luaState) C.int {
	// ok, err = go.tcp_setkeepalive(conn, _, _)

	C.lua_getfield(l, C.LUA_REGISTRYINDEX, cStrLuaRegKey_vmCtx)
	ctx := (*vmCtx)(C.lua_touserdata(l, -1))
	C.lua_getfield(l, 1, cStrLuaTableKey_tcp_conn)
	conn := (*tcpConn)(C.lua_touserdata(l, -1))

	var ok C.int = 1
	err := ctx.returnConnToPool(conn)
	if nil != err {
		ok = 0
	}
	C.lua_pushboolean(l, ok)
	luaApiPushErr(l, err)
	return 2
}

//export GoluaLuaApi_tcp_settimeout
func GoluaLuaApi_tcp_settimeout(l luaState) C.int {
	// go.tcp_settimeout(conn, ms)

	C.lua_pushinteger(l, C.lua_tointeger(l, 2))
	C.lua_setfield(l, 1, cStrLuaTableKey_tcp_ctimeout)
	C.lua_setfield(l, 1, cStrLuaTableKey_tcp_wtimeout)
	C.lua_setfield(l, 1, cStrLuaTableKey_tcp_rtimeout)
	return 0
}

//export GoluaLuaApi_tcp_settimeouts
func GoluaLuaApi_tcp_settimeouts(l luaState) C.int {
	// go.settimeouts(conn, connect_timeout_ms, write_timeout_ms, read_timeout_ms)

	C.lua_pushinteger(l, C.lua_tointeger(l, 2))
	C.lua_setfield(l, 1, cStrLuaTableKey_tcp_ctimeout)
	C.lua_pushinteger(l, C.lua_tointeger(l, 3))
	C.lua_setfield(l, 1, cStrLuaTableKey_tcp_wtimeout)
	C.lua_pushinteger(l, C.lua_tointeger(l, 4))
	C.lua_setfield(l, 1, cStrLuaTableKey_tcp_rtimeout)
	return 0
}
