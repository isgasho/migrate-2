# convert-workflows tool

This tool will generate a set of GitHub Actions YAML workflow files from your `.workflow` files.

## Conversion notes

The original beta version of Actions supported parallel Action execution while sharing a workspace. In V2, we do support parallel execution
of jobs, but jobs do not share workspaces.

