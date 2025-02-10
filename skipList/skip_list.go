package main

import (
	"fmt"
	"math/rand"

	"github.com/davecgh/go-spew/spew"
)

const (
	MaxLevel = 4
)

type levelsNode struct {
	Next *SkipNode
	Span int64
}
type SkipNode struct {
	Score    int64
	PlayerId string
	BackWard *SkipNode
	Levels   []levelsNode
}

type SkipList struct {
	Level  int8  //当前层高
	Length int64 //当前节点数
	Header *SkipNode
	Tail   *SkipNode
}

func NewSkipList() *SkipList {
	return &SkipList{}
}
func (p *SkipList) InitSkipList() {
	header := &SkipNode{
		Score:    0,
		PlayerId: "",
		BackWard: nil,
		Levels:   make([]levelsNode, MaxLevel),
	}

	for i := 0; i < MaxLevel; i++ {
		header.Levels[i].Span = 0
		header.Levels[i].Next = nil
	}
	p.Header = header
	p.Level = 1
}

// 插入节点
func (p *SkipList) Insert(playerId string, ele int64) {
	rank := make([]int64, MaxLevel)            //计算排位
	updateNodes := make([]*SkipNode, MaxLevel) //暂存插入点前一个元素
	for i := p.Level - 1; i >= 0; i-- {
		if i == p.Level-1 {
			rank[i] = 0
		} else {
			rank[i] = rank[i+1]
		}

		pNode := p.Header
		// 当前层寻找插入点
		for pNode.Levels[i].Next != nil &&
			(pNode.Levels[i].Next.Score >= ele) {
			rank[i] += pNode.Levels[i].Span
			pNode = pNode.Levels[i].Next
		}
		updateNodes[i] = pNode
	}

	// 获取随机层数
	newLevel := GetRandLevel()
	if newLevel > p.Level {
		for i := p.Level; i < newLevel; i++ {
			rank[i] = 0
			updateNodes[i] = p.Header
			updateNodes[i].Levels[i].Span = p.Length
		}
		p.Level = newLevel
	}

	newNode := &SkipNode{
		Score:    ele,
		PlayerId: playerId,
		BackWard: nil,
		Levels:   make([]levelsNode, newLevel),
	}
	for i := int8(0); i < newLevel; i++ {
		newNode.Levels[i].Next = updateNodes[i].Levels[i].Next
		updateNodes[i].Levels[i].Next = newNode
		newNode.Levels[i].Span = updateNodes[i].Levels[i].Span - (rank[0] - rank[i])
		updateNodes[i].Levels[i].Span = rank[0] - rank[i] + 1
	}

	for i := newLevel; i < p.Level; i++ {
		updateNodes[i].Levels[i].Span += 1
	}

	if updateNodes[0] != p.Header {
		newNode.BackWard = updateNodes[0]
	} else {
		newNode.BackWard = nil
	}

	if newNode.Levels[0].Next != nil {
		newNode.Levels[0].Next.BackWard = newNode
	} else {
		p.Tail = newNode
	}
	p.Length++
}

func (p *SkipList) Delete(player string, ele int64) bool {
	// 先查找 在删除
	flag := false
	updateNodes := make([]*SkipNode, MaxLevel)

	pNode := p.Header
	for i := p.Level - 1; i >= 0; i-- {
		for pNode.Levels[i].Next != nil && pNode.Levels[i].Next.Score > ele {
			pNode = pNode.Levels[i].Next
		}
		if pNode.Levels[i].Next != nil && pNode.Levels[i].Next.PlayerId == player &&
			pNode.Levels[i].Next.Score == ele {
			flag = true
		}
		updateNodes[i] = pNode
	}

	if !flag {
		return flag
	}

	// 获取待删除节点
	nodeToDelete := updateNodes[0].Levels[0].Next

	for i := 0; i < int(p.Level); i++ {
		if updateNodes[i].Levels[i].Next == nodeToDelete {
			updateNodes[i].Levels[i].Next = nodeToDelete.Levels[i].Next
			updateNodes[i].Levels[i].Span += nodeToDelete.Levels[i].Span - 1
		} else {
			updateNodes[i].Levels[i].Span -= 1
		}
	}

	if nodeToDelete.Levels[0].Next != nil {
		nodeToDelete.Levels[0].Next.BackWard = updateNodes[0]
	} else {
		p.Tail = updateNodes[0]
	}

	p.Length--
	return true
}

// 查找元素
func (p *SkipList) Find(playerId string, target int64) (res *SkipNode, rank int64) {
	pNode := p.Header

	for i := p.Level - 1; i >= 0; i-- {
		for pNode.Levels[i].Next != nil && pNode.Levels[i].Next.Score >= target {
			rank += pNode.Levels[i].Span
			if pNode.Levels[i].Next.PlayerId == playerId &&
				pNode.Levels[i].Next.Score == target {
				res = pNode.Levels[i].Next
				return
			}
			pNode = pNode.Levels[i].Next
		}
	}
	rank = 0
	return
}
func GetRandLevel() int8 {
	level := 1
	for level < MaxLevel && rand.Intn(100) < 50 {
		level++
	}
	return int8(level)
}

// 遍历节点
func (p *SkipList) Traverse() {
	pNode := p.Header
	// fmt.Printf("header:%v\n", pNode)
	fmt.Printf("header:")
	spew.Dump(pNode)
	n := 0
	for pNode.Levels[0].Next != nil {
		pNode = pNode.Levels[0].Next
		fmt.Printf("Top:%d,Member:%v\n", n, pNode)
		n++
	}
}

// 反向遍历
func (p *SkipList) TraverseBack() {
	pNode := p.Tail
	fmt.Printf("Tail:%v\n", pNode)
	for pNode.BackWard != nil {
		pNode = pNode.BackWard
		fmt.Printf("Member:%v\n", pNode)
	}
}