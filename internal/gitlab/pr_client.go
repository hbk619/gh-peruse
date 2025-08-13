package gitlab

import (
	"errors"
	"fmt"
	"maps"
	"path/filepath"
	"strconv"

	"github.com/hbk619/gh-peruse/internal/git"
	"github.com/hbk619/gh-peruse/internal/gitlab/graphql"
	"github.com/hbk619/gh-peruse/internal/requests"
	gitlab "gitlab.com/gitlab-org/api/client-go"
)

type PRClient struct {
	apiClient gitlab.GraphQLInterface
	repo      Repo
	runner    requests.CommandLine
}

func NewPRClient(apiClient gitlab.GraphQLInterface, repo Repo, runner requests.CommandLine) *PRClient {
	return &PRClient{
		apiClient: apiClient,
		repo:      repo,
		runner:    runner,
	}
}

func (gh *PRClient) GetCommentCountForOwnedPRs(repo *git.Repo) (map[int]int, error) {
	return nil, errors.New("not implemented")

}

func (gh *PRClient) GetPRDetails(repo *git.Repo, verbose bool) (*git.PR, error) {
	res := &GitlabResponse{}
	_, err := gh.apiClient.Do(graphql.GitlabPRQuery(repo), res)
	if err != nil {
		return nil, err
	}

	return &git.PR{
		Comments: gh.createComments(&res.Data.Project.MergeRequest),
		Title:    res.Data.Project.MergeRequest.Title,
		State:    gh.createState(verbose, &res.Data.Project.MergeRequest),
		Id:       res.Data.Project.MergeRequest.Id,
	}, nil
}

func (gh *PRClient) createComments(response *PR) []git.Comment {
	var commentList []git.Comment

	if response.Description != "" {
		commentList = append(commentList, git.Comment{
			Author: git.Author{Login: response.Author.Name},
			Body:   response.Description,
			File: git.File{
				FullPath: git.MainThread,
				FileName: git.MainThread,
			},
			CreatedAt: response.CreatedAt,
		})
	}

	discussionComments := map[string][]git.Comment{}
	for _, comment := range response.Notes.Nodes {
		if comment.Body != "" && !comment.System {
			localComment := git.Comment{
				Author: git.Author{Login: comment.Author.Name},
				Body:   comment.Body,
				File: git.File{
					FullPath: git.MainThread,
					FileName: git.MainThread,
				},
				CreatedAt: comment.CreatedAt,
				Thread: git.Thread{
					ID:         comment.Discussion.ReplyId,
					IsResolved: comment.Discussion.Resolved,
				},
			}
			if comment.Position.FilePath != "" {
				lineNumber := comment.Position.NewLine
				if lineNumber == 0 {
					lineNumber = comment.Position.OldLine
				}
				localComment.File.Path = comment.Position.FilePath
				localComment.File.FullPath = fmt.Sprintf("%s:%d", localComment.File.Path, lineNumber)
				localComment.File.FileName = filepath.Base(localComment.File.Path)
				localComment.File.Line = lineNumber
				discussionComments[comment.Discussion.ReplyId] = append(discussionComments[comment.Discussion.ReplyId], localComment)
			} else {
				commentList = append(commentList, localComment)
			}
		}
	}

	git.SortCommentsInPlace(commentList)
	for comments := range maps.Values(discussionComments) {
		git.SortCommentsInPlace(comments)
		commentList = append(commentList, comments...)
	}

	return commentList
}

var mergeStatuses = map[string]string{
	"UNCHECKED":                    "The mergeability of the pull request is still being calculated",
	"CHECKING":                     "The state cannot currently be determined",
	"MERGEABLE":                    "Mergeable",
	"COMMITS_STATUS":               "Branch exists with commits",
	"CI_MUST_PASS":                 "Pipeline must succeed before merging.",
	"CI_STILL_RUNNING":             "Pipeline is still running.",
	"DISCUSSIONS_NOT_RESOLVED":     "Discussions must be resolved before merging.",
	"DRAFT_STATUS":                 "Merge request must not be draft before merging.",
	"NOT_OPEN":                     "Merge request must be open before merging.",
	"NOT_APPROVED":                 "Merge request must be approved before merging.",
	"BLOCKED_STATUS":               "Merge request dependencies must be merged.",
	"EXTERNAL_STATUS_CHECKS":       "Status checks must pass.",
	"PREPARING":                    "Merge request diff is being created.",
	"JIRA_ASSOCIATION":             "Either the title or description must reference a Jira issue.",
	"CONFLICT":                     "Merge conflicts",
	"NEED_REBASE":                  "Merge request needs to be rebased.",
	"APPROVALS_SYNCING":            "Merge request approvals currently syncing.",
	"LOCKED_PATHS":                 "Merge request includes locked paths.",
	"LOCKED_LFS_FILES":             "Merge request includes locked LFS files.",
	"MERGE_TIME":                   "Merge request may not be merged until after the specified time.",
	"SECURITY_POLICIES_VIOLATIONS": "All policy rules must be satisfied.",
	"TITLE_NOT_MATCHING":           "Merge request title does not match required regex.",
	"REQUESTED_CHANGES":            "Indicates a reviewer has requested changes.",
}

func (gh *PRClient) createState(verbose bool, prDetails *PR) git.State {
	state := git.State{}
	if verbose {
		var statuses []git.Status
		for _, pipeline := range prDetails.Pipelines.Nodes {
			for _, job := range pipeline.Jobs.Nodes {
				statuses = append(statuses, git.Status{
					Name:       job.Name,
					Conclusion: job.Status,
				})
			}
		}
		reviewStatus := gh.getReviewStatuses(prDetails)
		state = git.State{
			Statuses:    statuses,
			MergeStatus: mergeStatuses[prDetails.DetailedMergeStatus],
			Reviews:     reviewStatus,
		}
	}
	return state
}

func (gh *PRClient) getReviewStatuses(response *PR) map[string][]string {
	reviewStatus := make(map[string][]string)
	alreadySeenReviewers := make(map[string]bool)
	for _, approver := range response.ApprovedBy.Nodes {
		reviewStatus["APPROVED"] = append(reviewStatus["APPROVED"], approver.Name)
		alreadySeenReviewers[approver.Name] = true
	}
	for _, reviewer := range response.Reviewers.Nodes {
		if !alreadySeenReviewers[reviewer.Name] {
			reviewStatus["PENDING"] = append(reviewStatus["PENDING"], reviewer.Name)
		}
	}

	return reviewStatus
}

func (gh *PRClient) DetectCurrentPR(repo *git.Repo) (int, error) {
	ref, err := gh.runner.Run("git", []string{"symbolic-ref", "--quiet", "--short", "HEAD"})
	if err != nil {
		return 0, err
	}
	res := &GitlabResponse{}
	_, err = gh.apiClient.Do(graphql.GitlabPRForBranch(ref, repo), res)
	if err != nil {
		return 0, err
	}

	if len(res.Data.Project.MergeRequests.Nodes) == 0 {
		return 0, fmt.Errorf("no merge request found for %s", ref)
	}

	if len(res.Data.Project.MergeRequests.Nodes) > 1 {
		return 0, fmt.Errorf("too many merge requests found for %s", ref)
	}

	number, err := strconv.Atoi(res.Data.Project.MergeRequests.Nodes[0].Iid)
	if err != nil {
		return 0, err
	}
	return number, nil

}

func (gh *PRClient) GetRepoDetails() (*git.Repo, error) {
	return gh.repo.Base("origin")
}

func (gh *PRClient) Reply(contents string, comment *git.Comment, prId string) error {
	res := &gitlab.GenericGraphQLErrors{}
	_, err := gh.apiClient.Do(graphql.AddNote(prId, contents, comment.Thread.ID), res)
	if err != nil {
		return err
	}

	if len(res.Errors) > 0 {
		msg := ""
		for _, errorDetails := range res.Errors {
			msg += errorDetails.Message
		}
		return errors.New(msg)
	}

	return nil
}

func (gh *PRClient) Resolve(comment *git.Comment) error {
	res := &gitlab.GenericGraphQLErrors{}
	_, err := gh.apiClient.Do(graphql.ResolveDiscussion(comment.Thread.ID), res)
	if err != nil {
		return err
	}

	if len(res.Errors) > 0 {
		msg := ""
		for _, errorDetails := range res.Errors {
			msg += errorDetails.Message
		}
		return errors.New(msg)
	}

	return nil
}
