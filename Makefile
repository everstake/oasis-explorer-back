include common.mk

# Check which tool to use for downloading.
HAVE_WGET := $(shell which wget > /dev/null && echo 1)
ifdef HAVE_WGET
DOWNLOAD := wget --quiet --show-progress --progress=bar:force:noscroll -O
else
HAVE_CURL := $(shell which curl > /dev/null && echo 1)
ifdef HAVE_CURL
DOWNLOAD := curl --progress-bar --location -o
else
$(error Please install wget or curl)
endif
endif

generate_api:
	swagger generate server -t gen -f ./swagger/swagger.yml --exclude-main -A OasisExplorer

serve_ui:
	swagger serve swagger/swagger.yml

validate:
	swagger validate swagger/swagger.yml

# Lint code, commits and documentation.
lint-targets := lint-go lint-docs lint-changelog lint-git lint-go-mod-tidy

lint-go:
	@$(ECHO) "$(CYAN)*** Running Go linters...$(OFF)"
	@env -u GOPATH golangci-lint run

lint-git:
	@$(CHECK_GITLINT)

lint-docs:
	@$(ECHO) "$(CYAN)*** Runnint markdownlint-cli...$(OFF)"
	@npx markdownlint-cli '**/*.md' --ignore swagger/

lint-changelog:
	@$(CHECK_CHANGELOG_FRAGMENTS)

lint-go-mod-tidy:
	@$(ECHO) "$(CYAN)*** Checking go mod tidy...$(OFF)"
	@$(ENSURE_GIT_CLEAN)
	@$(CHECK_GO_MOD_TIDY)

lint: $(lint-targets)
