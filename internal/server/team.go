package server

import (
	"context"

	"github.com/iskanye/avito-tech-internship/pkg/api"
)

// (POST /team/add)
func (serverAPI) PostTeamAdd(
	c context.Context,
	req api.PostTeamAddRequestObject,
) (api.PostTeamAddResponseObject, error) {
	return nil, nil
}

// (GET /team/get)
func (serverAPI) GetTeamGet(
	c context.Context,
	req api.GetTeamGetRequestObject,
) (api.GetTeamGetResponseObject, error) {
	return nil, nil
}
