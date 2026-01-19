# =================================================================================================
# Makefile para levantar un entorno Docker, ejecutar pruebas Hurl y limpiar.
# =================================================================================================

# Importación de variables de ambiente.
-include .env
export

# =================================================================================================

# Target por defecto que se ejecuta al correr `make`.
all: help

# =================================================================================================

run: up ## Construye y levanta los contenedores, esperando a que el servidor avise.

# =================================================================================================

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
	
	@echo "Esperando a que la base de datos esté healthy..."
	@until [ "$$(docker inspect -f '{{.State.Health.Status}}' database)" = "healthy" ]; do \
		sleep 1; \
	done
	@echo

	@echo "Base de datos lista. Esperando a que el servidor esté listo..."
	@until curl -f -s http://localhost:$(APP_PORT)/health > /dev/null 2>&1; do \
		sleep 1; \
	done
	@echo

	@echo "Servidor corriendo en http://localhost:$(APP_PORT)."
	@echo

# =================================================================================================

development: create-env ## Construye y levanta los contenedores en modo desarrollador.
	@docker compose up --build

# =================================================================================================

down: ## Detiene contenedores sin eliminar datos.
	@echo "Deteniendo y eliminando contenedores..."
	@docker compose down
	@echo "Contenedores detenidos. Los datos persisten en los volúmenes."

# =================================================================================================

clean: down ## Detiene contenedores y elimina volúmenes.
	@echo "Deteniendo servicios y eliminando volúmenes..."
	@docker compose -f docker-compose.yml down -v --rmi all
	@docker volume prune -f
	@echo "Limpieza completa realizada. Todos los datos fueron eliminados."

# =================================================================================================

is-running: ## Verifica que el servidor esté corriendo.
	@curl -f -s http://localhost:$(APP_PORT)/health > /dev/null \
		|| (echo "\033[31mERROR: El servidor no está corriendo. Usa 'make up' primero.\033[0m" \
			&& exit 1)

# =================================================================================================

create-env: ## Crea un archivo `.env` con valores por defecto si no existe.
	@test -f .env || cp resources/.env.example .env

# =================================================================================================

help: ## Muestra los comandos disponibles.
	@echo "Comandos disponibles:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort \
		| awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2}'

# =================================================================================================