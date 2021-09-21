package actions

import (
	"github.com/mchalski/gitstat/pkg/data"
)

const (
	topDefault = 10
)

// TopReposByWatchEvts represents listing the top repos by the number of watch events
type TopReposByWatchEvts struct {
	evts  data.Stream
	repos data.Stream
}

// TopReposByCommits represents listing the top repos by number of commits
type TopReposByCommits struct {
	evts    data.Stream
	commits data.Stream
	repos   data.Stream
}

// RepoRes is the 'top repos' results item (repo details + some counter)
type RepoRes struct {
	Id   string
	Name string

	// either of watch evts, or commits pushed
	Count int
}

func NewTopReposByWatchEvts(evts data.Stream, repos data.Stream) *TopReposByWatchEvts {
	return &TopReposByWatchEvts{
		evts:  evts,
		repos: repos,
	}
}

func NewTopReposByCommits(evts data.Stream, commits data.Stream, repos data.Stream) *TopReposByCommits {
	return &TopReposByCommits{
		evts:    evts,
		commits: commits,
		repos:   repos,
	}
}

func (a *TopReposByWatchEvts) Run() ([]RepoRes, error) {
	//TODO
	return nil, nil
}

func (a *TopReposByCommits) Run() ([]RepoRes, error) {
	//TODO
	return nil, nil
}
