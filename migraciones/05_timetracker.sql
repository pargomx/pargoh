INSERT INTO migraciones VALUES (05, CURRENT_TIMESTAMP, "Time tracker para tareas de gestión");

ALTER TABLE proyectos ADD COLUMN tiempo_gestion INT NOT NULL DEFAULT 0;
