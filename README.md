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
		- [Methods](#methods)
			- [Haschtag related methods](#haschtag-related-methods)
			- [ID related methods](#id-related-methods)
			- [Mentions related methods](#mentions-related-methods)
			- [Maintenance methods](#maintenance-methods)
		- [Basic Usage](#basic-usage)
	- [Libraries](#libraries)
	- [Licence](#licence)

----

## Purpose

Sometimes one might want to search and find socalled `#hashtags` or `@mentions` in one's texts (in a broader sense) and store them for later retrieval.
This package offers that facility.
It provides the `THashTags` class which can be used to parse texts for the occurrence of both `#hashtags` and `@mentions` and store the hits in an internal list for later lookup; that list can be stored in a file and later loaded from that file.

## Installation

You can use `Go` to install this package for you:

	go get -u github.com/mwat56/hashtags

## Usage

For each `#hashtag` or `@mention` a list of _IDs_ is maintained.
These _IDs_ can be any (`int64`) data that identifies the text in which the `#hashtag` or `@mention` was found, e.g. some database record reference or article ID.
The only condition is that it is unique as far as the program using this package is concerned.

_Note_ that both `#hashtag` and `@mention` are stored lower-cased to allow for case-insensitive searches.

To get a `THashTags` instance there's a simple way:

	fName := "mytags.lst"
	ht, err := hashtags.New(fName, true)
	if nil != err {
		log.PrintF("Problem loading file '%s': %v", fName, err)
	}

	// ...
	// do something with the list
	// ...

	written, err := ht.Store()
	if nil != err {
		log.PrintF("Problem writing file '%s': %v", fName, err)
	}

The constructor function `New()` takes two arguments: A `string` specifying the name of the file to use for loading/storing the list's data, and a `bool` value indicating whether the list should be thread-safe or not. The setting for the latter depends on the actual use-case.

The package provides a global boolean configuration variable called `UseBinaryStorage` which is `true` by default. It determines whether the data written by `Store()` and read by `Load()` use plain text (i.e. `hashtags.UseBinaryStorage = false`) or a binary data format.
The advantage of the _plain text_ format is that it can be inspected by any text related tool (like e.g. `grep` or `diff`).
The advantage of the _binary format_ is that it is about three to four times as fast when loading/storing data and it uses a few bytes less than the text format.
For this reasons it's used by default (i.e. `hashtags.UseBinaryStorage == true`). During development of your own application using this package, however, you might want to change to text format for diagnostic purposes.

For more details please refer to the [package documentation](https://godoc.org/github.com/mwat56/hashtags/).

### Methods

There are several kinds of methods provided:

#### Haschtag related methods

The following methods can be used to handle hashtags:

 - `HashAdd(aHash string, aID int64) bool` inserts `aHash` as used by document `aID`, returning whether anything changed.
 - `HashCount() int` returns the number of hashtags currently handled.
 - `HashLen(aHash string) int` returns the number of documents using `aHash`.
 - `HashList(aHash string) []int64` returns a list of all document IDs using `aHash`.
 - `HashRemove(aHash string, aID int64) bool` removes the document `aID` from the `aHash` list, returning whether anything changed.

#### ID related methods

The following methods can be used to handle the document IDs of the list entries.

 - `IDlist(aID int64) []string` returns a list of hashtags and mentions occurring in the document identified by `aID`.
 - `IDparse(aID int64, aText []byte) bool` parses the given `aText` for hashtags and mentions and stores `aID` in the respective hashtag/mention lists, returning whether anything changed.
 - `IDremove(aID int64) bool` deletes the given `aID` from all hashtag/mention lists, returning whether anything changed.
 - `IDrename(aOldID, aNewID int64) bool` changes the given `aOldID` to `aNewID` in the rare case that a document's ID changed, returning whether anything changed.
 - `IDupdate(aID int64, aText []byte) bool` replaces the current hashtags/mentions stored for `aID` with those found in `aText`, returning whether anything changed.

#### Mentions related methods

The following methods can be used to handle mentions:

- `MentionAdd(aMention string, aID int64) bool` inserts `aMention` as used by document `aID`, returning whether anything changed.
- `MentionCount() int` returns the number of mentions currently handled.
- `MentionLen(aMention string) int` returns the number of documents using `aMention`.
- `MentionList(aMention string) []int64` returns a list of all document IDs using `aMention`.
- `MentionRemove(aMention string, aID int64) bool` removes the document `aID` from the `aMention` list, returning whether anything changed.

#### Maintenance methods

 - `Clear() *THashTags` empties the internal data structures: all `#hashtags` and `@mentions` are deleted.
 - `Filename() string` returns the filename given to the initial `New()` call for reading/storing the list's contents.
 - `Len() int` returns the current length of the list i.e. how many #hashtags and @mentions are currently stored in the list.
 - `LenTotal() int` returns the length of all #hashtag/@mention lists and their respective number of source IDs stored in the list.
 - `List() TCountList`  returns a list of #hashtags/@mentions with their respective count of associated IDs.
 - `Load() (*THashTags, error)` reads the configured file returning the data structure read from the file given with the `New()` call and a possible error condition.
 - `SetFilename(aFilename string) *THashTags` sets the filename for loading/storing the hashtags, returning the updated list instance.
 - `Store() (int, error)` writes the whole list to the configured file returning the number of bytes written and a possible error.
 - `String() string` returns the whole list as a linefeed separated string.

### Basic Usage

Although there are a lot of options (methods) available, basically the module is quite straightforward to use.

1. Create a new instance:

		myList := hashtags.New("myFile.db", true)

2. Whenever your application receives a new document, retrieve or create it's ID and text, then call

		ok := myList.IDparse(docID, docText)

## Libraries

The following external libraries were used building `HashTags`:

- [SourceError](https://github.com/mwat56/sourceerror)

## Licence

	Copyright © 2019, 2024  M.Watermann, 10247 Berlin, Germany
			All rights reserved
		EMail : <support@mwat.de>

> This program is free software; you can redistribute it and/or modify it under the terms of the GNU General Public License as published by the Free Software Foundation; either version 3 of the License, or (at your option) any later version.
>
> This software is distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.
>
> You should have received a copy of the GNU General Public License along with this program. If not, see the [GNU General Public License](http://www.gnu.org/licenses/gpl.html) for details.

----
[![GFDL](https://www.gnu.org/graphics/gfdl-logo-tiny.png)](http://www.gnu.org/copyleft/fdl.html)
