package main

import (
	"encoding/json"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"sync"

	"github.com/ldeng7/golua"
)

func run(v *golua.VmState, code string, key string) {
	r := &golua.VmRunArgs{
		Src:      code,
		CacheKey: key,
		Name:     "primary",
		Ctx:      map[string]interface{}{},
	}
	v.Run(r)
	js, _ := json.Marshal(r.OutCtx)
	println(string(js))
}

func main() {
	path, _ := exec.LookPath(os.Args[0])
	path = filepath.Dir(path)
	v := golua.NewVm(&golua.VmInitArgs{
		LuaPath:  filepath.Join(path, "../lualib/?.lua") + ";;",
		LuaCpath: ";;",
		Logger:   log.New(os.Stdout, "", log.Ldate|log.Ltime),
		LogLevel: golua.LogLevelNotice,
	})
	go run(v, `
		--local c = go.tcp()
		--ctx.ok, ctx.err1 = go.tcp_connect(c, 'webtcp.tongxinmao.com', 10002)
		--ctx.n, ctx.err2 = go.tcp_send(c, 'ldeng')
		--ctx.pid = go.pid()

		local mysql = require "mysql"
		local db, _ = mysql:new()
		db:connect({host = "127.0.0.1", port = 3306, database = "test", user = "root", password = "abcabc"})
		local res, _ = db:query("SELECT count(1) AS cnt FROM tests LIMIT 2;")
		ctx.cnt = res[1].cnt

		return ctx`, "test")
	var wg sync.WaitGroup
	wg.Add(1)
	wg.Wait()
}
