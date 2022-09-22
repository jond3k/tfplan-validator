#!/usr/bin/env bash
set -euf -o pipefail

echo "First let's run a plan against the infrastructure"
echo "[press enter to continue..]"
read

pushd ./infrastructure/dev1
terraform plan -out ./plan.bin

echo
echo
echo "Now let's create a rules file from the plan.."
echo "[press enter to continue..]"
read

terraform show -json ./plan.bin > ./plan.json
tfplan-validator create ./plan.json --rules ../../rules.json

echo
popd
echo

echo "Using the rules we can now iterate through the environments and validate the plans for each"
echo "Note that we should repeat dev1 again because things may have changed since the creation of the rules"
echo "[press enter to continue..]"
read

for dev_env in dev1 dev2 dev3;
do
  pushd ./infrastructure/${dev_env}
  echo
  terraform plan -out ./plan.bin > /dev/null
  terraform show -json ./plan.bin > ./plan.json
  tfplan-validator check ./plan.json --rules ../../rules.json || continue
  echo
  terraform apply -auto-approve ./plan.bin
  popd
done
