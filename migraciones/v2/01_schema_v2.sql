INSERT INTO migraciones VALUES (2,1, CURRENT_TIMESTAMP, 'Esquema completo v2');

CREATE TABLE nodos (
  nodo_id INT NOT NULL,
  padre_id INT NOT NULL,
  tipo TEXT NOT NULL,
  posicion INT NOT NULL,

  titulo TEXT NOT NULL,
  descripcion TEXT NOT NULL DEFAULT '',
  objetivo TEXT NOT NULL DEFAULT '',
  notas TEXT NOT NULL DEFAULT '',

  color TEXT NOT NULL DEFAULT '',
  imagen TEXT NOT NULL DEFAULT '',

  prioridad INT NOT NULL DEFAULT 0,
  estatus INT NOT NULL DEFAULT 0,
  segundos INT NOT NULL DEFAULT 0,
  centavos INT NOT NULL DEFAULT 0,

  PRIMARY KEY (nodo_id),
  FOREIGN KEY (padre_id) REFERENCES nodos (nodo_id) ON DELETE RESTRICT ON UPDATE CASCADE
);

/* 
	Crear nodo root y gupos primigeneos.
*/
INSERT INTO nodos (nodo_id, padre_id, posicion, tipo, titulo) VALUES
( 1, 1, 1, 'ROOT', 'Nodo raíz'),
( 2, 1, 1, 'GRP', 'Proyectos activos'),
( 3, 1, 2, 'GRP', 'Proyectos archivados'),
( 5, 3, 1, 'PRY', 'Proyecto default')
;


/*
  Referencias no directas entre un nodo y otro.
*/
CREATE TABLE referencias (
  nodo_id INT NOT NULL,
  ref_nodo_id INT NOT NULL,

  PRIMARY KEY (nodo_id, ref_nodo_id),
  FOREIGN KEY (nodo_id) REFERENCES nodos (nodo_id) ON DELETE RESTRICT ON UPDATE CASCADE,
  FOREIGN KEY (ref_nodo_id) REFERENCES nodos (nodo_id) ON DELETE RESTRICT ON UPDATE CASCADE
);


/*
  Registro del tiempo
*/
CREATE TABLE intervalos (
  nodo_id INT NOT NULL,
  ts_ini TEXT NOT NULL,
  ts_fin TEXT NOT NULL,

  PRIMARY KEY (nodo_id, ts_ini),
  FOREIGN KEY (nodo_id) REFERENCES nodos (nodo_id) ON DELETE RESTRICT ON UPDATE CASCADE
);

CREATE TABLE latidos (
  ts_latido TEXT NOT NULL,
  nodo_id INT NOT NULL,
  segundos INT NOT NULL,

  PRIMARY KEY (ts_latido),
  FOREIGN KEY (nodo_id) REFERENCES nodos (nodo_id) ON DELETE RESTRICT ON UPDATE CASCADE
);
