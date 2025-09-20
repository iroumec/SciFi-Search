# UKI

Trabajo Práctico de Cursada para la materia Programación Web.

Uki es una aplicación web que permite a los usuarios puntuar contenido multimedia y compartir sus gustos con los demás usuarios de la plataforma.

## Configuración del entorno

Antes de ejecutar la aplicación, crea tu archivo de variables de entorno basado en el ejemplo:

```bash
cp resources/.env.example .env
```

Luego edita `.env` y completa tus credenciales de la base de datos y otros parámetros necesarios.

## Ejecución

Modo desarrollo (con AIR):

```bash
docker compose up --build
```

Modo producción:

```bash
docker compose -f docker-compose.yml up -d --build
```

## Limpieza

Ante cambios en la base de datos, es necesario eliminar el volumen y reconstruirlo. Para ello, ejecute el siguiente comando:

```sh
docker compose down -v --rmi all
```

IDEAS:

- Seguir deportes.
- Predicciones.
