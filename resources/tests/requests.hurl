# ==============================================================================
# Script de pruebas Hurl para la API CRUD de usuarios.
# Basado en la consigna de la API que devuelve JSON.
# ==============================================================================

# Prueba 1: Crear un nuevo usuario (POST)
# Verifica que el estado es 201 y captura el ID de la respuesta.
POST http://localhost:8080/users
Content-Type: application/json
{
    "name": "Alan",
    "middlename": "Mathison",
    "surname": "Turing"
}

HTTP 201
[Captures]
new_id: jsonpath "$.id"
[Asserts]
jsonpath "$.id" exists
jsonpath "$.name" == "Alan"
jsonpath "$.surname" == "Turing"


# Prueba 2: Listar todos los usuarios (GET)
# Verifica que el estado es 200 y que la respuesta es una lista (colección).
GET http://localhost:8080/users
HTTP 200
[Asserts]
body isCollection


# Prueba 3: Obtener el usuario recién creado (GET por ID)
# Usa la variable 'new_id' capturada en la Prueba 1.
GET http://localhost:8080/users/{{new_id}}
HTTP 200
[Asserts]
jsonpath "$.id" == {{new_id}}
jsonpath "$.name" == "Alan"


# Prueba 4: Actualizar el usuario (PUT)
# Usa 'new_id' para la URL y envía los datos actualizados.
PUT http://localhost:8080/users/{{new_id}}
Content-Type: application/json
{
    "name": "Alan",
    "middlename": "M.",
    "surname": "Turing"
}

# Asumimos que PUT devuelve 200 OK con el objeto actualizado
HTTP 200
[Asserts]
jsonpath "$.middlename" == "M."


# Verificación de la actualización (GET por ID)
# Volvemos a pedir el usuario para confirmar que los cambios se guardaron.
GET http://localhost:8080/users/{{new_id}}
HTTP 200
[Asserts]
jsonpath "$.middlename" == "M."


# Prueba 5: Eliminar el usuario (DELETE)
# La consigna pide un 204 No Content.
DELETE http://localhost:8080/users/{{new_id}}
HTTP 204


# Prueba 6: Verificar que el usuario fue eliminado (GET por ID)
# Pedir el mismo usuario ahora debe dar un 404.
GET http://localhost:8080/users/{{new_id}}
HTTP 404


# Prueba 7: Probar un 404 con un ID inexistente
GET http://localhost:8080/users/9999
HTTP 404
