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
var stateSticks [numSticks]bool
var statePh [numPh]int
var chopsticks [numSticks] chan bool
const THINKING = 2
const HUNGRY = 1
const EATING = 0

type Philosopher struct{
	id	int
	seat int
	//chopstick chan bool // associate chopsticks with Philosopher
	//neighbor *Philosopher // neighbor seated to the right
}

/*
Takes an ID number and a pointer to a neighbor
returns a pointer to an initialized Philosopher
*/
func MakePh( id int, neighbor *Philosopher) *Philosopher{
	ph:= &Philosopher{id, make(chan bool, 1), neighbor}
	ph.chopstick <- true // initialize chopstick as available
	return ph
}

func (ph *Philosopher) Think() {
	fmt.Printf("Philosopher %d is thinking.\n", ph.id)
	time.Sleep(time.Duration(rand.Int63n(quanta)))
}

func (ph *Philosopher) Eat() {
	fmt.Printf("Philosopher %d is eating.\n", ph.id)
	time.Sleep(time.Duration(rand.Int63n(quanta)))
}

func (ph *Philosopher) GetChops(){
	//timeout := make(chan bool,1) // channel for timeout signal
	//go func() {time.Sleep(quanta); timeout <-true}() // wait thensend timeout signal
	<- chopsticks[ph.seat] // us "right" or seat associatedchopstick
	fmt.Printf("Philosopher %d picked up his chopstick", ph.id)
	select {
	case <-ph.neighbor.chopstick:
		fmt.Printf("Philosopher %d got Philosopher %d's chopstick.\n", ph.id, ph.neighbor.id)
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

func (ph *Philosopher) ReturnChops(){
	chopsticks[ph.seat] <-true
	ph.chopstick <-true
	ph.neighbor.chopstick <- true
}

func (ph *Philosopher) dine(announce chan *Philosopher) {
	ph.Think()
	ph.GetChops()
	ph.Eat()
	ph.ReturnChops()
	announce <- ph
}

func main(){

	fmt.Print("Press 'n' then enter to end program...\n")
	go func() {

		bufio.NewReader(os.Stdin).ReadBytes('n')
		os.Exit(0)
	}()

	philosophers := make([]*Philosopher, numPh)
	var phil *Philosopher
	for i:=0; i<numPh; i++  {
		phil = MakePh(i, phil)
		philosophers[i] = phil
	}

	// complete the ring
	philosophers[0].neighbor = phil

	fmt.Printf("There are %d philosophers and %d seats at a table\n",numPh, numSeats)

	fmt.Printf("Seats have an associated chopstick. Philosophers must choose a seat and borrow from their neighbor to eat\n")


	for i:=0;i<10;i++{
		time.Sleep(time.Duration(rand.Int63n(quanta)))
	}
}
