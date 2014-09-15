package deptools

import (
	"time"
)

type PullRequest2 struct {
	Number   int       `bson:"number"`
	Title    string    `bson:"-"`
	HeadRef  string    `bson:"-"`
	HeadSHA  string    `bson:"header_sha"`
	Status   string    `bson:"-"`
	MergedAt time.Time `bson:"merged_at"`
}

type PrCommentedToTarget2 struct {
	Number  int    `bson:"number"`
	HeadSHA string `bson:"header_sha"`
	Target  string `bson:"target"`
}
