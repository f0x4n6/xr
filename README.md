NAME
====

**hee** — High Performance Event Extractor

SYNOPSIS
========

**STDIN** | **hee** > **STDOUT**

DESCRIPTION
===========
- [ ] Set process priority
- [ ] Set max threads
- [ ] Use byte pool
- [ ] Calculate pool limit on thread count
- [x] Read from `STDIN`
- [x] Write to `STDOUT`
- [ ] Ignore `STDERR`
- [x] Panic any time
- [x] Exit with `0` or `1`
- [x] No usage
- [ ] No Windows
- [ ] No dependencies
- [x] README as manpage
- [ ] Answer the four primary questions
  - What?
  - When?
  - Where?
  - Who?

EXAMPLES
========

cat nist.dd | hee | uniq > events.log

SEE ALSO
========

**cat(1)**, **uniq(1)**
