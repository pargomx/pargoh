package arbol

const (
	EstatusTareaNoIniciada = 0
	EstatusTareaEnCurso    = 1
	EstatusTareaEnPausa    = 2
	EstatusTareaFinalizada = 3
)

func (tar Tarea) NoIniciada() bool {
	return tar.Estatus == EstatusTareaNoIniciada
}
func (tar Tarea) EnCurso() bool {
	return tar.Estatus == EstatusTareaEnCurso
}
func (tar Tarea) EnPausa() bool {
	return tar.Estatus == EstatusTareaEnPausa
}
func (tar Tarea) Finalizada() bool {
	return tar.Estatus == EstatusTareaFinalizada
}
