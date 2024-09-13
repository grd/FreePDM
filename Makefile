FREEPDM=https://github.com/grd/FreePDM
FREECAD_FILES=$(FREEPDM)/ConceptOfDesign/TestFiles
BIN_DIR=${HOME}/bin


all: add_users createvault removevault testvault serverapp setup vcs vfs

add_users: 
	go build -o $(BIN_DIR)/add_users src/bin/add_users/main.go

createvault: 
	go build -o $(BIN_DIR)/createvault src/bin/createvault/main.go

add_samples: 
	go build -o $(BIN_DIR)/add_samples src/bin/createvault/add_samples/add_samples.go

removevault: 
	go build -o $(BIN_DIR)/removevault src/bin/removevault/main.go
	removevault

test:
	go test src/bin/testvault/main_test.go

serverapp: 
	go build -o $(BIN_DIR)/serverapp src/bin/serverapp/main.go

setup: 
	go build -o $(BIN_DIR)/setup src/bin/setup/main.go

vcs: 
	go build -o $(BIN_DIR)/vcs src/bin/vcs/main.go

vfs: 
	go build -o $(BIN_DIR)/vfs src/bin/vfs/main.go

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

