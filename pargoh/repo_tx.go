package main

import (
	"monorepo/arbol"
	"monorepo/dhistorias"
	"monorepo/sqlitearbol"

	"github.com/pargomx/gecko"
	"github.com/pargomx/gecko/eventsqlite"
	"github.com/pargomx/gecko/gko"
	"github.com/pargomx/gecko/gkoid"
	"github.com/pargomx/gecko/sqlitedb"
)

type readhdl struct {
	db      *sqlitedb.SqliteDB
	repo    arbol.ReadRepo
	repoOld dhistorias.Repo
}

type writehdl struct {
	db        *sqlitedb.SqliteDB
	eventRepo *eventsqlite.EventRepoSqlite
	app       *arbol.Servicio
	reloader  reloader // websocket.go
}

type handlerTx struct {
	app  *arbol.AppTx
	repo arbol.Repo
	db   *sqlitedb.Transaccion
}

type handlerTxFunc func(c *gecko.Context, tx *handlerTx) error

// Inicia una transacción en la base de datos, crea un repositorio y comienza
// una transacción de aplicación que es entregada al handler. Cuando el handler
// retorna: si no hay error hace Commit, si hay error o panic hace rollback.
func (s *writehdl) inTx(handler handlerTxFunc) gecko.HandlerFunc {

	ResponsableID := gkoid.Decimal(1)

	return func(c *gecko.Context) error {
		dbTx, dbErr := s.db.Begin()
		if dbErr != nil {
			return gko.ErrNoDisponible.Op("inTx.Begin").Err(dbErr).Msg("Servidor no dispoinible para esta transacción")
		}

		// Catch panic to do rollback
		alreadyRolledBack := false
		defer func() {
			if p := recover(); p != nil && !alreadyRolledBack {
				dbErr = dbTx.Rollback()
				if dbErr != nil {
					gko.Op("inTx.OnDeferPanic").Op("Rollback").Err(dbErr).Log()
				}
				panic(p) // re-throw panic after rollback
			}
		}()

		repoTx := sqlitearbol.NuevoRepo(dbTx)
		eventStore := &gko.EventStore{
			Repo:       s.eventRepo.NuevoRepoWrite(dbTx),
			Results:    &gko.TxResult{},
			ConsoleLog: true,
		}
		appTx := arbol.NewTx(ResponsableID, repoTx, eventStore)
		appErr := handler(c, &handlerTx{
			repo: repoTx,
			app:  appTx,
			db:   dbTx,
		})
		if appErr != nil {
			dbErr = dbTx.Rollback()
			if dbErr != nil {
				gko.Op("inTx.OnHandlerError").Op("Rollback").Err(dbErr).Log()
			}
			alreadyRolledBack = true
			return gko.Err(appErr)
		}

		if appTx.Rollback {
			dbErr = dbTx.Rollback()
			if dbErr != nil {
				gko.Op("inTx.EndWithRollback").Op("Rollback").Err(dbErr).Log()
			}

		} else {
			dbErr = dbTx.Commit() // defer necesita poder leer este error
			if dbErr != nil {
				return gko.Op("inTx.Commit").Err(dbErr)
			}
		}

		return nil
	}
}

// ================================================================ //
// ========== Transacción simple ================================== //

// Inicia una transacción en la base de datos, crea un repositorio y comienza
// una transacción de aplicación que es entregada a la función. Cuando ésta
// retorna: si no hay error hace Commit, si hay error o panic hace rollback.
func (s *servidor) inTx(ResponsableID gkoid.Decimal, function func(tx *arbol.AppTx) error) error {
	dbTx, dbErr := s.db.Begin()
	if dbErr != nil {
		return gko.ErrNoDisponible.Op("inTx.Begin").Err(dbErr).Msg("Servidor no dispoinible para esta transacción")
	}

	// Catch panic to do rollback
	alreadyRolledBack := false
	defer func() {
		if p := recover(); p != nil && !alreadyRolledBack {
			dbErr = dbTx.Rollback()
			if dbErr != nil {
				gko.Op("inTx.OnDeferPanic").Op("Rollback").Err(dbErr).Log()
			}
			panic(p) // re-throw panic after rollback
		}
	}()

	repoTx := sqlitearbol.NuevoRepo(dbTx)
	eventStore := &gko.EventStore{
		Repo:       s.eventRepo.NuevoRepoWrite(dbTx),
		Results:    &gko.TxResult{},
		ConsoleLog: true,
	}
	appTx := arbol.NewTx(ResponsableID, repoTx, eventStore)
	appErr := function(appTx)
	if appErr != nil {
		dbErr = dbTx.Rollback()
		if dbErr != nil {
			gko.Op("inTx.OnHandlerError").Op("Rollback").Err(dbErr).Log()
		}
		alreadyRolledBack = true
		return gko.Err(appErr)
	}

	if appTx.Rollback {
		dbErr = dbTx.Rollback()
		if dbErr != nil {
			gko.Op("inTx.EndWithRollback").Op("Rollback").Err(dbErr).Log()
		}

	} else {
		dbErr = dbTx.Commit() // defer necesita poder leer este error
		if dbErr != nil {
			return gko.Op("inTx.Commit").Err(dbErr)
		}
	}

	return nil
}
