package web

import (
	"fmt"
	"io"
	"os"
	"os/exec"

	"github.com/creack/pty"
	"github.com/gorilla/websocket"
)

type MinecraftServer struct {
	listener *Broker
	history  *History

	cmd *exec.Cmd
	pty *os.File
}

func NewMinecraftServer(cmd *exec.Cmd, listeners ...io.WriteCloser) *MinecraftServer {
	history := NewHistory()

	return &MinecraftServer{
		listener: NewBroker(append([]io.WriteCloser{history}, listeners...)...),
		history:  history,
		cmd:      cmd,
	}
}

func (mc *MinecraftServer) Start() (err error) {
	mc.cmd.Env = append(os.Environ(), "TERM=xterm-256color")

	mc.pty, err = pty.StartWithSize(mc.cmd, &pty.Winsize{Cols: 100, Rows: 24})
	if err != nil {
		return fmt.Errorf("failed to start minecraft server: %v", err)
	}

	// forward output to listeners
	go func() {
		_, _ = io.Copy(mc.listener, mc.pty)
	}()

	err = mc.cmd.Wait()
	if err != nil {
		return fmt.Errorf("minecraft server exited with error: %v", err)
	}

	err = mc.pty.Close()
	if err != nil {
		return fmt.Errorf("failed to close pty: %v", err)
	}

	mc.pty = nil
	mc.cmd = nil

	return
}

func (mc *MinecraftServer) HandleWebSocket(conn *websocket.Conn) {
	ws := &WSWriter{conn: conn}

	mc.history.Replay(ws)

	mc.AddListener(ws)
	defer mc.listener.Remove(ws)

	for {
		messageType, data, err := conn.ReadMessage()
		if err != nil {
			break
		}

		if messageType == websocket.TextMessage {
			_, err = mc.Write(data)
			if err != nil {
				break
			}
		}
	}
}

func (mc *MinecraftServer) Stop() {
	mc.Write([]byte("stop\n"))
	if mc.cmd != nil {
		mc.cmd.Wait()
	}
}

func (mc *MinecraftServer) Write(data []byte) (int, error) {
	if mc.pty == nil {
		return 0, fmt.Errorf("minecraft server is not running")
	}

	return mc.pty.Write(data)
}

func (mc *MinecraftServer) AddListener(writer io.WriteCloser) {
	mc.listener.Add(writer)
}

func (mc *MinecraftServer) RemoveListener(writer io.WriteCloser) {
	mc.listener.Remove(writer)
}
