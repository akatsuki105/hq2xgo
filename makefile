
ifdef COMSPEC
	EXE_EXT := .exe
else
	EXE_EXT := 
endif

.PHONY: build
build:
	go build -o hq2x$(EXE_EXT) ./cmd/main.go 

.PHONY: run
run:
	go run ./cmd/main.go
