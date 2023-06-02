.PHONY: all build fa lint test vendor

all: build

build:
	npx vite build

fa:
	@echo 'import { library } from "@fortawesome/fontawesome-svg-core";' > src/lib/font-awesome.js
	@find . -name "*.vue" -exec grep "fa-icon" {} \; \
		| sed -n 's/.*icon="\([^"]*\)".*/\1/p' \
		| sort -u \
		| awk -F"-" '{for(i=1;i<=NF;i++){$$i=toupper(substr($$i,1,1)) substr($$i,2)}} 1' OFS="" \
		| while read icon; do \
				echo ; \
				echo "import { fa$${icon} } from \"@fortawesome/free-solid-svg-icons\";" ; \
				echo "library.add(fa$${icon});" ; \
			done >> src/lib/font-awesome.js

lint:
	npm run lint

test:
	npm run test:unit

vendor:
	npm install
