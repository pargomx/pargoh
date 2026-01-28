package main

import (
	"fmt"
	"time"

	"github.com/pargomx/gecko"
	"github.com/pargomx/gecko/gko"
)

// Si no se pasa ningún rol, entonces siempre es permitido.
// Si el usuario es SuperAdmin todo le es permitido, no es necesario pasarlo como argumento.
func getResponsablePermitido(c *gecko.Context) (permitido bool) {
	if c.Sesion == nil {
		return false
	}
	// ses, ok := c.Sesion.(*Sesion)
	// if !ok {
	// 	gko.LogWarnf("Invalid session type %T", c.Sesion)
	// 	return false
	// }
	return true
	// ResponsableID = ses.Persona.PersonaID
	// if ResponsableID == 0 {
	// 	gko.LogWarnf("Empty ResponsableID on session '%v'", ses.SesionID)
	// 	return  false
	// }
	// if ses.Persona.Rol.EsSuper() {
	// 	return  true
	// }
	// for _, rolPermitido := range roles {
	// 	if rolPermitido.EsTodos() {
	// 		return  true
	// 	}
	// 	if ses.Persona.Rol.Es(rolPermitido) {
	// 		return  true
	// 	}
	// }
	// return false
}

// ================================================================ //
// ========== Middleware ========================================== //

func (s *servidor) GET(path string, authHandler gecko.HandlerFunc) {
	s.gecko.GET(path, s.auth.Auth(func(c *gecko.Context) error {
		logDevReq(c)
		if !getResponsablePermitido(c) {
			return gko.ErrNoAutorizado.Msg("Acceso no permitido")
		}
		// c.Response().Header().Set("Cache-Control", "no-store")
		return authHandler(c)
	}))
}

func (s *servidor) POSNoTx(path string, authHandler gecko.HandlerFunc) {
	s.gecko.POST(path, s.auth.Auth(func(c *gecko.Context) error {
		logDevReq(c)
		if !getResponsablePermitido(c) {
			return gko.ErrNoAutorizado.Msg("Acceso no permitido")
		}
		// c.Response().Header().Set("Cache-Control", "no-store")
		return authHandler(c)
	}))
}

func (s *servidor) POS(path string, authHandler handlerTxFunc) {
	s.gecko.POST(path, s.auth.Auth(func(c *gecko.Context) error {
		logDevReq(c)
		if !getResponsablePermitido(c) {
			return gko.ErrNoAutorizado.Msg("Acceso no permitido")
		}
		if AMBIENTE == "DEV" {
			time.Sleep(time.Millisecond * 400)
		}
		hdlr := s.w.inTx(authHandler)
		return hdlr(c)
	}))
}

func (s *servidor) PCH(path string, authHandler handlerTxFunc) {
	s.gecko.PATCH(path, s.auth.Auth(func(c *gecko.Context) error {
		logDevReq(c)
		if !getResponsablePermitido(c) {
			return gko.ErrNoAutorizado.Msg("Acceso no permitido")
		}
		if AMBIENTE == "DEV" {
			time.Sleep(time.Millisecond * 400)
		}
		hdlr := s.w.inTx(authHandler)
		return hdlr(c)
	}))
}

func (s *servidor) PUT(path string, authHandler handlerTxFunc) {
	s.gecko.PUT(path, s.auth.Auth(func(c *gecko.Context) error {
		logDevReq(c)
		if !getResponsablePermitido(c) {
			return gko.ErrNoAutorizado.Msg("Acceso no permitido")
		}
		if AMBIENTE == "DEV" {
			time.Sleep(time.Millisecond * 400)
		}
		hdlr := s.w.inTx(authHandler)
		return hdlr(c)
	}))
}

func (s *servidor) DEL(path string, authHandler handlerTxFunc) {
	s.gecko.DELETE(path, s.auth.Auth(func(c *gecko.Context) error {
		logDevReq(c)
		if !getResponsablePermitido(c) {
			return gko.ErrNoAutorizado.Msg("Acceso no permitido")
		}
		if AMBIENTE == "DEV" {
			time.Sleep(time.Millisecond * 400)
		}
		hdlr := s.w.inTx(authHandler)
		return hdlr(c)
	}))
}

// ================================================================ //
// ========== LOG ================================================= //

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
