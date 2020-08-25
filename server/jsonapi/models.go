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

//go:generate neurogns models methods --format=goimports --single-file .
//go:generate neurogns collections --format=goimports  --single-file .

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
		// Check if the error is about that no models were found for given query.
		if !errors.Is(err, query.ErrNoResult) {
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
	_ server.WithContextInserter = &BlogHandler{}
	_ server.InsertTransactioner = &BlogHandler{}
	_ server.AfterInsertHandler  = &BlogHandler{}
)

// BlogHandler is an API handler for the Blog model.
// This struct implements only the hooks before and after insertion of the blogs.
// All other handles would be done by the default json:api handler.
type BlogHandler struct{}

// InsertWithContext implements server.WithContextInserter.
func (b BlogHandler) InsertWithContext(ctx context.Context) (context.Context, error) {
	ts := time.Now()
	return context.WithValue(ctx, "timestamp", ts), nil
}

// InsertWithTransaction implements server.InsertTransactioner.
func (b BlogHandler) InsertWithTransaction() *query.TxOptions {
	return nil
}

func (b BlogHandler) HandleAfterInsert(ctx context.Context, db database.DB, input *codec.Payload) error {
	if input.Meta == nil {
		input.Meta = codec.Meta{}
	}
	// Add some metadata to the input.
	input.Meta["copyright"] = "neuron@neuronlabs.com"
	ts, ok := ctx.Value("timestamp").(time.Time)
	if !ok {
		log.Errorf("No timestamp found in the context")
	}
	log.Infof("Inserting blog taken: %s", time.Since(ts))
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
