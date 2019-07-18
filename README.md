# convert-workflows tool

This tool will generate a set of GitHub Actions YAML workflow files from your `.workflow` files.

## Conversion notes

The original beta version of Actions supported parallel Action execution while sharing a workspace. In V2, we do support parallel execution
of jobs, but jobs do not share workspaces.

## Install

Head over to the [releases](https://github.com/actions/migrate/releases) tab, and download the executable for your operating system.

Once you've downloaded it, navigate to a repository using Actions V1 and run the `migrate-actions` executable.

```sh
> cd path/to/your/repo
> cat ./.github/main.workflow
workflow "on push" {
    on = "push"
    resolves = ["say hi"]
}
workflow "on pull request" {
    on = "pull_request"
    resolves = ["say hi"]
}
action "say hi" {
    uses = "docker/whalesay@master"
    runs = "whalesay hello actions"
}
> path/to/migrate
Created workflow .github/workflows/push.yml
Created workflow .github/workflows/pull_request.yml
> tail -n +1 .github/workflows/*.yml
==> .github/workflows/pull_request.yml <==
on: pull_request
name: on pull request
jobs:
  sayHi:
    name: say hi
    steps:
    - name: say hi
      uses: docker/whalesay@master
      entrypoint: whalesay hello actions

==> .github/workflows/push.yml <==
on: push
name: on push
jobs:
  sayHi:
    name: say hi
    steps:
    - name: say hi
      uses: docker/whalesay@master
      entrypoint: whalesay hello actions
```

You can then delete your `main.workflow`. If you have any `.yml` files in `.github/workflows` your `main.workflow` will be ignored.
