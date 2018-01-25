# par - occam-style concurrency _primitives_

The par  package provides functions that implement  occam-style PAR
and replicated-PAR control structures. These provide synchronization
upon gorouting completion in the same way as idiomatic `sync.WaitGroup`
usage.

`par.DO` calls some number of function concurrently and waits for
all to complete before it returns. `par.FOR` calls a single function
N times concurrently, where N is defined by an integer range (with
each call being passed its _index_ in the range).

`par.DO` mimics the occam `PAR` keyword and `par.FOR` occam's
_replicated-PAR_ concurrent for-loop. The functions are implemented
using `sync.WaitGroup` and hides the repetitive _clutter_.

## An example

Imagine we have some functions that run loops to do some control
operation. In our system we run these concurrently, perhaps they
communicate but that's details. We run then concurrently and wait
for them to finish.  Which in this case they never do...

	par.DO(
		ControlFuelRods,
		MonitorCoolant,
		MoveDials,
		FlashLights,
		ControlSirens,
		func() {
			par.FOR(0, 10, func(number int) {
				MonitorDoor(number)
			})
		},
        )


## Hiding sync.WaitGroup

The `par` functions encapsulate  the, now common,  idiom of using  a 
`sync.WaitGroup`  to synchronize  goroutine  completion.   The `par`
functions offer  no actual new  functionality over what direct  use of
sync.WaitGroup affords, and actually provide  less, but their use does
make  for cleaner  code by  hiding the  implementation details  of the
synchronization.  The  functions eliminate clutter making  the process
structure  more obvious  and  therefore more  easily comprehended  and
maintained (i.e. not broken).

## Abusing Import

We can abuse Go's ability to import a package's symbols into the
current package namespace and drop the package qualification when
using DO() and FOR().  The lack of qualification makes using them
seem a little more like using a native language construct.

So, a user imports the package using ``import . "par"'', note
the dot. They can then call its functions without qualification.
The code above beomes,

	DO(
		ControlFuelRods,
		MonitorCoolant,
		MoveDials,
		FlashLights,
		ControlSirens,
		func() {
			FOR(0, 10, func(number int) {
				MonitorDoor(number)
		},
        )


This code looks okay, if you accept the namespace pollution, but the
unqualified names DO() and FOR() are a little too generic and the code
itself, especially "DO(", a little hard to read.

## Synonyms, PAR and PAR_FOR 

The package define synonyms for DO and FOR, PAR and PAR_FOR. When
imported into `.` code looks like,

	PAR(
		ControlFuelRods,
		MonitorCoolant,
		MoveDials,
		FlashLights,
		ControlSirens,
		func() {
			PAR_FOR(0, 10, func(number int) {
				MonitorDoor(number)
		},
        )
        

## Nested PAR

Each of the above examples shows nesting of PAR via the
function literal calling par.FOR. This pattern, a func()
that just calls par.FOR is common and luckily Go lets us
make it simpler.

The par package defines what it refers to as _fn_ functions.
These are,

	func DOfn(f ...func()) func()
	func FORfn(start, limit int, f func(int)) func()

The _fn_ functions return a func() intended to be passed to par.DO and
are used in created nested process structures.  FORfn is the most
useful.  

Armed with FORfn we can now write,

	PAR(
	    	ControlFuelRods,
	    	MonitorCoolant,
	    	MoveDials,
	    	FlashLights,
		ControlSirens,
	    	FORfn(0, 10, func(number int) {
			MonitorDoor(number)
    	    	}),
        )
