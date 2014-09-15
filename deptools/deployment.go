package deptools

import (
	"code.google.com/p/goauth2/oauth"
	"fmt"
	"github.com/google/go-github/github"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"time"
)

type Deployment struct {
	Owner                   string                 `bson:"owner"`
	Repository              string                 `bson:"repository"`
	AccessToken             string                 `bson:"-"`
	BaseType                string                 `bson:"base_type"`
	BaseName                string                 `bson:"base_name"`
	BaseTagSHA              string                 `bson:"-"`
	Target                  string                 `bson:"target"`
	PullRequests            map[string]PullRequest `bson:"-"`
	CommitsOnDeployedBase   map[string]string      `bson:"-"`
	PullRequestMergedOnBase map[string]PullRequest `bson:"-"`
	LastPrMergeDate         time.Time              `bson:"last_pr_merge_date"`
	CreatedAt               time.Time              `bson:"created_at"`
}

func (dpl *Deployment) BaseExist() (bool, error) {
	if dpl.BaseType == "tag" {
		exist, err := dpl.tagExist()
		return exist, err
	} else {
		exist, err := dpl.branchExist()
		return exist, err
	}
}

//TODO refactoring these two identical functions by regulating the type signature issues
func (dpl *Deployment) tagExist() (bool, error) {
	client := dpl.getGithubAccessClient()
	opt := &github.ListOptions{}
	tags, _, err := client.Repositories.ListTags(dpl.Owner, dpl.Repository, opt)
	if err != nil {
		return false, err
	} else {
		for _, tag := range tags {
			if *tag.Name == dpl.BaseName {
				dpl.BaseTagSHA = *tag.Commit.SHA
				return true, nil
			}
		}
	}
	tagNotFound := fmt.Errorf("%v %v not found on %v", dpl.BaseType, dpl.BaseName, dpl.Repository)
	return false, tagNotFound
}
func (dpl *Deployment) branchExist() (bool, error) {
	client := dpl.getGithubAccessClient()
	opt := &github.ListOptions{}
	branches, _, err := client.Repositories.ListBranches(dpl.Owner, dpl.Repository, opt)
	if err != nil {
		return false, err
	} else {
		for _, branch := range branches {
			if *branch.Name == dpl.BaseName {
				return true, nil
			}
		}
	}
	branchNotFound := fmt.Errorf("%v %v not found on %v", dpl.BaseType, dpl.BaseName, dpl.Repository)
	return false, branchNotFound
}

func (dpl *Deployment) CommentPrContainedInDeploy() (string, error) {
	if getPRError := dpl.getClosedPullRequests(); getPRError != nil {
		fmt.Println("bim une erreur getClosedPullRequests")
		return "", getPRError
	}

	if getCommitError := dpl.getCommitsOnBase(); getCommitError != nil {
		return "", getCommitError
	}

	dpl.setPrMergedOnBase()
	if nbPrToComment := len(dpl.PullRequestMergedOnBase); nbPrToComment < 1 {
		return "", nil
	}
	if commentError := dpl.commentMergedPR(); commentError != nil {
		return "", commentError
	}

	return "DONE", nil
}

func (dpl *Deployment) getClosedPullRequests() error {

	//TODO limit the numebr of result, with Base Option ?
	//TODO refactoring because if else -> for -> if ....
	client := dpl.getGithubAccessClient()
	opt := &github.PullRequestListOptions{State: "closed"}
	prs, _, err := client.PullRequests.List(dpl.Owner, dpl.Repository, opt)
	if err != nil {
		return err
	} else {
		dpl.PullRequests = make(map[string]PullRequest)
		for _, pr := range prs {
			if pr.MergedAt != nil {
				pullr := PullRequest{
					Number:   *pr.Number,
					Title:    *pr.Title,
					HeadRef:  *pr.Head.Ref,
					HeadSHA:  *pr.Head.SHA,
					Status:   *pr.State,
					MergedAt: *pr.MergedAt}
				dpl.PullRequests[pullr.HeadSHA] = pullr
			}
		}
	}
	return nil
}

func (dpl *Deployment) getCommitsOnBase() error {
	lastPrMergedDate := dpl.getLastPrMergeDate()
	var searchBaseName string
	if dpl.BaseType == "tag" {
		searchBaseName = string(dpl.BaseTagSHA)
	} else {
		searchBaseName = string(dpl.BaseName)
	}
	client := dpl.getGithubAccessClient()
	opt := &github.CommitsListOptions{SHA: searchBaseName, Since: lastPrMergedDate}
	commits, _, err := client.Repositories.ListCommits(dpl.Owner, dpl.Repository, opt)
	if err != nil {
		return err
	} else {
		dpl.CommitsOnDeployedBase = make(map[string]string)
		for _, commit := range commits {
			dpl.CommitsOnDeployedBase[*commit.SHA] = *commit.SHA
		}
	}
	return nil
}

func (dpl *Deployment) setPrMergedOnBase() {
	dpl.PullRequestMergedOnBase = make(map[string]PullRequest)
	for commitSha, _ := range dpl.CommitsOnDeployedBase {
		if mergedPullRequest, found := dpl.PullRequests[commitSha]; found {
			dpl.PullRequestMergedOnBase[mergedPullRequest.HeadSHA] = mergedPullRequest
			if mergedPullRequest.MergedAt.After(dpl.LastPrMergeDate) {
				dpl.LastPrMergeDate = mergedPullRequest.MergedAt
			}
		}
	}
}

func (dpl *Deployment) commentMergedPR() error {
	client := dpl.getGithubAccessClient()
	msg := fmt.Sprintf("This PR was deployed to %v (from the %v %v)", dpl.Target, dpl.BaseType, dpl.BaseName)
	comment := &github.IssueComment{Body: &msg}
	for _, prToComment := range dpl.PullRequestMergedOnBase {
		if hasAlreadyBeenCommented := prToComment.hasBeenDeployTo(dpl.Target); !hasAlreadyBeenCommented {
			_, _, err := client.Issues.CreateComment(dpl.Owner, dpl.Repository, prToComment.Number, comment)
			if err != nil {
				return err
			}
			prToComment.saveAsCommentToTarget(dpl.Target)
		}
	}
	dpl.save()
	return nil
}

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
		Find(bson.M{"owner": dpl.Owner, "repository": dpl.Repository, "base_type": dpl.BaseType, "base_name": dpl.BaseName, "target": dpl.Target}).
		Sort("-last_pr_merge_date").
		One(&saveDeploy)
	if err != nil {
		return time.Now().Add(-2 * 7 * 24 * time.Hour)
	}

	return saveDeploy.LastPrMergeDate
}

func (dpl *Deployment) getGithubAccessClient() *github.Client {
	t := &oauth.Transport{
		Token: &oauth.Token{AccessToken: dpl.AccessToken},
	}
	client := github.NewClient(t.Client())
	return client
}
