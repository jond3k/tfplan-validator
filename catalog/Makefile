SHELL   := /bin/bash

all: \
	generate-add-resource \
	generate-remove-resource \
	generate-update-resource

generate-%:
	cd ${*} ; \
	terraform init ; \
	terraform plan -out ./plan ; \
	terraform show ./plan > plan.txt ; \
	terraform show -json ./plan | jq '.' > plan.json

