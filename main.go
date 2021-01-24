package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	var githubType string // Organization, User
	var githubAcc string
	//var githubToken string // 기타허브는 퍼블릭만
	var giteaHost string
	var giteaToken string

	flag.StringVar(&githubType, "t", "Organization", "type of account : Organization/User \n")
	flag.StringVar(&githubAcc, "a", "", "(account name) : organization or user account name\n")
	//flag.StringVar(&githubToken, "ght", "", "GitHub Token\n")
	flag.StringVar(&giteaHost, "gth", "", "(GiTea Host) must provided\n")
	flag.StringVar(&giteaToken, "gtt", "", "(GiTea Token) not found\n")

//	flag.Usage = func() {
//		fmt.Printf(`
//Usage of gitea mirror:
//github-org -a steemit
//`)
//	}
	flag.Parse()

	if githubType == "" {
		githubType = os.Getenv("GITHUB_TYPE")
	}
	if githubAcc == "" {
		githubAcc = os.Getenv("GITHUB_ACC")
	}
	//if githubToken == "" {
	//	githubToken = os.Getenv("GITHUB_TOKEN")
	//}
	if giteaHost == "" {
		giteaHost = os.Getenv("GITEA_HOST")
	}
	if giteaToken == "" {
		giteaToken = os.Getenv("GITEA_TOKEN")
	}

	if githubAcc == "" || giteaHost == "" || giteaToken == "" {
		//fmt.Println("Usage of gitea mirror")
		flag.Usage()
		fmt.Println()
		fmt.Println("ex) giteamir -a github_org_or_user_name -gth gitea.domain.tld gtt tokentokentoken")
		return
	}

	if githubType == "Organization" {
		migrateOrgGithubToGitea(githubAcc, "", giteaHost, giteaToken)
	} else {
		migrateUsrGithubToGitea(githubAcc, "", giteaHost, giteaToken)
	}
}
