package store

import (
	"Go-mini-ranking/model"
)

// RankingStore 排行榜存储接口
type RankingStore interface {
	// UpdateScore 更新玩家分数，如果分数更高则更新，同分则保留先达到该分数的玩家
	UpdateScore(score model.PlayerScore)

	// GetPlayerRank 获取玩家当前排名
	GetPlayerRank(leaderboardID, playerID string) (int, int64, bool)

	// GetTopPlayers 获取前N名玩家
	GetTopPlayers(leaderboardID string, limit int) []model.RankingInfo

	// GetPlayersAround 获取玩家周围N名玩家（前后各N/2）
	GetPlayersAround(leaderboardID, playerID string, count int) []model.RankingInfo
}
