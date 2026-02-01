package mapper

import (
	"fmt"
	"time"

	"github.com/mgranderath/rlcs-cli/internal/api/blast"
	"github.com/mgranderath/rlcs-cli/internal/domain"
)

// ToDomainMatchesFromResponse converts API match responses to domain matches
func ToDomainMatchesFromResponse(apiMatches []blast.MatchResponse) ([]domain.Match, error) {
	result := make([]domain.Match, 0, len(apiMatches))

	for _, apiMatch := range apiMatches {
		match, err := toDomainMatchFromResponse(apiMatch)
		if err != nil {
			return nil, fmt.Errorf("failed to map match %s: %w", apiMatch.ID, err)
		}
		result = append(result, match)
	}

	return result, nil
}

func toDomainMatchFromResponse(api blast.MatchResponse) (domain.Match, error) {
	timeOfSeries, err := time.Parse(timeFormat, api.ScheduledAt)
	if err != nil {
		return domain.Match{}, fmt.Errorf("failed to parse scheduled time: %w", err)
	}

	maps := make([]domain.MatchMap, 0, len(api.Maps))
	for _, apiMap := range api.Maps {
		m, err := toDomainMatchMapFromResponse(apiMap)
		if err != nil {
			return domain.Match{}, fmt.Errorf("failed to map map %s: %w", apiMap.ID, err)
		}
		maps = append(maps, m)
	}

	// Infer status from map timestamps
	isCompleted, isLive := inferMatchStatus(api.Maps)

	return domain.Match{
		UUID:         api.ID,
		Type:         api.Type,
		Index:        api.Index,
		Name:         api.Name,
		TimeOfSeries: timeOfSeries,
		TeamA: domain.MatchTeam{
			UUID:      api.TeamA.ID,
			Name:      api.TeamA.Name,
			Shorthand: api.TeamA.ShortName,
			Location:  api.TeamA.Nationality,
		},
		TeamB: domain.MatchTeam{
			UUID:      api.TeamB.ID,
			Name:      api.TeamB.Name,
			Shorthand: api.TeamB.ShortName,
			Location:  api.TeamB.Nationality,
		},
		TeamAScore:  api.TeamAScore,
		TeamBScore:  api.TeamBScore,
		Maps:        maps,
		ExternalID:  api.ExternalID,
		IsLive:      isLive,
		IsCompleted: isCompleted,
	}, nil
}

func toDomainMatchMapFromResponse(api blast.MatchResponseMap) (domain.MatchMap, error) {
	scheduledStart, err := time.Parse(timeFormat, api.ScheduledAt)
	if err != nil {
		return domain.MatchMap{}, fmt.Errorf("failed to parse scheduled start time: %w", err)
	}

	// Handle empty time strings for fields that may not be set
	var actualStart time.Time
	if api.StartedAt != "" {
		parsedTime, err := time.Parse(timeFormat, api.StartedAt)
		if err != nil {
			return domain.MatchMap{}, fmt.Errorf("failed to parse started at time: %w", err)
		}
		actualStart = parsedTime
	}

	var matchEnded time.Time
	if api.EndedAt != "" {
		parsedTime, err := time.Parse(timeFormat, api.EndedAt)
		if err != nil {
			return domain.MatchMap{}, fmt.Errorf("failed to parse ended at time: %w", err)
		}
		matchEnded = parsedTime
	}

	return domain.MatchMap{
		UUID:               api.ID,
		ScheduledStartTime: scheduledStart,
		ActualStartTime:    actualStart,
		Name:               api.Name,
		MatchEndedTime:     matchEnded,
		TeamAScore:         api.TeamAScore,
		TeamBScore:         api.TeamBScore,
		ExternalID:         api.ExternalID,
	}, nil
}

// inferMatchStatus determines if a match is completed, live, or upcoming based on map timestamps
// - Completed: at least one map has started and all maps that have started have ended
// - Live: at least one map has started but not all have ended
// - Upcoming: no maps have started
func inferMatchStatus(maps []blast.MatchResponseMap) (isCompleted bool, isLive bool) {
	if len(maps) == 0 {
		return false, false
	}

	hasStarted := false
	allEnded := true

	for _, m := range maps {
		if m.StartedAt != "" {
			hasStarted = true
			if m.EndedAt == "" {
				allEnded = false
			}
		} else {
			// Map hasn't started, so match isn't complete yet
			allEnded = false
		}
	}

	if !hasStarted {
		return false, false // Upcoming
	}

	if allEnded {
		return true, false // Completed
	}

	return false, true // Live
}
