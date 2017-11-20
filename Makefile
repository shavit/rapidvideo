build:
	docker build -t itstommy/rv-inspector .

start_dev:
	docker run --rm \
		--env-file ${PWD}/.env \
		-v ${PWD}:/go/src/github.com/shavit/rapidvideo \
		-ti itstommy/rv-inspector
