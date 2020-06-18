package sdk

import (
	"fmt"

	"github.com/Kichiyaki/gqlgen-client/client"
	"github.com/pkg/errors"
	"github.com/tribalwarshelp/shared/models"
)

type Servers struct {
	sdk *SDK
}

type ServerInclude struct {
	LangVersion bool
}

func (incl ServerInclude) String() string {
	i := ""
	if incl.LangVersion {
		i += fmt.Sprintf(`
			langVersion {
				%s
			}
		`, langVersionFields)
	}
	return i
}

func (ss *Servers) Read(key string, incl *ServerInclude) (*models.Server, error) {
	if incl == nil {
		incl = &ServerInclude{}
	}
	resp := struct {
		Server *models.Server `json:"server" gqlgen:"server"`
	}{}
	query := fmt.Sprintf(`
		query server($key: String!) {
			server(key: $key) {
				id
				key
				status
				%s
			}
		}
	`, incl.String())
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

func (ss *Servers) Browse(filter *models.ServerFilter, incl *ServerInclude) (*ServersList, error) {
	if incl == nil {
		incl = &ServerInclude{}
	}
	if filter == nil {
		filter = &models.ServerFilter{}
	}
	resp := struct {
		Servers *ServersList `json:"servers" gqlgen:"servers"`
	}{}

	query := fmt.Sprintf(`
		query servers($filter: ServerFilter) {
			servers(filter: $filter) {
				items {
					id
					key
					status
					%s
				}
				total
			}
		}
	`, incl.String())

	err := ss.sdk.client.Post(minifyString(query), &resp, client.Var("filter", filter))
	if err != nil {
		return nil, errors.Wrap(err, "twhelp sdk")
	}
	return resp.Servers, nil
}
