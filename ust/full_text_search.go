package ust

import "html/template"

type SearchResult struct {
	HistoriaID int
	OtroID     string
	Origen     string
	Texto      string
	Subrallado template.HTML
}
