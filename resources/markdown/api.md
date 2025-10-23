# Uso de la API

Parámetros:

- `<name>`: nombre del usuario.
- `<middlename>`: segundo nombre del usuario. Puede estar vacío.
- `<surname>`: apellido del usuario.
- `<id>`: id de interés.
- `-i`: (Opcional, pero recomendado) Muestra los headers de la respuesta HTTP

## Añadido de un Usuario

```sh
curl -X POST -H "Content-Type: application/json" -d '{"name": "<name>", "middlename": "<middlename>", "surname": "<surname>"}' http://localhost:8080/api/users
```

## Actualización de un Usuario

```sh
curl -i -X PUT -H "Content-Type: application/json" -d '{"name": "<name>", "middlename": "<middlename>", "surname": "<surnam>"}' "http://localhost:8080/api/users?id=<id>"
```

## Eliminación de un Usuario

```sh
curl -i -X DELETE "http://localhost:8080/api/users?id=<id>"
```

## Obtener los Datos de un Usuario

```sh
curl -i "http://localhost:8080/api/users?id=<id>"
```

## Visualización de Todos los Usuarios

```sh
curl -i http://localhost:8080/api/users
```

Puede mejorarse la salida mediante el uso de `jq`:

```sh
curl http://localhost:8080/api/users | jq
```

Como punto en contra, no se mostrará el tipo de respuesta del servidor.
