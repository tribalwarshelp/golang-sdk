package sdk

import (
	"fmt"
	"io/ioutil"
	"net/http"

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
	url          string
	client       *client.Client
	httpClient   *http.Client
	LangVersions *LangVersions
	Servers      *Servers
	Players      *Players
	Tribes       *Tribes
	Villages     *Villages
	Ennoblements *Ennoblements
}

func New(url string) *SDK {
	sdk := &SDK{
		url:        url,
		httpClient: &http.Client{},
	}
	sdk.client = client.New(http.HandlerFunc(sdk.makeRequest))
	sdk.LangVersions = &LangVersions{sdk}
	sdk.Servers = &Servers{sdk}
	sdk.Players = &Players{sdk}
	sdk.Tribes = &Tribes{sdk}
	sdk.Villages = &Villages{sdk}
	sdk.Ennoblements = &Ennoblements{sdk}
	return sdk
}

func (sdk *SDK) makeRequest(w http.ResponseWriter, r *http.Request) {
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
