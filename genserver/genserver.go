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
 * Generic server framework inspired by Erlang/OTP.
 *
 * @file genserver.go
 * @author Menglong TAN <tanmenglong@gmail.com>
 * @date Thu Sep  4 17:54:31 2014
 *
 **/

package genserver

import (
	"github.com/crackcell/gotp"
	"log"
)

const (
	ReqCall = 1 << iota
	ReqCast
)

type GenServer struct {
	C         chan []interface{}
	start     bool
	hasServer bool
	callback  Callback
	state     interface{}
}

func (this *GenServer) init(args ...interface{}) bool {
	params := this.callback.Init(args...)
	gotp.AssertArrayArity(params, 2)
	tag := params[0].(int)
	state := params[1]
	switch tag {
	case gotp.Ok:
		this.state = state
	case gotp.Stop:
		log.Fatal(gotp.ErrInit)
		return false
	default:
		panic(gotp.ErrUnknownTag)
	}
	return true
}

func (this *GenServer) checkCallback() {
	if !this.hasServer {
		panic(gotp.ErrNoCallback)
	}
}

func (this *GenServer) handleReq() {
	for {
		req := <-this.C

		var reqType int
		var reqValue interface{}
		var reqRet chan []interface{}
		var callbackRet []interface{}
		var ok bool

		// Check message
		if len(req) < 2 {
			panic(gotp.ErrUnknownTag)
		}
		if reqType, ok = req[0].(int); !ok {
			panic(gotp.ErrUnknownTag)
		}
		reqValue = req[1]

		switch reqType {
		case ReqCall:
			gotp.Assert(len(req) == 3)
			reqRet = req[2].(chan []interface{})
			callbackRet = this.callback.HandleCall(this.state, reqValue.([]interface{})...)
		case ReqCast:
			gotp.Assert(len(req) == 2)
			callbackRet = this.callback.HandleCast(this.state, reqValue.([]interface{})...)
		//case ReqInfo:
		//	gotp.Assert(len(req) == 2) // TODO: support HandleInfo with Reply
		//	callbackRet = this.callback.HandleInfo(this.state, reqValue.([]interface{})...)
		default:
			panic(gotp.ErrUnknownTag)
		}

		var tag int
		var reply, state, reason interface{}

		tag = callbackRet[0].(int)
		switch tag {
		case gotp.Reply: // [Reply, $Reply, $NewState]
			gotp.Assert(len(callbackRet) == 3)
			//gotp.Assert(reqType != ReqInfo) // TODP: support HandleInfo with Reply
			reply = callbackRet[1]
			state = callbackRet[2]
			reqRet <- gotp.Pack(gotp.Reply, reply)
		case gotp.Noreply: // [Noreply, $NewState]
			gotp.Assert(len(callbackRet) == 2)
			state = callbackRet[1]
		case gotp.Stop: // [Stop, $Reason, $NewState]
			gotp.Assert(len(callbackRet) == 3)
			reason = callbackRet[1]
			state = callbackRet[2]
			this.callback.Terminate(reason, state) // $Reason, $NewState
			if reqType == ReqCall {
				reqRet <- gotp.Pack(gotp.Stop, reason)
			}
			break
		default:
			panic(gotp.ErrInvalidCallback)
		}

		this.state = state
	}
}

func (this *GenServer) Start(callback Callback, args ...interface{}) error {
	if this.start {
		return gotp.ErrAlreadyStarted
	}
	this.C = make(chan []interface{})
	this.callback = callback
	this.hasServer = true
	this.start = true
	if this.init(args...) {
		go this.handleReq()
	}
	return nil
}

func (this *GenServer) Call(args ...interface{}) (interface{}, error) {
	this.checkCallback()
	return SendCall(this.C, args...)
}

func (this *GenServer) Cast(args ...interface{}) {
	this.checkCallback()
	SendCast(this.C, args...)
}

//func (this *GenServer) Info(args ...interface{}) {
//	this.checkCallback()
//	this.C <- gotp.Pack(ReqInfo, args)
//}

func SendCall(c chan []interface{}, args ...interface{}) (interface{}, error) {
	ch := make(chan []interface{})
	c <- gotp.Pack(ReqCall, args, ch)
	ret := <-ch
	gotp.Assert(len(ret) == 2) // [Reply, $Reply] or [Stop, $Reason]
	switch ret[0].(int) {
	case gotp.Reply:
		return ret[1], nil
	case gotp.Stop:
		return ret[1], gotp.ErrStop
	default:
		return ret[1], gotp.ErrUnknownTag
	}
}

func SendCast(c chan []interface{}, args ...interface{}) {
	c <- gotp.Pack(ReqCast, args)
}

//func SendInfo(c chan []interface{}, args ...interface{}) {
//	c <- gotp.Pack(ReqInfo, args)
//}
