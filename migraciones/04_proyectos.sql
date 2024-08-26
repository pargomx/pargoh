INSERT INTO migraciones VALUES (04, CURRENT_TIMESTAMP, "Reglas de negocio");

CREATE TABLE reglas (
  historia_id INT NOT NULL,
  posicion INT NOT NULL,
  texto TEXT NOT NULL,
  implementada INT NOT NULL,
  probada INT NOT NULL,
  PRIMARY KEY (historia_id,posicion),
  FOREIGN KEY (historia_id) REFERENCES historias (historia_id) ON DELETE RESTRICT ON UPDATE CASCADE
);