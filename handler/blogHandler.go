package handler

import (
	"blog-microservice/model"
	"blog-microservice/repo"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type KeyProduct struct{}

type BlogHandler struct {
	logger *log.Logger
	repo   *repo.BlogRepo
}

func NewBlogHandler(l *log.Logger, r *repo.BlogRepo) *BlogHandler {
	return &BlogHandler{l, r}
}

func (b *BlogHandler) GetAllBlogs(rw http.ResponseWriter, h *http.Request) {
	blogs, err := b.repo.GetAll()
	if err != nil {
		b.logger.Print("Database exception: ", err)
	}

	if blogs == nil {
		return
	}

	err = blogs.ToJSON(rw)
	if err != nil {
		http.Error(rw, "Unable to convert to json", http.StatusInternalServerError)
		b.logger.Fatal("Unable to convert to json :", err)
		return
	}
}

func (b *BlogHandler) GetBlogById(rw http.ResponseWriter, h *http.Request) {
	vars := mux.Vars(h)
	id := vars["id"]

	blog, err := b.repo.Get(id)
	if err != nil {
		b.logger.Print("Database exception: ", err)
	}

	if blog == nil {
		http.Error(rw, "Blog with given id not found", http.StatusNotFound)
		b.logger.Printf("Blog with id: '%s' not found", id)
		return
	}

	err = blog.ToJSON(rw)
	if err != nil {
		http.Error(rw, "Unable to convert to json", http.StatusInternalServerError)
		b.logger.Fatal("Unable to convert to json :", err)
		return
	}
}

func (b *BlogHandler) PostBlog(rw http.ResponseWriter, h *http.Request) {
	blog := h.Context().Value(KeyProduct{}).(*model.Blog)
	b.repo.Insert(blog)
	rw.WriteHeader(http.StatusCreated)
}

func (b *BlogHandler) UpdateBlog(rw http.ResponseWriter, h *http.Request) {
	vars := mux.Vars(h)
	id := vars["id"]
	blog := h.Context().Value(KeyProduct{}).(*model.Blog)

	b.repo.Update(id, blog)
	rw.WriteHeader(http.StatusOK)
}

func (b *BlogHandler) DeleteBlog(rw http.ResponseWriter, h *http.Request) {
	vars := mux.Vars(h)
	id := vars["id"]

	b.repo.Delete(id)
	rw.WriteHeader(http.StatusNoContent)
}

func (h *BlogHandler) GetAllVotes(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	blogID := vars["id"]

	blog, err := h.repo.Get(blogID)
	if err != nil {
		http.Error(w, "Failed to retrieve the blog", http.StatusInternalServerError)
		return
	}

	var votes []*model.Vote
	if blog.Votes != nil {
		// Konvertujemo blog.Votes u []*model.Vote
		for _, v := range blog.Votes {
			votes = append(votes, &v)
		}
	} else {
		votes = []*model.Vote{}
	}

	jsonVotes, err := json.Marshal(votes)
	if err != nil {
		http.Error(w, "Failed to serialize votes to JSON", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	_, err = w.Write(jsonVotes)
	if err != nil {
		log.Println("Failed to write JSON response:", err)
	}
}

func (h *BlogHandler) GetVotesCount(w http.ResponseWriter, r *http.Request) {
	// Extract the blog ID from the request URL parameters
	vars := mux.Vars(r)
	blogID := vars["id"]

	// Retrieve the blog from the repository by its ID
	blog, err := h.repo.Get(blogID)
	if err != nil {
		// Handle the error (e.g., return an error response)
		http.Error(w, "Failed to retrieve the blog", http.StatusInternalServerError)
		return
	}

	if len(blog.Votes) == 0 || (blog.Votes == nil) {
		fmt.Fprint(w, "0")
		return
	}

	votesCount := 0

	for _, vote := range blog.Votes {
		if vote.IsUpvote {
			votesCount++
		} else {
			votesCount--
		}
	}

	fmt.Fprintf(w, "%d", votesCount)
}

func (b *BlogHandler) AddVote(rw http.ResponseWriter, h *http.Request) {
	vars := mux.Vars(h)
	id := vars["id"]
	b.logger.Print("Pre bodyija!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!: ")

	vote := h.Context().Value(KeyProduct{}).(*model.Vote)

	b.logger.Print("Pre ulaza u repo: ")
	b.repo.AddVote(id, vote)
	b.logger.Print("Posle ulaza u repo: ")
	rw.WriteHeader(http.StatusOK)
}

func (b *BlogHandler) ChangeVote(rw http.ResponseWriter, h *http.Request) {
	vars := mux.Vars(h)
	id := vars["id"]
	index, err := strconv.Atoi(vars["index"])
	if err != nil {
		http.Error(rw, "Unable to decode index", http.StatusBadRequest)
		b.logger.Fatal(err)
		return
	}

	var vote model.Vote
	d := json.NewDecoder(h.Body)
	d.Decode(&vote)

	b.repo.ChangeVote(id, index, &vote)
	rw.WriteHeader(http.StatusOK)
}

func (b *BlogHandler) GetBlogsByAuthorId(rw http.ResponseWriter, h *http.Request) {
	vars := mux.Vars(h)
	userId := vars["id"]
	id64, err := strconv.ParseUint(userId, 10, 32)
	id := uint32(id64)

	blogs, err := b.repo.GetByAuthorId(id)
	if err != nil {
		b.logger.Print("Database exception: ", err)
	}

	if blogs == nil {
		return
	}

	err = blogs.ToJSON(rw)
	if err != nil {
		http.Error(rw, "Unable to convert to json", http.StatusInternalServerError)
		b.logger.Fatal("Unable to convert to json :", err)
		return
	}
}

func (b *BlogHandler) MiddlewareBlogDeserialization(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, h *http.Request) {
		blog := &model.Blog{}        //pravim pokazivac na Patient strukturu
		err := blog.FromJSON(h.Body) // radimo deserijalizaciju is jsona iz sadrzaja koji nam stize u bodiju zahteva putem tog pokazivaca
		if err != nil {
			http.Error(rw, "Unable to decode json", http.StatusBadRequest)
			b.logger.Fatal(err)
			return
		}
		//ako deserijalizacija uspe treba da iskoristimo context with value
		//on koristi ovaj context iz requesta i  da na dati key napapira patienta
		//keyProduct - kljuc, patient - vrednost u okviru WithValue contexta
		//key je samo prazna struktura, a ovaj kontekst? je kao neka hesh mapa koja mapira kljuc na vrednost
		ctx := context.WithValue(h.Context(), KeyProduct{}, blog)
		// i onda kazemo da je request request sa novim kontekstom koji smo napravili
		h = h.WithContext(ctx)
		//potom zahtev prosledjujemo dalje na izvrsavanje (onoj metoui u mainu)
		next.ServeHTTP(rw, h)
	})
}

func (b *BlogHandler) MiddlewareVoteDeserialization(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, h *http.Request) {
		vote := &model.Vote{}
		err := vote.FromJSON(h.Body)
		if err != nil {
			http.Error(rw, "Unable to decode json", http.StatusBadRequest)
			b.logger.Fatal(err)
			return
		}

		ctx := context.WithValue(h.Context(), KeyProduct{}, vote)
		h = h.WithContext(ctx)

		next.ServeHTTP(rw, h)
	})
}

func (b *BlogHandler) MiddlewareContentTypeSet(next http.Handler) http.Handler { //setujemo content type u header
	//ova fja vraca http Handler, a njega ce instancirati http.HandlerFunc kome se prosledi fja sa ResponseWriterom i requestom
	return http.HandlerFunc(func(rw http.ResponseWriter, h *http.Request) {
		//ovo je nesto sto se desava pre gadjanja svakog od handlera u mainu
		b.logger.Println("Method [", h.Method, "] - Hit path :", h.URL.Path)

		rw.Header().Add("Content-Type", "application/json") // u headeru se kaci content type

		next.ServeHTTP(rw, h) // prosledi zahtev dalje na obradu kome treba
	})
}
