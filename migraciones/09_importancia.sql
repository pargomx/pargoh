INSERT INTO migraciones VALUES (09, CURRENT_TIMESTAMP, "Importancia de cada tarea");

ALTER TABLE tareas ADD COLUMN importancia TEXT NOT NULL DEFAULT '';

UPDATE tareas SET importancia = 'IDEA' WHERE importancia == '';