.PHONY: run help

help:
	@echo "Available commands:"
	@echo "  run   - Run Ruby script sitedog.rb"
	@echo "  help  - Show this help message"

run:
	ruby sitedog.rb

.DEFAULT_GOAL := help