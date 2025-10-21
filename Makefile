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

up: ## Construye y levanta los contenedores, esperando a que los servicios estén listos.
	@echo "--> Construyendo y levantando los contenedores de Docker..."
	# La opción --wait en el siguiente comando sirve para esperar a que los "healthchecks" pasen.
	@docker compose up --build -d --wait

run-tests: ## Ejecuta únicamente las pruebas Hurl (asume que el entorno ya está levantado).
	@echo "--> Verificando que el servidor esté corriendo en localhost:8080..."
	@curl -f -s http://localhost:8080/health > /dev/null || (echo "\033[31mERROR: El servidor no está corriendo. Usa 'make up' primero.\033[0m" && exit 1)
	@echo "--> Servidor detectado. Ejecutando pruebas Hurl desde el archivo: ${HURL_FILE}"
	@hurl --test ${HURL_FILE}

down: ## Detiene y elimina los contenedores, redes y volúmenes.
	@echo "--> Deteniendo y limpiando los contenedores de Docker..."
	@docker compose down --volumes

help: ## Muestra los comandos disponibles.
	@echo "Comandos disponibles:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2}'