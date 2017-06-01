// comparetest.go
package test

import (
	// "../lib/xx"
	"../card"
	// "../lib/xxio"
	// "encoding/json"
	"fmt"
	// "io/ioutil"
	// "os"
	// "strconv"
)

// var card_cls = []string{"diamond", "club", "heart", "spade"}

func CompareTest() {
	// var cards1, cards2 Cards
	var c [5]int
	for {
		cards1 := card.NewCards()
		cards2 := card.NewCards()
		fmt.Printf("input 5 cards:\n\t")
		fmt.Scanf("%d%d%d%d%d\n", &c[0], &c[1], &c[2], &c[3], &c[4])
		for j := 0; j < 5; j++ {
			val := c[j]
			cls := card_cls[val/100]
			idx := val % 100
			card := card.New(cls, idx)
			cards1.Append(card)
		}
		fmt.Printf("input another 5 cards:\n\t")
		fmt.Scanf("%d%d%d%d%d\n", &c[0], &c[1], &c[2], &c[3], &c[4])
		for j := 0; j < 5; j++ {
			val := c[j]
			cls := card_cls[val/100]
			idx := val % 100
			card := card.New(cls, idx)
			cards2.Append(card)
		}

		res := cards1.Compare(cards2)

		fmt.Println("The result is :", res)
	}
}
