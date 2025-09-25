package model

import "time"

// PlayerScore 玩家分数记录
type PlayerScore struct {
	PlayerID       string    // 玩家唯一标识
	Score          int64     // 分数
	LastUpdateTime time.Time // 最后更新时间（用于同分排序）
	LeaderboardID  string    // 所属排行榜ID
}

// RankingInfo 排名信息
type RankingInfo struct {
	PlayerID string // 玩家ID
	Score    int64  // 分数
	Rank     int    // 排名
}
