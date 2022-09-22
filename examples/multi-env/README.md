# tfplan-validator multi-env example

This folder contains a working example of plan validation where all environments are out of date with the code in `modules/environment`. 

Checking `infrastructure/dev1` and `infrastructure/dev2` with `rules.json` succeeds as they have identical signatures but `infrastructure/dev3` fails as updates will lead to unwanted changes, in this case deletions. 

In summary:

1. You can run `./example.sh` to see how a rule is created and checked on 3 similar local state files
1. The script will modify the local state files for `dev1` and `dev2`but reject changes to `dev3` as it has unexpected changes
1. Running a `git diff` will show you the applied changes to hte local state file
1. The rule file is already checked in as `./rules.json` 
