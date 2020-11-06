package sdk

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/Kichiyaki/gqlgen-client/client"
)

var (
	ErrServerNameIsEmpty = fmt.Errorf("twhelp sdk: Server name is empty")
	commonODFields       = `
		rankAtt
		scoreAtt
		rankDef
		scoreDef
		rankTotal
		scoreTotal
	`
	playerODFields = fmt.Sprintf(`
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
		exists
		dailyGrowth
		joinedAt
		deletedAt
		bestRank
		bestRankAt
		mostPoints
		mostPointsAt
		mostVillages
		mostVillagesAt
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
		exists
		dominance
		createdAt
		deletedAt
		bestRank
		bestRankAt
		mostPoints
		mostPointsAt
		mostVillages
		mostVillagesAt
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
	versionFields = `
		code
		name
		host
		timezone
	`
)

type SDK struct {
	url             string
	client          *client.Client
	httpClient      *http.Client
	Version         *Version
	Server          *Server
	Player          *Player
	Tribe           *Tribe
	Village         *Village
	LiveEnnoblement *LiveEnnoblement
}

func New(url string) *SDK {
	sdk := &SDK{
		url: url,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
	sdk.client = client.New(http.HandlerFunc(sdk.doRequest))
	sdk.Version = &Version{sdk}
	sdk.Server = &Server{sdk}
	sdk.Player = &Player{sdk}
	sdk.Tribe = &Tribe{sdk}
	sdk.Village = &Village{sdk}
	sdk.LiveEnnoblement = &LiveEnnoblement{sdk}
	return sdk
}

func (sdk *SDK) Post(query string, response interface{}, options ...client.Option) error {
	return sdk.client.Post(minifyString(query), response, options...)
}

func (sdk *SDK) doRequest(w http.ResponseWriter, r *http.Request) {
	resp, err := sdk.httpClient.Post(sdk.url, "application/json", r.Body)
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
