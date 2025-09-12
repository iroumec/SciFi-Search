# Comandos

## Docker

### Creación del Contenedor

`sudo docker run --name UKIDB -e POSTGRES_PASSWORD=postgres -p 5432:5432 -v "$(pwd)/database/schema:/database/schema" -d postgres`

Solo ejecutar una vez.

### Carga del Esquema

`sudo docker exec -i UKIDB psql -U postgres -d postgres -f /database/schema/schema.sql`

Siempre que haya cambios en el esquema.

### Activación del Contenedor

`sudo docker start UKIDB`

En otro caso, al interactuar con la base de datos, aparecerá el siguiente error:

_error creating user: dial tcp [::1]:5432: connect: connection refused_

Recordar luego desactivar el contenedor:

`sudo docker stop UKIDB`

### Visualización de Contenedores Activos

sudo docker ps

### Datos sobre Docker

No es posible usar un path relativo entre tu máquina y el contenedor Docker. Los paths relativos solo funcionan dentro del mismo sistema de archivos. El contenedor Docker tiene su propio sistema de archivos separado del host. Por eso, debes usar un path absoluto de tu máquina al montar el volumen con -v, para que Docker sepa exactamente qué carpeta compartir con el contenedor. Dentro del contenedor, puedes usar paths relativos, pero solo respecto al sistema de archivos del contenedor.

Resumen:

Path absoluto en el parámetro -v para Docker.
Path relativo solo dentro del contenedor, si el archivo ya está disponible ahí.

$(pwd) obtiene el directorio actual, así la ruta funciona en cualquier PC.
El archivo SQL se monta en /database/schema dentro del contenedor, y luego se usa esa ruta en el segundo comando.
Si quieres aún más flexibilidad, puedes usar variables de entorno para definir rutas o nombres de contenedores.

## Creación de Base de Datos

go mod init UKI
go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
export PATH="$PATH:$(go env GOPATH)/bin"
sqlc generate
go get github.com/lib/pq

## Error: error creating user:pq: relation "users" does not exist

sudo docker exec -i postgresDB psql -h localhost -p 5432 -U postgres -d postgres -f ./db/schema/schema.sql

(No se especificó la ruta absoluta y se creo el docker en home). Se debe especificar la ruta absoluta de schema.sql y crear el docker como se especifico arriba.

## Ver datos

sudo docker exec -it UKIDB psql -U postgres -d postgres
SELECT \* FROM users;

O exportando...

sudo docker exec -i UKIDB pg_dump -U postgres -d postgres > T2/E2/datos.sql

Comando `tree` para mostrar la estructura de directorio.
