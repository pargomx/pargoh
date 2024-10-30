INSERT INTO migraciones VALUES (13, CURRENT_TIMESTAMP, "Tabla para referenciar unas historias a otras");

CREATE TABLE referencias (
  historia_id INT NOT NULL,
  ref_historia_id INT NOT NULL,
  PRIMARY KEY (historia_id,ref_historia_id),
  FOREIGN KEY (historia_id) REFERENCES historias (historia_id) ON DELETE RESTRICT ON UPDATE CASCADE,
  FOREIGN KEY (ref_historia_id) REFERENCES historias (historia_id) ON DELETE RESTRICT ON UPDATE CASCADE
);
