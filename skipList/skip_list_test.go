package main

import (
	"fmt"
	"testing"
)

func TestSkipList(t *testing.T) {
	sl := NewSkipList()
	sl.InitSkipList()
	for i := 1; i <= 8; i++ {
		sl.Insert("1", int64(i))
	}
	sl.Traverse()
	res, rank := sl.Find("1", 5)
	if rank != 4 {
		t.Fatal("rank error:exprect 4,actual ", rank)
	}
	fmt.Println(res, rank)
	sl.Delete("1", 5)
	res, rank = sl.Find("1", 5)
	if rank != 0 {
		t.Fatal("rank error:exprect 0,actual ", rank)
	}
	fmt.Println("res:", res, "rank", rank)
	res, rank = sl.Find("1", 4)
	if rank != 4 {
		t.Fatal("rank error:exprect 4,actual ", rank)
	}
	fmt.Println("res:", res, "rank", rank)
}

func TestTraverseBack(t *testing.T) {
	sl := NewSkipList()
	sl.InitSkipList()
	for i := 1; i <= 10; i++ {
		sl.Insert("1", int64(i))
	}
	sl.TraverseBack()
	fmt.Println(1)
}
