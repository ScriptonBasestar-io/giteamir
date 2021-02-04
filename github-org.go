package main

import (
	"code.gitea.io/sdk/gitea"
	"context"
	"fmt"
	"github.com/google/go-github/github"
	"net/http"
	"os"
	"time"
)

func migrateOrgGithubToGitea(githubAccName, githubToken, giteaHost, giteaToken string) {
	ctx := context.Background()
	//ts := oauth2.StaticTokenSource(
	//	&oauth2.Token{AccessToken: githubToken},
	//)
	//tc := oauth2.NewClient(ctx, ts)
	//
	//githubClient := github.NewClient(tc)

	githubClient := github.NewClient(&http.Client{})
	orgOption := gitea.CreateOrgOption{}

	githubOrgObj, _, err := githubClient.Organizations.Get(ctx, githubAccName)
	if err != nil {
		fmt.Println("Org not exists")
		fmt.Println(err)
	} else {
		fmt.Print(githubOrgObj)
		description := ""
		if githubOrgObj.Description != nil { // will throw a nil pointer error if description is passed directly to the below struct
			description = *githubOrgObj.Description
		} else {
			description = ""
		}
		orgOption = gitea.CreateOrgOption{
			Name: *githubOrgObj.Login,
			//FullName: *githubOrgObj.
			Description: description,
			Website:     *githubOrgObj.HTMLURL,
			//Location:    *githubOrgObj.Location,
			Visibility: gitea.VisibleTypeLimited,
		}
	}
	// get all repositories from organization
	var allRepos []*github.Repository
	opt := &github.RepositoryListByOrgOptions{
		ListOptions: github.ListOptions{PerPage: 100},
	}
	for {
		repos, resp, err := githubClient.Repositories.ListByOrg(ctx, *githubOrgObj.Login, opt)
		//fmt.Println(repos)
		if err != nil {
			fmt.Println(err)
			//os.Exit(1)
			return
		}
		allRepos = append(allRepos, repos...)
		if resp.NextPage == 0 {
			break
		}
		opt.Page = resp.NextPage
	}
	fmt.Println("read repo success")

	// avatar 다운로드 수정 할까했는데 gitea에 변경api가 없음
	//err := downloadFile(*githubOrgObj.AvatarURL, "tmpavatar")
	//if err != nil {
	//	log.Fatal(err)
	//}

	giteaClient, _ := gitea.NewClient("https://"+giteaHost+"/", gitea.SetToken(giteaToken))

	// create org if not exists
	giteaOrgObj, res, err := giteaClient.GetOrg(*githubOrgObj.Login)
	fmt.Println(giteaOrgObj, res, err)
	if err != nil && err.Error() == "404 Not Found" {
		fmt.Println("organization not exists in gitea")
		fmt.Println("create org : " + *githubOrgObj.Login)
		giteaOrgObj, res, err = giteaClient.CreateOrg(orgOption)
		if err != nil {
			fmt.Println("exit. org id is ", giteaOrgObj.ID)
			os.Exit(1)
			return
		}
	}

	//if giteaOrgObj.ID == 1 {
	//	fmt.Println("exit. org id is ", giteaOrgObj.ID)
	//	os.Exit(1)
	//	return
	//}

	for i := 0; i < len(allRepos); i++ {
		fmt.Printf("repo name %d/%d  id: %d  %s\n", i, len(allRepos), giteaOrgObj.ID, *allRepos[i].Name)
		description := ""
		if allRepos[i].Description != nil { // will throw a nil pointer error if description is passed directly to the below struct
			description = *allRepos[i].Description
		}
		//giteaClient.TransferRepo("archmagece", *allRepos[i].Name, gitea.TransferRepoOption{
		//	NewOwner: giteaOrgObj.UserName,
		//})
		//res, err := giteaClient.DeleteRepo("archmagece", *allRepos[i].Name)
		//if err != nil {
		//	fmt.Println(res)
		//	fmt.Println("errorr")
		//}
		repo, _, _ := giteaClient.MigrateRepo(gitea.MigrateRepoOption{
			CloneAddr:   *allRepos[i].CloneURL,
			RepoOwnerID: giteaOrgObj.ID,
			RepoName:    *allRepos[i].Name,
			Mirror:      true,
			Private:     false,
			Description: description,
		})
		fmt.Println(repo)
		time.Sleep(100 * time.Millisecond) // THIS IS HERE SO THE GITEA SERVER DOESNT GET HAMMERED WITH REQUESTS
	}
}
