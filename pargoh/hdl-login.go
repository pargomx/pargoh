package main

import (
	"net/http"
	"time"

	gonanoid "github.com/matoous/go-nanoid/v2"
	"github.com/pargomx/gecko"
	"github.com/pargomx/gecko/gko"
)

type Sesion struct {
	SesionID  string
	Usuario   string
	IP        string
	UserAgent string
	ValidFrom time.Time
}

type authService struct {
	nombreCookie  string // Nombre del cookie de sesión.
	pathLoginPage string // default "/"
	pathLoginPost string // default "/login"
	pathHomePage  string
	pathLogout    string // default "/logout"

	vigencia time.Duration     // Vigencia de las sesiones.
	sesiones map[string]Sesion // Sesiones activas.
}

func NewAuthService() *authService {
	s := &authService{
		nombreCookie:  "pargotoken",
		pathLoginPage: "/",
		pathLoginPost: "/login",
		pathHomePage:  "/proyectos",
		pathLogout:    "/logout",
		sesiones:      make(map[string]Sesion),
		vigencia:      5 * 24 * time.Hour,
	}
	if s.pathLoginPage == "" {
		gko.LogWarn("No se ha definido la ruta para la página de inicio de sesión")
	}
	if AMBIENTE == "DEV" { // Sesión de prueba para no tener que loguearse en desarrollo.
		s.vigencia = 30 * 24 * time.Hour
		s.sesiones["DEV_MJv7kqY44Vnt2D8DbGBRBVqiaxN8bz70"] = Sesion{
			SesionID:  "DEV_MJv7kqY44Vnt2D8DbGBRBVqiaxN8bz70",
			Usuario:   "developer",
			IP:        "::1",
			UserAgent: "any",
			ValidFrom: time.Now(),
		}
	}
	return s
}

// Manda que el cliente elimine cookie de sesión de su navegador.
func (s *authService) limpiarCookie(c *gecko.Context) {
	c.SetCookie(&http.Cookie{
		Name:     s.nombreCookie,
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})
}

// Para no pedir autenticación basic en HAProxy.
func (s *authService) setCookieForHAProxy(c *gecko.Context) {
	c.SetCookie(&http.Cookie{
		Name:     "hapargo",
		Value:    "Z4v4!bsBM5BaeJ^Ryf6Pc*lfB",
		Path:     "/",
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Expires:  time.Now().Add(s.vigencia),
	})
}

// ================================================================ //
// ========== AUTENTICAR ========================================== //

func (s *authService) validarCredenciales(usuario, passwrd string) (string, error) {
	lenUsuario := len(usuario)
	lenPasswrd := len(passwrd)
	if lenUsuario > 25 || lenPasswrd > 25 {
		return "", gko.ErrDatoInvalido().Strf("creds_too_long: usuario(%d) passwd(%d)", lenUsuario, lenPasswrd)
	}
	if lenUsuario < 5 || lenPasswrd < 5 {
		return "", gko.ErrDatoInvalido().Strf("creds_too_short: usuario(%d) passwd(%d)", lenUsuario, lenPasswrd)
	}
	// TODO: Guardar credenciales en ambiente, archivo o base de datos con hash&salt.
	if usuario == "tulio" && passwrd == "flores99leetcode" {
		return usuario, nil
	}
	return "", gko.ErrDatoInvalido().Strf("creds_not_found: usuario[%s] passwd(%d)", usuario, lenPasswrd)
}

func (s *authService) registrarNuevaSesion(usuario, ip, userAgent string) (*Sesion, error) {
	sesionID, err := gonanoid.Generate("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ", 36)
	if err != nil {
		return nil, gko.ErrInesperado().Op("generarSesion").Err(err)
	}
	ses := Sesion{
		SesionID:  sesionID,
		Usuario:   usuario,
		IP:        ip,
		UserAgent: userAgent,
		ValidFrom: time.Now(),
	}
	if AMBIENTE == "DEV" {
		ses.SesionID = "DEV_MJv7kqY44Vnt2D8DbGBRBVqiaxN8bz70"
	}
	s.sesiones[ses.SesionID] = ses
	return &ses, nil
}

func (s *authService) validarSesion(sesionID string) (*Sesion, error) {
	if sesionID == "" {
		return nil, gko.ErrDatoIndef().Str("sesion_empty")
	}
	if len(sesionID) > 50 {
		return nil, gko.ErrDatoInvalido().Str("sesion_too_long")
	}
	ses, ok := s.sesiones[sesionID]
	if !ok {
		return nil, gko.ErrDatoInvalido().Strf("sesion_not_found: %s", sesionID)
	}
	if time.Since(ses.ValidFrom) > 5*24*time.Hour {
		delete(s.sesiones, sesionID)
		return nil, gko.ErrDatoInvalido().Strf("sesion_expired: %s", sesionID)
	}
	return &ses, nil // con el pointer se podría mutar la sesión original???
}

// ================================================================ //
// ================================================================ //

// Si la sesión es inválida manda limpiar el cookie.
func (s *authService) validarSesionCookie(c *gecko.Context) (*Sesion, error) {
	cookieSesion, err := c.Cookie(s.nombreCookie)
	if err != nil {
		return nil, gko.ErrDatoIndef().Str("sesion_cookie_missing")
	}
	if cookieSesion.Value == "" {
		s.limpiarCookie(c)
		return nil, gko.ErrDatoIndef().Str("sesion_cookie_empty")
	}
	ses, err := s.validarSesion(cookieSesion.Value)
	if err != nil {
		s.limpiarCookie(c)
		return nil, err
	}
	return ses, nil
}

// Valida la sesión del usuario y la pone en el contexto.
func (s *authService) Auth(next gecko.HandlerFunc) gecko.HandlerFunc {
	return func(c *gecko.Context) error {
		// Excepciones para evitar redirecciones infinitas.
		if c.Path() == "GET /{$}" {
			return next(c)
		} else if c.Path() == "POST /login" {
			return next(c)
		}
		ses, err := s.validarSesionCookie(c)
		if err != nil {
			gko.Err(err).Op("Auth").Log()
			return c.RedirFull(s.pathLoginPage) // ...sesión inválida
		}
		c.SesionID = ses.SesionID
		c.Sesion = ses
		return next(c)
	}
}

// ================================================================ //
// ========== HANDLERS ============================================ //

// Debe ser pública, sin pasar por el middleware de autenticación.
func (s *authService) getLogin(c *gecko.Context) error {
	_, err := s.validarSesionCookie(c)
	if err == nil {
		// gko.LogWarnf("Usuario %v ya tiene sesión %v", ses.Usuario, ses.SesionID)
		return c.RedirFull(s.pathHomePage) // ...ya tenía sesión
	}
	return c.Render(200, "app/login", nil)
}

func (s *authService) postLogin(c *gecko.Context) error {
	ses, err := s.validarSesionCookie(c)
	if err == nil {
		gko.LogWarnf("Usuario %v ya tenía sesión %v", ses.Usuario, ses.SesionID)
		return c.RedirFull(s.pathHomePage) // ...ya tenía sesión
	}
	usuario, err := s.validarCredenciales(c.FormVal("usuario"), c.FormValue("passwd"))
	if err != nil {
		gko.Err(err).Op("postLogin").Log()
		return c.RedirFull(s.pathLoginPage)
	}
	ses, err = s.registrarNuevaSesion(usuario, c.RealIP(), c.Request().UserAgent())
	if err != nil {
		gko.Err(err).Op("postLogin").Log()
		return c.RedirFull(s.pathLoginPage)
	}
	cookie := &http.Cookie{
		Name:     s.nombreCookie,
		Value:    ses.SesionID,
		Path:     "/",
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	}
	if c.FormBool("recordar") {
		cookie.Expires = time.Now().Add(s.vigencia)
		s.setCookieForHAProxy(c)
	}
	if AMBIENTE == "DEV" {
		cookie.Secure = false // para pruebas en red local
	}
	c.SetCookie(cookie)
	gko.LogInfof("Login '%s' (%s) %s [%s] recordar=%v", ses.SesionID, ses.Usuario, ses.ValidFrom.Format("2006-01-02 15:04:05"), ses.IP, !cookie.Expires.IsZero())
	return c.RedirFull(s.pathHomePage)
}

func (s *authService) logout(c *gecko.Context) error {
	ses, err := s.validarSesionCookie(c)
	if err != nil {
		gko.Err(err).Op("logout").Log()
		return c.RedirFull(s.pathLoginPage) // ...ya no tenía sesión válida
	}
	delete(s.sesiones, ses.SesionID)
	s.limpiarCookie(c)
	gko.LogInfof("Logout '%s' (%s) %s [%s]", ses.SesionID, ses.Usuario, ses.ValidFrom.Format("2006-01-02 15:04:05"), ses.IP)
	return c.RedirFull(s.pathLoginPage)
}

// ================================================================ //
// ========== MANTENIMIENTO ======================================= //

func (s *authService) printSesiones(c *gecko.Context) error {
	for _, ses := range s.sesiones {
		gko.LogInfof("Sesión '%s' (%s) %s [%s]", ses.SesionID, ses.Usuario, ses.ValidFrom.Format("2006-01-02 15:04:05"), ses.IP)
	}
	return c.StringOk("Nope")
}
