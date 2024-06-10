INSERT INTO migraciones VALUES (01, CURRENT_TIMESTAMP, "Esquema: historias de usuario y tareas");

CREATE TABLE nodos (
  nodo_id INT NOT NULL,
  nodo_tbl TEXT NOT NULL,
  padre_id INT NOT NULL,
  padre_tbl TEXT NOT NULL,
  nivel INT NOT NULL,
  posicion INT NOT NULL,
  PRIMARY KEY (nodo_id),
  FOREIGN KEY (padre_id) REFERENCES nodos (nodo_id) ON DELETE RESTRICT ON UPDATE CASCADE
);

/* Nodo root y padre del root no es sí mismo para no interferir con posiciones del nivel 1 */
INSERT INTO nodos VALUES (-1, 'nod', -1, 'nod', 0, 0);
INSERT INTO nodos VALUES (0, 'nod', -1, 'nod', 0, 0);

CREATE TABLE personas (
  persona_id INT NOT NULL,
  nombre TEXT NOT NULL,
  descripcion TEXT NOT NULL,
  PRIMARY KEY (persona_id)
);

CREATE TABLE historias (
  historia_id INT NOT NULL,
  titulo TEXT NOT NULL,
  objetivo TEXT NOT NULL,
  prioridad INT NOT NULL,
  completada INT NOT NULL,
  PRIMARY KEY (historia_id)
);


/* Tareas */

CREATE TABLE tareas (
  tarea_id INT NOT NULL,
  historia_id INT NOT NULL,
  tipo TEXT NOT NULL,
  descripcion TEXT NOT NULL,
  impedimentos TEXT NOT NULL,
  tiempo_estimado INT NOT NULL,
  tiempo_real INT NOT NULL,
  estatus INT NOT NULL,
  PRIMARY KEY (tarea_id),
  FOREIGN KEY (historia_id) REFERENCES historias (historia_id) ON DELETE RESTRICT ON UPDATE CASCADE
);

CREATE TABLE intervalos (
  tarea_id INT NOT NULL,
  inicio TEXT NOT NULL,
  fin TEXT NOT NULL,
  PRIMARY KEY (tarea_id,inicio),
  FOREIGN KEY (tarea_id) REFERENCES tareas (tarea_id) ON DELETE RESTRICT ON UPDATE CASCADE
);