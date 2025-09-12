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

Los siguientes comandos solo deben ejecutarse una vez:

```bash
# Borrado de contenedores y volúmenes.
sudo docker compose down -v

# Construcción del contenedor.
sudo docker compose build
```

Para correr la aplicación:

```bash
sudo docker compose up
```
