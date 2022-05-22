.PHONY: build cleanO clean run runDefault

all: build

build: 
	@go build -o build/coder/coder cmd/coderMain/coderMain.go

run:
	@./build/coder/coder $(IN) $(OUT) $(NUM)

runDefault:
	@./build/coder/coder data/input/testy4/example0.tga data/output/def 8

cleanO: 
	@rm data/output/*

clean:
	@rm build/coder/*

