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

func (v *Version) Read(code models.VersionCode) (*models.Version, error) {
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
	err := v.sdk.Post(query, &resp, client.Var("code", code))
	if err != nil {
		return nil, errors.Wrap(err, "twhelp sdk")
	}
	return &resp.Version, nil
}

type VersionList struct {
	Items []*models.Version `json:"items" gqlgen:"items"`
	Total int               `json:"total" gqlgen:"total"`
}

func (v *Version) Browse(limit,
	offset int,
	sort []string,
	filter *models.VersionFilter) (*VersionList, error) {
	if filter == nil {
		filter = &models.VersionFilter{}
	}
	resp := struct {
		Versions VersionList `json:"versions" gqlgen:"versions"`
	}{}
	query := fmt.Sprintf(`
		query versions($filter: VersionFilter, $limit: Int, $offset: Int, $sort: [String!]) {
			versions(filter: $filter, limit: $limit, offset: $offset, sort: $sort) {
				items {
					%s
				}
				total
			}
		}
	`, versionFields)

	err := v.
		sdk.
		Post(query,
			&resp,
			client.Var("filter", filter),
			client.Var("limit", limit),
			client.Var("offset", offset),
			client.Var("sort", sort))
	if err != nil {
		return nil, errors.Wrap(err, "twhelp sdk")
	}
	return &resp.Versions, nil
}
