package sqlitearbol

import (
	"monorepo/arbol"

	"github.com/pargomx/gecko/gko"
)

const NODO_ROOT = 1

func (s *Repositorio) GetRaiz() (*arbol.Raiz, error) {
	op := gko.Op("GetRaiz")
	desc, err := s.ListDescendientes(NODO_ROOT)
	if err != nil {
		return nil, op.Err(err)
	}
	raiz := arbol.Raiz{
		Grupos:    desc.Grupos,
		Proyectos: desc.Proyectos,
	}

	for i := range raiz.Grupos {
		err := s.AddHijosToGrupo(&raiz.Grupos[i])
		if err != nil {
			return nil, op.Err(err)
		}
	}

	for i := range raiz.Proyectos {
		err := s.AddHijosToProyecto(&raiz.Proyectos[i])
		if err != nil {
			return nil, op.Err(err)
		}
	}
	return &raiz, nil
}

func (s *Repositorio) AddHijosToGrupo(raiz *arbol.Grupo) error {
	op := gko.Op("addHijosToGrupo").Ctx("GrupoID", raiz.GrupoID)
	desc, err := s.ListDescendientes(raiz.GrupoID)
	if err != nil {
		return op.Err(err)
	}
	raiz.Grupos = desc.Grupos
	raiz.Proyectos = desc.Proyectos

	for i := range raiz.Grupos {
		err := s.AddHijosToGrupo(&raiz.Grupos[i])
		if err != nil {
			return op.Err(err)
		}
	}

	for i := range raiz.Proyectos {
		err := s.AddHijosToProyecto(&raiz.Proyectos[i])
		if err != nil {
			return op.Err(err)
		}
	}

	return nil
}

func (s *Repositorio) AddHijosToProyecto(raiz *arbol.Proyecto) error {
	op := gko.Op("addHijosToProyecto").Ctx("ProyectoID", raiz.ProyectoID)
	desc, err := s.ListDescendientes(raiz.ProyectoID)
	if err != nil {
		return op.Err(err)
	}
	raiz.Personas = desc.Personas
	raiz.HisUsuario = desc.HisUsuario
	raiz.HisTecnicas = desc.HisTecnicas
	raiz.HisGestion = desc.HisGestion

	for i := range raiz.Personas {
		err := s.AddHijosToPersona(&raiz.Personas[i])
		if err != nil {
			return op.Err(err)
		}
	}

	return nil
}

func (s *Repositorio) AddHijosToPersona(raiz *arbol.Persona) error {
	op := gko.Op("addHijosToPersona").Ctx("PersonaID", raiz.PersonaID)
	desc, err := s.ListDescendientes(raiz.PersonaID)
	if err != nil {
		return op.Err(err)
	}
	raiz.Personas = desc.Personas
	raiz.Historias = desc.HisUsuario
	raiz.HisTecnicas = desc.HisTecnicas
	raiz.HisGestion = desc.HisGestion

	for i := range raiz.Personas {
		err := s.AddHijosToPersona(&raiz.Personas[i])
		if err != nil {
			return op.Err(err)
		}
	}

	return nil
}

func (s *Repositorio) AddAncestrosToHisUsuario(raiz *arbol.HistoriaDeUsuario) error {
	op := gko.Op("addAncestrosToHisUsuario").Ctx("HistoriaID", raiz.HistoriaID)
	nextPadreID := raiz.PadreID
	loops := 0
	for {
		loops++ // prevent runaway
		if loops > 10 {
			gko.LogWarnf("Historia %v tiene demasiados ancestros", raiz.HistoriaID)
			break
		}
		// end when proyect is reached
		if raiz.Proyecto.ProyectoID != 0 {
			break
		}
		nod, err := s.GetNodo(nextPadreID)
		if err != nil {
			return op.Err(err).Ctx("id", nextPadreID)
		}
		nextPadreID = nod.PadreID
		switch {
		case nod.EsHistoriaDeUsuario(),
			nod.EsHistoriaTecnica(),
			nod.EsActividadDeGestión():
			raiz.Ancestros = append([]arbol.Nodo{*nod}, raiz.Ancestros...) // prepend

		case nod.EsPersona():
			raiz.Ancestros = append([]arbol.Nodo{*nod}, raiz.Ancestros...) // prepend
			if raiz.Persona.PersonaID == 0 {
				raiz.Persona = nod.ToPersona()
			}

		case nod.EsProyecto():
			raiz.Proyecto = nod.ToProyecto()

		default:
			return op.Msgf("Nodo %v %v no debería ser ancestro de historia %v", nod.Tipo, nextPadreID, raiz.HistoriaID)
		}
	}
	return nil
}

func (s *Repositorio) AddHijosToHisUsuario(raiz *arbol.HistoriaDeUsuario) error {
	op := gko.Op("addHijosToHisUsuario").Ctx("HistoriaID", raiz.HistoriaID)
	desc, err := s.ListDescendientes(raiz.HistoriaID)
	if err != nil {
		return op.Err(err)
	}
	raiz.Personas = desc.Personas
	raiz.HisUsuario = desc.HisUsuario
	raiz.HisTecnicas = desc.HisTecnicas
	raiz.HisGestion = desc.HisGestion

	raiz.Tramos = desc.Tramos
	raiz.Tareas = desc.Tareas
	raiz.Reglas = desc.Reglas

	for i := range raiz.HisUsuario {
		err := s.AddHijosToHisUsuario(&raiz.HisUsuario[i])
		if err != nil {
			return op.Err(err)
		}
	}

	relacionadas, err := s.ListNodosRelacionados(raiz.HistoriaID)
	if err != nil {
		return op.Err(err)
	}
	for _, nod := range relacionadas {
		raiz.Relacionadas = append(raiz.Relacionadas, nod.ToHistoriaDeUsuario())
	}

	return nil
}

// ================================================================ //
// ================================================================ //

func (s *Repositorio) GetProyecto(proyectoID int) (*arbol.Proyecto, error) {
	op := gko.Op("GetProyecto").Ctx("proyectoID", proyectoID)
	nod, err := s.GetNodo(proyectoID)
	if err != nil {
		return nil, op.Err(err)
	}
	pry := nod.ToProyecto()
	err = s.AddHijosToProyecto(&pry)
	if err != nil {
		return nil, op.Err(err)
	}
	return &pry, nil
}

func (s *Repositorio) GetPersona(personaID int) (*arbol.Persona, error) {
	op := gko.Op("GetPersona").Ctx("personaID", personaID)
	nod, err := s.GetNodo(personaID)
	if err != nil {
		return nil, op.Err(err)
	}
	per := nod.ToPersona()
	err = s.AddHijosToPersona(&per)
	if err != nil {
		return nil, op.Err(err)
	}
	return &per, nil
}

func (s *Repositorio) GetHistoria(historiaID int) (*arbol.HistoriaDeUsuario, error) {
	op := gko.Op("GetHistoria").Ctx("historiaID", historiaID)
	nodH, err := s.GetNodo(historiaID)
	if err != nil {
		return nil, op.Err(err)
	}
	his := nodH.ToHistoriaDeUsuario()
	err = s.AddHijosToHisUsuario(&his)
	if err != nil {
		return nil, op.Err(err)
	}
	err = s.AddAncestrosToHisUsuario(&his)
	if err != nil {
		return nil, op.Err(err)
	}
	return &his, nil
}

// ================================================================ //
// ================================================================ //

// Para obtener del repositorio todos los descendientes de un nodo y hacer lo
// que sea al respecto. Solo debe utilizarce como DTO interno entre los nodos y
// las entidades de dominio recursivas. No se debería filtrar a la aplicación
// porque es un detalle de implementación del repositorio.
type descendientes struct {
	Grupos    []arbol.Grupo
	Proyectos []arbol.Proyecto

	Personas    []arbol.Persona
	HisUsuario  []arbol.HistoriaDeUsuario
	HisTecnicas []arbol.HistoriaTecnica
	HisGestion  []arbol.ActividadDeGestión

	Reglas []arbol.Regla
	Tareas []arbol.Tarea
	Tramos []arbol.Tramo
}

func (s *Repositorio) ListDescendientes(padreID int) (descendientes, error) {
	op := gko.Op("ListDescendientes").Ctx("padreID", padreID)
	desc := descendientes{}
	nodos, err := s.ListNodosByPadreID(padreID)
	if err != nil {
		return desc, op.Err(err)
	}
	for _, nod := range nodos {
		switch nod.Tipo {

		case "GRP":
			if padreID != NODO_ROOT {
				gko.LogAlertf("Nodo descendiente %v es GRP", nod.NodoID)
			}
			desc.Grupos = append(desc.Grupos, nod.ToGrupo())

		case "PRY":
			desc.Proyectos = append(desc.Proyectos, nod.ToProyecto())

		case "PER":
			desc.Personas = append(desc.Personas, nod.ToPersona())

		case "HIS":
			desc.HisUsuario = append(desc.HisUsuario, nod.ToHistoriaDeUsuario())

		case "TEC":
			desc.HisTecnicas = append(desc.HisTecnicas, nod.ToHistoriaTecnica())

		case "GES":
			desc.HisGestion = append(desc.HisGestion, nod.ToActividadDeGestión())

		case "REG":
			desc.Reglas = append(desc.Reglas, nod.ToRegla())

		case "TAR":
			tar := nod.ToTarea()
			tar.Intervalos, err = s.ListIntervalosByNodoID(tar.TareaID)
			if err != nil {
				return desc, op.Err(err)
			}
			desc.Tareas = append(desc.Tareas, tar)

		case "VIA":
			desc.Tramos = append(desc.Tramos, nod.ToTramo())

		case "ROOT":
			// Ignorar raíz padre de sí misma.

		default:
			gko.LogAlertf("Nodo descendiente %v tipo %v desconocido", nod.NodoID, nod.Tipo)
		}
	}
	return desc, nil
}

// ================================================================ //
// ================================================================ //

// La cambia de posición respecto a sus hermanos, o sea en el mismo padre.
func (s *Repositorio) ReordenarNodo(nod arbol.Nodo, newPosicion int) error {
	var op = gko.Op("ReordenarNodo")
	if nod.NodoID == 0 {
		return op.Str("nodo inválido")
	}
	oldPosicion := nod.Posicion
	op.Ctx("id", nod.NodoID)

	if oldPosicion == newPosicion {
		return op.Msg("Posición sin cambios")
	}

	// Validar nueva posición
	var hermanos int
	err := s.db.QueryRow("SELECT COUNT(*) FROM nodos WHERE padre_id = ? AND tipo = ?",
		nod.PadreID, nod.Tipo).Scan(&hermanos)
	if err != nil {
		return op.Op("contar_hermanos").Err(err)
	}
	if newPosicion < 1 || newPosicion > hermanos {
		return op.Msgf("Posición %v inválida para nodo %v que tiene %v hermanos",
			newPosicion, nod.NodoID, hermanos)
	}

	// Se utilizan posiciones negativas temporales porque no se puede confiar en el orden
	// con el que el update modifica las filas para que mantenga el UNIQUE con la posición.
	_, err = s.db.Exec(
		"UPDATE nodos SET posicion = -(?) WHERE nodo_id = ?",
		newPosicion, nod.NodoID,
	)
	if err != nil {
		return op.Op("set_pos_negativa").Err(err)
	}

	// Dependiendo si se mueve hacia arriba o abajo, recorrer a los hermanos.
	if oldPosicion < newPosicion {
		_, err = s.db.Exec(
			"UPDATE nodos SET posicion = -(posicion - 1) WHERE padre_id = ? AND tipo = ? AND posicion > ? AND posicion <= ?",
			nod.PadreID, nod.Tipo, oldPosicion, newPosicion,
		)
	} else {
		_, err = s.db.Exec(
			"UPDATE nodos SET posicion = -(posicion + 1) WHERE padre_id = ? AND tipo = ? AND posicion >= ? AND posicion <= ?",
			nod.PadreID, nod.Tipo, newPosicion, oldPosicion,
		)
	}
	if err != nil {
		return op.Op("set_pos_hermanos").Err(err)
	}

	// Volver a positivo todas las posiciones.
	_, err = s.db.Exec(
		"UPDATE nodos SET posicion = -(posicion) WHERE padre_id = ? AND tipo = ? AND posicion < 0",
		nod.PadreID, nod.Tipo,
	)
	if err != nil {
		return op.Op("set_pos_positiva").Err(err)
	}
	return nil
}

// ================================================================ //
// ========== Mover =============================================== //

func (s *Repositorio) MoverNodo(nod arbol.Nodo, newPadreID int) error {
	op := gko.Op("MoverNodo")
	if nod.NodoID == 0 || newPadreID == 0 {
		return op.Str("nodo movido o nuevo padre indefinido")
	}

	// Poner como último hermano en el nuevo padre.
	_, err := s.db.Exec(
		"UPDATE nodos SET padre_id = ?, posicion = (SELECT count(nodo_id)+1 FROM nodos WHERE padre_id = ? AND tipo = ?) WHERE nodo_id = ?",
		newPadreID, newPadreID, nod.Tipo, nod.NodoID,
	)
	if err != nil {
		return op.Op("update_nodo").Err(err)
	}

	// Recorrer hermanos en el padre viejo.
	_, err = s.db.Exec(
		"UPDATE nodos SET posicion = posicion - 1 WHERE padre_id = ? AND tipo = ? AND posicion > ?",
		nod.PadreID, nod.Tipo, nod.Posicion,
	)
	if err != nil {
		return op.Err(err).Op("update_old_hermanos")
	}

	/*
		Ya no se usa esto porque la tabla nodos dejó de usar unique(nodo_id, tipo, posicion).
		// Se utilizan posiciones negativas temporales porque no se puede confiar en el orden
		// con el que el update modifica las filas para que mantenga el UNIQUE con la posición.
		_, err = s.db.Exec(
			"UPDATE nodos SET posicion = -(posicion - 1) WHERE padre_id = ? AND tipo = ? AND posicion > ?",
			nod.PadreID, nod.Tipo, nod.Posicion,
		)
		if err != nil {
			return op.Op("old_hermanos_negativo").Err(err)
		}
		_, err = s.db.Exec(
			"UPDATE nodos SET posicion = -(posicion) WHERE padre_id = ? AND tipo = ? AND posicion < 0",
			nod.PadreID, nod.Tipo,
		)
		if err != nil {
			return op.Op("old_hermanos_positivo").Err(err)
		}
	*/

	return nil
}

// ================================================================ //
// ========== DELETE ============================================== //

// Elimina todos los hijos (descendientes inmediatos) de un nodo.
func (s *Repositorio) DeleteHijos(NodoID int) error {
	op := gko.Op("DeleteHijos")
	if NodoID == 0 {
		return op.E(gko.ErrDatoIndef).Str("nodoID sin especificar")
	}
	err := s.ExisteNodo(NodoID)
	if err != nil {
		return op.Err(err)
	}
	_, err = s.db.Exec(
		"DELETE FROM nodos WHERE padre_id = ?",
		NodoID,
	)
	if err != nil {
		return op.E(gko.ErrAlEscribir).Err(err)
	}
	return nil
}

//  ================================================================  //
//  ========== LIST RELACIONADOS ===================================  //

func (s *Repositorio) ListNodosRelacionados(NodoID int) ([]arbol.Nodo, error) {
	const op string = "ListNodosRelacionados"
	rows, err := s.db.Query(
		"SELECT "+
			"n.nodo_id, n.padre_id, n.tipo, n.posicion, n.titulo, n.descripcion, n.objetivo, n.notas, n.color, n.imagen, n.prioridad, n.estatus, n.segundos, n.centavos "+
			" FROM nodos n JOIN referencias r ON r.nodo_id = n.nodo_id OR r.ref_nodo_id = n.nodo_id WHERE (r.nodo_id = ? OR r.ref_nodo_id = ?) AND n.nodo_id != ?",
		NodoID, NodoID, NodoID,
	)
	if err != nil {
		return nil, gko.ErrInesperado.Err(err).Op(op)
	}
	return s.scanRowsNodo(rows, op)
}
