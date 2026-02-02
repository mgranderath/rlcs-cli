# rlcs-cli

Command-line interface for exploring Rocket League Championship Series (RLCS) tournaments, matches, and brackets.

**Installation**

Primary install (recommended):

```bash
go install github.com/mgranderath/rlcs-cli/cmd/rlcs-cli@latest
```

Build from source (optional):

```bash
git clone https://github.com/mgranderath/rlcs-cli.git
cd rlcs-cli
go build -o rlcs-cli ./cmd/rlcs-cli
```

**Quickstart**

```bash
rlcs-cli tournaments list
```

**Command Structure**

```text
rlcs-cli [--debug] [--version|-v] <command>

commands:
  tournaments
    list
    matches
    brackets <tournamentID>
  matches
    list <tournamentID>
    get <matchID>
```

**Command Reference**

Top-level flags:
- `--debug` Enable debug mode.
- `--version`, `-v` Show version and exit.

`tournaments list` — List tournaments in a circuit/year.
- `--circuit` Circuit/year (e.g., `2025`, `2026`). Defaults to current year.
- `--region` Region filter: `NA`, `EU`, `APAC`, `SAM`, `OCE`, `MENA`, `SSA`
- `--online` Show only online tournaments.
- `--major` Show only majors (empty region/grouping).
- `--grouping` Filter by tournament grouping (partial name match).
- `--upcoming` Start date > today.
- `--ongoing` Start date <= today <= end date.
- `--past` End date < today.
- `--min-teams` Minimum number of teams.
- `--output`, `-o` Output format: `table`, `json`, `csv`, `yaml`.

`tournaments matches` — List matches across tournaments.
- `--circuit` Circuit/year (e.g., `2025`, `2026`). Defaults to current year.
- `--region` Region filter: `NA`, `EU`, `APAC`, `SAM`, `OCE`, `MENA`, `SSA`
- `--online` Show only online tournaments.
- `--major` Show only majors (empty region/grouping).
- `--grouping` Filter by tournament grouping (partial name match).
- `--min-teams` Minimum number of teams.
- `--live-only` Show only live matches.
- `--upcoming-only` Show only upcoming matches.
- `--completed-only` Show only completed matches.
- `--limit` Maximum number of matches to return.
- `--output`, `-o` Output format: `table`, `json`, `yaml`.

`tournaments brackets <tournamentID>` — Get brackets for a tournament.
- `--completed-only` Show only completed matches.
- `--live-only` Show only live matches.
- `--upcoming-only` Show only upcoming matches.
- `--team` Filter by team name (case-insensitive partial match).
- `--match-type` Filter by match type (e.g., `BO5`, `BO7`).
- `--output`, `-o` Output format: `table`, `json`, `yaml`.

`matches list <tournamentID>` — List matches for a tournament.
- `--completed-only` Show only completed matches.
- `--live-only` Show only live matches.
- `--upcoming-only` Show only upcoming matches.
- `--team` Filter by team name or shorthand (case-insensitive partial match).
- `--match-type` Filter by match type (e.g., `BO5`, `BO7`).
- `--output`, `-o` Output format: `table`, `json`, `yaml`.

`matches get <matchID>` — Get detailed information for a match.
- `--output`, `-o` Output format: `table`, `json`, `yaml`.

Notes:
- The status filters (`--live-only`, `--upcoming-only`, `--completed-only`) are mutually exclusive.

**Output Formats**

- `tournaments list`: `table`, `json`, `csv`, `yaml`
- `tournaments matches`, `tournaments brackets`, `matches list`, `matches get`: `table`, `json`, `yaml`

**Examples**

List tournaments in a specific region and circuit:

```bash
rlcs-cli tournaments list --circuit 2026 --region NA
```

Show only upcoming tournaments:

```bash
rlcs-cli tournaments list --upcoming
```

Show only ongoing tournaments in EU:

```bash
rlcs-cli tournaments list --ongoing --region EU
```

List live matches across tournaments (limit to 10):

```bash
rlcs-cli tournaments matches --live-only --limit 10
```

List matches for a tournament filtered by team name:

```bash
rlcs-cli matches list <tournamentID> --team "Falcons"
```

Get match details by match ID:

```bash
rlcs-cli matches get <matchID>
```

Get brackets filtered by team and match type:

```bash
rlcs-cli tournaments brackets <tournamentID> --team "G2" --match-type BO7
```

Output as JSON/YAML/CSV:

```bash
rlcs-cli tournaments list --output json
rlcs-cli tournaments matches --output yaml
rlcs-cli tournaments list --output csv
```

**Development**

Run tests:

```bash
go test ./...
```
