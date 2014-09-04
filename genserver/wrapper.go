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
 * Message process framework like gen_server from Erlang/OTP
 *
 * @file genserver.go
 * @author Menglong TAN <tanmenglong@gmail.com>
 * @date Thu Sep  4 17:54:31 2014
 *
 **/

package genserver

import (
	"sync"
)

// Message type
const (
	signalMessage = 1 << iota
	broadcastMessage
)

const (
	reqCall = 1 << iota
	reqInfo
	reqCast
)

type wrapperReq struct {
	typ   int
	value interface{}
	ret   chan wrapperResp
}

type wrapperResp struct {
	value interface{}
	err   error
}

type wrapper struct {
	ch        chan wrapperReq
	once      sync.Once
	hasServer bool
	server    GenServer
	state     interface{}
}

func (this *wrapper) init() {
	this.once.Do(func() {
		this.ch = make(chan wrapperReq)
		this.servers = []server{}
	})
}

func (this *wrapper) handleReq() {
	for {
		req := <-this.ch
		switch req.typ {
		case reqCall:
			reply, state := serv.HandleCall(req.value, this.state)
			req.ret <- reply
			this.state = state
		case reqInfo:
			this.state = serv.HandleInfo(req.value, this.state)
		case reqCast:
			this.state = serv.HandleCast(req.value, this.state)
		}
	}
}

func (this *wrapper) newServer(callback GenServer) {
	this.once.Do(func() {
		this.ch = make(chan wrapperReq)
		this.server = callback
		this.hasServer = true
	})
}

func (this *wrapper) checkInit() {
	if this.hasServer {
		panic("no callback")
	}
}

func (this *wrapper) call(msg interface{}) interface{} {
	this.checkInit()
	ret := make(chan interface{})
	this.ch <- wrapperReq{reqCall, msg, ret}
	v := <-ret
	return v
}

func (this *wrapper) info(name string, msg interface{}) {
	this.checkInit()
	this.ch <- wrapperReq{reqInfo, message{name, msg}, nil}
}

func (this *wrapper) cast(name string, msg interface{}) {
	this.checkInit()
	this.ch <- wrapperReq{reqCast, message{name, msg}, nil}
}
