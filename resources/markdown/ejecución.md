# Ejecución de la Aplicación

La aplicación cuenta con dos modos: el modo desarrollo y el modo producción.

## Modo Desarrollo

Este modo cuenta con la particularidad de que, además de utilizar una imagen que cuenta con las herramientas de Golang, se halla integrado Air, el cual permite que los cambios realizados en los archivos se reflejen automáticamente, lo que facilita el desarrollo.

Para iniciar la aplicación en este modo, debe ejecutar el siguiente comando:

```bash
docker compose up --build
```

## Modo Producción

Este modo se compone de una imagen liviana, compuesta únicamente de lo estrictamente necesario para correr la aplicación.

Para iniciar la aplicación en este modo, debe ejecutar el siguiente comando:

```bash
docker compose -f docker-compose.yml up -d --build
```

La ejecución inicial puede tardar unos pocos minutos. Finalizada esta, el servidor correrá en segundo plano y podrá acceder a él en: <http://localhost:8080/>.

Para detener el servidor, ejecute el siguiente comando:

```bash
docker compose -f docker-compose.yml down -v
```

De esta forma, el servidor se detiene y se borra el contenedor, volúmenes y redes asociados. No obstante, las imágenes continuarán existiendo.

> [!Important] No borra los volúmenes externos, solo los que fueron definidos dentro del docker-compose.yml.

## Limpieza

Para eliminar la imagen o volúmenes (útil en desarrollo ante cambios en la base de datos), ejecute el siguiente comando:

```sh
docker compose down -v --rmi all
```

**También se recomienda ejecutar el comando en caso de que no vaya a utilizar más la aplicación.**
