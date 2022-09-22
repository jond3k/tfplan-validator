#!/usr/bin/env bash
set -euf -o pipefail

cd ./infrastructure/dev1
terraform plan -out ./plan.bin

