package main

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/andrepriyanto10/go-stripe/internal/cards"
)

type stripePayload struct {
	Curreny string `json:"curreny"`
	Amount  string `json:"amount"`
}

// omitempty artinya fields tidak akan dikirim ke payload jika valuenya 0, false, or ""
type jsonResponse struct {
	Ok      bool   `json:"ok"`
	Message string `json:"message,omitempty"`
	Content string `json:"content,omitempty"`
	ID      int    `json:"id,omitempty"`
}

func (app *application) GetPaymentIntent(w http.ResponseWriter, r *http.Request) {
	// payload is data received from request or response
	var payload stripePayload

	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		app.errorLog.Println(err)
		return
	}

	// change type amount string to integer
	amount, err := strconv.Atoi(payload.Amount)
	if err != nil {
		app.errorLog.Println(err)
		return
	}

	// inject secret key currency for cards
	card := cards.Card{
		Secret:   app.config.stripe.secret,
		Key:      app.config.stripe.key,
		Currency: payload.Curreny,
	}

	okey := true

	// pi : payment intent
	pi, msg, err := card.Charge(payload.Curreny, amount)
	if err != nil {
		okey = false
	}

	if okey {
		out, err := json.MarshalIndent(pi, "", " ")
		if err != nil {
			app.errorLog.Println(err)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(out)
	} else {
		j := jsonResponse{
			Ok:      false,
			Message: msg,
			Content: "",
		}

		// Indentasi adalah tata letak atau format untuk membuat struktur data lebih mudah dibaca dan dimengerti.
		// MarshalIndent mengubah format menjadi json
		out, err := json.MarshalIndent(j, "", " ")
		if err != nil {
			app.errorLog.Println(err)
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(out)
	}

}
