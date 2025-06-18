package main

import (
	"monorepo/arbol"
	"monorepo/dhistorias"
	"monorepo/sqlitearbol"
	"monorepo/sqlitedb"
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

type handlerTx struct {
	app  *arbol.AppTx
	repo arbol.Repo
	db   *sqlitedb.Transaccion
}

// Inicia una transacción en la base de datos, crea un repositorio y comienza
// una transacción de aplicación que es entregada al handler. Cuando el handler
// retorna: si no hay error hace Commit, si hay error o panic hace rollback.
func (s *servidor) inTx(handler func(c *gecko.Context, tx *handlerTx) error) gecko.HandlerFunc {
	return func(c *gecko.Context) error {
		dbTx, err := s.db.Begin()
		if err != nil {
			return err
		}
		gko.LogDebug("TX Started")

		defer func() {
			if p := recover(); p != nil {
				err = dbTx.Rollback()
				if err != nil {
					gko.Err(err).Log()
				}
				gko.LogDebug("TX Rollback after panic")
				panic(p) // re-throw panic after rollback
			} else if err != nil {
				err = dbTx.Rollback()
				if err != nil {
					gko.Err(err).Log()
				}
				gko.LogDebug("TX Rollback in defer")
			}
		}()

		repoTx := sqlitearbol.NuevoRepo(dbTx)
		appTx := s.app.NewTx(repoTx)

		err = handler(c, &handlerTx{
			repo: repoTx,
			app:  appTx,
			db:   dbTx,
		})
		if err != nil {
			err = dbTx.Rollback()
			if err != nil {
				gko.Err(err).Log()
			}
			gko.LogDebug("TX Rollback")
		}

		err = dbTx.Commit() // necesario hacer en dos líneas para que tenga efecto el defer.
		if err != nil {
			gko.Err(err).Log()
		}
		gko.LogDebug("TX Commited")

		return err
	}
}
