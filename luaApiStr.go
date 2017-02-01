// Copyright 2017 Lynn Deng (ldeng7)

package golua

/*
#include <unistd.h>
#include <pthread.h>
#include <goluaCommon.h>
*/
import "C"
import (
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"hash/crc32"
)

//export GoluaLuaApi_crc32
func GoluaLuaApi_crc32(l luaState) C.int {
	// i = go.crc32(s)

	i := crc32.ChecksumIEEE(luaToGoBytes(l, 1))
	C.lua_pushinteger(l, C.lua_Integer(i))
	return 1
}

//export GoluaLuaApi_decode_base64
func GoluaLuaApi_decode_base64(l luaState) C.int {
	// s = go.decode_base64(digest)

	outBytes, _ := base64.StdEncoding.DecodeString(luaToGoString(l, 1))
	luaPushGoBytes(l, outBytes)
	return 1
}

//export GoluaLuaApi_encode_base64
func GoluaLuaApi_encode_base64(l luaState) C.int {
	// digest = go.encode_base64(s, no_padding)

	noPadding := 1 == C.lua_toboolean(l, 2)
	var enc *base64.Encoding
	if !noPadding {
		enc = base64.StdEncoding
	} else {
		enc = base64.RawStdEncoding
	}
	outStr := enc.EncodeToString(luaToGoBytes(l, 1))
	luaPushGoString(l, outStr)
	return 1
}

//export GoluaLuaApi_hmac_sha1
func GoluaLuaApi_hmac_sha1(l luaState) C.int {
	// digest = go.hmac_sha1(k, s)

	mac := hmac.New(sha1.New, luaToGoBytes(l, 1))
	mac.Write(luaToGoBytes(l, 2))
	luaPushGoBytes(l, mac.Sum(nil)[:])
	return 1
}

//export GoluaLuaApi_md5
func GoluaLuaApi_md5(l luaState) C.int {
	// digest = go.md5(s)

	outBytes := md5.Sum(luaToGoBytes(l, 1))
	luaPushGoString(l, fmt.Sprintf("%x", outBytes))
	return 1
}

//export GoluaLuaApi_md5_bin
func GoluaLuaApi_md5_bin(l luaState) C.int {
	// digest = go.md5_bin(s)

	outBytes := md5.Sum(luaToGoBytes(l, 1))
	luaPushGoBytes(l, outBytes[:])
	return 1
}

//export GoluaLuaApi_sha1_bin
func GoluaLuaApi_sha1_bin(l luaState) C.int {
	// digest = go.sha1_bin(s)

	outBytes := sha1.Sum(luaToGoBytes(l, 1))
	luaPushGoBytes(l, outBytes[:])
	return 1
}
