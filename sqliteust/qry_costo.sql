WITH RECURSIVE arbol(historia_id, padre_id, nivel, posicion, titulo, prioridad, completada, segundos_estimado, segundos_real) AS (

	SELECT h.historia_id, n.padre_id, n.nivel, n.posicion, h.titulo, h.prioridad, h.completada,
		coalesce((SELECT SUM(t.segundos_estimado) FROM tareas t WHERE t.historia_id = n.nodo_id),0),
		coalesce((SELECT SUM(t.segundos_real) FROM tareas t WHERE t.historia_id = n.nodo_id),0)
	FROM nodos n
	JOIN historias h ON n.nodo_id = h.historia_id
	WHERE n.padre_id = ?

	UNION ALL

	SELECT n.nodo_id, n.padre_id, n.nivel, n.posicion, h.titulo, h.prioridad, h.completada,
		coalesce((SELECT SUM(t.segundos_estimado) FROM tareas t WHERE t.historia_id = n.nodo_id),0),
		coalesce((SELECT SUM(t.segundos_real) FROM tareas t WHERE t.historia_id = n.nodo_id),0)
	FROM arbol a
	JOIN nodos n ON n.padre_id = a.historia_id
	JOIN historias h ON h.historia_id = n.nodo_id

) SELECT * FROM arbol ORDER BY nivel, padre_id, posicion

/*
) SELECT prioridad, completada, count(*) FROM arbol GROUP BY prioridad, completada;
*/