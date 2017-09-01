# fs

## Introduction
'fs' is a backend to tall. It is based on filesystem (so called 'fs').

Tall backend stores score=>data maps. In 'fs' data is stored in a file named with score.

To allow large amount of data to be stored, a tree structures is used. 'fs' keeps the depth of the tree.

For example:
For score `9f86d081884c7d659a2feaa0c55ad015a3bf4f1b2b0b822cd15d6c15b0f00a08`

if the depth is 1 (init value) it is stored in `data/9f/9f86d081884c7d659a2feaa0c55ad015a3bf4f1b2b0b822cd15d6c15b0f00a08`;

if the depth is 2 it is stored in `data/9f/86/9f86d081884c7d659a2feaa0c55ad015a3bf4f1b2b0b822cd15d6c15b0f00a08`.

Each level of tree stores up to 256 sublevels.


## Planned/TODO:
 - [x] Hash names files
 - [x] Write once, read many times
 - [x] Files in directory tree
 - [ ] Directory tree expansion
 - [ ] Delete files
 - [ ] Directory tree shinkage
 - [ ] Stats
 - [ ] Use journal
