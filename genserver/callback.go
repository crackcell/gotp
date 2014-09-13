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
	// args -> Ok, $NewState
	//      -> Stop, $Reason
	Init(args ...interface{}) []interface{}

	// msg, state -> Reply, $Reply, $NewState
	//            -> Stop, $Reason, $NewState
	HandleCall(msg, state interface{}) []interface{}

	// msg, state -> Noreply, $NewState
	//            -> Stop, $Reason, $NewState
	HandleCast(msg, state interface{}) []interface{}

	// reason, state
	Terminate(reason, state interface{})
}
