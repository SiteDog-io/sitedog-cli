.PHONY: help push push-install-prod build

help:
	@echo "Available commands:"
	@echo "  help   - Show this help message"
	@echo "  push   - Update files in gist"
	@echo "  build - Build Go binary"
	@echo "  install - Install Go binary globally (to ~/.local/bin or /usr/local/bin)"
	@echo "  uninstall - Remove Go binary from system"

push:
	rm -rf sitedog_gist
	git clone git@gist.github.com:fe278d331980a1ce09c3d946bbf0b83b.git sitedog_gist; \
	cp sitedog install.sh uninstall.sh demo.html.erb sitedog_gist/
	cd sitedog_gist && \
	if git diff --quiet; then \
		echo "No changes to deploy"; \
	else \
		git add . && \
		git commit -m "Update sitedog files" && \
		git push; \
	fi

build:
	go build -o sitedog.bin main.go

push-install-prod:
	# TODO: put install.sh to get.sitedog.io

.DEFAULT_GOAL := help