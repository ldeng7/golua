package main

import (
	"encoding/json"
	"log"
	"os"
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
	v := golua.NewVm(&golua.VmInitArgs{
		LuaPath:  "/home/ldeng/1;;",
		LuaCpath: ";;",
		Logger:   log.New(os.Stdout, "", log.Ldate|log.Ltime),
		LogLevel: 2,
	})
	go run(v, `--local c = go.tcp()
	    --ctx.ok, ctx.err1 = go.tcp_connect(c, 'webtcp.tongxinmao.com', 10002)
		--ctx.n, ctx.err2 = go.tcp_send(c, 'ldeng')
		--ctx.pid = go.pid()
		ctx.a=1
		go.log(1,0.0000000000023333)
		return ctx`, "test")
	var wg sync.WaitGroup
	wg.Add(1)
	wg.Wait()
}
