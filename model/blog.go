package model

import (
	"encoding/json"
	"io"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Status uint32

const (
	Draft     Status = 0
	Published Status = 1
	Closed    Status = 2
	Active    Status = 3
	Famous    Status = 4
)

type Category uint32

const (
	Destinations  Category = 0
	Travelogues   Category = 1
	Activities    Category = 2
	Gastronomy    Category = 3
	Tips          Category = 4
	Culture       Category = 5
	Accommodation Category = 6
)

type Blog struct {
	ID           primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID       uint32             `bson:"userId,omitempty" json:"userId"`
	Username     string             `bson:"username, omitempty" json:"username"`
	Title        string             `bson:"title" json:"title"`
	Description  string             `bson:"description,omitempty" json:"description"`
	CreationTime time.Time          `bson:"creationTime,omitempty" json:"creationTime"`
	Status       Status             `bson:"status,omitempty" json:"status"`
	Image        string             `bson:"image,omitempty" json:"image"`
	Category     Category           `bson:"category,omitempty" json:"category"`
	Comments     []Comment          `bson:"comments,omitempty" json:"comments"`
	Votes        []Vote             `bson:"votes,omitempty" json:"votes"`
}

type Blogs []*Blog

func (p *Blogs) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(p)
}

func (p *Blog) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(p)
}

func (p *Blog) FromJSON(r io.Reader) error {
	d := json.NewDecoder(r)
	return d.Decode(p)
}
