package main

import (
	"context"
	"time"

	"github.com/neuronlabs/neuron/codec"
	"github.com/neuronlabs/neuron/database"
	"github.com/neuronlabs/neuron/errors"
	"github.com/neuronlabs/neuron/log"
	"github.com/neuronlabs/neuron/query"
	"github.com/neuronlabs/neuron/server"
)

//go:generate neuron-generator models --format=goimports --single-file .
//go:generate neuron-generator collections --format=goimports  --single-file .

type Blog struct {
	ID        int
	CreatedAt time.Time
	Title     string
	Posts     []*Post
	TopPost   *Post
	TopPostID int
}

var (
	_ database.BeforeInserter = &Blog{}
	_ database.BeforeUpdater  = &Blog{}
)

// BeforeInsert is a model hook executed before insertion of blog.
func (b *Blog) BeforeInsert(ctx context.Context, db database.DB) error {
	ts, ok := ctx.Value("timestamp").(time.Time)
	if !ok {
		log.Errorf("no timestamp found in the context")
	}
	log.Debugf("Blog: %d is being inserted at: %s", ts.String())
	return nil
}

// BeforeUpdate is a hook executed before updating blog.
func (b *Blog) BeforeUpdate(ctx context.Context, db database.DB) error {
	// Get the most likeable post for given blog and set it to given blog top post.
	post, err := NRN_Posts.QueryCtx(ctx, db).
		Select("ID").
		Where("BlogID = ?", b.ID).
		OrderBy("-Likes").
		Get()
	if err != nil {
		cl, ok := err.(errors.ClassError)
		if !ok {
			return err
		}
		// Check if the error is about that no models were found for given query.
		if cl.Class() != query.ClassNoResult {
			return err
		}
	}
	// If the post is found set it's value to the top post.
	if post != nil {
		b.TopPost = post
		b.TopPostID = post.ID
	}
	return nil
}

var (
	_ server.BeforeInsertHandler = &BlogHandler{}
	_ server.AfterInsertHandler  = &BlogHandler{}
)

// BlogHandler is an API handler for the Blog model.
// This struct implements only the hooks before and after insertion of the blogs.
// All other handles would be done by the default json:api handler.
type BlogHandler struct{}

func (b BlogHandler) HandleAfterInsert(params *server.Params, input *codec.Payload) error {
	if input.Meta == nil {
		input.Meta = codec.Meta{}
	}
	// Add some metadata to the input.
	input.Meta["copyright"] = "neuron@neuronlabs.com"

	// Commit the transaction.
	tx, ok := params.DB.(*database.Tx)
	if !ok {
		return errors.NewDetf(server.ClassInternal, "internal error - db should be a transaction")
	}
	if !tx.Transaction.State.Done() {
		if err := tx.Commit(); err != nil {
			return err
		}
	}
	ts, ok := params.Ctx.Value("timestamp").(time.Time)
	if !ok {
		log.Errorf("No timestamp found in the context")
	}
	log.Infof("Inserting blog taken: %s", time.Since(ts))
	return nil
}

func (b BlogHandler) HandleBeforeInsert(params *server.Params, input *codec.Payload) (err error) {
	// Set the transaction at the beginning of the handle with the default options.
	params.DB, err = database.Begin(params.Ctx, params.DB, nil)
	if err != nil {
		return err
	}
	// Set some value in the context that would be used by the handler.
	startTimestamp := params.DB.Controller().Now()
	params.Ctx = context.WithValue(params.Ctx, "timestamp", startTimestamp)
	return nil
}

// Post is the example model that represents blog post.
// It is related to it's root blog, and may contains comments.
type Post struct {
	ID    int
	Title string
	Body  string
	Likes int
	// ForeignKey for the blog.
	Blog     *Blog
	BlogID   int
	Comments []*Comment
}

// Comment is the
type Comment struct {
	ID     int
	Body   string
	Post   *Post
	PostID *int
}
