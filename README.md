NAME
====

**xr** - experimental record analyzer

SYNOPSIS
========

```console
INPUT | xr | OUTPUT
```

DESCRIPTION
===========

**xr** is an experimental high performance event record analyzer for fast forensic triaging. It targets to answer two main questions about event logs: WHEN and WHAT did happen? Contrary to existing tools, it tries to answer these questions by analyzing the raw event log record structure, rather than parsing whole chunk blocks. By reading from any input stream, **xr** is capable of carving raw forensic disk and memory images.

EXAMPLES
========

```console
$ cat image.dd | xr | uniq | sort
```

REFERENCES
==========

- _Windows XML Event Log (EVTX) format_ - Metz, Joachim
- _Introducing the Microsoft Vista event log file format_ - Schuster, Andreas
- _Detection and recovery of NSA’s covered up tracks_ - Jansen, Wouter

SEE ALSO
========

[**dd(1)**](https://man7.org/linux/man-pages/man1/dd.1.html),
[**cat(1)**](https://man7.org/linux/man-pages/man1/cat.1.html),
[**uniq(1)**](https://man7.org/linux/man-pages/man1/uniq.1.html),
[**sort(1)**](https://man7.org/linux/man-pages/man1/sort.1.html)
