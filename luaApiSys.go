// Copyright 2017 Lynn Deng (ldeng7)

package golua

/*
#include <unistd.h>
#include <pthread.h>
#include <goluaCommon.h>
*/
import "C"
import (
	"fmt"
	"runtime"
	"strconv"
	"strings"
)

var (
	cStr_luaInfoOpts  = C.CString("Snl")
	cStr_luaMMTotring = C.CString("__tostring")
)

var luaApiSysLogLevelName = map[int]string{
	LogLevelEmerg:  "emerg",
	LogLevelAlert:  "alert",
	LogLevelCrit:   "crit",
	LogLevelErr:    "error",
	LogLevelWarn:   "warn",
	LogLevelNotice: "notice",
	LogLevelInfo:   "info",
	LogLevelDebug:  "debug",
}

//export GoluaLuaApi_log
func GoluaLuaApi_log(l luaState) C.int {
	// go.log(log_level, ...)

	nargs := C.lua_gettop(l)
	C.lua_getfield(l, C.LUA_REGISTRYINDEX, cStrLuaRegKey_vmCtx)
	ctx := (*vmCtx)(C.lua_touserdata(l, -1))
	level := int(C.lua_tonumber(l, 1))

	logger := ctx.logger.stdLogger
	if nil == logger || level > ctx.logger.minLevel {
		return 0
	}

	parts := make([]string, 9+nargs)
	parts[0] = "[lua "
	if levelName, ok := luaApiSysLogLevelName[level]; ok {
		parts[1] = levelName
	} else {
		parts[1] = "?"
	}
	parts[2] = "] "
	var ar C.lua_Debug
	C.lua_getstack(l, 1, &ar)
	C.lua_getinfo(l, cStr_luaInfoOpts, &ar)
	parts[3] = C.GoString((*C.char)(&(ar.short_src[0])))
	parts[4] = ":"
	line := ar.currentline
	if 0 == line {
		line = ar.linedefined
	}
	parts[5] = strconv.Itoa(int(line))
	parts[6] = ": "
	if nil != ar.name {
		parts[7] = C.GoString(ar.name)
		parts[8] = "(): "
	} else {
		parts[7] = ""
		parts[8] = ""
	}

	j := 9
	for i := C.int(2); i <= nargs; i = i + 1 {
		switch C.lua_type(l, i) {
		case C.LUA_TNIL:
			parts[j] = "nil"
		case C.LUA_TBOOLEAN:
			if 1 == C.lua_toboolean(l, i) {
				parts[j] = "true"
			} else {
				parts[j] = "false"
			}
		case C.LUA_TNUMBER, C.LUA_TSTRING:
			parts[j] = luaToGoString(l, i)
		case C.LUA_TTABLE:
			if 1 == C.luaL_callmeta(l, i, cStr_luaMMTotring) {
				parts[j] = luaToGoString(l, -1)
			} else {
				parts[j] = fmt.Sprintf("table@0x%016x", C.lua_topointer(l, i))
			}
		case C.LUA_TFUNCTION:
			parts[j] = fmt.Sprintf("function@0x%016x", C.lua_topointer(l, i))
		case C.LUA_TUSERDATA:
			parts[j] = fmt.Sprintf("userdate@0x%016x", C.lua_topointer(l, i))
		case C.LUA_TTHREAD:
			parts[j] = fmt.Sprintf("thread@0x%016x", C.lua_topointer(l, i))
		case C.LUA_TLIGHTUSERDATA:
			parts[j] = fmt.Sprintf("lightuserdate@0x%016x", C.lua_topointer(l, i))
		}
		j = j + 1
	}
	parts[j] = "\n"

	logger.Output(2, strings.Join(parts, ""))
	return 0
}

//export GoluaLuaApi_pid
func GoluaLuaApi_pid(l luaState) C.int {
	// pid = go.pid()

	C.lua_pushinteger(l, C.lua_Integer(C.getpid()))
	return 1
}

//export GoluaLuaApi_tid
func GoluaLuaApi_tid(l luaState) C.int {
	// tid = go.tid()

	C.lua_pushinteger(l, C.lua_Integer(C.pthread_self()))
	return 1
}

//export GoluaLuaApi_yield
func GoluaLuaApi_yield(l luaState) C.int {
	// go.yield()

	runtime.Gosched()
	return 0
}
