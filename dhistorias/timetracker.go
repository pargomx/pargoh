package dhistorias

import (
	"monorepo/ust"
	"time"

	"github.com/pargomx/gecko/gko"
)

type GestionTimeTracker struct {
	repo      gestionTimeTrackerRepo
	buffer    map[int]int // [personaID] segundos
	maxBuffer int         // cuántos segundos acumular antes de guardar en DB.
}

type gestionTimeTrackerRepo interface {
	GetPersona(PersonaID int) (*ust.Persona, error)
	UpdatePersona(per ust.Persona) error
	InsertLatido(lat ust.Latido) error
}

// Para llevar un registro de cuánto tiempo se ha invertido en gestión.
// El buffer acumula maxBuffer segundos antes de escribir en la DB.
// Si el buffer es <= 0 no se usa el buffer.
func NewGestionTimeTracker(repo gestionTimeTrackerRepo, maxBuffer int) *GestionTimeTracker {
	if repo == nil {
		gko.FatalExit("NewGestionTimeTracker: repo es nil")
	}
	return &GestionTimeTracker{
		repo:      repo,
		buffer:    make(map[int]int),
		maxBuffer: maxBuffer, // segundos
	}
}

func (s *GestionTimeTracker) AddTimeSpent(PersonaID int, segundos int) error {
	const op = "AddTimeSpent"
	per, err := s.repo.GetPersona(PersonaID)
	if err != nil {
		return gko.Err(err).Op(op)
	}
	if segundos < 0 {
		return gko.ErrDatoInvalido.Msg("El tiempo no puede ser negativo").Op(op)
	}
	err = s.repo.InsertLatido(ust.Latido{
		Timestamp: time.Now().In(locationMexicoCity).Format("2006-01-02 15:04:05"),
		Segundos:  segundos,
		PersonaID: per.PersonaID,
	})
	if err != nil {
		return gko.Err(err).Op(op)
	}

	s.buffer[per.PersonaID] += segundos
	if s.buffer[per.PersonaID] > s.maxBuffer {
		per.SegundosGestion += s.buffer[per.PersonaID]
		err = s.repo.UpdatePersona(*per)
		if err != nil {
			return gko.Err(err).Op(op)
		}
		s.buffer[per.PersonaID] = 0
	}

	return nil
}
