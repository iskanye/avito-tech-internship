package suite

import (
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/iskanye/avito-tech-internship/internal/models"
	"github.com/iskanye/avito-tech-internship/pkg/api"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	timeDelta = 1
)

func RandomPullRequest(authorID string) *api.PullRequest {
	now := time.Now().Truncate(1)
	return &api.PullRequest{
		PullRequestId:   gofakeit.UUID(),
		PullRequestName: gofakeit.Sentence(3),
		AuthorId:        authorID,
		Status:          api.PullRequestStatusOPEN,
		CreatedAt:       &now,
	}
}

func CheckPullRequestEqual(
	t *testing.T,
	pr1 *api.PullRequest,
	pr2 *api.PullRequest,
) {
	assert.Equal(t, pr1.PullRequestId, pr2.PullRequestId)
	assert.Equal(t, pr1.PullRequestName, pr2.PullRequestName)
	assert.Equal(t, pr1.AuthorId, pr2.AuthorId)
	assert.Equal(t, pr1.Status, pr2.Status)

	assert.InDelta(t, pr1.CreatedAt.Unix(), pr2.CreatedAt.Unix(), timeDelta)

	if pr1.MergedAt != nil && pr2.MergedAt != nil {
		assert.InDelta(t, pr1.MergedAt.Unix(), pr2.MergedAt.Unix(), timeDelta)
	} else {
		assert.Equal(t, pr1.MergedAt, pr2.MergedAt)
	}
}

func CheckPullRequestsEqual(
	t *testing.T,
	pullRequests1 []api.PullRequestShort,
	pullRequests2 ...models.PullRequest,
) {
	prs := make(map[string]api.PullRequestShort)
	for _, pr := range pullRequests1 {
		prs[pr.PullRequestId] = pr
	}

	for _, pr1 := range pullRequests2 {
		pr2, ok := prs[pr1.ID]

		require.True(t, ok)
		assert.Equal(t, pr1.Name, pr2.PullRequestName)
		assert.Equal(t, pr1.AuthorID, pr2.AuthorId)
		assert.Equal(t, string(pr1.Status), string(pr2.Status))
	}
}

func PullRequestCreateToModel(
	pr *api.PullRequest,
) models.PullRequest {
	return models.PullRequest{
		ID:       pr.PullRequestId,
		Name:     pr.PullRequestName,
		AuthorID: pr.AuthorId,
		Status:   models.PULLREQUEST_OPEN,
	}
}
