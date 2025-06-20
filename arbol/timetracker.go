package arbol

import (
	"github.com/pargomx/gecko/gko"
	"github.com/pargomx/gecko/gkt"
)

// Registra el tiempo que se pasa utilzando la aplicación.
type AppTimeTracker struct {
	repo      timeTrackerRepo
	buffer    map[int]int // [nodoID] segundos
	maxBuffer int         // cuántos segundos acumular antes de guardar en DB.
}

// Para llevar un registro de cuánto tiempo se ha invertido en gestión.
// El buffer acumula maxBuffer segundos antes de escribir en la DB.
// Si el buffer es <= 0 no se usa el buffer.
func NewAppTimeTracker(repo timeTrackerRepo, maxBuffer int) *AppTimeTracker {
	if repo == nil {
		gko.FatalExit("NewGestionTimeTracker: repo es nil")
	}
	return &AppTimeTracker{
		repo:      repo,
		buffer:    make(map[int]int),
		maxBuffer: maxBuffer, // segundos
	}
}

func (s *AppTimeTracker) AddTimeSpent(NodoID int, segundos int) error {
	op := gko.Op("AddTimeSpent")
	nod, err := s.repo.GetNodo(NodoID)
	if err != nil {
		return op.Err(err)
	}
	if segundos < 0 {
		return op.Str("el tiempo no puede ser negativo")
	}
	err = s.repo.InsertLatido(Latido{
		TsLatido: gkt.Now().Format(gkt.FormatoFechaHora),
		NodoID:   nod.NodoID,
		Segundos: segundos,
	})
	if err != nil {
		return op.Err(err)
	}
	/*
		Ya no se usa este campo para materializar los latidos acumulados.
		s.buffer[nod.NodoID] += segundos
		if s.buffer[nod.NodoID] > s.maxBuffer {
			nod.Segundos += s.buffer[nod.NodoID]
			err = s.repo.UpdateNodo(nod.NodoID, *nod)
			if err != nil {
				return op.Err(err)
			}
			s.buffer[nod.NodoID] = 0
		}
	*/
	// Sin evento por favor.
	return nil
}
