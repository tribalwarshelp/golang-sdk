package sdk

import (
	"fmt"

	"github.com/Kichiyaki/gqlgen-client/client"
	"github.com/pkg/errors"
	"github.com/tribalwarshelp/shared/models"
)

type Ennoblement struct {
	sdk *SDK
}

type EnnoblementInclude struct {
	NewOwner        bool
	NewOwnerInclude PlayerInclude
	OldOwner        bool
	OldOwnerInclude PlayerInclude
	Village         bool
}

func (incl EnnoblementInclude) String() string {
	i := ""
	if incl.NewOwner {
		i += fmt.Sprintf(`
			newOwner {
				%s
				%s
			}
		`, playerFields, incl.NewOwnerInclude.String())
	}
	if incl.OldOwner {
		i += fmt.Sprintf(`
			oldOwner {
				%s
				%s
			}
		`, playerFields, incl.OldOwnerInclude.String())
	}
	if incl.Village {
		i += fmt.Sprintf(`
			village {
				%s
			}
		`, villageFields)
	}
	return i
}

type EnnoblementList struct {
	Items []*models.Ennoblement `json:"items" gqlgen:"items"`
	Total int                   `json:"total" gqlgen:"total"`
}

func (en *Ennoblement) Browse(server string,
	limit,
	offset int,
	sort []string,
	filter *models.EnnoblementFilter,
	include *EnnoblementInclude) (*EnnoblementList, error) {
	if server == "" {
		return nil, ErrServerNameIsEmpty
	}
	if include == nil {
		include = &EnnoblementInclude{}
	}
	resp := struct {
		Ennoblements EnnoblementList `json:"ennoblements" gqlgen:"ennoblements"`
	}{}
	query := fmt.Sprintf(`
		query ennoblements($server: String!, $filter: EnnoblementFilter, $limit: Int, $offset: Int, $sort: [String!]) {
			ennoblements(server: $server, filter: $filter, limit: $limit, offset: $offset, sort: $sort) {
				items {
					ennobledAt
					%s
				}
				total
			}
		}
	`, include.String())
	err := en.sdk.Post(query,
		&resp,
		client.Var("filter", filter),
		client.Var("server", server),
		client.Var("limit", limit),
		client.Var("offset", offset),
		client.Var("sort", sort))
	if err != nil {
		return nil, errors.Wrap(err, "twhelp sdk")
	}
	return &resp.Ennoblements, nil
}
