# Variables
COMPOSE_FILE = docker-compose.yml
DOCKER_VOLUME = $(shell docker volume ls -q --filter name=mysql_data)
SERVICE_NAME = database

.PHONY: start stop reset rebuild

# Start all containers
start:
	docker-compose -f $(COMPOSE_FILE) up -d
	@echo "All containers started."

# Stop all containers
stop:
	docker-compose -f $(COMPOSE_FILE) down
	@echo "All containers stopped."

# Reset (clear database data)
reset:
	docker-compose -f $(COMPOSE_FILE) down --volumes
	@echo "Database data cleared and containers stopped."

# Rebuild a specific container (e.g., database)
rebuild:
	docker-compose -f $(COMPOSE_FILE) up --build $(SERVICE_NAME)
	@echo "$(SERVICE_NAME) container rebuilt."

# Rebuild all containers
rebuild_all:
	docker-compose -f $(COMPOSE_FILE) up --build --force-recreate
	@echo "All containers rebuilt."