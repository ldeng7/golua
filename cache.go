// Copyright 2017 Lynn Deng (ldeng7)

package golua

/*
#cgo CFLAGS: -I${SRCDIR}/c/include
#cgo LDFLAGS: -L${SRCDIR}/c/lib -lluajit-5.1 -ldl -lm
#include <goluaCommon.h>
#include <goluaClFactory.h>
*/
import "C"

func (vs *VmState) cacheLoadCode(key *C.char) bool {
	l := vs.Vm
	C.lua_getfield(l, C.LUA_REGISTRYINDEX, cStrLuaRegKey_codeCache)
	C.lua_getfield(l, -1, key)
	if lua_isfunction(l, -1) {
		C.lua_remove(l, -2)
		return true
	}
	lua_pop(l, 2)
	return false
}

func (vs *VmState) cacheStoreCode(key *C.char) {
	l := vs.Vm
	C.lua_getfield(l, C.LUA_REGISTRYINDEX, cStrLuaRegKey_codeCache)
	C.lua_pushvalue(l, -2)
	C.lua_setfield(l, -2, key)
	lua_pop(l, 1)
}

func (vs *VmState) cacheLoadBuffer(src string, cacheKey string, name string) {
	keyPtr := C.CString(cacheKey)
	defer freeCStr(keyPtr)
	if vs.cacheLoadCode(keyPtr) {
		return
	}
	srcPtr := C.CString(src)
	namePtr := C.CString(name)
	C.goluaClFactoryLoadBuffer(vs.Vm, srcPtr, C.size_t(len(src)), namePtr)
	vs.cacheStoreCode(keyPtr)
	freeCStr(srcPtr)
	freeCStr(namePtr)
}
