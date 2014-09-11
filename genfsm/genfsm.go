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
	callback     Callback
	data         interface{}
	state        string
	handlers     map[string]EventHandler
	syncHandlers map[string]SyncEventHandler
}

const (
	reqRegisterHandler = 1 << iota
	reqRegisterSyncHandler
	reqSendEvent
	reqSendSyncEvent
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
	return genserver.Ok, genFsmState{callback: args.(Callback), data: data,
		valid: false, state: nextState, currentHandler: nil, // nil handler, need init
		make(map[string]EventHandler),
		make(map[string]EventHandler)}
}

func (this genFsmCallback) HandleCall(msg, state interface{}) (int, interface{}, interface{}) {
	log.Println("[GenFsm] call")
	m := msg.(genFsmCallbackMsg)
	switch m.tag {
	case reqSendSyncEvent:
		// TODO
	default:
		panic(gotp.ErrUnknownTag)
	}
}

func (this genFsmCallback) HandleCast(msg, state interface{}) (int, interface{}, interface{}) {
	log.Println("[GenFsm] cast")
	m := msg.(genFsmCallbackMsg)
	switch m.tag {
	case reqRegisterHandler:
		// Args:
		//  0 - state name
		//  1 - handler
		if len(m.args) != 2 {
			panic(gotp.ErrInvalidArgs)
		}
		state := m.args[0].(int)
		handler := m.args[1].(EventHandler)
		this.handlers[stateName] = handler
	case reqRegisterSyncHandler:
		// Args:
		//  0 - state name
		//  1 - handler
		if len(m.args) != 2 {
			panic(gotp.ErrInvalidArgs)
		}
		state := m.args[0].(int)
		handler := m.args[1].(EventHandler)
		this.syncHandlers[stateName] = handler
	case reqSendEvent:
		// TODO
	default:
		panic(gotp.ErrUnknownTag)
	}
}

func (this genFsmCallback) Terminate(reason, state interface{}) {
	log.Printf("[GenFsm] Terminate: reason: %s\n", reason)
}

type GenFsm struct {
	server genserver.GenServer
}

func (this *GenFsm) Start(callback Callback, args interface{}) {
	this.server.Start(genFsmCallback, initArgs{args, callback})
}

func (this *GenFsm) SendEvent() {
	this.server.Cast(genFsmCallbackMsg{reqSendEvent})
}

func (this *GenFsm) SyncSendEvent() {}
