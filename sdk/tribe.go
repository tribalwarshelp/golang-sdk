package sdk

import (
	"fmt"
	"github.com/tribalwarshelp/shared/tw/twmodel"

	"github.com/Kichiyaki/gqlgen-client/client"
	"github.com/pkg/errors"
)

type Tribe struct {
	sdk *SDK
}

func (t *Tribe) Read(server string, id int) (*twmodel.Tribe, error) {
	resp := struct {
		Tribe twmodel.Tribe `json:"tribe" gqlgen:"tribe"`
	}{}
	query := fmt.Sprintf(`
		query tribe($server: String!, $id: Int!) {
			tribe(server: $server, id: $id) {
				%s
			}
		}
	`, tribeFields)
	err := t.sdk.Post(query, &resp, client.Var("server", server), client.Var("id", id))
	if err != nil {
		return nil, errors.Wrap(err, "twhelp sdk")
	}
	return &resp.Tribe, nil
}

type TribeList struct {
	Items []*twmodel.Tribe `json:"items" gqlgen:"items"`
	Total int              `json:"total" gqlgen:"total"`
}

func (t *Tribe) Browse(server string,
	limit,
	offset int,
	sort []string,
	filter *twmodel.TribeFilter) (*TribeList, error) {
	if filter == nil {
		filter = &twmodel.TribeFilter{}
	}
	resp := struct {
		Tribes TribeList `json:"tribes" gqlgen:"tribes"`
	}{}
	query := fmt.Sprintf(`
		query tribes($server: String!, $filter: TribeFilter, $limit: Int, $offset: Int, $sort: [String!]) {
			tribes(server: $server, filter: $filter, limit: $limit, offset: $offset, sort: $sort) {
				items {
					%s
				}
				total
			}
		}
	`, tribeFields)

	err := t.sdk.Post(query,
		&resp,
		client.Var("filter", filter),
		client.Var("server", server),
		client.Var("limit", limit),
		client.Var("offset", offset),
		client.Var("sort", sort))
	if err != nil {
		return nil, errors.Wrap(err, "twhelp sdk")
	}
	return &resp.Tribes, nil
}
