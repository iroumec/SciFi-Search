# ==============================================================================
# Makefile para levantar un entorno Docker, ejecutar pruebas Hurl y limpiar.
# ==============================================================================

# --- Variables de Configuración ---
HURL_FILE ?= resources/tests/requests.hurl

# --- Definiciones de Targets ---
.PHONY: all test up run-tests down help

# Target por defecto que se ejecuta al correr `make`.
all: help

test: up ## Levanta el entorno esperando a que esté sano, ejecuta las pruebas y lo detiene.
	@echo "--> Ejecutando pruebas Hurl desde el archivo: ${HURL_FILE}"
	@hurl --test ${HURL_FILE}
	@echo "--> Pruebas finalizadas. Deteniendo el entorno..."
	@$(MAKE) down

up: ## Construye y levanta los contenedores, esperando a que el servidor avise.
	@echo
	@echo "Construyendo y levantando los contenedores de Docker..."
	@echo
	
	@docker compose -f docker-compose.yml up -d --build
	
	@echo "Contenedores iniciados. Esperando a que el servidor esté listo..."
	@echo

	@# Bucle de espera: Intenta conectarse a /health cada segundo.
	@# 'until' sigue intentando HASTA QUE el comando curl tenga éxito (salga con 0).
	@# -f: Falla en silencio (no muestra HTML) si hay un error HTTP (como 404 o 500).
	@# -s: Modo silencioso (no muestra la barra de progreso).
	@until curl -f -s http://localhost:8080/health > /dev/null; do \
		sleep 1; \
	done

	@echo "Servidor corriendo en http://localhost:8000."
	@echo

development: ## Construye y levanta los contenedores en modo desarrollador (con air activo).
	@docker compose up --build

run-tests: ## Ejecuta únicamente las pruebas Hurl (asume que el entorno ya está levantado).
	@echo "--> Verificando que el servidor esté corriendo en localhost:8080..."
	@curl -f -s http://localhost:8080/health > /dev/null || (echo "\033[31mERROR: El servidor no está corriendo. Usa 'make up' primero.\033[0m" && exit 1)
	@echo "--> Servidor detectado. Ejecutando pruebas Hurl desde el archivo: ${HURL_FILE}"
	@hurl --test ${HURL_FILE}

down: ## Detiene y elimina los contenedores, redes y volúmenes.
	@echo "Deteniendo el servidor..."
	@docker compose down
	@docker compose -f docker-compose.yml down -v

clean: down ## Elimina la imagen y los volúmenes.
	@docker compose down -v --rmi all
	@docker volume prune -f

help: ## Muestra los comandos disponibles.
	@echo "Comandos disponibles:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2}'