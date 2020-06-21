package main

import (
	"fmt"
	"log"
	"os"

	"github.com/tribalwarshelp/golang-sdk/sdk"
	"github.com/tribalwarshelp/shared/models"
)

func init() {
	os.Setenv("TZ", "UTC")
}

func main() {
	api := sdk.New("http://localhost:8080/graphql")

	langVersion, err := api.LangVersions.Read(models.LanguageTagPL)
	if err != nil {
		log.Fatal(err)
	}
	if err != nil {
		log.Fatal(err)
	}
	log.Println(langVersion.Name, langVersion.Tag, langVersion.Host, langVersion.Timezone)

	langVersionsList, err := api.LangVersions.Browse(&models.LangVersionFilter{
		HostMATCH: "plemiona%",
	})
	if err != nil {
		log.Fatal(err)
	}
	for _, langVersion := range langVersionsList.Items {
		log.Println(langVersion.Name, langVersion.Tag, langVersion.Host, langVersion.Timezone)
	}

	server, err := api.Servers.Read("pl151", &sdk.ServerInclude{LangVersion: true})
	if err != nil {
		log.Fatal(err)
	}
	log.Println(server.Key, server.Status, server.LangVersion.Tag)

	serversList, err := api.Servers.Browse(nil, nil)
	if err != nil {
		log.Fatal(err)
	}
	for _, server := range serversList.Items {
		log.Print()
		log.Println(server.Key, server.Status)
	}

	player, err := api.Players.Read("pl151", 699813215, nil)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(player.ID, player.Name, player.RankAtt, player.RankDef, player.RankSup)

	playersList, err := api.Players.Browse("pl151", &models.PlayerFilter{
		Sort:  "rank ASC",
		Limit: 10,
	}, &sdk.PlayerInclude{
		Tribe: true,
	})
	if err != nil {
		log.Fatal(err)
	}
	for _, player := range playersList.Items {
		log.Print()
		log.Println(player.ID, player.Name, player.RankAtt, player.RankDef, player.RankSup)
		if player.Tribe != nil {
			log.Println(player.Tribe.ID, player.Tribe.Tag)
		}
	}

	tribe, err := api.Tribes.Read("pl151", 894)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(tribe.ID, tribe.Name, tribe.Tag, tribe.RankAtt, tribe.RankDef)

	tribesList, err := api.Tribes.Browse("pl151", &models.TribeFilter{
		TagIEQ: ":.+.:",
	})
	if err != nil {
		log.Fatal(err)
	}
	for _, tribe := range tribesList.Items {
		log.Println(tribe.ID, tribe.Name, tribe.Tag, tribe.RankAtt, tribe.RankDef)
	}

	village, err := api.Villages.Read("pl151", 28299, &sdk.VillageInclude{
		Player: true,
		PlayerInclude: sdk.PlayerInclude{
			Tribe: true,
		},
	})
	if err != nil {
		log.Fatal(err)
	}
	log.Println(village.ID, village.Name)
	if village.Player != nil {
		log.Println(village.Player.ID, village.Player.Name)
		if village.Player.Tribe != nil {
			log.Println(village.Player.Tribe.ID, village.Player.Tribe.Tag)
		}
	}

	villageslist, err := api.Villages.Browse("pl151", &models.VillageFilter{
		PlayerID: []int{699270453},
		Sort:     "id ASC",
		Limit:    10,
	}, &sdk.VillageInclude{
		Player: true,
		PlayerInclude: sdk.PlayerInclude{
			Tribe: true,
		},
	})
	if err != nil {
		log.Fatal(err)
	}

	for _, village := range villageslist.Items {
		fmt.Print("\n\n")
		log.Println(village.ID, village.Name)
		if village.Player != nil {
			log.Println(village.Player.ID, village.Player.Name)
			if village.Player.Tribe != nil {
				log.Println(village.Player.Tribe.ID, village.Player.Tribe.Tag)
			}
		}
	}

	ennoblements, err := api.LiveEnnoblements.Browse("pl151", &sdk.LiveEnnoblementInclude{
		NewOwner: true,
		NewOwnerInclude: sdk.PlayerInclude{
			Tribe: true,
		},
		OldOwner: true,
		OldOwnerInclude: sdk.PlayerInclude{
			Tribe: true,
		},
		Village: true,
	})
	if err != nil {
		log.Fatal(err)
	}
	for _, ennoblement := range ennoblements {
		fmt.Print("\n\n", ennoblement.EnnobledAt.String(), "\n")
		if ennoblement.NewOwner != nil {
			log.Println(ennoblement.NewOwner.ID, ennoblement.NewOwner.Name)
			if ennoblement.NewOwner.Tribe != nil {
				log.Println(ennoblement.NewOwner.Tribe.ID, ennoblement.NewOwner.Tribe.Tag)
			}
		}
		if ennoblement.OldOwner != nil {
			log.Println(ennoblement.OldOwner.ID, ennoblement.OldOwner.Name)
			if ennoblement.OldOwner.Tribe != nil {
				log.Println(ennoblement.OldOwner.Tribe.ID, ennoblement.OldOwner.Tribe.Tag)
			}
		}
		if ennoblement.Village != nil {
			log.Println(ennoblement.Village.ID, ennoblement.Village.Name)
		}
	}
}
