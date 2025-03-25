package main

import (
	"fmt"
	"math"
	"monorepo/dhistorias"
	"monorepo/ust"
	"time"

	"github.com/pargomx/gecko"
	"github.com/pargomx/gecko/gko"
	"github.com/pargomx/gecko/gkt"
)

type DiaReport struct {
	Fecha     string
	Segundos  int
	Proyectos map[string]ProyectoReport
}

type ProyectoReport struct {
	Proyecto        ust.Proyecto
	Segundos        int
	SegundosGestion int
	Latidos         []ust.Latido
	Historias       map[int]HistoriaReport
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

func (d *DiaReport) GetDiaSemanaAbrev() string {
	fecha, err := gkt.ToFecha(d.Fecha)
	if err != nil {
		gko.LogError(err)
		return "ERR"
	}
	switch fecha.Format("Mon") {
	case "Mon":
		return "Lun"
	case "Tue":
		return "Mar"
	case "Wed":
		return "Mié"
	case "Thu":
		return "Jue"
	case "Fri":
		return "Vie"
	case "Sat":
		return "Sáb"
	case "Sun":
		return "Dom"
	default:
		return "Err"
	}
}

func (d *DiaReport) GetFechaAbrev() string {
	fecha, err := gkt.ToFecha(d.Fecha)
	if err != nil {
		gko.LogError(err)
		return "ERR"
	}
	dia := fecha.Day()
	switch fecha.Month() {
	case time.January:
		return fmt.Sprintf("%d ene", dia)
	case time.February:
		return fmt.Sprintf("%d feb", dia)
	case time.March:
		return fmt.Sprintf("%d mar", dia)
	case time.April:
		return fmt.Sprintf("%d abr", dia)
	case time.May:
		return fmt.Sprintf("%d may", dia)
	case time.June:
		return fmt.Sprintf("%d jun", dia)
	case time.July:
		return fmt.Sprintf("%d jul", dia)
	case time.August:
		return fmt.Sprintf("%d ago", dia)
	case time.September:
		return fmt.Sprintf("%d sep", dia)
	case time.October:
		return fmt.Sprintf("%d oct", dia)
	case time.November:
		return fmt.Sprintf("%d nov", dia)
	case time.December:
		return fmt.Sprintf("%d dic", dia)
	default:
		return "Error"
	}
}

// ================================================================ //

type Proyecto1 struct {
	ust.Proyecto
	Historias []ust.NodoHistoria
	Segundos  int
}

func (s *servidor) getMétricas2(c *gecko.Context) error {
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
			segs += historia.SegundosUtilizado
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
	Historias, err := s.repo.ListNodoHistorias()
	if err != nil {
		return err
	}
	HistoriasMap := make(map[int]ust.NodoHistoria, len(Historias))
	for _, historia := range Historias {
		HistoriasMap[historia.HistoriaID] = historia
	}

	// Los días van desde las 6:00am hasta 5:59am (México central).
	ahora := gkt.Now().Truncate(time.Second)
	iniDia := time.Date(ahora.Year(), ahora.Month(), ahora.Day(), 6, 0, 0, 0, gkt.TzMexico)
	finDia := time.Date(ahora.Year(), ahora.Month(), ahora.Day(), 5, 59, 59, 0, gkt.TzMexico)
	if ahora.Hour() < 6 {
		iniDia = iniDia.AddDate(0, 0, -1) // antes 6am, contar desde 6:00am del día anterior.
	} else {
		finDia = finDia.AddDate(0, 0, 1) // después 6am, contar hasta 5:59am del día siguiente.
	}

	type AhoraStruct struct {
		Time            time.Time // Hora actual
		IniDia          time.Time // 6:00am
		FinDia          time.Time // 5:59am
		MinutosSince6am int
		DiaSemana       int
	}
	Ahora := AhoraStruct{
		Time:            ahora,
		MinutosSince6am: int(ahora.Sub(iniDia).Minutes()),
		DiaSemana:       int(ahora.Weekday()),
		IniDia:          iniDia,
		FinDia:          finDia,
	}

	Latidos, err := s.repo.ListLatidos(
		iniDia.AddDate(0, 0, -7).Format(gkt.FormatoFechaHora),
		finDia.Format(gkt.FormatoFechaHora),
	)
	if err != nil {
		return err
	}
	LatidosMapDia := make(map[string][]ust.Latido)
	for _, lati := range Latidos {
		latiTm := lati.Time()
		fechaLat := ""
		if latiTm.Hour() < 6 {
			fechaLat = latiTm.AddDate(0, 0, -1).Format(gkt.FormatoFecha) // antes 6am, contar como el día anterior.
		} else {
			fechaLat = latiTm.Format(gkt.FormatoFecha) // después 6am, contar como el día actual.
		}
		LatidosMapDia[fechaLat] = append(LatidosMapDia[fechaLat], lati)

	}

	Personas, err := s.repo.ListPersonas()
	if err != nil {
		return err
	}
	PersonasMap := make(map[int]ust.Persona, len(Personas))
	for _, per := range Personas {
		PersonasMap[per.PersonaID] = per
	}

	// Quitar día actual y tomar como día anterior si aún no son las 6am.
	if len(ListaDias) > 1 && ListaDias[len(ListaDias)-1] != iniDia.Format(gkt.FormatoFecha) {
		ListaDias = ListaDias[:len(ListaDias)-1]
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
			if tarea.HistoriaID == dhistorias.QUICK_TASK_HISTORIA_ID {
				continue
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
				tarea.SegundosUtilizado = itv.Segundos // reset para solo este día
			} else {
				tar.Segundos += itv.Segundos
				tarea.SegundosUtilizado += itv.Segundos
			}

			tar.Intervalos = append(tar.Intervalos, itv)
			his.Tareas[itv.TareaID] = tar
			pro.Historias[historia.HistoriaID] = his
			Dias[i].Proyectos[historia.ProyectoID] = pro
		}

		// TIEMPO DE GESTIÓN
		for _, lati := range LatidosMapDia[dia] {
			if Dias[i].Proyectos == nil {
				Dias[i].Proyectos = make(map[string]ProyectoReport)
			}
			per, ok := PersonasMap[lati.PersonaID]
			if !ok {
				gko.LogWarnf("getMétricas: persona_id %v no existe en el mapa para latidos en %v", lati.PersonaID, lati.Timestamp)
				continue
			}
			pro, ok := Dias[i].Proyectos[per.ProyectoID]
			if !ok {
				proyecto, err := s.repo.GetProyecto(per.ProyectoID)
				if err != nil {
					return err
				}
				pro = ProyectoReport{
					Proyecto:        *proyecto,
					SegundosGestion: lati.Segundos,
				}
			} else {
				pro.SegundosGestion += lati.Segundos
				pro.Latidos = append(pro.Latidos, lati)
			}
			Dias[i].Proyectos[per.ProyectoID] = pro
		}
	}
	// El último día puede ser en el futuro y estar vacío
	// if len(Dias) > 1 && Dias[len(Dias)-1].Segundos == 0 {
	// 	Dias = Dias[:len(Dias)-2]
	// }

	// Calendario trabajado: últimos 7 días
	Semana := make([]DiaReport, 7)
	if len(Dias) >= 7 {
		Semana = Dias[len(Dias)-7:]
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

	// TABLA DE DÍA TRABAJADO CONTRA PROYECTOS
	ProyectosSimple, err := s.repo.ListProyectos()
	if err != nil {
		return err
	}
	type DiaRow struct {
		Fecha         string
		SegundosLista []int
		SegundosTotal int
	}
	DiasRow := []DiaRow{}
	for _, dia := range Dias {
		row := DiaRow{
			Fecha:         dia.Fecha,
			SegundosLista: make([]int, len(ProyectosSimple)),
		}
		for _, p := range dia.Proyectos {
			segs := p.Segundos + p.SegundosGestion
			for idx, ps := range ProyectosSimple {
				if ps.ProyectoID == p.Proyecto.ProyectoID {
					row.SegundosLista[idx] = segs
					row.SegundosTotal += segs
					break
				}
			}
		}
		DiasRow = append(DiasRow, row)
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
		"DiasRow":             DiasRow,
		"ProyectosSimple":     ProyectosSimple,
		"Proyectos":           Proyectos,
		"Semana":              Semana,
		"Ahora":               Ahora,
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
