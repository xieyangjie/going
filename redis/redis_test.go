package redis

import (
	"testing"
	"fmt"
)

func TestRedis(t *testing.T) {
	var s = NewRedis("127.0.0.1:6379", "", 1, 30, 10)
	var c = s.GetSession()

	var p1 = Plan{Title:"pt1", Text:"t1"}
	c.HMSET("p1", &p1)

	var p2 Plan

	fmt.Println(c.HGETALL("p1", &p2))
	fmt.Println(p2)
}

type Plan struct {
	Title string `redis:"title"`
	Text  string `redis:"text"`
}