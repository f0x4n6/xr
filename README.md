NAME
====

**xr** - experimental record analysis

SYNOPSIS
========

INPUT | **xr** | OUTPUT

DESCRIPTION
===========
**xr** is an experimental high performance event record analyzer for fast forensic triaging. It targets to answer three main questions about event logs: *WHEN*, *WHERE* and *WHAT* did happen? Contrary to existing tools, it answers these question by analyzing the basic event log record structure.

it reads any input stream, including carved disk images.

- [ ] Use byte pool
  - [ ] Calculate pool limit on thread count
- [ ] Set process priority
- [ ] Set max threads

EXAMPLES
========

```
$ cat nist.dd | xr | uniq | sort
```

SEE ALSO
========

**cat(1)**, **uniq(1)**, **sort(1)**
