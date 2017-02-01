// Copyright 2017 Lynn Deng (ldeng7)

#ifndef goluaClFactory_h
#define goluaClFactory_h

#include "goluaCommon.h"

typedef struct {
    int         sent_begin;
    int         sent_end;
    const char *buf;
    size_t      size;
} GoluaClFactoryBufferCtx;

static const char* goluaClFactoryGetS(lua_State *L, void *ud, size_t *size) {
    GoluaClFactoryBufferCtx *ctx = ud;
    if (0 == ctx->sent_begin) {
        ctx->sent_begin = 1;
        *size = 21;
        return "return function(ctx)\n";
    }
    if (0 == ctx->size) {
        if (0 == ctx->sent_end) {
            ctx->sent_end = 1;
            *size = 4;
            return "\nend";
        }
        return NULL;
    }
    *size = ctx->size;
    ctx->size = 0;
    return ctx->buf;
}

static int goluaClFactoryLoadBuffer(lua_State *L, const char *buff, size_t size, const char *name) {
    GoluaClFactoryBufferCtx ctx;
    ctx.buf = buff;
    ctx.size = size;
    ctx.sent_begin = 0;
    ctx.sent_end = 0;
    return lua_load(L, goluaClFactoryGetS, &ctx, name);
}

#endif
