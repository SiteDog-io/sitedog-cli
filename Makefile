.PHONY: run help deploy

help:
	@echo "Available commands:"
	@echo "  run    - Run Ruby script sitedog.rb"
	@echo "  help   - Show this help message"
	@echo "  deploy - Update files in gist"

run:
	ruby sitedog.rb

deploy:
	@if [ ! -d "sitedog_gist" ]; then \
		git clone git@gist.github.com:fe278d331980a1ce09c3d946bbf0b83b.git sitedog_gist; \
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

.DEFAULT_GOAL := help