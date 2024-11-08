# Variables for the MongoDB container
MONGO_IMAGE = mongodb/mongodb-community-server
MONGO_CONTAINER_NAME = mongodb
MONGO_PORT = 27017
MONGO_DATA_PATH = $(PWD)/mongo_data
MONGO_INITDB_ROOT_USERNAME=admin
MONGO_INITDB_ROOT_PASSWORD=password

# Variables for the Redis container
REDIS_IMAGE = redis
REDIS_CONTAINER_NAME = redis
REDIS_PORT = 6379
REDIS_CONFIGURATION_PATH = $(PWD)/conf

# Start MongoDB container
.PHONY: start-mongo
start-mongo:
	@echo "Starting MongoDB container..."
	@echo "MONGO_CONTAINER_NAME=$(MONGO_CONTAINER_NAME)"
	@echo "MONGO_PORT=$(MONGO_PORT)"
	@echo "MONGO_DATA_PATH=$(MONGO_DATA_PATH)"
	@echo "MONGO_INITDB_ROOT_USERNAME=$(MONGO_INITDB_ROOT_USERNAME)"
	@echo "MONGO_INITDB_ROOT_PASSWORD=$(MONGO_INITDB_ROOT_PASSWORD)"
	@echo "MONGO_IMAGE=$(MONGO_IMAGE)"
	docker run -d --name $(MONGO_CONTAINER_NAME) \
		-p $(MONGO_PORT):27017 \
		-v $(MONGO_DATA_PATH):/data/db \
		-e MONGO_INITDB_ROOT_USERNAME=$(MONGO_INITDB_ROOT_USERNAME) \
		-e MONGO_INITDB_ROOT_PASSWORD=$(MONGO_INITDB_ROOT_PASSWORD) \
		$(MONGO_IMAGE)

# Stop MongoDB container
.PHONY: stop-mongo
stop-mongo:
	@echo "Stopping MongoDB container..."
	@docker stop $(MONGO_CONTAINER_NAME)

# Remove MongoDB container
.PHONY: remove-mongo
remove-mongo: stop-mongo
	@echo "Removing MongoDB container..."
	@docker rm $(MONGO_CONTAINER_NAME)

# View MongoDB container logs
.PHONY: logs-mongo
logs-mongo:
	@docker logs -f $(MONGO_CONTAINER_NAME)

# Connect to MongoDB shell
.PHONY: shell-mongo
shell-mongo:
	@docker exec -it $(MONGO_CONTAINER_NAME) mongo

# Clean up data
.PHONY: clean-mongo-data
clean-mongo-data:
	@echo "Removing MongoDB data..."
	@rm -rf $(MONGO_DATA_PATH)

# Start Redis container
.PHONY: start-redis
start-redis:
	@echo "Starting Redis container..."
	@echo "REDIS_CONTAINER_NAME=$(REDIS_CONTAINER_NAME)"
	@echo "REDIS_PORT=$(REDIS_PORT)"
	@echo "REDIS_CONFIGURATION_PATH=$(REDIS_CONFIGURATION_PATH)"
	@echo "REDIS_IMAGE=$(REDIS_IMAGE)"
	docker run -d --name $(REDIS_CONTAINER_NAME) \
		-p $(REDIS_PORT):6379 \
		-v $(REDIS_CONFIGURATION_PATH):/usr/local/etc/redis \
		$(REDIS_IMAGE)

# Stop Redis container
.PHONY: stop-redis
stop-redis:
	@echo "Stopping Redis container..."
	@docker stop $(REDIS_CONTAINER_NAME)

# Remove Redis container
.PHONY: remove-redis
remove-redis: stop-redis
	@echo "Removing Redis container..."
	@docker rm $(REDIS_CONTAINER_NAME)

# View Redis container logs
.PHONY: logs-redis
logs-redis:
	@docker logs -f $(REDIS_CONTAINER_NAME)

# Connect to Redis CLI
.PHONY: cli-redis
cli-redis:
	@docker exec -it $(REDIS_CONTAINER_NAME) redis-cli

# Clean up Redis data
.PHONY: clean-redis-data
clean-redis-data:
	@echo "Removing Redis data..."
	@rm -rf $(REDIS_CONFIGURATION_PATH)


# Start application
.PHONY: start
start:
	@echo "Starting application..."
	@MONGODB_URI="mongodb://$(MONGO_INITDB_ROOT_USERNAME):$(MONGO_INITDB_ROOT_PASSWORD)@localhost:$(MONGO_PORT)/test?authSource=admin" REDIS_ADDR="localhost:$(REDIS_PORT)" go run main.go

# Test application
.PHONY: test
test:
	@echo "Running tests..."
	@MONGODB_URI="mongodb://$(MONGO_INITDB_ROOT_USERNAME):$(MONGO_INITDB_ROOT_PASSWORD)@localhost:$(MONGO_PORT)/test?authSource=admin" go test -v ./...