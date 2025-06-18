package main

import (
	"monorepo/arbol"
	"monorepo/dhistorias"
	"monorepo/sqlitearbol"
	"monorepo/sqliteust"

	"github.com/pargomx/gecko"
	"github.com/pargomx/gecko/gko"
)

type serverTx struct {
	repo     arbol.Repo
	repoOld  dhistorias.Repo
	Commit   func() error
	Rollback func() error
}

func (s *servidor) newRepoTx() (*serverTx, error) {
	tx, err := s.db.Begin()
	if err != nil {
		return nil, gko.Err(err)
	}
	return &serverTx{
		repo:     sqlitearbol.NuevoRepo(tx),
		repoOld:  sqliteust.NuevoRepo(tx),
		Commit:   tx.Commit,
		Rollback: tx.Rollback,
	}, nil
}

// Inicia una transacción en la base de datos, crea un repositorio y comienza
// una transacción de aplicación que es entregada al handler. Cuando el handler
// retorna: si no hay error hace Commit, si hay error o panic hace rollback.
func (s *servidor) inTx(fn func(c *gecko.Context, tx *arbol.AppTx) error) gecko.HandlerFunc {
	return func(c *gecko.Context) error {
		dbTx, err := s.db.Begin()
		if err != nil {
			return err
		}
		defer func() {
			if p := recover(); p != nil {
				_ = dbTx.Rollback()
				panic(p) // re-throw panic after rollback
			} else if err != nil {
				_ = dbTx.Rollback()
			} else {
				err = dbTx.Commit()
			}
		}()
		repoTx := sqlitearbol.NuevoRepo(dbTx)
		appTx := s.app.NewTx(repoTx)
		err = fn(c, appTx)
		// necesario hacer en dos líneas para que tenga efecto el defer.
		return err
	}
}
