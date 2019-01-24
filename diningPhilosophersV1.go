/*
Implementation of the dining philosopher problem: https://en.wikipedia.org/wiki/Dining_philosophers_problem


1. Write a program (Version 1) to simulate the behavior of the philosophers, where each
philosopher is a thread and chopsticks are shared objects. Notice that you must prevent a
situation where two philosophers hold the same chopstick at the same time.

2. Write a program (Version 2) that modifies Version 1 so that it never reaches a state
where philosophers are deadlocked, that is, it is never the case that each philosopher
holds one chopstick and is stuck waiting for another to get the second chopstick.

3. Write a program (Version 3) so that no philosopher ever starves.

4. Write a program (Version 4) to provide a starvation-free solution for any number of
philosophers N.

Use multi-threading in your solution. You can choose any programming language
supporting threads for your implementation. Document your solution well. Keep all 4
versions of your solution into 4 separate files or folders. Provide a README.txt file with
detailed instructions on how to compile and run your program from the command
prompt.


Once a philosopher has chosen his seat, he can’t move to a new position
o As such he can only use the chopsticks to his left and right.
• Your	design	must	implement	the	chopsticks	as	shared	variables.
• The philosophers should only interact with each other through the chopsticks.
• Output	Format:
o When	a	philosopher	goes	from	eating	to	thinking	he	should	output: “%d	is	now	thinking.\n”
o When	a	philosopher	goes	from	thinking	to	hungry	he	should	output: “%d	is	now	hungry.\n”
o When	a	philosopher	goes	from	hungry	to	eating	he	should	output: “%d	is	now	eating.\n”
o Two	adjacent	philosophers	should	never	eat	at	the	same	time.
• The	programs	should	run	continuously	until	the	letter	‘n’	is	pressed.
• The	 last	 executable	 (starvation-free)	 should	 accept	 a	 command-line argument	that	represents	the	number	of	philosophers.
*/

package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"time"
)

const quanta = 2e9 // 2 seconds. timeout time. max time for eating and thinking
const numPh = 5
const numSticks = 5
const numSeats = 5
var chopsticks [numSticks] chan bool
var seats = make(chan int, numSeats)
const THINKING = 2
const HUNGRY = 1
const EATING = 0

type Philosopher struct{
	id	int
	seat int // -1 means not seated
	state int // what

}

func (ph *Philosopher) Eat() {
	ph.state = EATING
	fmt.Printf("Philosopher %d at seat %d is eating.\n", ph.id, ph.seat)
	time.Sleep(time.Duration(rand.Int63n(quanta)))
}


func (ph *Philosopher) Dine(announce chan *Philosopher) {
	ph.Think()
	// get hungry
	fmt.Printf("Philosopher %d is now hungry.\n", ph.id)
	ph.state = HUNGRY
	ph.Sit() // find a seat
	ph.GetChops()
	ph.Eat()
	ph.ReturnChops()
	ph.Stand() // give up seat
	announce <- ph
}

func (ph *Philosopher) GetChops(){
	//timeout := make(chan bool,1) // channel for timeout signal
	//go func() {time.Sleep(quanta); timeout <-true}() // wait then send timeout signal
	<- chopsticks[ph.seat] // use "right" or seat associated chopstick
	neighbor := GetNeighbor(ph.seat) // get the index of the neighbor chopstick
	fmt.Printf("Philosopher %d at seat %dpicked up his chopstick.\n", ph.id, ph.seat)
	select {
	case <-chopsticks[neighbor]:
		fmt.Printf("Philosopher %d at seat %d got seat %d's chopstick.\n", ph.id, ph.seat, neighbor)
		fmt.Printf("Philosopher %d has two chopsticks.\n", ph.id)
		return

		/*
		case <- timeout:
			ph.chopstick <- true
			ph.Think()
			ph.GetChop()
		*/
	}
}

func GetNeighbor(seat int)int{
	return (seat+1)%numSeats
}

/*
Takes an ID number and a pointer to a neighbor
returns a pointer to an initialized Philosopher
*/
func MakePh( id int) *Philosopher{
	ph:= &Philosopher{id, -1, THINKING}
	//ph.chopstick <- true // initialize chopstick as available
	return ph
}


func (ph *Philosopher) ReturnChops(){
	chopsticks[ph.seat] <-true
	chopsticks[GetNeighbor(ph.seat)]<-true

}

/*find an open seat*/
func (ph *Philosopher) Sit(){
	ph.seat = <- seats
}
/*give up seat*/
func (ph *Philosopher) Stand(){
	seats <- ph.seat
	ph.seat =-1
}

func (ph *Philosopher) Think() {
	ph.state = THINKING
	fmt.Printf("Philosopher %d is thinking.\n", ph.id)
	time.Sleep(time.Duration(rand.Int63n(quanta)))
}





func main() {

	fmt.Print("Press 'n' then enter to end program...\n")
	go func() {
		bufio.NewReader(os.Stdin).ReadBytes('n')
		os.Exit(0)
	}()

	//inti philosophers, chopsticks, and seats
	philosophers := make([]*Philosopher, numPh)
	var phil *Philosopher
	for i := 0; i < numPh; i++ {
		phil = MakePh(i)
		philosophers[i] = phil
	}

	for i := 0; i < numSticks; i++ {

		chopsticks[i] = make(chan bool,1)
		chopsticks[i] <- true
	}

	for i := 0; i < numSeats; i++ {
		seats <- i
	}

	fmt.Printf("There are %d philosophers and %d seats at a table.\n", numPh, numSeats)

	fmt.Printf("Seats have an associated chopstick. Philosophers must choose a seat and borrow from their neighbor to eat.\n")

	announce := make(chan *Philosopher)
	for _, phil := range philosophers {
		go phil.Dine(announce)
	}
	for i := 0; i < numPh; i++ {
		phil := <-announce
		fmt.Printf("%v is done dining.\n", phil.id)
	}
}