package deptools

import (
	"fmt"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"time"
)

func (dpl *Deployment) save() {
	sess, err := mgo.Dial("localhost")
	if err != nil {
		fmt.Printf("Erreur de connexion a Mongodb : %v", err)
	}
	defer sess.Close()
	sess.SetSafe(&mgo.Safe{})

	collection := sess.DB("deployedPullRequests").C("deployments")
	err = collection.Insert(dpl)
	if err != nil {
		fmt.Printf("Erreur a la sauvegarde du deploiement : %v", err)
	}
}

func (dpl *Deployment) getLastPrMergeDate() time.Time {
	sess, err := mgo.Dial("localhost")
	if err != nil {
		fmt.Printf("Erreur de connexion a Mongodb : %v", err)
	}
	defer sess.Close()
	sess.SetSafe(&mgo.Safe{})

	var saveDeploy Deployment
	err = sess.DB("deployedPullRequests").
		C("deployments").
		Find(bson.M{"Owner": dpl.Owner,
		"Repository": dpl.Repository,
		"base_type":  dpl.Base_type,
		"base_name":  dpl.Base_name,
		"Target":     dpl.Target}).
		Sort("-last_pr_merge_date").
		One(&saveDeploy)
	if err != nil {

		return time.Now().Add(-2 * 7 * 24 * time.Hour)
	}

	return saveDeploy.Last_pr_merge_date
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
