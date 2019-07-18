
workflow "workflow one" {
  on = "push"
  resolves = [
    "action one",
  ]
}

workflow "per minute" {
  on = "schedule(* * * * *)"
  resolves = [
    "action one",
  ]
}

action "action one" {
  uses = "docker://alpine"
  runs = ["sh", "-c", "echo $GITHUB_SHA"]
}
