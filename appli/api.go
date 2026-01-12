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
}

type List_artist struct {
	Name         string
	Members      []string
	CreationDate int
	FirstAlbum   string
	Image		string
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
		}
	}

	return list_artist
}
