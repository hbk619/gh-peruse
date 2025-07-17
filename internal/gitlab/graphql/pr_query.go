package graphql

import (
	"fmt"

	"github.com/hbk619/gh-peruse/internal/git"
	gitlab "gitlab.com/gitlab-org/api/client-go"
)

func GitlabPRQuery(repo *git.Repo) gitlab.GraphQLQuery {
	query := fmt.Sprintf(`
query {
  project(fullPath: "%s/%s") {
    mergeRequest(iid: "%d") {
      id
      reviewers {
        nodes {
          name

        }
      }
      approvedBy {
        nodes {
          name
        }
      }
      approvalsLeft
      approvalsRequired
      approved
      pipelines {
        nodes {
          jobs {
            nodes {
              name
              status
            }
          }
        }
      }
      title
      createdAt
      conflicts
      state
      detailedMergeStatus
      description
      author {
        name
      }
      notes {
        nodes {
          author {
            name
          }
          body
          createdAt
          discussion {
            id
            replyId
            resolved
            resolvable
          }
          position {
            filePath
            newLine
            oldLine
          }
          resolvable
          resolved
          system     
        }
      }
    }
  }
}
	`, repo.Owner, repo.Name, repo.PRNumber)

	return gitlab.GraphQLQuery{
		Query: query,
	}
}

func GitlabPRForBranch(branch string, repo *git.Repo) gitlab.GraphQLQuery {
	query := fmt.Sprintf(`
query {
  project(fullPath: "%s/%s") {
    mergeRequests(sourceBranches: "%s") {
      nodes {
        id
        iid
      }
    }
  }
}
	`, repo.Owner, repo.Name, branch)

	return gitlab.GraphQLQuery{
		Query: query,
	}
}
