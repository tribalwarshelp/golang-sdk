package sdk

import (
	"fmt"

	"github.com/Kichiyaki/gqlgen-client/client"
	"github.com/pkg/errors"
	"github.com/tribalwarshelp/shared/models"
)

type Version struct {
	sdk *SDK
}

func (lv *Version) Read(code models.VersionCode) (*models.Version, error) {
	resp := struct {
		Version models.Version `json:"version" gqlgen:"version"`
	}{}
	query := fmt.Sprintf(`
		query version($code: VersionCode!) {
			version(code: $code) {
				%s
			}
		}
	`, versionFields)
	err := lv.sdk.Post(query, &resp, client.Var("code", code))
	if err != nil {
		return nil, errors.Wrap(err, "twhelp sdk")
	}
	return &resp.Version, nil
}

type VersionList struct {
	Items []*models.Version `json:"items" gqlgen:"items"`
	Total int               `json:"total" gqlgen:"total"`
}

func (lv *Version) Browse(filter *models.VersionFilter) (*VersionList, error) {
	if filter == nil {
		filter = &models.VersionFilter{}
	}
	resp := struct {
		Versions VersionList `json:"versions" gqlgen:"versions"`
	}{}
	query := fmt.Sprintf(`
		query versions($filter: VersionFilter) {
			versions(filter: $filter) {
				items {
					%s
				}
				total
			}
		}
	`, versionFields)

	err := lv.sdk.Post(query, &resp, client.Var("filter", filter))
	if err != nil {
		return nil, errors.Wrap(err, "twhelp sdk")
	}
	return &resp.Versions, nil
}
