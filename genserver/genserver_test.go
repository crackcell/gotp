/* -*- encoding: utf-8; indent-tabs-mode: t -*- */

/***************************************************************
 *
 * Copyright (c) 2014, Menglong TAN <tanmenglong@gmail.com>
 *
 * This program is free software; you can redistribute it
 * and/or modify it under the terms of the GPL licence
 *
 **************************************************************/

/**
 * Unittest for genserver.
 *
 * @file genserver_test.go
 * @author Menglong TAN <tanmenglong@gmail.com>
 * @date Thu Sep  4 21:34:15 2014
 *
 **/

package genserver

import (
	"log"
	"testing"
	"time"
)

// Message tage
const (
	call1 = 1 << iota
	call2
	cast1
	cast2
)

type Msg struct {
	Type  int
	Value interface{}
}

type TestServer struct {
}

type testState struct {
	loopCount int
}

func (this TestServer) Init(args interface{}) (int, interface{}) {
	log.Println("[TestServer] args:", args)
	return Ok, testState{0}
}

func (this TestServer) HandleCall(msg, state interface{}) (int, interface{}, interface{}) {
	s := state.(testState)
	s.loopCount += 1
	m := msg.(Msg)
	log.Printf("[TestServer] HandleCall: recv: %s loopCount: %d\n", m.Value, s.loopCount)
	switch m.Type {
	case call1:
		return Reply, "reply", s
	case call2:
		return Stop, "call2", s
	default:
		panic("wrong case")
	}
}

func (this TestServer) HandleCast(msg, state interface{}) (int, interface{}, interface{}) {
	s := state.(testState)
	s.loopCount += 1
	m := msg.(Msg)
	log.Printf("[TestServer] HandleCast: recv: %s loopCount: %d\n", m.Value, s.loopCount)
	switch m.Type {
	case cast1:
		return Noreply, nil, s
	case cast2:
		return Stop, "cast2", s
	default:
		panic("wrong case")
	}
}

func (this TestServer) Terminate(reason, state interface{}) {
	log.Printf("[TestServer] Terminate: reason: %s\n", reason)
}

func TestStart(t *testing.T) {
	Start("TestServer", TestServer{}, "args")
}

func TestCall1(t *testing.T) {
	time.Sleep(2000)
	ret := Call("TestServer", Msg{call1, "call - 1"})
	log.Println("[TestCall]", ret)
	ret = Call("TestServer", Msg{call1, "call1 - 2"})
	log.Println("[TestCall]", ret)
}

func TestCast1(t *testing.T) {
	Cast("TestServer", Msg{cast1, "cast1 - 1"})
}

/*
func TestCast2(t *testing.T) {
	Cast("TestServer", Msg{cast2, "cast2 - 2"})
}
*/

func TestCall2(t *testing.T) {
	ret := Call("TestServer", Msg{call2, "call2"})
	log.Println("[TestCast]", ret)
}
