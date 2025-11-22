package suite

import (
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/iskanye/avito-tech-internship/pkg/api"
	"github.com/stretchr/testify/assert"
)

const (
	timeDelta = 1
)

func RandomPullRequest(authorID string) *api.PostPullRequestCreateJSONRequestBody {
	return &api.PostPullRequestCreateJSONRequestBody{
		PullRequestId:   gofakeit.UUID(),
		PullRequestName: gofakeit.Sentence(3),
		AuthorId:        authorID,
	}
}

func AssertPullRequestEqual(
	t *testing.T,
	pr1 *api.PostPullRequestCreateJSONRequestBody,
	pr2 *api.PullRequest,
) {
	assert.Equal(t, pr1.PullRequestId, pr2.PullRequestId)
	assert.Equal(t, pr1.PullRequestName, pr2.PullRequestName)
	assert.Equal(t, pr1.AuthorId, pr2.AuthorId)
	assert.Equal(t, api.PullRequestStatusOPEN, pr2.Status)

	assert.InDelta(t, time.Now().Unix(), pr2.CreatedAt.Unix(), timeDelta)
}
