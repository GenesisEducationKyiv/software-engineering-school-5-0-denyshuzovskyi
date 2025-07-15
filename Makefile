LINT_CONFIG=./../.golangci.yaml

.PHONY: sync docker-build-subs

sync:
	go work sync

docker-build-subs:
	docker build -f weather-upd-subscription-srv/Dockerfile -t weather-upd-subscription-srv .

docker-build-notification:
	docker build -f notification-srv/Dockerfile -t notification-srv .