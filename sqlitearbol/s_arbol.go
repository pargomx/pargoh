package sqlitearbol

import (
	"monorepo/arbol"

	"github.com/pargomx/gecko/gko"
)

func (s *Repositorio) GetProyecto(proyectoID int) (*arbol.Proyecto, error) {
	op := gko.Op("GetProyecto").Ctx("proyectoID", proyectoID)
	nod, err := s.GetNodo(proyectoID)
	if err != nil {
		return nil, op.Err(err)
	}
	pry := nod.ToProyecto()

	desc, err := s.listDescendientes(pry.ProyectoID)
	if err != nil {
		return nil, op.Err(err)
	}
	pry.Personas = desc.Personas
	pry.HisUsuario = desc.HisUsuario
	pry.HisTecnicas = desc.HisTecnicas
	pry.HisGestion = desc.HisGestion

	return &pry, nil
}

func (s *Repositorio) GetPersona(personaID int) (*arbol.Persona, error) {
	op := gko.Op("GetPersona").Ctx("personaID", personaID)
	nod, err := s.GetNodo(personaID)
	if err != nil {
		return nil, op.Err(err)
	}
	per := nod.ToPersona()

	desc, err := s.listDescendientes(per.PersonaID)
	if err != nil {
		return nil, op.Err(err)
	}
	per.Historias = desc.HisUsuario
	per.HisTecnicas = desc.HisTecnicas
	per.HisGestion = desc.HisGestion

	return &per, nil
}

func (s *Repositorio) GetHistoria(historiaID int) (*arbol.HistoriaDeUsuario, error) {
	op := gko.Op("GetHistoria").Ctx("historiaID", historiaID)
	nodH, err := s.GetNodo(historiaID)
	if err != nil {
		return nil, op.Err(err)
	}
	his := nodH.ToHistoriaDeUsuario()

	desc, err := s.listDescendientes(his.HistoriaID)
	if err != nil {
		return nil, op.Err(err)
	}
	his.HisUsuario = desc.HisUsuario
	his.HisTecnicas = desc.HisTecnicas
	his.HisGestion = desc.HisGestion

	his.Tramos = desc.Tramos
	his.Tareas = desc.Tareas
	his.Reglas = desc.Reglas

	nextPadreID := his.PadreID
	loops := 0
	for {
		loops++ // prevent runaway
		if loops > 10 {
			gko.LogWarnf("Historia %v tiene demasiados ancestros", his.HistoriaID)
			break
		}
		// end when proyect is reached
		if his.Proyecto.ProyectoID != 0 {
			break
		}
		nod, err := s.GetNodo(nextPadreID)
		if err != nil {
			return nil, op.Err(err).Ctx("id", nextPadreID)
		}
		nextPadreID = nod.PadreID
		switch {
		case nod.EsHistoriaDeUsuario(),
			nod.EsHistoriaTecnica(),
			nod.EsActividadDeGestión():
			his.Ancestros = append([]arbol.Nodo{*nod}, his.Ancestros...) // prepend

		case nod.EsPersona():
			his.Ancestros = append([]arbol.Nodo{*nod}, his.Ancestros...) // prepend
			if his.Persona.PersonaID == 0 {
				his.Persona = nod.ToPersona()
			}

		case nod.EsProyecto():
			his.Proyecto = nod.ToProyecto()

		default:
			return nil, op.Msgf("Nodo %v %v no debería ser ancestro de historia %v", nod.Tipo, nextPadreID, his.HistoriaID)
		}
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
	Personas    []arbol.Persona
	HisUsuario  []arbol.HistoriaDeUsuario
	HisTecnicas []arbol.HistoriaTecnica
	HisGestion  []arbol.ActividadDeGestión

	Tramos []arbol.Tramo
	Tareas []arbol.Tarea
	Reglas []arbol.Regla
}

func (s *Repositorio) listDescendientes(padreID int) (descendientes, error) {
	op := gko.Op("listDescendientes").Ctx("padreID", padreID)
	desc := descendientes{}
	nodos, err := s.ListNodosByPadreID(padreID)
	if err != nil {
		return desc, op.Err(err)
	}
	for _, nod := range nodos {
		switch nod.Tipo {
		case "GRP":
			gko.LogAlertf("Nodo descendiente %v es GRP", nod.NodoID)
		case "PRY":
			gko.LogAlertf("Nodo descendiente %v es PRY", nod.NodoID)

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
			desc.Tareas = append(desc.Tareas, nod.ToTarea())
		case "VIA":
			desc.Tramos = append(desc.Tramos, nod.ToTramo())

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
