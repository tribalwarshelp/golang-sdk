package sdk

import (
	"fmt"

	"github.com/Kichiyaki/gqlgen-client/client"
	"github.com/pkg/errors"
	"github.com/tribalwarshelp/shared/models"
)

type LangVersions struct {
	sdk *SDK
}

func (lv *LangVersions) Read(tag models.LanguageTag) (*models.LangVersion, error) {
	resp := struct {
		LangVersion *models.LangVersion `json:"langVersion" gqlgen:"langVersion"`
	}{}
	query := fmt.Sprintf(`
		query langVersion($tag: LanguageTag!) {
			langVersion(tag: $tag) {
				%s
			}
		}
	`, langVersionFields)
	err := lv.sdk.client.Post(minifyString(query), &resp, client.Var("tag", tag))
	if err != nil {
		return nil, errors.Wrap(err, "twhelp sdk")
	}
	return resp.LangVersion, nil
}

type LangVersionsList struct {
	Items []*models.LangVersion `json:"items" gqlgen:"items"`
	Total int                   `json:"total" gqlgen:"total"`
}

func (lv *LangVersions) Browse(filter *models.LangVersionFilter) (*LangVersionsList, error) {
	if filter == nil {
		filter = &models.LangVersionFilter{}
	}
	resp := struct {
		LangVersions *LangVersionsList `json:"langVersions" gqlgen:"langVersions"`
	}{}
	query := fmt.Sprintf(`
		query langVersions($filter: LangVersionFilter) {
			langVersions(filter: $filter) {
				items {
					%s
				}
				total
			}
		}
	`, langVersionFields)

	err := lv.sdk.client.Post(minifyString(query), &resp, client.Var("filter", filter))
	if err != nil {
		return nil, errors.Wrap(err, "twhelp sdk")
	}
	return resp.LangVersions, nil
}
