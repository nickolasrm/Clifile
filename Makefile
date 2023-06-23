## -----------------------
## Clifile common commands
## -----------------------


help: ## Show this help.
	@sed -n '/# IGNORE$$/!s/## *//p' $(MAKEFILE_LIST)  # IGNORE

install: ## Install package dependencies
	@go install github.com/onsi/ginkgo/v2/ginkgo
	@go install golang.org/x/tools/cmd/godoc
	@go install

lint: ## Run linter
	@golint ./**/*

test: ## Run tests
	@ginkgo -r -v --cover
