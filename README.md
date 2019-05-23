# HashTags

[![GoDoc](https://godoc.org/github.com/mwat56/go-hashtags?status.svg)](https://godoc.org/github.com/mwat56/go-hashtags)
[![License](https://img.shields.io/eclipse-marketplace/l/notepad4e.svg)](https://github.com/mwat56/go-hashtags/blob/master/LICENSE)

## Purpose

Sometimes one might want to search and find socalled `#hashtags` or `@mentions` in your texts (in a broader sense) and store them for later retrieval.
This package offers that facility.
It provides the `THashList` class which can be used to parse texts for the occurance of both `#hashtags` and `@mentions` and store the hits in an internal list for later lookup; that list can be both stored in a file and later read from a file.

## Installation

You can use `Go` to install this package for you:

    go get -u github.com/mwat56/go-hashtags

## Usage

In principle for each `#hashtag` or `@mention` a list of _IDs_ is maintained.
These _IDs_ can be any (string) data that identifies the text in which the `#hashtag` (or `@mention`) was found, e.g. a filename or some database record reference.
The only condition is that it is unique as far as the program using this package is concerned.

_Note_ that both `#hashtag` and `@mention` are stored lower-cased to allow for case-insensitive searches.

To get a `THashList` instance there's a simple way:

    fName := "mytags.lst"
    htl, err := hashtags.New(fName)
    if nil != err {
        log.PrintF("Problem loading file '%s': %v", fName, err)
    }
        // …
        // do something with the list
        // …
    written, err := htl.Store()
    if nil != err {
        log.PrintF("Problem writing file '%s': %v", fName, err)
    }

For more details please refer to the [package documentation](https://godoc.org/github.com/mwat56/go-hashtags).

## Licence

        Copyright © 2019 M.Watermann, 10247 Berlin, Germany
                        All rights reserved
                    EMail : <support@mwat.de>

> This program is free software; you can redistribute it and/or modify it under the terms of the GNU General Public License as published by the Free Software Foundation; either version 3 of the License, or (at your option) any later version.
>
> This software is distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.
>
> You should have received a copy of the GNU General Public License along with this program. If not, see the [GNU General Public License](http://www.gnu.org/licenses/gpl.html) for details.
