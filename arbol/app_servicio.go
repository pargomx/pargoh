package arbol

// ================================================================ //
// ========== Servicio ============================================ //

type Servicio struct {
	cfg Config
}

type Config struct {
	ImagesDir string
}

func NuevoServicio(cfg Config) (*Servicio, error) {
	// if cfg.Repo == nil {
	// 	return nil, gko.ErrNoDisponible.Str("NuevoServicio: repo es nil")
	// }

	return &Servicio{
		cfg: cfg,
	}, nil
}
