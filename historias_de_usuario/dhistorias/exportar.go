package dhistorias

import (
	"fmt"
	"io"
	"monorepo/historias_de_usuario/ust"
	"strings"
)

func ExportarMarkdown(w io.Writer, repo Repo) error {
	Personas, err := repo.ListNodosPersonas()
	if err != nil {
		return err
	}
	fmt.Fprint(w, "HISTORIAS DE USUARIO\n")
	for _, Persona := range Personas {
		fmt.Fprintf(w, "\n## %s\n", Persona.Nombre)

		Historias, err := repo.ListNodoHistoriasByPadreID(Persona.PersonaID)
		if err != nil {
			return err
		}
		for _, his := range Historias {
			printHistoria(w, his, repo)
		}

	}
	return nil
}

func printHistoria(w io.Writer, his ust.NodoHistoria, repo Repo) error {

	if his.Nivel == 2 {
		fmt.Fprintf(w, "\n#### %s", his.Titulo)
	} else if his.Nivel == 3 {
		fmt.Fprintf(w, "\n##### %s", his.Titulo)
	} else {
		fmt.Fprint(w, strings.Repeat("  ", his.Nivel-3))
		fmt.Fprintf(w, "+ %s", his.Titulo)
	}

	if his.Completada {
		fmt.Fprint(w, " ✔️")
	} else {
		fmt.Fprint(w, " ", his.PrioridadEmoji())
	}
	fmt.Fprint(w, "\n")

	if his.Objetivo != "" {
		// fmt.Fprint(w, strings.Repeat("  ", his.Nivel-2))
		fmt.Fprintf(w, "%s\n", his.Objetivo)
	}

	// Imprimir tareas
	Tareas, err := repo.ListTareasByHistoriaID(his.HistoriaID)
	if err != nil {
		return err
	}
	for _, tarea := range Tareas {
		if his.Nivel > 3 {
			fmt.Fprint(w, strings.Repeat("  ", his.Nivel-2))
		}

		if tarea.Finalizada() {
			fmt.Fprintf(w, "- [ ] %s\n", tarea.Descripcion)
		} else {
			fmt.Fprintf(w, "- [x] %s\n", tarea.Descripcion)
		}
	}

	Historias, err := repo.ListNodoHistoriasByPadreID(his.HistoriaID)
	if err != nil {
		return err
	}
	for _, his := range Historias {
		printHistoria(w, his, repo)
	}
	return nil
}
