package main

import (
	"monorepo/ust"

	"github.com/pargomx/gecko"
)

type DiaTrabajo struct {
	Fecha     string
	Segundos  int
	Proyectos map[int]DiaTrabajoPorProyecto
}

type DiaTrabajoPorProyecto struct {
	Proyecto  ust.Proyecto
	Segundos  int
	Historias map[int]DiaTrabajoPorHistoria
}

type DiaTrabajoPorHistoria struct {
	Historia ust.Historia
	Segundos int
	Tareas   map[int]DiaTrabajoPorTarea
}

type DiaTrabajoPorTarea struct {
	Tarea      ust.Tarea
	Segundos   int
	Intervalos []ust.IntervaloEnDia
}

func (s *servidor) getMétricas(c *gecko.Context) error {
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
	Historias, err := s.repo.ListHistorias()
	if err != nil {
		return err
	}
	HistoriasMap := make(map[int]ust.Historia, len(Historias))
	for _, historia := range Historias {
		HistoriasMap[historia.HistoriaID] = historia
	}

	// Popular la estructura de días trabajados.
	Dias := make([]DiaTrabajo, len(ListaDias))
	for i, dia := range ListaDias {

		Dias[i].Fecha = dia

		for _, itv := range IntervalosMapDia[dia] {
			// Conocer a quién pertenece.

			// tarea, ok := TareasMap[itv.TareaID]
			// if !ok {
			// 	return gko.ErrNoEncontrado().Msgf("Tarea %d no encontrada", itv.TareaID)
			// }
			// historia, ok := HistoriasMap[tarea.HistoriaID]
			// if !ok {
			// 	return gko.ErrNoEncontrado().Msgf("Historia %d no encontrada", tarea.HistoriaID)
			// }

			Dias[i].Segundos += itv.Segundos()

			if Dias[i].Proyectos == nil {
				Dias[i].Proyectos = make(map[int]DiaTrabajoPorProyecto)
			}

			// if pro, ok := Dias[i].Proyectos[itv.ProyectoID]; !ok {
			// 	proyecto, err := s.repo.GetProyecto(itv.ProyectoID)
			// 	if err != nil {
			// 		return err
			// 	}
			// 	pro.Proyecto = *proyecto
			// 	pro.Segundos = itv.Segundos()
			// 	Dias[i].Proyectos[itv.ProyectoID] = pro
			// } else {
			// 	pro.Segundos += itv.Segundos()
			// 	Dias[i].Proyectos[itv.ProyectoID] = pro
			// }

			// if tar, ok := Dias[i].Tareas[itv.TareaID]; !ok {
			// 	tarea, err := s.repo.GetTarea(itv.TareaID)
			// 	if err != nil {
			// 		return err
			// 	}
			// 	tarea.TiempoReal = itv.Segundos()
			// 	Dias[i].Tareas[itv.TareaID] = *tarea
			// } else {
			// 	tar.TiempoReal += itv.Segundos()
			// 	Dias[i].Tareas[itv.TareaID] = tar
			// }
		}
	}

	// Viejo recuento de horas por día.
	DiasTrabajoMapHoras := make(map[string]float64)
	for _, dia := range Intervalos {
		DiasTrabajoMapHoras[dia.Fecha] += float64(dia.Segundos()) / 60 / 60
	}

	data := map[string]any{
		"Titulo": "Métricas",

		"DiasTrabajoMapHoras": DiasTrabajoMapHoras,
		"DiasTrabajo":         Dias,
	}
	return c.RenderOk("metricas", data)
}

// ================================================================ //

func (d DiaTrabajo) Horas() float64 {
	return float64(d.Segundos) / 60 / 60
}

// ================================================================ //
