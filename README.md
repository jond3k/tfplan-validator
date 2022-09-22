# tfplan-validator

[![codecov](https://codecov.io/gh/fautom/tfplan-validator/branch/main/graph/badge.svg?token=1P6A5WBXOT)](https://codecov.io/gh/fautom/tfplan-validator)
![build](https://github.com/fautom/tfplan-validator/actions/workflows/test.yaml/badge.svg)

A simple way to validate Terraform plans. Designed to assist batch operations on large numbers of similar state files.

## Getting started

Run `go install github.com/fautom/tfplan-validator@latest` to use from the command line. You can now use commands like `tfplan-validator create` and `tfplan-validator check`

Read the below example section or see [test/multi-env](test/multi-env) for working code that uses local state files that you can run from your command line.

To use `tfplan-validator` as a library check the command implementations in [internal/app/tfplan-validator](internal/app/tfplan-validator)


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

So let's first generate a plan for `dev1`

    ➜ cd environments/dev1
    ➜ terraform plan -out plan.bin

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

This is a change we want and we'd like to apply it not only to `dev1` but `dev2` and `dev3` so we create a rules file to block any unwanted changes to other environments.

    ➜ terraform show -json plan.bin > plan.json
    ➜ tfplan-validator create plan.json --rules rules.json

    Created rules file rules.json that allows Terraform to perform the following actions:

    - module.environment.local_file.bar[0] can be created

Now we can generate plans for every other environment and validate them, starting with `dev1` and `dev2`

    ➜ tfplan-validator check --rules rules.json environments/dev1/plan.json environments/dev2/plan.json

    The plan environments/dev1/plan.json passes checks and will perform the following actions:

    - module.environment.local_file.bar[0] will be created

    The plan environments/dev2/plan.json passes checks and will perform the following actions:

    - module.environment.local_file.bar[0] will be created

Since both pass the plans can now be safely applied

    ➜ terraform apply -auto-approve environments/dev1/plan.bin
    ➜ terraform apply -auto-approve environments/dev2/plan.bin

However `dev3` has a resource that would be destroyed by its plan so `tfplan-validator check` rejects it with a non-zero status code

    ➜ tfplan-validator check --rules rules.json environments/dev3/plan.json

    The plan environments/dev3/plan.json has been rejected because it has the following actions:

    - module.environment.local_file.foo[0] cannot be deleted

It is now up to the developer to decide what to do:

1. Add the delete action to the rules by using `tfplan-validator merge` or `tfplan-validator create`
1. Fix the code so that the resource is not deleted
1. Apply the change anyway, destroying the resource

See [test/multi-env](test/multi-env) for a working example using local state files that you can run from your command line.

## Related projects

* [GoogleCloudPlatform/terraform-validator](https://github.com/GoogleCloudPlatform/terraform-validator)

