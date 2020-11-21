package sdk

import (
	"fmt"

	"github.com/Kichiyaki/gqlgen-client/client"
	"github.com/pkg/errors"
	"github.com/tribalwarshelp/shared/models"
)

type Server struct {
	sdk *SDK
}

type ServerInclude struct {
	Version bool
}

func (incl ServerInclude) String() string {
	i := ""
	if incl.Version {
		i += fmt.Sprintf(`
			version {
				%s
			}
		`, versionFields)
	}
	return i
}

func (s *Server) Read(key string, incl *ServerInclude) (*models.Server, error) {
	if incl == nil {
		incl = &ServerInclude{}
	}
	resp := struct {
		Server models.Server `json:"server" gqlgen:"server"`
	}{}
	query := fmt.Sprintf(`
		query server($key: String!) {
			server(key: $key) {
				key
				status
				dataUpdatedAt
				historyUpdatedAt
				statsUpdatedAt
				numberOfTribes
				numberOfPlayers
				numberOfVillages
				%s
			}
		}
	`, incl.String())
	err := s.sdk.Post(query, &resp, client.Var("key", key))
	if err != nil {
		return nil, errors.Wrap(err, "twhelp sdk")
	}
	return &resp.Server, nil
}

type ServerList struct {
	Items []*models.Server `json:"items" gqlgen:"items"`
	Total int              `json:"total" gqlgen:"total"`
}

func (s *Server) Browse(limit,
	offset int,
	sort []string,
	filter *models.ServerFilter,
	incl *ServerInclude) (*ServerList, error) {
	if incl == nil {
		incl = &ServerInclude{}
	}
	if filter == nil {
		filter = &models.ServerFilter{}
	}
	resp := struct {
		Servers ServerList `json:"servers" gqlgen:"servers"`
	}{}

	query := fmt.Sprintf(`
		query servers($filter: ServerFilter, $limit: Int, $offset: Int, $sort: [String!]) {
			servers(filter: $filter, limit: $limit, offset: $offset, sort: $sort) {
					items {
					key
					status
					dataUpdatedAt
					historyUpdatedAt
					statsUpdatedAt
					numberOfTribes
					numberOfPlayers
					numberOfVillages
					%s
				}
				total
			}
		}
	`, incl.String())

	err := s.sdk.Post(query,
		&resp,
		client.Var("filter", filter),
		client.Var("limit", limit),
		client.Var("offset", offset),
		client.Var("sort", sort))
	if err != nil {
		return nil, errors.Wrap(err, "twhelp sdk")
	}
	return &resp.Servers, nil
}
