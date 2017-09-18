package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"io"
	"os"
	"sort"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
)

var (
	//inFileName  = flag.String("filename", "./example.txt", "path to file =)")
	inFileName  = flag.String("in_filename", "./example_light.txt", "path to file =)")
	outFileName = flag.String("out_filename", "./sets.txt", "path to file =)")
	lineQueue   chan []byte
	wordQueue   chan []byte
	wg          sync.WaitGroup
	words       [][]byte
	cnt         int
)

func main() {
	flag.Parse()

	dt := time.Now()

	r, err := os.Open(*inFileName)
	if err != nil {
		panic(err)
	}

	br := bufio.NewReader(r)

	lineQueue = make(chan []byte, 100)
	wordQueue = make(chan []byte, 100)

	for i := 0; i < 3; i++ {
		go func() {
			for line := range lineQueue {
				parseLine(line)

				wg.Done()
			}
		}()
	}

	go wordProcessor()

	for {
		line, _, err := br.ReadLine()
		if err == io.EOF {
			break
		}
		wg.Add(1)

		newLine := make([]byte, 0, len(line))
		newLine = append(newLine, line...)

		lineQueue <- newLine
	}

	wg.Wait()

	sort.Slice(words, func(i, j int) bool {
		for k, c := range words[i] {
			if len(words[j]) <= k {
				return false
			}
			if c != words[j][k] {
				return c < words[j][k]
			}
		}
		return true
	})

	//fmt.Println(len(words))
	words = removeDuplicates(words)
	//fmt.Println(len(words))

	var maxLen int = -1
	var g [][]int
	var superWords map[int]bool = map[int]bool{}

	for i1, w1 := range words {
		for i2, w2 := range words {
			if i1 <= i2 {
				continue
			}
			if d := dist(w1, w2); d > maxLen {
				log.Info(d)
				// GC forever =)
				superWords = map[int]bool{}
				g = make([][]int, len(words))
				maxLen = d

				superWords[i1] = true
				superWords[i2] = true
				g[i1] = append(g[i1], i2)
				g[i2] = append(g[i2], i1)
			} else if d == maxLen {
				superWords[i1] = true
				superWords[i2] = true
				g[i1] = append(g[i1], i2)
				g[i2] = append(g[i2], i1)
			}
		}
	}

	//fmt.Println(maxLen)
	//fmt.Println(len(superWords))

	// Graph search
	w, err := os.Create(*outFileName)
	if err != nil {
		log.Error(err)
		return
	}
	defer w.Close()

	bw := bufio.NewWriter(w)
	defer bw.Flush()

	//s := 0
	//m := len(g[0])
	//i := 0
	//for v, d := range g {
	//	if len(d) > m {
	//		m = len(d)
	//		i = v
	//	}
	//	s += len(d)
	//}

	for _, d := range g {
		sort.Ints(d)
	}

	check := func(v int, sg []int) bool {
		for _, u := range sg {
			if dist(words[v], words[u]) < maxLen {
				return false
			}
		}
		return true
	}

	var build func([]int)

	build = func(sg []int) {
		from := sg[len(sg)-1]
		for _, to := range g[from] {
			if to <= from {
				continue
			}
			if check(to, sg) {
				sg = append(sg, to)
				build(sg)
				sg = sg[:len(sg)-1]
			}
		}
		if len(sg) > 1 {
			var w []string
			for _, v := range sg {
				w = append(w, string(words[v]))
			}
			fmt.Fprintln(bw, w)
		}
	}

	for v := range g {
		build([]int{v})
	}

	//fmt.Println(s)
	//fmt.Println(string(words[i]))
	//fmt.Println()
	//for _, j := range g[i] {
	//	fmt.Println(string(words[j]))
	//}
	//fmt.Println(m)

	log.Info(time.Now().Sub(dt))
}

func removeDuplicates(slice [][]byte) [][]byte {
	n := 0
	for i, word := range slice {
		if i == 0 {
			n++
			continue
		}
		if !bytes.Equal(word, slice[i-1]) {
			slice[n] = word
			n++
		}
	}
	return slice[:n]
}

func wordProcessor() {
	for word := range wordQueue {
		word = bytes.ToLower(word)
		words = append(words, word)
		cnt++

		//fmt.Println(string(word))
		wg.Done()
	}
}

func parseLine(line []byte) {
	start := 0
	end := 0
	for i, c := range line {
		var symbol bool
		if c >= 'a' && c <= 'z' || c >= 'A' && c <= 'Z' {
			symbol = true
		}
		if symbol {
			end++
		} else {
			if start < end {
				wg.Add(1)
				wordQueue <- line[start:end]
			}
			start = i + 1
			end = i + 1
		}
	}
}

func dist(a, b []byte) int {
	if len(a) > len(b) {
		a, b = b, a
	}

	var d int = len(b) - len(a)

	for i, c := range a {
		if b[i] != c {
			d++
		}
	}

	return d
}
