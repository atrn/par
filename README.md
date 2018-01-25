par
===

The par  package provides a  functions that implement  the "concurrent
control  structures"   from  the  occam  programming   language.   The
functions par.DO and  par.FOR provide Go analogues of  occam's PAR and
replicated-PAR control structures and  provide similar semantics which
are  the  same  as  those  of   the  idiomatic  use  of  the  standard
sync.WaitGroup type for goroutine synchronization.

In short, PAR.DO and par.FOR run some number of concurrent "processes"
(goroutines) and wait for them to  complete before they return and the
caller's execution proceeds.  They allow a program,  "process" in CSP,
to  "fan  out" execution  to  some  number  of concurrent  tasks  then
continue when they  finish. This turns out to be  quite a common thing
to want to do.

In terms  of Go the  functions encapsulate  the, now common,  idiom of
using  a  sync.WaitGroup  to synchronize  goroutine  completion.   The
functions offer  no actual new  functionality over what direct  use of
sync.WaitGroup affords, and actually provide  less, but their use does
make  for cleaner  code by  hiding the  implementation details  of the
synchronization.  The  functions eliminate clutter making  the process
structure  more obvious  and  therefore more  easily comprehended  and
maintained (i.e. not broken).

occam
-----

The  occam programming  language,  an  early CSP-inspired  programming
language used with  the Inmos transputer microprocessor,  used the PAR
and replicated-PAR  control structures to create  concurrent programs.
PAR runs one or more things  in parallel and replicated-PAR is occam's
verson of a parallel for-loop.

An  occam  program consists  of  "processes".   A  process is  like  a
statement in a programming language. There are a few special processes
such as  IF, STOP and SKIP  and others that define  program structure.
Unlike  most  languages  occam  requires programmers  define  how  the
processes of their program are executed.  Processes can either execute
sequentially or in parallel. The SEQ process is one that

The  simplest  occam   process  is  SEQ,  represented   by  the  "SEQ"
keyword.  SEQ  defines  sequential   execution  of  processes.  A  SEQ
statement  has  some  number  of  sub-statements,  or  processes,  and
executes each in turn.  I.e. SEQ defines the regular program execution
model we are all used to.

In addition to SEQ occam also has a PAR process that executes its
sub-processes in parallel.  Just like SEQ a PAR statement has some
number of sub-statements, or processes, but rather than executing them
sequentially, PAR executes them concurrently.  Each process runs in
what is essentially a goroutine.

Unlike  Go's "go"  statment occam's  PAR  waits for  its processes  to
complete - a PAR statement itself completes, and execution continuing,
only when its  sub-processes complete.  These semantics  allow for the
composition of process structures using PAR, SEQ et al.

occam for Go Programmers
------------------------

A Go programmer would not feel too out of place writing occam. They
would feel constrained but the core aspects of a CSP-based design are
similar in Go and occam despite Go's obvious advances. It is more
likely the expression syntax and overall static nature of occam would
be far more off putting for those accustomed to more forgiving dynamic
environments. A heap? What's that.

And while on this topic....an Inmos press release once stated that
they wanted to make occam "the FORTRAN of parallel processing". My
response upon reading this, noting its static nature, was "You have."

- channels are unbuffered, single reader, single writer
- shared read/write data is not permitted, the compiler checks
- occam has slices
- you have to use SEQ

Replacing sync.WaitGroup Use
----------------------------

PAR's synchronization sematics are the same those obtained with the
idiomatic use of a sync.WaitGroup. We want to block until all the
goroutines have finished.


Replicated-PAR
--------------

Occam also defines a parallel for-loop control structure called a
replicated-PAR.  A replicated PAR is the concurrent equivalent to a
normal for-loop. The loop is defined as an iteraton over an integer
range. Replicated-PAR creates zero or more processes, one for each
value in the range. The control structure is again synchronous and
blocks until all child processes complete.  Replicated-PAR nests


Benefits
--------

The synchronization semantics of PAR, SEQ and their replicated
counterparts allow for a slightly more structured approach to process
structure than that natively provided by Go's "go"
statement. Synchronization is an important part of concurrency and as
people have discovered not only involves protocols that provide
correct concurrent behaviour but also requires dynamic process
structures within their programs, at times adopting quite distinct
architectures for certain parts of the program's execution.


Using The Packge
----------------

The par package is used like all other packages by importing it.
Typically no renaming will occur and the package's exported symbols
are qualified by the package name, "par".

The pacakge exports, amongst others, the functiosn par.DO() and par.FOR().
In code these names read well. They are concise, direct and look almost
like language constructs. An example,

	par.DO(
	    controlFuelRods, // a func()
	    monitorCoolant, // ditto
	    moveDials, // and so on...
	    flashLights,
            controlSirens,
            func() {
		par.FOR(0, 10, func(int) {
		    ...  do work
		})
            },
        )

Each thing par.DO runs in parallel is a func().  In the above example
just assume there's a func controlFueldRods(), func monitorCoolant()
and so on (the example shows the use of named functions in place of
literal closures to clarify program structure).

In code the use of par.DO and par.FOR will read quite well. This is
important, the concurrent structure of a program can be a highly
important aspect of its design and often is not well communicated
between developers. I'm hoping these little functions an help.


Abusing Import

We can abuse Go's ability to import a package's symbols into the
current package namespace to drop the package qualification when
using DO() and FOR().  This makes using them seem a little more
like using a native language construct.

So, a user imports the package using ``import . "par"'', note
the dot. They can then call its functions without qualification.
The code above beomes,

	DO(
	    controlFuelRods,
	    monitorCoolant,
	    moveDials,
	    flashLights,
            controlSirens,
	    func() {
                FOR(0, 10, func(int) {
                    ...  do work
              	})
            },
        )


This code looks okay, if you accept the namespace pollution, but the
unqualified names DO() and FOR() are a little too generic and the code
itself, especially "DO(", a little hard to read.

So. We define synonyms, PAR and PAR_FOR, for DO() and FOR(). These names
mimic occam directly. Well, PAR does. A real replicated PAR requires
language changes. But I digress.

Code that uses the synonymous names, PAR and PAR_FOR, looks like,

	PAR(
	    controlFuelRods,
	    monitorCoolant,
	    moveDials,
	    flashLights,
            controlSirens,
            func() {
	        PAR_FOR(0, 10, func(int) {
	            ...  do work
    	        })
            },
        )
        
This, in my opinion, is much easy to read.

Why Not Both?

Users importing the package normally will also see the synonyms and
can replace uses of DO and FOR with PAR and PAR_FOR is desired.

	par.PAR(
	    controlFuelRods,
	    monitorCoolant,
	    moveDials,
	    flashLights,
            controlSirens,
            func() {
	        par.FOR(0, 10, func(int) {
	            ...  do work
    	        })
            },
        )


The par.PAR is a little redundant. And using par.PAR_FOR is right out.
It is for reaons such as this, readabilty, why I selected par.DO and
par.FOR as exported names. But also define PAR and PAR_FOR for the
import abusers, like myself.


PAR Nesting

Each of the above examples shows nesting of PAR via the
function literal calling par.FOR. This pattern, a func()
that just calls par.FOR is common and luckily Go lets us
make it simpler.

The par package defines what it refers to as "fn" functions. These are,

    func DOfn(f ...func()) func()
    func FORfn(start, limit int, f func(int)) func()

The "fn" functions return a func() intended to be passed to par.DO and
are used in created nested process structures.  FORfn is the most
useful.  

Armed with FORfn we can now write,

	PAR(
	    controlFuelRods,
	    monitorCoolant,
	    moveDials,
	    flashLights,
            controlSirens,
	    FORfn(0, 10, func(int) {
	            ...  do work
    	    }),
        )



SEQ

We provide a version of SEQ for completness. It's trivial so why not take the
emulation further. And there's a SEQfn too.

