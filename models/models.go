package models

type CurrencyResponse struct {
	Bitcoin struct {
		UAH float64 `json:"uah"`
	} `json:"bitcoin"`
}

type Subscriber struct {
	Email string `json:"email"`
}