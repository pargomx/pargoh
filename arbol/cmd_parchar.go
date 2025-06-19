package arbol

import (
	"fmt"
	"monorepo/ust"
	"strings"

	"github.com/pargomx/gecko/gko"
	"github.com/pargomx/gecko/gkt"
)

const EvParcharNodo gko.EventKey = "nodo_parchado"

type ArgsParcharNodo struct {
	NodoID int
	Campo  string // Campo a parchar
	NewVal string // Nuevo valor
}

func (s *AppTx) ParcharNodo(args ArgsParcharNodo) error {
	op := gko.Op("ParcharNodo")
	if args.NodoID == 0 {
		return op.E(gko.ErrDatoIndef).Str("NodoID indefinido")
	}
	if args.Campo == "" {
		return op.E(gko.ErrDatoIndef).Str("Campo a parchar indefinido")
	}
	nod, err := s.repo.GetNodo(args.NodoID)
	if err != nil {
		return op.Err(err)
	}

	switch args.Campo {

	case "nombre":
		nod.Titulo = gkt.SinEspaciosExtra(args.NewVal)
	case "titulo":
		nod.Titulo = gkt.SinEspaciosExtra(args.NewVal)
	case "texto":
		nod.Titulo = strings.TrimSpace(args.NewVal)
		if nod.Tipo == "" {
			// delete if nodo is empty.
		}

	case "objetivo":
		nod.Objetivo = gkt.SinEspaciosExtraConSaltos(args.NewVal)

	case "descripcion":
		nod.Descripcion = gkt.SinEspaciosExtraConSaltos(args.NewVal)

	case "notas":
		nod.Notas = gkt.SinEspaciosExtraConSaltos(args.NewVal)

	case "color":
		nod.Color = gkt.SinEspaciosNinguno(args.NewVal)

	case "prioridad":
		nod.Prioridad, _ = gkt.ToInt(args.NewVal)
		if !prioridadValida(nod.Prioridad) {
			return op.E(ErrPrioridadInvalida).Msg("La prioridad debe estar entre 0 y 3")
		}

	case "importancia":
		if args.NewVal != "" {
			nod.Prioridad, _ = gkt.ToInt(args.NewVal)
		} else {
			tarea := nod.ToTarea()
			tarea.Prioridad = ((tarea.Prioridad + 1) % 3) + 1
			nod.Prioridad = tarea.Prioridad
			args.NewVal = fmt.Sprint(tarea.Prioridad)
		}
		if !prioridadValida(nod.Prioridad) {
			return op.E(ErrPrioridadInvalida).Msg("La importancia debe estar entre 0 y 3")
		}

	case "completada":
		completada := gkt.ToBool(args.NewVal)
		if completada {
			nod.Estatus = 1
		} else {
			nod.Estatus = 0
		}

	case "marcar_regla":
		regla := nod.ToRegla()
		regla.Estatus = (regla.Estatus + 1) % 3
		nod.Estatus = regla.Estatus
		args.NewVal = fmt.Sprint(regla.Estatus)

	case "presupuesto":
		if args.NewVal == "" {
			nod.Segundos = 0
		} else {
			horas, err := gkt.ToInt(args.NewVal)
			if err != nil {
				return op.Err(err)
			}
			if horas < 0 {
				return op.E(gko.ErrDatoInvalido).Msg("El presupuesto debe ser positivo")
			}
			if horas > 30 {
				return op.E(gko.ErrDatoInvalido).Msgf("Establezca un presupuesto menor, %v son demasiadas horas para una sola historia.", horas)
			}
			nod.Segundos = horas * 3600
		}

	case "estimado":
		estimado, err := ust.NuevaDuraciónSegundos(args.NewVal)
		if err != nil {
			return op.Err(err)
		}
		if estimado <= 0 {
			return op.E(gko.ErrDatoInvalido).Msg("El estimado debe ser mayor a 0")
		}

	default:
		return op.E(gko.ErrDatoInvalido).Strf("campo no soportado: %v", args.Campo).Msg("No se pudo guardar el cambio")
	}

	err = s.repo.UpdateNodo(nod.NodoID, *nod)
	if err != nil {
		return op.Err(err)
	}

	s.Results.Add(EvParcharNodo.WithArgs(args))
	return nil
}
