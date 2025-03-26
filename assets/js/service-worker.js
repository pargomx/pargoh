// ================================================================ //
// ========== CACHE =============================================== //

const ASSETS_CACHE_NAME = 'assets-v1';
const IMAGES_CACHE_NAME = 'images-v1';
const ASSETS_CACHE_PREFIX = '/assets/';
const IMAGES_CACHE_PREFIX = '/imagenes/';
const OFFLINE_PATH = '/offline';

async function initAssetsCache() {
	sendToast('initAssetsCache');
	const cache = await caches.open(ASSETS_CACHE_NAME);
	return cache.addAll([OFFLINE_PATH]);
}

async function cleanOldCaches(cacheNames) {
	return Promise.all(
		cacheNames.map((cacheName) => {
			if (cacheName !== ASSETS_CACHE_NAME && cacheName !== IMAGES_CACHE_NAME) {
				sendToast('Cleaning old cache', cacheName);
				return caches.delete(cacheName);
			}
		})
	);
}

async function validateCaches() {
	sendToast('Validando caches...');
	const assetCache = await caches.open(ASSETS_CACHE_NAME);
	const imageCache = await caches.open(IMAGES_CACHE_NAME);
	await validateSpecificCache(assetCache, ASSETS_CACHE_PREFIX);
	await validateSpecificCache(imageCache, IMAGES_CACHE_PREFIX);
}

async function validateSpecificCache(cache, pathPrefix) {
	sendToast("Validando cache", pathPrefix)
	const keys = await cache.keys();
	for (const request of keys) {
		try {
			const cachedResponse = await cache.match(request);
			if (!cachedResponse) {
				continue;
			}
			const url = new URL(request.url);
			if (!url.pathname.startsWith(pathPrefix)) {
				continue;
			}
			const fetchOptions = {
				method: 'HEAD',
			};
			const serverResponse = await fetch(request.url, fetchOptions);
			if (!serverResponse.ok) {
				if (serverResponse.status === 404) {
					console.warn(`Removing ${request.url} from cache: Server returned 404`);
					await cache.delete(request);
				} else {
					console.warn(
						`Validation failed for ${request.url}: Server returned ${serverResponse.status}`
					);
				}
				continue;
			}
			const cachedETag = cachedResponse.headers.get('ETag');
			const cachedLastModified = cachedResponse.headers.get('Last-Modified');
			const serverETag = serverResponse.headers.get('ETag');
			const serverLastModified = serverResponse.headers.get('Last-Modified');
			if (
				(cachedETag && serverETag && cachedETag !== serverETag) ||
				(cachedLastModified &&
					serverLastModified &&
					cachedLastModified !== serverLastModified)
			) {
				sendToast(`Updating cache for ${request.url}`);
				const freshResponse = await fetch(request.url);
				if (freshResponse.ok) {
					await cache.put(request, freshResponse);
					sendToast(`${request.url} updated in cache`);
				} else {
					console.warn(
						`Failed to update ${request.url} in cache. Server returned ${freshResponse.status}`
					);
				}
			} else {
				sendToast(`${request.url} is up to date`);
			}
		} catch (error) {
			console.error(`Error validating ${request.url}:`, error);
		}
	}
}

// ================================================================ //
// ========== FETCH =============================================== //

async function fetchHandler(event) {
	const url = new URL(event.request.url);
	if (url.pathname.startsWith(ASSETS_CACHE_PREFIX)) {
		return handleCacheFirstFetch(event, ASSETS_CACHE_NAME, url);

	} else if (url.pathname.startsWith(IMAGES_CACHE_PREFIX)) {
		return handleCacheFirstFetch(event, IMAGES_CACHE_NAME, url);
		
	} else {
		return handleNetworkFirstFetch(event);
	}
}

/**
 * Handles fetch events with a network-first strategy.
 * 
 * This function attempts to retrieve the requested resource from the network.
 * If the network request fails (e.g., due to being offline), it retrieves the
 * offline page resource from the cache for GET requests. For other HTTP methods,
 * it tries to fetch them normally and only if that fails then returns a 503 Service Unavailable response.
 * 
 * @param {FetchEvent} event - The fetch event to handle.
 * @returns {Promise<Response>} A promise that resolves to the network or cached response.
 */
async function handleNetworkFirstFetch(event) {
	if (event.request.method === 'GET') {
		return fetch(event.request).catch(() => {
			return caches.match(OFFLINE_PATH);
		});
	} else {
		return fetch(event.request).catch(() => {
			return new Response(`App fuera de línea`, { status: 503, statusText: 'Service Unavailable' });
		});
	}
}

/**
 * Handles fetch events with a cache-first strategy.
 * 
 * This function attempts to retrieve the requested resource from the specified cache.
 * If the resource is not found in the cache, it fetches the resource from the network,
 * caches the response, and then returns the network response.
 * 
 * @param {FetchEvent} event - The fetch event to handle.
 * @param {string} cacheName - The name of the cache to use.
 * @param {URL} url - The URL of the resource being requested.
 * @returns {Promise<Response>} A promise that resolves to the cached or fetched response.
 */
async function handleCacheFirstFetch(event, cacheName, url) {
	return caches.open(cacheName).then(async (cache) => {
		const cachedResponse = await cache.match(event.request);
		if (cachedResponse) {
			return cachedResponse;
		}
		const fetchResponse = await fetch(event.request);
		sendToast("Guardando en cache", url.pathname)
		cache.put(event.request, fetchResponse.clone());
		return fetchResponse;
	});
}

// ================================================================ //
// ========== MENSAJES ============================================ //

function handleMessage(event) {
	if (event.data === 'validateCache') {
		validateCaches();
	}
}

function sendToast(...args) {
	const message = args.map(arg => typeof arg === 'object' ? JSON.stringify(arg) : arg).join(' ');
	self.clients.matchAll().then(clients => {
		clients.forEach(client => {
			client.postMessage({ tipo: 'toast', texto: message });
		});
	});
}

// ================================================================ //
// ========== EVENTOS ============================================= //

self.addEventListener('install', (event) => {
	sendToast('ServiceWorker installing');
	event.waitUntil(initAssetsCache());
});

self.addEventListener('activate', (event) => {
	sendToast('ServiceWorker activating');
	event.waitUntil(
		caches.keys().then((cacheNames) => cleanOldCaches(cacheNames))
	);
});

self.addEventListener('fetch', (event) => {
	event.respondWith(fetchHandler(event));
});

self.addEventListener('message', handleMessage);
