package api

import "context"

// (POST /team/add)
func (Server) PostTeamAdd(
	c context.Context,
	req PostTeamAddRequestObject,
) (PostTeamAddResponseObject, error) {
	return nil, nil
}

// (GET /team/get)
func (Server) GetTeamGet(
	c context.Context,
	req GetTeamGetRequestObject,
) (GetTeamGetResponseObject, error) {
	return nil, nil
}
