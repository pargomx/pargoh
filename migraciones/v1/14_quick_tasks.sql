INSERT INTO migraciones VALUES (1,14, CURRENT_TIMESTAMP, "Historia padre para quick tasks");

INSERT INTO historias(historia_id, titulo, objetivo, prioridad, completada)
 VALUES (1,'QuickTasksParent','Esta historia sirve de padre para las tareas sin proyecto',0,0);
