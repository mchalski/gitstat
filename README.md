# gitstat
Simple CLI for processing github event stream data.

## Build

```
// go 1.17
go build
```

## Usage

Top 10 repos by watch events:
```
./gitstat top-repos --events data/events.csv --repos data/repos.csv 

```

Top 10 repos by commits:
```
./gitstat top-repos --sort=commits --events data/events.csv --repos data/repos.csv --commits data/commits.csv
```

Top 10 users (UNIMPLEMENTED - SHORT ON TIME, SORRY!):
```
./gitstat top-users --events data/events.csv --commits data/commits.csv --actors data/actors.csv
```

## Comments

- 2 of the 3 commands are implemented
- written for decent testability, things wrapped in interfaces mostly - but tests are entirely TODO
  - including e2e tests (actually start the app with canned args and prepared files)
