package model

import (
	"encoding/json"
	"io"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Comment struct {
	ID               primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID           uint32             `bson:"userID,omitempty" json:"userID"`
	Username         string             `bson:"username,omitempty" json:"username"`
	Text             string             `bson:"text,omitempty" json:"text"`
	CreationTime     time.Time          `bson:"creationTime,omitempty" json:"creationTime"`
	LastModification time.Time          `bson:"lastModification,omitempty" json:"lastModification"`
}

type Comments []*Comment

func (p *Comment) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(p)
}

func (p *Comment) FromJSON(r io.Reader) error {
	d := json.NewDecoder(r)
	return d.Decode(p)
}

func (v *Comments) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(v)
}
