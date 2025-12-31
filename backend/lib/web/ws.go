package web

import "github.com/gorilla/websocket"

type WSWriter struct {
	conn *websocket.Conn
}

func NewWSWriter(conn *websocket.Conn) *WSWriter {
	return &WSWriter{conn: conn}
}

func (w *WSWriter) Write(data []byte) (int, error) {
	err := w.conn.WriteMessage(websocket.TextMessage, data)
	if err != nil {
		return 0, err
	}
	return len(data), nil
}

func (w *WSWriter) Close() error {
	return w.conn.Close()
}
