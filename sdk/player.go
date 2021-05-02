package sdk

import (
	"fmt"
	"github.com/tribalwarshelp/shared/tw/twmodel"

	"github.com/Kichiyaki/gqlgen-client/client"
	"github.com/pkg/errors"
)

type Player struct {
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

func (p *Player) Read(server string, id int, include *PlayerInclude) (*twmodel.Player, error) {
	if server == "" {
		return nil, ErrServerNameIsEmpty
	}
	if include == nil {
		include = &PlayerInclude{}
	}
	resp := struct {
		Player twmodel.Player `json:"player" gqlgen:"player"`
	}{}

	query := fmt.Sprintf(`
		query player($server: String!, $id: Int!) {
			player(server: $server, id: $id) {
				%s
				%s
			}
		}
	`, playerFields, include.String())
	err := p.sdk.Post(query, &resp, client.Var("server", server), client.Var("id", id))
	if err != nil {
		return nil, errors.Wrap(err, "twhelp sdk")
	}
	return &resp.Player, nil
}

type PlayerList struct {
	Items []*twmodel.Player `json:"items" gqlgen:"items"`
	Total int               `json:"total" gqlgen:"total"`
}

func (p *Player) Browse(server string,
	limit,
	offset int,
	sort []string,
	filter *twmodel.PlayerFilter,
	include *PlayerInclude) (*PlayerList, error) {
	if server == "" {
		return nil, ErrServerNameIsEmpty
	}
	if filter == nil {
		filter = &twmodel.PlayerFilter{}
	}
	if include == nil {
		include = &PlayerInclude{}
	}
	resp := struct {
		Players PlayerList `json:"players" gqlgen:"players"`
	}{}
	query := fmt.Sprintf(`
		query players($server: String!, $filter: PlayerFilter, $limit: Int, $offset: Int, $sort: [String!]) {
			players(server: $server, filter: $filter, limit: $limit, offset: $offset, sort: $sort) {
				items {
					%s
					%s
				}
				total
			}
		}
	`, playerFields, include.String())

	err := p.sdk.Post(query,
		&resp,
		client.Var("filter", filter),
		client.Var("server", server),
		client.Var("limit", limit),
		client.Var("offset", offset),
		client.Var("sort", sort))
	if err != nil {
		return nil, errors.Wrap(err, "twhelp sdk")
	}
	return &resp.Players, nil
}
