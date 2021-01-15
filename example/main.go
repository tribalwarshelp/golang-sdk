package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/tribalwarshelp/golang-sdk/sdk"
	"github.com/tribalwarshelp/shared/models"
)

func init() {
	os.Setenv("TZ", "UTC")
}

func main() {
	api := sdk.New("https://api.tribalwarshelp.com/graphql")

	version, err := api.Version.Read(models.VersionCodePL)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(version.Name, version.Code, version.Host, version.Timezone)

	versionList, err := api.Version.Browse(0, 0, []string{}, &models.VersionFilter{
		HostMATCH: "plemiona%",
	})
	if err != nil {
		log.Fatal(err)
	}
	for _, version := range versionList.Items {
		log.Println(version.Name, version.Code, version.Host, version.Timezone)
	}

	server, err := api.Server.Read("pl151", &sdk.ServerInclude{Version: true})
	if err != nil {
		log.Fatal(err)
	}
	log.Println(server.Key, server.Status, server.Version.Code)

	serverList, err := api.Server.Browse(3, 0, []string{}, nil, nil)
	if err != nil {
		log.Fatal(err)
	}
	for _, server := range serverList.Items {
		log.Println(server.Key, server.Status)
	}

	player, err := api.Player.Read("pl151", 699813215, nil)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(player.ID, player.Name, player.RankAtt, player.RankDef, player.RankSup)

	playerList, err := api.Player.Browse("pl151",
		10,
		0,
		[]string{"rank ASC"},
		&models.PlayerFilter{},
		&sdk.PlayerInclude{
			Tribe: true,
		})
	if err != nil {
		log.Fatal(err)
	}
	for _, player := range playerList.Items {
		log.Println(player.ID, player.Name, player.RankAtt, player.RankDef, player.RankSup)
		if player.Tribe != nil {
			log.Println(player.Tribe.ID, player.Tribe.Tag)
		}
	}

	tribe, err := api.Tribe.Read("pl151", 894)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(tribe.ID, tribe.Name, tribe.Tag, tribe.RankAtt, tribe.RankDef)

	tribeList, err := api.Tribe.Browse("pl151", 10, 0, []string{}, &models.TribeFilter{
		TagIEQ: ":.+.:",
	})
	if err != nil {
		log.Fatal(err)
	}
	for _, tribe := range tribeList.Items {
		log.Println(tribe.ID, tribe.Name, tribe.Tag, tribe.RankAtt, tribe.RankDef)
	}

	village, err := api.Village.Read("pl151", 28120, &sdk.VillageInclude{
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

	villagelist, err := api.Village.Browse("pl151",
		10,
		0,
		[]string{"id ASC"},
		&models.VillageFilter{
			PlayerID: []int{699270453},
		}, &sdk.VillageInclude{
			Player: true,
			PlayerInclude: sdk.PlayerInclude{
				Tribe: true,
			},
		})
	if err != nil {
		log.Fatal(err)
	}

	for _, village := range villagelist.Items {
		fmt.Print("\n\n")
		log.Println(village.ID, village.Name)
		if village.Player != nil {
			log.Println(village.Player.ID, village.Player.Name)
			if village.Player.Tribe != nil {
				log.Println(village.Player.Tribe.ID, village.Player.Tribe.Tag)
			}
		}
	}

	ennoblements, err := api.Ennoblement.Browse("pl151",
		100,
		0,
		[]string{},
		&models.EnnoblementFilter{
			EnnobledAtGTE: time.Now().Add(-1 * time.Hour),
		},
		&sdk.EnnoblementInclude{
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
	for _, ennoblement := range ennoblements.Items {
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
