INSERT INTO migraciones VALUES (1,06, CURRENT_TIMESTAMP, "Materializar proyecto_id y persona_id en tabla historias");

ALTER TABLE historias ADD COLUMN persona_id INT NOT NULL DEFAULT 0;
ALTER TABLE historias ADD COLUMN proyecto_id TEXT NOT NULL DEFAULT '';
