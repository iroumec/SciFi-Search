# Ejecución de la Aplicación

La aplicación cuenta con dos modos: el modo desarrollo y el modo producción.

## Modo Producción

Este modo se compone de una imagen liviana, compuesta únicamente de lo estrictamente necesario para correr la aplicación.

Para iniciar la aplicación en este modo, debe ejecutar el siguiente comando:

```bash
make up
```

La ejecución inicial puede tardar unos pocos minutos. Finalizada esta, el servidor correrá en segundo plano y podrá acceder a él en: <http://localhost:8080/>.

Para detener el servidor, ejecute el siguiente comando:

```bash
make down
```

De esta forma, el servidor se detiene y se borra el contenedor y redes asociadas. No obstante, las imágenes continuarán existiendo y los volúmenes persistirán.

> [!IMPORTANT]
> No borra los volúmenes externos, solo los que fueron definidos dentro del docker-compose.yml.

## Modo Desarrollo

Este modo cuenta con la particularidad de que, además de utilizar una imagen que cuenta con las herramientas de Golang, se halla integrado Air, el cual permite que los cambios realizados en los archivos se reflejen automáticamente, lo que facilita el desarrollo.

Para iniciar la aplicación en este modo, debe ejecutar el siguiente comando:

```bash
make development
```

En este modo, el servidor correrá en primer plano, por lo que no requiere de `make down` para detenerlo. Con `Ctrl + C` es suficiente.

## Limpieza

Para eliminar la imagen y volúmenes (útil en desarrollo ante cambios en la base de datos), ejecute el siguiente comando:

```sh
make clean
```

**También se recomienda ejecutar el comando en caso de que no vaya a utilizar más la aplicación.**
