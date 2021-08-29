package main

import "fmt"

type location struct {
  word string
  depth int
}

func ladderLength(beginWord string, endWord string, wordList []string) int {
    buckets := map[string][]string{}

    for _, word := range wordList {
        for i := 0; i < len(word); i++ {
            stubbed := word[:i] + "_" + word[i+1:]
            buckets[stubbed] = append(buckets[stubbed], word)
        }
    }

    queue := []location{{beginWord, 1}}
    seen := map[string]bool{}

    for len(queue) > 0 {
      loc := queue[0]
      queue = queue[1:]

      if loc.word == endWord {
        return loc.depth
      }

      seen[loc.word] = true
      for i := 0; i < len(loc.word); i++ {
          stubbed := loc.word[:i] + "_" + loc.word[i+1:]
          nextWords := buckets[stubbed]
          for _, nw := range nextWords {
            if seen[nw] {
              continue
            }
            queue = append(queue, location{nw, loc.depth + 1})
          }
      }
    }
    return 0
}

func main() {
  beginWord := "hit"
  endWord := "cog"
  wordList := []string{"hot","dot","dog","lot","log"}
  fmt.Println(ladderLength(beginWord, endWord, wordList))
}
