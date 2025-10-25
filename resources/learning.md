# Learning

En PostgreSQL **da lo mismo en rendimiento** usar `VARCHAR(255)` o `TEXT`. A diferencia de MySQL, Postgres no limita internamente `VARCHAR(n)` a `n` bytes: solo lo valida al insertar/actualizar.

## Diferencias

- `VARCHAR(255)`

  - Tiene una restricción de longitud máxima.
  - Si intentás guardar un hash mayor (no es el caso de bcrypt, siempre \~60 chars), fallará.
  - Es útil si querés **validar** que no se guarden cadenas más largas de lo esperado.

- `TEXT`

  - No tiene límite.
  - Ideal si no necesitás restricciones de longitud.
  - Más flexible para futuros cambios (si más adelante cambiás de bcrypt a Argon2, por ejemplo, que produce hashes más largos).

---

El navegador no peude acceder a ningún recurso que no esté en static (la ruta de archvios que se sirve). Por eso, se pone el css y laas imágenes ahí.

Cuidado con UTF-8 con BOM... No renderiza los HTML.

docker compose up --build
En lugar de:
docker compose build --no-cache
docker compose up
Esto hace build solo lo que cambió y levanta el contenedor. Mucho más rápido que --no-cache todo el tiempo.

En las funciones, minúscula inicial si es privado del paquete y mayúscula inicial si es público.

---

2025/10/16 14:56:02 cannot connect to db:sql: unknown driver "postgres" (forgotten import?)

Solución, falta:

\_ "github.com/lib/pq"

---

El html te dice qué hay y cómo buscarlo. El CSS es la forma.

el código 200 es el por defecto.

El patrón estándar de un middleware en Go es recibir un handler y devolver otro.

## Meili

go get github.com/meilisearch/meilisearch-go

go mod tidy # Para actualziar las dependnecias.

sudo rm -rf /home/iroumec/Documents/University/"Programación Web"/TPE/meili_data

## Supertokens

go get github.com/supertokens/supertokens-golang
go mod tidy

PLANEAR LA APLICACIÓN PARA TENER CONCURRENCIA A FUTURO.
