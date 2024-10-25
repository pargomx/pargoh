INSERT INTO migraciones VALUES (11, CURRENT_TIMESTAMP, "Asociar timetracker a persona en lugar de proyecto");

/*
	Agregar campo segundos_gestion a personas.
*/ 

ALTER TABLE personas RENAME TO personas_old;

CREATE TABLE personas (
  persona_id INT NOT NULL,
  proyecto_id TEXT NOT NULL,
  nombre TEXT NOT NULL,
  descripcion TEXT NOT NULL,
  segundos_gestion INT NOT NULL,
  PRIMARY KEY (persona_id),
  FOREIGN KEY (proyecto_id) REFERENCES proyectos (proyecto_id) ON DELETE RESTRICT ON UPDATE CASCADE
);

INSERT INTO personas SELECT persona_id, proyecto_id, nombre, descripcion, 0 FROM personas_old;



/* 
	Cambiar referencia en latidos de proyecto a persona.
	El tiempo registrado se mueve a la primera persona del proyecto.
*/
ALTER TABLE latidos RENAME TO latidos_old;

CREATE TABLE latidos (
  timestamp TEXT NOT NULL,
  persona_id INT NOT NULL,
  segundos INT NOT NULL,
  PRIMARY KEY (timestamp),
  FOREIGN KEY (persona_id) REFERENCES personas (persona_id) ON DELETE RESTRICT ON UPDATE CASCADE
);

INSERT INTO latidos (timestamp, persona_id, segundos)
SELECT
  latidos_old.timestamp,
  (SELECT personas.persona_id FROM personas JOIN nodos ON nodo_id = personas.persona_id WHERE personas.proyecto_id = latidos_old.proyecto_id ORDER BY nodos.posicion LIMIT 1) as persona_id,
  latidos_old.segundos
  FROM latidos_old;


DROP TABLE latidos_old;

DROP TABLE personas_old;

/*
	Agregar latido inicial con los segundos registrados antes de crear la tabla latidos.
	Se asocia igual a la primera persona del proyecto.
	Se suma el ROWID a la fecha para respetar el Unique del PK timestamp.
*/
INSERT INTO latidos
SELECT
	DATETIME(UNIXEPOCH('2024-10-01 12:00:00')+ROWID, 'unixepoch') AS timestmp,
	(SELECT personas.persona_id FROM personas JOIN nodos ON nodo_id = personas.persona_id WHERE personas.proyecto_id = proyectos.proyecto_id ORDER BY nodos.posicion LIMIT 1) as persona_id,
	tiempo_gestion - (SELECT SUM(segundos) FROM latidos WHERE latidos.persona_id = (SELECT personas.persona_id FROM personas JOIN nodos ON nodo_id = personas.persona_id WHERE personas.proyecto_id = proyectos.proyecto_id ORDER BY nodos.posicion LIMIT 1) GROUP BY persona_id) AS latido
FROM proyectos
WHERE latido IS NOT NULL;

/*
	Actualizar campo materializado con la suma de todos los latidos de manera
	que ahora ya coincida el número.
*/
UPDATE personas SET segundos_gestion = coalesce((SELECT SUM(segundos) FROM latidos WHERE latidos.persona_id = personas.persona_id GROUP BY persona_id), 0);

/* 
	Quitar campo ahora obsoleto "proyectos.tiempo_gestion"
*/
ALTER TABLE proyectos DROP COLUMN tiempo_gestion;