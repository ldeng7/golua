// Copyright 2017 Lynn Deng (ldeng7)

package golua

/*
#include <goluaCommon.h>
*/
import "C"
import (
	"regexp"
)

func luaApiRegexGetRe(l luaState, reStr string, doCache bool) (*regexp.Regexp, error) {
	if !doCache {
		return regexp.Compile(reStr)
	}
	C.lua_getfield(l, C.LUA_REGISTRYINDEX, cStrLuaRegKey_vmCtx)
	ctx := (*vmCtx)(C.lua_touserdata(l, -1))
	re, ok := ctx.regexCache[reStr]
	if ok {
		return re.Copy(), nil
	}
	var err error
	re, err = regexp.Compile(reStr)
	if nil == re {
		return nil, err
	}
	ctx.regexCache[reStr] = re
	return re, nil
}

//export GoluaLuaApi_re_find
func GoluaLuaApi_re_find(l luaState) C.int {
	// from, to, err = go.re_find(s, regex, do_catch, _, nth)

	re, err := luaApiRegexGetRe(l, luaToGoString(l, 2), 1 == C.lua_toboolean(l, 3))
	if nil == re {
		C.lua_pushnil(l)
		C.lua_pushnil(l)
		luaApiPushErr(l, err)
		return 3
	}
	nth := int(C.lua_tonumber(l, 5))
	if nth < 0 {
		nth = 0
	}

	ar := re.FindStringSubmatchIndex(luaToGoString(l, 1))
	from := -1
	if len(ar) >= (nth+1)*2 {
		from = ar[nth*2]
	}
	if from >= 0 {
		C.lua_pushinteger(l, C.lua_Integer(from+1))
		C.lua_pushinteger(l, C.lua_Integer(ar[nth*2+1]))
	} else {
		C.lua_pushboolean(l, 0)
		C.lua_pushboolean(l, 0)
	}
	C.lua_pushnil(l)
	return 3
}

//export GoluaLuaApi_re_gmatch
func GoluaLuaApi_re_gmatch(l luaState) C.int {
	// m, err = go.re_gmatch(s, regex, do_catch, limit)

	re, err := luaApiRegexGetRe(l, luaToGoString(l, 2), 1 == C.lua_toboolean(l, 3))
	if nil == re {
		C.lua_pushnil(l)
		luaApiPushErr(l, err)
		return 2
	}

	lim := int(C.lua_tonumber(l, 4))
	if lim <= 0 {
		lim = -1
	}
	ar := re.FindAllStringSubmatch(luaToGoString(l, 1), lim)
	if len(ar) >= 1 {
		C.lua_createtable(l, C.int(len(ar)), 0)
		for i, e := range ar {
			C.lua_createtable(l, C.int(len(e)), 0)
			for j, ee := range e {
				if len(ee) > 0 {
					luaPushGoString(l, ee)
				} else {
					C.lua_pushboolean(l, 0)
				}
				C.lua_rawseti(l, -2, C.int(j))
			}
			C.lua_rawseti(l, -2, C.int(i+1))
		}
	} else {
		C.lua_pushboolean(l, 0)
	}
	C.lua_pushnil(l)
	return 2
}

//export GoluaLuaApi_re_gsub
func GoluaLuaApi_re_gsub(l luaState) C.int {
	// m, err = go.re_gsub(s, regex, repl, do_catch)

	re, err := luaApiRegexGetRe(l, luaToGoString(l, 2), 1 == C.lua_toboolean(l, 4))
	if nil == re {
		C.lua_pushnil(l)
		luaApiPushErr(l, err)
		return 2
	}

	outStr := re.ReplaceAllString(luaToGoString(l, 1), luaToGoString(l, 3))
	luaPushGoString(l, outStr)
	C.lua_pushnil(l)
	return 2
}

//export GoluaLuaApi_re_match
func GoluaLuaApi_re_match(l luaState) C.int {
	// m, err = go.re_match(s, regex, do_catch)

	re, err := luaApiRegexGetRe(l, luaToGoString(l, 2), 1 == C.lua_toboolean(l, 3))
	if nil == re {
		C.lua_pushnil(l)
		luaApiPushErr(l, err)
		return 2
	}

	ar := re.FindAllStringSubmatch(luaToGoString(l, 1), 1)
	if len(ar) >= 1 {
		C.lua_createtable(l, C.int(len(ar[0])), 0)
		for i, e := range ar[0] {
			if len(e) > 0 {
				luaPushGoString(l, e)
			} else {
				C.lua_pushboolean(l, 0)
			}
			C.lua_rawseti(l, -2, C.int(i))
		}
	} else {
		C.lua_pushboolean(l, 0)
	}
	C.lua_pushnil(l)
	return 2
}
