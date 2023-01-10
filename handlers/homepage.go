package handlers

import (
	"net/http"

	"github.com/macedo/whatsappbot/whatsapp"
)

type homePageData struct {
	Clients []*whatsapp.Client
}

func HomePage() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data := &homePageData{
			Clients: whatsapp.Clients(),
		}

		err := renderPage(w, data)
		if err != nil {
			http.Error(w, "Ooops", http.StatusBadRequest)
			return
		}
	}
}
