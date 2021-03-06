package main

import (
	"code.gitea.io/sdk/gitea"
	"context"
	"fmt"
	"github.com/google/go-github/github"
	"github.com/zenthangplus/goccm"
	"net/http"
	"os"
)

func migrateUsrGithubToGitea(githubAccName, githubToken, giteaHost, giteaToken string) {
	ctx := context.Background()
	//ts := oauth2.StaticTokenSource(
	//	&oauth2.Token{AccessToken: githubToken},
	//)
	//tc := oauth2.NewClient(ctx, ts)
	//
	//githubClient := github.NewClient(tc)

	githubClient := github.NewClient(&http.Client{})
	orgOption := gitea.CreateOrgOption{}

	githubUsrObj, _, err := githubClient.Users.Get(ctx, githubAccName)
	if err != nil {
		fmt.Println("Org not exists")
		fmt.Println(err)
	} else {
		fmt.Println(githubUsrObj)
		orgOption = gitea.CreateOrgOption{
			Name: *githubUsrObj.Login,
			//FullName: *githubUsrObj.
			//Description: ,
			Website: *githubUsrObj.ReposURL,
			//Location: *githubUsrObj.Location,
			Visibility: gitea.VisibleTypeLimited,
		}
	}
	// get all repositories from organization
	var allRepos []*github.Repository
	opt := &github.RepositoryListOptions{
		ListOptions: github.ListOptions{PerPage: 100},
	}
	for {
		repos, resp, err := githubClient.Repositories.List(ctx, *githubUsrObj.Login, opt)
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
	giteaOrgObj, res, err := giteaClient.GetOrg(*githubUsrObj.Login)
	fmt.Println(giteaOrgObj, res, err)
	if err != nil && err.Error() == "404 Not Found" {
		fmt.Println("organization not exists in gitea")
		fmt.Println("create org : " + *githubUsrObj.Login)
		giteaOrgObj, res, err = giteaClient.CreateOrg(orgOption)
		if err != nil {
			fmt.Println("exit. org id is ", giteaOrgObj.ID)
			os.Exit(1)
			return
		}
	}

	c := goccm.New(10)

	// https://api.github.com/orgs/wikibook !?!!?!? https://api.github.com/users/wikibook ?!?!?!?!
	for i := 0; i < len(allRepos); i++ {
		c.Wait()
		fmt.Printf("repo name %d/%d  id: %d  %s\n", i, len(allRepos), giteaOrgObj.ID, *allRepos[i].Name)
		description := ""
		if allRepos[i].Description != nil { // will throw a nil pointer error if description is passed directly to the below struct
			description = *allRepos[i].Description
		}
		//res, err := giteaClient.DeleteRepo("archmagece", *allRepos[i].Name)
		//if err != nil {
		//	fmt.Println(res)
		//	fmt.Println("errorr")
		//}
		go func(i int, description string) {
			repo, _, _ := giteaClient.MigrateRepo(gitea.MigrateRepoOption{
				CloneAddr:   *allRepos[i].CloneURL,
				RepoOwnerID: giteaOrgObj.ID,
				RepoName:    *allRepos[i].Name,
				Mirror:      true,
				Private:     false,
				Description: description,
			})
			fmt.Println("finish", repo.Name, repo.CloneURL)
			c.Done()
		}(i, description)
	}
	c.WaitAllDone()
}
