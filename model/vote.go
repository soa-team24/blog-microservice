package model

import (
	"encoding/json"
	"io"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Vote struct {
	ID           primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	IsUpvote     bool               `bson:"isUpvote, omitempty" json:"isUpvote"`
	UserID       uint32             `bson:"userId,omitempty" json:"userId"`
	CreationTime time.Time          `bson:"creationTime,omitempty" json:"creationTime"`
}

type Votes []*Vote

func (v *Votes) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(v)
}

func (v *Vote) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(v)
}

func (v *Vote) FromJSON(r io.Reader) error {
	d := json.NewDecoder(r)
	return d.Decode(v)
}
