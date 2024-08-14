package dhistorias

import (
	"fmt"
	"io"
	"monorepo/ust"
	"strings"

	"github.com/gingfrederik/docx"
	"github.com/pargomx/gecko/gko"
)

type ÁrbolCompleto struct {
	Proyectos []ProyectoExport
}

type ProyectoExport struct {
	Proyecto ust.Proyecto
	Personas []PersonaExport
}

type PersonaExport struct {
	Persona   ust.NodoPersona
	Historias []HistoriaExport
}

type HistoriaExport struct {
	Historia  ust.NodoHistoria
	Tareas    []ust.Tarea
	Tramos    []ust.Tramo
	Historias []HistoriaExport
}

// ================================================================ //

func GetArbolCompleto(repo Repo) ([]ProyectoExport, error) {
	op := gko.Op("GetArbolCompleto")
	proyectos, err := repo.ListProyectos()
	if err != nil {
		return nil, op.Err(err)
	}
	Proyectos := make([]ProyectoExport, len(proyectos))
	for i, p := range proyectos {
		personas, err := repo.ListNodosPersonasByProyecto(p.ProyectoID)
		if err != nil {
			return nil, op.Err(err)
		}
		Proyectos[i] = ProyectoExport{
			Proyecto: p,
			Personas: make([]PersonaExport, len(personas)),
		}
		for j, p := range personas {
			historias, err := repo.ListNodoHistorias(p.PersonaID)
			if err != nil {
				return nil, op.Err(err)
			}
			Proyectos[i].Personas[j] = PersonaExport{
				Persona:   p,
				Historias: make([]HistoriaExport, len(historias)),
			}
			for k, h := range historias {
				Proyectos[i].Personas[j].Historias[k] = getHistoriaExportsRecursiva(h, repo)
			}
		}
	}
	return Proyectos, nil
}

// ================================================================ //

func ExportarProyecto(proyectoID string, repo Repo) (*ProyectoExport, error) {
	Proyecto, err := repo.GetProyecto(proyectoID)
	if err != nil {
		return nil, err
	}
	Personas, err := repo.ListNodosPersonasByProyecto(Proyecto.ProyectoID)
	if err != nil {
		return nil, err
	}
	proyectoExport := ProyectoExport{
		Proyecto: *Proyecto,
		Personas: make([]PersonaExport, len(Personas)),
	}
	for i, per := range Personas {
		historias, err := repo.ListNodoHistorias(per.PersonaID)
		if err != nil {
			return nil, err
		}
		Persona := PersonaExport{
			Persona:   per,
			Historias: make([]HistoriaExport, len(historias)),
		}
		for j, his := range historias {
			Persona.Historias[j] = getHistoriaExportsRecursiva(his, repo)
		}
		proyectoExport.Personas[i] = Persona
	}
	return &proyectoExport, nil
}

func getHistoriaExportsRecursiva(his ust.NodoHistoria, repo Repo) HistoriaExport {
	historia := HistoriaExport{
		Historia:  his,
		Historias: nil,
	}
	hijos, err := repo.ListNodoHistorias(his.HistoriaID)
	if err != nil {
		fmt.Println("getHistoriaExportsRecursiva: %w", err)
	}
	historia.Tareas, err = repo.ListTareasByHistoriaID(his.HistoriaID)
	if err != nil {
		fmt.Println("getHistoriaExportsRecursiva: %w", err)
	}
	historia.Tramos, err = repo.ListTramosByHistoriaID(his.HistoriaID)
	if err != nil {
		fmt.Println("getHistoriaExportsRecursiva: %w", err)
	}
	for _, hijo := range hijos {
		historia.Historias = append(historia.Historias, getHistoriaExportsRecursiva(hijo, repo))
	}
	return historia
}

// ================================================================ //

func Importar(p ProyectoExport, repo Repo) error {
	err := repo.InsertProyecto(p.Proyecto)
	if err != nil {
		return err
	}
	for _, per := range p.Personas {
		err = InsertarPersona(ust.Persona{
			PersonaID:   per.Persona.PersonaID,
			ProyectoID:  per.Persona.ProyectoID,
			Nombre:      per.Persona.Nombre,
			Descripcion: per.Persona.Descripcion,
		}, repo)
		if err != nil {
			return err
		}
		for _, his := range per.Historias {
			err = insertHistoriaRecursiva(per.Persona.PersonaID, his, repo)
			if err != nil {
				return err
			}
		}
	}
	gko.LogOkeyf("Importado proyecto %s", p.Proyecto.ProyectoID)
	return nil
}

func insertHistoriaRecursiva(padreID int, his HistoriaExport, repo Repo) error {
	err := AgregarHistoria(padreID, ust.Historia{
		HistoriaID: his.Historia.HistoriaID,
		Titulo:     his.Historia.Titulo,
		Objetivo:   his.Historia.Objetivo,
		Prioridad:  his.Historia.Prioridad,
		Completada: his.Historia.Completada,
	}, repo)
	if err != nil {
		return err
	}
	for _, tarea := range his.Tareas {
		err = repo.InsertTarea(tarea)
		if err != nil {
			return err
		}
	}
	for _, tramo := range his.Tramos {
		err = repo.InsertTramo(tramo)
		if err != nil {
			return err
		}
	}
	for _, h := range his.Historias {
		err = insertHistoriaRecursiva(his.Historia.HistoriaID, h, repo)
		if err != nil {
			return err
		}
	}
	return nil
}

// ================================================================ //

func ExportarMarkdown(proyectoID string, w io.Writer, repo Repo) error {
	Proyecto, err := repo.GetProyecto(proyectoID)
	if err != nil {
		return err
	}
	Personas, err := repo.ListNodosPersonasByProyecto(proyectoID)
	if err != nil {
		return err
	}
	fmt.Fprintf(w, "# %s\n\n%s\n", Proyecto.Titulo, Proyecto.Descripcion)
	for _, Persona := range Personas {
		fmt.Fprintf(w, "\n## %s\n", Persona.Nombre)

		Historias, err := repo.ListNodoHistorias(Persona.PersonaID)
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

	// Imprimir tramos
	Tramos, err := repo.ListTramosByHistoriaID(his.HistoriaID)
	if err != nil {
		return err
	}
	if len(Tramos) > 0 {
		fmt.Fprintln(w, "\nViaje:")
		for _, tramo := range Tramos {
			fmt.Fprintf(w, "%d.%v\n", tramo.Posicion, tramo.Texto)
		}
	}

	// Imprimir tareas
	Tareas, err := repo.ListTareasByHistoriaID(his.HistoriaID)
	if err != nil {
		return err
	}
	if len(Tareas) > 0 {
		fmt.Fprintln(w, "\nTareas:")
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
	}

	Historias, err := repo.ListNodoHistorias(his.HistoriaID)
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

func ExportarDocx(proyectoID string, repo Repo, filepath string) error {
	if filepath == "" {
		return fmt.Errorf("ruta de archivo vacía")
	}
	if !strings.HasSuffix(filepath, ".docx") {
		filepath += ".docx"
	}
	Proyecto, err := repo.GetProyecto(proyectoID)
	if err != nil {
		return err
	}
	Personas, err := repo.ListNodosPersonasByProyecto(proyectoID)
	if err != nil {
		return err
	}
	f := docx.NewFile()
	f.AddParagraph().AddText(Proyecto.Titulo).Size(24).Color("3a344a") // TITULO
	f.AddParagraph().AddText(Proyecto.Descripcion).Size(12)

	for _, Persona := range Personas {
		f.AddParagraph().AddText("").Size(12)
		f.AddParagraph().AddText(Persona.Nombre).Size(22).Color("0b3d42") // PERSONA

		Historias, err := repo.ListNodoHistorias(Persona.PersonaID)
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

	Historias, err := repo.ListNodoHistorias(his.HistoriaID)
	if err != nil {
		return err
	}
	for _, his := range Historias {
		printHistoriaDocx(f, his, repo)
	}
	return nil
}
