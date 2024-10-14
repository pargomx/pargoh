package ust_test

import (
	"monorepo/ust"
	"testing"
)

func TestNuevaDuración(t *testing.T) {
	t.Parallel()

	tests := []struct {
		txt  string
		want int // en minutos
	}{
		{"", 0},
		{"15", 15},
		{"90", 90},
		{"90m", 90},
		{"90 m", 90},
		{"90 min", 90},
		{"1h", 60},
		{"1h15", 75},
		{"1h 15 min", 75},
		{"2:30", 150},
		{"4h:1", 241},
		{"01:90", 150},
	}

	for _, tt := range tests {
		t.Run(tt.txt, func(t *testing.T) {
			got, err := ust.NuevaDuraciónSegundos(tt.txt)
			if err != nil {
				t.Fatalf("NuevaDuración() error = %v", err)
			}
			if got != tt.want*60 {
				t.Errorf("NuevaDuración() = %v, want %v", got, tt.want)
			}
		})
	}
}
