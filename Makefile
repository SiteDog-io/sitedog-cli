.PHONY: run help push push-install

help:
	@echo "Available commands:"
	@echo "  run    - Run Ruby script sitedog.rb"
	@echo "  help   - Show this help message"
	@echo "  push   - Update files in gist"
	@echo "  push-install - TODO: put install.sh to get.sitedog.io"

run:
	ruby sitedog.rb

push:
	@if [ ! -d "sitedog_gist" ]; then \
		git clone git@gist.github.com:a85deab4772d1c825602ea64e0c035bc.git sitedog_gist; \
	fi
	cp sitedog.rb demo.html.erb sitedog_gist/
	cd sitedog_gist && \
	if git diff --quiet; then \
		echo "No changes to deploy"; \
	else \
		git add . && \
		git commit -m "Update sitedog.rb and demo.html.erb" && \
		git push; \
	fi

push-install:
	# TODO: put install.sh to get.sitedog.io

.DEFAULT_GOAL := help