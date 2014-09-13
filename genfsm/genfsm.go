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

//type genFsmCallbackMsg struct {
//	tag  int
//	args []interface{}
//}

func (this genFsmCallback) Init(args ...interface{}) []interface{} {
	// args:
	//   0 - args []interface{}
	//   1 - callback
	log.Println("[GenFsm] init:", args)
	gotp.AssertArrayArity(args, 2)
	callback := args[1].(Callback)
	params := callback.Init(args[0])
	if len(params) < 2 {
		panic(gotp.ErrInvalidCallback)
	}
	tag := params[0].(int)
	switch tag {
	case Ok: // params: [Ok, $NextState, $InitData]
		return gotp.Pack(genserver.Ok,
			genFsmState{
				callback: callback,
				data:     params[1],
				state:    params[0].(string),
			})
	case Stop:
		return gotp.Pack(genserver.Stop, gotp.ErrInit)
	default:
		panic(gotp.ErrUnknownTag)
	}
}

func (this genFsmCallback) HandleCall(msg, state interface{}) []interface{} {
	log.Println("[GenFsm] call")
	m := msg.([]interface{})
	s := state.(genFsmState)
	if len(m) < 3 {
		panic(gotp.ErrInvalidArgs)
	}
	switch m[0].(int) {
	case reqSendSyncEvent: // [reqSendEvent, state, args[] ]
		// args:
		//   0 - tag
		//   1 - state
		//   2 - args[]
		gotp.AssertArrayArity(m, 3)
		stateName := m[1].(string)
		args := m[2].([]interface{})
		ret := reflect.ValueOf(s.callback).MethodByName(stateName).Call([]reflect.Value{args})
		if len(ret) < 2 {
			panic(gotp.ErrInvalidCallback)
		}
		// Callback:
		// $Args -> [NextState, $Reply, $NextState, $NewData]
		//       -> [Stop, $Reason, $NewData]
		switch ret[0] {
		case NextState:
			log.Printf("send_sysnc_event: %s -> %s\n", s.state, ret[2])
			s.state = ret[2]
			s.data = ret[3]
			return gotp.Pack(genserver.Reply, ret[1], s)
		case Stop:
			log.Println("stop")
			return gotp.Pack(genserver.Stop, ret[1], s)
		default:
			panic(gotp.ErrInvalidCallback)
		}
	default:
		panic(gotp.ErrUnknownTag)
	}
}

func (this genFsmCallback) HandleCast(msg, state interface{}) []interface{} {
	log.Println("[GenFsm] cast")
	m := msg.([]interface{})
	s := state.(genFsmState)
	if len(m) < 2 {
		panic(gotp.ErrInvalidArgs)
	}
	switch m[0].(int) {
	case reqSendEvent:
		// args:
		//   0 - tag
		//   1 - state
		//   2 - args[]
		if len(m) < 2 {
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
			return gotp.Pack(genserver.Noreply, s)
		case Stop: // Stop, $Reason, $NewData
			gotp.AssertArrayArity(values, 2)
			return gotp.Pack(genserver.Stop, values[1].Interface(), s)
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
	this.server.Start(genFsmCallback{}, []interface{}{args, callback})
}

func (this *GenFsm) SendEvent(state string, args ...interface{}) {
	this.server.Cast(reqSendEvent, state, args)
}

func (this *GenFsm) SyncSendEvent(state string, args ...interface{}) interface{} {
	return this.server.Call(reqSendSyncEvent, state, args)
}
