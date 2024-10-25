# Variables for the MongoDB container
MONGO_IMAGE = mongodb/mongodb-community-server
CONTAINER_NAME = mongodb
MONGO_PORT = 27017
MONGO_DATA_PATH = $(PWD)/mongo_data
MONGO_INITDB_ROOT_USERNAME=admin
MONGO_INITDB_ROOT_PASSWORD=password

# Start MongoDB container
.PHONY: start-mongo
start-mongo:
	@echo "Starting MongoDB container..."
	@echo "CONTAINER_NAME=$(CONTAINER_NAME)"
	@echo "MONGO_PORT=$(MONGO_PORT)"
	@echo "MONGO_DATA_PATH=$(MONGO_DATA_PATH)"
	@echo "MONGO_INITDB_ROOT_USERNAME=$(MONGO_INITDB_ROOT_USERNAME)"
	@echo "MONGO_INITDB_ROOT_PASSWORD=$(MONGO_INITDB_ROOT_PASSWORD)"
	@echo "MONGO_IMAGE=$(MONGO_IMAGE)"
	docker run -d --name $(CONTAINER_NAME) \
		-p $(MONGO_PORT):27017 \
		-v $(MONGO_DATA_PATH):/data/db \
		-e MONGO_INITDB_ROOT_USERNAME=$(MONGO_INITDB_ROOT_USERNAME) \
		-e MONGO_INITDB_ROOT_PASSWORD=$(MONGO_INITDB_ROOT_PASSWORD) \
		$(MONGO_IMAGE)

# Stop MongoDB container
.PHONY: stop-mongo
stop-mongo:
	@echo "Stopping MongoDB container..."
	@docker stop $(CONTAINER_NAME)

# Remove MongoDB container
.PHONY: remove-mongo
remove-mongo: stop-mongo
	@echo "Removing MongoDB container..."
	@docker rm $(CONTAINER_NAME)

# View MongoDB container logs
.PHONY: logs-mongo
logs-mongo:
	@docker logs -f $(CONTAINER_NAME)

# Connect to MongoDB shell
.PHONY: shell-mongo
shell-mongo:
	@docker exec -it $(CONTAINER_NAME) mongo

# Clean up data
.PHONY: clean-mongo-data
clean-mongo-data:
	@echo "Removing MongoDB data..."
	@rm -rf $(MONGO_DATA_PATH)

# Start application
.PHONY: start
start:
	@echo "Starting application..."
	@MONGODB_URI="mongodb://$(MONGO_INITDB_ROOT_USERNAME):$(MONGO_INITDB_ROOT_PASSWORD)@localhost:$(MONGO_PORT)/test?authSource=admin" go run main.go