package mapper

import (
	"fmt"
	"strings"
	"time"

	"github.com/mgranderath/rlcs-cli/internal/api/blast"
	"github.com/mgranderath/rlcs-cli/internal/domain"
)

const dateFormat = "2006-01-02"

// ToDomainTournament converts a Blast API tournament to the domain model
func ToDomainTournament(api blast.Tournament) (domain.Tournament, error) {
	startDate, err := time.Parse(dateFormat, api.StartDate)
	if err != nil {
		return domain.Tournament{}, fmt.Errorf("failed to parse start date: %w", err)
	}

	endDate, err := time.Parse(dateFormat, api.EndDate)
	if err != nil {
		return domain.Tournament{}, fmt.Errorf("failed to parse end date: %w", err)
	}

	region := parseRegion(api.Region)
	tournamentType := determineTournamentType(api.Name, api.Grouping, api.Region)
	isMajor := api.Region == "" && api.Grouping == ""

	return domain.Tournament{
		ID:          api.ID,
		Name:        api.Name,
		StartDate:   startDate,
		EndDate:     endDate,
		CircuitID:   api.CircuitID,
		PrizePool:   api.PrizePool,
		Location:    api.Location,
		TeamCount:   api.NumberOfTeams,
		Region:      region,
		Type:        tournamentType,
		Description: api.Description,
		IsOnline:    api.Location == "Online",
		IsMajor:     isMajor,
	}, nil
}

// ToDomainTournaments converts a slice of Blast API tournaments to domain models
func ToDomainTournaments(apiTournaments []blast.Tournament) ([]domain.Tournament, error) {
	result := make([]domain.Tournament, 0, len(apiTournaments))

	for _, api := range apiTournaments {
		domain, err := ToDomainTournament(api)
		if err != nil {
			return nil, fmt.Errorf("failed to map tournament %s: %w", api.ID, err)
		}
		result = append(result, domain)
	}

	return result, nil
}

// parseRegion converts API region string to domain Region type
func parseRegion(region string) domain.Region {
	switch strings.ToUpper(region) {
	case "NA":
		return domain.RegionNA
	case "EU":
		return domain.RegionEU
	case "APAC":
		return domain.RegionAPAC
	case "SAM":
		return domain.RegionSAM
	case "OCE":
		return domain.RegionOCE
	case "MENA":
		return domain.RegionMENA
	case "SSA":
		return domain.RegionSSA
	default:
		return domain.RegionNone
	}
}

// determineTournamentType figures out the tournament type from the name and grouping
func determineTournamentType(name, grouping, region string) domain.TournamentType {
	lowerName := strings.ToLower(name)

	if strings.Contains(lowerName, "world championship") {
		return domain.TypeWorldChampionship
	}
	if strings.Contains(lowerName, "kick-off") || strings.Contains(lowerName, "kickoff") {
		return domain.TypeKickoff
	}
	if region == "" && grouping == "" {
		return domain.TypeMajor
	}
	if strings.Contains(lowerName, "open") {
		return domain.TypeOpen
	}

	return domain.TypeOpen
}
