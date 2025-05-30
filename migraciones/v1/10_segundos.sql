INSERT INTO migraciones VALUES (1,10, CURRENT_TIMESTAMP, "Todos los campos de duración deben ser segundos");

ALTER TABLE tareas RENAME COLUMN tiempo_real TO segundos_real;

/* Ya no serán minutos, sino segundos */
ALTER TABLE tareas RENAME COLUMN tiempo_estimado TO segundos_estimado;
UPDATE tareas SET segundos_estimado = segundos_estimado*60;


ALTER TABLE historias RENAME COLUMN minutos_estimado TO segundos_presupuesto;
UPDATE historias SET segundos_presupuesto = segundos_presupuesto*60;