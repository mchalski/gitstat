package actions

import (
	"sort"

	"github.com/mchalski/gitstat/pkg/data"
)

const (
	topDefault = 10

	// tuple field constants
	evtTypeWatch = "WatchEvent"
	evtTypePush  = "PushEvent"

	// tuple field indices
	evtIdxId     = 0
	evtIdxType   = 1
	evtIdxRepoId = 3

	repoIdxId   = 0
	repoIdxName = 1

	commitIdxEvtId = 2
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
	// map repo ids to repo results
	cntMap := make(map[string]RepoRes)

	// count watch events in the 'events' stream
	evtChan := a.evts.C()
	for e := range evtChan {
		if e.Err != nil {
			return nil, e.Err
		}

		evtType, repoId := e.Tuple[evtIdxType], e.Tuple[evtIdxRepoId]

		if evtType == evtTypeWatch {
			repo, ok := cntMap[repoId]
			if !ok {
				repo = RepoRes{
					Id:    repoId,
					Count: 0,
				}
			}

			repo.Count += 1
			cntMap[repoId] = repo
		}
	}

	// pull repo results into a slice, sort, get top
	res := make([]RepoRes, len(cntMap))

	i := 0
	for _, v := range cntMap {
		res[i] = v
		i += 1
	}

	sort.Slice(res, func(i, j int) bool {
		return res[i].Count > res[j].Count
	})

	topIdx := topDefault
	if topIdx > len(res) {
		topIdx = len(res)
	}

	res = res[:topIdx]

	// find repo names in the 'repos' stream
	repoChan := a.repos.C()
	for e := range repoChan {
		if e.Err != nil {
			return nil, e.Err
		}

		repoId, repoName := e.Tuple[repoIdxId], e.Tuple[repoIdxName]

		// finish early if we have all top names
		if fillName(repoId, repoName, res) {
			break
		}
	}

	return res, nil
}

// fillName tries to match/insert a name into a set of repo results described by ids so far
// returns a flag telling if all repos have names already set
func fillName(id, name string, repos []RepoRes) bool {
	allFilled := true

	for i, r := range repos {
		if r.Name == "" {
			allFilled = false
		}

		if r.Id == id {
			repos[i].Name = name
		}
	}

	return allFilled
}

func (a *TopReposByCommits) Run() ([]RepoRes, error) {
	// map event ids to repo results
	// (counts correspond to # commits per event initially)
	evtCntMap := make(map[string]RepoRes)

	// record push events
	evtChan := a.evts.C()
	for e := range evtChan {
		if e.Err != nil {
			return nil, e.Err
		}

		evtId, evtType, repoId := e.Tuple[evtIdxId],
			e.Tuple[evtIdxType],
			e.Tuple[evtIdxRepoId]

		if evtType == evtTypePush {
			_, ok := evtCntMap[repoId]

			if !ok {
				repo := RepoRes{
					Id:    repoId,
					Count: 0,
				}
				evtCntMap[evtId] = repo
			}
		}
	}

	// count commits for each event
	commitChan := a.commits.C()
	for e := range commitChan {
		if e.Err != nil {
			return nil, e.Err
		}

		evtId := e.Tuple[commitIdxEvtId]

		// if we're interested in this event - it should already
		// be tracked, bump the commit counter if so
		res, ok := evtCntMap[evtId]
		if !ok {
			continue
		}

		res.Count += 1
		evtCntMap[evtId] = res
	}

	// at this point - we have a slice of partial repo counts
	// need to sum them up; use a count map again

	// map repo ids to repo results
	cntMap := make(map[string]RepoRes)

	for _, r := range evtCntMap {
		repo, ok := cntMap[r.Id]
		if !ok {
			repo = RepoRes{
				Id:    r.Id,
				Count: 0,
			}
		}

		// bump the final counter by the partial count
		repo.Count += r.Count
		cntMap[repo.Id] = repo

	}

	// pull repo results into a slice, sort, get top
	res := make([]RepoRes, len(cntMap))

	i := 0
	for _, v := range cntMap {
		res[i] = v
		i += 1
	}

	sort.Slice(res, func(i, j int) bool {
		return res[i].Count > res[j].Count
	})

	topIdx := topDefault
	if topIdx > len(res) {
		topIdx = len(res)
	}

	res = res[:topIdx]

	// find repo names in the 'repos' stream
	repoChan := a.repos.C()
	for e := range repoChan {
		if e.Err != nil {
			return nil, e.Err
		}

		repoId, repoName := e.Tuple[repoIdxId], e.Tuple[repoIdxName]

		// finish early if we have all top names
		if fillName(repoId, repoName, res) {
			break
		}
	}

	return res, nil
}
