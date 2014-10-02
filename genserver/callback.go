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

// Callback interface
type Callback interface {
	// args -> Ok, $NewState
	//      -> Stop, $Reason
	Init(args ...interface{}) []interface{}

	// state, args[] -> Reply, $Reply, $NewState
	//               -> Stop, $Reason, $NewState
	HandleCall(state interface{}, args ...interface{}) []interface{}

	// state, args[] -> Noreply, $NewState
	//               -> Stop, $Reason, $NewState
	HandleCast(state interface{}, args ...interface{}) []interface{}

	// state, args[] -> Reply, $Reply, $NewState
	//               -> Noreply, $NewState
	//               -> Stop, $Reason, $NewState
	HandleInfo(state interface{}, args ...interface{}) []interface{}

	// reason, state
	Terminate(state interface{}, reason interface{})
}
