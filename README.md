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
- [ ] Read from `STDIN`
- [ ] Write to `STDOUT`
- [ ] Ignore `STDERR`
- [ ] Panic any time
- [ ] Exit with `0` or `1`
- [ ] No usage
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
