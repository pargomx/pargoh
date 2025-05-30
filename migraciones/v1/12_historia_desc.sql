INSERT INTO migraciones VALUES (1,12, CURRENT_TIMESTAMP, "Agregar descripción para historias dirigida a usuario");

ALTER TABLE historias ADD COLUMN descripcion TEXT NOT NULL DEFAULT '';