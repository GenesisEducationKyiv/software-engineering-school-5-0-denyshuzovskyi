.PHONY: sync
sync:
	go work sync

.PHONY: docker-build-subs
docker-build-subs:
	docker build -f weather-upd-subscription-srv/Dockerfile -t weather-upd-subscription-srv .

.PHONY: docker-build-notification
docker-build-notification:
	docker build -f notification-srv/Dockerfile -t notification-srv .