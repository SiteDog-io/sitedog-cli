.PHONY: help push push-install-prod build-docker version

help:
	@echo "Available commands:"
	@echo "  help                - Show this help message"
	@echo "  build               - Build all Go binaries for all platforms (in Docker, output to ./dist)"
	@echo "  push                - Update files in gist (binaries from ./dist, install/uninstall scripts, etc.)"
	@echo "  push!               - build + push"
	@echo "  push-install-prod   - TODO: put install.sh to get.sitedog.io"
	@echo "  version             - Update version in main.go and create git tag"

push:
	rm -rf sitedog_gist
	git clone git@gist.github.com:fe278d331980a1ce09c3d946bbf0b83b.git sitedog_gist
	rm -rf sitedog_gist/*
	cp dist/* sitedog_gist/
	cp demo.html.erb scripts/install.sh scripts/uninstall.sh sitedog_gist/
	cd sitedog_gist && \
	if git diff --quiet; then \
		echo "No changes to deploy"; \
	else \
		git add . && \
		git commit -m "Update sitedog files" && \
		git push; \
	fi

push-install-prod:
	# TODO: put install.sh to get.sitedog.io

build:
	docker run --rm -v $(PWD):/app -w /app golang:1.20-alpine sh -c "./scripts/build.sh"

push!: build push

version:
	@if [ -z "$(v)" ]; then \
		echo "Usage: make version v=x.y.z"; \
		exit 1; \
	fi; \
	file=main.go; \
	ver=$(v); \
	sed -i "s/Version[ ]*=[ ]*\"[^"]*\"/Version = \"$$ver\"/" $$file; \
	git add $$file; \
	git commit -m "bump version to $$ver"; \
	git tag $$ver; \
	echo "Version updated to $$ver and git tag created."

.DEFAULT_GOAL := help