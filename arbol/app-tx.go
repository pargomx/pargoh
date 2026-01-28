package arbol

import (
	"github.com/pargomx/gecko/gko"
	"github.com/pargomx/gecko/gkoid"
)

// ================================================================ //
// ========== Transacción ========================================= //

type AppTx struct {
	repo     Repo // tx repo
	events   *gko.EventStore
	Rollback bool // if true, rollback when finished

	ResponsableID gkoid.Decimal // UsuarioID del responsable.
}

func NewTx(ResponsableID gkoid.Decimal, repoTx Repo, eventStore *gko.EventStore) *AppTx {
	return &AppTx{
		repo:          repoTx,
		events:        eventStore,
		ResponsableID: ResponsableID,
	}
}

func (s *AppTx) riseEvent(key gko.EventKey, data gko.EventData) error {
	_, err := s.events.Rise(s.ResponsableID, key, data)
	if err != nil {
		gko.LogError(err)
		s.Rollback = true // Rollback si hay error.
	}
	return err
}
