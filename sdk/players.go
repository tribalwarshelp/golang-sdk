package sdk

import (
	"fmt"

	"github.com/Kichiyaki/gqlgen-client/client"
	"github.com/pkg/errors"
	"github.com/tribalwarshelp/shared/models"
)

type Players struct {
	sdk *SDK
}

type PlayerInclude struct {
	Tribe bool
}

func (incl PlayerInclude) String() string {
	i := ""
	if incl.Tribe {
		i += fmt.Sprintf(`
			tribe {
				%s
			}
		`, tribeFields)
	}
	return i
}

func (ps *Players) Read(server string, id int, include *PlayerInclude) (*models.Player, error) {
	if server == "" {
		return nil, ErrServerNameIsEmpty
	}
	if include == nil {
		include = &PlayerInclude{}
	}
	resp := struct {
		Player models.Player `json:"player" gqlgen:"player"`
	}{}

	query := fmt.Sprintf(`
		query player($server: String!, $id: Int!) {
			player(server: $server, id: $id) {
				%s
				%s
			}
		}
	`, playerFields, include.String())
	err := ps.sdk.Post(minifyString(query), &resp, client.Var("server", server), client.Var("id", id))
	if err != nil {
		return nil, errors.Wrap(err, "twhelp sdk")
	}
	return &resp.Player, nil
}

type PlayerList struct {
	Items []*models.Player `json:"items" gqlgen:"items"`
	Total int              `json:"total" gqlgen:"total"`
}

func (ps *Players) Browse(server string, filter *models.PlayerFilter, include *PlayerInclude) (*PlayerList, error) {
	if server == "" {
		return nil, ErrServerNameIsEmpty
	}
	if filter == nil {
		filter = &models.PlayerFilter{}
	}
	if include == nil {
		include = &PlayerInclude{}
	}
	resp := struct {
		Players PlayerList `json:"players" gqlgen:"players"`
	}{}
	query := fmt.Sprintf(`
		query players($server: String!, $filter: PlayerFilter) {
			players(server: $server, filter: $filter) {
				items {
					%s
					%s
				}
				total
			}
		}
	`, playerFields, include.String())

	err := ps.sdk.Post(minifyString(query), &resp, client.Var("filter", filter), client.Var("server", server))
	if err != nil {
		return nil, errors.Wrap(err, "twhelp sdk")
	}
	return &resp.Players, nil
}
