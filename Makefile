FREEPDM=https://github.com/grd/FreePDM
BIN_DIR=${HOME}/bin
CONTAINER_NAME := freepdm


all: add_users createvault removevault testvault pdmserver pdmterm pdmclient vcs vfs 

add_users: 
	go build -o $(BIN_DIR)/add_users ./cmd/add_users

createvault: 
	go build -o $(BIN_DIR)/createvault ./cmd/createvault
	createvault

removevault: 
	go build -o $(BIN_DIR)/removevault ./cmd/removevault
	removevault

test:
	go test -v ./...

vaultstest:
	go test -failfast internal/vaults/vaults_test.go

pdmserver: 
	go build -o $(BIN_DIR)/pdmserver ./cmd/pdmserver
	pdmserver

pdmterm: 
	go build -o $(BIN_DIR)/pdmterm ./cmd/pdmterm
	pdmterm

pdmclient:
	go build -o $(BIN_DIR)/pdmclient ./cmd/pdmclient
	pdmclient

vcs: 
	go build -o $(BIN_DIR)/vcs ./cmd/vcs

vfs: 
	go build -o $(BIN_DIR)/vfs ./cmd/vfs

smb2:
	go build -o $(BIN_DIR)/smb2 ./temp/smb2.go
	
error-test:
	go build -o $(BIN_DIR)/error-test ./temp/error-test/error-test.go

file-manipulation:
	go build -o $(BIN_DIR)/file-manipulation ./temp/file-manipulation/file-manipulation.go

glob-walkdir:
	go build -o $(BIN_DIR)/glob-walkdir ./temp/glob-walkdir/glob-walkdir.go

list-share-names:
	go build -o $(BIN_DIR)/list-share-names ./temp/list-share-names/list-share-names.go

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
