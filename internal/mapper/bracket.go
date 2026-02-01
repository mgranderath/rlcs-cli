package mapper

import (
	"fmt"
	"time"

	"github.com/mgranderath/rlcs-cli/internal/api/blast"
	"github.com/mgranderath/rlcs-cli/internal/domain"
)

const timeFormat = "2006-01-02T15:04:05.000Z"

// ToDomainBracket converts a Blast API bracket to the domain model
func ToDomainBracket(api blast.Bracket) (domain.Bracket, error) {
	startDate, err := time.Parse(timeFormat, api.StartDate)
	if err != nil {
		return domain.Bracket{}, fmt.Errorf("failed to parse start date: %w", err)
	}

	endDate, err := time.Parse(timeFormat, api.EndDate)
	if err != nil {
		return domain.Bracket{}, fmt.Errorf("failed to parse end date: %w", err)
	}

	matches := make([]domain.Match, 0, len(api.Matches))
	for _, apiMatch := range api.Matches {
		match, err := ToDomainMatch(apiMatch)
		if err != nil {
			return domain.Bracket{}, fmt.Errorf("failed to map match %s: %w", apiMatch.UUID, err)
		}
		matches = append(matches, match)
	}

	return domain.Bracket{
		TournamentUUID:         api.TournamentUUID,
		TournamentName:         api.TournamentName,
		ParentTournamentName:   api.ParentTournamentName,
		ParentTournamentFormat: api.ParentTournamentFormat,
		CircuitName:            api.CircuitName,
		StartDate:              startDate,
		EndDate:                endDate,
		Index:                  api.Index,
		Label:                  api.Label,
		Format:                 api.Format,
		NumberOfTeams:          api.NumberOfTeams,
		Matches:                matches,
	}, nil
}

// ToDomainMatch converts a Blast API match to the domain model
func ToDomainMatch(api blast.Match) (domain.Match, error) {
	timeOfSeries, err := time.Parse(timeFormat, api.TimeOfSeries)
	if err != nil {
		return domain.Match{}, fmt.Errorf("failed to parse time of series: %w", err)
	}

	maps := make([]domain.MatchMap, 0, len(api.Maps))
	for _, apiMap := range api.Maps {
		m, err := ToDomainMap(apiMap)
		if err != nil {
			return domain.Match{}, fmt.Errorf("failed to map map %s: %w", apiMap.UUID, err)
		}
		maps = append(maps, m)
	}

	var winnerGoesTo *domain.BracketDestination
	if api.WinnerGoesTo != nil {
		winnerGoesTo = &domain.BracketDestination{
			TournamentUUID:  api.WinnerGoesTo.TournamentUUID,
			SeriesUUID:      api.WinnerGoesTo.SeriesUUID,
			BracketPosition: api.WinnerGoesTo.BracketPosition,
		}
	}

	var loserGoesTo *domain.BracketDestination
	if api.LoserGoesTo != nil {
		loserGoesTo = &domain.BracketDestination{
			TournamentUUID:  api.LoserGoesTo.TournamentUUID,
			SeriesUUID:      api.LoserGoesTo.SeriesUUID,
			BracketPosition: api.LoserGoesTo.BracketPosition,
		}
	}

	return domain.Match{
		UUID:         api.UUID,
		Type:         api.Type,
		Index:        api.Index,
		Name:         api.Name,
		TimeOfSeries: timeOfSeries,
		TeamA: domain.MatchTeam{
			UUID:         api.TeamA.UUID,
			Name:         api.TeamA.Name,
			Shorthand:    api.TeamA.Shorthand,
			Location:     api.TeamA.Location,
			IsEliminated: api.TeamA.IsEliminated,
		},
		TeamB: domain.MatchTeam{
			UUID:         api.TeamB.UUID,
			Name:         api.TeamB.Name,
			Shorthand:    api.TeamB.Shorthand,
			Location:     api.TeamB.Location,
			IsEliminated: api.TeamB.IsEliminated,
		},
		TeamAScore:   api.TeamAScore,
		TeamBScore:   api.TeamBScore,
		Maps:         maps,
		ExternalID:   api.ExternalID,
		WinnerGoesTo: winnerGoesTo,
		LoserGoesTo:  loserGoesTo,
		IsLive:       api.IsLive,
		IsCompleted:  api.IsCompleted,
	}, nil
}

// ToDomainMap converts a Blast API map to the domain model
func ToDomainMap(api blast.Map) (domain.MatchMap, error) {
	scheduledStart, err := time.Parse(timeFormat, api.ScheduledStartTime)
	if err != nil {
		return domain.MatchMap{}, fmt.Errorf("failed to parse scheduled start time: %w", err)
	}

	actualStart, err := time.Parse(timeFormat, api.ActualStartTime)
	if err != nil {
		return domain.MatchMap{}, fmt.Errorf("failed to parse actual start time: %w", err)
	}

	matchEnded, err := time.Parse(timeFormat, api.MatchEndedTime)
	if err != nil {
		return domain.MatchMap{}, fmt.Errorf("failed to parse match ended time: %w", err)
	}

	return domain.MatchMap{
		UUID:               api.UUID,
		ScheduledStartTime: scheduledStart,
		ActualStartTime:    actualStart,
		Name:               api.Name,
		MatchEndedTime:     matchEnded,
		TeamAScore:         api.TeamAScore,
		TeamBScore:         api.TeamBScore,
		ExternalID:         api.ExternalID,
	}, nil
}

// ToDomainBrackets converts a slice of Blast API brackets to domain models
func ToDomainBrackets(apiBrackets []blast.Bracket) ([]domain.Bracket, error) {
	result := make([]domain.Bracket, 0, len(apiBrackets))

	for _, api := range apiBrackets {
		bracket, err := ToDomainBracket(api)
		if err != nil {
			return nil, fmt.Errorf("failed to map bracket %s: %w", api.TournamentUUID, err)
		}
		result = append(result, bracket)
	}

	return result, nil
}
