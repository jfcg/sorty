// +build !tuneparam

/*	Copyright (c) 2019, Serhat Şevki Dinçer.
	This Source Code Form is subject to the terms of the Mozilla Public
	License, v. 2.0. If a copy of the MPL was not distributed with this
	file, You can obtain one at http://mozilla.org/MPL/2.0/.
*/

package sorty

const (
	// Mli is the maximum slice length for insertion sort in
	// Sort*() except SortS(), SortB() and Sort().
	Mli = 100

	// Hmli is the maximum slice length for insertion sort in
	// SortS(), SortB() and Sort().
	Hmli = 40

	// Mlr is the maximum slice length for recursion when there are available
	// goroutines. So Mlr+1 is the minimum slice length for new sorting goroutines.
	Mlr = 496
)
