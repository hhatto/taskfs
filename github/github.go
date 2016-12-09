package github

import (
	"net/http"
	"net/url"
	"time"

	"github.com/google/go-github/github"
	"github.com/lufia/ac2016/fs"
	"golang.org/x/oauth2"
)

type Issue github.Issue

func (p *Issue) ArticleID() uint32 {
	return uint32(*p.ID)
}

func (p *Issue) Subject() string {
	return *p.Title
}

func (p *Issue) Message() string {
	return *p.Body
}

func (p *Issue) Creation() time.Time {
	return *p.CreatedAt
}

func (p *Issue) LastMod() time.Time {
	return *p.UpdatedAt
}

type Config struct {
	BaseURL string
	Token   string
}

func (c *Config) authorizedClient() *http.Client {
	if c.Token == "" {
		return nil
	}
	token := &oauth2.Token{
		AccessToken: c.Token,
	}
	s := oauth2.StaticTokenSource(token)
	return oauth2.NewClient(oauth2.NoContext, s)
}

type Service struct {
	c *github.Client
}

func NewService(config *Config) (*Service, error) {
	var client *http.Client
	if config.Token != "" {
		client = config.authorizedClient()
	}
	c := github.NewClient(client)
	if config.BaseURL != "" {
		u, err := url.Parse(config.BaseURL)
		if err != nil {
			return nil, err
		}
		c.BaseURL = u
	}
	return &Service{c: c}, nil
}

func (p *Service) List() ([]fs.Article, error) {
	var a []fs.Article
	var opt github.IssueListOptions
	for {
		b, resp, err := p.c.Issues.List(true, &opt)
		if err != nil {
			return nil, err
		}
		a = p.appendIssues(a, b)
		if resp.NextPage == 0 {
			break
		}
		opt.Page = resp.NextPage
	}
	return a, nil
}

func (*Service) appendIssues(a []fs.Article, b []*github.Issue) []fs.Article {
	for _, v := range b {
		a = append(a, (*Issue)(v))
	}
	return a
}