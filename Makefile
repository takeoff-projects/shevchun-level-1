SHELL := /bin/bash

docker-build:
	docker build \
		-f Dockerfile \
		-t go-website:latest \
		.

gcloud-build: 
	gcloud builds submit --tag=gcr.io/roi-takeoff-user51/go-website:v1.8 .

deploy:
	terraform -chdir="./terraform" apply --auto-approve -var="project_id=roi-takeoff-user51"

destory:
	terraform -chdir="./terraform" destroy --auto-approve -var="project_id=roi-takeoff-user51"
