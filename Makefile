# ==============================================================================
# Makefile para levantar un entorno Docker, ejecutar pruebas Hurl y limpiar.
# ==============================================================================

# Target por defecto que se ejecuta al correr `make`.
all: help

run: up ## Construye y levanta los contenedores, esperando a que el servidor avise.

up: create-env ## Construye y levanta los contenedores, esperando a que el servidor avise.
	@echo
	@echo "Construyendo y levantando los contenedores de Docker..."
	@echo
	
	@docker compose -f docker-compose.yml up -d --build
	
	@echo "Contenedores iniciados."
	@echo

	@# Bucle de espera: Intenta conectarse a /health cada segundo.
	@# 'until' sigue intentando HASTA QUE el comando curl tenga éxito (salga con 0).
	@# -f: Falla en silencio (no muestra HTML) si hay un error HTTP (como 404 o 500).
	@# -s: Modo silencioso (no muestra la barra de progreso).
	
	@echo "Esperando a que la base de datos esté lista..."
	@until docker exec scifi-search-db pg_isready -U postgres > /dev/null 2>&1; do \
		sleep 1; \
	done
	@echo

	@echo "Aplicando migraciones..."
	@docker exec -i scifi-search-db psql -U postgres -d postgres  < database/schema/schema.sql > /dev/null 2>&1
	@echo

	@echo "Base de datos lista. Esperando a que el servidor esté listo..."
	@until curl -f -s http://localhost:8080/health > /dev/null 2>&1; do \
		sleep 1; \
	done
	@echo

	@echo "Servidor corriendo en http://localhost:8080."
	@echo

development: create-env ## Construye y levanta los contenedores en modo desarrollador (con air activo).
	@docker compose up --build

down: ## Detiene los contenedores y redes, sin eliminar volúmenes.
	@echo "Deteniendo el servidor..."
	@ #Se detiene la versión de producción. La versión no sobrescritra.
	@docker compose -f docker-compose.yml down

clean: down ## Elimina la imagen y los volúmenes.
	@docker compose -f docker-compose.yml down -v --rmi all
	@docker volume prune -f

is-running: ## Verifica que el servidor esté corriendo.
	@curl -f -s http://localhost:8080/health > /dev/null || (echo "\033[31mERROR: El servidor no está corriendo. Usa 'make up' primero.\033[0m" && exit 1)

create-env: ## Crea un archivo `.env` con valores por defecto para las variables de ambiente.
	@cp resources/.env.example .env

help: ## Muestra los comandos disponibles.
	@echo "Comandos disponibles:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2}'