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
	
	@echo
	@echo "Contenedores iniciados. Esperando a que el servidor esté listo..."
	@until curl -f -s http://localhost:$(APP_PORT)/health > /dev/null 2>&1; do \
		sleep 1; \
	done
	@echo

	@# Bucle de espera: Intenta conectarse a /health cada segundo.
	@# 'until' sigue intentando HASTA QUE el comando curl tenga éxito (salga con 0).
	@# -f: Falla en silencio (no muestra HTML) si hay un error HTTP (como 404 o 500).
	@# -s: Modo silencioso (no muestra la barra de progreso).

	@echo "Servidor corriendo en http://localhost:$(APP_PORT)."
	@echo

# =================================================================================================

development: create-env ## Construye y levanta los contenedores en modo desarrollador.
	@docker compose up --build

# =================================================================================================

down: ## Detiene contenedores sin eliminar datos.
	@echo "Deteniendo contenedores..."
	@docker compose down
	@echo "Contenedores detenidos."

# =================================================================================================

clean: down ## Detiene contenedores y elimina volúmenes.
	@echo "Deteniendo servicios y eliminando volúmenes..."
	@docker compose -f docker-compose.yml down -v --rmi all
	@docker volume prune -f
	@echo "Limpieza completa realizada. Todos los datos fueron eliminados."

# =================================================================================================

clear-dependencies: ## Limpia las dependencias no utilizadas de la aplicación.
	@echo "Actualizando dependencias..."
	@go mod tidy
	@echo
	@echo "Dependencias actualizadas."

# =================================================================================================

is-running: ## Verifica que el servidor esté corriendo.
	@curl -f -s http://localhost:$(APP_PORT)/health > /dev/null \
		|| (echo "\033[31mERROR: El servidor no está corriendo. Usa 'make up' primero.\033[0m" \
			&& exit 1)

# =================================================================================================

create-env: ## Crea un archivo `.env` con valores por defecto si no existe.
	@test -f .env || cp resources/.env.example .env

# =================================================================================================

absolute-clean: ## Deja Docker como recién instalado.
	@echo "Stopping containers..."
	@docker ps -aq | xargs -r docker stop 2>/dev/null || true
	@echo "Deleting containers..."
	@docker ps -aq | xargs -r docker rm 2>/dev/null || true
	@echo "Deleting images..."
	@docker images -q | xargs -r docker rmi -f 2>/dev/null || true
	@echo "Deleting volumes..."
	@docker volume ls -q | xargs -r docker volume rm 2>/dev/null || true
	@echo "Cleaning networks..."
	@docker network prune -f
	@echo "Final cleaning of the system..."
	@docker system prune -a --volumes -f
	@echo "Docker completely cleaned."

# =================================================================================================

help: ## Muestra los comandos disponibles.
	@echo "Comandos disponibles:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort \
		| awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2}'

# =================================================================================================