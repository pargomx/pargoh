package main

import (
	"monorepo/arbol"

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
		data["ScriptsHistoria"] = true
		return c.RenderOk("historia", data)

	case arbol.TipoTarea:
		raiz := nodo.ToTarea()
		intervalos, err := s.repo.ListIntervalosByNodoID(raiz.TareaID)
		if err != nil {
			return err
		}
		data["Titulo"] = raiz.Descripcion
		data["Tarea"] = raiz
		data["Intervalos"] = intervalos
		return c.RenderOk("tarea", data)

	case "ROOT":
		// Ignorar raíz padre de sí misma.

	default:
		data["Titulo"] = nodo.Titulo
		data["Nodo"] = nodo
		return c.RenderOk("nodo-raw", data)

	}
	return gko.ErrNoSoportado.Msg("Nodo inválido")
}

func (s *writehdl) postAppTime(c *gecko.Context, tx *handlerTx) error {
	err := tx.app.AddTimeSpent(c.PathInt("nodo_id"), c.PathInt("seg"))
	if err != nil {
		return err
	}
	return c.StringOk("ok")
}

func (s *readhdl) getDocumentacion(c *gecko.Context) error {
	nodo, err := s.repo.GetNodo(c.PathInt("nodo_id"))
	if err != nil {
		return err
	}
	data := map[string]any{
		"Titulo": nodo.Titulo,
		"Nodo":   nodo,
	}
	if nodo.EsPersona() {
		per := nodo.ToPersona()
		err = s.repo.AddHijosToPersona(&per)
		if err != nil {
			return err
		}
		data["Nodo"] = per
		return c.Render(200, "persona_doc", data)
	} else {
		pry := nodo.ToProyecto()
		err = s.repo.AddHijosToProyecto(&pry)
		if err != nil {
			return err
		}
		data["Nodo"] = pry
		return c.RenderOk("proyecto_doc", data)
	}
}
