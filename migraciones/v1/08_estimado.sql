INSERT INTO migraciones VALUES (1,08, CURRENT_TIMESTAMP, "Estimación de costo para cada historia");

ALTER TABLE historias ADD COLUMN minutos_estimado INT NOT NULL DEFAULT 0;