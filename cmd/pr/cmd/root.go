package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/cli/cli/v2/git"
	"github.com/cli/go-gh/v2/pkg/api"
	"github.com/hbk619/gh-peruse/cmd/pr/internal"
	common "github.com/hbk619/gh-peruse/internal"
	"github.com/hbk619/gh-peruse/internal/filesystem"
	"github.com/hbk619/gh-peruse/internal/github"
	"github.com/hbk619/gh-peruse/internal/gitlab"
	"github.com/hbk619/gh-peruse/internal/history"
	internal_os "github.com/hbk619/gh-peruse/internal/os"
	"github.com/hbk619/gh-peruse/internal/requests"
	"github.com/spf13/cobra"

	gitlabApi "gitlab.com/gitlab-org/api/client-go"
)

var PRCmd = &cobra.Command{
	Use:   "pr [number]",
	Args:  cobra.MaximumNArgs(1),
	Short: "Browse Github PR comments",
	Long:  `View comments from a PR one by one and reply to them`,
	Run: func(cmd *cobra.Command, args []string) {
		historyService, err := history.NewHistoryService(os.Getenv("HOME"), filesystem.NewFS())
		if err != nil {
			fmt.Println(err)
			return
		}
		runner := requests.NewCommandRunner()
		repo := gitlab.NewGLRepo(runner)

		gitClient := &git.Client{}
		var prClient github.PullRequestClient

		if strings.ToLower(os.Getenv("PERUSE_TYPE")) == "gitlab" {
			prClient, err = createGitlabClient(repo, runner)
			if err != nil {
				fmt.Println(err)
				return
			}
		} else {
			graphQlClient, err := api.DefaultGraphQLClient()
			if err != nil {
				fmt.Println(err)
				return
			}
			prClient = github.NewPRClient(graphQlClient, gitClient)
		}
		clipboard := internal_os.NewClipboard()
		output := filesystem.NewStdOut()
		prompt := common.NewPrompt(os.Stdin, output)
		pr := internal.NewPRAction(prClient, historyService, output, clipboard, prompt)
		verbose, err := cmd.Flags().GetBool("verbose")
		if err != nil {
			fmt.Println(err)
			return
		}
		err = pr.Init(args, verbose)
		if err != nil {
			fmt.Println(err)
			return
		}
		pr.Run()
	},
}

func createGitlabClient(repo gitlab.Repo, runner requests.CommandLine) (github.PullRequestClient, error) {
	token := os.Getenv("GITLAB_TOKEN")
	var c *gitlabApi.Client

	c, err := gitlabApi.NewClient(token)
	if err != nil {
		return nil, err
	}
	return gitlab.NewPRClient(c.GraphQL, repo, runner), nil
}

func Execute() {
	err := PRCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	PRCmd.AddCommand(CheckCommentCountCmd)
	PRCmd.Flags().BoolP("verbose", "v", false, "Verbose mode")
}
