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
type GenServer interface {
	// args -> tag, state
	Init(args interface{}) (int, interface{})
	// msg, state -> tag, reply, state
	HandleCall(msg, state interface{}) (int, interface{}, interface{})
	// msg, state -> tag, state
	HandleCast(msg, state interface{}) (int, interface{})
}

const (
	reqInit = 1 << iota
	reqCall
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
	callback  GenServer
	state     interface{}
}

func (this *wrapper) handleReq() {
	//log.Println("handleReq starts")
	for {
		//log.Println("handleReq")
		req := <-this.ch
		switch req.typ {
		case reqInit:
			tag, state := this.callback.Init(req.value)
			log.Println("init")
			switch tag {
			case Ok:
				this.state = state
			case Stop:
				log.Fatal(errInit)
				break
			default:
				panic(errUnknownTag)
			}
		case reqCall:
			tag, reply, state := this.callback.HandleCall(req.value, this.state)
			switch tag {
			case Reply:
				req.ret <- wrapperResp{reply, nil}
				this.state = state
			case Stop:
				log.Fatal(errStop)
				break
			default:
				panic(errUnknownTag)
			}
		case reqCast:
			tag, state := this.callback.HandleCast(req.value, this.state)
			switch tag {
			case Noreply:
				this.state = state
			case Stop:
				log.Fatal(errStop)
				break
			default:
				panic(errUnknownTag)
			}
		}
	}
}

func (this *wrapper) start(name string, callback GenServer, args interface{}) {
	this.once.Do(func() {
		this.ch = make(chan wrapperReq)
		this.callback = callback
		this.hasServer = true
		go this.handleReq()
		log.Println("start")
		this.ch <- wrapperReq{reqInit, args, nil}
	})
}

func (this *wrapper) checkInit() {
	if !this.hasServer {
		panic(errNoCallback)
	}
}

func (this *wrapper) call(msg interface{}) interface{} {
	this.checkInit()
	ret := make(chan wrapperResp)
	this.ch <- wrapperReq{reqCall, msg, ret}
	v := <-ret
	return v.value
}

func (this *wrapper) cast(msg interface{}) {
	this.checkInit()
	this.ch <- wrapperReq{reqCast, msg, nil}
}

var w wrapper

func Start(name string, callback GenServer, args interface{}) {
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
