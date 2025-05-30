INSERT INTO migraciones VALUES (1,17, CURRENT_TIMESTAMP, "Agregar notas de implementación a historias");

ALTER TABLE historias ADD notas TEXT NOT NULL DEFAULT '';
