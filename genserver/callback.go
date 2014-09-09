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
 * @date Tue Sep  9 16:12:29 2014
 *
 **/

package genserver

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
