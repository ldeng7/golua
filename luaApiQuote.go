// Copyright 2017 Lynn Deng (ldeng7)

package golua

/*
#include <goluaCommon.h>
*/
import "C"
import (
	"net/url"
)

//export GoluaLuaApi_decode_args
func GoluaLuaApi_decode_args(l luaState) C.int {
	// t = go.decode_args(s)

	m, _ := url.ParseQuery(luaToGoString(l, 1))
	C.lua_createtable(l, 0, C.int(len(m)))
	for k, v := range m {
		lv := C.int(len(v))
		if 1 == lv {
			luaPushGoString(l, v[0])
		} else if lv > 1 {
			C.lua_createtable(l, lv, 0)
			for i := C.int(0); i < lv; i = i + 1 {
				luaPushGoString(l, v[i])
				C.lua_rawseti(l, -2, i+1)
			}
		} else {
			luaPushGoString(l, "")
		}
		kPtr := C.CString(k)
		C.lua_setfield(l, -2, kPtr)
		freeCStr(kPtr)
	}
	return 1
}

func luaApiStrEncodeArgsGetPrimitiveValue(l luaState, typ C.int) string {
	switch typ {
	case C.LUA_TBOOLEAN:
		if 1 == C.lua_toboolean(l, -1) {
			return "true"
		} else {
			return "false"
		}
	case C.LUA_TNUMBER, C.LUA_TSTRING:
		return C.GoString(lua_tostring(l, -1))
	}
	return ""
}

//export GoluaLuaApi_encode_args
func GoluaLuaApi_encode_args(l luaState) C.int {
	// s = go.encode_args(t)

	if !lua_istable(l, 1) {
		luaPushGoString(l, "")
		return 1
	}
	C.lua_pushnil(l)
	m := url.Values{}
	for {
		if 0 == C.lua_next(l, -2) {
			break
		}
		if !lua_isstring(l, -2) {
			lua_pop(l, 1)
			continue
		}
		k := C.GoString(lua_tostring(l, -2))
		typ := C.lua_type(l, -1)
		if C.LUA_TTABLE != typ {
			v := luaApiStrEncodeArgsGetPrimitiveValue(l, typ)
			m[k] = []string{v}
		} else {
			C.lua_pushnil(l)
			ar := []string{}
			for {
				if 0 == C.lua_next(l, -2) {
					break
				}
				v := luaApiStrEncodeArgsGetPrimitiveValue(l, C.lua_type(l, -1))
				ar = append(ar, v)
				lua_pop(l, 1)
			}
			m[k] = ar
		}
		lua_pop(l, 1)
	}
	luaPushGoString(l, m.Encode())
	return 1
}

//export GoluaLuaApi_escape_uri
func GoluaLuaApi_escape_uri(l luaState) C.int {
	// digest = go.escape_uri(s)

	outStr := url.QueryEscape(luaToGoString(l, 1))
	luaPushGoString(l, outStr)
	return 1
}

//export GoluaLuaApi_quote_sql
func GoluaLuaApi_quote_sql(l luaState) C.int {
	// digest = go.quote_sql(s)

	s := luaToGoBytes(l, 1)
	outBytes := make([]byte, len(s)*2+2)
	j := 0
	outBytes[j] = '\''
	j = j + 1
	for _, c := range s {
		switch c {
		case 0:
			outBytes[j] = '\\'
			outBytes[j+1] = '0'
			j = j + 2
		case '\b':
			outBytes[j] = '\\'
			outBytes[j+1] = 'b'
			j = j + 2
		case '\n':
			outBytes[j] = '\\'
			outBytes[j+1] = 'n'
			j = j + 2
		case '\r':
			outBytes[j] = '\\'
			outBytes[j+1] = 'r'
			j = j + 2
		case '\t':
			outBytes[j] = '\\'
			outBytes[j+1] = 't'
			j = j + 2
		case 26:
			outBytes[j] = '\\'
			outBytes[j+1] = 'Z'
			j = j + 2
		case '\\':
			outBytes[j] = '\\'
			outBytes[j+1] = '\\'
			j = j + 2
		case '\'':
			outBytes[j] = '\\'
			outBytes[j+1] = '\''
			j = j + 2
		case '"':
			outBytes[j] = '\\'
			outBytes[j+1] = '"'
			j = j + 2
		default:
			outBytes[j] = c
			j = j + 1
		}
	}
	outBytes[j] = '\''
	j = j + 1
	luaPushGoBytes(l, outBytes[:j])
	return 1
}

//export GoluaLuaApi_unescape_uri
func GoluaLuaApi_unescape_uri(l luaState) C.int {
	// s = go.unescape_uri(digest)

	outStr, _ := url.QueryUnescape(luaToGoString(l, 1))
	luaPushGoString(l, outStr)
	return 1
}
