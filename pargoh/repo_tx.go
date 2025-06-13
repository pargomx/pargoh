package main

import (
	"monorepo/dhistorias"
	"monorepo/sqlitepuente"

	"github.com/pargomx/gecko/gko"
)

type repoTx struct {
	repo     dhistorias.Repo
	Commit   func() error
	Rollback func() error
}

func (s *servidor) newRepoTx() (*repoTx, error) {
	tx, err := s.db.Begin()
	if err != nil {
		return nil, gko.Err(err)
	}
	return &repoTx{
		repo:     sqlitepuente.NuevoRepo(tx),
		Commit:   tx.Commit,
		Rollback: tx.Rollback,
	}, nil
}
