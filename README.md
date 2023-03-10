# go-rendezvous

This package provides a simplistic implementation of Rendezvous Hashing with a pluggable `hash.Hash` function.

From [Wikipedia](https://en.wikipedia.org/wiki/Rendezvous_hashing):

> Rendezvous or highest random weight (HRW) hashing is an algorithm that allows clients to achieve
> distributed agreement on a set of k options out of a possible set of n options.
>
> A typical application is when clients need to agree on which sites (or proxies) objects are assigned to.

Rendezvous Hashing is simple in principle, but wide-reaching in application.


For example, to determine ownership of an item:

```go
h := rendezvous.NewHasher(
    rendezvous.WithMembers("node1", "node2", "node3", "node4", "node5"),
)

fmt.Printf("The owner of item is %s\n", h.Owner("item"))
// Output: The owner of item is node3
```

Or to fairly place replicas of a file in a storage cluster:

```go
h := rendezvous.NewHasher(
    rendezvous.WithMembers("node1", "node2", "node3", "node4", "node5"),
)

fmt.Printf("the nodes to replicate data to are: %v\n", h.Place("item", 3))
// Output: the nodes to replicate data to are: [node3 node5 node1]
```

Or maybe to implement a priority list:

```go
h := rendezvous.NewHasher(
    rendezvous.WithMembers("node1", "node2", "node3", "node4", "node5"),
)

fmt.Printf("priority list: %v\n", h.Prioritise("item"))
// Output: priority list: [node3 node5 node1 node2 node4]
```
