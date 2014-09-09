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
	"sync"
)

// Message tag
const (
	ok = 1 << iota
	reply
	noreply
	nextState
	stop
)

type GenFsm interface {
	// args -> Ok, $InitStateName, $InitData
	//      -> Stop, $Reason
	Init(args interface{}) (bool, string, interface{})

	// msg, state_name, data -> NextState, $NextStateName, $NewData
	//                       -> Stop, $Reason
	HandleEvent(msg interface{}, statename string, data interface{}) (int, string, interface{})
	HandleAllEvent(msg interface{})
}

func (this *GenFsm) SendEvent() {}

func (this *GenFsm) SyncSendEvent() {}
