package main

import (
	"sync"
)

var (
	ScoreRankServre *ScoreRank
)

func init() {
	ScoreRankServre = &ScoreRank{
		Meta:     NewSkipList(),
		UserMark: make(map[string]*PlayerScoreInfo),
	}
	ScoreRankServre.Meta.InitSkipList()
}

type PlayerScoreInfo struct {
	Score     int64
	Timestamp int64
}
type RankInfo struct {
	PlayerId string
	Score    int64
}

// 排行榜
type ScoreRank struct {
	Locker   sync.RWMutex
	Meta     *SkipList
	UserMark map[string]*PlayerScoreInfo //时间戳 + 积分
}

// 添加玩家积分
func (p *ScoreRank) AddPlayScore(playerId string, score, timestamp int64) (succ bool) {
	if playerId == "" || score < 0 || timestamp <= 0 {
		return
	}
	p.Locker.Lock()
	defer p.Locker.Unlock()
	_, ok := p.UserMark[playerId]
	if ok {
		return
	}
	p.UserMark[playerId] = &PlayerScoreInfo{
		Score:     score,
		Timestamp: timestamp,
	}
	p.Meta.Insert(playerId, score)
	succ = true
	return
}

// 更新玩家积分
func (p *ScoreRank) UpdateScore(playerId string, score, timestamp int64) (succ bool) {
	if playerId == "" || score < 0 || timestamp < 0 {
		return
	}
	p.Locker.Lock()
	defer p.Locker.Unlock()
	scoreInfo, ok := p.UserMark[playerId]

	if !ok || scoreInfo.Timestamp >= timestamp {
		return
	}
	if p.Meta.Delete(playerId, scoreInfo.Score) {
		p.Meta.Insert(playerId, score)
		scoreInfo.Score = score
		scoreInfo.Timestamp = timestamp
		succ = true
	}
	return
}

// 获取玩家排名
func (p *ScoreRank) GetPlayerRank(playerId string) (rank int64) {
	rank = -1
	if playerId == "" {
		return
	}
	p.Locker.RLock()
	defer p.Locker.RUnlock()
	scoreInfo, ok := p.UserMark[playerId]
	if !ok {
		return
	}
	playScore, rank := p.Meta.Find(playerId, scoreInfo.Score)
	if playScore == nil {
		return
	}
	// rank = playScore.Rank
	return
}

// 获取排行榜前N名
func (p *ScoreRank) GetTopN(n int64) (topN []RankInfo) {
	if n <= 0 || p.Meta.Header == nil {
		return
	}
	p.Locker.RLock()
	defer p.Locker.RUnlock()
	if p.Meta.Length < n {
		n = p.Meta.Length
	}

	topN = make([]RankInfo, 0, n)
	pNode := p.Meta.Header
	for n > 0 && pNode.Levels[0].Next != nil {
		pNode = pNode.Levels[0].Next
		topN = append(topN, RankInfo{
			PlayerId: pNode.PlayerId,
			Score:    pNode.Score,
		})
		n--
	}
	return
}

// // 获取玩家周边排名
func (p *ScoreRank) GetPlayerRankRange(playerId string, n int64) (aroundRank []RankInfo) {
	if n <= 0 || p.Meta.Header == nil {
		return
	}
	p.Locker.RLock()
	defer p.Locker.RUnlock()
	aroundRank = make([]RankInfo, 0, n*2+1)
	var (
		beforeN, afterN = n, n

		pNode *SkipNode
	)
	scoreInfo, ok := p.UserMark[playerId]
	if !ok {
		return
	}
	playScore, _ := p.Meta.Find(playerId, scoreInfo.Score)
	if playScore == nil {
		return
	}
	pNode = playScore

	for pNode.BackWard != nil && beforeN > 0 {
		pNode = pNode.BackWard
		beforeN--
	}
	for pNode != playScore {
		aroundRank = append(aroundRank, RankInfo{
			PlayerId: pNode.PlayerId,
			Score:    pNode.Score,
		})
		pNode = pNode.Levels[0].Next
	}
	aroundRank = append(aroundRank, RankInfo{
		PlayerId: pNode.PlayerId,
		Score:    pNode.Score,
	})
	for pNode.Levels[0].Next != nil && afterN > 0 {
		aroundRank = append(aroundRank, RankInfo{
			PlayerId: pNode.Levels[0].Next.PlayerId,
			Score:    pNode.Levels[0].Next.Score,
		})
		pNode = pNode.Levels[0].Next
		afterN--
	}

	return
}
