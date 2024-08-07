package dhistorias

import (
	"monorepo/ust"

	"github.com/pargomx/gecko/gko"
)

func ImportarFake(repo Repo) error {
	op := gko.Op("CrearFake")

	personas := []ust.Persona{
		{PersonaID: 11, Nombre: "Juan"},
		{PersonaID: 22, Nombre: "Pedro"},
		{PersonaID: 33, Nombre: "Maria"},
	}
	for _, per := range personas {
		err := InsertarPersona(per, repo)
		if err != nil {
			return op.Err(err)
		}
	}
	historias := []ust.NodoHistoria{
		{HistoriaID: 1, PadreID: 11, Titulo: "Hacer una grande", Objetivo: "Objetivo 1", Prioridad: 1, Completada: false},
		{HistoriaID: 2, PadreID: 22, Titulo: "Lograr un objetivo", Objetivo: "Objetivo 2", Prioridad: 2, Completada: false},
		{HistoriaID: 3, PadreID: 33, Titulo: "Ingresar a un edificio", Objetivo: "Objetivo 3", Prioridad: 3, Completada: false},

		{HistoriaID: 101, PadreID: 1, Titulo: "Meter dos niveles", Objetivo: "Objetivo 1", Prioridad: 1, Completada: false},
		{HistoriaID: 102, PadreID: 1, Titulo: "Encontrar a dos personas", Objetivo: "Objetivo 2", Prioridad: 2, Completada: false},
		{HistoriaID: 103, PadreID: 1, Titulo: "Leer un par de PadreID:0,historias", Objetivo: "Objetivo 3", Prioridad: 3, Completada: false},

		{HistoriaID: 201, PadreID: 101, Titulo: "Vencer doble vez", Objetivo: "Objetivo 1", Prioridad: 1, Completada: false},
		{HistoriaID: 202, PadreID: 101, Titulo: "Pensarlo dos veces antes de hablar", Objetivo: "Objetivo 2", Prioridad: 2, Completada: false},

		{HistoriaID: 1101, PadreID: 201, Titulo: "Decir tres tristes tigres", Objetivo: "", Prioridad: 1, Completada: false},
		{HistoriaID: 1102, PadreID: 201, Titulo: "Montar un triciclo", Objetivo: "Para ir más rápido al mercado", Prioridad: 2, Completada: false},
		{HistoriaID: 1103, PadreID: 201, Titulo: "Hacer llamada tripartita", Objetivo: "Objetivo 3", Prioridad: 3, Completada: false},
		{HistoriaID: 1103, PadreID: 201, Titulo: "Pensar en las tres naranjas", Objetivo: "Objetivo 3", Prioridad: 3, Completada: false},

		{HistoriaID: 11101, PadreID: 1101, Titulo: "Construir una cuatrimoto", Objetivo: "", Prioridad: 1, Completada: false},
		{HistoriaID: 11102, PadreID: 1101, Titulo: "Rodear las cuatro esquinas de un cubo", Objetivo: "", Prioridad: 2, Completada: false},
	}
	for _, his := range historias {
		err := AgregarHistoria(his.PadreID, ust.Historia{
			HistoriaID: his.HistoriaID,
			Titulo:     his.Titulo,
			Objetivo:   his.Objetivo,
			Prioridad:  his.Prioridad,
			Completada: his.Completada,
		}, repo)
		if err != nil {
			return op.Err(err)
		}
	}

	tareas := []ust.Tarea{
		{TareaID: 1, HistoriaID: 1, Descripcion: "Hacer una grande", Tipo: ust.TipoTareaConf},
		{TareaID: 2, HistoriaID: 1, Descripcion: "Hacer una pequeña", Tipo: ust.TipoTareaDominio},
		{TareaID: 3, HistoriaID: 101, Descripcion: "Hacer una mediana", Tipo: ust.TipoTareaHandlr},
		{TareaID: 4, HistoriaID: 101, Descripcion: "Hacer una chiquita", Tipo: ust.TipoTareaIndefinido},
	}
	for _, tar := range tareas {
		err := AgregarTarea(tar, repo)
		if err != nil {
			return op.Err(err)
		}
	}

	tramos := []ust.Tramo{
		{HistoriaID: 1, Texto: "Iniciar sesión"},
		{HistoriaID: 1, Texto: "Hacer una pequeña"},
		{HistoriaID: 1, Texto: "Hacer una mediana"},
		{HistoriaID: 1, Texto: "Ver el mensaje de una grande"},
		{HistoriaID: 101, Texto: "Hacer una mediana"},
		{HistoriaID: 101, Texto: "Hacer una chiquita"},
	}
	for _, tra := range tramos {
		err := NuevoTramoDeViaje(repo, tra.HistoriaID, tra.Texto)
		if err != nil {
			return op.Err(err)
		}
	}

	return nil
}
