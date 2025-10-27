# Instrucciones

## TP3

Para ejecutar el servidor y correr el archivo `.hurl` con las pruebas, ejecute el siguiente comando en la carpeta del proyecto:

```sh
make test
```

Probablemente, la primera vez que lo ejecute deba esperar unos pocos minutos a que todas las imágenes se descarguen.

Si desea ver cuáles son las pruebas, puede hallar el archivo `.hurl` [aquí](../tests/requests.hurl).

## TP4

Para ejecutar el servidor en segundo plano, ejecute, en la carpeta del proyecto, el siguiente comando:

```sh
make up
```

Una vez terminada la inicialización, puede acceder a la aplicación _web_ desde un navegador, ingresando a `http://localhost:8080/`.

En la página principal, hallará un formulario que le permitirá cargar usuarios, cuyo listado se mostrará dinámicamente a la derecha. Adicionalmente, al lado de cada usuario cargado encontrará un botón que le brinda la posibilidad de eliminarlo.

Mientras el servidor está corriendo, puede ejecutar, utilizando Curl, los siguientes [comandos](api.md) para probar la API.

Para detener el servidor, ejecute el siguiente comando:

```sh
make down
```

Este detendrá los servicios. No obstante, no eliminará las imágenes ni volúmenes. Si requiere hacer una limpieza completa, ejecute:

```sh
make clean
```
