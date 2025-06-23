.PHONY: help push build-docker bump-version

help:
	@echo "Available commands:"
	@echo "  help                - Show this help message"
	@echo "  build               - Build all Go binaries for all platforms (in Docker, output to ./dist)"
	@echo "  push                - Update files in gist (binaries from ./dist, install/uninstall scripts, etc.)"
	@echo "  push!               - build + push"
	@echo "  bump-version        - Update version in main.go and create git tag"
	@echo "  push-version        - Push changes and tags to remote repository"
	@echo "  show-versions       - Display all git tags"
	@echo "  install             - Install sitedog"
	@echo "  uninstall           - Uninstall sitedog"
	@echo "  reinstall           - Uninstall and install sitedog"

push:
	rm -rf sitedog_gist
	git clone git@gist.github.com:fe278d331980a1ce09c3d946bbf0b83b.git --depth 1 sitedog_gist
	rm -rf sitedog_gist/*
	cp demo.html.tpl scripts/install.sh scripts/uninstall.sh sitedog_gist/
	cd sitedog_gist && \
	if git diff --quiet; then \
		echo "No changes to deploy"; \
	else \
		git add . && \
		git commit -m "Update sitedog files" && \
		git push; \
	fi

build:
	docker run --rm -v $(PWD):/app -w /app golang:1.20-alpine sh -c "./scripts/build.sh"

push!: build push

bump-version:
	@if [ -z "$(v)" ]; then \
		echo "Usage: make bump-version v=x.y.z"; \
		exit 1; \
	fi; \
	sed -i 's/Version[ ]*=[ ]*".*"/Version = "$(v)"/' main.go; \
	go fmt main.go; \
	git add main.go; \
	git commit -m "bump version to $(v)"; \
	git tag $(v); \
	echo "Version updated to $(v) and git tag created."

push-version:
	git push
	git push --tags

show-versions:
	@git tag -l

install:
	scripts/install.sh

dev-install:
	@echo "Detecting platform and architecture..."
	@PLATFORM=$$(uname -s | tr '[:upper:]' '[:lower:]'); \
	ARCH=$$(uname -m); \
	if [ "$$ARCH" = "x86_64" ]; then \
		ARCH="amd64"; \
	elif [ "$$ARCH" = "aarch64" ]; then \
		ARCH="arm64"; \
	fi; \
	BINARY="sitedog-$$PLATFORM-$$ARCH"; \
	echo "Installing $$BINARY as sitedog-dev..."; \
	sudo ln -sf $(PWD)/dist/$$BINARY /usr/local/bin/sitedog-dev; \
	echo "Development version installed successfully!"

dev-uninstall:
	sudo rm -f /usr/local/bin/sitedog-dev

dev-install!: dev-uninstall dev-install

uninstall:
	scripts/uninstall.sh

reinstall: uninstall install

.DEFAULT_GOAL := help