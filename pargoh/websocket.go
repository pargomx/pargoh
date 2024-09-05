package main

import (
	"strconv"
	"sync"

	"github.com/pargomx/gecko"
	"github.com/pargomx/gecko/gko"
	"golang.org/x/net/websocket"
)

type reloader struct {
	sockets []socket
	lastID  int
	mu      sync.Mutex
}

type socket struct {
	id         string
	ws         *websocket.Conn
	historiaID int
}

func (s *reloader) nuevoWS(c *gecko.Context) error {
	websocket.Handler(func(ws *websocket.Conn) {
		defer ws.Close()
		s.mu.Lock()
		s.lastID++
		socket := socket{
			id:         strconv.Itoa(s.lastID),
			ws:         ws,
			historiaID: c.PathInt("historia_id"),
		}
		s.sockets = append(s.sockets, socket)
		// gko.LogDebugf("socket(%v) creado", socket.id)
		s.mu.Unlock()

		// Enviar ID del socket como json
		err := websocket.Message.Send(socket.ws, `{"id":`+socket.id+`}`)
		if err != nil {
			gko.Err(err).Msgf("socket(%v) send", socket.id).Log()
			s.quitar(socket)
		}
		// Recibir mensajes para conservar el socket
		for {
			var msg string
			err := websocket.Message.Receive(ws, &msg)
			if err != nil {
				if err.Error() != "EOF" {
					gko.Err(err).Msgf("socket(%v) receive", socket.id).Log()
				}
				break
			}
			gko.LogDebugf("socket(%v) recived: %s", socket.id, msg)
		}
		// Una vez que se cierra la conexión...
		s.quitar(socket)
	}).ServeHTTP(c.Response(), c.Request())
	return nil
}

func (s *reloader) quitar(ws socket) {
	for i, socket := range s.sockets { // remove closed connection
		if socket.id == ws.id {
			s.sockets = append(s.sockets[:i], s.sockets[i+1:]...)
			// gko.LogDebugf("socket(%v) eliminado2", socket.id)
			break
		}
	}
}

func (s *servidor) brodcastReload(c *gecko.Context) error {
	s.reloader.brodcastReload(c)
	return c.StatusOkf("Reload sent to %d connections", len(s.reloader.sockets))
}

func (s *reloader) brodcastReload(c *gecko.Context) {
	s.mu.Lock()
	defer s.mu.Unlock()
	brodcasted := 0
	for _, socket := range s.sockets {
		if socket.historiaID == c.PathInt("historia_id") ||
			socket.historiaID == c.FormInt("historia_id") {

			id := c.Request().Header.Get("X-SocketID")
			if id == "" {
				gko.LogWarn("X-SocketID empty")
				continue
			}
			if socket.id == id {
				continue
			}
			err := websocket.Message.Send(socket.ws, `{"id":`+socket.id+`,"reload":true}`)
			if err != nil {
				gko.Err(err).Msgf("socket(%v) send", socket.id).Log()
				s.quitar(socket)
			}
			// gko.LogDebugf("socket(%v) reload", socket.id)
			brodcasted++
		}
	}
	// gko.LogDebugf("Broadcasting reload to %d/%d connections", brodcasted, len(s.sockets))
}
