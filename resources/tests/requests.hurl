# ==============================================================================
# Script de pruebas Hurl para la API de usuarios.
# ==============================================================================

# ------------------------------------------------------------------------------
# Prueba 1: Creación de usuarios (POST)
# ------------------------------------------------------------------------------

POST http://localhost:8080/api/users
Content-Type: application/json
{
    "name": "Alan",
    "surname": "Turing"
}

# Se verifica que el estado es 201 y se captura la respuesta.
HTTP 201

# Se captura el id devuelto.
[Captures]
new_id: jsonpath "$.user_id"

[Asserts]
jsonpath "$.user_id" exists
jsonpath "$.name" == "Alan"
jsonpath "$.surname" == "Turing"

# ------------------------------------------------------------------------------

# Creación del segundo usuario.
POST http://localhost:8080/api/users
Content-Type: application/json
{
    "name": "Marie",
    "surname": "Curie"
}

# Se verifica que la respuesta sea 201.
HTTP 201

[Asserts]
jsonpath "$.user_id" exists
jsonpath "$.name" == "Marie"
jsonpath "$.surname" == "Curie"

# ------------------------------------------------------------------------------
# Prueba 2: Intento de creación de usuario con datos faltantes (POST).
# ------------------------------------------------------------------------------

POST http://localhost:8080/api/users
Content-Type: application/json
{
    "name": "Albert"
}

# Se verifica que la respuesta sea 400.
HTTP 400

# ------------------------------------------------------------------------------
# Prueba 3: Listado de todos los usuarios (GET)
# ------------------------------------------------------------------------------

GET http://localhost:8080/api/users

# Se verifica que el estado es 200.
HTTP 200

# Se verifica que la respuesta es una lista (colección).
[Asserts]
jsonpath "$" isCollection

# ------------------------------------------------------------------------------
# Prueba 4: Obtención de uno de los usuarios recién creados (GET por ID)
# ------------------------------------------------------------------------------

# Se utiliza la variable `new_id` capturada en la primera prueba.
GET http://localhost:8080/api/users/{{new_id}}

HTTP 200

[Asserts]
jsonpath "$.user_id" == {{new_id}}
jsonpath "$.name" == "Alan"

# ------------------------------------------------------------------------------
# Prueba 5: Actualización de un usuario (PUT)
# ------------------------------------------------------------------------------

# Se actualiza el usuario cuya id fue capturada en `new_id`.
PUT http://localhost:8080/api/users/{{new_id}}
Content-Type: application/json
{
    "name": "Alan",
    "surname": "Rickman"
}

# Se verifica que se devuelva 200.
HTTP 200

# ------------------------------------------------------------------------------
# Prueba 6: Verificación de la actualización (GET por ID)
# ------------------------------------------------------------------------------

# Se vuelve a pedir el usuario para confirmar que los cambios se guardaron.
GET http://localhost:8080/api/users/{{new_id}}

HTTP 200

# Se verifica que el apellida se haya actualizado.
[Asserts]
jsonpath "$.surname" == "Rickman"

# ------------------------------------------------------------------------------
# Prueba 7: Eliminación de un usuario (DELETE)
# ------------------------------------------------------------------------------

# Prueba 5: Eliminar el usuario (DELETE)
# La consigna pide un 204 No Content.
DELETE http://localhost:8080/api/users/{{new_id}}

# Se verifica que el mensaje sea 204 NO Content.
HTTP 204

# ------------------------------------------------------------------------------
# Prueba 8: Verificación de que el usuario fue eliminado (GET por ID)
# ------------------------------------------------------------------------------

# Se pide el mismo usuario que se acaba de eliminar.
GET http://localhost:8080/api/users/{{new_id}}

# Se verifica que el mensaje sea 404.
HTTP 404

# ------------------------------------------------------------------------------
# Prueba 9: Obtención de usuario utilizando un ID inexistente.
# ------------------------------------------------------------------------------

# Se intenta obtener un usuario con un ID inexistente.
GET http://localhost:8080/api/users/99999

# Se verifica que el mensaje sea 404.
HTTP 404

# ------------------------------------------------------------------------------
