package main

import (
	"math"
	"monorepo/ust"

	"github.com/pargomx/gecko"
	"github.com/pargomx/gecko/gko"
)

type DiaReport struct {
	Fecha     string
	Segundos  int
	Proyectos map[string]ProyectoReport
}

type ProyectoReport struct {
	Proyecto  ust.Proyecto
	Segundos  int
	Historias map[int]HistoriaReport
}

type HistoriaReport struct {
	Historia ust.NodoHistoria
	Segundos int
	Tareas   map[int]TareaReport
}

type TareaReport struct {
	Tarea      ust.Tarea
	Segundos   int
	Intervalos []ust.IntervaloEnDia
}

// ================================================================ //

type Proyecto1 struct {
	ust.Proyecto
	Historias []ust.NodoHistoria
	Segundos  int
}

func (s *servidor) getMétricas1(c *gecko.Context) error {
	proyectos, err := s.repo.ListProyectos()
	if err != nil {
		return err
	}
	Proyectos := make([]Proyecto1, len(proyectos))
	for i, proyecto := range proyectos {
		historias, err := s.repo.ListNodoHistoriasByProyectoID(proyecto.ProyectoID)
		if err != nil {
			return err
		}
		segs := 0
		for _, historia := range historias {
			segs += historia.Segundos
		}
		Proyectos[i] = Proyecto1{
			Proyecto:  proyecto,
			Historias: historias,
			Segundos:  segs,
		}
	}
	data := map[string]any{
		"Titulo":    "Métricas",
		"Proyectos": Proyectos,
	}
	return c.RenderOk("metricas2", data)
}

func (s *servidor) getMétricas2(c *gecko.Context) error {

	// Traer todos los días trabajados hasta el presente.
	ListaDias, err := s.repo.ListDias()
	if err != nil {
		return err
	}
	// Traer todos los recursos para no estar buscando en DB y poner en mapa.
	Intervalos, err := s.repo.ListIntervalosEnDias()
	if err != nil {
		return err
	}
	IntervalosMapDia := make(map[string][]ust.IntervaloEnDia)
	for _, interv := range Intervalos {
		IntervalosMapDia[interv.Fecha] = append(IntervalosMapDia[interv.Fecha], interv)
	}

	Tareas, err := s.repo.ListTareas()
	if err != nil {
		return err
	}
	TareasMap := make(map[int]ust.Tarea, len(Tareas))
	for _, tarea := range Tareas {
		TareasMap[tarea.TareaID] = tarea
	}
	Historias, err := s.repo.ListNodoHistorias()
	if err != nil {
		return err
	}
	HistoriasMap := make(map[int]ust.NodoHistoria, len(Historias))
	for _, historia := range Historias {
		HistoriasMap[historia.HistoriaID] = historia
	}

	// Popular la estructura de días trabajados.
	Dias := make([]DiaReport, len(ListaDias))
	for i, dia := range ListaDias {

		Dias[i].Fecha = dia

		for _, itv := range IntervalosMapDia[dia] {

			// Conocer a quién pertenece.
			tarea, ok := TareasMap[itv.TareaID]
			if !ok {
				return gko.ErrNoEncontrado().Msgf("Tarea %d no encontrada", itv.TareaID)
			}
			historia, ok := HistoriasMap[tarea.HistoriaID]
			if !ok {
				return gko.ErrNoEncontrado().Msgf("Historia %d no encontrada", tarea.HistoriaID)
			}

			Dias[i].Segundos += itv.Segundos

			if Dias[i].Proyectos == nil {
				Dias[i].Proyectos = make(map[string]ProyectoReport)
			}
			pro, ok := Dias[i].Proyectos[historia.ProyectoID]
			if !ok {
				proyecto, err := s.repo.GetProyecto(historia.ProyectoID)
				if err != nil {
					return err
				}
				pro = ProyectoReport{
					Proyecto: *proyecto,
					Segundos: itv.Segundos,
				}
			} else {
				pro.Segundos += itv.Segundos
			}

			if pro.Historias == nil {
				pro.Historias = make(map[int]HistoriaReport)
			}
			his, ok := pro.Historias[historia.HistoriaID]
			if !ok {
				his = HistoriaReport{
					Historia: historia,
					Segundos: itv.Segundos,
				}
			} else {
				his.Segundos += itv.Segundos
			}

			if his.Tareas == nil {
				his.Tareas = make(map[int]TareaReport)
			}
			tar, ok := his.Tareas[itv.TareaID]
			if !ok {
				tar = TareaReport{
					Tarea:    tarea,
					Segundos: itv.Segundos,
				}
				tarea.TiempoReal = itv.Segundos // reset para solo este día
			} else {
				tar.Segundos += itv.Segundos
				tarea.TiempoReal += itv.Segundos
			}

			tar.Intervalos = append(tar.Intervalos, itv)
			his.Tareas[itv.TareaID] = tar
			pro.Historias[historia.HistoriaID] = his
			Dias[i].Proyectos[historia.ProyectoID] = pro
		}
	}

	Proyectos := map[string]ProyectoReport{}
	for _, dia := range Dias {
		for _, p := range dia.Proyectos {
			pry, ok := Proyectos[p.Proyecto.ProyectoID]
			if !ok {
				Proyectos[p.Proyecto.ProyectoID] = ProyectoReport{
					Proyecto: p.Proyecto,
					Segundos: p.Segundos,
				}
			} else {
				pry.Segundos += p.Segundos
				Proyectos[p.Proyecto.ProyectoID] = pry
			}
		}
	}

	// Viejo recuento de horas por día.
	DiasTrabajoMapHoras := make(map[string]float64)
	for _, dia := range Intervalos {
		DiasTrabajoMapHoras[dia.Fecha] += float64(dia.Segundos) / 60 / 60
	}

	data := map[string]any{
		"Titulo": "Métricas",

		"DiasTrabajoMapHoras": DiasTrabajoMapHoras,
		"DiasTrabajo":         Dias,
		"Proyectos":           Proyectos,
	}
	return c.RenderOk("metricas", data)
}

// ================================================================ //

func (d DiaReport) Horas() float64 {
	return math.Round(float64(d.Segundos)/3600*100) / 100
}
func (d ProyectoReport) Horas() float64 {
	return math.Round(float64(d.Segundos)/3600*100) / 100
}
func (d HistoriaReport) Horas() float64 {
	return math.Round(float64(d.Segundos)/3600*100) / 100
}
func (d TareaReport) Horas() float64 {
	return math.Round(float64(d.Segundos)/3600*100) / 100
}

// ================================================================ //

func (s *servidor) getMetricasProyecto(c *gecko.Context) error {
	Proyecto, err := s.repo.GetProyecto(c.PathVal("proyecto_id"))
	if err != nil {
		return err
	}
	Personas, err := s.repo.ListNodosPersonas(Proyecto.ProyectoID)
	if err != nil {
		return err
	}
	Proyectos, err := s.repo.ListProyectos()
	if err != nil {
		return err
	}
	TareasEnCurso, err := s.repo.ListTareasEnCurso()
	if err != nil {
		return err
	}

	Historias, err := s.repo.ListNodoHistoriasByProyectoID(Proyecto.ProyectoID)
	if err != nil {
		return err
	}

	data := map[string]any{
		"Titulo":        Proyecto.Titulo,
		"Proyecto":      Proyecto,
		"Personas":      Personas,
		"Proyectos":     Proyectos, // Para cambiar de proyecto a una persona.
		"TareasEnCurso": TareasEnCurso,
		"Historias":     Historias,
	}
	return c.RenderOk("proyecto", data)
}
