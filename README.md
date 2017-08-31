# tall
tall, grande, ways to venti: a persistent filesystem

## Venti: a legacy from Plan9 Operating System
See: [Wikipedia](https://en.wikipedia.org/wiki/Venti)

## Features:
Tall is not suitable for mutation-heavy situation: it just generate too much history.
It is best used for archival: data stored are final and not likely to mutate too often.
In the design:
 - data is write-once, all histroy perserved
 - data storage is indexed by its digest hash value, called 'score' in Venti
 - metadata is also write-once except for the supermetadata, which is a fixed entry of all metadata
