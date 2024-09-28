package main

import (
	"math"
	"monorepo/ust"

	"github.com/pargomx/gecko"
	"github.com/pargomx/gecko/gko"
)

type DiaTrabajo struct {
	Fecha     string
	Segundos  int
	Proyectos map[string]DiaTrabajoPorProyecto
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

type ProyectoTime struct {
	ProyectoID string
	Segundos   int
	Proyecto   ust.Proyecto
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
			tarea, ok := TareasMap[itv.TareaID]
			if !ok {
				return gko.ErrNoEncontrado().Msgf("Tarea %d no encontrada", itv.TareaID)
			}
			historia, ok := HistoriasMap[tarea.HistoriaID]
			if !ok {
				return gko.ErrNoEncontrado().Msgf("Historia %d no encontrada", tarea.HistoriaID)
			}

			Dias[i].Segundos += itv.Segundos()

			if Dias[i].Proyectos == nil {
				Dias[i].Proyectos = make(map[string]DiaTrabajoPorProyecto)
			}
			pro, ok := Dias[i].Proyectos[historia.ProyectoID]
			if !ok {
				proyecto, err := s.repo.GetProyecto(historia.ProyectoID)
				if err != nil {
					return err
				}
				pro = DiaTrabajoPorProyecto{
					Proyecto: *proyecto,
					Segundos: itv.Segundos(),
				}
			} else {
				pro.Segundos += itv.Segundos()
			}

			if pro.Historias == nil {
				pro.Historias = make(map[int]DiaTrabajoPorHistoria)
			}
			his, ok := pro.Historias[historia.HistoriaID]
			if !ok {
				his = DiaTrabajoPorHistoria{
					Historia: historia,
					Segundos: itv.Segundos(),
				}
			} else {
				his.Segundos += itv.Segundos()
			}

			if his.Tareas == nil {
				his.Tareas = make(map[int]DiaTrabajoPorTarea)
			}
			tar, ok := his.Tareas[itv.TareaID]
			if !ok {
				tar = DiaTrabajoPorTarea{
					Tarea:    tarea,
					Segundos: itv.Segundos(),
				}
				tarea.TiempoReal = itv.Segundos() // reset para solo este día
			} else {
				tar.Segundos += itv.Segundos()
				tarea.TiempoReal += itv.Segundos()
			}

			tar.Intervalos = append(tar.Intervalos, itv)
			his.Tareas[itv.TareaID] = tar
			pro.Historias[historia.HistoriaID] = his
			Dias[i].Proyectos[historia.ProyectoID] = pro
		}
	}

	Proyectos := map[string]ProyectoTime{}
	for _, dia := range Dias {
		for _, p := range dia.Proyectos {
			pry, ok := Proyectos[p.Proyecto.ProyectoID]
			if !ok {
				Proyectos[p.Proyecto.ProyectoID] = ProyectoTime{
					ProyectoID: p.Proyecto.ProyectoID,
					Proyecto:   p.Proyecto,
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
		DiasTrabajoMapHoras[dia.Fecha] += float64(dia.Segundos()) / 60 / 60
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

func (d DiaTrabajo) Horas() float64 {
	return math.Round(float64(d.Segundos)/3600*100) / 100
}
func (d DiaTrabajoPorProyecto) Horas() float64 {
	return math.Round(float64(d.Segundos)/3600*100) / 100
}
func (d DiaTrabajoPorHistoria) Horas() float64 {
	return math.Round(float64(d.Segundos)/3600*100) / 100
}
func (d DiaTrabajoPorTarea) Horas() float64 {
	return math.Round(float64(d.Segundos)/3600*100) / 100
}
func (d ProyectoTime) Horas() float64 {
	return math.Round(float64(d.Segundos)/3600*100) / 100
}

// ================================================================ //
