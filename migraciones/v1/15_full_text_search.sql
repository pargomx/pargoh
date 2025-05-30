INSERT INTO migraciones VALUES (1,15, CURRENT_TIMESTAMP, "Full text search virtual table & triggers");

CREATE VIRTUAL TABLE historias_fts USING fts5 (
	historia_id UNINDEXED,
	otro_id UNINDEXED,
	origen UNINDEXED,
	texto,
	tokenize = 'trigram remove_diacritics 1'
);

/*
	Historias
*/
INSERT INTO historias_fts (historia_id, otro_id, origen, texto)
	SELECT historia_id, '', 'tit', titulo FROM historias
	/* WHERE titulo != '' */
;

INSERT INTO historias_fts (historia_id, otro_id, origen, texto)
	SELECT historia_id, '', 'obj', objetivo FROM historias
	/* WHERE objetivo != '' */
;

INSERT INTO historias_fts (historia_id, otro_id, origen, texto)
	SELECT historia_id, '', 'des', descripcion FROM historias
	/* WHERE descripcion != '' */
;

CREATE TRIGGER fts_insert_historias AFTER INSERT ON historias
BEGIN
	INSERT INTO historias_fts (historia_id, otro_id, origen, texto)
		VALUES (NEW.historia_id, '', 'tit', NEW.titulo)
	;
	INSERT INTO historias_fts (historia_id, otro_id, origen, texto)
		VALUES (NEW.historia_id, '', 'obj', NEW.objetivo)
	;
	INSERT INTO historias_fts (historia_id, otro_id, origen, texto)
		VALUES (NEW.historia_id, '', 'des', NEW.descripcion)
	;
END;

CREATE TRIGGER fts_update_historias_titulo AFTER UPDATE ON historias
WHEN NEW.titulo != OLD.titulo
BEGIN
	UPDATE historias_fts
		SET texto = NEW.titulo
		WHERE historia_id = NEW.historia_id AND origen = 'tit'
	;
END;

CREATE TRIGGER fts_update_historias_objetivo AFTER UPDATE ON historias
WHEN NEW.objetivo != OLD.objetivo
BEGIN
	UPDATE historias_fts
		SET texto = NEW.objetivo
		WHERE historia_id = NEW.historia_id AND origen = 'obj'
	;
END;

CREATE TRIGGER fts_update_historias_descripcion AFTER UPDATE ON historias
WHEN NEW.descripcion != OLD.descripcion
BEGIN
	UPDATE historias_fts
		SET texto = NEW.descripcion
		WHERE historia_id = NEW.historia_id AND origen = 'des'
	;
END;

CREATE TRIGGER fts_delete_historias AFTER DELETE ON historias
BEGIN
	DELETE FROM historias_fts
		WHERE historia_id = OLD.historia_id
	;
END;


/*
	Tramos
*/
INSERT INTO historias_fts (historia_id, otro_id, origen, texto)
	SELECT historia_id, posicion, 'via', texto FROM tramos WHERE texto != ''
;

CREATE TRIGGER fts_insert_tramos AFTER INSERT ON tramos
BEGIN
	INSERT INTO historias_fts (historia_id, otro_id, origen, texto)
		VALUES (NEW.historia_id, NEW.posicion, 'via', NEW.texto)
	;
END;

CREATE TRIGGER fts_update_tramos AFTER UPDATE ON tramos
WHEN NEW.texto != OLD.texto
BEGIN
	UPDATE historias_fts
		SET texto = NEW.texto
		WHERE historia_id = NEW.historia_id AND otro_id = NEW.posicion AND origen = 'via'
	;
END;

CREATE TRIGGER fts_delete_tramos AFTER DELETE ON tramos
BEGIN
	DELETE FROM historias_fts
		WHERE historia_id = OLD.historia_id AND otro_id = OLD.posicion AND origen = 'via'
	;
END;

/*
	Reglas
*/
INSERT INTO historias_fts (historia_id, otro_id, origen, texto)
	SELECT historia_id, posicion, 'reg', texto FROM reglas WHERE texto != ''
;

CREATE TRIGGER fts_insert_reglas AFTER INSERT ON reglas
BEGIN
	INSERT INTO historias_fts (historia_id, otro_id, origen, texto)
		VALUES (NEW.historia_id, NEW.posicion, 'reg', NEW.texto)
	;
END;

CREATE TRIGGER fts_update_reglas AFTER UPDATE ON reglas
WHEN NEW.texto != OLD.texto
BEGIN
	UPDATE historias_fts
		SET texto = NEW.texto
		WHERE historia_id = NEW.historia_id AND otro_id = NEW.posicion AND origen = 'reg'
	;
END;

CREATE TRIGGER fts_delete_reglas AFTER DELETE ON reglas
BEGIN
	DELETE FROM historias_fts
		WHERE historia_id = OLD.historia_id AND otro_id = OLD.posicion AND origen = 'reg'
	;
END;



/*
	Tareas
*/
INSERT INTO historias_fts (historia_id, otro_id, origen, texto)
	SELECT historia_id, tarea_id, 'tar', descripcion FROM tareas WHERE descripcion != ''
;

INSERT INTO historias_fts (historia_id, otro_id, origen, texto)
	SELECT historia_id, tarea_id, 'imp', impedimentos FROM tareas WHERE impedimentos != ''
;

CREATE TRIGGER fts_insert_tareas AFTER INSERT ON tareas
BEGIN
	INSERT INTO historias_fts (historia_id, otro_id, origen, texto)
		VALUES (NEW.historia_id, NEW.tarea_id, 'tar', NEW.descripcion)
	;
	INSERT INTO historias_fts (historia_id, otro_id, origen, texto)
		VALUES (NEW.historia_id, NEW.tarea_id, 'imp', NEW.impedimentos)
	;
END;

CREATE TRIGGER fts_update_tareas_descripcion AFTER UPDATE ON tareas
WHEN NEW.descripcion != OLD.descripcion
BEGIN
	UPDATE historias_fts
		SET texto = NEW.descripcion
		WHERE historia_id = NEW.historia_id AND otro_id = NEW.tarea_id AND origen = 'tar'
	;
END;

CREATE TRIGGER fts_update_tareas_impedimentos AFTER UPDATE ON tareas
WHEN NEW.impedimentos != OLD.impedimentos
BEGIN
	UPDATE historias_fts
		SET texto = NEW.impedimentos
		WHERE historia_id = NEW.historia_id AND otro_id = NEW.tarea_id AND origen = 'imp'
	;
END;

CREATE TRIGGER fts_delete_tareas AFTER DELETE ON tareas
BEGIN
	DELETE FROM historias_fts
		WHERE historia_id = OLD.historia_id AND otro_id = OLD.tarea_id AND origen IN('tar','imp')
	;
END;



/* 
	Probar triggers
*/
INSERT INTO historias(historia_id, titulo, objetivo, prioridad, completada)
	VALUES (2, 'Test historia','Historia de prueba para triggers de FTS',0,0);
INSERT INTO tramos(historia_id, posicion, texto, imagen)
	VALUES (2, 3,'test', '');
INSERT INTO reglas(historia_id, posicion, texto, implementada, probada)
	VALUES (2, 3,'test', 0, 0);
INSERT INTO tareas(historia_id, tarea_id, descripcion, tipo, impedimentos,segundos_estimado,segundos_real,estatus)
	VALUES (2, 3,'test', '','',0,0,'');

UPDATE historias SET titulo = 'Historia updated' WHERE historia_id = 2;
UPDATE historias SET objetivo = 'Historia updated' WHERE historia_id = 2;
UPDATE historias SET descripcion = 'Historia updated' WHERE historia_id = 2;

UPDATE tramos SET texto = 'updated' WHERE historia_id=2 AND posicion=3;
UPDATE reglas SET texto = 'updated' WHERE historia_id=2 AND posicion=3;

UPDATE tareas SET descripcion = 'updated' WHERE tarea_id=3;
UPDATE tareas SET impedimentos = 'updated' WHERE tarea_id=3;

DELETE FROM tramos WHERE historia_id=2 AND posicion=3;
DELETE FROM reglas WHERE historia_id=2 AND posicion=3;
DELETE FROM tareas WHERE tarea_id=3;
DELETE FROM historias WHERE historia_id = 2;
