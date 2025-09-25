package main

import (
	"Go-mini-ranking/service"
	"Go-mini-ranking/store"
	"fmt"
	"time"
)

func main() {
	// 初始化排行榜系统
	memStore := store.NewMemoryStore()
	rankService := service.NewRankingService(memStore)

	// 定义一个排行榜ID
	leaderboardID := "game_ranking_001"

	// 模拟一些玩家分数更新
	players := []struct {
		id    string
		score int64
		delay time.Duration
	}{
		{"player_01", 950, 0 * time.Millisecond},
		{"player_02", 1000, 10 * time.Millisecond},
		{"player_03", 900, 20 * time.Millisecond},
		{"player_04", 1050, 30 * time.Millisecond},
		{"player_05", 950, 40 * time.Millisecond}, // 与player_01分数相同但时间更晚
		{"player_06", 980, 50 * time.Millisecond},
		{"player_07", 920, 60 * time.Millisecond},
		{"player_08", 1020, 70 * time.Millisecond},
		{"player_09", 970, 80 * time.Millisecond},
		{"player_10", 990, 90 * time.Millisecond},
	}

	// 提交玩家分数（添加延迟以模拟不同时间提交）
	fmt.Println("=== 提交玩家分数 ===")
	for _, p := range players {
		time.Sleep(p.delay)
		rankService.UpdatePlayerScore(leaderboardID, p.id, p.score)
		fmt.Printf("玩家 %s 提交分数: %d\n", p.id, p.score)
	}

	// 演示1: 获取前三名玩家
	fmt.Println("\n=== 排行榜前三名 ===")
	top3 := rankService.GetTopPlayers(leaderboardID, 3)
	for _, info := range top3 {
		fmt.Printf("第 %d 名: 玩家 %s, 分数: %d\n", info.Rank, info.PlayerID, info.Score)
	}

	// 演示2: 查询特定玩家排名
	fmt.Println("\n=== 玩家排名查询 ===")
	playerID := "player_01"
	rank, score, exists := rankService.GetPlayerRank(leaderboardID, playerID)
	if exists {
		fmt.Printf("玩家 %s 的排名: 第 %d 名, 分数: %d\n", playerID, rank, score)
	}

	// 演示3: 查询玩家周围的排名（前后各2名，共5名）
	fmt.Println("\n=== 玩家周围排名 ===")
	around := rankService.GetPlayersAround(leaderboardID, playerID, 5)
	for _, info := range around {
		fmt.Printf("第 %d 名: 玩家 %s, 分数: %d\n", info.Rank, info.PlayerID, info.Score)
	}

	// 演示4: 玩家更新分数后再次查询
	fmt.Println("\n=== 玩家更新分数后 ===")
	rankService.UpdatePlayerScore(leaderboardID, playerID, 1010)
	fmt.Printf("玩家 %s 更新分数为: 1010\n", playerID)

	newRank, newScore, _ := rankService.GetPlayerRank(leaderboardID, playerID)
	fmt.Printf("玩家 %s 的新排名: 第 %d 名, 新分数: %d\n", playerID, newRank, newScore)

	// 再次查询前三名
	fmt.Println("\n=== 更新后的排行榜前三名 ===")
	newTop3 := rankService.GetTopPlayers(leaderboardID, 3)
	for _, info := range newTop3 {
		fmt.Printf("第 %d 名: 玩家 %s, 分数: %d\n", info.Rank, info.PlayerID, info.Score)
	}
}
