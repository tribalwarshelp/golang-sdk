package sdk

import (
	"fmt"

	"github.com/Kichiyaki/gqlgen-client/client"
	"github.com/pkg/errors"
	"github.com/tribalwarshelp/shared/models"
)

type Tribes struct {
	sdk *SDK
}

func (ts *Tribes) Read(server string, id int) (*models.Tribe, error) {
	resp := struct {
		Tribe models.Tribe `json:"tribe" gqlgen:"tribe"`
	}{}
	query := fmt.Sprintf(`
		query tribe($server: String!, $id: Int!) {
			tribe(server: $server, id: $id) {
				%s
			}
		}
	`, tribeFields)
	err := ts.sdk.Post(query, &resp, client.Var("server", server), client.Var("id", id))
	if err != nil {
		return nil, errors.Wrap(err, "twhelp sdk")
	}
	return &resp.Tribe, nil
}

type TribeList struct {
	Items []*models.Tribe `json:"items" gqlgen:"items"`
	Total int             `json:"total" gqlgen:"total"`
}

func (ts *Tribes) Browse(server string, filter *models.TribeFilter) (*TribeList, error) {
	if filter == nil {
		filter = &models.TribeFilter{}
	}
	resp := struct {
		Tribes TribeList `json:"tribes" gqlgen:"tribes"`
	}{}
	query := fmt.Sprintf(`
		query tribes($server: String!, $filter: TribeFilter) {
			tribes(server: $server, filter: $filter) {
				items {
					%s
				}
				total
			}
		}
	`, tribeFields)

	err := ts.sdk.Post(query, &resp, client.Var("server", server), client.Var("filter", filter))
	if err != nil {
		return nil, errors.Wrap(err, "twhelp sdk")
	}
	return &resp.Tribes, nil
}
