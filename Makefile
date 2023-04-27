pkg_name = $(shell basename $$PWD)

.SHELLFLAGS = -ec
.ONESHELL:

.PHONY: build
build:
	docker build --ssh default -t ${pkg_name} .

.PHONY: push
push:
	tag=$$(git describe --tags --exact)

	docker tag ${pkg_name} europe-west1-docker.pkg.dev/tmaas-dev-dev/images/${pkg_name}:$${tag}
	docker push europe-west1-docker.pkg.dev/tmaas-dev-dev/images/${pkg_name}:$${tag}

.PHONY: doc
doc:
	swag init --pd

.PHONY: run
run:
	docker run --rm -it \
		-p 8080 \
		-e PORT=8080 \
		-e GOOGLE_CLOUD_PROJECT=tmaas-dev-dev \
		-e PROJECT_ID=tmaas-dev-dev \
		-e PUBSUB_USER_CHANGES_TOPIC=user-changes \
		-e SERVICE_ACCOUNT_NAME=${pkg_name} \
		-e SERVICE_ACCOUNT_PERMISSIONS='{}' \
		-e TEST_ACCOUNT_EMAIL_REGEXES='^$$' \
		-e LOG_LEVEL=debug \
		-e CORS_CONFIG='{"AllowOrigins":["*"]}' \
		-v $$HOME/.config/gcloud/application_default_credentials.json:/creds.json \
		-e GOOGLE_APPLICATION_CREDENTIALS=/creds.json \
		${pkg_name}

.PHONY: run-local
run-local:
	go build \
	&& \
		PORT=8080 \
		GOOGLE_CLOUD_PROJECT=tmaas-dev-dev \
		PROJECT_ID=tmaas-dev-dev \
		PUBSUB_USER_CHANGES_TOPIC=user-changes \
		SERVICE_ACCOUNT_NAME=${pkg_name} \
		SERVICE_ACCOUNT_PERMISSIONS='{}' \
		TEST_ACCOUNT_EMAIL_REGEXES='^$$' \
		LOG_LEVEL=debug \
		CORS_CONFIG='{"AllowOrigins":["*"]}' \
		GOOGLE_APPLICATION_CREDENTIALS=$$HOME/.config/gcloud/application_default_credentials.json \
		./tmaas-${pkg_name}

.PHONY: debug-local
debug-local:
	go build \
	&& \
		PORT=8080 \
		GOOGLE_CLOUD_PROJECT=tmaas-dev-dev \
		PROJECT_ID=tmaas-dev-dev \
		PUBSUB_USER_CHANGES_TOPIC=user-changes \
		SERVICE_ACCOUNT_NAME=${pkg_name} \
		SERVICE_ACCOUNT_PERMISSIONS='{}' \
		TEST_ACCOUNT_EMAIL_REGEXES='^$$' \
		LOG_LEVEL=debug \
		CORS_CONFIG='{"AllowOrigins":["*"]}' \
		GOOGLE_APPLICATION_CREDENTIALS=$$HOME/.config/gcloud/application_default_credentials.json \
		dlv debug
