package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"url-shortner/internal/models"
	"url-shortner/internal/services"
)

type Handler struct {
	services *services.URLService
}

func NewHandler(services services.URLService) *Handler {
	return &Handler{
		services: &services,
	}
}

// POST /
func (h *Handler) GenerateUrl(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var longUrl models.Url

	err :=json.NewDecoder(r.Body).Decode(&longUrl)

	shortUrl, err := h.services.Shorten(r.Context(), longUrl.Url)
	if err != nil {
		log.Printf("ERROR: generate short url failed: %v", err)
		http.Error(w, "cannot generate short url", http.StatusInternalServerError)
		return
	}

	var response models.Url
	response.Url = shortUrl

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"status":   "success",
		"shortUrl": "http://localhost:3000/" + response.Url,
	})
}

// GET /
func (h *Handler) GetUrl(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	shortUrl := r.PathValue("shortUrl")

	originalUrl, err := h.services.Resolve(r.Context(), shortUrl)
	if err != nil {
		http.Error(w, "url does not exist", http.StatusNotFound)
		return
	}

	http.Redirect(w, r, originalUrl, http.StatusPermanentRedirect)

}
