package blog

import (
	"github.com/ezbuy/ezorm/cache"
	"gopkg.in/mgo.v2/bson"
)

type Blog struct {
	ID bson.ObjectId `bson:"_id,omitempty"`

	Title string `bson:"Title"`

	Hits int32 `bson:"Hits"`

	Slug string `bson:"Slug"`

	Body string `bson:"Body"`

	User int32 `bson:"User"`

	IsPublished bool `bson:"IsPublished"`
	isNew       bool
}

func (p *Blog) GetNameSpace() string {
	return "blog"
}

func (p *Blog) GetClassName() string {
	return "Blog"
}

type _BlogMgr struct {
}

var BlogMgr *_BlogMgr

var BlogCache cache.Cache

func (m *_BlogMgr) NewBlog() *Blog {
	rval := new(Blog)
	rval.isNew = true
	rval.ID = bson.NewObjectId()

	return rval
}
