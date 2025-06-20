package arbol

import (
	"monorepo/ust"
	"strings"
	"time"

	"github.com/pargomx/gecko/gko"
	"github.com/pargomx/gecko/gkt"
)

const EvTareaIniciada gko.EventKey = "tarea_iniciada"
const EvTareaPausada gko.EventKey = "tarea_pausada"
const EvTareaFinalizada gko.EventKey = "tarea_finalizada"

type argsTareaChangedState struct {
	NodoID int    // TareaID
	TS     string // Timestamp
}

func (s *AppTx) IniciarTarea(TareaID int) error {
	op := gko.Op("IniciarTarea")
	nod, err := s.repo.GetNodo(TareaID)
	if err != nil {
		return op.Err(err)
	}
	if !nod.EsTarea() {
		return op.Msg("Solo se puede registrar un intervalo de trabajo para las tareas")
	}

	// No debe haber otro intervalo abierto (en curso).
	intervalos, err := s.repo.ListIntervalosByNodoID(nod.NodoID)
	if err != nil {
		return op.Err(err)
	}
	for _, itv := range intervalos {
		if itv.TsFin == "" {
			return op.Msgf("Esta tarea ya se había iniciado en %v", itv.TsIni)
		}
	}

	// Iniciar intervalo abierto.
	interv := Intervalo{
		NodoID: nod.NodoID,
		TsIni:  gkt.Now().Format(gkt.FormatoFechaHora),
	}
	err = s.repo.InsertIntervalo(interv)
	if err != nil {
		return op.Err(err)
	}

	// Declarar tarea en curso.
	nod.Estatus = ust.EstatusTareaEnCurso
	err = s.repo.UpdateNodo(nod.NodoID, *nod)
	if err != nil {
		return op.Err(err)
	}

	s.Results.Add(EvTareaIniciada.WithArgs(argsTareaChangedState{
		NodoID: interv.NodoID,
		TS:     interv.TsIni,
	}))
	return nil
}

func (s *AppTx) PausarTarea(TareaID int) error {
	op := gko.Op("PausarTarea")
	nod, err := s.repo.GetNodo(TareaID)
	if err != nil {
		return op.Err(err)
	}
	if !nod.EsTarea() {
		return op.Msg("Solo se puede registrar un intervalo de trabajo para las tareas")
	}

	// Debe haber un intervalo en curso.
	intervalos, err := s.repo.ListIntervalosByNodoID(nod.NodoID)
	if err != nil {
		return op.Err(err)
	}
	var interv *Intervalo
	for _, itv := range intervalos {
		if itv.TsFin == "" {
			interv = &itv
			break
		}
	}
	if interv == nil {
		return op.Msg("No hay ningún intervalo en curso para esta tarea")
	}

	// Finalizar intervalo.
	interv.TsFin = gkt.Now().Format(gkt.FormatoFechaHora)
	err = s.repo.UpdateIntervalo(interv.NodoID, interv.TsIni, *interv)
	if err != nil {
		return op.Err(err)
	}

	// Declarar tarea en pausa.
	nod.Estatus = ust.EstatusTareaEnPausa
	err = s.repo.UpdateNodo(nod.NodoID, *nod)
	if err != nil {
		return op.Err(err)
	}

	s.Results.Add(EvTareaPausada.WithArgs(argsTareaChangedState{
		NodoID: interv.NodoID,
		TS:     interv.TsIni,
	}))
	return nil
}

func (s *AppTx) FinalizarTarea(TareaID int) error {
	op := gko.Op("FinalizarTarea")
	nod, err := s.repo.GetNodo(TareaID)
	if err != nil {
		return op.Err(err)
	}
	if !nod.EsTarea() {
		return op.Msg("Solo se puede registrar un intervalo de trabajo para las tareas")
	}

	// Probablemente haya un intervalo en curso, pero no es necesario.
	intervalos, err := s.repo.ListIntervalosByNodoID(nod.NodoID)
	if err != nil {
		return op.Err(err)
	}
	var interv *Intervalo
	for _, itv := range intervalos {
		if itv.TsFin == "" {
			interv = &itv
			break
		}
	}
	if interv != nil {
		// Finalizar intervalo en curso.
		interv.TsFin = gkt.Now().Format(gkt.FormatoFechaHora)
		err = s.repo.UpdateIntervalo(interv.NodoID, interv.TsIni, *interv)
		if err != nil {
			return op.Err(err)
		}
	}

	// Declarar tarea finalizada.
	nod.Estatus = ust.EstatusTareaFinalizada
	err = s.repo.UpdateNodo(nod.NodoID, *nod)
	if err != nil {
		return op.Err(err)
	}

	s.Results.Add(EvTareaFinalizada.WithArgs(argsTareaChangedState{
		NodoID: interv.NodoID,
		TS:     interv.TsIni,
	}))
	return nil
}

// ================================================================ //

const EvIntervaloParchado gko.EventKey = "intervalo_parchado"

type ArgsParcharIntervalo struct {
	NodoID  int    // TareaID
	TsID    string // Identificar intervalo
	Cambiar string // "INI" o "FIN"
	NewTS   string // New timestamp
}

func (s *AppTx) ParcharIntervalo(args ArgsParcharIntervalo) error {
	op := gko.Op("ParcharIntervalo")
	interv, err := s.repo.GetIntervalo(args.NodoID, args.TsID)
	if err != nil {
		return op.Err(err)
	}

	// Validar new timestamp
	if args.NewTS == "" {
		return op.Msg("Nada que cambiar")
	}
	args.NewTS = strings.Replace(args.NewTS, "T", " ", 1)
	_, err = time.Parse(gkt.FormatoFechaHora, args.NewTS)
	if err != nil {
		return op.Err(err).
			Msgf("Timestamp '%v' inválida. Debe tener formato AAAA-MM-DD HH:MM:SS",
				args.NewTS)
	}

	// Hacer el cambio
	if args.Cambiar == "INI" {
		interv.TsIni = args.NewTS
	} else if args.Cambiar == "FIN" {
		interv.TsFin = args.NewTS
	} else {
		return op.Str("no se sabe si cambiar INI o FIN de intervalo")
	}

	// Validar intervalo: final debe ser después de inicio
	inicio, err := time.Parse(gkt.FormatoFechaHora, interv.TsIni)
	if err != nil {
		return op.Op("ParseToCheckIntervalo").Err(err)
	}
	fin, err := time.Parse(gkt.FormatoFechaHora, interv.TsFin)
	if err != nil {
		return op.Op("ParseToCheckIntervalo").Err(err)
	}
	if fin.Before(inicio) {
		return op.Msg("La fecha de final debe ser posterior a la de inicio")
	}

	// Guardar cambios
	err = s.repo.UpdateIntervalo(args.NodoID, args.TsID, *interv)
	if err != nil {
		return op.Err(err)
	}

	s.Results.Add(EvIntervaloParchado.WithArgs(args))
	return nil
}
