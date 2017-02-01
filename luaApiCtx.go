// Copyright 2017 Lynn Deng (ldeng7)

package golua

/*
#cgo CFLAGS: -I${SRCDIR}/c/include
#cgo LDFLAGS: -L${SRCDIR}/c/lib -lluajit-5.1 -ldl -lm
#include <goluaCommon.h>
*/
import "C"

func luaApiCtxInjectValue(l luaState, v interface{}) bool {
	valid := true
	switch vv := v.(type) {
	case nil:
		C.lua_pushnil(l)
	case bool:
		if vv {
			C.lua_pushboolean(l, 1)
		} else {
			C.lua_pushboolean(l, 0)
		}
	case int:
		C.lua_pushinteger(l, C.lua_Integer(vv))
	case int64:
		C.lua_pushinteger(l, C.lua_Integer(vv))
	case float64:
		C.lua_pushnumber(l, C.lua_Number(vv))
	case string:
		luaPushGoString(l, vv)
	case []interface{}:
		C.lua_createtable(l, C.int(len(vv)), 0)
		for i, e := range vv {
			if !luaApiCtxInjectValue(l, e) {
				continue
			}
			C.lua_rawseti(l, -2, C.int(i+1))
		}
	case map[int]interface{}:
		C.lua_createtable(l, 0, C.int(len(vv)))
		for k, e := range vv {
			if !luaApiCtxInjectValue(l, e) {
				continue
			}
			C.lua_rawseti(l, -2, C.int(k))
		}
	case map[string]interface{}:
		C.lua_createtable(l, 0, C.int(len(vv)))
		for k, e := range vv {
			if !luaApiCtxInjectValue(l, e) {
				continue
			}
			kPtr := C.CString(k)
			C.lua_setfield(l, -2, kPtr)
			freeCStr(kPtr)
		}
	default:
		valid = false
	}
	return valid
}

func luaApiCtxInject(l luaState, ctx map[string]interface{}) {
	C.lua_createtable(l, 0, C.int(len(ctx)))
	for k, v := range ctx {
		if !luaApiCtxInjectValue(l, v) {
			continue
		}
		kPtr := C.CString(k)
		C.lua_setfield(l, -2, kPtr)
		freeCStr(kPtr)
	}
}

func luaApiCtxExjectValue(l luaState) interface{} {
	var out interface{}
	switch C.lua_type(l, -1) {
	case C.LUA_TNIL:
		out = nil
	case C.LUA_TBOOLEAN:
		out = 1 == C.lua_toboolean(l, -1)
	case C.LUA_TNUMBER:
		f := float64(C.lua_tonumber(l, -1))
		i := int64(f)
		if float64(i) == f {
			out = i
		} else {
			out = f
		}
	case C.LUA_TSTRING:
		out = luaToGoString(l, -1)
	case C.LUA_TTABLE:
		m := map[string]interface{}{}
		C.lua_pushnil(l)
		for {
			if 0 == C.lua_next(l, -2) {
				break
			}
			var k string
			switch C.lua_type(l, -2) {
			case C.LUA_TNUMBER, C.LUA_TSTRING:
				k = luaToGoString(l, -2)
			default:
				lua_pop(l, 1)
				continue
			}
			m[k] = luaApiCtxExjectValue(l)
			lua_pop(l, 1)
		}
		out = m
	default:
		out = nil
	}
	return out
}

func luaApiCtxExject(l luaState) map[string]interface{} {
	if !lua_istable(l, -1) {
		return nil
	}
	C.lua_pushnil(l)
	out := map[string]interface{}{}
	for {
		if 0 == C.lua_next(l, -2) {
			break
		}
		if lua_isstring(l, -2) {
			k := C.GoString(lua_tostring(l, -2))
			out[k] = luaApiCtxExjectValue(l)
		}
		lua_pop(l, 1)
	}
	return out
}
