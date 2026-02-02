// logTimeTotal("historia.js started")

// ================================================================ //
// ========== REGLAS ============================================== //

function initSortableReglas() {
	let sortableElement = document.querySelector("[tipo='sort_reglas']")
	if (!sortableElement) {
		return
	}
	let sortableReglas = new Sortable(sortableElement, {
		group: "sort_reglas", animation: 150, swapThreshold: 0.50,
		draggable: "[tipo='sort_regla']",
		handle: "[tipo='sort_handle']",
		ghostClass: "opacity-50",
		onEnd: function (event) {
			if (event.oldIndex == event.newIndex){ return }
			event.item.querySelector("input[name='new_pos']").value = event.newIndex + 1;
			event.item.querySelector("form[tipo='sort_form']").dispatchEvent(new Event("reordenarEnd"));
		},
	});
	logTimeHasta("initSortableReglas", sortableElement)
}

// ================================================================ //
// ========== VIAJE =============================================== //

function initSortableTramos() {
	let sortableElement = document.querySelector("[tipo='sort_tramos']")
	if (!sortableElement) {
		return
	}
	let sortableTramos = new Sortable(sortableElement, {
		group: "sort_tramos", animation: 150, swapThreshold: 0.50,
		draggable: "[tipo='sort_tramo']",
		handle: "[tipo='sort_handle']",
		ghostClass: "opacity-50",
		onEnd: function (event) {
			if (event.oldIndex == event.newIndex){ return }
			event.item.querySelector("input[name='new_pos']").value = event.newIndex + 1;
			event.item.querySelector("form[tipo='sort_form']").dispatchEvent(new Event("reordenarEnd"));
		},
	});
	logTimeHasta("initSortableTramos", sortableElement)
}

function showMoverTramo(tramoID, posicion, descripcion) {
	console.log(`showMoverTramo(${tramoID},${posicion})`);
	let mov = document.getElementById('moverTramo');
	mov.showModal();
	mov.querySelector('input[name=nodo_id]').value = tramoID;
	mov.querySelector('input[name=posicion]').value = posicion;
	mov.querySelector('p').innerText = descripcion;
	document.getElementById('navTreeTramo').dispatchEvent(new Event('navTreeTramoShown'));
}

// ================================================================ //

/**
 * Capturar imagen con CTRL + V para subir al tramo.
 * 
 * `window.addEventListener('paste', handlePasteImagen)`
 * @param {ClipboardEvent} event
 */
function handlePasteImagen(event) {
	let tramoImgModal = document.getElementById("TramoImgModal");
	if (!tramoImgModal.open) {
		return
	}
	const imgFileInput = tramoImgModal.querySelector("input[type='file']")
	const items = (event.clipboardData || event.originalEvent.clipboardData).items

	// console.log("Pegando", tramoImgModal.open, imgFileInput)

	let isGoogleImage = false
	let isImageFile = false
	
	// Averiguar qué hacer.
	for (const item of items) {
		// item.getAsString(url => { console.log(`ITEM kind:${item.kind} type:${item.type}`, url) })
		if (item.type.startsWith('image/')) {
			isImageFile = true

		} else if (item.type === 'application/x-vnd.google-docs-image-clip+wrapped') {
			isGoogleImage = true // Imágenes desde google docs.
		}
	}

	if (isImageFile) {
		// TODO: solo aceptar un archivo y que sea imagen.
		imgFileInput.files = event.clipboardData.files
		imgFileInput.dispatchEvent(new Event('change'));

	} else if (isGoogleImage) {
		for (const item of items) {
			if (item.type === 'text/html') {
				item.getAsString(htmlString => {
					// Todo se hace aquí adentro porque getAsString es async y hay que esperarla.
					try {
						const imgUrl = getGoogleusercontentImageUrl(htmlString)
						fetchImageFromURL(imgUrl).then(file => {
							const dataTransfer = new DataTransfer();
							dataTransfer.items.add(file);
							imgFileInput.files = dataTransfer.files;
							imgFileInput.dispatchEvent(new Event('change'));
						});
					} catch (error) {
						console.error(error)
					}

				})
			}
		}
	}
}

function getGoogleusercontentImageUrl(htmlString) {
	if (htmlString == "") {
		throw new Error('Invalid HTML: Empty string');
	}
	const parser = new DOMParser();
	const doc = parser.parseFromString(htmlString, 'text/html');
	const imgTag = doc.querySelector('img');
	if (!imgTag) {
		throw new Error('Invalid HTML: Must contain an image tag');
	}
	const imgUrl = imgTag.getAttribute('src');
	const regex = /^https?:\/\/([a-zA-Z0-9-]+\.)*googleusercontent\.com/;
	if (!regex.test(imgUrl)) {
		throw new Error('Invalid URL: Must be an image hosted on googleusercontent.com');
	}
	return imgUrl;
}

async function fetchImageFromURL(imgUrl) {
	if (imgUrl == "") {
		throw new Error('fetchImageFromURL: empty url');
	}
	console.log("Fetching image:", imgUrl)
	const response = await fetch(imgUrl);
	const blob = await response.blob();
	const file = new File([blob], 'image.jpg', { type: blob.type });
	return file;
}

function showTramoImgModal(liElement) {
	// Obtener datos del tramo seleccionado.
	let historiaID= liElement.querySelector("input[name='historia_id']").value
	let posicion  =	liElement.querySelector("input[name='old_pos']").value
	let texto     =	liElement.querySelector("textarea").value
	let imgElem   =	liElement.querySelector("img")
	let imgSrc    =	imgElem ? imgElem.src : ""
	
	// Poner datos en el dialog.
	let dialog = document.getElementById("TramoImgModal")
	let btnEliminarImg = dialog.querySelector("button[tipo='eliminar_img']")
	dialog.querySelector("h2").textContent = `Viaje de usuario (${posicion})`
	dialog.querySelector("p").textContent = texto
	dialog.querySelector("input[name='historia_id']").value = historiaID
	dialog.querySelector("input[name='posicion']").value = posicion
	dialog.querySelector("input[type='file']").value = ""
	dialog.querySelector("button[type='submit']").textContent = imgSrc ? "Remplazar" : "Subir imagen"
	if (imgSrc != "") {
		dialog.querySelector("img").src = imgSrc
		dialog.querySelector("img").classList.remove("hidden")
		btnEliminarImg.setAttribute("hx-delete", "/imagenes/" + historiaID + "/" + posicion)
		btnEliminarImg.classList.remove("hidden")
		dialog.querySelector("button[tipo='cancelar']").classList.add('hidden')
	} else {
		dialog.querySelector("img").src = ""
		dialog.querySelector("img").classList.add("hidden")
		btnEliminarImg.removeAttribute("hx-delete")
		btnEliminarImg.classList.add("hidden")
		dialog.querySelector("button[tipo='cancelar']").classList.remove('hidden')
	}
	htmx.process(btnEliminarImg)
	
	// Botón para ir a tramo anterior
	prevTramo = liElement.previousElementSibling
	let btnPrev = dialog.querySelector("button[title='Anterior']")
	if (prevTramo && prevTramo.hasAttribute('tipo') && prevTramo.getAttribute('tipo') === 'sort_tramo') {
		btnPrev.classList.remove('hidden');
		btnPrev.onclick = function() { showTramoImgModal(prevTramo); }
		dialog.addEventListener('keydown', showPrevTramoHdlrKey);
	} else {
		btnPrev.classList.add('hidden');
		btnPrev.onclick = function() { console.log("imposible") }
		dialog.removeEventListener('keydown', showPrevTramoHdlrKey);
	}
	// Botón para ir a tramo siguiente
	nextTramo = liElement.nextElementSibling
	let btnNext = dialog.querySelector("button[title='Siguiente']")
	if (nextTramo && nextTramo.hasAttribute('tipo') && nextTramo.getAttribute('tipo') === 'sort_tramo') {
		btnNext.classList.remove('hidden');
		btnNext.onclick = function() { showTramoImgModal(nextTramo); }
		dialog.addEventListener('keydown', showNextTramoHdlrKey);
	} else {
		btnNext.classList.add('hidden');
		btnNext.onclick = function() { console.log("imposible") }
		dialog.removeEventListener('keydown', showNextTramoHdlrKey);
	}
	// Mostrar dialog antes de hacer focus.
	dialog.showModal()
	dialog.querySelector("button[title='Cerrar']").focus()
}

// Handlers como funciones utilizando globales para poder removerlos.
function showPrevTramoHdlrKey(event) {
	if (event.key === 'ArrowLeft') { showTramoImgModal(prevTramo); }
}
function showNextTramoHdlrKey(event) {
	if (event.key === 'ArrowRight') { showTramoImgModal(nextTramo); }
}
let prevTramo = null
let nextTramo = null

// Mostrar imagen cuando se selecciona un archivo.		
function showTramoImagePreview() {
	let dialog = document.getElementById("TramoImgModal")
	const [file] = dialog.querySelector("input[type='file']").files
	if (file) {
		// console.log('Archivo: ' + file.name);
		dialog.querySelector("img").src = URL.createObjectURL(file)
		dialog.querySelector("img").classList.remove("hidden")
		dialog.querySelector("button[tipo='eliminar_img']").classList.add('hidden')
		dialog.querySelector("button[tipo='cancelar']").classList.remove("hidden")
		dialog.querySelector("button[type='submit']").focus()
	}
}

// ================================================================ //
// ========== DESCENDIENTES ======================================= //

function initSortableDescendientes() {
	let sortableElement = document.querySelector("[tipo='sort_descendientes']")
	if (!sortableElement) {
		return
	}
	let sortableDescend = new Sortable(sortableElement, {
		group: "sort_descendientes", animation: 150, swapThreshold: 0.50,
		draggable: "[tipo='sort_descendiente']",
		handle: "[tipo='sort_handle']",
		ghostClass: "opacity-50",
		onEnd: function (event) {
			if (event.oldIndex == event.newIndex){ return }
			event.item.querySelector("input[name='new_pos']").value = event.newIndex + 1;
			event.item.querySelector("form[tipo='sort_form']").dispatchEvent(new Event("reordenarEnd"));
		},
	});
	logTimeHasta("initSortableDescendientes", sortableElement)
}

/**
 * Mostrar context menu de una historia descendiente.
 * @param {PointerEvent} event 
 * @param {number} historiaID 
 */
function showMenuHistoria(event, historiaID) {
	let menu = document.getElementById('ctxMenuHistoria' + historiaID);
	let menuW = menu.offsetWidth // Dimensiones del menú
	let menuH = menu.offsetHeight
	let winW = window.innerWidth // Dimensiones de la ventana
	let winH = window.innerHeight
	let x = event.clientX // Coordenadas dentro de la ventana
	let y = event.clientY
	if (x + menuW > winW) { x = winW - menuW } // Que no salga de la ventana
	if (y + menuH > winH) { y = winH - menuH }
	menu.style.left = `${x}px`
	menu.style.top = `${y}px`
	menu.showModal();
}

// Cerrar context menu al hacer click en él.
document.querySelectorAll("dialog[tipo='ctxMenuHistoria']").forEach((menu) => {
	menu.addEventListener('click', (event) => menu.close());
});

// Mostrar dialog para mover historia descendiente.
function showMoverHistoria(historiaID, descripcion) {
	console.log(`showMoverHistoria(${historiaID})`);
	let mov = document.getElementById('moverHistoria');
	mov.showModal();
	mov.querySelector('input[name=nodo_id]').value = historiaID;
	mov.querySelector('p').innerText = descripcion;
	document.getElementById('navTreeHistoria').dispatchEvent(new Event('navTreeHistoriaShown'));
}

// ================================================================ //
// ========== RELACIONADAS ======================================== //

function showAddReferencia() {
	let mov = document.getElementById('dialogAddReferencia');
	mov.showModal();
	document.getElementById('navAddReferencia').dispatchEvent(new Event('navAddReferenciaShown'));
}

// ================================================================ //
// ========== NOTAS DE IMPLEMENTACIÓN ============================= //

/**
 * Mostrar textarea para editar notas de implementación.
 * @param {HTMLButtonElement} btn 
 */
function editNotasImplementacion(btn) {
	const editNotasDiv = document.getElementById('notasImplEdit')
	editNotasDiv.classList.remove('hidden');
	document.getElementById('notasImplView').classList.add('hidden');
	btn.classList.add('hidden');
	const textarea = editNotasDiv.querySelector('textarea');
	textarea.focus();
	textarea.setSelectionRange(textarea.value.length, textarea.value.length);
	document.getElementById('viewNotasBtn').classList.remove('hidden');
}


// ================================================================ //
// ========== TAREAS ============================================== //

function showMoverTarea(tareaID, descripcion) {
	let mov = document.getElementById('moverTareaDialog');
	mov.showModal();
	mov.querySelector('input[name=nodo_id]').value = tareaID;
	mov.querySelector('p').innerText = descripcion;
}

/**
 * Mostrar y cargar dialog con intervalos de trabajo de una tarea.
 * @param {number} tareaID 
 */
function showIntervalos(tareaID) {
	let intervalosList = document.getElementById("intervalosList")
	intervalosList.setAttribute("hx-get", "/h/" + tareaID)
	htmx.process(intervalosList)
	intervalosList.dispatchEvent(new Event("cargarIntervalos"))
	document.getElementById("modalIntervalos").showModal()
}

// Recargar intervalos luego de modificar alguno.
document.getElementById("tareasList").addEventListener("htmx:afterSwap", (event) => {
	if (document.getElementById("modalIntervalos").open) {
		document.getElementById('intervalosList').dispatchEvent(new Event('cargarIntervalos'));
	}
})

// ================================================================ //
// ========== TIME TRACKER ======================================== //

let uri = 'ws:'; if (window.location.protocol === 'https:') { uri = 'wss:'; }; uri += '//' + window.location.host + `/h/${historiaID}/ws`;
const ws = new WebSocket(uri)
let wsID = 0
ws.onopen = function() {
	// console.log('Connected to ' + uri);
}
ws.onclose = function(event) {
	console.log("WebSocket closed:", event);
};
ws.onerror = function(error) {
	console.log("WebSocket error:", error);
};
ws.onmessage = function(event) {
	let data = JSON.parse(event.data)
	if (data.id > 0 && wsID == 0) {
		wsID = data.id
		// console.log("WebsocketID: " + wsID)

	} else if (data.reload) {
		document.querySelector("main").dispatchEvent(new Event("reloadHistoria"))
		// console.log("Recargar historia")
		
	} else {
		// console.log("Mensaje recibido: ", data)
	}
};

// Enviar wsID con solicitudes de escritura.
document.body.addEventListener('htmx:configRequest', function(event) {
	if (event.detail.verb != 'get'){
		event.detail.headers['X-SocketID'] = wsID;
		// console.log("Send wsID: " + wsID + " " + event.detail.verb)
	}
});

// ================================================================ //
// ========== INICIALIZAR ========================================= //

window.addEventListener('paste', handlePasteImagen)

document.addEventListener('DOMContentLoaded', (event) => {
	initSortableReglas()
	initSortableTramos()
	initSortableDescendientes()
});


document.body.addEventListener('htmx:load', (event) => {
	if (primeraCarga) {
		return // No inicializar doble
	}

	if (event.detail.elt.id === 'reglasDeNegocio') {
		initSortableReglas()
	}
	if (event.detail.elt.id === 'viajeList') {
		initSortableTramos()
	}
	if (event.detail.elt.id === 'descendientesList') {
		initSortableDescendientes()
	}
});

logTimeTotal("historia.js finished")