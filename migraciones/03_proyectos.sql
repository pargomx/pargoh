INSERT INTO migraciones VALUES (03, CURRENT_TIMESTAMP, "Proyectos");

CREATE TABLE proyectos (
  proyecto_id TEXT NOT NULL,
  titulo TEXT NOT NULL,
  descripcion TEXT NOT NULL,
  imagen TEXT NOT NULL,
  PRIMARY KEY (proyecto_id)
);

INSERT INTO proyectos VALUES ("default","Proyecto default","","");

PRAGMA foreign_keys = OFF;

CREATE TABLE personas_new (
  persona_id INT NOT NULL,
  proyecto_id TEXT NOT NULL,
  nombre TEXT NOT NULL,
  descripcion TEXT NOT NULL,
  PRIMARY KEY (persona_id,proyecto_id),
  FOREIGN KEY (proyecto_id) REFERENCES proyectos (proyecto_id) ON DELETE RESTRICT ON UPDATE CASCADE
);

INSERT INTO personas_new SELECT persona_id, "default", nombre, descripcion FROM personas;

DROP TABLE personas;

ALTER TABLE personas_new RENAME TO personas;

PRAGMA foreign_keys = ON;
