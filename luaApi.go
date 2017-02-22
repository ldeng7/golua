// Copyright 2017 Lynn Deng (ldeng7)

package golua

/*
#cgo CFLAGS: -I${SRCDIR}/c/include
#cgo LDFLAGS: -L${SRCDIR}/c/lib -lluajit-5.1 -ldl -lm
#include <goluaCommon.h>
#include <goluaLuaApi.h>
*/
import "C"
import (
	"unsafe"
)

var luaApiInts = map[*C.char]int{
	C.CString("EMERG"):  LogLevelEmerg,
	C.CString("ALERT"):  LogLevelAlert,
	C.CString("CRIT"):   LogLevelCrit,
	C.CString("ERR"):    LogLevelErr,
	C.CString("WARN"):   LogLevelWarn,
	C.CString("NOTICE"): LogLevelNotice,
	C.CString("INFO"):   LogLevelInfo,
	C.CString("DEBUG"):  LogLevelDebug,
}

var luaApiFuncs = map[*C.char]unsafe.Pointer{
	// quote
	C.CString("decode_args"):  C.goluaLuaApiC_decode_args,
	C.CString("encode_args"):  C.goluaLuaApiC_encode_args,
	C.CString("escape_uri"):   C.goluaLuaApiC_escape_uri,
	C.CString("quote_sql"):    C.goluaLuaApiC_quote_sql,
	C.CString("unescape_uri"): C.goluaLuaApiC_unescape_uri,

	// regex
	C.CString("re_find"):   C.goluaLuaApiC_re_find,
	C.CString("re_gmatch"): C.goluaLuaApiC_re_gmatch,
	C.CString("re_gsub"):   C.goluaLuaApiC_re_gsub,
	C.CString("re_match"):  C.goluaLuaApiC_re_match,

	// str
	C.CString("crc32"):         C.goluaLuaApiC_crc32,
	C.CString("decode_base64"): C.goluaLuaApiC_decode_base64,
	C.CString("encode_base64"): C.goluaLuaApiC_encode_base64,
	C.CString("hmac_sha1"):     C.goluaLuaApiC_hmac_sha1,
	C.CString("md5"):           C.goluaLuaApiC_md5,
	C.CString("md5_bin"):       C.goluaLuaApiC_md5_bin,
	C.CString("sha1_bin"):      C.goluaLuaApiC_sha1_bin,

	// sys
	C.CString("log"):   C.goluaLuaApiC_log,
	C.CString("pid"):   C.goluaLuaApiC_pid,
	C.CString("tid"):   C.goluaLuaApiC_tid,
	C.CString("yield"): C.goluaLuaApiC_yield,

	// tcp
	C.CString("tcp"):                C.goluaLuaApiC_tcp,
	C.CString("tcp_close"):          C.goluaLuaApiC_tcp_close,
	C.CString("tcp_connect"):        C.goluaLuaApiC_tcp_connect,
	C.CString("tcp_getreusedtimes"): C.goluaLuaApiC_tcp_getreusedtimes,
	C.CString("tcp_receive"):        C.goluaLuaApiC_tcp_receive,
	C.CString("tcp_send"):           C.goluaLuaApiC_tcp_send,
	C.CString("tcp_setkeepalive"):   C.goluaLuaApiC_tcp_setkeepalive,
	C.CString("tcp_settimeout"):     C.goluaLuaApiC_tcp_settimeout,
	C.CString("tcp_settimeouts"):    C.goluaLuaApiC_tcp_settimeouts,
}

func luaApiPushErr(l luaState, err error) {
	if nil != err {
		luaPushGoString(l, err.Error())
	} else {
		C.lua_pushnil(l)
	}
}

func (vs *VmState) injectLuaApi() {
	l := vs.Vm
	C.lua_createtable(l, 0, 32)

	C.lua_pushlightuserdata(l, unsafe.Pointer(uintptr(0)))
	C.lua_setfield(l, -2, C.CString("null"))
	for name, i := range luaApiInts {
		C.lua_pushinteger(l, C.lua_Integer(i))
		C.lua_setfield(l, -2, name)
	}
	for name, fun := range luaApiFuncs {
		C.lua_pushcclosure(l, C.lua_CFunction(fun), 0)
		C.lua_setfield(l, -2, name)
	}

	lua_getglobal(l, cStr_package)
	C.lua_getfield(l, -1, cStr_loaded)
	C.lua_pushvalue(l, -3)
	C.lua_setfield(l, -2, cStr_go)
	lua_pop(l, 2)
	lua_setglobal(l, cStr_go)
}
