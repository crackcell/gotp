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

type Msg struct {
	Type  int
	Value interface{}
}

type testServer struct{}

type testState struct {
	loopCount int
}

func (this testServer) Init(args ...interface{}) []interface{} {
	log.Println("[testServer] args:", args)
	return gotp.Pack(Ok, testState{0})
}

func (this testServer) HandleCall(msg, state interface{}) []interface{} {

	s := state.(testState)
	s.loopCount += 1
	m := msg.(Msg)
	log.Printf("[testServer] HandleCall: recv: %s loopCount: %d\n", m.Value, s.loopCount)
	switch m.Type {
	case call1:
		return gotp.Pack(Reply, "reply", s)
	case call2:
		return gotp.Pack(Stop, "call2", s)
	default:
		panic("wrong case")
	}
}

func (this testServer) HandleCast(msg, state interface{}) []interface{} {
	s := state.(testState)
	s.loopCount += 1
	m := msg.(Msg)
	log.Printf("[testServer] HandleCast: recv: %s loopCount: %d\n", m.Value, s.loopCount)
	switch m.Type {
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
	ret := server.Call(Msg{call1, "call - 1"})
	log.Println("[TestCall]", ret)
	ret = server.Call(Msg{call1, "call1 - 2"})
	log.Println("[TestCall]", ret)
}

func TestCast1(t *testing.T) {
	server.Cast(Msg{cast1, "cast1 - 1"})
}

func TestCast2(t *testing.T) {
	server.Cast(Msg{cast2, "cast2 - 2"})
}

/*
func TestCall2(t *testing.T) {
	ret := server.Call(Msg{call2, "call2"})
	log.Println("[TestCast]", ret)
}
*/
