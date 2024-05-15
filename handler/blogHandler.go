package handler

import (
	"blog-microservice/mapper"
	"blog-microservice/model"
	"blog-microservice/repo"
	"context"

	//"encoding/json"
	"fmt"
	"log"
	"soa/grpc/proto/blog"
	"strconv"
)

type KeyProduct struct{}

type BlogHandler struct {
	blog.UnimplementedBlogServiceServer
	logger *log.Logger
	repo   *repo.BlogRepo
}

func NewBlogHandler(l *log.Logger, r *repo.BlogRepo) *BlogHandler {
	return &BlogHandler{logger: l, repo: r}
}

// ???????
func (b *BlogHandler) GetAllBlogs(ctx context.Context, request *blog.GetAllRequest) (*blog.GetBlogsResponse, error) {
	modelBlogs, err := b.repo.GetAll()
	if err != nil {
		b.logger.Print("Database exception: ", err)
	}

	var blogs []model.Blog

	for _, b := range modelBlogs {
		blogs = append(blogs, *b)
	}

	protoBlogs := mapper.MapSliceToProtoBlogs(blogs)
	response := &blog.GetBlogsResponse{
		Blogs: protoBlogs,
	}
	return response, nil

}

func (b *BlogHandler) GetBlogById(ctx context.Context, request *blog.GetByIdRequest) (*blog.BlogResponse, error) {
	id := request.Id

	blogM, err := b.repo.Get(id)
	if err != nil {
		b.logger.Print("Database exception: ", err)
	}

	if blogM == nil {
		b.logger.Printf("Blog with id: '%s' not found", id)
		return nil, fmt.Errorf("blog with given id not found")
	}

	protoBlog := mapper.MapToPBlog(blogM)
	response := &blog.BlogResponse{
		Blog: protoBlog,
	}
	return response, nil
}

/*
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
}*/

func (b *BlogHandler) PostBlog(ctx context.Context, request *blog.CreateBlogRequest) (*blog.BlogResponse, error) {
	newBlog := mapper.MapToBlog(request.Blog)

	err := b.repo.Insert(newBlog)
	if err != nil {
		return nil, err
	}

	protoBlog := mapper.MapToPBlog(newBlog)
	response := &blog.BlogResponse{
		Blog: protoBlog,
	}

	return response, nil
}

/*
func (b *BlogHandler) PostBlog(rw http.ResponseWriter, h *http.Request) {
	blog := h.Context().Value(KeyProduct{}).(*model.Blog)
	b.repo.Insert(blog)
	rw.WriteHeader(http.StatusCreated)
}
*/

func (b *BlogHandler) UpdateBlog(ctx context.Context, request *blog.UpdateBlogRequest) (*blog.BlogResponse, error) {
	id := request.Id

	blogUpdate := ctx.Value(KeyProduct{}).(*model.Blog)

	b.repo.Update(id, blogUpdate)

	response := &blog.BlogResponse{
		Blog: nil,
	}
	return response, nil
}

/*
func (b *BlogHandler) UpdateBlog(rw http.ResponseWriter, h *http.Request) {
	vars := mux.Vars(h)

	id := vars["id"]
	blog := h.Context().Value(KeyProduct{}).(*model.Blog)

	b.repo.Update(id, blog)
	rw.WriteHeader(http.StatusOK)
}
*/
// ?????????????????????????
func (b *BlogHandler) DeleteBlog(ctx context.Context, request *blog.GetByIdRequest) (*blog.BlogResponse, error) {
	id := request.Id

	b.repo.Delete(id)
	response := &blog.BlogResponse{
		Blog: nil,
	}
	return response, nil
}

/*
func (b *BlogHandler) DeleteBlog(rw http.ResponseWriter, h *http.Request) {
	vars := mux.Vars(h)
	id := vars["id"]

	b.repo.Delete(id)
	rw.WriteHeader(http.StatusNoContent)
}
*/

func (h *BlogHandler) GetAllVotes(ctx context.Context, request *blog.GetByIdRequest) (*blog.GetAllVotesResponse, error) {
	blogID := request.Id

	modelBlog, err := h.repo.Get(blogID)
	if err != nil {
		return nil, err
	}

	var votes []model.Vote
	if modelBlog.Votes != nil {
		for _, v := range modelBlog.Votes {
			votes = append(votes, v)
		}
	} else {
		votes = []model.Vote{}
	}

	protoVotes := mapper.MapSliceToProtoVotes(votes)

	response := &blog.GetAllVotesResponse{
		Votes: protoVotes,
	}
	return response, nil
}

/*
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
*/

func (h *BlogHandler) GetVotesCount(ctx context.Context, request *blog.GetByIdRequest) (*blog.GetVotesCountResponse, error) {
	blogID := request.Id

	modelBlog, err := h.repo.Get(blogID)
	if err != nil {
		return nil, err
	}

	if len(modelBlog.Votes) == 0 || (modelBlog.Votes == nil) {
		return nil, nil
	}

	votesCount := 0

	for _, vote := range modelBlog.Votes {
		if vote.IsUpvote {
			votesCount++
		} else {
			votesCount--
		}
	}

	response := &blog.GetVotesCountResponse{
		Count: uint32(votesCount),
	}

	return response, nil
}

/*
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
*/

func (b *BlogHandler) AddVote(ctx context.Context, request *blog.AddVoteRequest) (*blog.VoteResponse, error) {
	id := request.Id
	b.logger.Print("Pre bodyija!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!: ")

	vote := mapper.MapToVote(request.Vote)

	b.logger.Print("Pre ulaza u repo: ")
	b.repo.AddVote(id, vote)
	b.logger.Print("Posle ulaza u repo: ")
	protoVote := mapper.MapToPVote(vote)
	response := &blog.VoteResponse{
		Vote: protoVote,
	}

	return response, nil
}

/*
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
*/

func (b *BlogHandler) ChangeVote(ctx context.Context, request *blog.ChangeVoteRequest) (*blog.VoteResponse, error) {
	id := request.Id
	index := int(request.Index)
	vote := mapper.MapToVote(request.Vote)

	b.repo.ChangeVote(id, index, vote)

	protoVote := mapper.MapToPVote(vote)

	response := &blog.VoteResponse{
		Vote: protoVote,
	}

	return response, nil
}

/*
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
*/

func (b *BlogHandler) GetBlogsByAuthorId(ctx context.Context, request *blog.GetByIdRequest) (*blog.GetBlogsResponse, error) {
	userId := request.Id
	id64, err := strconv.ParseUint(userId, 10, 32)
	if err != nil {
		b.logger.Print("Database exception: ", err)
	}
	id := uint32(id64)

	modelBlogs, err := b.repo.GetByAuthorId(id)
	if err != nil {
		b.logger.Print("Database exception: ", err)
	}

	if modelBlogs == nil {
		return nil, nil
	}

	var blogs []model.Blog

	for _, b := range modelBlogs {
		blogs = append(blogs, *b)
	}

	protoBlogs := mapper.MapSliceToProtoBlogs(blogs)
	response := &blog.GetBlogsResponse{
		Blogs: protoBlogs,
	}
	return response, nil
}

/*
func (b *BlogHandler) GetBlogsByAuthorId(rw http.ResponseWriter, h *http.Request) {
	vars := mux.Vars(h)
	userId := vars["id"]
	id64, err := strconv.ParseUint(userId, 10, 32)
	if err != nil {
		b.logger.Print("Database exception: ", err)
	}
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
*/

func (b *BlogHandler) AddComment(ctx context.Context, request *blog.AddCommentRequest) (*blog.CommentResponse, error) {
	id := request.Id
	comment := mapper.MapToComment(request.Comment)

	b.logger.Print("Pre ulaza u repo: ")
	b.repo.AddComment(id, comment)
	b.logger.Print("Posle ulaza u repo: ")

	protoComment := mapper.MapToPComment(comment)

	response := &blog.CommentResponse{
		Comment: protoComment,
	}

	return response, nil

}

func (b *BlogHandler) UpdateComment(ctx context.Context, request *blog.UpdateCommentRequest) (*blog.CommentResponse, error) {
	id := request.Id
	comment := mapper.MapToComment(request.Comment)
	index := int(request.Index)

	b.repo.UpdateComment(id, index, comment)

	protoComment := mapper.MapToPComment(comment)

	response := &blog.CommentResponse{
		Comment: protoComment,
	}

	return response, nil
}

func (b *BlogHandler) DeleteComment(ctx context.Context, request *blog.DeleteCommentRequest) (*blog.CommentResponse, error) {
	id := request.Id
	index := string(request.Index)

	b.repo.DeleteComment(id, index)
	response := &blog.CommentResponse{
		Comment: nil,
	}

	return response, nil

}

/*
func (b *BlogHandler) DeleteComment(rw http.ResponseWriter, h *http.Request) {
	vars := mux.Vars(h)
	blogId := vars["blogId"]
	commentId := vars["commentId"]

	err := b.repo.DeleteComment(blogId, commentId)

	if err != nil {
		http.Error(rw, "Unable to delete comment", http.StatusBadRequest)
		b.logger.Fatal(err)
		return
	}
	rw.WriteHeader(http.StatusOK)
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

func (b *BlogHandler) MiddlewareCommentDeserialization(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, h *http.Request) {
		comment := &model.Comment{}
		err := comment.FromJSON(h.Body)
		if err != nil {
			http.Error(rw, "Unable to decode JSON", http.StatusBadRequest)
			b.logger.Fatal(err)
			return
		}

		ctx := context.WithValue(h.Context(), KeyProduct{}, comment)
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
*/
