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

// $Msg, $Data -> NextState, {$NextState, $NewData}
//             -> Stop, {$Reason, $NewData}
type EventHandler func(msg, data interface{}) (int, []interface{})

// $Msg, $Data -> NextState, {$Reply, $NextState, $NewData}
//             -> Stop, {$Reason, $NewData}
type SyncEventHandler func(msg, data interface{}) (int, []interface{})

type Callback interface {
	// args -> Ok, {$NextState, $InitData}
	//      -> Stop, {$Reason}
	Init(args interface{}) (int, []interface{})
}
