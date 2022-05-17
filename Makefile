.PHONY: build cleanO clean 

all: build

build: 
	@go build -o build/coder/coder cmd/coderMain/coderMain.go

run:
	@./build/coder/coder $(IN)

cleanO: 
	@rm data/output/*

clean:
	@rm build/coder/*

