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
	"github.com/crackcell/gotp/genserver"
	"log"
	"sync"
)

type genFsmCallback struct{}

type genFsmState struct {
	callback Callback
	state    interface{}
}

type initArgs struct {
	args interface{}
}

func (this GenFsmCallback) Init(args interface{}) (int, interface{}) {
	log.Println("[GenFsm] init:", args)
	c := args.(Callback)
	tag, nextState, data := c.Init()
	return Ok, genFsmState{callback: args.(Callback)}
}

// GenFsm message tag
const (
	reqSend = 1 << iota
	reqSyncSend
)

type GenFsm struct {
	server GenServer
}

func (this *GenFsm) Start(callback Callback, args interface{}) {

}

func (this *GenFsm) SendEvent() {}

func (this *GenFsm) SyncSendEvent() {}
