package sdk

import (
	"github.com/Kichiyaki/gqlgen-client/client"
	"github.com/pkg/errors"
	"github.com/tribalwarshelp/shared/models"
)

type Servers struct {
	sdk *SDK
}

func (ss *Servers) Read(key string) (*models.Server, error) {
	resp := struct {
		Server *models.Server `json:"server" gqlgen:"server"`
	}{}
	query := `
		query server($key: String!) {
			server(key: $key) {
				id
				key
				status
				langVersionTag
			}
		}
	`
	err := ss.sdk.client.Post(minifyString(query), &resp, client.Var("key", key))
	if err != nil {
		return nil, errors.Wrap(err, "twhelp sdk")
	}
	return resp.Server, nil
}

type ServersList struct {
	Items []*models.Server `json:"items" gqlgen:"items"`
	Total int              `json:"total" gqlgen:"total"`
}

func (ss *Servers) Browse(filter *models.ServerFilter) (*ServersList, error) {
	if filter == nil {
		filter = &models.ServerFilter{}
	}
	resp := struct {
		Servers *ServersList `json:"servers" gqlgen:"servers"`
	}{}

	query := `
		query servers($filter: ServerFilter) {
			servers(filter: $filter) {
				items {
					id
					key
					status
					langVersionTag
				}
				total
			}
		}
	`

	err := ss.sdk.client.Post(minifyString(query), &resp, client.Var("filter", filter))
	if err != nil {
		return nil, errors.Wrap(err, "twhelp sdk")
	}
	return resp.Servers, nil
}
