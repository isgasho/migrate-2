workflow "Send Message On Push" {
  on = "push"
  resolves = [
    "run",
    "new-task",
  ]
}

action "run" {
  uses = "./.github/actions/irc/"
}

action "new-task" {
  uses = "docker://alpine"
  args = "cat /github/workflow/event.json"
}