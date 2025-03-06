INSERT INTO migraciones VALUES (16, CURRENT_TIMESTAMP, "Agregar posicion, color, fecha_registro a proyectos");

ALTER TABLE proyectos ADD posicion INT NOT NULL DEFAULT 1;
ALTER TABLE proyectos ADD color TEXT NOT NULL DEFAULT '';
ALTER TABLE proyectos ADD fecha_registro TEXT NOT NULL DEFAULT '';

