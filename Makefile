.PHONY: build
build:
	go build ./cmd/vending

.PHONY: install
install:
	go install ./cmd/vending

.PHONY: clean
clean:
	go clean ./...
	rm -f coverage.out
	rm -f vending

.PHONY: test
test:
	go test -count=1 ./...

.PHONY: cover
cover:
	go test -count=1 -cover -coverprofile coverage.out ./...
	go tool cover -html coverage.out

.PHONY: mod-update
mod-update:
	go get -u ./...
	go mod tidy

git_dirty := $(shell git status -s)

.PHONY: git-clean-check
git-clean-check:
ifneq ($(git_dirty),)
	@echo "Git repository is dirty!"
	@false
else
	@echo "Git repository is clean."
endif

.PHONY: format
format:
	yq -i .github/workflows/test.yml
	for _file in $$(gofmt -s -l . | grep -vE '^vendor/'); do \
		gofmt -s -w $$_file ; \
	done

.PHONY: format-check
format-check:
ifneq ($(git_dirty),)
	$(error format-check must be invoked on a clean repository)
endif
	$(MAKE) format
	$(MAKE) git-clean-check

.PHONY: authors
authors:
	scripts/generate-authors.sh
