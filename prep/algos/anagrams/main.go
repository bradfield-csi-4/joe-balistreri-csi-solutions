package main

import(
  "fmt"
)

func is_anagram(a string, b string) bool {
  if len(a) != len(b) {
    return false
  }
  letters := [26]int{}
  for _, l := range a {
    letters[l-'a']++
  }
  for _, l := range b {
    letters[l-'a']--  
  }
  for _, v := range letters {
    if v != 0 {
      return false
    }
  }
  return true
}

func main() {
  a, b := "heart", "earth"
  fmt.Printf("%s, %s = %t\n", a, b, is_anagram(a, b))

  a, b = "dog", "cat"
  fmt.Printf("%s, $%s = %t\n", a, b, is_anagram(a, b))
}
