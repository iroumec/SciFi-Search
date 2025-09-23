# Estructura del Proyecto

La aplicación está distribuida en cuatro directorios principales.

## `app`

Contiene todo el código Go de la aplicación. Dentro de él, pueden hallarse tres subdirectorios:

- `database`: contiene todos los archivos generados por SQLC. Estos permiten interactuar con la base de datos de forma sencilla desde el código Go.
- `handlers`: contiene todos los _handlers_ definidos para servir las rutas de la aplicación, tales como `/registrarse`, `/noticias`, y demás.
- `utils`: contiene funcionalidades adicionales, tales como la compresión de archivos y la validación de las constancias de alumnos.

## `database`

Adicionalmente al archivo `sqlc.yaml`, el cual permite generar el código Go mencionado con anterioridad, es posible hallar dos subdirectorios: `queries`, en la que se definen todas las consultas que pueden realizarse a la base de datos; y `schema`, en la que se define el esquema de nuestra base de datos, los _alter table_, _triggers_, entre otras funcionalidades.

## `static`

Contiene todos los archivos que se sirven de forma estática en la aplicación. Esto incluye: imágenes, archivos CSS, archivos PDF y código JavaScript.

## `template`

Contiene todos los archivos `.html` servidos en la aplicación. Estos están modularizados de forma de evitar la redundancia y, agrupados de acuerdo al propósito que sirven.

## `resources`

Este directorio no forma parte de la aplicación en sí. Únicamente contiene comandos, _scripts_ y recursos adicionales usados durante el planeamiento y desarrollo de la aplicación.
