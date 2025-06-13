package sqlitepuente

import (
	"monorepo/sqlitearbol"
	"monorepo/sqliteust"

	"github.com/pargomx/gecko/sqlitedb"
)

type Repositorio struct {
	db  sqlitedb.Ejecutor
	old *sqliteust.Repositorio
	nvo *sqlitearbol.Repositorio
}

func NuevoRepo(db sqlitedb.Ejecutor) *Repositorio {
	return &Repositorio{
		db:  db,
		old: sqliteust.NuevoRepo(db),
		nvo: sqlitearbol.NuevoRepo(db),
	}
}
