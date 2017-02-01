// Copyright 2017 Lynn Deng (ldeng7)

package golua

/*
#cgo CFLAGS: -I${SRCDIR}/c/include
#cgo LDFLAGS: -L${SRCDIR}/c/lib -lluajit-5.1 -ldl -lm
#include <goluaCommon.h>
*/
import "C"
import (
	"log"
	"sync"
	"unsafe"
)

type VmState struct {
	Vm        *C.lua_State
	newCoLock sync.Mutex
	ctx       vmCtx
}

type VmInitArgs struct {
	LuaPath  string
	LuaCpath string
	Logger   *log.Logger
	LogLevel int
}

type VmRunArgs struct {
	Src      string
	CacheKey string
	Name     string
	Ctx      map[string]interface{}
	OutCtx   map[string]interface{}
}

// private

type luaState *C.lua_State

var (
	cStr_package = C.CString("package")
	cStr_loaded  = C.CString("loaded")
	cStr_path    = C.CString("path")
	cStr_cpath   = C.CString("cpath")
	cStr_go      = C.CString("go")

	cStrLuaRegKey_codeCache = C.CString("golua_cc")
	cStrLuaRegKey_vmCtx     = C.CString("golua_vc")
	cStrLuaRegKey_coRef     = C.CString("golua_cr")
)

func lua_isnumber(l luaState, index C.int) bool {
	return C.lua_type(l, index) == C.LUA_TNUMBER
}

func lua_isstring(l luaState, index C.int) bool {
	return C.lua_type(l, index) == C.LUA_TSTRING
}

func lua_istable(l luaState, index C.int) bool {
	return C.lua_type(l, index) == C.LUA_TTABLE
}

func lua_isfunction(l luaState, index C.int) bool {
	return C.lua_type(l, index) == C.LUA_TFUNCTION
}

func lua_pop(l luaState, index C.int) {
	C.lua_settop(l, -index-1)
}

func lua_getglobal(l luaState, s *C.char) {
	C.lua_getfield(l, C.LUA_GLOBALSINDEX, s)
}

func lua_setglobal(l luaState, s *C.char) {
	C.lua_setfield(l, C.LUA_GLOBALSINDEX, s)
}

func lua_tostring(l luaState, index C.int) *C.char {
	return C.lua_tolstring(l, index, (*C.size_t)(unsafe.Pointer(uintptr(0))))
}

func luaPushGoString(l luaState, str string) {
	cStr := C.CString(str)
	C.lua_pushlstring(l, cStr, C.size_t(len(str)))
	C.free(unsafe.Pointer(cStr))
}

func luaPushGoBytes(l luaState, bytes []byte) {
	cStr := C.CBytes(bytes)
	C.lua_pushlstring(l, (*C.char)(cStr), C.size_t(len(bytes)))
	C.free(unsafe.Pointer(cStr))
}

func luaToGoString(l luaState, index C.int) string {
	var sz C.size_t
	cStr := C.lua_tolstring(l, index, &sz)
	return C.GoStringN(cStr, C.int(sz))
}

func luaToGoBytes(l luaState, index C.int) []byte {
	var sz C.size_t
	cStr := C.lua_tolstring(l, index, &sz)
	return C.GoBytes(unsafe.Pointer(cStr), C.int(sz))
}

func freeCStr(s *C.char) {
	C.free(unsafe.Pointer(s))
}
