.PHONY: help push push-install-prod

help:
	@echo "Available commands:"
	@echo "  help   - Show this help message"
	@echo "  push   - Update files in gist"
	@echo "  push-install-prod - TODO: put install.sh to get.sitedog.io"

push:
	rm -rf sitedog_gist
	git clone git@gist.github.com:fe278d331980a1ce09c3d946bbf0b83b.git sitedog_gist; \
	cp sitedog.rb demo.html.erb install.sh uninstall.sh sitedog_gist/
	cd sitedog_gist && \
	if git diff --quiet; then \
		echo "No changes to deploy"; \
	else \
		git add . && \
		git commit -m "Update sitedog.rb and demo.html.erb" && \
		git push; \
	fi

push-install-prod:
	# TODO: put install.sh to get.sitedog.io

.DEFAULT_GOAL := help