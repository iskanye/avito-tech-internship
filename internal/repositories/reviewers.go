package repositories

import "context"

func (s *Storage) AssignReviewers(
	ctx context.Context,
	pullRequestID string,
	authorID string,
) error {
	const op = "repositories.postgres.AssignReviewers"
	return nil
}
