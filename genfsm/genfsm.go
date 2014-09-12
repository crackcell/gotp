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
	"reflect"
)

type genFsmCallback struct{}

type genFsmState struct {
	callback Callback
	data     interface{}
	state    string
}

const (
	reqSendEvent = 1 << iota
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
	if len(params) == 0 {
		panic(gotp.ErrInvalidArgs)
	}
	switch tag {
	case Ok: // Ok, {$NextState, $InitData}
		return genserver.Ok, genFsmState{
			callback: args.(Callback),
			data:     params[1],
			state:    params[0].(string)}
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
		// args:
		//   1 - state
		//   2 - args
		if len(m.args) != 1 {
			panic(gotp.ErrInvalidArgs)
		}
		stateName := m.args[0].(string)
		values := reflect.ValueOf(s.callback).MethodByName(stateName).Call([]reflect.Value{m.args[1]})

		handler, ok := s.syncHandlers[s.state]
		if !ok {
			panic(gotp.ErrNoHandler)
		}

		tag, params := handler(msg, s.data)
		if len(params) < 2 {
			panic(gotp.ErrInvalidArgs)
		}
		switch tag {
		case NextState: // params = {reply, next_state, new_data}
			reply := params[0]
			state := params[1].(string)
			log.Printf("send_sync_event: %d -> $d\n", s.state, state)
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
	case reqSendEvent:
		// args:
		//   0 - state
		//   1 - args
		if len(m.args) != 2 {
			panic(gotp.ErrInvalidArgs)
		}
		stateName := m.args[0].(string)
		values := reflect.ValueOf(s.callback).MethodByName(stateName).Call([]reflect.Value{m.args[1]})
		tag := values[0].Int()
		switch tag {
		case NextState: // NextState, $NextState, $NewData
			gotp.AssertArrayArity(values, 3)
			s.state = values[1].String()
			s.data = values[2].Interface()
			return genserver.Noreply, nil, s
		case Stop: // Stop, $Reason, $NewData
			gotp.AssertArrayArity(values, 2)
			return genserver.Stop, values[1].Interface(), s
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
	this.server.Start(genFsmCallback{}, initArgs{args, callback})
}

func (this *GenFsm) SendEvent(state string) {
	this.server.Cast(genFsmCallbackMsg{reqSendEvent, gotp.Pack(state)})
}

func (this *GenFsm) SyncSendEvent(state string) interface{} {
	return this.server.Call(genFsmCallbackMsg{reqSendSyncEvent, gotp.Pack(state)})
}
