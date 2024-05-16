package mapper

import (
	"blog-microservice/model"
	p "blog-microservice/proto/blog"

	"google.golang.org/protobuf/types/known/timestamppb"
)

func MapToPBlog(blog *model.Blog) *p.Blog {
	blogP := &p.Blog{
		Id:           blog.ID.String(),
		UserId:       blog.UserID,
		Username:     blog.Username,
		Title:        blog.Title,
		Description:  blog.Description,
		CreationTime: timestamppb.New(blog.CreationTime),
		Status:       p.Blog_Status(blog.Status),
		Image:        blog.Image,
		Category:     p.Blog_Category(blog.Category),
	}

	for _, comment := range blog.Comments {
		blogP.Comments = append(blogP.Comments, MapToPComment(&comment))
	}

	for _, vote := range blog.Votes {
		blogP.Votes = append(blogP.Votes, MapToPVote(&vote))
	}

	return blogP
}

func MapToPComment(comment *model.Comment) *p.Comment {
	commentP := &p.Comment{
		Id:               comment.ID.String(),
		UserId:           comment.UserID,
		Username:         comment.Username,
		Text:             comment.Text,
		CreationTime:     timestamppb.New(comment.CreationTime),
		LastModification: timestamppb.New(comment.LastModification),
	}

	return commentP

}

func MapToPVote(vote *model.Vote) *p.Vote {
	voteP := &p.Vote{
		Id:           vote.ID.String(),
		IsUpvote:     vote.IsUpvote,
		UserId:       vote.UserID,
		CreationTime: timestamppb.New(vote.CreationTime),
	}

	return voteP
}

func MapToBlog(blogP *p.Blog) *model.Blog {
	blog := &model.Blog{
		UserID:       blogP.UserId,
		Username:     blogP.Username,
		Title:        blogP.Title,
		Description:  blogP.Description,
		CreationTime: blogP.CreationTime.AsTime(),
		Status:       model.Status(blogP.Status),
		Image:        blogP.Image,
		Category:     model.Category(blogP.Category),
	}

	for _, commentP := range blogP.Comments {
		blog.Comments = append(blog.Comments, *MapToComment(commentP))
	}

	for _, voteP := range blogP.Votes {
		blog.Votes = append(blog.Votes, *MapToVote(voteP))
	}

	return blog
}

func MapToComment(commentP *p.Comment) *model.Comment {
	comment := &model.Comment{
		UserID:           commentP.UserId,
		Username:         commentP.Username,
		Text:             commentP.Text,
		CreationTime:     commentP.CreationTime.AsTime(),
		LastModification: commentP.LastModification.AsTime(),
	}

	return comment

}

func MapToVote(voteP *p.Vote) *model.Vote {
	vote := &model.Vote{
		IsUpvote:     voteP.IsUpvote,
		UserID:       voteP.UserId,
		CreationTime: voteP.CreationTime.AsTime(),
	}

	return vote
}

func MapSliceToProtoBlogs(modelBlogs []model.Blog) []*p.Blog {
	var protoBlogs []*p.Blog

	for _, modelBlog := range modelBlogs {
		protoBlog := MapToPBlog(&modelBlog)
		protoBlogs = append(protoBlogs, protoBlog)
	}

	return protoBlogs
}

func MapSliceToProtoVotes(modelVotes []model.Vote) []*p.Vote {
	var protoVotes []*p.Vote

	for _, modelVote := range modelVotes {
		protoVote := MapToPVote(&modelVote)
		protoVotes = append(protoVotes, protoVote)
	}

	return protoVotes
}
