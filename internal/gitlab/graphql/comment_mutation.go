package graphql

import (
	"fmt"

	gitlab "gitlab.com/gitlab-org/api/client-go"
)

func AddNote(mergeRequestId string, body string, discussionId string) gitlab.GraphQLQuery {
	input := fmt.Sprintf(`{
			noteableId: "%s"
			body: "%s"
		}
	`, mergeRequestId, body)

	if discussionId != "" {
		input = fmt.Sprintf(`{
			noteableId: "%s"
			discussionId: "%s"
			body: "%s"
		}
	`, mergeRequestId, discussionId, body)
	}
	query := fmt.Sprintf(`
		mutation Comment() {
		createNote(input: %s) {
			errors
		}
	}
`, input)

	return gitlab.GraphQLQuery{
		Query: query,
	}
}

func ResolveDiscussion(id string) gitlab.GraphQLQuery {
	query := fmt.Sprintf(`
	mutation {
		discussionToggleResolve(input: {
			resolve: true,
			id: "%s"
		}) {
			clientMutationId
		}
	}
	`, id)
	return gitlab.GraphQLQuery{
		Query: query,
	}
}
