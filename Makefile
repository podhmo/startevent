test:
	go test .

format:
	goimports -w .

lint:
	go vet .
