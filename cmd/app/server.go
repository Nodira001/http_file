package app

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"path"
	"strconv"

	"github.com/Nodira001/http/pkg/banners"
)

type Server struct {
	mux     *http.ServeMux
	service *banners.Service
}

func NewServer(mux *http.ServeMux, service *banners.Service) *Server {
	return &Server{mux: mux, service: service}
}

func (s *Server) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	s.mux.ServeHTTP(writer, request)
}

func (s *Server) Init() {
	s.mux.HandleFunc("/banners.getAll", s.handleGetAllBanners)
	s.mux.HandleFunc("/banners.getById", s.handleGetBannerByID)
	s.mux.HandleFunc("/banners.save", s.handleSaveBanner)
	s.mux.HandleFunc("/banners.removeById", s.handleRemoveBannerByID)
}
func (s *Server) handleGetAllBanners(writer http.ResponseWriter, request *http.Request) {
	items, err := s.service.All(request.Context())
	if err != nil {
		log.Println(err)
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
	}
	data, err := json.Marshal(items)
	if err != nil {
		log.Println(err)
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	writer.Header().Set("Content-Type", "application/json")
	_, err = writer.Write(data)
	if err != nil {
		log.Println(err)
	}
}
func (s *Server) handleGetBannerByID(writer http.ResponseWriter, request *http.Request) {
	idParam := request.URL.Query().Get("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		log.Println(err)
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	item, err := s.service.ByID(request.Context(), id)
	if err != nil {
		log.Println(err)
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	data, err := json.Marshal(item)
	if err != nil {
		log.Println(err)
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	writer.Header().Set("Content-Type", "application/json")
	_, err = writer.Write(data)
	if err != nil {
		log.Println(err)
	}
}
func (s *Server) handleSaveBanner(writer http.ResponseWriter, request *http.Request) {
	idParam := request.FormValue("id")
	title := request.FormValue("title")
	content := request.FormValue("content")
	button := request.FormValue("button")
	link := request.FormValue("link")

	id, err := strconv.ParseInt(idParam, 10, 64)

	if err != nil {
		log.Println(err)
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	//10 - mb

	err = request.ParseMultipartForm(10 * 1024 * 1024)

	if err != nil {
		log.Println(err)

	}

	request.ParseForm()

	image, handler, _ := request.FormFile("image")

	if image != nil {
		defer image.Close()
	}

	img := ""

	if handler != nil {
		img = fmt.Sprint(id) + path.Ext(handler.Filename)
	}

	newBanner := &banners.Banner{ID: id, Button: button, Content: content, Link: link, Title: title, Image: img}
	banner, err := s.service.Save(request.Context(), newBanner, image, handler)

	if err != nil {
		log.Println(err)
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	data, err := json.Marshal(banner)
	if err != nil {
		log.Println(err)
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	writer.Header().Set("Content-Type", "application/json")
	_, err = writer.Write(data)
	if err != nil {
		log.Println(err)
	}

	return
}

func (s *Server) handleRemoveBannerByID(writer http.ResponseWriter, request *http.Request) {
	idParam := request.URL.Query().Get("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		log.Println(err)
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	banner, e := s.service.RemoveByID(request.Context(), id)
	if e != nil {
		log.Println(e)
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	data, err := json.Marshal(banner)
	if err != nil {
		log.Println(err)
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusBadRequest)
		return
	}
	writer.Header().Set("Content-Type", "application/json")
	_, err = writer.Write(data)
	if err != nil {
		log.Println(err)
	}

	return
}
