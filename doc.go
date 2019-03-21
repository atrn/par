/*
 *  Copyright © 2013. Andy Newman.
 */

// Package par provides functions to structure concurrent programs.
//
// The package defines the functions par.DO and par.FOR that implement
// synchronized "processes" (goroutines) in manner similar to the occam
// programming language's PAR and replicated-PAR control structures.
// Synchronized goroutines are a common process structure used in many
// concurrent programs. In Go this process structure is iodiomatically
// implemented using the sync package's WaitGroup type.
//
// The par.DO function runs some number of functions concurrently and
// waits for all to terminate before it returns. par.DO mimics occam's
// PAR control structure. In occam PAR applies to statements, in Go
// functions are used as the unit of code with closures relied upon for
// binding data values to specific function invocations.
//
// The par.FOR function concurrenly calls a single function a number
// of times defined by a range described by two integer "loop control"
// values. par.FOR iterates over this range, with a step of 1, calling
// the supplied function within a separate goroutine for each
// iteration of the loop. The function is passed the value of the
// loop control variable for its iteration, a form of "id". par.FOR
// then waits for all functions to complete, and their goroutines
// terminate, before it returns. par.FOR is a form of occam's
// "replicated-PAR" control structure.
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
//	    ControlFuelRods,
//	    MonitorCoolant,
//	    MoveDials,
//	    FlashLights,
//	    ControlSirens,
//	    func() {
//		par.FOR(0, 10, func(int) {
// 		    ...  do work
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
// 	    runSirens,
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
