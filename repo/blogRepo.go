package repo

import (
	"blog-microservice/model"
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type BlogRepo struct {
	cli    *mongo.Client
	logger *log.Logger
}

// NoSQL: Constructor which reads db configuration from environment
func New(ctx context.Context, logger *log.Logger) (*BlogRepo, error) {
	dburi := os.Getenv("MONGO_DB_URI")

	clientOptions := options.Client().ApplyURI(dburi)

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, err
	}

	return &BlogRepo{ //sve se pakuje u ovu strukturu, ciju adresu vracam kao rezultat
		cli:    client,
		logger: logger,
	}, nil
}

// Disconnect from database
func (pr *BlogRepo) Disconnect(ctx context.Context) error {
	err := pr.cli.Disconnect(ctx)
	if err != nil {
		return err
	}
	return nil
}

// Check database connection
func (pr *BlogRepo) Ping() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Check connection -> if no error, connection is established
	err := pr.cli.Ping(ctx, readpref.Primary())
	if err != nil {
		pr.logger.Println(err)
	}

	// Print available databases
	databases, err := pr.cli.ListDatabaseNames(ctx, bson.M{})
	if err != nil {
		pr.logger.Println(err)
	}
	fmt.Println("databases: ")
	fmt.Println(databases)
}

func (br *BlogRepo) GetAll() (model.Blogs, error) {
	// Initialise context (after 5 seconds timeout, abort operation)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	blogsCollection := br.getCollection()

	var blogs model.Blogs
	blogCursor, err := blogsCollection.Find(ctx, bson.M{})
	if err != nil {
		br.logger.Println(err)
		return nil, err
	}
	if err = blogCursor.All(ctx, &blogs); err != nil {
		br.logger.Println(err)
		return nil, err
	}
	return blogs, nil
}

func (br *BlogRepo) Get(id string) (*model.Blog, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	blogsCollection := br.getCollection()

	var blog model.Blog
	objID, _ := primitive.ObjectIDFromHex(id)
	err := blogsCollection.FindOne(ctx, bson.M{"_id": objID}).Decode(&blog)
	if err != nil {
		br.logger.Println(err)
		return nil, err
	}
	return &blog, nil
}

func (br *BlogRepo) Insert(blog *model.Blog) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	blogsCollection := br.getCollection()

	result, err := blogsCollection.InsertOne(ctx, &blog)
	if err != nil {
		br.logger.Println(err)
		return err
	}
	br.logger.Printf("Documents ID: %v\n", result.InsertedID)
	return nil
}

func (br *BlogRepo) Update(id string, blog *model.Blog) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	blogsCollection := br.getCollection()

	objID, _ := primitive.ObjectIDFromHex(id)
	filter := bson.M{"_id": objID}
	update := bson.M{"$set": bson.M{
		"title":       blog.Title,
		"description": blog.Description,
		"image":       blog.Image,
	}}
	result, err := blogsCollection.UpdateOne(ctx, filter, update)
	br.logger.Printf("Documents matched: %v\n", result.MatchedCount)
	br.logger.Printf("Documents updated: %v\n", result.ModifiedCount)

	if err != nil {
		br.logger.Println(err)
		return err
	}
	return nil
}

func (br *BlogRepo) Delete(id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	blogsCollection := br.getCollection()

	objID, _ := primitive.ObjectIDFromHex(id)
	filter := bson.D{{Key: "_id", Value: objID}}
	result, err := blogsCollection.DeleteOne(ctx, filter)
	if err != nil {
		br.logger.Println(err)
		return err
	}
	br.logger.Printf("Documents deleted: %v\n", result.DeletedCount)
	return nil
}

func (br *BlogRepo) AddVote(id string, vote *model.Vote) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	blogsCollection := br.getCollection()

	objID, _ := primitive.ObjectIDFromHex(id)
	filter := bson.D{{Key: "_id", Value: objID}}
	update := bson.M{"$push": bson.M{
		"votes": vote,
	}}
	result, err := blogsCollection.UpdateOne(ctx, filter, update)
	br.logger.Printf("Documents matched: %v\n", result.MatchedCount)
	br.logger.Printf("Documents updated: %v\n", result.ModifiedCount)

	if err != nil {
		br.logger.Println(err)
		return err
	}
	return nil
}

func (br *BlogRepo) ChangeVote(id string, index int, vote *model.Vote) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	blogsCollection := br.getCollection()

	// What happens if set value for index=10, but we only have 3 phone numbers?
	// -> Every value in between will be set to an empty string
	objID, _ := primitive.ObjectIDFromHex(id)
	filter := bson.D{{Key: "_id", Value: objID}}
	update := bson.M{"$set": bson.M{
		fmt.Sprintf("votes.%d.isUpvote", index): vote.IsUpvote,
	}}
	result, err := blogsCollection.UpdateOne(ctx, filter, update)
	br.logger.Printf("Documents matched: %v\n", result.MatchedCount)
	br.logger.Printf("Documents updated: %v\n", result.ModifiedCount)

	if err != nil {
		br.logger.Println(err)
		return err
	}
	return nil
}

func (br *BlogRepo) GetByAuthorId(userID uint32) (model.Blogs, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	blogsCollection := br.getCollection()

	var blogs model.Blogs
	blogsCursor, err := blogsCollection.Find(ctx, bson.M{"userId": userID})
	if err != nil {
		br.logger.Println(err)
		return nil, err
	}
	if err = blogsCursor.All(ctx, &blogs); err != nil {
		br.logger.Println(err)
		return nil, err
	}
	return blogs, nil
}

func (br *BlogRepo) GetByStatus(status model.Status) (model.Blogs, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	blogsCollection := br.getCollection()

	var blogs model.Blogs
	blogsCursor, err := blogsCollection.Find(ctx, bson.M{"status": status})
	if err != nil {
		br.logger.Println(err)
		return nil, err
	}
	if err = blogsCursor.All(ctx, &blogs); err != nil {
		br.logger.Println(err)
		return nil, err
	}
	return blogs, nil
}

func (br *BlogRepo) AddComment(blogID string, comment *model.Comment) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	comment.ID = primitive.NewObjectID()
	blogsCollection := br.getCollection()

	objID, _ := primitive.ObjectIDFromHex(blogID)

	filter := bson.D{{Key: "_id", Value: objID}}
	update := bson.M{"$push": bson.M{"comments": comment}}

	result, err := blogsCollection.UpdateOne(ctx, filter, update)
	br.logger.Printf("Documents matched: %v\n", result.MatchedCount)
	br.logger.Printf("Documents updated: %v\n", result.ModifiedCount)
	if err != nil {
		br.logger.Println(err)
		return err
	}

	return nil

}

func (br *BlogRepo) UpdateComment(id string, index int, updatedComment *model.Comment) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	blogsCollection := br.getCollection()

	objID, _ := primitive.ObjectIDFromHex(id)
	filter := bson.D{{Key: "_id", Value: objID}}
	br.logger.Print("comment: ", updatedComment.Text)
	update := bson.M{"$set": bson.M{
		fmt.Sprintf("comments.%d.text", index): updatedComment.Text,
	}}
	result, err := blogsCollection.UpdateOne(ctx, filter, update)
	br.logger.Printf("Documents matched: %v\n", result.MatchedCount)
	br.logger.Printf("Documents updated: %v\n", result.ModifiedCount)

	if err != nil {
		br.logger.Println(err)
		return err
	}
	return nil

}

func (br *BlogRepo) DeleteComment(blogID string, commentID string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	blogsCollection := br.getCollection()

	blogObjID, _ := primitive.ObjectIDFromHex(blogID)
	commentObjID, _ := primitive.ObjectIDFromHex(commentID)

	filter := bson.D{{Key: "_id", Value: blogObjID}}
	update := bson.M{
		"$pull": bson.M{"comments": bson.M{"_id": commentObjID}},
	}

	result, err := blogsCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		br.logger.Println(err)
		return err
	}

	br.logger.Printf("Documents matched: %v\n", result.MatchedCount)
	br.logger.Printf("Documents updated: %v\n", result.ModifiedCount)

	return nil

}

func (br *BlogRepo) getCollection() *mongo.Collection {
	blogDatabase := br.cli.Database("mongoDemo")
	blogsCollection := blogDatabase.Collection("blogs")
	return blogsCollection
}
