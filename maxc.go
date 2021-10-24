//go:build !tuneparam
// +build !tuneparam

/*	Copyright (c) 2019, Serhat Şevki Dinçer.
	This Source Code Form is subject to the terms of the Mozilla Public
	License, v. 2.0. If a copy of the MPL was not distributed with this
	file, You can obtain one at http://mozilla.org/MPL/2.0/.
*/

package sorty

const (
	// MaxLenIns is the default maximum slice length for insertion sort.
	MaxLenIns = 100

	// MaxLenInsFC is the maximum slice length for insertion sort when
	// sorting strings or calling Sort().
	MaxLenInsFC = 40

	// MaxLenRec is the maximum slice length for recursion when there is goroutine
	// quota. So MaxLenRec+1 is the minimum slice length for new sorting goroutines.
	MaxLenRec = 496
)
