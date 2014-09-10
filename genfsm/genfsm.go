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

func (this genFsmCallback) Init(args interface{}) (int, interface{}) {
	log.Println("[GenFsm] init:", args)
	a := args.(initArgs)
	tag, nextState, data := a.callback.Init(a.args)
	return genserver.Ok, genFsmState{args.(Callback), nextState, data}
}

func (this genFsmCallback) HandleCall(msg, state interface{}) (int, interface{}, interface{}) {
	log.Println("[GenFsm] call")
	m := msg.(genFsmCallbackMsg)
	switch m.tag {
	case reqTransferState:
	default:
		panic(gotp.ErrUnknownTag)
	}
}

func (this genFsmCallback) HandleCast(msg, state interface{}) (int, interface{}, interface{}) {
	log.Println("[GenFsm] cast")
	m := msg.(genFsmCallbackMsg)
	switch m.tag {
	case reqTransferState:
	default:
		panic(gotp.ErrUnknownTag)
	}
}

func (this genFsmCallback) Terminate(reason, state interface{}) {
	log.Printf("[GenFsm] Terminate: reason: %s\n", reason)
}

// GenFsm message tag
const (
	reqSend = 1 << iota
	reqSyncSend
)

type GenFsm struct {
	server genserver.GenServer
}

func (this *GenFsm) Start(callback Callback, args interface{}) {
	this.server.Start(genFsmCallback, initArgs{args, callback})
}

func (this *GenFsm) SendEvent() {}

func (this *GenFsm) SyncSendEvent() {}
