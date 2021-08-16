package main

import "fmt"

const NULL_CHAR = '0'
var pairs = map[rune]rune{
  '}': '{',
  ']': '[',
  ')': '(',
  '{': NULL_CHAR,
  '[': NULL_CHAR,
  '(': NULL_CHAR,
}

func is_balanced(s string) bool {
  var stack []rune
  for _, c := range s {
      opening_pair, ok := pairs[c]
      if !ok {
        continue
      }
      if opening_pair == NULL_CHAR {
        stack = append(stack, c)
        continue
      }
      if len(stack) == 0 || stack[len(stack)-1] != opening_pair {
        return false
      }
      stack = stack[:len(stack)-1]
  }
  return true
}

func main() {

  fmt.Println(is_balanced("(())"))
  fmt.Println(is_balanced("[((){}){}[]]"))
  fmt.Println(is_balanced("}((){}){}[]]"))
  fmt.Println(is_balanced("(())())"))
}
