INSERT INTO main.migraciones VALUES (2,0, CURRENT_TIMESTAMP, 'Migrar datos de v1 a v2');

SELECT printf('     %d nodos', COUNT(*)) FROM main.nodos;

/*
  PROYECTOS:
  Convertir proyectos en nodos y poner en el grupo de los activos.
  Transformar la clave alfanumérica humana en un simple entero.
  Arreglar de una la posición haciéndola consecutiva.
*/
INSERT INTO main.nodos (
  nodo_id,
  padre_id,
  tipo,
  posicion,

  titulo,
  descripcion,
  objetivo,
  notas,

  color,
  imagen,

  prioridad,
  estatus,
  segundos,
  centavos

) SELECT /* proyectos */

  unicode(substr(proyecto_id,1))*3 * unicode(substr(proyecto_id,3))*1 + unicode(substr(proyecto_id,4))*100 - unicode(substr(proyecto_id,2)) AS nodo_id,
  2, /* grupo proyectos activos */
  'PRY',
  ROW_NUMBER() OVER (ORDER BY posicion),

  titulo,
  descripcion,
  '',
  CASE WHEN fecha_registro != '' THEN 'Creado: ' || fecha_registro ELSE '' END,

  color,
  imagen,

  0,
  0,
  0,
  0

  FROM old_schema.proyectos
  ORDER BY posicion ASC
;

SELECT printf('+%03d proyectos migrados', COUNT(*)) FROM old_schema.proyectos;
SELECT printf('     %d nodos', COUNT(*)) FROM main.nodos;

/*
  PERSONAS:
  Convertir tabla personas en nodos tipo persona.
  Usar la misma transformación para convertir el proyecto_id.
	Notar que para arreglar de una vez la posición se usa
	ROW_NUMBER OVER (PARTITION BY p.proyecto_id ORDER BY n.posicion)
	que genera un número consecutivo desde 1 para cada proyecto.
*/
INSERT INTO main.nodos (
  nodo_id,
  padre_id,
  tipo,
  posicion,

  titulo,
  descripcion,
  objetivo,
  notas,

  color,
  imagen,

  prioridad,
  estatus,
  segundos,
  centavos

) SELECT /* personas */

  p.persona_id,
  unicode(substr(proyecto_id,1))*3 * unicode(substr(proyecto_id,3))*1 + unicode(substr(proyecto_id,4))*100 - unicode(substr(proyecto_id,2)) AS nodo_id,
  'PER',
  ROW_NUMBER() OVER (PARTITION BY p.proyecto_id ORDER BY n.posicion),

  p.nombre,
  p.descripcion,
  '',
  CASE WHEN segundos_gestion != 0 THEN 'tiempo_gestion_v1: ' || segundos_gestion || 's' ELSE '' END,

  '',
  '',

  0,
  0,
  0,
  0

  FROM old_schema.personas p
  LEFT JOIN old_schema.nodos n ON n.nodo_id = p.persona_id
;

SELECT printf('+%03d personas migradas', COUNT(*)) FROM old_schema.personas;
SELECT printf('     %d nodos', COUNT(*)) FROM main.nodos;


/*
  Cambiar el historia_id(1) utilizada para QuickTasks, ya que el nodo_id(1)
  es del nodo ROOT en el nuevo esquema. El ID se actualiza automáticamente
  en la tabla tareas, que es el único lugar donde usa esta historia.
*/
UPDATE old_schema.historias SET historia_id = 9 WHERE historia_id = 1;

/*
  Historia(1) QuickTasks.
  Se pone de padre al proyecto default (nodo_id 5), ya que QuickTasks no tenía ninguno.
  El insert es igual al anterior, exepto por el 'WHERE n.nodo_id IS NULL'
  y el primer bloque donde se pone padre_id(5) y posición(1).
*/
INSERT INTO main.nodos (
  nodo_id,
  padre_id,
  tipo,
  posicion,

  titulo,
  descripcion,
  objetivo,
  notas,

  color,
  imagen,

  prioridad,
  estatus,
  segundos,
  centavos

) SELECT /* historia QuickTasks */

  h.historia_id,
  5, /* Proyecto default */
  'HIS',
  1, /* Debería ser la única historia en el proyecto default en este punto */

  h.titulo,
  h.descripcion,
  h.objetivo,
  h.notas,

  '',
  '',

  h.prioridad,
  h.completada,
  h.segundos_presupuesto,
  0

  FROM old_schema.historias h
  LEFT JOIN old_schema.nodos n ON n.nodo_id = h.historia_id
  WHERE n.nodo_id IS NULL
  ORDER BY n.nivel, n.padre_id, n.posicion
;

/*
  HISTORIAS (all):
  En esta consulta migra todas las historias menos la de QuickTasks.
*/
INSERT INTO main.nodos (
  nodo_id,
  padre_id,
  tipo,
  posicion,

  titulo,
  descripcion,
  objetivo,
  notas,

  color,
  imagen,

  prioridad,
  estatus,
  segundos,
  centavos

) SELECT /* historias (1) */

  h.historia_id,
  n.padre_id,
  'HIS',
  n.posicion,

  h.titulo,
  h.descripcion,
  h.objetivo,
  h.notas,

  '',
  '',

  h.prioridad,
  h.completada,
  h.segundos_presupuesto,
  0

  FROM old_schema.historias h
  LEFT JOIN old_schema.nodos n ON n.nodo_id = h.historia_id
  WHERE n.nodo_id IS NOT NULL
  ORDER BY n.nivel, n.padre_id, n.posicion
;

SELECT printf('+%03d historias migradas', COUNT(*)) FROM old_schema.historias;
SELECT printf('     %d nodos', COUNT(*)) FROM main.nodos;


/*
	TRAMOS:
	Convertir los tramos de viaje de usuario en un nodo más.
*/
INSERT INTO main.nodos (
  nodo_id,
  padre_id,
  tipo,
  posicion,

  titulo,
  descripcion,
  objetivo,
  notas,

  color,
  imagen,

  prioridad,
  estatus,
  segundos,
  centavos

) SELECT /* tramos */

  abs(random() %9999999999)+10000000000, /* Make unique */
  t.historia_id,
  'VIA',
  t.posicion,

  t.texto,
  '',
  '',
  '',

  '',
  t.imagen,

  0,
  0,
  0,
  0

  FROM old_schema.tramos t
;

SELECT printf('+%03d tramos migrados', COUNT(*)) FROM old_schema.tramos;
SELECT printf('     %d nodos', COUNT(*)) FROM main.nodos;

/*
	REGLAS:
	Convertir las reglas de negocio en un nodo más.
*/
INSERT INTO main.nodos (
  nodo_id,
  padre_id,
  tipo,
  posicion,

  titulo,
  descripcion,
  objetivo,
  notas,

  color,
  imagen,

  prioridad,
  estatus,
  segundos,
  centavos

) SELECT /* reglas */

  abs(random() %9999999999)+1000000000, /* Make unique */
  r.historia_id,
  'REG',
  r.posicion,

  r.texto,
  '',
  '',
  '',

  '',
  '',

  0,
  r.implementada + r.probada, /* Siempre suman 0, 1, 2 */
  0,
  0

  FROM old_schema.reglas r
;

SELECT printf('+%03d reglas migrados', COUNT(*)) FROM old_schema.reglas;
SELECT printf('     %d nodos', COUNT(*)) FROM main.nodos;

UPDATE tareas SET tarea_id = tarea_id + abs(random() %999999999);

/*
	TAREAS:
	Convertir las tareas de desarrollo en un nodo más.
  Guardar el tipo de tarea como texto.
  También guardar como texto el segundos_real que fue remplazado por los intervalos.
*/
INSERT INTO main.nodos (
  nodo_id,
  padre_id,
  tipo,
  posicion,

  titulo,
  descripcion,
  objetivo,
  notas,

  color,
  imagen,

  prioridad,
  estatus,
  segundos,
  centavos

) SELECT /* tareas */

  tarea_id,
  historia_id,
  'TAR',
  ROW_NUMBER() OVER (PARTITION BY historia_id),

  descripcion,
  impedimentos,
  tipo,
  CASE WHEN segundos_real != 0 THEN 'tiempo_v1: ' || segundos_real || 's' ELSE '' END,

  '',
  '',

  CASE WHEN importancia = 'IDEA' THEN 1 WHEN importancia = 'MEJORA' THEN 2 WHEN importancia = 'NECESARIA' THEN 3 ELSE 0 END,
  estatus,
  segundos_estimado,
  0

  FROM old_schema.tareas
;

SELECT printf('+%03d tareas migradas', COUNT(*)) FROM old_schema.tareas;
SELECT printf('     %d nodos', COUNT(*)) FROM main.nodos;

SELECT 'Migrados nodos (check sum total)';

/*
  FUERA DEL ARBOL (1): Referencias
*/
INSERT INTO main.referencias (
  nodo_id,
  ref_nodo_id
) SELECT
  historia_id,
  ref_historia_id
  FROM old_schema.referencias
;

SELECT printf('Migradas referencias %d/%d',
  (SELECT COUNT(*) FROM main.referencias),
  (SELECT COUNT(*) FROM old_schema.referencias)
);


/*
  FUERA DEL ARBOL (2): Intervalos
*/
INSERT INTO main.intervalos (
  nodo_id,
  ts_ini,
  ts_fin
) SELECT
  tarea_id,
  inicio,
  fin
  FROM old_schema.intervalos
;

SELECT printf('Migrados intervalos %d/%d',
  (SELECT COUNT(*) FROM main.intervalos),
  (SELECT COUNT(*) FROM old_schema.intervalos)
);

/*
  FUERA DEL ARBOL (3): Latidos
*/
INSERT INTO main.latidos (
  ts_latido,
  nodo_id,
  segundos
) SELECT
  timestamp,
  persona_id,
  segundos
  FROM old_schema.latidos
;

SELECT printf('Migrados latidos %d/%d',
  (SELECT COUNT(*) FROM main.latidos),
  (SELECT COUNT(*) FROM old_schema.latidos)
);

SELECT 'Done!';