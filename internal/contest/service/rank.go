package service

import (
	"context"
	"net/http"
	"sort"
	"time"

	"github.com/f-code-club/rode-battle-api/internal/contest/repository"
	"github.com/f-code-club/rode-battle-api/internal/shared/errors"
	"github.com/google/uuid"
)

type Detail struct {
	ProblemID       uuid.UUID `json:"problem_id"`
	ProblemPosition int       `json:"problem_position"`
	SubmissionCount int       `json:"submission_count"`
	Score           int       `json:"score"`
	LastSubmit      time.Time `json:"last_submit"`
}

type Ranking struct {
	Name    string   `json:"name"`
	Score   float64  `json:"score"`
	Penalty float64  `json:"penalty"`
	Details []Detail `json:"details"`
}

type submissionRow = repository.GetContestSubmissionsForRankingRow

const (
	PENALTY_PER_SUBMISSION = 10
	SCORE_PER_PROBLEM      = 1
)

func (s *Service) GetRank(
	ctx context.Context,
	contestID uuid.UUID,
) ([]Ranking, error) {
	queries := repository.New(s.pool)

	rows, err := queries.GetContestSubmissionsForRanking(ctx, contestID)
	if err != nil {
		return nil, errors.Wrap(
			http.StatusInternalServerError,
			"failed to get contest ranking",
			err,
		)
	}

	accountNames, grouped := groupSubmissions(rows)

	result := make([]Ranking, 0, len(grouped))

	for accountID, problems := range grouped {
		details, score, penalty := buildRanking(problems)

		result = append(result, Ranking{
			Name:    accountNames[accountID],
			Score:   score,
			Penalty: penalty,
			Details: details,
		})
	}

	sortRankings(result)

	return result, nil
}

func groupSubmissions(
	rows []submissionRow,
) (map[uuid.UUID]string, map[uuid.UUID]map[uuid.UUID][]submissionRow) {
	accountNames := make(map[uuid.UUID]string)
	grouped := make(map[uuid.UUID]map[uuid.UUID][]submissionRow)

	for _, row := range rows {
		if _, ok := grouped[row.AccountID]; !ok {
			grouped[row.AccountID] = make(map[uuid.UUID][]submissionRow)
			accountNames[row.AccountID] = row.AccountName
		}

		grouped[row.AccountID][row.ProblemID] = append(
			grouped[row.AccountID][row.ProblemID],
			row,
		)
	}

	return accountNames, grouped
}

func buildRanking(
	problems map[uuid.UUID][]submissionRow,
) ([]Detail, float64, float64) {
	details := make([]Detail, 0, len(problems))

	var totalScore, totalPenalty float64

	for _, submissions := range problems {
		detail, penalty := calculateProblemResult(submissions)

		details = append(details, detail)
		totalScore += float64(detail.Score)
		totalPenalty += penalty
	}

	sort.Slice(details, func(i, j int) bool {
		return details[i].ProblemPosition < details[j].ProblemPosition
	})

	return details, totalScore, totalPenalty
}

func calculateProblemResult(
	submissions []submissionRow,
) (Detail, float64) {
	sort.Slice(submissions, func(i, j int) bool {
		return submissions[i].CreatedAt.Time.Before(submissions[j].CreatedAt.Time)
	})

	truncated, hasAccepted := truncateAtFirstAccepted(submissions)

	last := truncated[len(truncated)-1]

	submissionCount := len(truncated)

	score := 0
	if hasAccepted {
		score = SCORE_PER_PROBLEM
	}

	detail := Detail{
		ProblemID:       last.ProblemID,
		ProblemPosition: intValue(last.ProblemPosition),
		SubmissionCount: submissionCount,
		Score:           score,
		LastSubmit:      last.CreatedAt.Time,
	}

	penalty := calculatePenalty(
		last.Language,
		score,
		submissionCount,
		last.CreatedAt.Time,
		last.ContestStartTime.Time,
	)

	return detail, penalty
}

func truncateAtFirstAccepted(
	submissions []submissionRow,
) ([]submissionRow, bool) {
	for i, sub := range submissions {
		if sub.Verdict != nil && *sub.Verdict == repository.VerdictAccepted {
			return submissions[:i+1], true
		}
	}

	return submissions, false
}

func calculatePenalty(
	language repository.Language,
	score int,
	submissionCount int,
	lastSubmit time.Time,
	contestStart time.Time,
) float64 {
	if language == repository.LanguageHtml {
		return float64(submissionCount * PENALTY_PER_SUBMISSION)
	}

	if score == 0 {
		return 0
	}

	minutes := lastSubmit.Sub(contestStart).Minutes()

	return minutes + float64(submissionCount*PENALTY_PER_SUBMISSION)
}

func sortRankings(rankings []Ranking) {
	sort.Slice(rankings, func(i, j int) bool {
		if rankings[i].Score != rankings[j].Score {
			return rankings[i].Score > rankings[j].Score
		}

		return rankings[i].Penalty < rankings[j].Penalty
	})
}

func intValue(value *int32) int {
	if value == nil {
		return 0
	}

	return int(*value)
}
