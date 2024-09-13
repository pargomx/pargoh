
// ================================================================ //
// ========== PLATAFORMA ========================================== //
function getPlatform() {
    // Modern way of detecting
    if (typeof navigator.userAgentData !== 'undefined' && navigator.userAgentData != null) {
        return navigator.userAgentData.platform;
    }
    // Deprecated fallback
    if (typeof navigator.platform !== 'undefined') {
        if (typeof navigator.userAgent !== 'undefined' && /android/.test(navigator.userAgent.toLowerCase())) {
            // android device’s navigator.platform is often set as 'linux', so let’s use userAgent for them
            return 'android';
        }
        return navigator.platform;
    }
    return 'unknown';
}
// let platform = getPlatform().toLowerCase();
let esMac = /mac/.test(getPlatform().toLowerCase()); // Mac desktop
// let esIOS = ['iphone', 'ipad', 'ipod'].indexOf(platform) >= 0; // Mac iOs
// let esApple = esMac || esIOS; // Apple device (desktop or iOS)
// let esWindows = /win/.test(platform); // Windows
// let esAndroid = /android/.test(platform); // Android
// let esLinux = /linux/.test(platform); // Linux



// ================================================================ //
// ========== SHORTCUTS =========================================== //

// Guardar con "CTRL + S" y hx-trigger="submit,cmdGuardar"
// Se debe prevenir default en keydown pero lanzar el evento
// en keyup para solo lanzarlo una vez.
document.addEventListener("keydown", function(e) {
	if ((esMac ? e.metaKey : e.ctrlKey) && e.code === 'KeyS') {
		e.preventDefault();
	}
}, false);
document.addEventListener("keyup", function(e) {
	if ((esMac ? e.metaKey : e.ctrlKey) && e.code === 'KeyS' ) {
		e.target.dispatchEvent(new Event("cmdGuardar", { bubbles: true }))
	}
}, false);


// ================================================================ //
// ================================================================ //


function delay(time) {
	return new Promise(resolve => setTimeout(resolve, time));
}

function quitar(element) {
	console.log("quitando")
	setTimeout(() => {
		console.log("bye")
		element.remove();
	}, 200);
}

window.onload = () => {
	// console.log("page is fully loaded");
	
	htmx.on("htmx:responseError", function(event) {
		alert(`Error ${event.detail.xhr.response}`)
	});
	htmx.on("htmx:sendError", function(event) {
		alert(`Error de red: no se puede conectar con el servidor.`)
	});
	htmx.on("htmx:timeout", function(event) {
		alert(`Error: se agotó el tiempo de espera para la respuesta del servidor.`)
	});
};




/**
 * 
 * @param {string} str string para quitar acentos, espacios, trim, diacríticos.
 * @returns string
 */
// normalizar quita acentos, espacios adicionales y mayúsculas.
function normalizar(str) {
	return str.normalize('NFD').replace(/\p{Diacritic}/gu, '').toLowerCase().trim().replace(/\s+/g, ' ')
}


// Filtrar entidades por nombre
// 
// Uso:
// <input type="search" class="" oninput="filtrarArticles(this.value)" placeholder="Buscar entidad...">
//
function filtrarArticles(qryText) {
	let cards = document.getElementsByTagName("article")
	if (cards.length < 1) {
		console.log("no hay article para filtrar")
		return
	}
	qryText = normalizar(qryText)
	for (let card of cards) {
		if ( qryText.length < 2 ) {
			card.classList.remove("hidden")
			continue
		}
		let cardText = normalizar(card.getElementsByTagName("h2")[0].textContent)
		if ( cardText.includes(qryText) ) {
			card.classList.remove("hidden")
		} else {
			card.classList.add("hidden")
		}
	}
}


function filtrarRows(qryText, tableID) {
	let rows = document.getElementById(tableID).getElementsByTagName("tbody")[0].getElementsByTagName("tr")
	if (rows.length < 1) {
		return
	}
	qryText = normalizar(qryText)
	for (let row of rows) {
		if ( qryText.length < 2 ) {
			row.classList.remove("hidden")
			continue
		}
		let rowText = normalizar(row.textContent)
		if ( rowText.includes(qryText) ) {
			row.classList.remove("hidden")
		} else {
			row.classList.add("hidden")
		}
	}
}

// Mostar checkboxes para seleccionar varias entidades.
function showCheckboxs() {
	let chbxs = document.getElementsByClassName("filtro_chkbox")
	for (let cb of chbxs) {
		cb.classList.remove("hidden")
	}
	document.getElementById("showCheckboxsBtn").classList.add("hidden")
	document.getElementById("applyCheckboxsBtn").classList.remove("hidden")
}


/**
 * Ordena una tabla alfabéticamente por una columna.
 * Fuente: https://youtu.be/8SL_hM1a0yo
 * TODO: Ordenar campos numéricos correctamente.
 * TODO: Poner campos vacíos al final.
 * 
 * @param {string} tblID El id de la tabla que se va a ordenar
 * @param {number} colIdx El index de la columna por la que ordenar
 */
function sortTableByIdAndColumn(tblID, colIdx) {
	const table = document.getElementById(tblID)
	if (table === null) {
		console.warn("sortable: la tabla ", tblID, " no existe")
		return
	}
	const tbody = table.tBodies[0]
	const rows = Array.from(tbody.querySelectorAll("tr"))

	let ordenASC = true
	if (table.querySelector(`th:nth-child(${ colIdx + 1 })`).classList.contains("th-sort-asc")) {
		ordenASC = false
	}
	
	const sortedRows = rows.sort((a, b) => {
		aColText = normalizar(a.querySelector(`td:nth-child(${ colIdx + 1 })`).textContent)
		bColText = normalizar(b.querySelector(`td:nth-child(${ colIdx + 1 })`).textContent)
		// console.log(bColText + " - " + aColText)
		
		// Opción A: usando matemáticas solamente.
		// const direccion = ordenASC ? 1 : -1
		// return aColText > bColText ? (1 * direccion) : (-1 * direccion)

		// Opción B: usando localCompare para ordenar números correctamente.
		if (ordenASC) {
			return aColText.localeCompare(bColText, undefined, { numeric: true, sensitivity: 'base' });
		} else {
			return bColText.localeCompare(aColText, undefined, { numeric: true, sensitivity: 'base' });
		}
	})

	while (tbody.firstChild) {
		tbody.removeChild(tbody.firstChild)
	}
	tbody.append(...sortedRows)

	// Agregar clase en la celda del encabezado
	table.querySelectorAll("th").forEach(th => th.classList.remove("th-sort-asc", "th-sort-desc"))
	table.querySelector(`th:nth-child(${ colIdx + 1 })`).classList.toggle("th-sort-asc", ordenASC)
	table.querySelector(`th:nth-child(${ colIdx + 1 })`).classList.toggle("th-sort-desc", !ordenASC)
}

// Que todas las tablas se puedan ordenar por las columnas que tengan el atributo "sortable" en sus <th>.
document.querySelectorAll('table').forEach(tbl => {
	if (tbl.id == "") {
		return
	}	
	for(let col of tbl.querySelectorAll("th[sortable]").entries()) {
		col[1].addEventListener("click", () => sortTableByIdAndColumn(tbl.id, col[0]) )
	}
});

// ================================================================ //
// ========== TEXTAREAS =========================================== //

// Autosize: puede tener cualquier border-width, max-height o rows.
function autosizeTextarea(textarea) {
	let style = window.getComputedStyle(textarea); // obtener y considerar border para evitar scrollbars.
    let bTop = parseFloat(style.getPropertyValue('border-top-width'));
    let bBottom = parseFloat(style.getPropertyValue('border-bottom-width'));
    let borderPx = Math.ceil(bTop + bBottom);
    textarea.setAttribute("style", "height:" + (textarea.scrollHeight + borderPx) + "px; resize: none;");
    textarea.addEventListener("input", function() {
        this.style.height = 'auto';
        this.style.height = (this.scrollHeight + borderPx) + "px";
    }, false);
}

// Aplicar autosize cuando el textarea se hace visible, no solo al cargar la página.
// Nesesario para textareas ocultos dentro de modales <dialog>.
const observer = new IntersectionObserver((entries, observer) => {
	entries.forEach(entry => {
		if (entry.isIntersecting) {
			// console.log('Textarea is visible');
			entry.target.style.backgroundColor = 'lightyellow';
			autosizeTextarea(entry.target);
			observer.unobserve(entry.target);
		}
	});
}, { threshold: 0 });

// Guardar con Enter; new line con Ctrl+Enter; cancelar con Esc;
function hdlTextAreaEnter(event) {
	if (event.key === 'Enter') {
		if (event.ctrlKey) {
			event.target.value += '\n';
			autosizeTextarea(event.target);
			event.preventDefault();
		} else {
			event.preventDefault();
			event.target.blur();
			event.target.focus(); // para que htmx sepa qué elemento enfocar after swap.
		}
	} else if (event.key === 'Escape') {
		event.target.value = event.target.defaultValue; // restaurar valor original
		event.target.blur();
	}
}

// ================================================================ //
// ========== INICIALIZAR CONTENIDO =============================== //

// Evento htmx:load para inicializar cosas después de cargar contenido.
// https://htmx.org/api/#onLoad
htmx.onLoad(function(content) {

	// console.log("htmx:onLoad", content);

	// Los input tienen autocomplete="off" a menos que se especifique lo contrario.
	content.querySelectorAll('input').forEach(input => {
		if (input.getAttribute("autocomplete") == null) {
			input.setAttribute("autocomplete", "off")
		}
	});

	// Los textarea se ajustan automáticamente a su contenido.
	const textareas = content.getElementsByTagName("textarea");
	for (let i = 0; i < textareas.length; i++) {
		autosizeTextarea(textareas[i])
		observer.observe(textareas[i]);
		textareas[i].addEventListener('keydown', hdlTextAreaEnter);
		// textareas[i].setAttribute("autocomplete", "off")
		// textareas[i].setAttribute("spellcheck", "true")
		// textareas[i].setAttribute("autocorrect", "on")
		// textareas[i].setAttribute("autocapitalize", "on")
	}

	
	// Si se declara una función "onLoad" en el contenido, se ejecuta.
	if (typeof onLoad === 'function') { 
		onLoad(content);
	}
})

// ================================================================ //
// ========== TimeTracker para gestión de proyecto ================ //

const segundosParaInactividad = 20;
const segundosParaEnviarHeartbeat = 5;
let segundosContados = 0;
let proyectoID = "_none";
let timeCounterIntvl = null;
let interactionTimeout;

// Mostrar cuántos segundos se han contado.
function setCounterDisplay(text) {
	// window.document.title = text;
	document.querySelector("footer small").innerText = text;
}

// Enviar un pulso al servidor cada x segundos.
function sendHeartbeat() {
	// Enviar segundos al servidor
	fetch(`/proyectos/${proyectoID}/time/${segundosParaEnviarHeartbeat}`, { method: 'POST' }).then(response => {
	    if (!response.ok) {
			// TODO: no dar error al usuario, pero guardar en localStorage y enviar cuando se pueda.
	        throw new Error('Network response was not ok');
	    }
	})
	// Cuenta local para mostrar al usuario.
	segundosContados += segundosParaEnviarHeartbeat;
    // localStorage.setItem('timeActive', segundosContados);
	setCounterDisplay(`🌿 ${segundosContados}s ${proyectoID}`);
}

// Contar el tiempo que se trabaja en un proyecto. Idempotente.
function startHeartbeat(razon) {
	if (!document.querySelector("[data-proyecto-id]")) {
		proyectoID = 0;
		return // Solo contar cuando se trabaja en un proyecto.
	}
	if (timeCounterIntvl) {
		console.log(razon + " [already started]");
		return // Idempotente si ya está contando.
	}
	console.log(razon);
	proyectoID = document.querySelector("[data-proyecto-id]").getAttribute("data-proyecto-id");
	// segundosContados = parseInt(localStorage.getItem('timeActive')) || segundosContados
	timeCounterIntvl = setInterval(sendHeartbeat, segundosParaEnviarHeartbeat * 1000);
	setCounterDisplay(`🌿 Start: ${segundosContados}s ${proyectoID}`);
}

// Pausar el contador de tiempo.
function stopHeartbeat(razon) {
	if (!timeCounterIntvl) {
		console.log(razon + " [already stopped]");
		return // Idempotente si ya está detenido.
    }
	console.log(razon);
	clearInterval(timeCounterIntvl);
	timeCounterIntvl = null;
	// localStorage.setItem('timeActive', segundosContados); // inecesario?
	setCounterDisplay(`⏸️ Cuenta detenida ${segundosContados}s ${proyectoID}`);
}

// Detectar cuando la pestaña está enfocada o si deja de estarlo.
function handleVisibilityChange() {
    if (document.visibilityState === 'visible' && document.hasFocus()) {
        startHeartbeat("🌳 Pestaña enfocada"); // Cubierto por UserInteraction.
    } else {
        stopHeartbeat("⚠️ Pestaña desenfocada");
		clearTimeout(interactionTimeout); // Detener el contador de inactividad.
    }
}
window.addEventListener('focus', handleVisibilityChange);
window.addEventListener('blur', handleVisibilityChange);
document.addEventListener('visibilitychange', handleVisibilityChange);

if (document.visibilityState === 'visible' && document.hasFocus()) {
    startHeartbeat("🌳 Already focused and visible when script loaded");
}

// Detectar cuando el usuario interactúa con la página.
function handleUserInteraction() {
	clearTimeout(interactionTimeout);
	interactionTimeout = setTimeout(() => {
		stopHeartbeat("⚠️ Inactividad detectada por " + segundosParaInactividad + " segundos");
    }, segundosParaInactividad * 1000);
	startHeartbeat("🌳 Actividad detectada");
}
document.onkeyup = handleUserInteraction;
document.onclick = handleUserInteraction;
document.ontouchstart = handleUserInteraction;
// window.onload = handleUserInteraction; // cubierto por VisibilityChange.
// document.onscroll = resetInteractionTimer; // meh...
// document.onmousemove = resetInteractionTimer; // demasiado sensible

// ================================================================ //