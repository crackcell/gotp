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
 * Errors
 *
 * @file errors.go
 * @author Menglong TAN <tanmenglong@gmail.com>
 * @date Thu Sep  4 23:21:48 2014
 *
 **/

package gotp

import (
	"errors"
)

var (
	ErrInit             = errors.New("init failed")
	ErrNoCallback       = errors.New("no callback")
	ErrStop             = errors.New("server stopped")
	ErrUnknownTag       = errors.New("unknown tag")
	ErrAlreadyRegisterd = errors.New("already registered")
	ErrAlreadyStarted   = errors.New("already started")
	ErrInvalidArgs      = errors.New("invalid arguments")
	ErrUnInit           = errors.New("module uninitialized")
)
