package main

import (
	"fmt"
	"time"

	"github.com/pargomx/gecko"
)

func logDevReq(c *gecko.Context) bool {
	if AMBIENTE == "DEV" {
		htmx := "->"
		if c.EsHTMX() {
			htmx = "hx"
		}
		params := ""
		for k, v := range c.Request().URL.Query() {
			params += k + "=" + v[0] + " "
		}
		fmt.Println(
			"\033[32m"+htmx+"\033[0m",
			"\033[2m"+time.Now().Format("15:04:05.000")+"\033[0m",
			c.Path()+"\033[2m",
			c.Request().URL.String(),
			params,
			"\033[0m",
		)
	}
	return true
}

func (s *servidor) GET(path string, authHandler gecko.HandlerFunc) {
	s.gecko.GET(path, func(c *gecko.Context) error {
		logDevReq(c)
		return authHandler(c)
	})
}

func (s *servidor) POS(path string, authHandler gecko.HandlerFunc) {
	s.gecko.POST(path, func(c *gecko.Context) error {
		logDevReq(c)
		return authHandler(c)
	})
}

func (s *servidor) PCH(path string, authHandler gecko.HandlerFunc) {
	s.gecko.PATCH(path, func(c *gecko.Context) error {
		logDevReq(c)
		return authHandler(c)
	})
}

func (s *servidor) PUT(path string, authHandler gecko.HandlerFunc) {
	s.gecko.PUT(path, func(c *gecko.Context) error {
		logDevReq(c)
		return authHandler(c)
	})
}

func (s *servidor) DEL(path string, authHandler gecko.HandlerFunc) {
	s.gecko.DELETE(path, func(c *gecko.Context) error {
		logDevReq(c)
		return authHandler(c)
	})
}
