package ust

import (
	"fmt"
	"strconv"
	"strings"
)

// Duración representa la duración de una tarea.
// Acepta strings en las siguientes formas, tanto minutos como horas:
// - "15" para un minutos
// - "90" para 90 minutos
// - "90m" para 90 minutos
// - "90 m" para 90 minutos
// - "90 min" para 90 minutos
//
// - "1h" para una hora
// - "1h15" para una hora y 15 minutos
// - "1h 15 min" para una hora y 15 minutos
// - "2:30" para dos horas y 30 minutos
func NuevaDuración(txt string) (int, error) {
	if txt == "" {
		return 0, nil
	}

	txt = strings.ToLower(txt)
	txt = strings.ReplaceAll(txt, " ", "")

	txt = strings.ReplaceAll(txt, "horas", ":")
	txt = strings.ReplaceAll(txt, "hora", ":")
	txt = strings.ReplaceAll(txt, "hrs", ":")
	txt = strings.ReplaceAll(txt, "hr", ":")
	txt = strings.ReplaceAll(txt, "h", ":")

	txt = strings.ReplaceAll(txt, "minutos", "")
	txt = strings.ReplaceAll(txt, "mins", "")
	txt = strings.ReplaceAll(txt, "min", "")
	txt = strings.ReplaceAll(txt, "m", "")

	txt = strings.ReplaceAll(txt, "::", ":")

	split := strings.Split(txt, ":")
	switch len(split) {
	case 1:
		mins, err := strconv.Atoi(split[0])
		if err != nil {
			return 0, fmt.Errorf("duración inválida: %w", err)
		}
		return mins, nil
	case 2:
		hrs, err := strconv.Atoi(split[0])
		if err != nil {
			return 0, fmt.Errorf("duración inválida: %w", err)
		}
		mins := 0
		if split[1] != "" {
			mins, err = strconv.Atoi(split[1])
			if err != nil {
				return 0, fmt.Errorf("duración inválida: %w", err)
			}
		}
		return hrs*60 + mins, nil
	}
	return 0, fmt.Errorf("duración inválida: %s", txt)
}

func MinutosToString(mins int) string {
	if mins == 0 {
		return ""
	}
	if mins < 60 {
		return fmt.Sprintf("%dm", mins)
	}
	h := mins / 60
	m := mins % 60
	if m == 0 {
		return fmt.Sprintf("%d h", h)
	}
	return fmt.Sprintf("%d:%02dm", h, m)
}

func SegundosToString(segs int) string {
	if segs == 0 {
		return "-"
	}
	if segs < 60 {
		return fmt.Sprintf("%ds", segs)
	}
	h := segs / 3600
	m := (segs % 3600) / 60
	s := segs % 60
	if h == 0 {
		return fmt.Sprintf("%dm %ds", m, s)
	}
	if m == 0 {
		return fmt.Sprintf("%dh", h)
	}
	if s == 0 {
		return fmt.Sprintf("%d:%02dm", h, m)
	}
	return fmt.Sprintf("%d:%02d:%02ds", h, m, s)
}

func (tar Tarea) TiempoEstimadoString() string {
	return MinutosToString(tar.TiempoEstimado)
}
func (tar Tarea) TiempoRealString() string {
	return SegundosToString(tar.TiempoReal)
}

func (tar IntervaloReciente) TiempoEstimadoString() string {
	return MinutosToString(tar.TiempoEstimado)
}
func (tar IntervaloReciente) TiempoRealString() string {
	return SegundosToString(tar.TiempoReal)
}

func (pry Proyecto) TiempoGestionString() string {
	return SegundosToString(pry.TiempoGestion)
}
