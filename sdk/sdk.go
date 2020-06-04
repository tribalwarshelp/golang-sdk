package sdk

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/Kichiyaki/gqlgen-client/client"

	"github.com/pkg/errors"
	"github.com/tribalwarshelp/shared/models"
)

const (
	commonODFields = `rankAtt
scoreAtt
rankDef
scoreDef
rankTotal
scoreTotal`
)

var (
	ErrServerNameIsEmpty = fmt.Errorf("twhelp sdk: Server name is empty")
	playerODFields       = fmt.Sprintf(`
		%s
		rankSup
		scoreSup
	`, commonODFields)
	playerFields = fmt.Sprintf(`
		id
		name
		totalVillages
		points
		rank
		exist
		%s
	`, playerODFields)
	tribeFields = fmt.Sprintf(`
		id
		name
		tag
		totalMembers
		totalVillages
		points
		allPoints
		rank
		exist
		%s
	`, commonODFields)
	villageFields = `
		id
		name
		bonus
		points
		x
		y
	`
)

type SDK struct {
	uri          string
	client       *client.Client
	httpClient   *http.Client
	LangVersions *LangVersions
	Servers      *Servers
	Players      *Players
	Tribes       *Tribes
	Villages     *Villages
	Ennoblements *Ennoblements
}

func New(uri string) *SDK {
	sdk := &SDK{
		uri:        uri,
		httpClient: &http.Client{},
	}
	sdk.client = client.New(http.HandlerFunc(sdk.handler))
	sdk.LangVersions = &LangVersions{sdk}
	sdk.Servers = &Servers{sdk}
	sdk.Players = &Players{sdk}
	sdk.Tribes = &Tribes{sdk}
	sdk.Villages = &Villages{sdk}
	sdk.Ennoblements = &Ennoblements{sdk}
	return sdk
}

func (sdk *SDK) handler(w http.ResponseWriter, r *http.Request) {
	resp, err := sdk.httpClient.Post(sdk.uri, "application/json", r.Body)
	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte(err.Error()))
		return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte(err.Error()))
		return
	}
	w.WriteHeader(resp.StatusCode)
	w.Write(body)
}

type LangVersions struct {
	sdk *SDK
}

func (lv *LangVersions) Read(tag models.LanguageTag) (*models.LangVersion, error) {
	resp := struct {
		LangVersion *models.LangVersion `json:"langVersion" gqlgen:"langVersion"`
	}{}
	query := `
		query langVersion($tag: LanguageTag!) {
			langVersion(tag: $tag) {
				tag
				name
				host
				timezone
			}
		}
	`
	err := lv.sdk.client.Post(minify(query), &resp, client.Var("tag", tag))
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
	query := `
		query langVersions($filter: LangVersionFilter) {
			langVersions(filter: $filter) {
				items {
					tag
					name
					host
					timezone
				}
				total
			}
		}
	`

	err := lv.sdk.client.Post(minify(query), &resp, client.Var("filter", filter))
	if err != nil {
		return nil, errors.Wrap(err, "twhelp sdk")
	}
	return resp.LangVersions, nil
}

type Servers struct {
	sdk *SDK
}

func (ss *Servers) Read(key string) (*models.Server, error) {
	resp := struct {
		Server *models.Server `json:"server" gqlgen:"server"`
	}{}
	query := `
		query server($key: String!) {
			server(key: $key) {
				id
				key
				status
				langVersionTag
			}
		}
	`
	err := ss.sdk.client.Post(minify(query), &resp, client.Var("key", key))
	if err != nil {
		return nil, errors.Wrap(err, "twhelp sdk")
	}
	return resp.Server, nil
}

type ServersList struct {
	Items []*models.Server `json:"items" gqlgen:"items"`
	Total int              `json:"total" gqlgen:"total"`
}

func (ss *Servers) Browse(filter *models.ServerFilter) (*ServersList, error) {
	if filter == nil {
		filter = &models.ServerFilter{}
	}
	resp := struct {
		Servers *ServersList `json:"servers" gqlgen:"servers"`
	}{}

	query := `
		query servers($filter: ServerFilter) {
			servers(filter: $filter) {
				items {
					id
					key
					status
					langVersionTag
				}
				total
			}
		}
	`

	err := ss.sdk.client.Post(minify(query), &resp, client.Var("filter", filter))
	if err != nil {
		return nil, errors.Wrap(err, "twhelp sdk")
	}
	return resp.Servers, nil
}

type Players struct {
	sdk *SDK
}

type PlayerInclude struct {
	Tribe bool
}

func (incl PlayerInclude) String() string {
	i := ""
	if incl.Tribe {
		i += fmt.Sprintf(`
			tribe {
				%s
			}
		`, tribeFields)
	}
	return i
}

func (ps *Players) Read(server string, id int, include *PlayerInclude) (*models.Player, error) {
	if server == "" {
		return nil, ErrServerNameIsEmpty
	}
	if include == nil {
		include = &PlayerInclude{}
	}
	resp := struct {
		Player *models.Player `json:"player" gqlgen:"player"`
	}{}

	query := fmt.Sprintf(`
		query player($server: String!, $id: Int!) {
			player(server: $server, id: $id) {
				%s
				%s
			}
		}
	`, playerFields, include.String())
	err := ps.sdk.client.Post(minify(query), &resp, client.Var("server", server), client.Var("id", id))
	if err != nil {
		return nil, errors.Wrap(err, "twhelp sdk")
	}
	return resp.Player, nil
}

type PlayersList struct {
	Items []*models.Player `json:"items" gqlgen:"items"`
	Total int              `json:"total" gqlgen:"total"`
}

func (ps *Players) Browse(server string, filter *models.PlayerFilter, include *PlayerInclude) (*PlayersList, error) {
	if server == "" {
		return nil, ErrServerNameIsEmpty
	}
	if filter == nil {
		filter = &models.PlayerFilter{}
	}
	if include == nil {
		include = &PlayerInclude{}
	}
	resp := struct {
		Players *PlayersList `json:"players" gqlgen:"players"`
	}{}
	query := fmt.Sprintf(`
		query players($server: String!, $filter: PlayerFilter) {
			players(server: $server, filter: $filter) {
				items {
					%s
					%s
				}
				total
			}
		}
	`, playerFields, include.String())

	err := ps.sdk.client.Post(minify(query), &resp, client.Var("filter", filter), client.Var("server", server))
	if err != nil {
		return nil, errors.Wrap(err, "twhelp sdk")
	}
	return resp.Players, nil
}

type Tribes struct {
	sdk *SDK
}

func (ss *Tribes) Read(server string, id int) (*models.Tribe, error) {
	resp := struct {
		Tribe *models.Tribe `json:"tribe" gqlgen:"tribe"`
	}{}
	query := fmt.Sprintf(`
		query tribe($server: String!, $id: Int!) {
			tribe(server: $server, id: $id) {
				%s
			}
		}
	`, tribeFields)
	err := ss.sdk.client.Post(minify(query), &resp, client.Var("server", server), client.Var("id", id))
	if err != nil {
		return nil, errors.Wrap(err, "twhelp sdk")
	}
	return resp.Tribe, nil
}

type TribesList struct {
	Items []*models.Tribe `json:"items" gqlgen:"items"`
	Total int             `json:"total" gqlgen:"total"`
}

func (ss *Tribes) Browse(server string, filter *models.TribeFilter) (*TribesList, error) {
	if filter == nil {
		filter = &models.TribeFilter{}
	}
	resp := struct {
		Tribes *TribesList `json:"tribes" gqlgen:"tribes"`
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

	err := ss.sdk.client.Post(minify(query), &resp, client.Var("server", server), client.Var("filter", filter))
	if err != nil {
		return nil, errors.Wrap(err, "twhelp sdk")
	}
	return resp.Tribes, nil
}

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

func (ps *Villages) Read(server string, id int, include *VillageInclude) (*models.Village, error) {
	if server == "" {
		return nil, ErrServerNameIsEmpty
	}
	if include == nil {
		include = &VillageInclude{}
	}
	resp := struct {
		Village *models.Village `json:"village" gqlgen:"village"`
	}{}

	query := fmt.Sprintf(`
		query village($server: String!, $id: Int!) {
			village(server: $server, id: $id) {
				%s
				%s
			}
		}
	`, villageFields, include.String())
	err := ps.sdk.client.Post(minify(query), &resp, client.Var("server", server), client.Var("id", id))
	if err != nil {
		return nil, errors.Wrap(err, "twhelp sdk")
	}
	return resp.Village, nil
}

type VillagesList struct {
	Items []*models.Village `json:"items" gqlgen:"items"`
	Total int               `json:"total" gqlgen:"total"`
}

func (ps *Villages) Browse(server string, filter *models.VillageFilter, include *VillageInclude) (*VillagesList, error) {
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
		Villages *VillagesList
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

	err := ps.sdk.client.Post(minify(query), &resp, client.Var("filter", filter), client.Var("server", server))
	if err != nil {
		return nil, errors.Wrap(err, "twhelp sdk")
	}
	return resp.Villages, nil
}

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

func (ps *Ennoblements) Browse(server string, include *EnnoblementInclude) ([]*models.Ennoblement, error) {
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
	err := ps.sdk.client.Post(minify(query), &resp, client.Var("server", server))
	if err != nil {
		return nil, errors.Wrap(err, "twhelp sdk")
	}
	return resp.Ennoblements, nil
}

func minify(str string) string {
	return strings.Join(strings.Fields(str), " ")
}
