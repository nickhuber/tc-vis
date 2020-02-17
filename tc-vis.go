package main

import (
	"fmt"
	"os"

	"github.com/alecthomas/kong"
	"github.com/vishvananda/netlink"
)

func HandleQdisc(qdiscs []netlink.Qdisc, link netlink.Link, qdisc netlink.Qdisc, depth int) {
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
				HandleQdisc(qdiscs, link, cur_qdisc, depth+2)
			}
		}
	}
}

func HandleLink(link netlink.Link, extra_nesting bool) {
	var offset = 0
	if extra_nesting {
		offset = 1
	}
	qdiscs, _ := netlink.QdiscList(link)
	for idx, qdisc := range qdiscs {
		switch qdisc.Attrs().Parent {
		case netlink.HANDLE_ROOT:
			if extra_nesting {
				fmt.Printf("    ")
			}
			fmt.Printf("Qdisc %s %x: root\n", qdisc.Type(), qdisc.Attrs().Handle>>16)
			HandleQdisc(qdiscs, link, qdisc, idx+offset)
		case netlink.HANDLE_INGRESS:
			if extra_nesting {
				fmt.Printf("    ")
			}
			fmt.Printf("Qdisc %s %x:\n", qdisc.Type(), qdisc.Attrs().Handle>>16)
			HandleQdisc(qdiscs, link, qdisc, idx+offset)
		}
	}
}

var cli struct {
	Interface string `arg name:"interface" help:"Interface to query." type:"string" default:"" optional`
}

func main() {
	kong.Parse(&cli,
		kong.Name("tc-vis"),
		kong.Description("A hierarchical viewer for tc qdiscs and classes."))

	if cli.Interface != "" {
		link, err := netlink.LinkByName(cli.Interface)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Unknown interface %s\n", cli.Interface)
			os.Exit(1)
		}
		HandleLink(link, false)
	} else {
		links, _ := netlink.LinkList()
		for _, link := range links {
			fmt.Printf("%s\n", link.Attrs().Name)
			HandleLink(link, true)
		}
	}
}
