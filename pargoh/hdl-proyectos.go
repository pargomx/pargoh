package main

import (
	"monorepo/arbol"
	"monorepo/dhistorias"
	"monorepo/ust"
	"strings"

	"github.com/pargomx/gecko"
	"github.com/pargomx/gecko/gko"
)

func (s *readhdl) getProyectosActivos(c *gecko.Context) error {
	nodo, err := s.repo.GetNodo(arbol.NODO_PROYECTOS_ACTIVOS)
	if err != nil {
		return err
	}
	raiz := nodo.ToGrupo()
	err = s.repo.AddHijosToGrupo(&raiz)
	if err != nil {
		return err
	}
	data := map[string]any{
		"Titulo": raiz.Nombre,
		"Grupo":  raiz,
	}
	return c.RenderOk("grupo", data)
}

func (s *readhdl) getNodoCualquiera(c *gecko.Context) error {

	data := map[string]any{
		"Titulo": "Pargo",
	}

	nodo, err := s.repo.GetNodo(c.PathInt("nodo_id"))
	if err != nil {
		return err
	}
	switch nodo.Tipo {

	case arbol.TipoGrupo:
		raiz := nodo.ToGrupo()
		err = s.repo.AddHijosToGrupo(&raiz)
		if err != nil {
			return err
		}
		data["Titulo"] = raiz.Nombre
		data["Grupo"] = raiz
		return c.RenderOk("grupo", data)

	case arbol.TipoProyecto:
		raiz := nodo.ToProyecto()
		err = s.repo.AddHijosToProyecto(&raiz)
		if err != nil {
			return err
		}
		data["Titulo"] = raiz.Titulo
		data["Proyecto"] = raiz
		return c.RenderOk("proyecto", data)

	case arbol.TipoPersona:
		raiz := nodo.ToPersona()
		err = s.repo.AddHijosToPersona(&raiz)
		if err != nil {
			return err
		}
		data["Titulo"] = raiz.Nombre
		data["Persona"] = raiz
		return c.RenderOk("persona", data)

	case arbol.TipoHistoria:
		raiz := nodo.ToHistoriaDeUsuario()
		err = s.repo.AddHijosToHisUsuario(&raiz)
		if err != nil {
			return err
		}
		err = s.repo.AddAncestrosToHisUsuario(&raiz)
		if err != nil {
			return err
		}
		data["Titulo"] = raiz.Titulo
		data["Historia"] = raiz
		return c.RenderOk("historia", data)

	case "ROOT":
		// Ignorar raíz padre de sí misma.

	default:
		data["Titulo"] = nodo.Titulo
		data["Nodo"] = nodo
		return c.RenderOk("nodo-raw", data)

	}
	return gko.ErrNoSoportado.Msg("Nodo inválido")
}

func (s *readhdl) listaProyectosOld(c *gecko.Context) error {
	type Pry struct {
		ust.Proyecto
		Personas []ust.NodoPersona
	}
	Proyectos, err := s.repoOld.ListProyectos()
	if err != nil {
		return err
	}

	res := make([]Pry, len(Proyectos))
	for i, p := range Proyectos {
		res[i].Proyecto = p
		res[i].Personas, err = s.repoOld.ListNodosPersonas(p.ProyectoID)
		if err != nil {
			return err
		}
	}
	Bugs, err := s.repoOld.ListTareasBugs()
	if err != nil {
		return err
	}
	TareasEnCurso, err := s.repoOld.ListTareasEnCurso()
	if err != nil {
		return err
	}
	data := map[string]any{
		"Titulo":        "Pargo",
		"Proyectos":     res,
		"Bugs":          Bugs,
		"TareasEnCurso": TareasEnCurso,
	}
	return c.RenderOk("proyectos", data)
}

func (s *servidor) setImagenProyecto(c *gecko.Context, tx *handlerTx) error {
	hdr, err := c.FormFile("imagen")
	if err == nil {
		file, err := hdr.Open()
		if err != nil {
			return err
		}
		defer file.Close()
		gko.LogDebugf("Imagen recibida: %v\t Tamaño: %v\t MIME:%v", hdr.Filename, hdr.Size, hdr.Header.Get("Content-Type"))
		err = dhistorias.SetImagenProyecto(c.PathVal("proyecto_id"), strings.TrimPrefix(hdr.Header.Get("Content-Type"), "image/"), file, s.cfg.ImagesDir, s.repoOld)
		if err != nil {
			return err
		}
	}
	return c.AskedFor("Proyecto actualizado")
}

func (s *writehdl) patchProyecto(c *gecko.Context, tx *handlerTx) error {
	err := tx.app.ParcharNodo(arbol.ArgsParcharNodo{
		NodoID: c.PathInt("proyecto_id"),
		Campo:  c.PathVal("param"),
		NewVal: c.FormValue("value"),
	})
	if err != nil {
		return err
	}
	return c.RefreshHTMX()
}

func (s *writehdl) postAppTime(c *gecko.Context, tx *handlerTx) error {
	err := tx.app.AddTimeSpent(c.PathInt("nodo_id"), c.PathInt("seg"))
	if err != nil {
		return err
	}
	return c.StringOk("ok")
}

func (s *readhdl) getProyecto(c *gecko.Context) error {
	Pry, err := s.repo.GetProyecto(c.PathInt("proyecto_id"))
	if err != nil {
		return err
	}
	TareasEnCurso, err := s.repoOld.ListTareasEnCurso()
	if err != nil {
		return err
	}
	data := map[string]any{
		"Titulo":   Pry.Titulo,
		"Proyecto": Pry,
		"Personas": Pry.Personas,
		// "Proyectos":     Proyectos, // Para cambiar de proyecto a una persona.
		"TareasEnCurso": TareasEnCurso,
	}
	return c.RenderOk("proyecto", data)
}

func (s *readhdl) getDocumentacionProyecto(c *gecko.Context) error {
	Proyecto, err := s.repoOld.GetProyecto(c.PathVal("proyecto_id"))
	if err != nil {
		return err
	}
	Personas, err := s.repoOld.ListNodosPersonas(Proyecto.ProyectoID)
	if err != nil {
		return err
	}
	type Personaje struct {
		Persona   ust.NodoPersona
		Historias []ust.Historia
	}
	Personajes := make([]Personaje, len(Personas))
	for i, p := range Personas {
		hists, err := s.repoOld.ListHistoriasByPadreID(p.PersonaID)
		if err != nil {
			return err
		}
		Personajes[i] = Personaje{
			Persona:   Personas[i],
			Historias: hists,
		}
	}
	data := map[string]any{
		"Titulo":     Proyecto.Titulo,
		"Proyecto":   Proyecto,
		"Personajes": Personajes,
	}
	return c.RenderOk("proyecto_doc", data)
}
