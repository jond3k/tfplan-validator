SHELL   := /bin/bash
TFPV    := ../../bin/tfplan_validator

all: \
	gen-create \
	gen-delete \
	gen-delete-create \
	gen-create-delete \
	gen-update

gen-%:
	cd ${*} && \
	terraform init && \
	terraform plan -out plan.bin && \
	terraform show plan.bin -no-color > plan.txt && \
	terraform show -json plan.bin | jq '.' > plan.json && \
	$(TFPV) create plan.json filter.json

