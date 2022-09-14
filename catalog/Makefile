SHELL   := /bin/bash

all: \
	generate-create \
	generate-delete \
	generate-delete-create

generate-%:
	cd ${*} && \
	terraform init && \
	terraform plan -out plan.bin && \
	terraform show plan.bin -no-color > plan.txt && \
	terraform show -json plan.bin | jq '.' > plan.json

