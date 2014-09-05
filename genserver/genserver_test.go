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
	log.Printf("[TestServer] HandleCall: recv: %s loopCount: %d\n", msg, s.loopCount)
	return Reply, "reply", s
}

func (this TestServer) HandleCast(msg, state interface{}) (int, interface{}) {
	s := state.(testState)
	s.loopCount += 1
	log.Printf("[TestServer] HandleCast: recv: %s loopCount: %d\n", msg, s.loopCount)
	return Noreply, s
}

func TestStart(t *testing.T) {
	Start("TestServer", TestServer{}, "args")
}

func TestCall(t *testing.T) {
	time.Sleep(2000)
	ret := Call("TestServer", "call args")
	log.Println("[TestCall]", ret)
	ret = Call("TestServer", "call args")
	log.Println("[TestCall]", ret)
}

func TestCast(t *testing.T) {
	Cast("TestServer", "cast args")
	log.Println("[TestCast]")
}
