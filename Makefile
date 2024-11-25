FREEPDM=https://github.com/grd/FreePDM
FREECAD_FILES=$(FREEPDM)/ConceptOfDesign/TestFiles
BIN_DIR=${HOME}/bin


all: add_users createvault removevault testvault pdmserver pdmterm pdmclient vcs vfs 

add_users: 
	go build -o $(BIN_DIR)/add_users cmd/add_users/main.go

createvault: 
	go build -o $(BIN_DIR)/createvault cmd/createvault/main.go

removevault: 
	go build -o $(BIN_DIR)/removevault cmd/removevault/main.go
	removevault

test:
	go test -v ./...

fstest:
	go test -failfast pkg/filesystem/main_test.go

pdmserver: 
	go build -o $(BIN_DIR)/pdmserver cmd/pdmserver/main.go

pdmterm: 
	go build -o $(BIN_DIR)/pdmterm cmd/pdmterm/main.go

pdmclient:
	go build -o $(BIN_DIR)/pdmclient cmd/pdmclient/main.go

vcs: 
	go build -o $(BIN_DIR)/vcs cmd/vcs/main.go

vfs: 
	go build -o $(BIN_DIR)/vfs cmd/vfs/main.go

smb2:
	go build -o $(BIN_DIR)/smb2 temp/smb2.go
	
error-test:
	go build -o $(BIN_DIR)/error-test temp/error-test/error-test.go

file-manipulation:
	go build -o $(BIN_DIR)/file-manipulation temp/file-manipulation/file-manipulation.go

glob-walkdir:
	go build -o $(BIN_DIR)/glob-walkdir temp/glob-walkdir/glob-walkdir.go

list-share-names:
	go build -o $(BIN_DIR)/list-share-names temp/list-share-names/list-share-names.go

docker_all:
	docker_rm
	docker

docker:
	docker_compose pull
	docker-compose up --build -d

docker_stop:
	docker-compose down

docker_rm:
	docker-compose down -v