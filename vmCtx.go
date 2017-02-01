// Copyright 2017 Lynn Deng (ldeng7)

package golua

import (
	"log"
	"regexp"
)

type logger struct {
	stdLogger *log.Logger
	minLevel  int
}

const (
	LogLevelEmerg = iota
	LogLevelAlert
	LogLevelCrit
	LogLevelErr
	LogLevelWarn
	LogLevelNotice
	LogLevelInfo
	LogLevelDebug
)

func (logger *logger) init(args *VmInitArgs) {
	logger.stdLogger = args.Logger
	logger.minLevel = args.LogLevel
}

type goObjRefs struct {
	conns map[*tcpConn]bool
}

func (refs *goObjRefs) init() {
	refs.conns = map[*tcpConn]bool{}
}

func (refs *goObjRefs) deinit() {
	refs.conns = nil
}

type vmCtx struct {
	logger     logger
	goObjRefs  goObjRefs
	connPools  map[string]*connPool
	regexCache map[string]*regexp.Regexp
}

func (ctx *vmCtx) init(args *VmInitArgs) {
	(&ctx.logger).init(args)
	(&ctx.goObjRefs).init()
	ctx.connPools = map[string]*connPool{}
	ctx.regexCache = map[string]*regexp.Regexp{}
}

func (ctx *vmCtx) deinit() {
	(&ctx.goObjRefs).deinit()
}
