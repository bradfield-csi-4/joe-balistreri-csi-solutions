package main

import "fmt"

func hot_potato(people []string, num int) string {
  for len(people) > 1 {
    for i := 0; i < num; i++ {
      person := people[len(people)-1]
      copy(people[1:], people[0:])
      people[0] = person
    }
    people = people[:len(people)-1]
  }
  return people[0]
}

func main() {
  people = []string{"Brad", "Kent", "jane", "susan", "david", "bill"}
  fmt.Println(hot_potato(people, 9))
}
