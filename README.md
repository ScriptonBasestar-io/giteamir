# giteamir

gitea mirror command line

## Install

```shell
go install github.com/ScriptonBasestar-io/giteamir
```

## Caution

github api limit

## Usage

```shell
# default
expoet GITHUB_TYPE=Organization
export GITHUB_ACCOUNT
# 필요없음. 퍼블릭만 지원하니까
export GITHUB_TOKEN
export GITEA_HOST
export GITEA_TOKEN

gitmir 
```

```shell
gitmir -t User -a personaiam -ght gitea.domain.tld -gtt giteatokentokentoken
```