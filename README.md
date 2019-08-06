# migrate-actions tool

Converts GitHub Actions `main.workflow` files into the new `.yml` syntax. Pre-built binaries are available for Linux, OSX and Windows, and it can be built for any environment that Go supports.

## Install

Head over to the [releases](https://github.com/actions/migrate/releases) tab, and download the archive for your operating system:

- Windows: migrate-actions-windows.zip
- Linux: migrate-actions-linux.tar
- OSX: migrate-actions-osx.tar

and unarchive it.

## Usage

Once you've downloaded and unarchived the tool, navigate to a repository using Actions via a `main.workflow` file and run the `migrate-actions` executable (`migrate-actions.exe` on windows),
locating it by the full path to the executable e.g `~/Downloads/migrate-actions-osx/migrate-actions` on OSX.

Given an existing `.github/main.workflow`:

```hcl
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
    runs = "cowsay hello actions"
}
```

Running `migrate-actions`:

```
Created workflow .github/workflows/push.yml
Created workflow .github/workflows/pull_request.yml
```

This will produce two GitHub Actions V2 YAML configuration files.
First, it will produce a configuration for pull requests, `.github/workflows/pull_request.yml`:

```yaml
on: pull_request
name: on pull request
jobs:
  sayHi:
    name: say hi
    runs-on: ubuntu-latest
    steps:
    - uses: actions/checkout@master
    - name: say hi
      uses: docker/whalesay@master
      with:
        entrypoint: cowsay
        args: hello actions
```

Next, it will produce a configuration file for pushes, `.github/workflows/push.yml`:

```yaml
on: push
name: on push
jobs:
  sayHi:
    name: say hi
    runs-on: ubuntu-latest
    steps:
    - uses: actions/checkout@master
    - name: say hi
      uses: docker/whalesay@master
      with:
        entrypoint: cowsay
        args: hello actions
```

You can then delete your `main.workflow`. If you have any `.yml` files in `.github/workflows` your `main.workflow` will be ignored.

## Limitations

- Comments are dropped, sorry! 😭
- Parallel workflows are serialized (see Conversion notes)

## Conversion notes

The original beta version of Actions supported parallel Action execution while sharing a workspace. In V2, we do support parallel execution of jobs, but jobs do not share workspaces.

This means parallel workflows will be linearized:

```
      ┌───────┐             ┌───────┐
      │   A   │             │   A   │
      └───────┘             └───────┘
          │                     │    
    ┌─────┴─────┐               ▼    
    │           │           ┌───────┐
    ▼           ▼           │   C   │
┌───────┐   ┌───────┐       └───────┘
│   B   │   │   C   │           │    
└───────┘   └───────┘           ▼    
    │           │           ┌───────┐
    └─────┬─────┘           │   B   │
          ▼                 └───────┘
      ┌───────┐                 │    
      │   D   │                 ▼    
      └───────┘             ┌───────┐
          │                 │   D   │
    ┌─────┴─────┐           └───────┘
    │           │               │    
    ▼           ▼               ▼    
┌───────┐   ┌───────┐       ┌───────┐
│   E   │   │   F   │       │   F   │
└───────┘   └───────┘       └───────┘
                                │    
                                ▼    
                            ┌───────┐
                            │   E   │
                            └───────┘
```

Since `jobs` can be parallel, if there isn't a hard dependency on shared workspaces, you can regain paralleism by duplicating work. For instance, we can make E and F run in parallel by creating two jobs, both linearizing the dependencies leading up to the final action:

<table>
      <tr>
            <th>Workflow file</th>
            <th>Resulting build</th>
      </tr>
      <tr> 
           <td>
           
```
on: push
jobs:
  E:
    runs-on: ubuntu-latest
    steps:
    - run: ./A
    - run: ./B
    - run: ./C
    - run: ./D
    - run: ./E
  F:
    runs-on: ubuntu-latest
    steps:
    - run: ./A
    - run: ./B
    - run: ./C
    - run: ./D
    - run: ./F
```

</td>

<td>
      
```
┌───────┐     ┌───────┐
│   A   │     │   A   │
└───────┘     └───────┘
    │             │    
    ▼             ▼    
┌───────┐     ┌───────┐
│   C   │     │   C   │
└───────┘     └───────┘
    │             │    
    ▼             ▼    
┌───────┐     ┌───────┐
│   B   │     │   B   │
└───────┘     └───────┘
    │             │    
    ▼             ▼    
┌───────┐     ┌───────┐
│   D   │     │   D   │
└───────┘     └───────┘
    │             │    
    ▼             ▼    
┌───────┐     ┌───────┐
│   E   │     │   F   │
└───────┘     └───────┘
```

</td>

</tr>

</table>

## Contributing

You don't need the following to use the tool, but if you'd like to contribute here's how.

### Build

- Prerequesites: Go 1.12.x, `dep`
- Clone project
- Run `./script/bootstrap`
- Run `./script/build`

### Releasing

- Make your change following the [contribution guide](CONTRIBUTING.md)
- Once you've merged your PR to master run `./script/release` locally, and CI will create a new release with the binaries

