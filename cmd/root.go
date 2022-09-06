package cmd

import (
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/cli/go-gh"
	"github.com/cli/go-gh/pkg/api"
	"github.com/spf13/cobra"
)

type cmdFlags struct {
	reposExclude []string
	reportFile   string
	debug        bool
}

func NewCmd() *cobra.Command {

	// Instantiate struct to contain values from cobra flags; arguments are handled within RunE
	cmdFlags := cmdFlags{}

	// Instantiate cobra command driving work from package
	// Closures are used for cobra command lifecycle hooks for access to cobra flags struct
	cmd := cobra.Command{
		Use:   "gh-dependency-report [flags] owner [repo ...]",
		Short: "Generate report of repository manifests and dependencies discovered through the dependency graph",
		Long:  "Generate report of repository manifests and dependencies discovered through the dependency graph",
		Args:  cobra.MinimumNArgs(1),
		// Setup command lifecycle handler; cmd representing the cobra.Command being instantiated now
		RunE: func(cmd *cobra.Command, args []string) error {

			var err error
			var client api.GQLClient

			// Reinitialize logging if debugging was enabled
			if cmdFlags.debug {
				logger, _ := log.NewLogger(cmdFlags.debug)
				defer logger.Sync() // nolint:errcheck // not sure how to errcheck a deferred call like this
				zap.ReplaceGlobals(logger)
			}

			client, err = gh.GQLClient(&api.ClientOptions{
				Headers: map[string]string{
					"Accept": "application/vnd.github.hawkgirl-preview+json",
				},
			})

			if err != nil {
				zap.S().Errorf("Error arose retrieving graphql client")
				return err
			}

			owner := args[0]
			repos := args[1:]

			if _, err := os.Stat(cmdFlags.reportFile); errors.Is(err, os.ErrExist) {
				return err
			}

			reportWriter, err := os.OpenFile(cmdFlags.reportFile, os.O_WRONLY|os.O_CREATE, 0644)

			if err != nil {
				return err
			}

			return runCmd(owner, repos, cmdFlags.reposExclude, newAPIGetter(client), reportWriter)
		},
	}

	// Determine default report file based on current timestamp; for more info see https://pkg.go.dev/time#pkg-constants
	reportFileDefault := fmt.Sprintf("report-%s.csv", time.Now().Format("20060102150405"))

	// Configure flags for command
	cmd.Flags().StringSliceVarP(&cmdFlags.reposExclude, "exclude", "e", []string{}, "Repositories to exclude from report")
	cmd.Flags().StringVarP(&cmdFlags.reportFile, "output-file", "o", reportFileDefault, "Name of file to write CSV report")
	cmd.PersistentFlags().BoolVarP(&cmdFlags.debug, "debug", "d", false, "Whether to debug logging")

	return &cmd
}

func runCmd(owner string, repos []string, reposExclude []string, apiGetter *APIGetter, reportWriter *os.File) error {
	fmt.Printf('Dylan was here')
}
type APIGetter struct {
	client api.GQLClient
}

func newAPIGetter(client api.GQLClient) *APIGetter {
	return &APIGetter{
		client: client,
	}
}

// client, err := gh.RESTClient(nil)
// if err != nil {
// 	fmt.Println(err)
// 	return
// }
// response := struct{ Login string }{}
// err = client.Get("user", &response)
// if err != nil {
// 	fmt.Println(err)
// 	return
// }
// fmt.Printf("running as %s\n", response.Login)

// GraphQL example
// client, err := gh.GQLClient(nil)
// if err != nil {
// 	log.Fatal(err)
// }
// var query struct {
// 	Repository struct {
// 		Refs struct {
// 			Nodes []struct {
// 				Name string
// 			}
// 		} `graphql:"refs(refPrefix: $refPrefix, last: $last)"`
// 	} `graphql:"repository(owner: $owner, name: $name)"`
// }
// variables := map[string]interface{}{
// 	"refPrefix": graphql.String("refs/tags/"),
// 	"last":      graphql.Int(30),
// 	"owner":     graphql.String("cli"),
// 	"name":      graphql.String("cli"),
// }

// var query struct {
// 	Enterprise struct {

// }

// err = client.Query("RepositoryTags", &query, variables)
// if err != nil {
// 	log.Fatal(err)
// }
// fmt.Println(query)
