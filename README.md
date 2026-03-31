NAME
====

**tri** — High Performance Event Triage

SYNOPSIS
========

**STDIN** | **tri** > **STDOUT**

DESCRIPTION
===========

- [ ] Answer the four primary questions
  - [x] What?
  - [x] When?
  - [x] Where?
  - [ ] Who?
- [x] README as manpage
- [x] Read from `STDIN`
- [x] Write to `STDOUT`
- [x] Debug to `STDERR`
- [x] Exit with `0` or `1`
- [x] Panic any time
- [ ] Error handling with recover and exit
- [x] No usage
- [ ] No Windows
- [x] No dependencies
- [ ] Use byte pool
  - [ ] Calculate pool limit on thread count
- [ ] Set process priority
- [ ] Set max threads

EXAMPLES
========

```
$ cat nist.dd | tri | uniq | sort > events.txt
```

SEE ALSO
========

**cat(1)**, **uniq(1)**, **sort(1)**
