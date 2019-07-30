

workflow "args" {
  on = "push"
  resolves = [
    "array args",
  ]
}

action "string args" {
  uses = "docker://alpine"
  args = "echo hi"
}

action "array args" {
  needs = ["string args"]
  uses = "docker://alpine"
  args = ["echo", "hi"]
}


workflow "secrets wf" {
  on = "push"
  resolves = [
    "secrets",
  ]
}

action "secrets" {
  uses = "docker://alpine"
  args = "echo hi"
  secrets = ["SECRET"]
}

workflow "much parallel" {
  on = "push"
  resolves = [
    "E",
    "F",
  ]
}

action "A" {
  uses = "docker://alpine"
}
action "B" {
  uses = "docker://alpine"
  needs = ["A"]
}
action "C" {
  uses = "docker://alpine"
  needs = ["A"]
}
action "D" {
  uses = "docker://alpine"
  needs = ["B", "C"]
}
action "E" {
  uses = "docker://alpine"
  needs = ["D"]
}
action "F" {
  uses = "docker://alpine"
  needs = ["D"]
}

workflow "   $.yml annoYthin✅ name.js.yml" {
  on = "push"
  resolves = [
    "one"
  ]
}

action "one" {
  uses = "docker://alpine"
}

workflow "Top 5" {
  resolves = ["Create an issue"]
  on = "schedule(00 08 * * 1)"
}

action "Create an issue" {
  uses = "bdougie/create-an-issue@22c7dda012351b778b392d257999fb782cb44041"
  secrets = ["GITHUB_TOKEN"]
  args = ".github/ISSUE_TEMPLATE/TOP5.md"
}

workflow "   nonense stress test :key value {k:1,}" {
  on = "push"
  resolves = [
    "emoji actions 🙈 yay    ' : "
  ]
}

action "emoji actions 🙈 yay    ' : " {
  uses = "docker://alpine"
}

workflow "no resolves" {
	on = "push"
}