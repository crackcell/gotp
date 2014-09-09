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
 * Name server.
 *
 * @file nameserver.go
 * @author Menglong TAN <tanmenglong@gmail.com>
 * @date Fri Sep  5 15:48:57 2014
 *
 **/

package naming

import (
	"github.com/crackcell/gotp"
	"github.com/crackcell/gotp/genserver"
	"log"
	"sync"
)

// Request type
const (
	reqRegisterProc = 1 << iota
)

type nameServer struct {
	ch      chan gotp.Req
	once    sync.Once
	servers map[int]*genserver.GenServer
}

func (this *nameServer) init() {
	this.once.Do(func() {
		this.ch = make(chan gotp.Req)
		this.servers = make(make[string] * genserver.GenServer)
		go this.handleReq()
	})
}

func (this *nameServer) handleReq() {
	for {
		req := <-this.ch
		switch req.Type {
		case reqRegister:
		default:
			panic(gotp.errUnknownTag)
		}
	}
}
