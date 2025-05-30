INSERT INTO migraciones VALUES (1,02, CURRENT_TIMESTAMP, "Viajes de usuario");

CREATE TABLE tramos (
  historia_id INT NOT NULL,
  posicion INT NOT NULL,
  texto TEXT NOT NULL,
  imagen TEXT NOT NULL,
  PRIMARY KEY (historia_id, posicion),
  FOREIGN KEY (historia_id) REFERENCES historias (historia_id) ON DELETE RESTRICT ON UPDATE CASCADE
);