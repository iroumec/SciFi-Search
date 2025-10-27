// ======================================================================================
// Configuración Global
// ======================================================================================

const API_BASE = 'http://localhost:8080/api';
const USERS_ENDPOINT = `${API_BASE}/users`;

// ======================================================================================
// Funciones Utilitarias
// ======================================================================================

async function apiFetch(url, options = {}) {
    const response = await fetch(url, options);
    if (!response.ok) throw new Error(`Error HTTP ${response.status}`);
    // Si la respuesta tiene contenido JSON, se parsea; si no, se devuelve null.
    return response.status !== 204 ? response.json() : null;
}

// --------------------------------------------------------------------------------------

function mostrarError(mensaje, destino = 'user-list') {
    console.error(mensaje);
    const contenedor = document.getElementById(destino);
    if (contenedor) contenedor.innerHTML = `<li>${mensaje}</li>`;
}

// ======================================================================================
// Obtención de Usuarios (GET)
// ======================================================================================

async function obtenerEntidades() {
    const lista = document.getElementById('user-list');

    try {
        const data = await apiFetch(USERS_ENDPOINT);
        console.log('Datos recibidos de la API:', data);

        lista.innerHTML = ''; // Se vacía la lista.

        if (Array.isArray(data) && data.length > 0) {
            for (const { name, surname, user_id } of data) {
                const li = document.createElement('li');
                li.textContent = `${name} ${surname} `;

                // Creación del botón de eliminar.
                const btnEliminar = document.createElement('button');
                btnEliminar.textContent = 'Eliminar';
                btnEliminar.classList.add('delete-button');
                btnEliminar.dataset.id = user_id; // Guardado del ID.

                li.appendChild(btnEliminar);
                lista.appendChild(li);
            }
        } else {
            lista.innerHTML = '<li>No hay usuarios registrados.</li>';
        }

    } catch (error) {
        mostrarError('Error al cargar los usuarios.');
    }
}

// ======================================================================================
// Eliminación de Usuario (DELETE)
// ======================================================================================

async function eliminarUsuario(id) {
    try {
        await apiFetch(`${USERS_ENDPOINT}/${id}`, { method: 'DELETE' });
        console.log(`Usuario ${id} eliminado correctamente`);
        await obtenerEntidades();
    } catch (error) {
        console.error('Error al eliminar usuario:', error);
        alert('Hubo un problema al eliminar el usuario.');
    }
}

// ======================================================================================
// Envío de Nuevo Usuario (POST)
// ======================================================================================

async function enviarUsuario(event) {
    event.preventDefault();

    const form = event.target;
    const { name, surname } = Object.fromEntries(new FormData(form));

    if (!name.trim() || !surname.trim()) {
        alert('Por favor, completa ambos campos.');
        return;
    }

    try {
        const nuevoUsuario = { name: name.trim(), surname: surname.trim() };
        const data = await apiFetch(USERS_ENDPOINT, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(nuevoUsuario)
        });

        console.log('Usuario agregado exitosamente:', data);
        form.reset();
        await obtenerEntidades();
    } catch (error) {
        console.error('Error al enviar usuario:', error);
        alert('Hubo un problema al enviar el usuario.');
    }
}

// ======================================================================================
// Búsqueda
// ======================================================================================

const input = document.getElementById('first-query');
input.addEventListener('keypress', (e) => {
    if (e.key === 'Enter') {
        const query = input.value.trim();
        if (!query) return;
        window.location.href = `/search?query=${encodeURIComponent(query)}`;
    }
});

// ======================================================================================
// Inicialización
// ======================================================================================

(async () => {
    await obtenerEntidades();

    document.getElementById('signup-form')
        .addEventListener('submit', enviarUsuario);

    document.getElementById('user-list')
        .addEventListener('click', async (e) => {
            if (e.target.matches('.delete-button')) {
                const id = e.target.dataset.id;
                if (id && confirm('¿Está seguro de que desea eliminar este usuario?')) {
                    await eliminarUsuario(id);
                }
            }
        });
})();
