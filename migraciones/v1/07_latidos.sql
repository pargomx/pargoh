INSERT INTO migraciones VALUES (1,07, CURRENT_TIMESTAMP, "Time tracker de gestión con latidos");

CREATE TABLE latidos (
  timestamp TEXT NOT NULL,
  segundos INT NOT NULL,
  proyecto_id TEXT NOT NULL,
  historia_id INT,
  PRIMARY KEY (timestamp),
  FOREIGN KEY (proyecto_id) REFERENCES proyectos (proyecto_id) ON DELETE RESTRICT ON UPDATE CASCADE,
  FOREIGN KEY (historia_id) REFERENCES historias (historia_id) ON DELETE RESTRICT ON UPDATE CASCADE
);