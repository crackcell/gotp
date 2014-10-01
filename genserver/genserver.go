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
	"sync"
)

const (
	reqCall = 1 << iota
	reqCast
	reqInfo
)

type GenServer struct {
	C         chan []interface{}
	once      sync.Once
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

func (this *GenServer) checkInit() {
	if !this.hasServer {
		panic(gotp.ErrNoCallback)
	}
}

func (this *GenServer) handleReq() {
	//log.Println("handleReq starts")
	for {
		//log.Println("handleReq")
		req := <-this.C
		var tag int
		var reply, state, reason interface{}

		gotp.Assert(len(req) >= 2)

		reqType := req[0].(int)
		reqValue := req[1]
		var reqRet chan []interface{}

		switch reqType {
		case reqCall:
			reqRet = req[2].(chan []interface{})
			// tag, reply, state
			params := this.callback.HandleCall(this.state, reqValue.([]interface{})...)
			if len(params) != 3 {
				panic(gotp.ErrInvalidCallback)
			}
			tag = params[0].(int)
			reply = params[1]
			state = params[2]
		case reqCast:
			// tag, reply, state
			params := this.callback.HandleCast(this.state, reqValue.([]interface{})...)
			switch len(params) {
			case 2: // Noreply, $NewState
				tag = params[0].(int)
				state = params[1]
			case 3: // Stop, $Reason, $NewState
				tag = params[0].(int)
				reason = params[1]
				state = params[2]
			default:
				panic(gotp.ErrInvalidCallback)
			}
		}

		this.state = state
		switch tag {
		case gotp.Reply:
			reqRet <- gotp.Pack(reply)
		case gotp.Noreply: // DO NOTHING
		case gotp.Stop:
			this.callback.Terminate(reason, this.state) // Reason, state
			if reqType == reqCall {
				reqRet <- gotp.Pack()
			}
			//log.Print(gotp.ErrStop, ", reason: ", reason)
			break
		default:
			panic(gotp.ErrUnknownTag)
		}
	}
}

func (this *GenServer) Start(callback Callback, args ...interface{}) {
	this.once.Do(func() {
		this.C = make(chan []interface{})
		this.callback = callback
		this.hasServer = true
		if this.init(args...) {
			go this.handleReq()
		}
	})
}

func (this *GenServer) Call(args ...interface{}) interface{} {
	this.checkInit()
	ret := make(chan []interface{})
	this.C <- gotp.Pack(reqCall, args, ret)
	v := <-ret
	return v[0]
}

func (this *GenServer) Cast(args ...interface{}) {
	this.checkInit()
	this.C <- gotp.Pack(reqCast, args)
}
