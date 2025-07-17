package gitlab

import (
	"encoding/json"
	"errors"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/hbk619/gh-peruse/internal"
	"github.com/hbk619/gh-peruse/internal/git"
	"github.com/hbk619/gh-peruse/internal/gitlab/graphql"
	mock_gitlab "github.com/hbk619/gh-peruse/internal/gitlab/mocks"
	mock_requests "github.com/hbk619/gh-peruse/internal/requests/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	gitlab "gitlab.com/gitlab-org/api/client-go"
)

type GitlabClientTestSuite struct {
	suite.Suite
	mockGraphQL     *mock_gitlab.MockGraphQLInterface
	mockRepo        *mock_gitlab.MockRepo
	mockCommandLine *mock_requests.MockCommandLine
	ctrl            *gomock.Controller
	repo            *git.Repo
	prService       *PRClient
}

func (suite *GitlabClientTestSuite) BeforeTest(string, string) {
	suite.ctrl = gomock.NewController(suite.T())
	suite.mockGraphQL = mock_gitlab.NewMockGraphQLInterface(suite.ctrl)
	suite.mockRepo = mock_gitlab.NewMockRepo(suite.ctrl)
	suite.mockCommandLine = mock_requests.NewMockCommandLine(suite.ctrl)
	suite.repo = &git.Repo{
		Owner:    "luigi",
		Name:     "castle",
		PRNumber: 123,
	}
	suite.prService = NewPRClient(
		suite.mockGraphQL,
		suite.mockRepo,
		suite.mockCommandLine,
	)
}

func (suite *GitlabClientTestSuite) TestPRService_getPrDetails_error() {
	expected := errors.New("failed to graphql")
	suite.mockGraphQL.EXPECT().
		Do(graphql.GitlabPRQuery(suite.repo), gomock.Any(), gomock.Any()).
		DoAndReturn(func(query gitlab.GraphQLQuery, gr any, _ ...any) (gitlab.Response, error) {
			return gitlab.Response{}, expected
		})
	details, err := suite.prService.GetPRDetails(suite.repo, false)
	assert.ErrorContains(suite.T(), err, expected.Error())
	assert.Nil(suite.T(), details)
}

func (suite *GitlabClientTestSuite) TestPRService_getPrDetails_no_comments() {
	prDetails := `{
  "data": {
    "project": {
      "mergeRequests": {
        "nodes": []
      },
      "mergeRequest": {
        "id": "gid://gitlab/MergeRequest/11111111",
        "reviewers": {
          "nodes": []
        },
        "approvedBy": {
          "nodes": []
        },
        "state": "opened",
        "mergeable": true,
        "mergeableDiscussionsState": true,
        "mergeStatusEnum": "CAN_BE_MERGED",
        "detailedMergeStatus": "MERGEABLE",
        "conflicts": false,
        "approvalsLeft": 0,
        "approvalsRequired": 0,
        "approved": false,
        "pipelines": {
          "nodes": []
        },
        "title": "Test pr",
        "createdAt": "2025-06-26T14:16:54Z",
        "description": "",
        "author": {
          "name": "Mario"
        },
        "notes": null
      }
    }
  },
  "correlationId": "2423kn3kjn4kj32n42kjn3k4n"
}`

	expected := &git.PR{
		Comments: nil,
		State:    git.State{},
		Title:    "Test pr",
		Id:       "gid://gitlab/MergeRequest/11111111",
	}

	suite.mockGraphQL.EXPECT().
		Do(graphql.GitlabPRQuery(suite.repo), gomock.Any(), gomock.Any()).
		DoAndReturn(func(query gitlab.GraphQLQuery, gr any, _ ...any) (gitlab.Response, error) {

			err := json.Unmarshal([]byte(prDetails), &gr)
			suite.NoError(err)
			return gitlab.Response{}, nil
		})
	details, err := suite.prService.GetPRDetails(suite.repo, false)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), expected, details)
}

func (suite *GitlabClientTestSuite) TestPRService_getPrDetails_comments() {
	prDetails := `{
  "data": {
    "project": {
      "mergeRequests": {
        "nodes": []
      },
      "mergeRequest": {
        "id": "gid://gitlab/MergeRequest/11111111",
        "reviewers": {
          "nodes": []
        },
        "approvedBy": {
          "nodes": []
        },
        "state": "opened",
        "mergeable": true,
        "mergeableDiscussionsState": true,
        "mergeStatusEnum": "CAN_BE_MERGED",
        "detailedMergeStatus": "MERGEABLE",
        "conflicts": false,
        "approvalsLeft": 0,
        "approvalsRequired": 0,
        "approved": false,
        "pipelines": {
          "nodes": []
        },
        "title": "Test pr",
        "createdAt": "2025-06-26T14:16:54Z",
        "description": "",
        "author": {
          "name": "Mario"
        },
        "notes": {
          "nodes": [
            {
              "author": {
                "name": "Luigi"
              },
              "body": "assigned to @hbk619",
              "createdAt": "2025-06-26T14:16:54Z",
              "discussion": {
                "id": "gid://gitlab/Discussion/mvcxnfi23938",
                "replyId": "gid://gitlab/IndividualNoteDiscussion/mvcxnfi23938",
                "resolved": false,
                "resolvable": false
              },
              "position": null,
              "resolvable": false,
              "resolved": false,
              "system": true
            },
            {
              "author": {
                "name": "Wario"
              },
              "body": "but why??",
              "createdAt": "2025-06-24T11:17:27Z",
              "discussion": {
                "id": "gid://gitlab/Discussion/adsa89u92iajdsamd",
                "replyId": "gid://gitlab/DiffDiscussion/adsa89u92iajdsamd",
                "resolved": false,
                "resolvable": true
              },
              "position": {
                "filePath": "test.md",
                "newLine": 8,
                "oldLine": 2,
                "positionType": "text"
              },
              "resolvable": true,
              "resolved": false,
              "system": false
            },
            {
              "author": {
                "name": "Luigi"
              },
              "body": "single comment",
              "createdAt": "2025-06-26T14:17:06Z",
              "discussion": {
                "id": "gid://gitlab/Discussion/asd9980adsasdmklm",
                "replyId": "gid://gitlab/Discussion/asd9980adsasdmklm",
                "resolved": false,
                "resolvable": true
              },
              "position": null,
              "resolvable": true,
              "resolved": false,
              "system": false
            },
            {
              "author": {
                "name": "Peach"
              },
              "body": "'cause it's cool",
              "createdAt": "2025-06-24T14:19:27Z",
              "discussion": {
                "id": "gid://gitlab/Discussion/adsa89u92iajdsamd",
                "replyId": "gid://gitlab/DiffDiscussion/adsa89u92iajdsamd",
                "resolved": false,
                "resolvable": true
              },
              "position": {
                "filePath": "test.md",
                "newLine": 8,
                "oldLine": 2,
                "positionType": "text"
              },
              "resolvable": true,
              "resolved": false,
              "system": false
            },
            {
              "author": {
                "name": "Toad"
              },
              "body": "delete delete delete",
              "createdAt": "2025-06-21T14:19:27Z",
              "discussion": {
                "id": "gid://gitlab/Discussion/mklmlk32mlk",
                "replyId": "gid://gitlab/DiffDiscussion/mklmlk32mlk",
                "resolved": false,
                "resolvable": true
              },
              "position": {
                "filePath": "test.md",
                "oldLine": 2,
                "positionType": "text"
              },
              "resolvable": true,
              "resolved": false,
              "system": false
            },
            {
              "author": {
                "name": "Luigi"
              },
              "body": "a review comment",
              "createdAt": "2025-06-26T14:17:27Z",
              "discussion": {
                "id": "gid://gitlab/Discussion/sad09099009i90n3243",
                "replyId": "gid://gitlab/Discussion/sad09099009i90n3243",
                "resolved": false,
                "resolvable": true
              },
              "position": null,
              "resolvable": true,
              "resolved": false,
              "system": false
            },
            {
              "author": {
                "name": "Luigi"
              },
              "body": "lovely",
              "createdAt": "2025-06-26T14:17:27Z",
              "discussion": {
                "id": "gid://gitlab/Discussion/dsoijierew22ukc",
                "replyId": "gid://gitlab/DiffDiscussion/dsoijierew22ukc",
                "resolved": false,
                "resolvable": true
              },
              "position": {
                "filePath": "test.md",
                "newLine": 3,
                "oldLine": null,
                "positionType": "text"
              },
              "resolvable": true,
              "resolved": false,
              "system": false
            },
            {
              "author": {
                "name": "Mario"
              },
              "body": "thanks I think so too",
              "createdAt": "2025-06-26T14:19:27Z",
              "discussion": {
                "id": "gid://gitlab/Discussion/dsoijierew22ukc",
                "replyId": "gid://gitlab/DiffDiscussion/dsoijierew22ukc",
                "resolved": false,
                "resolvable": true
              },
              "position": {
                "filePath": "test.md",
                "newLine": 3,
                "oldLine": null,
                "positionType": "text"
              },
              "resolvable": true,
              "resolved": false,
              "system": false
            },
            {
              "author": {
                "name": "Luigi"
              },
              "body": "left review comments",
              "createdAt": "2025-06-26T14:17:27Z",
              "discussion": {
                "id": "gid://gitlab/Discussion/vuiew93ejklm",
                "replyId": "gid://gitlab/IndividualNoteDiscussion/vuiew93ejklm",
                "resolved": false,
                "resolvable": false
              },
              "position": null,
              "resolvable": false,
              "resolved": false,
              "system": true
            },
            {
              "author": {
                "name": "Luigi"
              },
              "body": "second single comment",
              "createdAt": "2025-06-26T14:47:22Z",
              "discussion": {
                "id": "gid://gitlab/Discussion/sad0908n32432jnk2n4838",
                "replyId": "gid://gitlab/IndividualNoteDiscussion/sad0908n32432jnk2n4838",
                "resolved": false,
                "resolvable": false
              },
              "position": null,
              "resolvable": false,
              "resolved": false,
              "system": false
            },
            {
              "author": {
                "name": "Luigi"
              },
              "body": "this is a commit comment",
              "createdAt": "2025-06-26T14:58:07Z",
              "discussion": {
                "id": "gid://gitlab/Discussion/asd909i309ikm",
                "replyId": "gid://gitlab/IndividualNoteDiscussion/asd909i309ikm",
                "resolved": false,
                "resolvable": false
              },
              "position": null,
              "resolvable": false,
              "resolved": false,
              "system": false
            },
            {
              "author": {
                "name": "Luigi"
              },
              "body": "added 1 commit\n\n<ul><li>b39ad224 - second commit</li></ul>\n\n[Compare with previous version](/hbk619-group/hbk619-project/-/merge_requests/1/diffs?diff_id=1403811057&start_sha=e4cd115966cc3347614cce819e416ab53da8b5b6)",
              "createdAt": "2025-06-26T15:29:56Z",
              "discussion": {
                "id": "gid://gitlab/Discussion/sdjao99jcm",
                "replyId": "gid://gitlab/IndividualNoteDiscussion/sdjao99jcm",
                "resolved": false,
                "resolvable": false
              },
              "position": null,
              "resolvable": false,
              "resolved": false,
              "system": true
            }
		]
		}
      }
    }
  },
  "correlationId": "2423kn3kjn4kj32n42kjn3k4n"
}`

	expected := &git.PR{
		Comments: []git.Comment{
			{
				Author: git.Author{
					Login: "Luigi",
				},
				Body: "single comment",
				File: git.File{
					FullPath: git.MainThread,
					FileName: git.MainThread,
				},
				Thread: git.Thread{
					ID:         "gid://gitlab/Discussion/asd9980adsasdmklm",
					IsResolved: false,
				},
				CreatedAt: internal.TimeMustParse(time.RFC3339, "2025-06-26T14:17:06Z"),
			},
			{
				Author: git.Author{
					Login: "Luigi",
				},
				Body: "a review comment",
				File: git.File{
					FullPath: git.MainThread,
					FileName: git.MainThread,
				},
				Thread: git.Thread{
					ID:         "gid://gitlab/Discussion/sad09099009i90n3243",
					IsResolved: false,
				},
				CreatedAt: internal.TimeMustParse(time.RFC3339, "2025-06-26T14:17:27Z"),
			},
			{
				Author: git.Author{
					Login: "Luigi",
				},
				Body: "second single comment",
				File: git.File{
					FullPath: git.MainThread,
					FileName: git.MainThread,
				},
				CreatedAt: internal.TimeMustParse(time.RFC3339, "2025-06-26T14:47:22Z"),
				Thread: git.Thread{
					ID:         "gid://gitlab/IndividualNoteDiscussion/sad0908n32432jnk2n4838",
					IsResolved: false,
				},
			},
			{
				Author: git.Author{
					Login: "Luigi",
				},
				Body: "this is a commit comment",
				File: git.File{
					FullPath: git.MainThread,
					FileName: git.MainThread,
				},
				Thread: git.Thread{
					ID:         "gid://gitlab/IndividualNoteDiscussion/asd909i309ikm",
					IsResolved: false,
				},
				CreatedAt: internal.TimeMustParse(time.RFC3339, "2025-06-26T14:58:07Z"),
			},
			{
				Author: git.Author{
					Login: "Wario",
				},
				Body:      "but why??",
				CreatedAt: internal.TimeMustParse(time.RFC3339, "2025-06-24T11:17:27Z"),
				File: git.File{
					Path:     "test.md",
					FullPath: "test.md:8",
					Line:     8,
					FileName: "test.md",
				},
				Thread: git.Thread{
					ID:         "gid://gitlab/DiffDiscussion/adsa89u92iajdsamd",
					IsResolved: false,
				},
			},
			{
				Author: git.Author{
					Login: "Peach",
				},
				Body:      "'cause it's cool",
				CreatedAt: internal.TimeMustParse(time.RFC3339, "2025-06-24T14:19:27Z"),
				File: git.File{
					Path:     "test.md",
					FullPath: "test.md:8",
					Line:     8,
					FileName: "test.md",
				},
				Thread: git.Thread{
					ID:         "gid://gitlab/DiffDiscussion/adsa89u92iajdsamd",
					IsResolved: false,
				},
			},
			{
				Author: git.Author{
					Login: "Toad",
				},
				Body:      "delete delete delete",
				CreatedAt: internal.TimeMustParse(time.RFC3339, "2025-06-21T14:19:27Z"),
				File: git.File{
					Path:     "test.md",
					FullPath: "test.md:2",
					Line:     2,
					FileName: "test.md",
				},
				Thread: git.Thread{
					ID:         "gid://gitlab/DiffDiscussion/mklmlk32mlk",
					IsResolved: false,
				},
			},
			{
				Author: git.Author{
					Login: "Luigi",
				},
				Body:      "lovely",
				CreatedAt: internal.TimeMustParse(time.RFC3339, "2025-06-26T14:17:27Z"),
				File: git.File{
					Path:     "test.md",
					FullPath: "test.md:3",
					Line:     3,
					FileName: "test.md",
				},
				Thread: git.Thread{
					ID:         "gid://gitlab/DiffDiscussion/dsoijierew22ukc",
					IsResolved: false,
				},
			},
			{
				Author: git.Author{
					Login: "Mario",
				},
				Body:      "thanks I think so too",
				CreatedAt: internal.TimeMustParse(time.RFC3339, "2025-06-26T14:19:27Z"),
				File: git.File{
					Path:     "test.md",
					FullPath: "test.md:3",
					Line:     3,
					FileName: "test.md",
				},
				Thread: git.Thread{
					ID:         "gid://gitlab/DiffDiscussion/dsoijierew22ukc",
					IsResolved: false,
				},
			},
		},
		State: git.State{},
		Title: "Test pr",
		Id:    "gid://gitlab/MergeRequest/11111111",
	}

	suite.mockGraphQL.EXPECT().
		Do(graphql.GitlabPRQuery(suite.repo), gomock.Any(), gomock.Any()).
		DoAndReturn(func(query gitlab.GraphQLQuery, gr any, _ ...any) (gitlab.Response, error) {

			err := json.Unmarshal([]byte(prDetails), &gr)
			suite.NoError(err)
			return gitlab.Response{}, nil
		})
	details, err := suite.prService.GetPRDetails(suite.repo, false)
	assert.NoError(suite.T(), err)
	assert.Len(suite.T(), details.Comments, len(details.Comments))
	for index, comment := range expected.Comments {
		assert.Equal(suite.T(), comment, details.Comments[index])
	}
}

func (suite *GitlabClientTestSuite) TestPRService_getPrDetails_with_verbose() {
	prDetails := `{
  "data": {
    "project": {
      "mergeRequests": {
        "nodes": []
      },
      "mergeRequest": {
        "id": "gid://gitlab/MergeRequest/11111111",
        "state": "opened",
        "mergeable": true,
        "mergeableDiscussionsState": true,
        "mergeStatusEnum": "CAN_BE_MERGED",
        "detailedMergeStatus": "NOT_APPROVED",
        "conflicts": false,
        "approvalsLeft": 0,
        "approvalsRequired": 0,
        "approved": false,
        "pipelines": {
          "nodes": [
            {
              "complete": true,
              "status": "FAILED",
              "name": null,
              "active": false,
              "failureReason": null,
              "jobs": {
					"nodes": [
						{
							"name": "run_tests",
							"status": "FAILED"
						},
						{
							"name": "lint",
							"status": "SUCCEEDED"
						}
					]
				}
			}
			]
        },
		"approvedBy": {
          "nodes": [
		  	{
		  		"name": "Luigi"
		  	}
		  ]
        },
		"reviewers": {
          "nodes": [
            {
              "name": "GitLab Duo"
            },
            {
              "name": "Luigi"
            }
          ]
        },
        "title": "Test pr",
        "createdAt": "2025-06-26T14:16:54Z",
        "description": "",
        "author": {
          "name": "Mario"
        },
        "notes": null
      }
    }
  },
  "correlationId": "2423kn3kjn4kj32n42kjn3k4n"
}`

	expected := &git.PR{
		Comments: nil,
		State: git.State{
			Statuses: []git.Status{
				{
					Name:       "run_tests",
					Conclusion: "FAILED",
				}, {
					Name:       "lint",
					Conclusion: "SUCCEEDED",
				},
			},
			Reviews: map[string][]string{
				"APPROVED": {"Luigi"},
				"PENDING":  {"GitLab Duo"},
			},
			MergeStatus: "Merge request must be approved before merging.",
		},
		Title: "Test pr",
		Id:    "gid://gitlab/MergeRequest/11111111",
	}

	suite.mockGraphQL.EXPECT().
		Do(graphql.GitlabPRQuery(suite.repo), gomock.Any(), gomock.Any()).
		DoAndReturn(func(query gitlab.GraphQLQuery, gr any, _ ...any) (gitlab.Response, error) {

			err := json.Unmarshal([]byte(prDetails), &gr)
			suite.NoError(err)
			return gitlab.Response{}, nil
		})
	details, err := suite.prService.GetPRDetails(suite.repo, true)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), expected, details)
}

func (suite *GitlabClientTestSuite) TestReply_with_thread() {
	suite.mockGraphQL.EXPECT().
		Do(graphql.AddNote("mergeID", "my comment", "discussionID"), gomock.Any(), gomock.Any()).
		DoAndReturn(func(query gitlab.GraphQLQuery, gr any, _ ...any) (gitlab.Response, error) {
			response := `{
        "Errors": []
      }`
			err := json.Unmarshal([]byte(response), &gr)
			suite.NoError(err)
			return gitlab.Response{}, nil
		})

	commentToReplyTo := &git.Comment{
		Thread: git.Thread{
			ID: "discussionID",
		},
	}
	err := suite.prService.Reply("my comment", commentToReplyTo, "mergeID")
	assert.NoError(suite.T(), err)
}

func (suite *GitlabClientTestSuite) TestReply_without_thread() {
	suite.mockGraphQL.EXPECT().
		Do(graphql.AddNote("mergeID", "my comment", ""), gomock.Any(), gomock.Any()).
		DoAndReturn(func(query gitlab.GraphQLQuery, gr any, _ ...any) (gitlab.Response, error) {
			response := `{
        "Errors": []
      }`
			err := json.Unmarshal([]byte(response), &gr)
			suite.NoError(err)
			return gitlab.Response{}, nil
		})

	commentToReplyTo := &git.Comment{}
	err := suite.prService.Reply("my comment", commentToReplyTo, "mergeID")
	assert.NoError(suite.T(), err)
}

func (suite *GitlabClientTestSuite) TestReply_errors() {
	suite.mockGraphQL.EXPECT().
		Do(graphql.AddNote("mergeID", "my comment", ""), gomock.Any(), gomock.Any()).
		DoAndReturn(func(query gitlab.GraphQLQuery, gr any, _ ...any) (gitlab.Response, error) {
			response := `{
        "Errors": [
          {
            "Message": "one of the errors"
          },
          {
            "Message": "second error"
          }
        ]
      }`
			err := json.Unmarshal([]byte(response), &gr)
			suite.NoError(err)
			return gitlab.Response{}, nil
		})

	commentToReplyTo := &git.Comment{}
	err := suite.prService.Reply("my comment", commentToReplyTo, "mergeID")
	assert.ErrorContains(suite.T(), err, "one of the errors")
	assert.ErrorContains(suite.T(), err, "second error")
}

func (suite *GitlabClientTestSuite) TestResolve() {
	suite.mockGraphQL.EXPECT().
		Do(graphql.ResolveDiscussion("threadID"), gomock.Any(), gomock.Any()).
		DoAndReturn(func(query gitlab.GraphQLQuery, gr any, _ ...any) (gitlab.Response, error) {
			response := `{
        "Errors": []
      }`
			err := json.Unmarshal([]byte(response), &gr)
			suite.NoError(err)
			return gitlab.Response{}, nil
		})

	commentToReplyTo := &git.Comment{Thread: git.Thread{ID: "threadID"}}
	err := suite.prService.Resolve(commentToReplyTo)
	assert.NoError(suite.T(), err)
}

func (suite *GitlabClientTestSuite) TestResolve_errors() {
	suite.mockGraphQL.EXPECT().
		Do(graphql.ResolveDiscussion("threadID"), gomock.Any(), gomock.Any()).
		DoAndReturn(func(query gitlab.GraphQLQuery, gr any, _ ...any) (gitlab.Response, error) {
			response := `{
        "Errors": [
          {
            "Message": "one of the errors"
          },
          {
            "Message": "second error"
          }
        ]
      }`
			err := json.Unmarshal([]byte(response), &gr)
			suite.NoError(err)
			return gitlab.Response{}, nil
		})

	commentToReplyTo := &git.Comment{Thread: git.Thread{ID: "threadID"}}
	err := suite.prService.Resolve(commentToReplyTo)
	assert.ErrorContains(suite.T(), err, "one of the errors")
	assert.ErrorContains(suite.T(), err, "second error")
}

func TestGitlabClientTestSuite(t *testing.T) {
	suite.Run(t, new(GitlabClientTestSuite))
}
