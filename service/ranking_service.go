package service

import (
	"Go-mini-ranking/model"
	"Go-mini-ranking/store"
	"time"
)

// RankingService 排行榜服务
type RankingService struct {
	store store.RankingStore
}

// NewRankingService 创建新的排行榜服务
func NewRankingService(store store.RankingStore) *RankingService {
	return &RankingService{
		store: store,
	}
}

// UpdatePlayerScore 更新玩家分数
func (s *RankingService) UpdatePlayerScore(leaderboardID, playerID string, score int64) {
	// 创建分数记录，包含当前时间作为更新时间
	playerScore := model.PlayerScore{
		LeaderboardID:  leaderboardID,
		PlayerID:       playerID,
		Score:          score,
		LastUpdateTime: time.Now(),
	}

	s.store.UpdateScore(playerScore)
}

// GetPlayerRank 获取玩家排名
func (s *RankingService) GetPlayerRank(leaderboardID, playerID string) (int, int64, bool) {
	return s.store.GetPlayerRank(leaderboardID, playerID)
}

// GetTopPlayers 获取前N名玩家
func (s *RankingService) GetTopPlayers(leaderboardID string, limit int) []model.RankingInfo {
	if limit <= 0 {
		limit = 10 // 默认返回前10名
	}
	return s.store.GetTopPlayers(leaderboardID, limit)
}

// GetPlayersAround 获取玩家周围的玩家
func (s *RankingService) GetPlayersAround(leaderboardID, playerID string, count int) []model.RankingInfo {
	if count <= 0 {
		count = 5 // 默认返回前后共5名
	}
	return s.store.GetPlayersAround(leaderboardID, playerID, count)
}
