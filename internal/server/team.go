package server

import (
	"context"

	"github.com/iskanye/avito-tech-internship/pkg/api"
)

// (POST /team/add)
func (Server) PostTeamAdd(
	c context.Context,
	req api.PostTeamAddRequestObject,
) (api.PostTeamAddResponseObject, error) {
	return nil, nil
}

// (GET /team/get)
func (Server) GetTeamGet(
	c context.Context,
	req api.GetTeamGetRequestObject,
) (api.GetTeamGetResponseObject, error) {
	return nil, nil
}
