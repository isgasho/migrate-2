package commands

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path"

	"github.com/actions/migrate/converter"
	"github.com/spf13/cobra"
)

const workflowFilePath = ".github/main.workflow"
const workflowDirectory = ".github/workflows"

var rootCmd = &cobra.Command{
	Use:   "migrate-actions",
	Short: "CLI for migrating Actions main.workflow files to the new YAML syntax.",
	Run: func(cmd *cobra.Command, args []string) {
		// Do logic here.
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

		converted, err := converter.Parse(f)
		if err != nil {
			failed(fmt.Sprintf("Failed to convert workflow file: %s", err.Error()))
			return
		}

		files, err := converted.Files()
		if err != nil {
			failed(fmt.Sprintf("Failed to convert workflow file: %s", err.Error()))
			return
		}

		if len(files) == 0 {
			failed("Could not find any workflows to convert")
			return
		}

		for _, file := range files {
			writeFile(file)
			fmt.Printf("Created workflow %s\n", file.Path)
		}

		fmt.Println("")
		fmt.Println("You can now delete your main.workflow file. If you have any .yml files in .github/workflows your main.workflow will be ignored.")
		fmt.Println("    rm .github/main.workflow")
		fmt.Println("    git add -A .github/main.workflow .github/workflows/*.yml")
		fmt.Println("    git commit -m 'converted main.workflow to Actions V2 yml files'")
		fmt.Println("")

		fmt.Println("Thanks for being a")
		fmt.Println(banner)
	},
}

// Execute executes the root command of converting the main.workflow file.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func writeFile(file converter.OutputFile) {
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

const banner = `  ________.__  __     ___ ___      ___.        _____          __  .__                      
 /  _____/|__|/  |_  /   |   \ __ _\_ |__     /  _  \   _____/  |_|__| ____   ____   ______
/   \  ___|  \   __\/    ~    \  |  \ __ \   /  /_\  \_/ ___\   __\  |/  _ \ /    \ /  ___/
\    \_\  \  ||  |  \    Y    /  |  / \_\ \ /    |    \  \___|  | |  (  <_> )   |  \\___ \ 
 \______  /__||__|   \___|_  /|____/|___  / \____|__  /\___  >__| |__|\____/|___|  /____  >
        \/                 \/           \/          \/     \/                    \/     \/ 
__________        __           ___________              __                                 
\______   \ _____/  |______    \__    ___/___   _______/  |_  ___________                  
 |    |  _// __ \   __\__  \     |    |_/ __ \ /  ___/\   __\/ __ \_  __ \                 
 |    |   \  ___/|  |  / __ \_   |    |\  ___/ \___ \  |  | \  ___/|  | \/                 
 |______  /\___  >__| (____  /   |____| \___  >____  > |__|  \___  >__|                    
        \/     \/          \/               \/     \/            \/        `
