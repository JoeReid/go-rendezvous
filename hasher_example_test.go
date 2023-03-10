package rendezvous_test

import (
	"fmt"

	"github.com/JoeReid/go-rendezvous"
)

func ExampleHasher_ownership() {
	h := rendezvous.NewHasher(
		rendezvous.WithMembers("node1", "node2", "node3", "node4", "node5"),
	)

	fmt.Printf("The owner of item is %s\n", h.Owner("item"))
	// Output: The owner of item is node3
}

func ExampleHasher_replication() {
	h := rendezvous.NewHasher(
		rendezvous.WithMembers("node1", "node2", "node3", "node4", "node5"),
	)

	fmt.Printf("the nodes to replicate data to are: %v\n", h.Place("item", 3))
	// Output: the nodes to replicate data to are: [node3 node5 node1]
}

func ExampleHasher_priorityList() {
	h := rendezvous.NewHasher(
		rendezvous.WithMembers("node1", "node2", "node3", "node4", "node5"),
	)

	fmt.Printf("priority list: %v\n", h.Prioritise("item"))
	// Output: priority list: [node3 node5 node1 node2 node4]
}
