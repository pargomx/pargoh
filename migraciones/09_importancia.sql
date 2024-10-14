INSERT INTO migraciones VALUES (09, CURRENT_TIMESTAMP, "Importancia de cada tarea");

ALTER TABLE tareas ADD COLUMN importancia INT NOT NULL DEFAULT 0;