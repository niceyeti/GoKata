# GoKata

This repo contains algorithms, patterns, data structures in golang for personal exercise.

This code is free to use, but should be evaluated, extended, and tested extensively before use.
Much of this code is merely first-pass, with minimal testing.

Points of concern and required rewrite I have left incomplete:
* memory leaks
* concurrency/races
* benchmarks: on some inputs the builtin sort runs faster than concurrent qsort; it may be using a
  memory, language, or runtime strategy that could be copied.
* grep for TODOs in the code