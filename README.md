# HashTags

[![Golang](https://img.shields.io/badge/Language-Go-green.svg)](https://golang.org/)
[![GoDoc](https://godoc.org/github.com/mwat56/hashtags?status.svg)](https://godoc.org/github.com/mwat56/hashtags/)
[![Go Report](https://goreportcard.com/badge/github.com/mwat56/hashtags)](https://goreportcard.com/report/github.com/mwat56/hashtags)
[![Issues](https://img.shields.io/github/issues/mwat56/hashtags.svg)](https://github.com/mwat56/hashtags/issues?q=is%3Aopen+is%3Aissue)
[![Size](https://img.shields.io/github/repo-size/mwat56/hashtags.svg)](https://github.com/mwat56/hashtags/)
[![Tag](https://img.shields.io/github/tag/mwat56/hashtags.svg)](https://github.com/mwat56/hashtags/tags)
[![License](https://img.shields.io/github/license/mwat56/hashtags.svg)](https://github.com/mwat56/hashtags/blob/main/LICENSE)

- [HashTags](#hashtags)
	- [Purpose](#purpose)
	- [Installation](#installation)
	- [Usage](#usage)
	- [Libraries](#libraries)
	- [Licence](#licence)

----

## Purpose

Sometimes one might want to search and find socalled `#hashtags` or `@mentions` in one's texts (in a broader sense) and store them for later retrieval.
This package offers that facility.
It provides the `THashList` class which can be used to parse texts for the occurrence of both `#hashtags` and `@mentions` and store the hits in an internal list for later lookup; that list can be stored in a file and later read from that file.

## Installation

You can use `Go` to install this package for you:

	go get -u github.com/mwat56/hashtags

## Usage

In principle for each `#hashtag` or `@mention` a list of _IDs_ is maintained.
These _IDs_ can be any (string) data that identifies the text in which the `#hashtag` or `@mention` was found, e.g. a filename or some database record reference.
The only condition is that it is unique as far as the program using this package is concerned.

_Note_ that both `#hashtag` and `@mention` are stored lower-cased to allow for case-insensitive searches.

To get a `THashList` instance there's a simple way:

	fName := "mytags.lst"
	ht, err := hashtags.New(fName)
	if nil != err {
		log.PrintF("Problem loading file '%s': %v", fName, err)
	}

	// …
	// do something with the list
	// …

	written, err := ht.Store()
	if nil != err {
		log.PrintF("Problem writing file '%s': %v", fName, err)
	}

The package provides a boolean configuration variable called `UseBinaryStorage` which is `true` by default.
It determines whether the data written by `Store()` and read by `Load()` use plain text (i.e. `hashtags.UseBinaryStorage = false`) or a binary data format.
The advantage of the plain text format is that it can be inspected by any text related tool (like e.g. `grep` or `diff`).
The advantage of the binary format is that it is about three to four times as fast when loading/storing data and it uses a few bytes less than the text format.
For this reasons it's used by default (i.e. `hashtags.UseBinaryStorage == true`); during development of your own application using this package, however, you might want to change to text format for diagnostic purposes.

For more details please refer to the [package documentation](https://godoc.org/github.com/mwat56/hashtags/).

## Libraries

The following external libraries were used building `HashTags`:

- [SourceError](https://github.com/mwat56/sourceerror)

## Licence

	Copyright © 2019, 2024 M.Watermann, 10247 Berlin, Germany
			All rights reserved
		EMail : <support@mwat.de>

> This program is free software; you can redistribute it and/or modify it under the terms of the GNU General Public License as published by the Free Software Foundation; either version 3 of the License, or (at your option) any later version.
>
> This software is distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.
>
> You should have received a copy of the GNU General Public License along with this program. If not, see the [GNU General Public License](http://www.gnu.org/licenses/gpl.html) for details.

----
[![GFDL](https://www.gnu.org/graphics/gfdl-logo-tiny.png)](http://www.gnu.org/copyleft/fdl.html)
