/*
 *  Copyright © 2013. Andy Newman.
 */

// Package par provides functions to structure concurrent programs.
//
// The package defines the functions par.DO and par.FOR that implement
// synchronized "processes" (goroutines) in a manner similar to the
// PAR and replicated-PAR control structures found in the occam
// programming language.
//
// Synchronizing upon the complementation of groups of goroutines is a
// common process structure used in many concurrent programs. And in
// Go the structure is idiomatically implemented via sync.WaitGroup
// (and the fact the idiom exists demonstrates its commonality).
//
// The par.DO function mimics occam's PAR and runs some number of
// functions concurrently then waits for them to complete before it
// oompletes and returns to the caller.
//
// par.DO implements the CSP _PAR_ construct, part of what some people
// are calling _structured concurrency_.
//
// In occam, and CSP, a PAR applies to statements. In Go functions,
// closures, are used. A par.DO calls zero or more functions with each
// call within a separate goroutine.
//
// The par.FOR function is a concurrent for-loop and mimics occam's
// replicated-PAR statement. par.FOR calls a single function N times,
// concurrently, where N is defined by two integer _control_ values, a
// start value and a limit value. par.FOR iterates over the range
// defined by these values, using a step of 1, and calls the supplied
// function with the current loop index as its _identifier_.  Each
// call occuring within a separate goroutine. Like par.DO the call
// to par.FOR only returns when all goroutines complete.
//
// Implementation
//
// par.DO and par.FOR are implemented sync.WaitGroup and consolidate
// WaitGroup manipulation within the package which helps remove repetition
// from user code. Subjectively the functions also improve readability
// and remove clutter caused by the required WaitGroup manipulations
// which can often obscure the user's actual code.
//
// Usage
//
// 	par.DO(
//	    controlFuelRods,
//	    monitorCoolant,
//	    moveDials,
//	    flashLights,
//	    func() {
//		par.FOR(0, 10, func(i int) {
// 		    ...  run generator i
// 		})
//          },
//      )
//
// Each thing par.DO runs in parallel is a func().  In the above example
// just assume there's a func controlFueldRods(), func monitorCoolant()
// and so on (the example shows the use of named functions in place of
// literal closures to clarify program structure).
//
// In code the use of par.DO and par.FOR will read quite well. This is
// important, the concurrent structure of a program can be a highly
// important aspect of its design and often is not well communicated
// between developers. I'm hoping these little functions an help.
//
// PAR Nesting
//
// Each of the above examples shows nesting of PAR via the
// function literal calling par.FOR. This pattern, a func()
// that just calls par.FOR is common and luckily Go lets us
// make it simpler.
//
// The par package defines what it refers to as "fn" functions. These are,
//
//     func DOfn(f ...func()) func()
//     func FORfn(start, limit int, f func(int)) func()
//
// The "fn" functions return a func() intended to be passed to par.DO and
// are used in created nested process structures.  FORfn is the most
// useful.
//
// Armed with FORfn we can now write,
//
// 	par.DO(
// 	    controlFuelRods,
// 	    monitorCoolant,
// 	    moveDials,
// 	    flashLights,
// 	    par.FORfn(0, 10, func(i int) {
// 	        ...  run generator i
// 	    }),
//     )
//
// Dynamic par
//
// To support more dynamic use the package defines a type, Group, that
// wraps a sync.WaitGroup and uses methods specific to starting
// goroutines and waiting for them.
//
// The user defines a variable of type par.Group and then uses
// the Add and Wait methods to start new goroutines and to wait
// for them to complete.
//
//	var g par.Group
//	for in := range channel {
//		g.Add(fn, in)
//	}
//	g.Wait()
//
package par
