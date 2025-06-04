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
	rm -rf sitedog_gist
	git clone git@gist.github.com:fe278d331980a1ce09c3d946bbf0b83b.git sitedog_gist; \
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