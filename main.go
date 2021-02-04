package main

import (
	"flag"
	"fmt"
	"os"
)

//type GithubType int
//
//const (
//    User GithubType = 1 + iota
//    Organization
//)
//
//var githubType = [...]string{
//    "user",
//    "organization",
//}

//func (g GithubType) Create(gt string) GithubType { return githubType[gt]  }
//func (g GithubType) String() string { return githubType[(g-1)%2] }

func main() {
	var githubAccountType string // Organization, User
	var githubAccount string
	//var githubToken string // 기타허브는 퍼블릭만
	var giteaHost string
	var giteaToken string

	flag.StringVar(&githubAccountType, "t", "", "type of account : Organization/User \n")
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

	if githubAccountType == "" {
		githubAccountType = os.Getenv("GITHUB_ACCOUNT_TYPE")
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

	if githubAccount == "" || giteaHost == "" || giteaToken == "" || githubAccountType == "" {
		//fmt.Println("Usage of gitea mirror")
		flag.Usage()
		fmt.Println()

		if githubAccount == "" {
			fmt.Println("github account is not specified. ex) -a scriptonbasestar ,  env not supported")
		}
		if giteaHost == "" {
			fmt.Println("gitea host is not specified. ex) -gth gitea.domain.tld or export GITEA_HOST=gitea.domain.tld")
		}
		if giteaToken == "" {
			fmt.Println("gitea token is not specified. ex) -gtt 'tokentokentoken' or export GITEA_TOKEN=tokentokentoken")
		}
		if githubAccountType == "" {
			fmt.Println("gitea account type is not specified. ex) -t Organization or export GITHUB_ACCOUNT_TYPE=Organization")
		}

		fmt.Println("ex) giteamir -a github_org_or_user_name -gth gitea.domain.tld -gtt tokentokentoken")
		os.Exit(1)
	}

	if githubAccountType == "Organization" {
		migrateOrgGithubToGitea(githubAccount, "", giteaHost, giteaToken)
	} else if githubAccountType == "User" {
		migrateUsrGithubToGitea(githubAccount, "", giteaHost, giteaToken)
	} else {
		fmt.Printf("githubType %s is not supported", githubAccountType)
	}
}
