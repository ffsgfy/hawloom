package main

import (
	"time"
)

type timerAcc struct {
	Num   float64 `json:"n"`
	Denom float64 `json:"d"`
}

type Timer map[string]timerAcc

func NewTimer() Timer {
	return map[string]timerAcc{}
}

func (t Timer) add(key string, num, denom float64) {
	acc := t[key]
	acc.Num += num
	acc.Denom += denom
	t[key] = acc
}

type timerHandle struct {
	timer Timer
	start time.Time
	key   string
}

func (t Timer) Start(key string) *timerHandle {
	return &timerHandle{
		timer: t,
		start: time.Now(),
		key:   key,
	}
}

func (th *timerHandle) Stop() {
	dt := time.Now().Sub(th.start).Seconds()
	th.timer.add(th.key, dt, 1.0)
}

func (t Timer) Merge(other Timer) {
	for key, acc := range other {
		t.add(key, acc.Num, acc.Denom)
	}
}

func (t Timer) Drop(keys ...string) {
	for _, key := range keys {
		delete(t, key)
	}
}
