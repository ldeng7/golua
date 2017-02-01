// Copyright 2017 Lynn Deng (ldeng7)

#ifndef goluaLuaApi_h
#define goluaLuaApi_h

// quote

extern int GoluaLuaApi_decode_args(lua_State *l);
int goluaLuaApiC_decode_args(lua_State *l) {
    return GoluaLuaApi_decode_args(l);
}

extern int GoluaLuaApi_encode_args(lua_State *l);
int goluaLuaApiC_encode_args(lua_State *l) {
    return GoluaLuaApi_encode_args(l);
}

extern int GoluaLuaApi_escape_uri(lua_State *l);
int goluaLuaApiC_escape_uri(lua_State *l) {
    return GoluaLuaApi_escape_uri(l);
}

extern int GoluaLuaApi_quote_sql(lua_State *l);
int goluaLuaApiC_quote_sql(lua_State *l) {
    return GoluaLuaApi_quote_sql(l);
}

extern int GoluaLuaApi_unescape_uri(lua_State *l);
int goluaLuaApiC_unescape_uri(lua_State *l) {
    return GoluaLuaApi_unescape_uri(l);
}

// regex

extern int GoluaLuaApi_re_find(lua_State *l);
int goluaLuaApiC_re_find(lua_State *l) {
    return GoluaLuaApi_re_find(l);
}

extern int GoluaLuaApi_re_gmatch(lua_State *l);
int goluaLuaApiC_re_gmatch(lua_State *l) {
    return GoluaLuaApi_re_gmatch(l);
}

extern int GoluaLuaApi_re_gsub(lua_State *l);
int goluaLuaApiC_re_gsub(lua_State *l) {
    return GoluaLuaApi_re_gsub(l);
}

extern int GoluaLuaApi_re_match(lua_State *l);
int goluaLuaApiC_re_match(lua_State *l) {
    return GoluaLuaApi_re_match(l);
}

// str

extern int GoluaLuaApi_crc32(lua_State *l);
int goluaLuaApiC_crc32(lua_State *l) {
    return GoluaLuaApi_crc32(l);
}

extern int GoluaLuaApi_decode_base64(lua_State *l);
int goluaLuaApiC_decode_base64(lua_State *l) {
    return GoluaLuaApi_decode_base64(l);
}

extern int GoluaLuaApi_encode_base64(lua_State *l);
int goluaLuaApiC_encode_base64(lua_State *l) {
    return GoluaLuaApi_encode_base64(l);
}

extern int GoluaLuaApi_hmac_sha1(lua_State *l);
int goluaLuaApiC_hmac_sha1(lua_State *l) {
    return GoluaLuaApi_hmac_sha1(l);
}

extern int GoluaLuaApi_md5(lua_State *l);
int goluaLuaApiC_md5(lua_State *l) {
    return GoluaLuaApi_md5(l);
}

extern int GoluaLuaApi_md5_bin(lua_State *l);
int goluaLuaApiC_md5_bin(lua_State *l) {
    return GoluaLuaApi_md5_bin(l);
}

extern int GoluaLuaApi_sha1_bin(lua_State *l);
int goluaLuaApiC_sha1_bin(lua_State *l) {
    return GoluaLuaApi_sha1_bin(l);
}

// sys

extern int GoluaLuaApi_log(lua_State *l);
int goluaLuaApiC_log(lua_State *l) {
    return GoluaLuaApi_log(l);
}

extern int GoluaLuaApi_pid(lua_State *l);
int goluaLuaApiC_pid(lua_State *l) {
    return GoluaLuaApi_pid(l);
}

extern int GoluaLuaApi_tid(lua_State *l);
int goluaLuaApiC_tid(lua_State *l) {
    return GoluaLuaApi_tid(l);
}

extern int GoluaLuaApi_yield(lua_State *l);
int goluaLuaApiC_yield(lua_State *l) {
    return GoluaLuaApi_yield(l);
}

// tcp

extern int GoluaLuaApi_tcp(lua_State *l);
int goluaLuaApiC_tcp(lua_State *l) {
    return GoluaLuaApi_tcp(l);
}

extern int GoluaLuaApi_tcp_close(lua_State *l);
int goluaLuaApiC_tcp_close(lua_State *l) {
    return GoluaLuaApi_tcp_close(l);
}

extern int GoluaLuaApi_tcp_connect(lua_State *l);
int goluaLuaApiC_tcp_connect(lua_State *l) {
    return GoluaLuaApi_tcp_connect(l);
}

extern int GoluaLuaApi_tcp_getreusedtimes(lua_State *l);
int goluaLuaApiC_tcp_getreusedtimes(lua_State *l) {
    return GoluaLuaApi_tcp_getreusedtimes(l);
}

extern int GoluaLuaApi_tcp_receive(lua_State *l);
int goluaLuaApiC_tcp_receive(lua_State *l) {
    return GoluaLuaApi_tcp_receive(l);
}

extern int GoluaLuaApi_tcp_send(lua_State *l);
int goluaLuaApiC_tcp_send(lua_State *l) {
    return GoluaLuaApi_tcp_send(l);
}

extern int GoluaLuaApi_tcp_setkeepalive(lua_State *l);
int goluaLuaApiC_tcp_setkeepalive(lua_State *l) {
    return GoluaLuaApi_tcp_setkeepalive(l);
}

extern int GoluaLuaApi_tcp_settimeout(lua_State *l);
int goluaLuaApiC_tcp_settimeout(lua_State *l) {
    return GoluaLuaApi_tcp_settimeout(l);
}

extern int GoluaLuaApi_tcp_settimeouts(lua_State *l);
int goluaLuaApiC_tcp_settimeouts(lua_State *l) {
    return GoluaLuaApi_tcp_settimeouts(l);
}

#endif
