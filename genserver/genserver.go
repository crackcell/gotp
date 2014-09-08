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
	"github.com/crackcell/goserv"
	"log"
	"sync"
)

// Message tag
const (
	Ok = 1 << iota
	Reply
	Noreply // for cast
	Stop
)

// Callback interface
type Callback interface {
	// args -> Ok, $State
	//      -> Stop, $Reason
	Init(args interface{}) (int, interface{})
	// msg, state -> Reply, $Reply, $State
	//            -> Stop, $Reason, $State
	HandleCall(msg, state interface{}) (int, interface{}, interface{})
	// msg, state -> Reply, nil, $State
	//            -> Stop, $Reason, $State
	HandleCast(msg, state interface{}) (int, interface{}, interface{})
	// reason, state
	Terminate(reason, state interface{})
}

const (
	reqCall = 1 << iota
	reqCast
)

type GenServer struct {
	ch        chan goserv.Req
	once      sync.Once
	hasServer bool
	callback  Callback
	state     interface{}
}

func (this *GenServer) handleReq() {
	//log.Println("handleReq starts")
	for {
		//log.Println("handleReq")
		req := <-this.ch
		var tag int
		var reply, state interface{}

		switch req.Type {
		case reqCall:
			tag, reply, state = this.callback.HandleCall(req.Value, this.state)
		case reqCast:
			tag, reply, state = this.callback.HandleCast(req.Value, this.state)
		}

		this.state = state
		switch tag {
		case Reply:
			req.Ret <- goserv.Resp{reply, nil}
		case Noreply: // DO NOTHING
		case Stop:
			this.callback.Terminate(reply, this.state) // Reason, state
			if req.Type == reqCall {
				req.Ret <- goserv.Resp{nil, nil}
			}
			log.Print(goserv.ErrStop, ", reason: ", reply)
			break
		default:
			panic(goserv.ErrUnknownTag)
		}
	}
}

func (this *GenServer) init(args interface{}) bool {
	tag, state := this.callback.Init(args)
	switch tag {
	case Ok:
		this.state = state
	case Stop:
		log.Fatal(goserv.ErrInit)
		return false
	default:
		panic(goserv.ErrUnknownTag)
	}
	return true
}

func (this *GenServer) start(name string, callback Callback, args interface{}) {
	this.once.Do(func() {
		this.ch = make(chan goserv.Req)
		this.callback = callback
		this.hasServer = true
		if this.init(args) {
			go this.handleReq()
		}
	})
}

func (this *GenServer) checkInit() {
	if !this.hasServer {
		panic(goserv.ErrNoCallback)
	}
}

func (this *GenServer) call(msg interface{}) interface{} {
	this.checkInit()
	ret := make(chan goserv.Resp)
	this.ch <- goserv.Req{reqCall, msg, ret}
	v := <-ret
	return v.Value
}

func (this *GenServer) cast(msg interface{}) {
	this.checkInit()
	this.ch <- goserv.Req{reqCast, msg, nil}
}

var w GenServer

func Start(name string, callback Callback, args interface{}) {
	w.start(name, callback, args)
}

func Call(name string, req interface{}) interface{} {
	ret := w.call(req)
	//log.Println("call: ret:", ret)
	return ret
}

func Cast(name string, req interface{}) {
	w.cast(req)
}
