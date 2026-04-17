SHELL := /bin/bash

.PHONY: dev-web dev-api build-web build-api gen gen-apply

dev-web:
	npm run dev:web

dev-api:
	npm run dev:api

build-web:
	npm run build:web

build-api:
	npm run build:api

gen:
	npm run module:plan -- --spec $(SPEC)

gen-apply:
	npm run module:apply -- --spec $(SPEC)

