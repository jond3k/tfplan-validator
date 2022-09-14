SHELL   := /bin/bash

all: \
	gen-create \
	gen-delete \
	gen-delete-create \
	gen-create-delete

gen-%:
	cd ${*} && \
	terraform init && \
	terraform plan -out plan.bin && \
	terraform show plan.bin -no-color > plan.txt && \
	terraform show -json plan.bin | jq '.' > plan.json

