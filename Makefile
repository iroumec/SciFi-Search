# ==============================================================================
# Makefile para levantar un entorno Docker, ejecutar pruebas Hurl y limpiar.
# ==============================================================================

# --- Variables de Configuración ---
HURL_VERSION ?= 7.0.0
INSTALL_DIR ?= /tmp
HURL_BIN_DIR = $(INSTALL_DIR)/hurl-$(HURL_VERSION)-x86_64-unknown-linux-gnu/bin
HURL_FILE ?= resources/tests/requests.hurl

# Target por defecto que se ejecuta al correr `make`.
all: help

test: up ## Levanta el entorno, ejecuta las pruebas y lo detiene.
	@$(MAKE) -s run-tests
	@echo
	@echo "--> Pruebas finalizadas. Deteniendo el entorno..."
	@echo
	@$(MAKE) -s down

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

run-tests: is-running ## Ejecuta únicamente las pruebas Hurl (asume que el entorno ya está levantado).
	@echo
	@# Instalación universal de hurl, independientemente de la distribución de Linux.
	@if [ ! -x "$(HURL_BIN_DIR)/hurl" ]; then \
		echo "Hurl no encontrado en $(HURL_BIN_DIR). Instalando versión $(HURL_VERSION)..."; \
		curl --silent --location https://github.com/Orange-OpenSource/hurl/releases/download/$(HURL_VERSION)/hurl-$(HURL_VERSION)-x86_64-unknown-linux-gnu.tar.gz \
			| tar xvz -C $(INSTALL_DIR); \
		echo "Hurl instalado en $(HURL_BIN_DIR)"; \
	fi

	@echo "--> Ejecutando pruebas con hurl..."
	@echo
	@PATH=$(HURL_BIN_DIR):$$PATH hurl --test ${HURL_FILE}

down: ## Detiene y elimina los contenedores, redes y volúmenes.
	@echo "Deteniendo el servidor..."
	@docker compose down
	@docker compose -f docker-compose.yml down -v

clean: down ## Elimina la imagen y los volúmenes.
	@docker compose down -v --rmi all
	@docker volume prune -f

is-running: ## Verifica que el servidor esté corriendo.
	@curl -f -s http://localhost:8080/health > /dev/null || (echo "\033[31mERROR: El servidor no está corriendo. Usa 'make up' primero.\033[0m" && exit 1)

help: ## Muestra los comandos disponibles.
	@echo "Comandos disponibles:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2}'