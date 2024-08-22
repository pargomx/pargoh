package main

import (
	"sync"

	"github.com/pargomx/gecko"
	"github.com/pargomx/gecko/gko"
	"golang.org/x/net/websocket"
)

type reloader struct {
	counter int
	sockets []socket
	idcount int
	mu      sync.Mutex
}

type socket struct {
	id int
	ws *websocket.Conn
}

func (s *reloader) nuevoWS(c *gecko.Context) error {
	websocket.Handler(func(ws *websocket.Conn) {
		defer ws.Close()
		s.mu.Lock()
		s.idcount++
		id := s.idcount
		s.sockets = append(s.sockets, socket{id: id, ws: ws})
		gko.LogDebugf("socket(%d) nuevo", id)
		s.mu.Unlock()
		for {
			var msg string
			err := websocket.Message.Receive(ws, &msg)
			if err != nil {
				if err.Error() != "EOF" {
					gko.Err(err).Msgf("socket(%d) receive", id).Log()
				}
				break
			}
			gko.LogDebugf("socket(%d) recived: %s", id, msg)
		}
		// gko.LogDebugf("socket(%d) terminado", id)
		for i, socket := range s.sockets { // remove closed connection
			if socket.ws == ws {
				// gko.LogDebugf("socket(%d) eliminado", id)
				s.sockets = append(s.sockets[:i], s.sockets[i+1:]...)
				break
			}
		}
	}).ServeHTTP(c.Response(), c.Request())
	return nil
}

func (s *reloader) brodcastReload(c *gecko.Context) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.counter++
	gko.LogDebugf("Broadcasting reload #%d to %d connections", s.counter, len(s.sockets))
	for _, ws := range s.sockets {
		err := websocket.Message.Send(ws.ws, "reload")
		if err != nil {
			gko.Err(err).Msgf("socket(%d) send", ws.id).Log()
			for i, socket := range s.sockets { // remove closed connection
				if socket == ws {
					s.sockets = append(s.sockets[:i], s.sockets[i+1:]...)
					gko.LogDebugf("socket(%d) eliminado2", socket.id)
					break
				}
			}
		}
	}
	return c.StatusOkf("Reload #%d sent to %d connections", s.counter, len(s.sockets))
}
