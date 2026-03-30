NAME
====

**tri** — High Performance Event Triage

SYNOPSIS
========

**STDIN** | **tri** > **STDOUT**

DESCRIPTION
===========

- [ ] Set process priority
- [ ] Set max threads
- [ ] Use byte pool
- [ ] Calculate pool limit on thread count
- [x] Read from `STDIN`
- [x] Write to `STDOUT`
- [ ] Debug to `STDERR`
- [x] Panic any time
- [x] Exit with `0` or `1`
- [x] No usage
- [ ] No Windows
- [x] No dependencies
- [x] README as manpage
- [ ] Answer the four primary questions
  - [x] What?
  - [x] When?
  - [ ] Where?
  - [ ] Who?

EXAMPLES
========

cat nist.dd | tri | uniq > events.log

SEE ALSO
========

**cat(1)**, **uniq(1)**
