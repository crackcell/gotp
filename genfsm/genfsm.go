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
	state        int
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
	tag, params := a.callback.Init(a.args)
	if len(params) < 2 {
		panic(gotp.ErrInvalidArgs)
	}
	switch params[0].(int) {
	case Ok:
		return genserver.Ok, genFsmState{
			callback: args.(Callback),
			data:     data,
			state:    nextState,
			make(map[string]EventHandler),
			make(map[string]EventHandler)}
	case Stop:
		return genserver.Stop, gotp.ErrInit
	default:
		panic(gotp.ErrUnknownTag)
	}
}

func (this genFsmCallback) HandleCall(msg, state interface{}) (int, interface{}, interface{}) {
	log.Println("[GenFsm] call")
	m := msg.(genFsmCallbackMsg)
	s := state.(genFsmState)
	switch m.tag {
	case reqSendSyncEvent:
		handler, err := this.syncHandlers[this.state]
		if err != nil {
			panic(gotp.ErrNoHandler)
		}
		tag, params := handler(msg, this.data)
		if len(params) < 2 {
			panic(gotp.ErrInvalidArgs)
		}
		switch tag {
		case NextState: // params = {reply, next_state, new_data}
			reply := params[0]
			state := params[1].(int)
			log.Printf("send_sync_event: %d -> $d\n", this.state, state)
			s.state = state
			s.data = params[2]
			return genserver.Reply, reply, s
		case Stop:
			return genserver.Stop, gotp.ErrStop, s
		default:
			panic(gotp.ErrUnknownTag)
		}
	default:
		panic(gotp.ErrUnknownTag)
	}
}

func (this genFsmCallback) HandleCast(msg, state interface{}) (int, interface{}, interface{}) {
	log.Println("[GenFsm] cast")
	m := msg.(genFsmCallbackMsg)
	s := state.(genFsmState)
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
		s.handlers[stateName] = handler
	case reqRegisterSyncHandler:
		// Args:
		//  0 - state name
		//  1 - handler
		if len(m.args) != 2 {
			panic(gotp.ErrInvalidArgs)
		}
		state := m.args[0].(int)
		handler := m.args[1].(EventHandler)
		s.handlers[state] = handler
	case reqSendEvent:
		handler, err := this.handlers[this.state]
		if err != nil {
			panic(gotp.ErrNoHandler)
		}
		tag, params := handler(msg, this.data)
		if len(params) != 2 { // not {$NextState, $NewData} or {$Reason, $NewData}
			panic(gotp.ErrInvalidArgs)
		}
		switch tag {
		case NextState:
			state := params[0].(int)
			log.Printf("send_event: %d -> %d\n", this.state, state)
			s.state = state
			s.data = params[1]
			return genserver.Noreply, nil, s
		case Stop:
			return genserver.Stop, gotp.ErrStop, s
		default:
			panic(gotp.ErrUnknownTag)
		}
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

func (this *GenFsm) RegisterHandler()
