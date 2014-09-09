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
 * Generic finite state machine framework inspired by Erlang/OTP.
 *
 * @file genfsm.go
 * @author Menglong TAN <tanmenglong@gmail.com>
 * @date Thu Sep  4 23:35:57 2014
 *
 **/

package genfsm

import (
	"github.com/crackcell/gotp"
	"github.com/crackcell/gotp/genserver"
	"log"
	"sync"
)

type genFsmCallback struct{}

type genFsmState struct {
	callback Callback
	state    interface{}
	data     interface{}
}

const (
	reqTransferState = 1 << iota
)

type genFsmCallbackMsg struct {
	tag  int
	args []interface{}
}

type initArgs struct {
	args     interface{}
	callback Callback
}

func (this GenFsmCallback) Init(args interface{}) (int, interface{}) {
	log.Println("[GenFsm] init:", args)
	a := args.(initArgs)
	tag, nextState, data := a.callback.Init(a.args)
	return Ok, genFsmState{callback: args.(Callback), state: nextState, data: data}
}

func (this GenFsmCallback) HandleCall(msg, state interface{}) (int, interface{}, interface{}) {
	log.Println("[GenFsm] call")
	m := msg.(genFsmCallbackMsg)
	switch m.tag {
	case reqTransferState:
	default:
		panic(gotp.ErrUnknownTag)
	}
}

// GenFsm message tag
const (
	reqSend = 1 << iota
	reqSyncSend
)

type GenFsm struct {
	server GenServer
}

func (this *GenFsm) Start(callback Callback, args interface{}) {

}

func (this *GenFsm) SendEvent() {}

func (this *GenFsm) SyncSendEvent() {}
