include app.env
export

FREEPDM=https://github.com/grd/FreePDM
BIN_DIR=${HOME}/bin
CONTAINER_NAME := freepdm


all: run createvault removevault test pdmserver gui

fpg:
	go build -o $(BIN_DIR)/fpg ./apps/fpg/main.go

fpg-run:
	go run ./apps/fpg/main.go

run: 
	air

createvault: 
	go build -o $(BIN_DIR)/createvault ./cmd/createvault
	createvault

removevault: 
	go build -o $(BIN_DIR)/removevault ./cmd/removevault
	removevault

test:
	go test -failfast -v ./...

vaultstest:
	go test -failfast internal/vaults/vaults_test.go

pdmserver: 
	go build -o $(BIN_DIR)/pdmserver ./cmd/pdmserver
	pdmserver

# SQLite database clean-up
reset-db-sqlite:
	rm -f $(SQLITE_PATH)

# PostgreSQL database clean-up
reset-db-postgres:
	PGPASSWORD=$(PG_PASSWORD) dropdb -h $(PG_HOST) -U $(PG_USER) $(PG_DB)
	PGPASSWORD=$(PG_PASSWORD) createdb -h $(PG_HOST) -U $(PG_USER) $(PG_DB)

# Stop interfering activities. Run `make docker` after this.
docker-up:
	sudo systemctl stop postgresql
	sudo systemctl stop smbd

# And restart the local interfering activities again.
local-up:
	docker-compose down
	sudo systemctl start postgresql
	sudo systemctl start smbd

docker:
	@$(MAKE) docker_stop
	@$(MAKE) docker_start

docker_start:
	@echo "Updating Docker containers..."
	docker-compose pull --ignore-pull-failures
	docker-compose up --build -d
	@echo "Docker containers updated successfully."

docker_stop:
	docker-compose down || true

docker_rm:
	docker-compose down -v || true
	docker rm -f ${CONTAINER_NAME} || true

docker_shell:
	docker exec -it ${CONTAINER_NAME} /bin/sh

docker_logs:
	docker logs ${CONTAINER_NAME}

