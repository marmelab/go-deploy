package deptools

import (
	"fmt"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"time"
)

type PullRequest struct {
	Number   int       `bson:"number"`
	Title    string    `bson:"-"`
	HeadRef  string    `bson:"-"`
	HeadSHA  string    `bson:"header_sha"`
	Status   string    `bson:"-"`
	MergedAt time.Time `bson:"merged_at"`
}

type PrCommentedToTarget struct {
	Number  int    `bson:"number"`
	HeadSHA string `bson:"header_sha"`
	Target  string `bson:"target"`
}

func (pr *PullRequest) hasBeenDeployTo(target string) bool {
	sess, err := mgo.Dial("localhost")
	if err != nil {
		fmt.Printf("Erreur de connexion a Mongodb : %v", err)
		return false
	}
	defer sess.Close()
	sess.SetSafe(&mgo.Safe{})

	var prDeployed PrCommentedToTarget
	err = sess.DB("deployedPullRequests").
		C("commented_prs").
		Find(bson.M{"number": pr.Number, "header_sha": pr.HeadSHA, "target": target}).
		One(&prDeployed)
	if err != nil {
		return false
	}

	return true
}

func (pr *PullRequest) saveAsCommentToTarget(target string) {
	sess, err := mgo.Dial("localhost")
	if err != nil {
		fmt.Printf("Erreur de connexion a Mongodb : %v", err)
	}
	defer sess.Close()
	sess.SetSafe(&mgo.Safe{})

	collection := sess.DB("deployedPullRequests").C("commented_prs")
	prDeployed := PrCommentedToTarget{
		Number:  pr.Number,
		HeadSHA: pr.HeadSHA,
		Target:  target,
	}
	err = collection.Insert(prDeployed)
	if err != nil {
		fmt.Printf("Erreur a la sauvegarde de la pull request : %v", err)
	}
}
