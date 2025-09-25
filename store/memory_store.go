package store

import (
	"Go-mini-ranking/model"
	"sort"
	"sync"
)

// MemoryStore 基于内存的排行榜实现
type MemoryStore struct {
	// 存储结构: leaderboardID -> map[playerID]PlayerScore
	leaderboards map[string]map[string]model.PlayerScore
	mu           sync.RWMutex
}

// NewMemoryStore 创建新的内存存储
func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		leaderboards: make(map[string]map[string]model.PlayerScore),
	}
}

// UpdateScore 更新玩家分数
func (s *MemoryStore) UpdateScore(score model.PlayerScore) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// 确保排行榜存在
	if _, exists := s.leaderboards[score.LeaderboardID]; !exists {
		s.leaderboards[score.LeaderboardID] = make(map[string]model.PlayerScore)
	}

	leaderboard := s.leaderboards[score.LeaderboardID]

	// 检查玩家是否已存在
	if existing, exists := leaderboard[score.PlayerID]; exists {
		// 只在新分数更高时更新，同分不更新（保持先达到该分数的状态）
		if score.Score > existing.Score {
			leaderboard[score.PlayerID] = score
		}
	} else {
		// 新玩家直接添加
		leaderboard[score.PlayerID] = score
	}
}

// GetPlayerRank 获取玩家当前排名
func (s *MemoryStore) GetPlayerRank(leaderboardID, playerID string) (int, int64, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	// 检查排行榜是否存在
	leaderboard, exists := s.leaderboards[leaderboardID]
	if !exists {
		return 0, 0, false
	}

	// 检查玩家是否存在
	playerScore, exists := leaderboard[playerID]
	if !exists {
		return 0, 0, false
	}

	// 获取排序后的玩家列表
	sortedPlayers := s.getSortedPlayers(leaderboardID)

	// 查找玩家排名
	for rank, info := range sortedPlayers {
		if info.PlayerID == playerID {
			return rank + 1, playerScore.Score, true // 排名从1开始
		}
	}

	return 0, 0, false
}

// GetTopPlayers 获取前N名玩家
func (s *MemoryStore) GetTopPlayers(leaderboardID string, limit int) []model.RankingInfo {
	s.mu.RLock()
	defer s.mu.RUnlock()

	sortedPlayers := s.getSortedPlayers(leaderboardID)

	// 限制返回数量
	if limit <= 0 || limit > len(sortedPlayers) {
		limit = len(sortedPlayers)
	}

	return sortedPlayers[:limit]
}

// GetPlayersAround 获取玩家周围的玩家
func (s *MemoryStore) GetPlayersAround(leaderboardID, playerID string, count int) []model.RankingInfo {
	s.mu.RLock()
	defer s.mu.RUnlock()

	// 检查玩家是否存在
	leaderboard, exists := s.leaderboards[leaderboardID]
	if !exists {
		return []model.RankingInfo{}
	}

	if _, exists := leaderboard[playerID]; !exists {
		return []model.RankingInfo{}
	}

	// 获取排序后的玩家列表
	sortedPlayers := s.getSortedPlayers(leaderboardID)

	// 找到玩家位置
	playerIndex := -1
	for i, info := range sortedPlayers {
		if info.PlayerID == playerID {
			playerIndex = i
			break
		}
	}

	if playerIndex == -1 {
		return []model.RankingInfo{}
	}

	// 计算需要返回的范围
	half := count / 2
	start := playerIndex - half
	end := playerIndex + half + (count % 2) // 处理奇数情况

	// 确保范围在有效范围内
	if start < 0 {
		start = 0
	}
	if end > len(sortedPlayers) {
		end = len(sortedPlayers)
	}

	return sortedPlayers[start:end]
}

// getSortedPlayers 获取排序后的玩家列表
func (s *MemoryStore) getSortedPlayers(leaderboardID string) []model.RankingInfo {
	leaderboard, exists := s.leaderboards[leaderboardID]
	if !exists || len(leaderboard) == 0 {
		return []model.RankingInfo{}
	}

	// 转换为切片以便排序
	players := make([]model.PlayerScore, 0, len(leaderboard))
	for _, score := range leaderboard {
		players = append(players, score)
	}

	// 排序：先按分数降序，分数相同则按更新时间升序（先达到该分数的排前面）
	sort.Slice(players, func(i, j int) bool {
		if players[i].Score != players[j].Score {
			return players[i].Score > players[j].Score
		}
		return players[i].LastUpdateTime.Before(players[j].LastUpdateTime)
	})

	// 转换为排名信息
	rankingList := make([]model.RankingInfo, 0, len(players))
	for i, player := range players {
		rankingList = append(rankingList, model.RankingInfo{
			PlayerID: player.PlayerID,
			Score:    player.Score,
			Rank:     i + 1, // 排名从1开始
		})
	}

	return rankingList
}
