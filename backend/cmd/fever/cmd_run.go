package main

import (
	"fmt"
	"os/signal"
	"syscall"

	"github.com/robertkozin/feverpanel/backend/lib/tr"
)

type RunCmd struct {
}

func (runCmd *RunCmd) Run(cmdEnv CmdEnv) (err error) {
	ctx, span := tracer.Start(cmdEnv.ctx, "cmd_run")
	defer tr.End(span, &err)

	ctx, cancel := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	fmt.Println("RUN!!!")
	return nil
}
