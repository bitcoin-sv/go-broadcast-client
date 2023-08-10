.PHONY: test lint
test:
    go test ./...
lint:
    golangci-lint run

.PHONY: godocs
godocs: ## Sync the latest tag with GoDocs
	@echo "syndicating to GoDocs..."
	@test $(GIT_DOMAIN)
	@test $(REPO_OWNER)
	@test $(REPO_NAME)
	@test $(VERSION_SHORT)
	@curl https://proxy.golang.org/$(GIT_DOMAIN)/$(REPO_OWNER)/$(REPO_NAME)/@v/$(VERSION_SHORT).info 
