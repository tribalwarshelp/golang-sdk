package sdk

import (
	"fmt"

	"github.com/Kichiyaki/gqlgen-client/client"
	"github.com/pkg/errors"
	"github.com/tribalwarshelp/shared/models"
)

type Villages struct {
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

func (vs *Villages) Read(server string, id int, include *VillageInclude) (*models.Village, error) {
	if server == "" {
		return nil, ErrServerNameIsEmpty
	}
	if include == nil {
		include = &VillageInclude{}
	}
	resp := struct {
		Village models.Village `json:"village" gqlgen:"village"`
	}{}

	query := fmt.Sprintf(`
		query village($server: String!, $id: Int!) {
			village(server: $server, id: $id) {
				%s
				%s
			}
		}
	`, villageFields, include.String())
	err := vs.sdk.Post(minifyString(query), &resp, client.Var("server", server), client.Var("id", id))
	if err != nil {
		return nil, errors.Wrap(err, "twhelp sdk")
	}
	return &resp.Village, nil
}

type VillageList struct {
	Items []*models.Village `json:"items" gqlgen:"items"`
	Total int               `json:"total" gqlgen:"total"`
}

func (vs *Villages) Browse(server string, filter *models.VillageFilter, include *VillageInclude) (*VillageList, error) {
	if server == "" {
		return nil, ErrServerNameIsEmpty
	}
	if filter == nil {
		filter = &models.VillageFilter{}
	}
	if include == nil {
		include = &VillageInclude{}
	}
	resp := struct {
		Villages VillageList
	}{}
	query := fmt.Sprintf(`
		query villages($server: String!, $filter: VillageFilter) {
			villages(server: $server, filter: $filter) {
				items {
					%s
					%s
				}
				total
			}
		}
	`, villageFields, include.String())

	err := vs.sdk.Post(minifyString(query), &resp, client.Var("filter", filter), client.Var("server", server))
	if err != nil {
		return nil, errors.Wrap(err, "twhelp sdk")
	}
	return &resp.Villages, nil
}
