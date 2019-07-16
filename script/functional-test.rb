#!/usr/bin/env ruby
require 'fileutils'
require 'minitest/autorun'


describe "functional test" do
  before do
    FileUtils.mkdir_p(".github")

    File.write(".github/main.workflow", %q(
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
))

    `./bin/migrate-actions`
  end


  it "outputs expected files" do
    actual = Dir[".github/workflows/*.yml"]
    assert_equal actual.sort, [".github/workflows/schedule.yml", ".github/workflows/push.yml"].sort
  end

  it "handles push correctly" do
    push = File.read(".github/workflows/push.yml")

    assert_equal %q(on: push
name: workflow one
jobs:
  actionOne:
    name: action one
    steps:
    - name: action one
      uses: docker://alpine
      entrypoint: sh -c echo ${{ github.sha }}
), push
  end

  it "handles schedule correctly" do
    schedule = File.read(".github/workflows/schedule.yml")

    assert_equal %q(on:
  schedules:
  - cron: '* * * * *'
name: per minute
jobs:
  actionOne:
    name: action one
    steps:
    - name: action one
      uses: docker://alpine
      entrypoint: sh -c echo ${{ github.sha }}
), schedule
  end
end
