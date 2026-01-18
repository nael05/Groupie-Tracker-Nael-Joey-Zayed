package groupie

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strconv"
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

type geocodeResponse []struct {
	Lat         string `json:"lat"`
	Lon         string `json:"lon"`
	DisplayName string `json:"display_name"`
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

func GeocodeLocation(query string) (float64, float64, string, error) {
	values := url.Values{}
	values.Set("q", query)
	values.Set("format", "json")
	values.Set("limit", "1")

	req, err := http.NewRequest("GET", "https://nominatim.openstreetmap.org/search", nil)
	if err != nil {
		return 0, 0, "", err
	}
	req.URL.RawQuery = values.Encode()
	req.Header.Set("User-Agent", "GroupieTracker/1.0")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return 0, 0, "", err
	}
	defer resp.Body.Close()

	var geo geocodeResponse
	if err := json.NewDecoder(resp.Body).Decode(&geo); err != nil {
		return 0, 0, "", err
	}
	if len(geo) == 0 {
		return 0, 0, "", nil
	}

	lat, err := strconv.ParseFloat(geo[0].Lat, 64)
	if err != nil {
		return 0, 0, "", err
	}
	lon, err := strconv.ParseFloat(geo[0].Lon, 64)
	if err != nil {
		return 0, 0, "", err
	}

	return lat, lon, geo[0].DisplayName, nil
}
