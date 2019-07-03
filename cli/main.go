package main

import (
	"flag"
	"fmt"
	"github.com/actions/migrate"
	"io/ioutil"
	"os"
	"path"
)

const workflowFilePath = ".github/main.workflow"
const workflowDirectory = ".github/workflows"

const usage = `./actions-migrate`

func main() {

	helpFlag := flag.Bool("help", false, "outputs help")
	flag.Parse()

	if *helpFlag {
		flag.Usage()
		os.Exit(0)
	}

	// find root
	rootDir := "."
	workflowFile := path.Join(rootDir, workflowFilePath)

	f, err := os.Open(workflowFile)
	if err != nil {
		userError(fmt.Sprintf("No `%s' file to convert", workflowFilePath))
		return
	}

	err = ensureDirectory()
	if err != nil {
		failed(fmt.Sprintf("Failed to create directory: %s", err))
	}

	converted, err := migrate.Parse(f)
	if err != nil {
		failed(fmt.Sprintf("Failed to convert workflow file: %s", err.Error()))
		return
	}

	files, err := converted.Files()
	if err != nil {
		failed(fmt.Sprintf("Failed to convert workflow file: %s", err.Error()))
		return
	}

	for _, file := range files {
		writeFile(file)
	}
}

func writeFile(file migrate.OutputFile) {
	err := ioutil.WriteFile(file.Path, []byte(file.Content), 0644)
	if err != nil {
		failed(fmt.Sprintf("Failed to write `%s'", file.Path))
	}
}

func ensureDirectory() error {
	return os.MkdirAll(workflowDirectory, 0755)
}

func failed(msg string) {
	// TODO report
	exitWithMessage(msg)
}

func userError(msg string) {
	exitWithMessage(msg)
}

func exitWithMessage(msg string) {
	fmt.Fprintln(os.Stderr, msg)
	os.Exit(1)

}
