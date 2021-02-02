package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	var githubType string // Organization, User
	var githubAccount string
	//var githubToken string // 기타허브는 퍼블릭만
	var giteaHost string
	var giteaToken string

	flag.StringVar(&githubType, "t", "", "type of account : Organization/User \n")
	flag.StringVar(&githubAccount, "a", "", "(account name) : organization or user account name\n")
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
	//if githubAccount == "" {
	//	githubAccount = os.Getenv("GITHUB_ACCOUNT")
	//}
	//if githubToken == "" {
	//	githubToken = os.Getenv("GITHUB_TOKEN")
	//}
	if giteaHost == "" {
		giteaHost = os.Getenv("GITEA_HOST")
	}
	if giteaToken == "" {
		giteaToken = os.Getenv("GITEA_TOKEN")
	}

	if githubAccount == "" || giteaHost == "" || giteaToken == "" {
		//fmt.Println("Usage of gitea mirror")
		flag.Usage()
		fmt.Println()
		fmt.Println("ex) giteamir -a github_org_or_user_name -gth gitea.domain.tld -gtt tokentokentoken")
		os.Exit(1)
		//return
	}

	if githubType == "Organization" {
		migrateOrgGithubToGitea(githubAccount, "", giteaHost, giteaToken)
	} else if githubType == "User" {
		migrateUsrGithubToGitea(githubAccount, "", giteaHost, giteaToken)
	} else {
		fmt.Printf("githubType %s is not supported", githubType)
	}
}
