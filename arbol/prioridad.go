package arbol

import "github.com/pargomx/gecko/gko"

const ErrPrioridadInvalida gko.ErrorKey = "prioridad_invalida"

func prioridadValida(prioridad int) bool {
	return prioridad >= 0 && prioridad <= 3
}

func (his HistoriaDeUsuario) EsPrioridadMust() bool {
	return his.Prioridad == 3
}
func (his HistoriaDeUsuario) EsPrioridadShould() bool {
	return his.Prioridad == 2
}
func (his HistoriaDeUsuario) EsPrioridadCould() bool {
	return his.Prioridad == 1
}
func (his HistoriaDeUsuario) EsPrioridadWont() bool {
	return his.Prioridad == 0
}
