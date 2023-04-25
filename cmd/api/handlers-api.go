package main

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/andrepriyanto10/go-stripe/internal/cards"
	"github.com/go-chi/chi/v5"
)

type stripePayload struct {
	Currency string `json:"currency"`
	Amount   string `json:"amount"`
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

	//Encoder digunakan untuk mengambil data dari sebuah struktur Golang dan mengubahnya menjadi format JSON.
	//decoder digunakan untuk mengambil data dari format JSON dan mengembalikannya ke struktur Golang.
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
		Currency: payload.Currency,
	}

	okey := true

	// pi : payment intent
	pi, msg, err := card.Charge(payload.Currency, amount)
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

func (app *application) GetWidgetById(w http.ResponseWriter, r *http.Request) {
	// get id from url
	id := chi.URLParam(r, "id")
	// change type string to int
	widgetID, _ := strconv.Atoi(id)
	widget, err := app.DB.GetWidget(widgetID)
	if err != nil {
		app.errorLog.Println(err)
		return
	}

	out, err := json.MarshalIndent(widget, "", " ")
	if err != nil {
		app.errorLog.Println(err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(out)
}
