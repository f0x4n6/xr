NAME
====

**xr** - experimental record analyzer

SYNOPSIS
========

```console
$ cat FILE | xr | uniq | sort
```

DESCRIPTION
===========

**XR** is an experimental high performance event record analyzer for fast forensic triaging. It targets to answer two main questions about event logs: WHAT and WHEN did it happen? Contrary to existing tools, it tries to answer these questions by analyzing the raw event record structure, rather than parsing whole chunks. By reading from any input stream, **XR** is capable of carving raw forensic disk and memory images.

INSTALLATION
============

```console
$ go install go.foxforensics.dev/xr@latest
```

REFERENCES
==========

- _Introducing the Microsoft Vista event log file format_ - Schuster, Andreas
- _Windows XML Event Log (EVTX) format_ - Metz, Joachim

SEE ALSO
========

[**dd(1)**](https://man7.org/linux/man-pages/man1/dd.1.html),
[**cat(1)**](https://man7.org/linux/man-pages/man1/cat.1.html),
[**uniq(1)**](https://man7.org/linux/man-pages/man1/uniq.1.html),
[**sort(1)**](https://man7.org/linux/man-pages/man1/sort.1.html)
