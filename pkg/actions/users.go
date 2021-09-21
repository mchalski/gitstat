package actions

import (
	"github.com/mchalski/gitstat/pkg/data"
)

// TopUsers represents listing the top users by the number of created PRs and pushed commits
type TopUsers struct {
	evts    data.Stream
	commits data.Stream
	actors  data.Stream
}

func NewTopUsers(evts data.Stream, commits data.Stream, actors data.Stream) *TopUsers {
	return &TopUsers{
		evts:    evts,
		commits: commits,
		actors:  actors,
	}
}

func (a *TopUsers) Run() {
	//TODO
}
