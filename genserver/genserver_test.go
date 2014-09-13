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
	"github.com/crackcell/gotp"
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

type testServer struct{}

type testState struct {
	loopCount int
}

func (this testServer) Init(args ...interface{}) []interface{} {
	log.Println("[testServer] args:", args)
	return gotp.Pack(Ok, testState{0})
}

func (this testServer) HandleCall(state interface{}, args ...interface{}) []interface{} {
	s := state.(testState)
	s.loopCount += 1
	log.Printf("[testServer] HandleCall: recv: %s loopCount: %d\n", args[1], s.loopCount)
	switch args[0].(int) {
	case call1:
		return gotp.Pack(Reply, "reply", s)
	case call2:
		return gotp.Pack(Stop, "call2", s)
	default:
		panic("wrong case")
	}
}

func (this testServer) HandleCast(state interface{}, args ...interface{}) []interface{} {
	s := state.(testState)
	s.loopCount += 1
	log.Printf("[testServer] HandleCast: recv: %s loopCount: %d\n", args[1], s.loopCount)
	switch args[0].(int) {
	case cast1:
		return gotp.Pack(Noreply, s)
	case cast2:
		return gotp.Pack(Stop, "cast2", s)
	default:
		panic("wrong case")
	}
}

func (this testServer) Terminate(reason, state interface{}) {
	log.Printf("[testServer] Terminate: reason: %s\n", reason)
}

var server GenServer

func TestStart(t *testing.T) {
	server.Start(testServer{}, "args")
}

func TestCall1(t *testing.T) {
	time.Sleep(2000)
	ret := server.Call(call1, "call - 1")
	log.Println("[TestCall]", ret)
	ret = server.Call(call1, "call1 - 2")
	log.Println("[TestCall]", ret)
}

func TestCast1(t *testing.T) {
	server.Cast(cast1, "cast1 - 1")
}

func TestCast2(t *testing.T) {
	server.Cast(cast2, "cast2 - 2")
}

/*
func TestCall2(t *testing.T) {
	ret := server.Call(call2, "call2")
	log.Println("[TestCast]", ret)
}
*/
