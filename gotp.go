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
 * @file gotp.go
 * @author Menglong TAN <tanmenglong@gmail.com>
 * @date Fri Sep  5 16:07:59 2014
 *
 **/

package gotp

// Message tag
const (
	Ok = 1 << iota
	Reply
	Noreply // for cast
	Stop
)

func Pack(args ...interface{}) []interface{} {
	return args
}

func Assert(cond bool) {
	if !cond {
		panic("Assertion failed")
	}
}

func AssertArrayArity(array []interface{}, arity int) {
	if len(array) != arity {
		panic("Bad arity")
	}
}
