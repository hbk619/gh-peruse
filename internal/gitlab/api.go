package gitlab

import (
	"time"

	"github.com/hbk619/gh-peruse/internal/git"
)

type GLComments struct {
	Nodes []GLComment
}
type Discussion struct {
	ReplyId    string
	Id         string
	Resolved   bool
	Resolvable bool
}
type GLFile struct {
	FilePath string
	NewLine  int
	OldLine  int
}

type GLComment struct {
	git.Comment
	System     bool
	Author     GLAuthor
	Discussion Discussion
	Position   GLFile
}

type GLAuthor struct {
	Name string
}

type Approvers struct {
	Nodes []GLAuthor
}

type Reviewers struct {
	Nodes []GLAuthor
}

type Job struct {
	Name   string
	Status string
}

type Jobs struct {
	Nodes []Job
}
type Pipeline struct {
	Status string
	Jobs   Jobs
}
type Pipelines struct {
	Nodes []Pipeline
}

type PR struct {
	Id                  string
	Iid                 string
	Title               string
	Author              GLAuthor
	Notes               GLComments
	Description         string
	CreatedAt           time.Time
	ApprovedBy          Approvers
	Reviewers           Reviewers
	Conflicts           bool
	State               string
	DetailedMergeStatus string
	Pipelines           Pipelines
}

type MergeRequests struct {
	Nodes []PR
}

type GitlabProject struct {
	MergeRequest  PR
	MergeRequests MergeRequests
}

type GitlabData struct {
	Project GitlabProject
}

type GitlabResponse struct {
	Data GitlabData
}
