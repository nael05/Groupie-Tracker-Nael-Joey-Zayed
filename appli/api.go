package groupie

import (
	"encoding/json"
	"net/http"
)

type Artist struct {
	ID           int      `json:"id"`
	Image        string   `json:"image"`
	Name         string   `json:"name"`
	Members      []string `json:"members"`
	CreationDate int      `json:"creationDate"`
	FirstAlbum   string   `json:"firstAlbum"`
	Relations    string   `json:"relations"`
}

type List_artist struct {
	Name         string
	Members      []string
	CreationDate int
	FirstAlbum   string
	Image        string
	RelationsUrl string
}

type RelationsData struct {
	ID             int                 `json:"id"`
	DatesLocations map[string][]string `json:"datesLocations"`
}

func Api() map[int]List_artist {
	resp, err := http.Get("https://groupietrackers.herokuapp.com/api/artists")
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	var artists []Artist
	if err := json.NewDecoder(resp.Body).Decode(&artists); err != nil {
		panic(err)
	}

	list_artist := make(map[int]List_artist)

	for _, a := range artists {
		list_artist[a.ID] = List_artist{
			Name:         a.Name,
			Members:      a.Members,
			CreationDate: a.CreationDate,
			FirstAlbum:   a.FirstAlbum,
			Image:        a.Image,
			RelationsUrl: a.Relations,
		}
	}

	return list_artist
}

func GetRelations(url string) (map[string][]string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var relData RelationsData
	if err := json.NewDecoder(resp.Body).Decode(&relData); err != nil {
		return nil, err
	}
	return relData.DatesLocations, nil
}
