I end up using TC a lot to manage many rules that have a few levels of parents
and it can be hard to keep track of how everything relates so this tool just
formats the Qdiscs and classes into an hierarchical format.

The only requirement is github.com/vishvananda/netlink for fetching the data from netlink, so installation is just

    go get github.com/vishvananda/netlink
    go build
    ./ts-vis eth0

Or if you want you can just use `make install` to put it in /usr/local/bin

Turns this

    [nick@fedora ~]$ tc qdisc show dev eth0
    qdisc htb 25: root refcnt 65 r2q 10 default 0 direct_packets_stat 92720 direct_qlen 1000
    qdisc sfq 8001: parent 25:1 limit 127p quantum 1514b depth 127 divisor 1024
    qdisc htb 8003: parent 25:420 r2q 10 default 0 direct_packets_stat 0 direct_qlen 1000
    qdisc fq_codel 8002: parent 25:60 limit 10240p flows 1024 quantum 1514 target 5.0ms interval 100.0ms memory_limit 32Mb ecn
    qdisc htb 8004: parent 8003:42 r2q 10 default 0 direct_packets_stat 0 direct_qlen 1000
    qdisc ingress ffff: parent ffff:fff1 ----------------
    [nick@fedora ~]$ tc class show dev eth0
    class htb 25:1 root leaf 8001: prio 0 rate 5Mbit ceil 5Mbit burst 1280b cburst 1600b
    class htb 25:36 root prio 0 rate 420Mbit ceil 420Mbit burst 53760b cburst 1522b
    class htb 25:60 root leaf 8002: prio 0 rate 69Mbit ceil 69Mbit burst 88320b cburst 1587b
    class htb 25:420 root leaf 8003: prio 0 rate 420Mbit ceil 420Mbit burst 53760b cburst 1522b
    class htb 8003:42 root leaf 8004: prio 0 rate 2Mbit ceil 2Mbit burst 256b cburst 1600b
    class htb 8003:43 root prio 0 rate 2Mbit ceil 2Mbit burst 256b cburst 1600b
    class htb 8003:42 root leaf 8004: prio 0 rate 2Mbit ceil 2Mbit burst 256b cburst 1600b
    class htb 8003:43 root prio 0 rate 2Mbit ceil 2Mbit burst 256b cburst 1600b

into this

    [nick@fedora ~]$ tc-vis eth0
    Qdisc htb 25: root
        class htb 25:1
            Qdisc sfq 8001:
        class htb 25:36
        class htb 25:60
            Qdisc fq_codel 8002:
        class htb 25:420
            Qdisc htb 8003:
                class htb 8003:42
                    Qdisc htb 8004:
                class htb 8003:43
                class htb 8003:42
                    Qdisc htb 8004:
                class htb 8003:43
    Qdisc ingress ffff:
