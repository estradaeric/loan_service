# Makefile for Loan Service Project
# Dev (hot reload) & Prod (build/run) workflows

# Container names
DEV_SERVICE=loan_service_dev
PROD_SERVICE=loan_service

# -----------------------------------
# 🧪 Development Commands (Air Hot Reload)
# -----------------------------------

# Clean & rebuild dev container
dev-build:
	docker-compose -f docker-compose.dev.yml --env-file .env.dev down -v
	docker-compose -f docker-compose.dev.yml --env-file .env.dev build --no-cache
	docker-compose -f docker-compose.dev.yml --env-file .env.dev up

# View logs from dev container
dev-logs:
	docker logs -f $(DEV_SERVICE)

# Stop all dev containers and volumes
dev-stop:
	docker-compose -f docker-compose.dev.yml down -v --remove-orphans

# -----------------------------------
# 🚀 Production Commands (Dockerized)
# -----------------------------------

# Full prod clean + rebuild + run
prod-build:
	docker-compose --env-file .env.prod down -v --remove-orphans
	docker-compose --env-file .env.prod build --no-cache
	docker-compose --env-file .env.prod up -d

# View logs from prod container
prod-logs:
	docker logs -f $(PROD_SERVICE)

# Stop all prod containers and volumes
prod-stop:
	docker-compose --env-file .env.prod down -v --remove-orphans