// Copyright 2017 Lynn Deng (ldeng7)

package golua

/*
#cgo CFLAGS: -I${SRCDIR}/c/include
#cgo LDFLAGS: -L${SRCDIR}/c/lib -lluajit-5.1 -ldl -lm
#include <goluaCommon.h>
*/
import "C"
import (
	"fmt"
	"regexp"
	"unsafe"
)

func setLuaPath(l luaState, field *C.char, s string) {
	C.lua_getfield(l, -1, field)
	re, _ := regexp.Compile(";;")
	s = re.ReplaceAllLiteralString(s, fmt.Sprintf(";%s;", C.GoString(lua_tostring(l, -1))))
	luaPushGoString(l, s)
	C.lua_setfield(l, -3, field)
	lua_pop(l, 1)
}

func (vs *VmState) initRegistry() {
	l := vs.Vm
	C.lua_createtable(l, 0, 8)
	C.lua_setfield(l, C.LUA_REGISTRYINDEX, cStrLuaRegKey_codeCache)
	C.lua_pushlightuserdata(l, unsafe.Pointer(&vs.ctx))
	C.lua_setfield(l, C.LUA_REGISTRYINDEX, cStrLuaRegKey_vmCtx)
	C.lua_createtable(l, 0, 64)
	C.lua_setfield(l, C.LUA_REGISTRYINDEX, cStrLuaRegKey_coRef)
}

func (vs *VmState) initGlobals() {
	vs.injectLuaApi()
}

func (vs *VmState) newState(args *VmInitArgs) {
	l := C.luaL_newstate()
	vs.Vm = l
	C.luaL_openlibs(l)

	lua_getglobal(l, cStr_package)
	setLuaPath(l, cStr_path, args.LuaPath)
	setLuaPath(l, cStr_cpath, args.LuaCpath)
	lua_pop(l, 1)

	vs.initRegistry()
	vs.initGlobals()
}

func NewVm(args *VmInitArgs) *VmState {
	vs := &VmState{}
	vs.newState(args)
	(&vs.ctx).init(args)
	return vs
}

func (vs *VmState) DeinitVm() {
	(&vs.ctx).deinit()
	C.lua_close(vs.Vm)
}

func (vs *VmState) newThread(args *VmRunArgs) (luaState, C.int) {
	l := vs.Vm
	vs.newCoLock.Lock()
	co := C.lua_newthread(l)
	vs.cacheLoadBuffer(args.Src, args.CacheKey, args.Name)
	C.lua_xmove(l, co, 1)

	C.lua_getfield(l, C.LUA_REGISTRYINDEX, cStrLuaRegKey_coRef)
	C.lua_pushvalue(l, -2)
	coRef := C.luaL_ref(l, -2)
	lua_pop(l, 2)
	vs.newCoLock.Unlock()
	return co, coRef
}

func (vs *VmState) Run(args *VmRunArgs) {
	l := vs.Vm
	co, coRef := vs.newThread(args)
	C.lua_pcall(co, 0, 1, 0)
	luaApiCtxInject(co, args.Ctx)
	for {
		if C.LUA_YIELD != C.lua_resume(co, 1) {
			break
		}
	}
	args.OutCtx = luaApiCtxExject(co)
	C.lua_getfield(co, C.LUA_REGISTRYINDEX, cStrLuaRegKey_coRef)
	C.luaL_unref(l, -1, coRef)
}
