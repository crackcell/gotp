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
 * Callbacks.
 *
 * @file callback.go
 * @author Menglong TAN <tanmenglong@gmail.com>
 * @date Tue Sep  9 16:11:45 2014
 *
 **/

package genfsm

// Message tag
const (
	ok = 1 << iota
	reply
	noreply
	nextState
	stop
)

type Callback interface {
	// args -> Ok, $NextStateName, $InitData
	//      -> Stop, nil, $Reason
	Init(args interface{}) (int, interface{}, interface{})

	// msg, state_name, data -> NextState, $NextStateName, $NewData
	//                       -> Stop, $Reason, $NewData
	HandleEvent(msg interface{}, state interface{}, data interface{}) (int, interface{}, interface{})
	HandleAllEvent(msg interface{})
}
