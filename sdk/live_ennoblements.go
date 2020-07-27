package sdk

import (
	"fmt"

	"github.com/Kichiyaki/gqlgen-client/client"
	"github.com/pkg/errors"
	"github.com/tribalwarshelp/shared/models"
)

type LiveEnnoblements struct {
	sdk *SDK
}

type LiveEnnoblementInclude struct {
	NewOwner        bool
	NewOwnerInclude PlayerInclude
	OldOwner        bool
	OldOwnerInclude PlayerInclude
	Village         bool
}

func (incl LiveEnnoblementInclude) String() string {
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

type liveEnnoblementsResponse struct {
	LiveEnnoblements []*models.LiveEnnoblement `json:"liveEnnoblements" gqlgen:"liveEnnoblements"`
}

func (en *LiveEnnoblements) Browse(server string, include *LiveEnnoblementInclude) ([]*models.LiveEnnoblement, error) {
	if server == "" {
		return nil, ErrServerNameIsEmpty
	}
	if include == nil {
		include = &LiveEnnoblementInclude{}
	}
	resp := &liveEnnoblementsResponse{}
	query := fmt.Sprintf(`
		query liveEnnoblements($server: String!) {
			liveEnnoblements(server: $server) {
				ennobledAt
				%s
			}
		}
	`, include.String())
	err := en.sdk.Post(minifyString(query), &resp, client.Var("server", server))
	if err != nil {
		return nil, errors.Wrap(err, "twhelp sdk")
	}
	return resp.LiveEnnoblements, nil
}
