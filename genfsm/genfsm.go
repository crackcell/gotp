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
	// args -> succ, init_state_name, init_data
	Init(args interface{}) (bool, string, interface{})

	// msg, state_name, data -> action, next_state_name, data
	HandleEvent()
}
