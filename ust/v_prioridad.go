package ust

func (his Historia) PrioridadEmoji() string {
	switch his.Prioridad {
	case 0:
		return ""
	case 1:
		return "😶‍🌫️"
	case 2:
		return "🤔"
	case 3:
		return "🤩"
	default:
		return "🤯"
	}
}

func (his NodoHistoria) PrioridadEmoji() string {
	switch his.Prioridad {
	case 0:
		return ""
	case 1:
		return "😶‍🌫️"
	case 2:
		return "🤔"
	case 3:
		return "🤩"
	default:
		return "🤯"
	}
}

func (his Historia) EsPrioridadMust() bool {
	return his.Prioridad == 3
}
func (his Historia) EsPrioridadShould() bool {
	return his.Prioridad == 2
}
func (his Historia) EsPrioridadCould() bool {
	return his.Prioridad == 1
}
func (his Historia) EsPrioridadWont() bool {
	return his.Prioridad == 0
}

func (his NodoHistoria) EsPrioridadMust() bool {
	return his.Prioridad == 3
}
func (his NodoHistoria) EsPrioridadShould() bool {
	return his.Prioridad == 2
}
func (his NodoHistoria) EsPrioridadCould() bool {
	return his.Prioridad == 1
}
func (his NodoHistoria) EsPrioridadWont() bool {
	return his.Prioridad == 0
}
