// ======================================================================================
// Obtención de Usuarios (GET)
// ======================================================================================

function obtenerEntidades() {
    const endpoint = 'http://localhost:8080/api/users';

    fetch(endpoint)
        .then(response => {
            if (!response.ok) {
                throw new Error(`Error HTTP! Estado: ${response.status}`);
            }
            return response.json();
        })
        .then(data => {
            console.log('Datos recibidos de la API:', data);

            const lista = document.getElementById('user-list');
            lista.innerHTML = ''; // Se limpia la lista previa.

            if (Array.isArray(data) && data.length > 0) {
                data.forEach(usuario => {
                    const li = document.createElement('li');
                    li.textContent = `${usuario.name} ${usuario.surname} `;

                    // Creación del botón de eliminar.
                    const btnEliminar = document.createElement('button');
                    btnEliminar.textContent = 'Eliminar';
                    btnEliminar.classList.add('delete-button');
                    btnEliminar.dataset.id = usuario.user_id;  // Se guarda el id.
                    li.appendChild(btnEliminar);

                    li.appendChild(btnEliminar);
                    lista.appendChild(li);
                });
            } else {
                lista.innerHTML = '<li>No hay usuarios registrados.</li>';
            }
        })
        .catch(error => {
            console.error('Error en la petición Fetch:', error);
            const lista = document.getElementById('user-list');
            lista.innerHTML = '<li>Error al cargar los usuarios.</li>';
        });
}

// ======================================================================================
// Eliminación de Usuario (DELETE)
// ======================================================================================

function eliminarUsuario(id) {
    const endpoint = `http://localhost:8080/api/users/${id}`;

    if (!confirm('¿Estás seguro de que deseas eliminar este usuario?')) {
        return;
    }

    fetch(endpoint, {
        method: 'DELETE'
    })
    .then(response => {
        if (!response.ok) {
            throw new Error(`Error HTTP! Estado: ${response.status}`);
        }
        console.log(`Usuario ${id} eliminado correctamente`);
        obtenerEntidades(); // Actualiza la lista
    })
    .catch(error => {
        console.error('Error al eliminar usuario:', error);
        alert('Hubo un problema al eliminar el usuario.');
    });
}

// ======================================================================================
// Envío de Nuevo Usuario (POST)
// ======================================================================================

function enviarUsuario(event) {
    event.preventDefault(); // Se evita que el formulario recargue la página.

    const endpoint = 'http://localhost:8080/api/users';
    const name = document.getElementById('name').value.trim();
    const surname = document.getElementById('surname').value.trim();

    if (!name || !surname) {
        alert('Por favor, completa ambos campos.');
        return;
    }

    const nuevoUsuario = { name, surname };

    fetch(endpoint, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(nuevoUsuario)
    })
    .then(response => {
        if (!response.ok) {
            throw new Error(`Error HTTP! Estado: ${response.status}`);
        }
        return response.json();
    })
    .then(data => {
        console.log('Usuario agregado exitosamente:', data);
        document.getElementById('signup-form').reset();
        obtenerEntidades();
    })
    .catch(error => {
        console.error('Error al enviar usuario:', error);
        alert('Hubo un problema al enviar el usuario.');
    });
}

// ======================================================================================
// Inicialización
// ======================================================================================
// --- Inicialización ---
document.addEventListener('DOMContentLoaded', () => {
    obtenerEntidades();
    document.getElementById('signup-form').addEventListener('submit', enviarUsuario);

    // Delegación.
    document.getElementById('user-list').addEventListener('click', (event) => {
        if (event.target.classList.contains('delete-button')) {
            const id = event.target.dataset.id;
            if (!id) return;
            if (confirm('¿Estás seguro de que deseas eliminar este usuario?')) {
                eliminarUsuario(id);
            }
        }
    });
});

