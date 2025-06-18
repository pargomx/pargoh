package arbol

import (
	"github.com/pargomx/gecko/gko"
)

// ================================================================ //
// ========== Servicio ============================================ //

type Servicio struct {
	cfg Config
}

type Config struct{}

func NuevoServicio(cfg Config) (*Servicio, error) {
	// if cfg.Repo == nil {
	// 	return nil, gko.ErrNoDisponible.Str("NuevoServicio: repo es nil")
	// }

	return &Servicio{
		cfg: cfg,
	}, nil
}

// ================================================================ //
// ========== Transacción ========================================= //

func (s *Servicio) NewTx(repoTx Repo) *AppTx {
	return &AppTx{
		s:       s,
		repo:    repoTx,
		Results: &gko.TxResult{},
	}
}

type AppTx struct {
	s       *Servicio
	repo    Repo          // Podría ser db.Tx
	Results *gko.TxResult // Eventos y errores
}
