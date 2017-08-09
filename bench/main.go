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

		local mysql = require "redis"
			local db, err = mysql:new()
			go.log(1,err)
			_,err=db:connect("127.0.0.1", 6379)
			go.log(1,err)
			local res, err = db:get("ldeng")
			go.log(1,err)
			ctx.res=res

		return ctx`, "test")
	var wg sync.WaitGroup
	wg.Add(1)
	wg.Wait()
}
