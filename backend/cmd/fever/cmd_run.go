package main

import (
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"syscall"

	"github.com/robertkozin/feverpanel/backend/lib/tr"
	"github.com/robertkozin/feverpanel/backend/lib/web"
)

type RunCmd struct {
}

func (runCmd *RunCmd) Run(cmdEnv CmdEnv) (err error) {
	ctx, span := tracer.Start(cmdEnv.ctx, "cmd_run")
	defer tr.End(span, &err)

	ctx, cancel := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	fmt.Println("RUN!!!")

	cmd := exec.Command("java", "-Xmx4g", "-jar", "server.jar", "nogui")
	mc := web.NewMinecraftServer(cmd, os.Stdout)

	go func() {
		if err := mc.Start(); err != nil {
			fmt.Printf("Minecraft server exited with error: %v\n", err)
		}
	}()

	http.HandleFunc("/ws", web.ServeWebSocket(mc))

	server := &http.Server{Addr: ":8080"}

	go func() {
		fmt.Println("Starting web server on :8080")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Printf("Web server error: %v\n", err)
		}
	}()

	<-ctx.Done()
	fmt.Println("\nShutting down...")

	mc.Stop()
	server.Close()

	return nil
}
