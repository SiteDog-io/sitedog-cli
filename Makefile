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
	@echo "  test                - Run all detector tests"
	@echo "  stats               - Show detector statistics"
	@echo "  stats!              - Show detailed detector information"
	@echo "  detectors           - List all available detectors"

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

uninstall:
	scripts/uninstall.sh

reinstall: uninstall install

test:
	cd tests && go run run_tests.go

stats:
	cd tests && go run run_tests.go | grep -E "(Testing|Found)" | awk '/Testing/ {name=$$0} /Found/ {print name " - " $$0}' | sed 's/=== Testing \(.*\) Detector ===/\1:/' | sed 's/ - Found \([0-9]*\) services:/ → \1 services/'

stats!:
	cd tests && go run run_tests.go | grep -E "(Testing|Detector:|Description:|Should run:|Found|^  [0-9]+\.)" | sed 's/=== Testing \(.*\) Detector ===/\n🔍 \1/' | sed 's/Detector: \(.*\)/   Name: \1/' | sed 's/Description: \(.*\)/   Description: \1/' | sed 's/Should run: \(.*\)/   Status: \1/' | sed 's/Found \([0-9]*\) services:/   Services (\1):/' | sed 's/^  \([0-9]*\)\. \(.*\)/     • \2/'

detectors:
	@echo "📋 All Available Detectors:"
	@echo ""
	@make stats!

.DEFAULT_GOAL := help