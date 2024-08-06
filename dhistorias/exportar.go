package dhistorias

import (
	"fmt"
	"io"
	"monorepo/ust"
	"strings"

	"github.com/gingfrederik/docx"
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
			printHistoriaMarkdown(w, his, repo)
		}

	}
	return nil
}

func printHistoriaMarkdown(w io.Writer, his ust.NodoHistoria, repo Repo) error {

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
		printHistoriaMarkdown(w, his, repo)
	}
	return nil
}

// ================================================================ //
// ================================================================ //

func ExportarDocx(repo Repo, filepath string) error {
	if filepath == "" {
		return fmt.Errorf("ruta de archivo vacía")
	}
	if !strings.HasSuffix(filepath, ".docx") {
		filepath += ".docx"
	}
	Personas, err := repo.ListNodosPersonas()
	if err != nil {
		return err
	}
	f := docx.NewFile()
	f.AddParagraph().AddText("HISTORIAS DE USUARIO").Size(24).Color("3a344a") // TITULO

	for _, Persona := range Personas {
		f.AddParagraph().AddText("").Size(12)
		f.AddParagraph().AddText(Persona.Nombre).Size(22).Color("0b3d42") // PERSONA

		Historias, err := repo.ListNodoHistoriasByPadreID(Persona.PersonaID)
		if err != nil {
			return err
		}
		for _, his := range Historias {
			printHistoriaDocx(f, his, repo)
		}

	}
	// f.Write(w)
	return f.Save(filepath)
}

func printHistoriaDocx(f *docx.File, his ust.NodoHistoria, repo Repo) error {

	txt := ""
	txt += strings.Repeat(">", his.Nivel-1)
	txt += " "
	txt += his.Titulo
	if his.Completada {
		txt += " ✔️"
	} else {
		txt += " " + his.PrioridadEmoji()
	}

	color := "3c207d"
	switch his.Nivel {
	case 2:
		color = "4c22a1"
	case 3:
		color = "1d517a"
	case 4:
		color = "5c7516"
	case 5:
		color = "6b250c"
	case 6:
		color = "571056"
	}

	f.AddParagraph().AddText("").Size(12)
	f.AddParagraph().AddText(txt).Size(22 - 2*his.Nivel).Color(color) // HISTORIA

	if his.Objetivo != "" {
		f.AddParagraph().AddText(his.Objetivo).Size(12).Color("3c3b40") // OBJETIVO
	}

	Tareas, err := repo.ListTareasByHistoriaID(his.HistoriaID)
	if err != nil {
		return err
	}
	for _, tarea := range Tareas {
		if tarea.Finalizada() {
			f.AddParagraph().AddText("[ ] " + tarea.Descripcion).Size(12) // TAREAS
		} else {
			f.AddParagraph().AddText("[x] " + tarea.Descripcion).Size(12)
		}
	}

	Historias, err := repo.ListNodoHistoriasByPadreID(his.HistoriaID)
	if err != nil {
		return err
	}
	for _, his := range Historias {
		printHistoriaDocx(f, his, repo)
	}
	return nil
}
