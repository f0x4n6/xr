NAME
====

**triage**

triforce?
esp - event stream processor?
rsp - record stream processor?

SYNOPSIS
========

**STDIN** | **triage** > **STDOUT**

DESCRIPTION
===========
**triage** is an experimental high performance event record stream processor for fast forensic triaging.

- [x] Answer three primary questions
  - [x] What?
  - [x] When?
  - [x] Where?
- [x] README as manpage
- [x] Read from `STDIN`
- [x] Write to `STDOUT`
- [x] Debug to `STDERR`
- [x] Exit with `0` or `1`
- [x] Panic any time
- [x] Error handling with recover and exit
- [x] No usage
- [x] No dependencies
- [ ] Use byte pool
  - [ ] Calculate pool limit on thread count
- [ ] Set process priority
- [ ] Set max threads

EXAMPLES
========

```
$ cat nist.dd | triage | uniq | sort > events.txt
```

SEE ALSO
========

**cat(1)**, **uniq(1)**, **sort(1)**
