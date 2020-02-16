package main

import (
	"fmt"
	"os"

	"github.com/vishvananda/netlink"
)

func handleQdisc(qdiscs []netlink.Qdisc, link netlink.Link, qdisc netlink.Qdisc, depth int) {
	classes, _ := netlink.ClassList(link, qdisc.Attrs().Handle)
	for _, class := range classes {
		if class.Attrs().Parent == netlink.HANDLE_ROOT && qdisc.Attrs().Parent != netlink.HANDLE_ROOT {
			continue
		}
		for i := 0; i < depth+1; i++ {
			fmt.Printf("    ")
		}
		fmt.Printf(
			"class %s %x:%x\n",
			class.Type(),
			class.Attrs().Handle>>16,
			class.Attrs().Handle&0x0000ffff,
		)
		for _, cur_qdisc := range qdiscs {
			if cur_qdisc.Attrs().Parent == class.Attrs().Handle {
				for i := 0; i < depth+2; i++ {
					fmt.Printf("    ")
				}
				fmt.Printf(
					"Qdisc %s %x:\n",
					cur_qdisc.Type(),
					cur_qdisc.Attrs().Handle>>16,
				)
				handleQdisc(qdiscs, link, cur_qdisc, depth+2)
			}
		}
	}
}

func handleLink(link netlink.Link, extra_nesting bool) {
	qdiscs, _ := netlink.QdiscList(link)
	for idx, qdisc := range qdiscs {
		switch qdisc.Attrs().Parent {
		case netlink.HANDLE_ROOT:
			if extra_nesting {
				fmt.Printf("    ")
			}
			fmt.Printf("Qdisc %s %x: root\n", qdisc.Type(), qdisc.Attrs().Handle>>16)
			handleQdisc(qdiscs, link, qdisc, idx+1)
		case netlink.HANDLE_INGRESS:
			if extra_nesting {
				fmt.Printf("    ")
			}
			fmt.Printf("Qdisc %s %x:\n", qdisc.Type(), qdisc.Attrs().Handle>>16)
			handleQdisc(qdiscs, link, qdisc, idx+1)
		}
	}
}

func main() {
	if len(os.Args) > 3 {
		fmt.Printf("usage: %s [interface]\n", os.Args[0])
		os.Exit(1)
	}
	if len(os.Args) == 2 {
		link, _ := netlink.LinkByName(os.Args[1])
		handleLink(link, false)
	} else {
		links, _ := netlink.LinkList()
		for _, link := range links {
			fmt.Printf("%s\n", link.Attrs().Name)
			handleLink(link, true)
		}
	}
}
