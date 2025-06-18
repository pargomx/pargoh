package sqlitepuente

import (
	"fmt"
	"monorepo/arbol"
	"monorepo/ust"

	"github.com/pargomx/gecko/gko"
	"github.com/pargomx/gecko/gkt"
)

// ================================================================ //
// ========== Traductor del nuevo repo al viejo =================== //

var ErrMientrasMigramos = gko.ErrNoDisponible.Msg("No disponible durante migración")

func strToInt(txt string) int {
	num, _ := gkt.ToInt(txt)
	return num
}

// Proyectos
func (s *Repositorio) InsertProyecto(pro ust.Proyecto) error {
	return ErrMientrasMigramos.Copy().Op("InsertProyecto")
}

func (s *Repositorio) GetProyecto(ProyectoID string) (*ust.Proyecto, error) {
	nod, err := s.nvo.GetNodo(strToInt(ProyectoID))
	if err != nil {
		return nil, err
	}
	return &ust.Proyecto{
		ProyectoID:  fmt.Sprint(nod.NodoID),
		Posicion:    nod.Posicion,
		Titulo:      nod.Titulo,
		Color:       nod.Color,
		Imagen:      nod.Imagen,
		Descripcion: nod.Descripcion,
	}, nil
}

func (s *Repositorio) UpdateProyecto(ProyectoID string, pro ust.Proyecto) error {
	return ErrMientrasMigramos.Copy().Op("UpdateProyecto")
}

func (s *Repositorio) ExisteProyecto(ProyectoID string) error {
	return s.nvo.ExisteNodo(strToInt(ProyectoID))
}

func (s *Repositorio) DeleteProyecto(ProyectoID string) error {
	return s.nvo.DeleteNodo(strToInt(ProyectoID))
}

func (s *Repositorio) ListProyectos() ([]ust.Proyecto, error) {
	nodos, err := s.nvo.ListNodosByPadreIDTipo(2, "PRY")
	if err != nil {
		return nil, err
	}
	lista := []ust.Proyecto{}
	for _, nod := range nodos {
		lista = append(lista, ust.Proyecto{
			ProyectoID:  fmt.Sprint(nod.NodoID),
			Posicion:    nod.Posicion,
			Titulo:      nod.Titulo,
			Color:       nod.Color,
			Imagen:      nod.Imagen,
			Descripcion: nod.Descripcion,
		})
	}
	return lista, nil
}

// Nodos
func (s *Repositorio) InsertNodo(nod ust.Nodo) error {
	return ErrMientrasMigramos.Copy().Op("InsertNodo")
}

func (s *Repositorio) EliminarNodo(nodoID int) error {
	return s.nvo.DeleteNodo(nodoID)
}

func (s *Repositorio) MoverNodo(nodoID int, nuevoPadreID int) error {
	return ErrMientrasMigramos.Copy().Op("MoverNodo")
}

func (s *Repositorio) ReordenarNodo(nodoID int, oldPosicion int, newPosicion int) error {
	return ErrMientrasMigramos.Copy().Op("ReordenarNodo")
}

func (s *Repositorio) GetNodo(nodoID int) (*ust.Nodo, error) {
	nod, err := s.nvo.GetNodo(nodoID)
	if err != nil {
		return nil, err
	}
	return &ust.Nodo{
		NodoID:   nod.NodoID,
		NodoTbl:  nod.Tipo,
		PadreID:  nod.PadreID,
		PadreTbl: "???",
		Nivel:    2,
		Posicion: nod.Posicion,
	}, nil
}

func (s *Repositorio) ListNodosByPadreID(PadreID int) ([]ust.Nodo, error) {
	nodos, err := s.nvo.ListNodosByPadreIDTipo(PadreID, "HIS")
	if err != nil {
		return nil, err
	}
	lista := []ust.Nodo{}
	for _, nod := range nodos {
		lista = append(lista, ust.Nodo{
			NodoID:   nod.NodoID,
			NodoTbl:  nod.Tipo,
			PadreID:  nod.PadreID,
			PadreTbl: "???",
			Nivel:    2,
			Posicion: nod.Posicion,
		})
	}
	return lista, nil
}

// Personas
func (s *Repositorio) InsertPersona(per ust.Persona) error {
	return ErrMientrasMigramos.Copy().Op("InsertPersona")
}

func (s *Repositorio) GetPersona(personaID int) (*ust.Persona, error) {
	nod, err := s.nvo.GetNodo(personaID)
	if err != nil {
		return nil, err
	}
	return &ust.Persona{
		PersonaID:   nod.NodoID,
		ProyectoID:  fmt.Sprint(nod.PadreID),
		Nombre:      nod.Titulo,
		Descripcion: nod.Descripcion,
	}, nil
}

func (s *Repositorio) UpdatePersona(per ust.Persona) error {
	return ErrMientrasMigramos.Copy().Op("UpdatePersona")
}

func (s *Repositorio) DeletePersona(personaID int) error {
	return s.nvo.DeleteNodo(personaID)
}

func (s *Repositorio) ListNodosPersonas(ProyectoID string) ([]ust.NodoPersona, error) {
	nodos, err := s.nvo.ListNodosByPadreIDTipo(strToInt(ProyectoID), "PER")
	if err != nil {
		return nil, err
	}
	lista := []ust.NodoPersona{}
	for _, nod := range nodos {
		lista = append(lista, ust.NodoPersona{
			PersonaID:  nod.NodoID,
			ProyectoID: fmt.Sprint(nod.PadreID),
			Posicion:   nod.Posicion,

			Nombre:      nod.Titulo,
			Descripcion: nod.Descripcion,

			PadreID:  nod.PadreID,
			PadreTbl: "PRY",
			Nivel:    2,
		})
	}
	return lista, nil
}

// Historias
func (s *Repositorio) ExisteHistoria(HistoriaID int) error {
	return s.nvo.ExisteNodo(HistoriaID)
}

func (s *Repositorio) InsertHistoria(his ust.Historia) error {
	return ErrMientrasMigramos.Copy().Op("InsertHistoria")
}

func (s *Repositorio) UpdateHistoria(ust.Historia) error {
	return ErrMientrasMigramos.Copy().Op("UpdateHistoria")
}

func (s *Repositorio) DeleteHistoria(historiaID int) error {
	return s.nvo.DeleteNodo(historiaID)
}

func (s *Repositorio) GetHistoria(historiaID int) (*ust.Historia, error) {
	nod, err := s.nvo.GetNodo(historiaID)
	if err != nil {
		return nil, err
	}
	return &ust.Historia{
		HistoriaID:          nod.NodoID,
		Titulo:              nod.Titulo,
		Objetivo:            nod.Objetivo,
		Prioridad:           nod.Prioridad,
		Completada:          nod.Estatus > 0,
		SegundosPresupuesto: nod.Segundos,
		Descripcion:         nod.Descripcion,
		Notas:               nod.Notas,
	}, nil
}

func convertNodoHistoria(nod arbol.Nodo) ust.NodoHistoria {
	return ust.NodoHistoria{
		HistoriaID:          nod.NodoID,
		ProyectoID:          "",
		PersonaID:           0,
		Titulo:              nod.Titulo,
		Objetivo:            nod.Objetivo,
		Prioridad:           nod.Prioridad,
		Completada:          nod.Estatus > 0,
		SegundosPresupuesto: nod.Segundos,
		Descripcion:         nod.Descripcion,
		Notas:               nod.Notas,
		PadreID:             nod.PadreID,
		PadreTbl:            "HIS", // mentira
		Nivel:               0,
		Posicion:            nod.Posicion,
	}
}

func (s *Repositorio) GetNodoHistoria(nodoID int) (*ust.NodoHistoria, error) {
	nod, err := s.nvo.GetNodo(nodoID)
	if err != nil {
		return nil, err
	}
	his := convertNodoHistoria(*nod)
	return &his, nil
}

func (s *Repositorio) ListHistorias() ([]ust.Historia, error) {
	return nil, ErrMientrasMigramos.Copy().Op("ListHistorias")
}

func (s *Repositorio) ListNodoHistorias() ([]ust.NodoHistoria, error) {
	return nil, ErrMientrasMigramos.Copy().Op("ListNodoHistorias")
}

func (s *Repositorio) ListNodoHistoriasByPadreID(PadreID int) ([]ust.NodoHistoria, error) {
	nodos, err := s.nvo.ListNodosByPadreIDTipo(PadreID, "HIS")
	if err != nil {
		return nil, err
	}
	lista := []ust.NodoHistoria{}
	for _, nod := range nodos {
		lista = append(lista, convertNodoHistoria(nod))
	}
	return lista, nil
}

// Materializados
func (s *Repositorio) CambiarProyectoDeHistoriasByPersonaID(personaID int, proyectoID string) error {
	return ErrMientrasMigramos.Copy().Op("CambiarProyectoDeHistoriasByPersonaID")
}

// Tareas
func (s *Repositorio) InsertTarea(tar ust.Tarea) error {
	return ErrMientrasMigramos.Copy().Op("InsertTarea")
}

func (s *Repositorio) UpdateTarea(tar ust.Tarea) error {
	return ErrMientrasMigramos.Copy().Op("UpdateTarea")
}

func (s *Repositorio) DeleteTarea(tareaID int) error {
	return s.nvo.DeleteNodo(tareaID)
}

func (s *Repositorio) DeleteAllTareas(HistoriaID int) error {
	return ErrMientrasMigramos.Copy().Op("DeleteAllTareas")
}

func (s *Repositorio) GetTarea(tareaID int) (*ust.Tarea, error) {
	return nil, ErrMientrasMigramos.Copy().Op("GetTarea")
}

func (s *Repositorio) ListTareas() ([]ust.Tarea, error) {
	return nil, ErrMientrasMigramos.Copy().Op("ListTareas")
}

func (s *Repositorio) ListTareasByHistoriaID(historiaID int) ([]ust.Tarea, error) {
	nodos, err := s.nvo.ListNodosByPadreIDTipo(historiaID, "TAR")
	if err != nil {
		return nil, err
	}
	lista := []ust.Tarea{}
	for _, nod := range nodos {
		lista = append(lista, ust.Tarea{
			TareaID:          nod.NodoID,
			HistoriaID:       nod.PadreID,
			Descripcion:      nod.Titulo,
			Importancia:      ust.ImportanciaTareaNecesaria, // mentira
			Tipo:             ust.TipoTareaIndefinido,       // mentira
			Estatus:          nod.Estatus,
			Impedimentos:     nod.Objetivo,
			SegundosEstimado: nod.Segundos,
		})
	}
	return lista, nil
}

// Intervalos
func (s *Repositorio) InsertIntervalo(interv ust.Intervalo) error {
	return ErrMientrasMigramos.Copy().Op("InsertIntervalo")
}

func (s *Repositorio) UpdateIntervalo(TareaID int, Inicio string, interv ust.Intervalo) error {
	return ErrMientrasMigramos.Copy().Op("UpdateIntervalo")
}

func (s *Repositorio) DeleteIntervalo(TareaID int, Inicio string) error {
	return ErrMientrasMigramos.Copy().Op("DeleteIntervalo")
}

func (s *Repositorio) GetIntervalo(TareaID int, Inicio string) (*ust.Intervalo, error) {
	return nil, ErrMientrasMigramos.Copy().Op("GetIntervalo")
}

func (s *Repositorio) ListIntervalosByTareaID(TareaID int) ([]ust.Intervalo, error) {
	return []ust.Intervalo{}, nil
}

// Viajes
func (s *Repositorio) InsertTramo(tra ust.Tramo) error {
	return ErrMientrasMigramos.Copy().Op("InsertTramo")
}

func (s *Repositorio) UpdateTramo(tra ust.Tramo) error {
	return ErrMientrasMigramos.Copy().Op("UpdateTramo")
}

func (s *Repositorio) ExisteTramo(HistoriaID int, Posicion int) error {
	return ErrMientrasMigramos.Copy().Op("ExisteTramo")
}

func (s *Repositorio) DeleteTramo(HistoriaID int, Posicion int) error {
	return ErrMientrasMigramos.Copy().Op("DeleteTramo")
}

func (s *Repositorio) DeleteAllTramos(HistoriaID int) error {
	return ErrMientrasMigramos.Copy().Op("DeleteAllTramos")
}

func (s *Repositorio) GetTramo(HistoriaID int, Posicion int) (*ust.Tramo, error) {
	return nil, ErrMientrasMigramos.Copy().Op("GetTramo")
}

func (s *Repositorio) ListTramosByHistoriaID(HistoriaID int) ([]ust.Tramo, error) {
	nodos, err := s.nvo.ListNodosByPadreIDTipo(HistoriaID, "VIA")
	if err != nil {
		return nil, err
	}
	lista := []ust.Tramo{}
	for _, nod := range nodos {
		lista = append(lista, ust.Tramo{
			HistoriaID: nod.PadreID,
			Posicion:   nod.Posicion,
			Texto:      nod.Titulo,
			Imagen:     nod.Imagen,
		})
	}
	return lista, nil
}

func (s *Repositorio) MoverTramo(historiaID int, posicion int, newHistoriaID int) error {
	return ErrMientrasMigramos.Copy().Op("MoverTramo")
}

// Reglas
func (s *Repositorio) InsertRegla(reg ust.Regla) error {
	return ErrMientrasMigramos.Copy().Op("InsertRegla")
}

func (s *Repositorio) UpdateRegla(reg ust.Regla) error {
	return ErrMientrasMigramos.Copy().Op("UpdateRegla")
}

func (s *Repositorio) ExisteRegla(HistoriaID int, Posicion int) error {
	return ErrMientrasMigramos.Copy().Op("ExisteRegla")
}

func (s *Repositorio) DeleteRegla(HistoriaID int, Posicion int) error {
	return ErrMientrasMigramos.Copy().Op("DeleteRegla")
}

func (s *Repositorio) DeleteAllReglas(HistoriaID int) error {
	return ErrMientrasMigramos.Copy().Op("DeleteAllReglas")
}

func (s *Repositorio) GetRegla(HistoriaID int, Posicion int) (*ust.Regla, error) {
	return nil, ErrMientrasMigramos.Copy().Op("GetRegla")
}

func (s *Repositorio) ListReglasByHistoriaID(HistoriaID int) ([]ust.Regla, error) {
	nodos, err := s.nvo.ListNodosByPadreIDTipo(HistoriaID, "REG")
	if err != nil {
		return nil, err
	}
	lista := []ust.Regla{}
	for _, nod := range nodos {
		lista = append(lista, ust.Regla{
			HistoriaID:   nod.PadreID,
			Posicion:     nod.Posicion,
			Texto:        nod.Titulo,
			Implementada: nod.Estatus > 0,
			Probada:      nod.Estatus > 1,
		})
	}
	return lista, nil
}

func (s *Repositorio) ReordenarRegla(HistoriaID, oldPos, newPos int) error {
	return ErrMientrasMigramos.Copy().Op("ReordenarRegla")
}

// Referencias
func (s *Repositorio) InsertReferencia(ref ust.Referencia) error {
	return ErrMientrasMigramos.Copy().Op("InsertReferencia")
}

func (s *Repositorio) DeleteReferencia(HistoriaID int, RefHistoriaID int) error {
	return ErrMientrasMigramos.Copy().Op("DeleteReferencia")
}

func (s *Repositorio) ListNodoHistoriasRelacionadas(HistoriaID int) ([]ust.NodoHistoria, error) {
	return []ust.NodoHistoria{}, nil
}

// ================================================================ //
// ========== ADICIONALES READ ONLY =============================== //

func (s *Repositorio) FullTextSearch(search string) ([]ust.SearchResult, error) {
	return nil, ErrMientrasMigramos.Copy().Op("FullTextSearch")
}

func (s *Repositorio) ListNodoHistoriasByProyectoID(ProyectoID string) ([]ust.NodoHistoria, error) {
	return nil, ErrMientrasMigramos.Copy().Op("ListNodoHistoriasByProyectoID")
}

func (s *Repositorio) ListIntervalosEnDias() ([]ust.IntervaloEnDia, error) {
	return nil, ErrMientrasMigramos.Copy().Op("ListIntervalosEnDias")
}
func (s *Repositorio) ListIntervalosEnDiasByProyectoID(ProyectoID string) ([]ust.IntervaloEnDia, error) {
	return nil, ErrMientrasMigramos.Copy().Op("ListIntervalosEnDiasByProyectoID")
}
func (s *Repositorio) ListIntervalosEnDiasEntre(desde string, hasta string) ([]ust.IntervaloEnDia, error) {
	return nil, ErrMientrasMigramos.Copy().Op("ListIntervalosEnDiasEntre")
}

func (s *Repositorio) ListLatidos(desde, hasta string) ([]ust.Latido, error) {
	return nil, ErrMientrasMigramos.Copy().Op("ListLatidos")
}

func (s *Repositorio) ListHistoriasByPadreID(nodoID int) ([]ust.Historia, error) {
	nodos, err := s.nvo.ListNodosByPadreIDTipo(nodoID, "HIS")
	if err != nil {
		return nil, err
	}
	lista := []ust.Historia{}
	for _, nod := range nodos {
		lista = append(lista, ust.Historia{
			HistoriaID:          nod.NodoID,
			Titulo:              nod.Titulo,
			Objetivo:            nod.Objetivo,
			Prioridad:           nod.Prioridad,
			Completada:          nod.Estatus > 0,
			SegundosPresupuesto: nod.Segundos,
			Descripcion:         nod.Descripcion,
			Notas:               nod.Notas,
		})
	}
	return lista, nil
}

func (s *Repositorio) ListPersonas() ([]ust.Persona, error) {
	return nil, ErrMientrasMigramos.Copy().Op("ListPersonas")
}

func (s *Repositorio) ListTareasEnCurso() ([]ust.Tarea, error) {
	return []ust.Tarea{}, nil
}

func (s *Repositorio) ListTareasBugs() ([]ust.Tarea, error) {
	return []ust.Tarea{}, nil
}

func (s *Repositorio) ListHistoriasCosto(personaID int) ([]ust.HistoriaCosto, error) {
	return nil, ErrMientrasMigramos.Copy().Op("ListHistoriasCosto")
}

func (s *Repositorio) ListIntervalosRecientes() ([]ust.IntervaloReciente, error) {
	return nil, ErrMientrasMigramos.Copy().Op("ListIntervalosRecientes")
}
func (s *Repositorio) ListIntervalosRecientesAbiertos() ([]ust.IntervaloReciente, error) {
	return nil, ErrMientrasMigramos.Copy().Op("ListIntervalosRecientesAbiertos")
}

func (s *Repositorio) InsertLatido(lat ust.Latido) error {
	return ErrMientrasMigramos.Copy().Op("InsertLatido")
}

// ================================================================ //
// ========== DIAS ================================================ //

const qryDias = `WITH RECURSIVE date_range AS (
    SELECT (SELECT date(min(inicio), "-5 hours") FROM intervalos) AS dia
    UNION ALL
    SELECT date(dia, '+1 day')
    FROM date_range
    WHERE dia < date('now')
)
SELECT dia FROM date_range;`

func (s *Repositorio) ListDias() ([]string, error) {
	const op string = "ListDias"
	rows, err := s.db.Query(qryDias)
	if err != nil {
		return nil, gko.ErrInesperado.Err(err).Op(op)
	}
	defer rows.Close()
	items := []string{}
	for rows.Next() {
		var item string
		err := rows.Scan(&item)
		if err != nil {
			return nil, gko.ErrInesperado.Err(err).Op(op)
		}
		items = append(items, item)
	}
	return items, nil
}
