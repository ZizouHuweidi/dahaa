package service

import (
	"testing"

	"github.com/zizouhuweidi/dahaa/internal/domain"
)

func TestCalculateRoundScores(t *testing.T) {
	service := &GameService{}
	game := &domain.Game{
		Players: []domain.Player{
			{ID: "p1", Name: "One"},
			{ID: "p2", Name: "Two"},
			{ID: "p3", Name: "Three"},
		},
	}
	round := &domain.Round{
		AnswerPool: domain.AnswerPool{
			Options: []domain.Answer{
				{ID: "correct", Kind: domain.AnswerKindCorrect, PlayerID: "system", Text: "Cairo", Votes: []string{"p1", "p2"}},
				{ID: "fake", Kind: domain.AnswerKindFake, PlayerID: "p1", Text: "Istanbul", Votes: []string{"p3"}},
				{ID: "filler", Kind: domain.AnswerKindFiller, PlayerID: "system", Text: "Rabat", Votes: []string{"p2"}},
			},
		},
	}

	service.calculateRoundScores(game, round)

	assertScore(t, game, "p1", 3) // 2 for correct vote, 1 for fooling p3.
	assertScore(t, game, "p2", 2) // 2 for correct vote, no points for filler vote.
	assertScore(t, game, "p3", 0)
}

func assertScore(t *testing.T, game *domain.Game, playerID string, want int) {
	t.Helper()
	for _, player := range game.Players {
		if player.ID == playerID {
			if player.Score != want {
				t.Fatalf("player %s score = %d, want %d", playerID, player.Score, want)
			}
			return
		}
	}
	t.Fatalf("player %s not found", playerID)
}
