package structures

type ArtistsAPI struct {
	Artists   string `json:"artists"`   //api URL
	Locations string `json:"locations"` //api URL
	Dates     string `json:"dates"`     //api URL
	Relation  string `json:"relation"`  //api URL
}

type Artists struct {
	Error        string
	ArtistsAPI   ArtistsAPI
	ArtistsArray []Artist
}

type Artist struct {
	Error           string
	Id              int      `json:"id"`
	Image           string   `json:"image"`
	Name            string   `json:"name"`
	Members         []string `json:"members"`
	CreationDate    int      `json:"creationDate"`
	FirstAlbum      string   `json:"firstAlbum"`
	LocationsAPI    string   `json:"locations"`    //api URL
	ConcertDatesAPI string   `json:"concertDates"` //api URL
	RelationsAPI    string   `json:"relations"`    //api URL
}

type ArtistsLocation struct {
	Id        int      `json:"id"`
	Locations []string `json:"locations"`
	DatesAPI  string   `json:"dates"` //api URL
}

type ArtistsConcertDates struct {
	Id    int      `json:"id"`
	Dates []string `json:"dates"`
}

type ArtistsRealtion struct {
	Id            int                 `json:"id"`
	DatesLocation map[string][]string `json:"datesLocations"`
}

type ArtistFullData struct {
	Error        string
	Id           int                 `json:"id"`
	Image        string              `json:"image"`
	Name         string              `json:"name"`
	Members      []string            `json:"members"`
	CreationDate int                 `json:"creationDate"`
	FirstAlbum   string              `json:"firstAlbum"`
	Locations    ArtistsLocation     `json:"artistsLocations"`
	ConcertDates ArtistsConcertDates `json:"artistsConcertDates"`
	Relations    ArtistsRealtion     `json:"artistsRelations"`
}
