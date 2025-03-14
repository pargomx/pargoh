
// Force reload when BFCache activates despite Cache-Control: no-store.
window.addEventListener('pageshow', (event) => {
	if (event.persisted) {
		console.log("persisted aborted")
		location.reload();
	}
});



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

// Preparar todas las textareas al cargar el contenido especificado.
function prepararTextareasEn(contenido) {
	const textareas = contenido.getElementsByTagName("textarea");
	for (let i = 0; i < textareas.length; i++) {
		// Ajustar al contenido manualmente, en espera de soporte para field-sizing.
		textareas[i].classList.add('resize-none');
		autosizeTextarea(textareas[i]);
		observer.observe(textareas[i]);
		textareas[i].addEventListener("input", handleTextareaInputWithoutJumps);

		// Usar [Enter] para enviar y [Shift+Enter] para nueva línea.
		textareas[i].addEventListener('keydown', hdlTextAreaKeyDown);
		textareas[i].addEventListener('beforeinput', hdlTextAreaBeforeInput);
		
		// textareas[i].setAttribute("autocomplete", "off")
		// textareas[i].setAttribute("spellcheck", "true")
		// textareas[i].setAttribute("autocorrect", "on")
		// textareas[i].setAttribute("autocapitalize", "on")
	}
}

// Aplicar autosize cuando el textarea se hace visible, no solo al cargar contenido.
// Nesesario para textareas que comienzan ocultos, como dentro de modales <dialog>
// ya que al inicio no se puede obtener su altura (están ocultos).
const observer = new IntersectionObserver((entries, observer) => {
	entries.forEach(entry => {
		if (entry.isIntersecting) {
			autosizeTextarea(entry.target);
			observer.unobserve(entry.target);
		}
	});
}, { threshold: 0 });

// ================================================================ //

// Autosize: obtener y considerar altura del border para evitar scrollbars.
function getTextareaBorderHeight(textarea) {
	let style = window.getComputedStyle(textarea);
    let bTop = parseFloat(style.getPropertyValue('border-top-width'));
    let bBottom = parseFloat(style.getPropertyValue('border-bottom-width'));
    let borderPx = Math.ceil(bTop + bBottom); // nunca quedarse cortos en px
	return borderPx
}

// Autosize: puede tener cualquier border-width, max-height o rows.
function autosizeTextarea(textarea) {
	if (!(textarea instanceof HTMLTextAreaElement)) {
		console.warn(`autosizeTextarea: elemento no es un textarea`, textarea)
		return
	}
	// En chrome android dentro de dialogs el textarea.scrollHeight incrementa
	// 1px con cada caracter introducido. Ignorar esos incrementos minúsculos.
	const oldH = parseInt(textarea.style.height) || 0;
	const newH = textarea.scrollHeight + getTextareaBorderHeight(textarea)
	if (Math.abs(oldH - newH) <= 2) {
		return
	}
	textarea.style.height = newH + "px";
}

// Hacky way de no saltar en textareas muy grandes al editar y al pegar texto.
function handleTextareaInputWithoutJumps(event) {
	autosizeTextarea(event.target)
	let scrollTopRestore = null;
	if (document.getElementById('contenido')) {
		scrollTopRestore = document.getElementById('contenido').scrollTop;
	}
	if (scrollTopRestore) {
		contenido.scrollTop = scrollTopRestore;
	}
}

// ================================================================ //

// Agregar un salto de línea en la posición actual del cursor.
function addNewLineAtCursor(textarea) {
	if (textarea.selectionStart || textarea.selectionStart == '0') {
		var startPos = textarea.selectionStart;
		var endPos = textarea.selectionEnd;
		textarea.value = textarea.value.substring(0, startPos)
			+ '\n'
			+ textarea.value.substring(endPos, textarea.value.length);
		textarea.selectionStart = startPos+1;
		textarea.selectionEnd = endPos+1;
	} else {
		textarea.value += '\n';
	}
}

// Guardar con Enter; new line con Shift+Enter; cancelar con Esc;
function hdlTextAreaKeyDown(event) {
	if (event.key === 'Enter') {
		if (event.shiftKey) {
			event.preventDefault();
			addNewLineAtCursor(event.target);
			autosizeTextarea(event.target);
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

// En Chrome con Gboard no se envía event.shiftKey en el keyDown event. Esto es
// un hack para permitir introducir saltos de línea desde android mediante la
// introducción compuesta de la palabra "break", con swipe por ejemplo.
function hdlTextAreaBeforeInput(event) {
	if (event.data && event.data.trim().toLowerCase() === 'break') {
		event.preventDefault();
		addNewLineAtCursor(event.target);
		autosizeTextarea(event.target);
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

	prepararTextareasEn(content)
	
	// Restablecer scroll position después de autosizeTextarea, ya que HTMX
	// provoca layout shift al hacer scrollIntoView sin considerar autosize.
	var contenido = document.getElementById('contenido');
	if (contenido && contenido.dataset.scrollPosition) {
		contenido.scrollTop = contenido.dataset.scrollPosition;
		// console.log("Restored scroll: " + contenido.dataset.scrollPosition)
	}

	// Si se declara una función "onLoad" en el contenido, se ejecuta.
	if (typeof onLoad === 'function') { 
		onLoad(content);
	}
})

// Guardar scroll position antes de swap.
document.addEventListener('htmx:beforeSwap', function(evt) {
    var contenido = document.getElementById('contenido');
    if (contenido) {
      contenido.dataset.scrollPosition = contenido.scrollTop;
	//   console.log("Saved scroll: " + contenido.scrollTop)
    }
});

// Guardar scroll position antes de navegar a otra página.
window.addEventListener('beforeunload', function() {
    var contenido = document.getElementById('contenido');
    if (contenido) {
        var scrollTop = contenido.scrollTop;
        var currentUrl = window.location.href;
        // Get the stored data
        var scrollData = JSON.parse(localStorage.getItem('scrollData')) || [];
        // Remove the current URL if it already exists in the array
        scrollData = scrollData.filter(item => item.url !== currentUrl);
        // Add the current URL and scroll position to the beginning of the array
        scrollData.unshift({ url: currentUrl, scrollTop: scrollTop });
        // Keep only the last 5 entries
        if (scrollData.length > 5) {
            scrollData.pop();
        }
        // Save the updated data back to localStorage
        localStorage.setItem('scrollData', JSON.stringify(scrollData));
    }
});

// Restore scroll position al cargar sin HTMX.
window.addEventListener('load', function() {
    var contenido = document.getElementById('contenido');
    if (contenido) {
        var currentUrl = window.location.href;
        // Get the stored data
        var scrollData = JSON.parse(localStorage.getItem('scrollData')) || [];
        // Find the scroll position for the current URL
        var scrollItem = scrollData.find(item => item.url === currentUrl);
        // Do de scroll
		if (scrollItem) {
            contenido.scrollTop = scrollItem.scrollTop;
			// console.log("scrolled "+scrollItem.scrollTop + " for "+ scrollItem.url)
        }
    }
});

// ================================================================ //
// ========== TimeTracker para gestión de proyecto ================ //

const segundosParaInactividad = 20;
const segundosParaEnviarHeartbeat = 5;
let segundosContados = 0;
let timePersonaID = "0";
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
	fetch(`/personas/${timePersonaID}/time/${segundosParaEnviarHeartbeat}`, { method: 'POST' }).then(response => {
	    if (!response.ok) {
			// TODO: no dar error al usuario, pero guardar en localStorage y enviar cuando se pueda.
	        throw new Error('Network response was not ok');
	    }
	})
	// Cuenta local para mostrar al usuario.
	segundosContados += segundosParaEnviarHeartbeat;
    // localStorage.setItem('timeActive', segundosContados);
	// setCounterDisplay(`🌿 ${segundosContados}s ${timePersonaID}`);
}

// Contar el tiempo que se trabaja en un personaje. Idempotente.
function startHeartbeat(razon) {
	if (!document.querySelector("[data-persona-id]")) {
		timePersonaID = 0;
		return // Solo contar cuando se trabaja con un personaje.
	}
	if (timeCounterIntvl) {
		// console.log(razon + " [already started]");
		return // Idempotente si ya está contando.
	}
	// console.log(razon);
	timePersonaID = document.querySelector("[data-persona-id]").getAttribute("data-persona-id");
	// segundosContados = parseInt(localStorage.getItem('timeActive')) || segundosContados
	timeCounterIntvl = setInterval(sendHeartbeat, segundosParaEnviarHeartbeat * 1000);
	// setCounterDisplay(`🌿 Start: ${segundosContados}s ${timePersonaID}`);
}

// Pausar el contador de tiempo.
function stopHeartbeat(razon) {
	if (!timeCounterIntvl) {
		// console.log(razon + " [already stopped]");
		return // Idempotente si ya está detenido.
    }
	// console.log(razon);
	clearInterval(timeCounterIntvl);
	timeCounterIntvl = null;
	// localStorage.setItem('timeActive', segundosContados); // inecesario?
	// setCounterDisplay(`⏸️ Cuenta detenida ${segundosContados}s ${timePersonaID}`);
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

function clickSubmit(form) {
	if (form.tagName !== "FORM") {
		form = form.closest("form");
	}
	form.querySelector('button[type="submit"]').click();
}

// ================================================================ //

// Navegar en el árbol de historias con CTRL + UP
document.addEventListener('keydown', function(event) {
    if ((event.altKey || event.ctrlKey) && event.key === 'ArrowUp') {
        event.preventDefault();
		let ancestroDirectoLink = document.querySelector('a[tipo="ancestro_directo"]');
		if (ancestroDirectoLink) {
			ancestroDirectoLink.click();
		}
    }
});

// ================================================================ //
// ========== DATE TIME INIPUT ==================================== //

const regIsDigit = /^[0-9]$/;
const regIsDatetime = /^[0-9]{4}-[0-9]{2}-[0-9]{2} [0-9]{2}:[0-9]{2}:[0-9]{2}$/;
function isDigit(str) { return regIsDigit.test(str); }
function isDatetime(str) { return regIsDatetime.test(str); }

/**
 * Reemplazar un digito a la vez para solo tener que usar
 * números al introducir fecha y hora.
 * 
 * Uso: hx-trigger="keyup[key=='Enter']"
 *		onkeydown="facilitarInputDatetime(event)"
 * 
 * @param {KeyboardEvent} event The keyboard event object.
 */
function facilitarInputDatetime(event) {
	const key = event.key; // The key that was pressed
	const inputField = event.target; // The input field element
	const value = inputField.value; // The current value of the input field
	let start = inputField.selectionStart; // The cursor position

	if (!isDatetime(value)) {
		// showToast("Formato de fecha incorrecto")
		return
	}
	// showToast("'"+key+"'")
	
	// Sustituir número en la posición actual.
	if (isDigit(key)) {
		event.preventDefault();

		// Ensure cursor is within a numeric character
		if (start < value.length && isDigit(value[start])) {
			// Replace the digit at the cursor position
			const newValue = value.substring(0, start) + key + value.substring(start + 1);
			inputField.value = newValue;
			inputField.selectionStart = start + 1;
			inputField.selectionEnd = start + 1;

			// If the next char is not a number, jump to the next number
			if (start + 1 < newValue.length && !isDigit(newValue[start + 1])) {
				let nextNumStart = start + 1;
				while (
					nextNumStart < newValue.length &&
					!isDigit(newValue[nextNumStart])
				) {
					nextNumStart++;
				}
				inputField.selectionStart = nextNumStart;
				inputField.selectionEnd = nextNumStart;
			}
		} else {
			// If not on a number, move to the next number or do nothing
			let nextNumStart = start;
			while (nextNumStart < value.length && !isDigit(value[nextNumStart])) {
				nextNumStart++;
			}
			inputField.selectionStart = nextNumStart;
			inputField.selectionEnd = nextNumStart;
		}
	
	// Jump to next number
	} else if (key === ' ' || key === 'ArrowRight') {
		event.preventDefault();
		let nextNumStart = start +1;
		while (nextNumStart < value.length && !isDigit(value[nextNumStart])) {
			nextNumStart++;
		}
		nextNumStart = nextNumStart > value.length ? value.length : nextNumStart
		inputField.selectionStart = nextNumStart;
		inputField.selectionEnd = nextNumStart;
	
	// Jump to previous number
	} else if (key === 'ArrowLeft') {
		event.preventDefault();
		let prevNumPos = start - 1;
		while (prevNumPos >= 0 && !isDigit(value[prevNumPos])) {
			prevNumPos--;
		}
		prevNumPos = prevNumPos < 0 ? 0 : prevNumPos
		inputField.selectionStart = prevNumPos;
		inputField.selectionEnd = prevNumPos;

	// Replace previous number with "0"
	} else if (key === 'Backspace') {
		event.preventDefault();
		let prevNumPos = start - 1;
		while (prevNumPos >= 0 && !isDigit(value[prevNumPos])) {
			prevNumPos--;
		}
		if (prevNumPos >= 0) {
			const newValue = value.substring(0, prevNumPos) + '0' + value.substring(prevNumPos + 1);
			inputField.value = newValue;
			inputField.selectionStart = prevNumPos; // Move cursor back to replaced digit
			inputField.selectionEnd = prevNumPos;
		}

	// Replace the current number with "0"
	} else if (key === 'Delete') {
		event.preventDefault();
		if (start < value.length && isDigit(value[start])) {
			const newValue = value.substring(0, start) + '0' + value.substring(start + 1);
			inputField.value = newValue;
			inputField.selectionStart = start + 1;
			inputField.selectionEnd = start + 1;

			// If the next char is not a number, jump to the next number
			if (start + 1 < newValue.length && !isDigit(newValue[start + 1])) {
				let nextNumStart = start + 1;
				while (nextNumStart < newValue.length && !isDigit(newValue[nextNumStart])) {
					nextNumStart++;
				}
				inputField.selectionStart = nextNumStart;
				inputField.selectionEnd = nextNumStart;
			}
		}

	}
}