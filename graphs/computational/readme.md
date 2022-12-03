This work is not a single graph data structure of algorithm,
but an attempt at an observable graph library in golang with linq'ish semantics.
These go by many names, e.g. computational graphs, but the
scope here is limited to those that compose observer patterns
into graphs, like mobx. Others can be used for autograd functionality;
but this is not such a numerical library. Mobx style graphs are
more for generating data structures and state, as opposed to computing 
mathematical functions or linear algebraic operators.

Note: this library is intentionally exotic. The goal here is to push differentiation/optimization
of data structure and algorithms for computational graphs, to see what's possible and what isn't,
and to derive some code patterns/views. The lessons apply to other things, such as optimizing
system operations and db storage, e.g., replace operators in some graph with I/O or other
system operations to get feedback on how quickly an entire system scheme is working. But also
note that some builtin golang code tools already provide flamecharts for perf monitoring.

Example:
    - On some change in input data
    - Transform the data into some other form (usually from dto's to frontend business objects)
    - Filter by some criterion ('Where')
    - Compute something
    - Generate or update a new list of info
    - Requirement: be able to intercept any of the above steps to add or implement other features

Straightforward example: a dynamically updated list
    - Some network component signals a response/message of data objects, LDO
    - From LDO, derive a list of front end objects, LFO, through some transformation function f
    - Given LFO, update a list of ui components:
        - add new items, remove old ones (not in current input)
        - update existing items, using efficient delta techniques
    The library should be intelligent enough to figure out what changes are needed, but nothing
    more. For instance, many observable-chain libraries simply apply a 'dumb' function at each stage.
    This library should compile down such that it can intelligently infer the most efficient way to
    update each object. For example, assume we want a chain like:
        1) receive objects and convert them to object type D2
        2) update a list with the new D2s
        3) modify some other list of ui objects
    And now assume that only one object in the input list (1) has actually changed. Steps 2/3 should
    not execute or only partially execute, w.r.t the minimal changes needed to actually update the
    single object. A dumb, striaghtforward library implementation would simply implement the steps
    above in order, e.g. convert every object again, reset every downstream, data structure

Requirements:
- build observable chains of data structures
- ensure they (the code) is open to modification and extension (api cannot be too rigid)
- ensure that changes can be throttled, and overall have good control mechanisms/features,
  such that the developer-user can easily determine how/when to update things, throttle things, etc
- efficient list resets: data structure updates should be as efficient as possible, and should
  minimize computation/new memory. No "reset everything" logic.
- no decorator magic: I'm not a fan. Decorators (e.g. struct tags) make logic less readable and usually
  rely on runtime magic.

Inspiration: being intentionally ethereal here, but the autograd assumption is that all operations
in a computer are differentiable, hence derivatives can be computed based on a structured description
of a math expression. We wish to borrow similar ideas here: can operations on data structures be
differentiated? Of course, this is merely borrowing inspiration. Differentiation is done with respect
to a specific variables/variables, for finding the zeroes (minmax) of functions, and so on, whereas the
optimization objective to apply to data structures minimizing the number of operations required to update
them to a new state. For the simplest example, pretend some arbitrary number of operations output
an xml string describing a list of structured objects; a change library applies a specific list of operations
(sort this, move that, update this object's field), the goal is to minimize this list. For instance,
say a number of operations are inverses or override one another's after state: set field f1 to zero,
set f1 to one, set f1 back to zero. The library should observe this, and only execute the required operations.
The code requirement is that the internal intelligence span these calls:
    .Reset(...
    .Convert(...
    .SetValue(...
    .Where(...
The calls would not be independent calls in the library; instead, the library would be aware of the entire
chain of calls, and capable of minimizing them.

Some possible hacks:
    - treat everything as a vector
    - treat everything as a string, and run virtual-dom string-diff ops on a composition of string transformations
    - require observable objects to implement some interface to outsource some work, e.g. Equals(other) or Delta(other)
    - require observable objects to implement an efficient String() method and operate in string space
    - end-to-end minimization: rather than intermediate optimizations and tons of rules and one-off query optimizations,
      one might use some strategy as a trade-off: given S1 is the current derived view of the data after all
      operations, and S2 is the next (pending) state, compare S1 and S2 as strings to derive only the required
      changes to map S1 -> S2. This gets rid of some costly or repeated operations 

Challenges:
- autograd is amenable to python since it is a dynamic and introspective language, not statically compiled like golang
- async? chans? They are highly amenable to building computational graphs, cancellation, etc. 
  However, the external api could be synchronous (define your graph), then internally implement the graph requirements
  and channels, using channerics
- testability and observability: be able to observe changes when and how they occur, such as in a web ui. This is
  just the generic feature provided by tools like pperf: you really don't know how things are working/benchmarks
  until you can observe it introspectively.
- plan viewer: like sql-plan viewers, provide a view of the computational graph and benchmarks (percentage time spent)
- state historian: be able to rewind through old states, like a transaction db. Note this could be configurable,
  for testing/debugging. This might be like a httptest-recorder or Memento pattern.
- context and cancellation: pass in context at any stage, causing changes to be cancelled or aborted as needed.
    Maybe it is taking too long for a service to respond, or a ui to refresh.
- lazy evaluation (compare linq vs observable patterns: lazy evaluation is done in linq because the expression is the
  unit of transaction, whereas observable patterns are long-lived)
- declarative derivations are more amenable to optimization (the same as with declarative linq expressions)

Conclusions:
- Given all of the above and the history of the problem, I favor a heuristic approach: provide custom-behavior
  hooks where needed, to allow the developer to control how changes are performed, such that they can override
  default behavior, i.e. perform Swap behavior rather than resetting an entire list. Then, apply end to end
  diff'ing of some kind, and capture as much low-hanging optimization as possible.
- Per the previous, automatic optimization is more of an at-scale concern. If someone paid for this,
  seems worth it. But not for a first pass. Calm down son.

Descriptive problem case: a major bug occurred that describes the needs for this library to address. A large
graph was implemented in code; upstream, when two items were swapped, this caused a downstream list to perform
two complete resets on a more expensive set of objects. To get around this, a flag was set at the site of the
upstream operation, indices set, and then read downstream to perform a swap instead of a reset. The issue
is that intervening code was modified later, affecting the invariants of the indices, and throwing an out of range.
    The issue is twofold, maintainability and optimization/performance:
    - poor performance for simple operations (simple item swaps triggering full resets)
    - ease of modifying the method by which updates occur, instead of bulk operations
    Conclusion: the item swap case is a good test case for design. How does your design address this?



References and projects:
- https://github.com/ahmetb/go-linq/issues/96:
    'Based on my last experimentation, Go generics unfortunately isn't advanced enough to support this by design.
    But I also encourage others to try Select().Order() where each method takes a generic type and returns another generic type processed further down the chain.'
- https://www.microsoft.com/en-us/research/wp-content/uploads/2011/06/paper-pldi.pdf
    Steno, gets into depth on the compiled view of the expression
    https://www.microsoft.com/en-us/research/publication/steno-automatic-optimization-of-declarative-queries/

Sample code:
    var graph Graph<DTO>
    func NewComponent() {
        // init the static computation definition
        graph = NewGraph[DTO]()
            .Reset(onReset)                 // func onReset(items []DTO)
            .Convert[T1,T2](toTargetType)   // func toTargetType(item T1) T2
            .WithThrottle(time.Second * 5)  // internally queue items on the graph to reduce downstream change-signals
            .Merge(externalGraph, mergeFoo) // combine internal graph at this point, with external graph, using mergeFoo

        g2 := graph.Split()                          // split output of the graph to a second channel, such as for another ui
            .Convert[TGraph,TOtherUI](toOtherModel) // convert TGraph to a model for a different ui

        ui1 := graph.Convert[TGraph, TUIComponent](toUIComponent1)
        ui2 := graph.Convert[TGraph, TUIComponent](toUIComponent2)
    }

    func handler(freshItems []DTO) Observable {
        newState := graph.Next(freshItems) // Get new state, may be same as the old state.
        ...

        ui1String := ui1.Render() // Get the final output strings for each component
        ui2String := ui2.Render()

        diffs := htmlDiffer(ui1String, ui2String) // some view diffing magic, a la tabular/fastview stuff

        throttle.Pause()   // pause some throttle; assume the graph is aware of this manual throttle, and will not execute after Pause() is called
        done := clientSock.Send(diffs) // send the diffs to the client to update
        defer func() {
            <-done
            throttle.Continue()  // notify the throttle the client is ready for more input; this is merely to demo such structure
        }
    }

Questions
1) recall that linq is heavily sql-influenced; sql in turn utilizes query-optimization engines. These properties
may be the same/overlapping requirements as this library, e.g. query optimization.

This is the linq api; how much is 'differentiable' per the criterion defined earlier?
    Sum / Min / Max / Average
    Aggregate
    Join / GroupJoin
    Take / TakeWhile
    Skip / SkipWhile
    OfType
    Concat
    OrderBy / ThenBy
    Reverse
    GroupBy
    Distinct
    Union / Intersect / Except
    SequenceEqual
    First / FirstOrDefault / Last / LastOrDefault
    Single
    SingleOrDefault
    ElementAt
    Any / All
    Contains
    Count
Also these conversions/data structure ops:
    AsEnumerable
    AsQueryable
    ToArray
    ToList
    ToDictionary
    ToLookup
    Cast
    OfType

Scratch design work:
    1) 
        eles := NewGraph[TInput](inputChan).
            // These calls are non-sensical examples, except for those commented.
            Convert(convFoo).
            Distinct(idFunc).
            Swap(i,j).          // a) Some sort of swapping operation (possibly custom imp).
            GroupBy(groupFunc).
            Select(subFunc).
            Sort(sortFunc).
            Select(eleFunc)     // b) How can this output list ensure it is not completely reset because of a mere swap above?
        One off optimizations:
            List-change event carries forward the changes to be applied: they can carry forward only to the next
            linear/list operation (Select, etc). 
                Events: reset, clear, swap, concat, convert... hmm, most of these are just first-class methods, not events

            - Perform operations in place: reset, clear, swap, etc., and merely signal the change.
                This seems like a clever way to perform a carry-forward version of these calls:
                    list.Swap().Concat().Take()
                    Become merely a sequence of operations, then enumerate them once a terminal call is reached.

