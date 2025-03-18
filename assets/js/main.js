// Ejecutado con defer luego de html pero antes de DOMContentLoaded.
// logTimeTotal("main.js started")

/**
 * Elemento con el contenido dentro de layout.html
 * Guarda scrollPosition.
 */
let contenido = document.getElementById('contenido');

// ================================================================ //
// ========== SHORTCUTS GLOBALES ================================== //

/**
 * Navegar hacia arriba en el árbol de historias con CTRL + UP.
 * @param {KeyboardEvent} event keydown
 */
function handleShortcutAncestro(event) {
	if ((event.altKey || event.ctrlKey) && event.key === 'ArrowUp') {
        event.preventDefault();
		let ancestroDirectoLink = document.querySelector('a[tipo="ancestro_directo"]');
		if (ancestroDirectoLink) {
			ancestroDirectoLink.click();
		}
    }
}

// ================================================================ //
// ========== INPUT HELPERS ======================================= //

/**
 * Eliminar acentos, diacríticos, espacios adicionales y mayúsculas de un texto.
 * @param {string} str string para quitar acentos, espacios, trim, diacríticos.
 * @returns string
 */
function normalizar(str) {
	return str.normalize('NFD').replace(/\p{Diacritic}/gu, '').toLowerCase().trim().replace(/\s+/g, ' ')
}

/**
 * Hacer click en botón de submit del formulario.
 * @param {HTMLElement} element Formulario o elemento dentro de un formulario.
 */
function clickSubmit(element) {
	if (element.tagName !== "FORM") {
		element = element.closest("form");
	}
	const submitButton = element.querySelector('button[type="submit"]');
	if (submitButton) {
		submitButton.click();
	}
}

/**
 * Preparar todos los inputs que estén dentro del elemento dado.
 * 
 * Set autocomplete="off" a menos que se especifique lo contrario.
 * @param {HTMLElement} element
 */
function prepararInputsEn(element) {
	element.querySelectorAll('input').forEach(input => {
		if (input.getAttribute("autocomplete") == null) {
			input.setAttribute("autocomplete", "off")
		}
	});
}

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

// ================================================================ //
// ========== TEXTAREAS =========================================== //

/**
 * Agregar un salto de línea en la posición actual del cursor.
 * @param {HTMLTextAreaElement} textarea 
 */
function textareaAddNewLineAtCursor(textarea) {
	if (textarea.selectionStart || textarea.selectionStart == 0) {
		let startPos = textarea.selectionStart;
		let endPos = textarea.selectionEnd;
		textarea.value = textarea.value.substring(0, startPos)
			+ '\n'
			+ textarea.value.substring(endPos, textarea.value.length);
		textarea.selectionStart = startPos+1;
		textarea.selectionEnd = endPos+1;
	} else {
		textarea.value += '\n';
	}
}

/**
 * Guardar con Enter; new line con Shift+Enter; cancelar con Esc;
 * @param {KeyboardEvent} event KeyDown
 */
function textareaHandleKeyDown(event) {
	let textarea = event.target
	if (!(textarea instanceof HTMLTextAreaElement)) {
		return
	}
	if (event.key === 'Enter') {
		event.preventDefault()
		if (event.shiftKey) {
			textareaAddNewLineAtCursor(textarea)
			textareaAutosize(textarea)
		} else {
			textarea.blur()  // trigger change para enviar con htmx.
			textarea.focus() // para que htmx sepa qué elemento enfocar after swap.
		}
	} else if (event.key === 'Escape') {
		textarea.value = textarea.defaultValue // restaurar valor original
		textarea.blur()
	}
}

/**
 * En Chrome con Gboard no se envía event.shiftKey en el keyDown event. Esto es
 * un hack para permitir introducir saltos de línea desde android mediante la
 * introducción compuesta de la palabra "break", con swipe por ejemplo.
 * @param {InputEvent} event 
 */
function textareaHandleBeforeInput(event) {
	let textarea = event.target
	if (!(textarea instanceof HTMLTextAreaElement)) {
		return
	}
	if (event.data && event.data.trim().toLowerCase() === 'break') {
		event.preventDefault();
		textareaAddNewLineAtCursor(textarea);
		textareaAutosize(textarea);
	}
}

// ================================================================ //

/**
 * Autosize: adaptar altura del cuadro de texto al contenido.
 * Toma en cuenta cualquier border-width, max-height o rows.
 * @param {HTMLTextAreaElement} textarea 
 */
function textareaAutosize(textarea) {
	if (!(textarea instanceof HTMLTextAreaElement)) {
		console.warn(`autosizeTextarea: elemento no es un textarea`, textarea)
		return
	}
	const oldH = parseInt(textarea.style.height) || 0;
	const newH = textarea.scrollHeight + textareaGetBorderHeight(textarea)
	// En chrome android dentro de dialogs el textarea.scrollHeight incrementa
	// 1px con cada caracter introducido. Ignorar esos incrementos minúsculos.
	if (Math.abs(oldH - newH) <= 2) {
		return
	}
	textarea.style.height = newH + "px";
}

/**
 * Obtener altura del border para evitar scrollbars en autosizeTextarea.
 * @param {HTMLTextAreaElement} textarea 
 */
function textareaGetBorderHeight(textarea) {
	let style = window.getComputedStyle(textarea);
    let bTop = parseFloat(style.getPropertyValue('border-top-width'));
    let bBottom = parseFloat(style.getPropertyValue('border-bottom-width'));
    let borderPx = Math.ceil(bTop + bBottom); // nunca quedarse cortos en px
	return borderPx
}

/**
 * Hacky way de no saltar en textareas muy grandes al editar y al pegar texto.
 * @param {InputEvent} event 
 */
function textareaHandleInputSinScrollJumps(event) {
	textareaAutosize(event.target)
	
	// TODO: NO RECUERDO CUÁL ERA LA DIFERENCIA :V
	// console.log("scrollTop1", contenido.scrollTop)
	// let scrollTop = contenido.scrollTop;
	// console.log("scrollTop3", document.getElementById('contenido').scrollTop)
	// console.log("scrollTop2", contenido.scrollTop)
	// contenido.scrollTop = scrollTop

	// let scrollTopRestore = null;
	// if (contenido) {
	// 	scrollTopRestore = contenido.scrollTop;
	// }
	// if (scrollTopRestore) {
	// 	contenido.scrollTop = scrollTopRestore;
	// }
	// console.log("scrollTop2", contenido.scrollTop)
}

/**
 * Aplicar autosize cuando el textarea se hace visible, no solo al cargar contenido.
 * Nesesario para textareas que comienzan ocultos, como dentro de modales <dialog>
 * ya que al inicio no se puede obtener su altura porque están ocultos.
 */
const textareasObserver = new IntersectionObserver((entries, observer) => {
	entries.forEach(entry => {
		if (entry.isIntersecting) {
			textareaAutosize(entry.target);
			observer.unobserve(entry.target);
		}
	});
}, { threshold: 0 });

// ================================================================ //

/**
 * Preparar todas las textareas que estén dentro del elemento dado.
 * @param {HTMLElement} element
 */
function prepararTextareasEn(element) {
	const textareas = element.getElementsByTagName("textarea");
	if (textareas.length == 0){
		return
	}
	for (let i = 0; i < textareas.length; i++) {
		// Ajustar al contenido manualmente, en espera de soporte para field-sizing.
		textareas[i].classList.add('resize-none');
		textareaAutosize(textareas[i]);
		textareasObserver.observe(textareas[i]);
		textareas[i].addEventListener("input", textareaHandleInputSinScrollJumps);

		// Usar [Enter] para enviar y [Shift+Enter] para nueva línea.
		textareas[i].addEventListener('keydown', textareaHandleKeyDown);
		textareas[i].addEventListener('beforeinput', textareaHandleBeforeInput);
		
		// textareas[i].setAttribute("autocomplete", "off")
		// textareas[i].setAttribute("spellcheck", "true")
		// textareas[i].setAttribute("autocorrect", "on")
		// textareas[i].setAttribute("autocapitalize", "on")
	}
	if (textareas.length == 1){
		logTimeHasta(`textarea preparada`, textareas[0])
	} else {
		logTimeHasta(`textareas preparadas`, [...textareas])
	}
}

// ================================================================ //
// ========== SCROLL POSITION ===================================== //

/**
 * Guardar scroll position antes de navegar a otra página con o sin HTMX.
 * 
 * body - htmx:beforeSwap | window - beforeunload
 */
function guardarScrollPosition() {
    if (!contenido) {
		return
	}
	let scrollTop = contenido.scrollTop;
	contenido.dataset.scrollPosition = scrollTop; // para htmx:beforeSwap
	let scrollData = JSON.parse(localStorage.getItem('scrollData')) || []; // Para navegador nativo
	let currentURL = window.location.href;
	scrollData = scrollData.filter(item => item.url !== currentURL); // Quitar si ya existe en el arreglo
	scrollData.unshift({ url: currentURL, scrollTop: scrollTop }); // Agregar al inicio del arreglo
	if (scrollData.length > 5) {
		scrollData.pop(); // Recordar hasta 5 páginas
	}
	localStorage.setItem('scrollData', JSON.stringify(scrollData));
}

/**
 * Restaurar scroll position al navegar con o sin HTMX.
 * Invocar después de textareaAutosize.
 */
function restoreScrollPosition() {
    // No restablecer si no es una página de aplicación normal.
	if (!contenido) {
		return
	}
	// Si el atributo scrollPosition existe significa que se navegó con HTMX.
	if (contenido.dataset.scrollPosition) {
		if (contenido.dataset.scrollPosition == contenido.scrollTop){
			return
		}
		logTimeHasta(`scrollRestored (htmx) from ${contenido.scrollTop} to ${contenido.dataset.scrollPosition}`)
		contenido.scrollTop = contenido.dataset.scrollPosition;
		return
	}
	// De lo contrario se navegó a otra página y se busca en localStorage.
	let currentURL = window.location.href;
	let scrollData = JSON.parse(localStorage.getItem('scrollData')) || [];
	let savedPos = scrollData.find(item => item.url === currentURL);
	if (savedPos) {
		logTimeHasta(`scrollRestored (native) from ${contenido.scrollTop} to ${savedPos.scrollTop}`)
		contenido.scrollTop = savedPos.scrollTop;
	}
}

// Esperar a que se haga el restoreScrollPosition para no mostrar
// al usuario el movimiento y que el contenido haga flash.
function mostrarContenidoListo() {
	if (!contenido) {
		return
	}
	contenido.classList.remove('opacity-0', 'pointer-events-none')
}

/**
 * Resalta con un border el fragmento de la URL para ubicarlo fácilmente.
 */
function resaltarFragmento() {
	if (window.location.hash) {
		let elemento = document.getElementById(window.location.hash.substr(1));
		if (elemento != null) {
			elemento.classList.add("border-indigo-600");
			elemento.classList.add("border-2");
			elemento.scrollIntoView();
		}
	}
}

// ================================================================ //
// ========== TimeTracker para gestión de proyecto ================ //

const segundosParaInactividad = 20;
const segundosParaEnviarHeartbeat = 5;
let segundosContados = 0;
let timePersonaID = "0";
let timeCounterIntvl = null;
let interactionTimeout;

// Enviar un pulso de actividad al servidor.
function sendHeartbeat() {
	fetch(`/personas/${timePersonaID}/time/${segundosParaEnviarHeartbeat}`, { method: 'POST' }).then(response => {
	    if (!response.ok) {
			// TODO: no dar error al usuario, pero guardar en localStorage y enviar cuando se pueda.
			// localStorage.setItem('timeActive', segundosContados);
	        throw new Error('Network response was not ok');
	    }
	})
	segundosContados += segundosParaEnviarHeartbeat; // Cuenta local para mostrar al usuario.
}

// Continuar contando el tiempo que se trabaja en un personaje.
function startHeartbeat(razon) {
	if (!document.querySelector("[data-persona-id]")) {
		timePersonaID = 0;
		return // Solo contar cuando se trabaja con un personaje.
	}
	if (timeCounterIntvl) {
		return // Puede ya estar contando.
	}
	timePersonaID = document.querySelector("[data-persona-id]").getAttribute("data-persona-id");
	timeCounterIntvl = setInterval(sendHeartbeat, segundosParaEnviarHeartbeat * 1000);
}

// Pausar el contador de tiempo.
function stopHeartbeat(razon) {
	if (!timeCounterIntvl) {
		return // Puede ya estar detenido.
    }
	clearInterval(timeCounterIntvl);
	timeCounterIntvl = null;
}

// Reiniciar timeout de inactividad cuando el usuario interactúa con la página.
// Si se pasa de los segundosParaInactividad entonces se detiene el heartbeat.
function keepHeartbeatAfterUserInteraction() {
	clearTimeout(interactionTimeout);
	interactionTimeout = setTimeout(() => {
		stopHeartbeat("⚠️ Inactividad " + segundosParaInactividad + " segundos");
    }, segundosParaInactividad * 1000);
	startHeartbeat("🌳 Actividad detectada");
}

// Detectar cuando la pestaña se enfoca y se desenfoca para pausar inmediatamente.
function handleVisibilityChange() {
    if (document.visibilityState === 'visible' && document.hasFocus()) {
        startHeartbeat("🌳 Pestaña enfocada"); // Cubierto por keepHeartbeatAfterUserInteraction.
    } else {
        stopHeartbeat("⚠️ Pestaña desenfocada");
		clearTimeout(interactionTimeout); // Detener el contador de inactividad.
    }
}


document.addEventListener('keyup', keepHeartbeatAfterUserInteraction);
document.addEventListener('click', keepHeartbeatAfterUserInteraction);
document.addEventListener('touchstart', keepHeartbeatAfterUserInteraction);
// window.onload = keepHeartbeatAfterUserInteraction; // cubierto por handleVisibilityChange.
// document.onscroll = keepHeartbeatAfterUserInteraction; // meh...
// document.onmousemove = keepHeartbeatAfterUserInteraction; // demasiado sensible
window.addEventListener('focus', handleVisibilityChange);
window.addEventListener('blur', handleVisibilityChange);
document.addEventListener('visibilitychange', handleVisibilityChange);

// Iniciar heartbeat al ejecutar script.
if (document.visibilityState === 'visible' && document.hasFocus()) {
    startHeartbeat("🌳 Already focused and visible when script loaded");
}

// ================================================================ //
// ========== Eventos HTMX ======================================== //

document.body.addEventListener('htmx:responseError', function(event) {
	alert(`Error ${event.detail.xhr.response}`)
});

document.body.addEventListener('htmx:sendError', function(event) {
	alert(`Error de red: no se puede conectar con el servidor.`)
});

document.body.addEventListener('htmx:timeout', function(event) {
	alert(`Error: se agotó el tiempo de espera para la respuesta del servidor.`)
});

// Debug perfomrance
document.body.addEventListener('htmx:configRequest', function(event) {
	logTimeSince('htmx:configRequest')
});

// Guardar scroll position antes de swap.
document.body.addEventListener('htmx:beforeSwap', function() {
	guardarScrollPosition()
});

/**
 * Para que HTMX no repita el setup de los elementos en la primera
 * carga de la página, sino que lo haga cuando se hace un swap.
 */
let primeraCarga = true

// Cuando htmx termina de cargar contenido (event.detail.elt)
document.body.addEventListener('htmx:load', (event) => {
	if (primeraCarga) {
		return // No inicializar doble
	}
	logTimeHasta("htmx:load (main.js)", event.detail.elt)
	prepararInputsEn(event.detail.elt)
	prepararTextareasEn(event.detail.elt)
	restoreScrollPosition()
});

// ================================================================ //
// ========== INICIALIZAR PÁGINA ================================== //

document.addEventListener('keydown', handleShortcutAncestro);

window.addEventListener('beforeunload', guardarScrollPosition);

prepararInputsEn(document.body)
prepararTextareasEn(document.body)
restoreScrollPosition()

// Cuando se termine la carga de recursos incluyendo imágenes y estilos.
window.addEventListener("load", (event) => {
	logTimeTotal("window:load (main.js)")
	mostrarContenidoListo()
	primeraCarga = false
});

// El navegador puede ignorar "Cache-Control: no-store" y activar
// el BFCache al navegar adelante o atrás. Forzar la recarga lo evita.
window.addEventListener('pageshow', (event) => {
	if (event.persisted) {
		//? console.log("show page persisted in BFCache aborted")
		location.reload();
	}
});

logTimeTotal("main.js finished")