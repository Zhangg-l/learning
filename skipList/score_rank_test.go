package main

import (
	"fmt"
	"testing"
	"time"
)

func TestScoreRank(t *testing.T) {
	for i := 0; i < 100; i++ {
		userID := fmt.Sprintf("user%d", i)
		res := ScoreRankServre.AddPlayScore(userID, int64(i), time.Now().Unix())
		fmt.Printf("%s add res %v\n", userID, res)
	}
	fmt.Println(1)
	// time.Sleep(1 * time.Second)
	// ts := time.Now().Unix()
	// res := ScoreRankServre.UpdateScore("user1", 101, ts)
	// fmt.Printf("update res %v\n", res)
	// res = ScoreRankServre.UpdateScore("user1", 102, ts)
	// fmt.Printf("update res %v\n", res)
	playerRank := ScoreRankServre.GetPlayerRank("user1")
	if playerRank != 99 {
		t.Errorf("player rank error")
	}
	arroundRanks := ScoreRankServre.GetPlayerRankRange("user5", 10)
	fmt.Printf("arroundRanks %v\n", arroundRanks)

	topNList := ScoreRankServre.GetTopN(100)
	fmt.Printf("topN %v\n", topNList)

}
