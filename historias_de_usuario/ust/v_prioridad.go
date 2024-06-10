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
