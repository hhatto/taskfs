package gitlab

import (
	"fmt"
	"time"

	"github.com/lufia/taskfs/fs"
	"github.com/xanzy/go-gitlab"
)

type Comment struct {
	num  int
	note *gitlab.Note
}

func (p *Comment) Key() string {
	return fmt.Sprintf("%d", p.num)
}

func (p *Comment) Message() string {
	return p.note.Body
}

func (p *Comment) Creation() time.Time {
	return *p.note.CreatedAt
}

func (p *Comment) LastMod() time.Time {
	return *p.note.UpdatedAt
}

type Issue struct {
	issue *gitlab.Issue
	svc   *Service
}

func (p *Issue) Key() string {
	// TODO: implement
	//group
	//project
	return fmt.Sprintf("%d", p.issue.ID)
}

func (p *Issue) Subject() string {
	return p.issue.Title
}

func (p *Issue) Message() string {
	return p.issue.Description
}

func (p *Issue) PermaLink() string {
	return p.issue.WebURL
}

func (p *Issue) Creation() time.Time {
	return *p.issue.CreatedAt
}

func (p *Issue) LastMod() time.Time {
	return *p.issue.UpdatedAt
}

func (p *Issue) Comments() (a []fs.Comment, err error) {
	var buf []*gitlab.Note
	page := 0
	for {
		var b []*gitlab.Note
		b, page, err = p.fetchNotes(page)
		if err != nil {
			return
		}
		buf = append(buf, b...)
		if page == 0 {
			break
		}
	}
	a = make([]fs.Comment, len(buf))
	for i, v := range buf {
		a[i] = &Comment{num: i + 1, note: v}
	}
	return a, nil
}

func (p *Issue) fetchNotes(page int) ([]*gitlab.Note, int, error) {
	pid := p.issue.ProjectID
	n := p.issue.ID
	var opt gitlab.ListIssueNotesOptions
	opt.Page = page
	b, resp, err := p.svc.c.Notes.ListIssueNotes(pid, n, &opt)
	if err != nil {
		return nil, 0, err
	}
	return b, resp.NextPage, nil
}

type Config struct {
	BaseURL string
	Token   string
}

type Service struct {
	c *gitlab.Client
}

func NewService(config *Config) (*Service, error) {
	c := gitlab.NewClient(nil, config.Token)
	c.SetBaseURL(config.BaseURL)
	return &Service{c: c}, nil
}

func (p *Service) Name() string {
	return "gitlab"
}

func (p *Service) List() ([]fs.Task, error) {
	var a []fs.Task
	var opt gitlab.ListIssuesOptions
	for {
		b, resp, err := p.c.Issues.ListIssues(&opt)
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

func (p *Service) appendIssues(a []fs.Task, b []*gitlab.Issue) []fs.Task {
	for _, v := range b {
		a = append(a, &Issue{issue: v, svc: p})
	}
	return a
}
