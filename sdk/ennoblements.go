package sdk

import (
	"fmt"

	"github.com/Kichiyaki/gqlgen-client/client"
	"github.com/pkg/errors"
	"github.com/tribalwarshelp/shared/models"
)

type Ennoblements struct {
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

type ennoblementsResponse struct {
	Ennoblements []*models.Ennoblement `json:"ennoblements" gqlgen:"ennoblements"`
}

func (en *Ennoblements) Browse(server string, include *EnnoblementInclude) ([]*models.Ennoblement, error) {
	if server == "" {
		return nil, ErrServerNameIsEmpty
	}
	if include == nil {
		include = &EnnoblementInclude{}
	}
	resp := &ennoblementsResponse{}
	query := fmt.Sprintf(`
		query ennoblements($server: String!) {
			ennoblements(server: $server) {
				ennobledAt
				%s
			}
		}
	`, include.String())
	err := en.sdk.client.Post(minifyString(query), &resp, client.Var("server", server))
	if err != nil {
		return nil, errors.Wrap(err, "twhelp sdk")
	}
	return resp.Ennoblements, nil
}
