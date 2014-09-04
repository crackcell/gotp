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
 *
 *
 * @file genserver.go
 * @author Menglong TAN <tanmenglong@gmail.com>
 * @date Thu Sep  4 18:09:12 2014
 *
 **/

package genserver

// Callback interface
type GenServer interface {
	Init(args interface{})

	// msg, state -> reply, state
	HandleCall(msg, state interface{}) (interface{}, interface{})
	// msg, state -> state
	HandleInfo(msg, state interface{}) interface{}
	// msg, state -> state
	HandleCast(msg, state interface{}) interface{}
}
