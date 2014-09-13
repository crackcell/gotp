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
 * Unittest for genfsm.
 *
 * @file genfsm_test.go
 * @author Menglong TAN <tanmenglong@gmail.com>
 * @date Fri Sep 12 03:10:51 2014
 *
 **/

package genfsm

import (
	"github.com/crackcell/gotp"
	"log"
	"testing"
)

type dog struct{}

func (this dog) Init(args interface{}) (int, []interface{}) {
	log.Println("dog - init:", args)
	return NextState, gotp.Pack("Sit", 1)
}

// sync handler
func (this dog) Sit(args []interface{}) []interface{} {
	count := data.(int)
	log.Println("dog - pet")
	return NextState, gotp.Pack(count+1, "Stand")
}

var fsm GenFsm
