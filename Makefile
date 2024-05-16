build:
	cd src/api && npm install

run:
	cd src/api && npm start

re: build run

.PHONY: build run re
