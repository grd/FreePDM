FREEPDM=https://github.com/grd/FreePDM
FREECAD_FILES=$(FREEPDM)/ConceptOfDesign/TestFiles
BIN_DIR=${HOME}/bin


all: add_users createvault removevault testvault serverapp setup vcs vfs

add_users: 
	go build -o $(BIN_DIR)/add_users src/bin/add_users/main.go

createvault: 
	go build -o $(BIN_DIR)/create_vault src/bin/create_new_vault/main.go

add_samples: 
	go build -o $(BIN_DIR)/add_samples src/bin/create_vault/add_samples/add_samples.go

removevault: 
	go build -o $(BIN_DIR)/removevault src/bin/removevault/main.go
	removevault


testvault: 
	rm -f $(BIN_DIR)/testvault
	make build_testvault runtestvault

build_testvault: 
	go build -o $(BIN_DIR)/testvault src/bin/testvault/main.go


run_testvault: 
	$(BIN_DIR)/testvault

serverapp: 
	go build -o $(BIN_DIR)/serverapp src/bin/serverapp/main.go

setup: 
	go build -o $(BIN_DIR)/setup src/bin/setup/main.go

vcs: 
	go build -o $(BIN_DIR)/vcs src/bin/vcs/main.go

vfs: 
	go build -o $(BIN_DIR)/vfs src/bin/vfs/main.go

