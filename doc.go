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
// 	    ControlFuelRods, // a func()
// 	    MonitorCoolant,
// 	    MoveDials,
// 	    FlashLights,
//          ControlSirens,
//          func() {
// 		par.FOR(0, 10, func(int) {
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
//
// Import Abuses
//
// We can abuse Go's ability to import a package's symbols into the
// current package namespace to drop the package qualification when
// using DO() and FOR().  This makes using them seem a little more
// like using a native language construct.
//
// So, a user imports the package using ``import . "par"'', note
// the dot. They can then call its functions without qualification.
// The code above beomes,
//
// 	DO(
// 	    controlFuelRods,
// 	    monitorCoolant,
// 	    moveDials,
// 	    flashLights,
//          controlSirens,
// 	    func() {
//                 FOR(0, 10, func(int) {
//                     ...  do work
//               	})
//             },
//         )
//
//
// This code looks okay, if you accept the namespace pollution, but the
// unqualified names DO() and FOR() are a little too generic and the code
// itself, especially "DO(", a little hard to read.
//
// So. We define synonyms, PAR and PAR_FOR, for DO() and FOR(). These names
// mimic occam directly. Well, PAR does. A real replicated PAR requires
// language changes. But I digress.
//
// Code that uses the synonymous names, PAR and PAR_FOR, looks like,
//
// 	PAR(
// 	    controlFuelRods,
// 	    monitorCoolant,
// 	    moveDials,
// 	    flashLights,
//             controlSirens,
//             func() {
// 	        PAR_FOR(0, 10, func(int) {
// 	            ...  do work
//     	        })
//             },
//         )
//
// This, in my opinion, is much easy to read.
//
// Why Not Both?
//
// Users importing the package normally will also see the synonyms and
// can replace uses of DO and FOR with PAR and PAR_FOR is desired.
//
// 	par.PAR(
// 	    controlFuelRods,
// 	    monitorCoolant,
// 	    moveDials,
// 	    flashLights,
//             controlSirens,
//             func() {
// 	        par.FOR(0, 10, func(int) {
// 	            ...  do work
//     	        })
//             },
//         )
//
//
// The par.PAR is a little redundant. And using par.PAR_FOR is right out.
// It is for reaons such as this, readabilty, why I selected par.DO and
// par.FOR as exported names. But also define PAR and PAR_FOR for the
// import abusers, like myself.
//
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
// 	PAR(
// 	    controlFuelRods,
// 	    monitorCoolant,
// 	    moveDials,
// 	    flashLights,
//          controlSirens,
// 	    FORfn(0, 10, func(i int) {
// 	            ...  run generator i
//     	    }),
//         )
//
//
package par
