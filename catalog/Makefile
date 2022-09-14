SHELL   := /bin/bash

all: \
	generate-add-resource \
	generate-remove-resource \
	generate-replace-resource

generate-%:
	cd ${*} && \
	terraform init && \
	terraform plan -out ./plan && \
	terraform show ./plan -no-color > plan.txt && \
	terraform show -json ./plan | jq '.' > plan.json

