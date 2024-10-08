package dhistorias

import (
	"monorepo/ust"
	"time"

	"github.com/pargomx/gecko/gko"
)

type GestionTimeTracker struct {
	repo      gestionTimeTrackerRepo
	buffer    map[string]int // [proyectoID] segundos
	maxBuffer int            // cuántos segundos acumular antes de guardar en DB.
}

type gestionTimeTrackerRepo interface {
	GetProyecto(ProyectoID string) (*ust.Proyecto, error)
	UpdateProyecto(pro ust.Proyecto) error
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
		buffer:    make(map[string]int),
		maxBuffer: maxBuffer, // segundos
	}
}

// Agregar tiempo utilizado en gestionar un proyecto.
func (s *GestionTimeTracker) AddTimeSpent(ProyectoID string, segundos int) error {
	const op = "AddTimeSpent"
	pry, err := s.repo.GetProyecto(ProyectoID)
	if err != nil {
		return gko.Err(err).Op(op)
	}
	if segundos < 0 {
		return gko.ErrDatoInvalido().Msg("El tiempo no puede ser negativo").Op(op)
	}
	err = s.repo.InsertLatido(ust.Latido{
		Timestamp:  time.Now().In(locationMexicoCity).Format("2006-01-02 15:04:05"),
		Segundos:   segundos,
		ProyectoID: pry.ProyectoID,
	})
	if err != nil {
		return gko.Err(err).Op(op)
	}

	s.buffer[pry.ProyectoID] += segundos
	if s.buffer[pry.ProyectoID] > s.maxBuffer {
		pry.TiempoGestion += s.buffer[pry.ProyectoID]
		err = s.repo.UpdateProyecto(*pry)
		if err != nil {
			return gko.Err(err).Op(op)
		}
		s.buffer[pry.ProyectoID] = 0
	}

	return nil
}
