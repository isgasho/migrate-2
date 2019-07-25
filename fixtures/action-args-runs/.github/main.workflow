workflow "args and runs" {
  on = "push"
  resolves = [
    "args",
    "args runs",
    "args array",
    "args runs array",
  ]
}

action "args" {
  uses = "docker://alpine"
  args = "echo hi"
}

action "args runs" {
  uses = "docker://alpine"
  args = "echo hi"
  runs = "some thing"
}

action "args array" {
  uses = "docker://alpine"
  args = ["echo","hi"]
}

action "args runs array" {
  uses = "docker://alpine"
  args = ["echo","hi"]
  runs = ["some","thing"]
}
