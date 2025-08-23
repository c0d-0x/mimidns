all:
	@go build ./cmd/mimidns/
	@echo "build complete"

clean:
	@rm mimidns
	@echo "cleaned up"
