package rate

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/FEDQK/btc-api/constants"
	"github.com/FEDQK/btc-api/models"
)

func Get(w http.ResponseWriter, r *http.Request) {
	resp, err := http.Get(constants.PRICE_API_URL)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	var currencyResponse models.CurrencyResponse
	dec := json.NewDecoder(resp.Body)
	err = dec.Decode(&currencyResponse)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "BTC to UAH rate: %.2f", currencyResponse.Bitcoin.UAH)
}