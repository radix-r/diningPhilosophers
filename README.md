## Install golang 1.11.4 or later

https://golang.org/dl/

Ubuntu:
sudo snap install --classic go

##Run
go run diningPhilosophersV\<version number>.go \<command line args>

Only version 4 takes command line arguments

example: go run diningPhilosophersV4.go 10

runs dining philosophers program with 10 philosophers.

pressing 'n' followed by the enter key sends the finish signal. It may take a few seconds for the philosophers to finish dining. Ctrl + c stops the program immediately.

##Overview

####V1: 
This is my first pass at the dining philosophers problem. Im implementing a modified version of  Dijkstra's proposed solution. Each philosopher eats and thinks for an amount of tim limited by a quanta, 2 seconds in this case. This out of the gate prevents starvation because if a philosopher is waiting on a chopstick it is guaranteed to be free in the amount of time dictated by quanta. This version is still prone to deadlock if all philosophers get their first chopstick at the same time. Influenced by https://github.com/doug/go-dining-philosophers.



####V2: 
This program is very similar to V1 except it has a timeout built into the GetChops() function. This means that in the situation where all seats and chopsticks are taken simultaneously a few philosophers will time out and put down there chopsticks to let the others eat. This prevents deadlock.

####V3: 
This program is identical V2. Because of the limited eat and think times dictated by quanta and the timeout feature in the GetChops() function starvation is prevented.
 
####V4: 
 This program is almost identical to V3 except it takes a single command line input, the number of philosophers to simulate.
