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

package genserver

import (
	"errors"
)

var (
	ErrInit       = errors.New("init failed")
	ErrNoCallback = errors.New("no callback")
)
