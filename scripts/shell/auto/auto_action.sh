#!/usr/bin/expect
set password "Mkhahah@1032"
set ipaddr [lrange $argv 0 0]
set cmds [lrange $argv 1 end]
set timeout -1
set port 22
spawn ssh -p $port $ipaddr
match_max 100000
expect "*?\(yes/no\)*" {
    send -- "yes\r"
        expect "*?assword:*"
   send -- "$password\r"
send -- "\r"
} "*?assword:*"  {
    send -- "$password\r"
send -- "\r"
}
send -- "\r"
expect "*#*"
set cmds [split "$cmds" ";"]
foreach cmd $cmds {
    set cmd [string trim $cmd]
    send -- "\r"
    expect "*#"
    send -- "$cmd\r"
    expect "*#"
    send -- "\r"
    sleep 1
}
exit 0
expect eof

