package sdk

import (
	"fmt"
	"github.com/tribalwarshelp/shared/tw/twmodel"

	"github.com/Kichiyaki/gqlgen-client/client"
	"github.com/pkg/errors"
)

type Village struct {
	sdk *SDK
}

type VillageInclude struct {
	Player        bool
	PlayerInclude PlayerInclude
}

func (incl VillageInclude) String() string {
	i := ""
	if incl.Player {
		i += fmt.Sprintf(`
			player {
				%s
				%s
			}
		`, playerFields, incl.PlayerInclude.String())
	}
	return i
}

func (v *Village) Read(server string, id int, include *VillageInclude) (*twmodel.Village, error) {
	if server == "" {
		return nil, ErrServerNameIsEmpty
	}
	if include == nil {
		include = &VillageInclude{}
	}
	resp := struct {
		Village twmodel.Village `json:"village" gqlgen:"village"`
	}{}

	query := fmt.Sprintf(`
		query village($server: String!, $id: Int!) {
			village(server: $server, id: $id) {
				%s
				%s
			}
		}
	`, villageFields, include.String())
	err := v.sdk.Post(query, &resp, client.Var("server", server), client.Var("id", id))
	if err != nil {
		return nil, errors.Wrap(err, "twhelp sdk")
	}
	return &resp.Village, nil
}

type VillageList struct {
	Items []*twmodel.Village `json:"items" gqlgen:"items"`
	Total int                `json:"total" gqlgen:"total"`
}

func (v *Village) Browse(server string,
	limit,
	offset int,
	sort []string,
	filter *twmodel.VillageFilter,
	include *VillageInclude) (*VillageList, error) {
	if server == "" {
		return nil, ErrServerNameIsEmpty
	}
	if filter == nil {
		filter = &twmodel.VillageFilter{}
	}
	if include == nil {
		include = &VillageInclude{}
	}
	resp := struct {
		Villages VillageList
	}{}
	query := fmt.Sprintf(`
		query villages($server: String!, $filter: VillageFilter, $limit: Int, $offset: Int, $sort: [String!]) {
			villages(server: $server, filter: $filter, limit: $limit, offset: $offset, sort: $sort) {
				items {
					%s
					%s
				}
				total
			}
		}
	`, villageFields, include.String())

	err := v.sdk.Post(query,
		&resp,
		client.Var("filter", filter),
		client.Var("server", server),
		client.Var("limit", limit),
		client.Var("offset", offset),
		client.Var("sort", sort))
	if err != nil {
		return nil, errors.Wrap(err, "twhelp sdk")
	}
	return &resp.Villages, nil
}
