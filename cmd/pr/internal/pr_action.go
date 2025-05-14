package internal

import (
	"errors"
	"fmt"
	"github.com/hbk619/gh-peruse/internal"
	"github.com/hbk619/gh-peruse/internal/filesystem"
	"github.com/hbk619/gh-peruse/internal/git"
	"github.com/hbk619/gh-peruse/internal/github"
	"github.com/hbk619/gh-peruse/internal/history"
	"os"
	"strconv"
	"strings"
)

type PRAction struct {
	Id                  string
	Repo                *git.Repo
	Results             []git.Comment
	PrintedPathLastTime bool
	LastFullPath        string
	HelpText            string
	State               git.State
	client              github.PullRequestClient
	history             history.Storage
	output              filesystem.Output
	internal.Interactive
}

func NewPRAction(client github.PullRequestClient, history history.Storage, output filesystem.Output) *PRAction {
	return &PRAction{
		Repo:                &git.Repo{},
		PrintedPathLastTime: true,
		LastFullPath:        "",
		HelpText:            "Type c to comment",
		client:              client,
		history:             history,
		output:              output,
	}
}

func (pr *PRAction) Init(args []string, verbose bool) error {
	repoDetails, err := pr.client.GetRepoDetails()
	if err != nil {
		return err
	}
	pr.Repo.Owner = repoDetails.Owner
	pr.Repo.Name = repoDetails.Name

	pr.Repo.PRNumber, err = pr.getPRNumber(args)
	if err != nil {
		return err
	}

	prDetails, err := pr.client.GetPRDetails(pr.Repo, verbose)
	if err != nil {
		return err
	}
	pr.Results = prDetails.Comments
	pr.State = prDetails.State
	pr.Id = prDetails.Id
	if verbose {
		pr.PrintState()
	}

	commentCount := len(pr.Results)
	pr.updateHistory(pr.Repo.PRNumber, commentCount)

	if commentCount == 0 {
		return errors.New("no comments found")
	}

	pr.Interactive.MaxIndex = commentCount - 1
	pr.Print()
	return nil
}

func (pr *PRAction) getPRNumber(args []string) (int, error) {
	if len(args) > 0 {
		number, err := strconv.Atoi(args[0])
		if err != nil {
			return 0, fmt.Errorf("please provide a valid PR number")
		}
		return number, nil
	} else {
		return pr.client.DetectCurrentPR(pr.Repo)
	}
}

func (pr *PRAction) updateHistory(prNumber int, commentCount int) {
	prHistory, err := pr.history.Load()
	if err != nil {
		pr.output.Println(fmt.Sprintf("Warning failed to load comments to history: %s", err.Error()))
	}

	if err == nil {
		existingPrHistory := prHistory.Prs[prNumber]
		if existingPrHistory.CommentCount != commentCount {
			pr.output.Println("New comments ahead!")
		}

		existingPrHistory.CommentCount = commentCount
		prHistory.Prs[prNumber] = existingPrHistory
		err = pr.history.Save(prHistory)
		if err != nil {
			pr.output.Println(fmt.Sprintf("Warning failed to save comments to history: %s", err.Error()))
		}
	}
}

func (pr *PRAction) Reply(contents string) {
	err := pr.client.Reply(contents, &pr.Results[pr.Interactive.Index], pr.Id)
	if err != nil {
		pr.output.Println(fmt.Sprintf("Warning failed to comment: %s", err.Error()))
	} else {
		pr.output.Println("Posted comment")
	}
}

func (pr *PRAction) Resolve() {
	err := pr.client.Resolve(&pr.Results[pr.Interactive.Index])
	if err != nil {
		pr.output.Println(fmt.Sprintf("Warning failed to resolve thread: %s", err.Error()))
	} else {
		pr.output.Println("Conversation resolved")
	}
}

func (pr *PRAction) Run() {
	for {
		prompt := "n to go to the next result, p for previous, r to repeat or q to quit"
		currentComment := pr.Results[pr.Interactive.Index]
		pr.LastFullPath = currentComment.File.FullPath
		if currentComment.Thread.ID != "" && !currentComment.Thread.IsResolved {
			prompt += ", res to resolve"
		}
		if currentComment.Thread.IsResolved || currentComment.Outdated {
			prompt += ", e to expand"
		}
		result := internal.StringPrompt(prompt)
		switch result {
		case "n":
			pr.Interactive.Next(pr.Print)
		case "p":
			pr.Interactive.Previous(pr.Print)
		case "r":
			pr.Interactive.Repeat(pr.Print)
		case "e":
			pr.LastFullPath = ""
			pr.printContents(currentComment)
		case "res":
			pr.Resolve()
		case "c":
			comment := internal.StringPrompt("Type comment and press enter")
			pr.Reply(comment)
		case "q":
			os.Exit(0)
		default:
			pr.output.Println("Invalid choice")
		}

	}
}

func (pr *PRAction) Print() {
	current := pr.Results[pr.Interactive.Index]
	if current.Thread.IsResolved {
		pr.output.Println("This comment is resolved")
		return
	}
	if current.Outdated {
		pr.output.Println("This comment is outdated")
		return
	}
	pr.printContents(current)
}

func (pr *PRAction) printContents(current git.Comment) {
	if pr.LastFullPath != current.File.FullPath {
		pr.output.Print(current.File.FileName)
		if current.File.Path != "" {
			pr.output.Print(current.File.Path)
			pr.output.Print(strconv.Itoa(current.File.Line))
			pr.output.Print(current.File.LineContents)
		}
	}
	pr.output.Print(current.Author.Login)
	pr.output.Print(current.Body)
}

func (pr *PRAction) PrintState() {
	pr.output.Print(pr.State.MergeStatus)
	pr.output.Print(pr.State.ConflictStatus)
	for reviewState, names := range pr.State.Reviews {
		pr.output.Print(fmt.Sprintf("%s %s", reviewState, strings.Join(names, " ")))
	}
	for _, status := range pr.State.Statuses {
		pr.output.Print(fmt.Sprintf("Check %s %s", status.Name, status.Conclusion))
	}
}
