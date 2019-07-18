# Build Process

action "Install Dependencies" {
	uses = "actions/npm@master"
	args = "install"
}

action "Lint Code" {
	needs = ["Install Dependencies"]
	uses = "actions/npm@master"
	args = "run lint"
}

action "Build Project" {
	needs = ["Lint Code"]
	uses = "actions/npm@master"
	args = "run build"
}

action "Docker Build" {
	needs = ["Build Project"]
	uses = "actions/docker/cli@master"
	args = "build -t roryclaasen/reactchess ."
}

# Heroku

action "Login Heroku" {
	uses = "actions/heroku@master"
	args = "container:login"
	secrets = ["HEROKU_API_KEY"]
}

# Test Workflow

workflow "Develop Test, Build And Deploy" {
	on = "push"
	resolves = [
		"Docker Build",
		"Release Heroku Develop",
	]
}

action "Filter Develop" {
	needs = ["Docker Build"]
	uses = "actions/bin/filter@master"
	args = "branch develop"
}

action "Push Heroku Develop" {
	needs = ["Filter Develop", "Login Heroku"]
	uses = "actions/heroku@master"
	args = "container:push -a react-chessgame-dev web"
	secrets = ["HEROKU_API_KEY"]
}

action "Release Heroku Develop" {
	needs = "Push Heroku Develop"
	uses = "actions/heroku@master"
	args = "container:release -a react-chessgame-dev web"
	secrets = ["HEROKU_API_KEY"]
}

# Production Workflow

workflow "Production Develop Build And Deploy" {
	on = "push"
	resolves = [
		"Release Heroku Production"
	]
}

action "Filter Production" {
	needs = ["Docker Build"]
	uses = "actions/bin/filter@master"
	args = "branch master"
}

action "Push Heroku Production" {
	needs = ["Filter Production", "Login Heroku"]
	uses = "actions/heroku@master"
	args = "container:push -a react-chessgame web"
	secrets = ["HEROKU_API_KEY"]
}

action "Release Heroku Production" {
	needs = ["Push Heroku Production"]
	uses = "actions/heroku@master"
	args = "container:release -a react-chessgame web"
	secrets = ["HEROKU_API_KEY"]
}