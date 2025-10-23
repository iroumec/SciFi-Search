/* DEfino la comunicacion con la API */
function obtenerEntidades() {
    const endpoint = '/users';


//Hago FETCH para la peticion del GET
    fetch(endpoint)
        .then(response => {
            // Verifica si la respuesta fue exitosa (código 200-299)
            if (!response.ok) {
                // Si la respuesta no es OK, lanza un error para ser capturado por el catch
                throw new Error(`Error HTTP! Estado: ${response.status}`);
            }
            // Convierte la respuesta a formato JSON
            return response.json();
        })

        .then(data => {
            //maneja los datos recibidios
            console.log('Datos recibidos de la API', data);
            
            //como mostrar los datos en el DOM. Filtro mediante el ID
            const container = document.getElementById('data-container');
            container.innerHTML = 'Datos cargados exitosamente. REvisar la consola para observar el JSON completo';

        })
        .catch(error => {
            //handler para cualquier tipo de error durante la peticion o procesamiento
            console.error('HUbo un problema con la peticion Fetch', error);
            document.getElementById('data-container').innerHTML = 'Error al cargar los datos.';
        });
}

document.addEventListener ('DOMContentLoaded', obtenerEntidades);
//EL DOMContentLoaded espera a que el HTML este completamente cargado antes de intentear ejecutarser