SHELL   := /bin/bash

all: \
	generate-create \
	generate-delete \
	generate-delete-create

generate-%:
	cd ${*} && \
	terraform init && \
	terraform plan -out ./plan && \
	terraform show ./plan -no-color > plan.txt && \
	terraform show -json ./plan | jq '.' > plan.json

