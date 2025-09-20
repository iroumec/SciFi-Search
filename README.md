# Olimpiadas UNICEN

Aplicación web de las Olimpiadas Interfacultativas de la UNICEN.

Trabajo Práctico Integrador para las materias Programación Web y Sistemas Operativos.

Alumnos:
- Roumec, Iñaki.
- Velis, Ulises.
- Zaffaroni, Gerónimo.

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
