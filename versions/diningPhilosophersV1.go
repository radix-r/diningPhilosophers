/*
Implementation of the dining philosopher problem: https://en.wikipedia.org/wiki/Dining_philosophers_problem


V1: This is my first pass at the dining philosophers problem. Im implementing a modified version of  Dijkstra's proposed solution. Each philosopher eats and thinks for an amount of tim limited by a quanta, 2 seconds in this case. This out of the gate prevents starvation because if a philosopher is waiting on a chopstick it is guaranteed to be free in the amount of time dictated by quanta. This version is still prone to deadlock if all philosophers get their first chop stick at the same time. Influenced by https://github.com/doug/go-dining-philosophers.

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
var done = false

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
	for !done{
		ph.Think()
		// get hungry
		fmt.Printf("Philosopher %d is now hungry.\n", ph.id)
		ph.state = HUNGRY
		ph.Sit() // find a seat
		ph.GetChops()
		ph.Eat()
		ph.ReturnChops()
		ph.Stand() // give up seat
	}
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
		done = true
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
