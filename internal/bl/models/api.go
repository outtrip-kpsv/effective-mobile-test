package models

type Nationalize struct {
  Count   int    `json:"count"`
  Name    string `json:"name"`
  Country []struct {
    CountryId   string  `json:"country_id"`
    Probability float64 `json:"probability"`
  } `json:"country"`
}

type Genderize struct {
  Count       int     `json:"count"`
  Name        string  `json:"name"`
  Gender      string  `json:"gender"`
  Probability float64 `json:"probability"`
}

type Agify struct {
  Count int    `json:"count"`
  Name  string `json:"name"`
  Age   int    `json:"age"`
}
