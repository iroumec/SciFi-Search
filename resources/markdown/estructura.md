# Estructura del Proyecto

La aplicación está distribuida en cuatro directorios principales.

## `app`

Contiene todo el código Go de la aplicación. Dentro de él, pueden hallarse tres subdirectorios:

- `database`: contiene todos los archivos generados por SQLC. Estos permiten interactuar con la base de datos de forma sencilla desde el código Go.
- `handlers`: contiene todos los _handlers_ definidos para servir las rutas de la aplicación, tales como `/users`, `/search`, y demás.
- `utils`: contiene funcionalidades adicionales, tales como la compresión de archivos, _middlewares_ de _logging_, etc.
- `views`: contiene todos los códigos `.templ`, los cuales permiten servir páginas dinámicamente a partir de renderizado del lado del server.

## `database`

Adicionalmente al archivo `sqlc.yaml`, el cual permite generar el código Go mencionado con anterioridad, es posible hallar dos subdirectorios: `queries`, en la que se definen todas las consultas que pueden realizarse a la base de datos; y `schema`, en la que se define el esquema de nuestra base de datos, los _alter table_, _triggers_, entre otras funcionalidades.

## `static`

Contiene todos los archivos que se sirven de forma estática en la aplicación. Esto incluye: imágenes, archivos CSS, PDF y HTML y código JavaScript.

## `resources`

Este directorio no forma parte de la aplicación en sí. Únicamente contiene comandos, _scripts_ y recursos adicionales usados durante el planeamiento y desarrollo de la aplicación.
