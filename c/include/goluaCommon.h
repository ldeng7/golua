// Copyright 2017 Lynn Deng (ldeng7)

#ifndef goluaCommon_h
#define goluaCommon_h

#include <stdlib.h>
#include "luajit/lua.h"
#include "luajit/luaconf.h"
#include "luajit/lauxlib.h"
#include "luajit/lualib.h"
#include "luajit/luajit.h"

static void goluaLua_pushlightuserdata(lua_State *L, unsigned long p) {
    lua_pushlightuserdata(L, (void*)p);
}

#endif
