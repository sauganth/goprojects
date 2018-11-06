package main

import (
  "fmt"
  "math"
)

type Circle struct {
  x,y,radius float64
}

func getSequence() func() int {
  i:=0
  nextNumber := func() int { //function assigned to variable and used as closure of i
    i = i + 1
    return i
  }
  return nextNumber
}

func (circle Circle) area() float64 { //method
  return math.Pi * circle.radius * circle.radius
}

func main() {
  fmt.Println("Hello")
  circle := Circle{x:0, y:0, radius:25}
  fmt.Printf("Area of circle is %f\n",circle.area())
  nextNumber := getSequence()
  fmt.Println(nextNumber())
  fmt.Println(nextNumber())

  var k [10]int
  k[0] = 5
  fmt.Printf("The value of k[0] is %d\n",k[0])

  var k1 = []int{2,3,5,6,7}
  fmt.Printf("The length of k1 is %d\n",len(k1))
}
