URLook
======

[![Build Status](https://api.travis-ci.org/olshevskiy87/urlook.svg?branch=master)](https://travis-ci.org/olshevskiy87/urlook) [![MIT license](https://img.shields.io/badge/License-MIT-blue.svg)](https://lbesson.mit-license.org/) [![Go Report Card](https://goreportcard.com/badge/github.com/olshevskiy87/urlook)](https://goreportcard.com/report/github.com/olshevskiy87/urlook)

Find and check http(s) links in the text file.

Install
-------

```bash
$ go get github.com/olshevskiy87/urlook
```

Usage
-----

```bash
$ urlook -h
Usage: urlook FILENAME

Positional arguments:
  FILENAME               filename with links to check

Options:
  --help, -h             display this help and exit
```

```bash
$ urlook test_links
URLs to check: 6
  1. http://cb.vu/unixtoolbox.xhtml
  2. https://eax.me
  3. http://eax.me
  4. https://google.com
  5. https://github.com
  6. https://github.com/non-existent-url-non-existent-url-non-existent-url
✓→✓→✓x
  1. http://eax.me [301, Moved Permanently] -> https://eax.me/
  2. https://google.com [302, Found] -> https://www.google.ru/?gfe_rd=cr&dcr=0&ei=EFSZWqWQAcGDtAGArrLYBw
  3. https://github.com/non-existent-url-non-existent-url-non-existent-url [404, Not Found]
issues found: 3
```

```
$ grep github test_links | urlook
URLs to check: 2
  1. https://github.com
  2. https://github.com/non-existent-url-non-existent-url-non-existent-url
x✓
  1. https://github.com/non-existent-url-non-existent-url-non-existent-url [404, Not Found]
issues found: 1
```

```
$ echo 'check these links: https://ya.ru, https://www.reddit.com' | urlook
URLs to check: 2
  1. https://ya.ru
  2. https://www.reddit.com
✓✓
no issues found
```

Todo
----

- [x] read input from stdin
- [ ] try HEAD http request before GET
- [ ] add tests
- [ ] add white list option
- [ ] remove duplicate urls (a CLI option)

Motivations
-----------

Inspired by ruby gem [awesome\_bot](https://rubygems.org/gems/awesome_bot)

License
-------

MIT. See LICENSE for details.
