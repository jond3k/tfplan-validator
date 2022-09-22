# tfplan-validator

[![codecov](https://codecov.io/gh/fautom/tfplan-validator/branch/main/graph/badge.svg?token=1P6A5WBXOT)](https://codecov.io/gh/fautom/tfplan-validator)
![build](https://github.com/fautom/tfplan-validator/actions/workflows/test.yaml/badge.svg)

A simple way to validate Terraform plans. Designed to assist batch operations on large numbers of similar state files.

## Example

Suppose we have multiple development environments all created from the same Terraform module as our production environment. This means developers do not interfere with each others work and deployments are predictable but over time these environments will diverge due to changes to branches, variables, or state drift from manual changes.

    environments/
                 dev1/
                 dev2/
                 dev3/
                 prod/

    modules/
            environment/

A simple solution to this is to regularly run a `terraform apply` or a `terragrunt run-all apply` but this could be quite time consuming if the plans need to be manually reviewed and we risk overriding intentional differences if `-auto-approve` is used, including destructive changes like deleting databases. Instead let's modify a single development environment by hand and then use the plan for that to validate proposed changes to other environments.

First, let's generate a plan which we can manually inspect.

    > cd ../environments/dev1
    > terraform plan -out ./plan

    Terraform will perform the following actions:

      # local_file.foo will be created
      + resource "local_file" "foo" {
          + content              = "foo!"
          + directory_permission = "0777"
          + file_permission      = "0777"
          + filename             = "./foo.bar"
          + id                   = (known after apply)
        }

    Plan: 1 to add, 0 to change, 0 to destroy.

If we like the results we can create a validator that will only accept plans with this create operation. The validator currently only accepts plans in json format. 

    > terraform show -json ./plan > ./plan.json
    > tfplan-validator create ./plan.json ../rules.json

    Created rules file ../rules.json that allows Terraform to perform the following actions:

    - local_file.foo can be created

Now we can safely auto-approve the other plans knowing that the validator will reject other changes.

    set -e
    for dev_env in dev1 dev2 dev3;
    do
      cd ../${dev_env}/
      terragrunt plan -out ./plan
      terragrunt show ./plan > ./plan.json
      tfplan-validator check ./plan.json ../rules.json && \
      terraform apply -auto-approve ./plan
    done

## Related projects

* [GoogleCloudPlatform/terraform-validator](https://github.com/GoogleCloudPlatform/terraform-validator)

