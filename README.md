URLook
======

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

Motivations
-----------

Inspired by ruby gem [awesome\_bot](https://rubygems.org/gems/awesome_bot)

License
-------

MIT. See LICENSE for details.
