package main

import (
	"fmt"
	"os"

	"github.com/vishvananda/netlink"
)

func handleQdisc(qdiscs []netlink.Qdisc, iface netlink.Link, qdisc netlink.Qdisc, depth int) {
	classes, _ := netlink.ClassList(iface, qdisc.Attrs().Handle)
	for _, class := range classes {
		for i := 0; i < depth+1; i++ {
			fmt.Printf("    ")
		}
		fmt.Printf(
			"class %s %x:%x\n",
			class.Type(),
			class.Attrs().Handle>>16,
			class.Attrs().Handle&0x0000ffff,
		)
		for _, parent_qdisc := range qdiscs {
			if parent_qdisc.Attrs().Parent == class.Attrs().Handle {
				for i := 0; i < depth+2; i++ {
					fmt.Printf("    ")
				}
				fmt.Printf(
					"Qdisc %s %x:\n",
					parent_qdisc.Type(),
					parent_qdisc.Attrs().Handle>>16,
				)
				handleQdisc(qdiscs, iface, parent_qdisc, depth+2)
			}
		}
	}
}

func main() {
	if len(os.Args) < 2 {
		fmt.Printf("usage: %s interface\n", os.Args[0])
		os.Exit(1)
	}
	iface_name := os.Args[1]
	iface, _ := netlink.LinkByName(iface_name)
	if iface == nil {
		fmt.Printf("%s does not exist\n", iface_name)
		os.Exit(2)
	}
	qdiscs, _ := netlink.QdiscList(iface)
	for idx, qdisc := range qdiscs {
		switch qdisc.Attrs().Parent {
		case netlink.HANDLE_ROOT:
			fmt.Printf("Qdisc %s %x: root\n", qdisc.Type(), qdisc.Attrs().Handle>>16)
			handleQdisc(qdiscs, iface, qdisc, idx)
		case netlink.HANDLE_INGRESS:
			fmt.Printf("Qdisc %s %x:\n", qdisc.Type(), qdisc.Attrs().Handle>>16)
			handleQdisc(qdiscs, iface, qdisc, idx)
		}
	}
}
