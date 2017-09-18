package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"sort"

	log "github.com/sirupsen/logrus"
)

var (
	inFileName  = flag.String("in_filename", "./sm6.txt", "path to file =)")
	inFileName2 = flag.String("in_filename2", "./sm7.txt", "path to file =)")
	outFileName = flag.String("out_filename", "./sm8.txt", "path to file =)")
)

type (
	StateMachine struct {
		WordsCount   int
		StartState   int
		FinishStates []int
		Graph        [][]Edge
	}

	Edge struct {
		V      int
		Symbol int
	}
)

// Minimaize splits k-equivalence class into k+1-equivalence classes if it's possible.
// It takes k-equivalence classes of the StateMachine
func (m StateMachine) Minimaize() StateMachine {
	var (
		n         = len(m.Graph)
		csk       = make([]int, n)
		csk1      = make([]int, n)
		colorsOld int
		colors    = 2
		moves     = make([][]int, n)
	)

	for _, v := range m.FinishStates {
		csk1[v] = 1
	}

	for u, es := range m.Graph {
		moves[u] = make([]int, m.WordsCount)
		for i := range moves[u] {
			moves[u][i] = -1
		}
		for _, e := range es {
			moves[u][e.Symbol] = e.V
		}
	}

	for colors != colorsOld {
		colorsOld, colors = colors, 0
		csk, csk1 = csk1, csk

		//fmt.Println(colorsOld, csk)

		for i := range csk1 {
			csk1[i] = -1
		}

		for u := range m.Graph {
			if csk1[u] > -1 {
				continue
			}

			csk1[u] = colors
			colors++

			for v := u + 1; v < n; v++ {
				if csk[u] != csk[v] {
					continue
				}

				var eq = true
				for s := 0; s < m.WordsCount; s++ {
					v1, v2 := moves[u][s], moves[v][s]
					//fmt.Println(v1, v2)
					if v1 == -1 {
						if v2 != -1 {
							eq = false
							break
						}
					} else {
						if v2 == -1 {
							eq = false
							break
						} else if csk[v1] != csk[v2] {
							eq = false
							break
						}
					}
				}
				//fmt.Println(u, v, eq)

				if eq {
					csk1[v] = csk1[u]
				}
			}
		}
	}

	var m2 StateMachine
	m2.StartState = csk1[m.StartState]
	m2.WordsCount = m.WordsCount
	m2.Graph = make([][]Edge, colors)

	for u, es := range m.Graph {
		for _, e := range es {
			m2.Graph[csk1[u]] = append(m2.Graph[csk1[u]], Edge{V: csk1[e.V], Symbol: e.Symbol})
		}
	}

	for _, s := range m.FinishStates {
		m2.FinishStates = append(m2.FinishStates, csk1[s])
	}

	return m2
}

// Read reads StateMachine from io.Reader
func (m *StateMachine) Read(r io.Reader) {
	var n int

	fmt.Fscan(r, &m.WordsCount)

	fmt.Fscan(r, &n)
	m.Graph = make([][]Edge, n)

	fmt.Fscan(r, &m.StartState)

	fmt.Fscan(r, &n)
	m.FinishStates = make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(r, &m.FinishStates[i])
	}

	for {
		var (
			u int
			e Edge
		)

		_, err := fmt.Fscanf(r, "%d %d %d", &u, &e.Symbol, &e.V)
		if err == io.EOF {
			break
		}

		if err != nil {
			log.Error(err)
		}

		m.Graph[u] = append(m.Graph[u], e)
	}
}

// Write writes StateMachine into io.Writer
func (m *StateMachine) Write(w io.Writer) {
	fmt.Fprintln(w, m.WordsCount)
	fmt.Fprintln(w, len(m.Graph))
	fmt.Fprintln(w, m.StartState)
	fmt.Fprint(w, len(m.FinishStates))
	for _, s := range m.FinishStates {
		fmt.Fprintf(w, " %d", s)
	}
	fmt.Fprintln(w)
	for u, edges := range m.Graph {
		for _, edge := range edges {
			fmt.Fprintln(w, u, edge.Symbol, edge.V)
		}
	}
}

func (m StateMachine) Determine() StateMachine {
	var m2 StateMachine
	m2.Graph = make([][]Edge, encodeOne(len(m.Graph)))
	m2.WordsCount = m.WordsCount
	m2.StartState = encode(m.StartState)

	for i, a := 1, encode(m.FinishStates...)+1; i <= len(m2.Graph); i++ {
		if i&a > 0 {
			m2.FinishStates = append(m2.FinishStates, i-1)
		}
	}

	for nu, n := 0, len(m2.Graph); nu < n; nu++ {
		var edges []int = make([]int, m2.WordsCount)
		for _, u := range decode(nu) {
			for _, e := range m.Graph[u] {
				edges[e.Symbol] |= 1 << uint(e.V)
			}
		}

		for symbol, nv := range edges {
			if nv == 0 {
				continue
			}
			m2.Graph[nu] = append(m2.Graph[nu], Edge{Symbol: symbol, V: nv - 1})
		}
	}

	return m2
}

func (m StateMachine) RemoveUnreachable() StateMachine {
	var color []int = make([]int, len(m.Graph))
	var q []int = []int{m.StartState}
	color[m.StartState] = 1

	//fmt.Println(q)
	//fmt.Println(m.Graph)
	for ; len(q) > 0; q = q[1:] {
		s := q[0]
		//fmt.Println(s)
		//fmt.Println(m.Graph[s])
		for _, e := range m.Graph[s] {
			//fmt.Println(" ", e.V)
			if color[e.V] == 0 {
				q = append(q, e.V)
				//fmt.Println(" ", q, q[1:])
				color[e.V] = 1
			}
		}
	}

	var n int
	for k, v := range color {
		if v == 0 {
			color[k] = -1
		} else {
			color[k] = n
			n++
		}
	}

	//fmt.Println(n, color)

	var m2 StateMachine

	m2.StartState = color[m.StartState]
	m2.WordsCount = m.WordsCount
	for _, s := range m.FinishStates {
		if c := color[s]; c != -1 {
			m2.FinishStates = append(m2.FinishStates, c)
		}
	}

	m2.Graph = make([][]Edge, n)
	for u, edges := range m.Graph {
		if nu := color[u]; nu != -1 {
			for _, e := range edges {
				if nv := color[e.V]; nv != -1 {
					m2.Graph[nu] = append(m2.Graph[nu], Edge{Symbol: e.Symbol, V: nv})
				}
			}
		}
	}

	//m2.OrganizeEdges()
	return m2
}

func (m *StateMachine) OrganizeEdges() {
	sort.Ints(m.FinishStates)
	n := 1
	for i := range m.FinishStates {
		if m.FinishStates[i] != m.FinishStates[n-1] {
			m.FinishStates[n] = m.FinishStates[i]
			n++
		}
	}
	m.FinishStates = m.FinishStates[:n]

	for u, es := range m.Graph {
		if len(es) == 0 {
			continue
		}

		sort.Slice(es, func(i, j int) bool {
			if es[i].Symbol < es[j].Symbol {
				return true
			} else if es[i].Symbol == es[j].Symbol && es[i].V < es[j].V {
				return true
			}

			return false
		})

		n := 1
		for i := range es {
			if es[i] != es[n-1] {
				es[n] = es[i]
				n++
			}
		}
		m.Graph[u] = es[:n]
	}
}

func Equivalent(m1, m2 StateMachine) bool {
	if len(m1.Graph) != len(m2.Graph) ||
		len(m1.FinishStates) != len(m2.FinishStates) ||
		m1.WordsCount != m2.WordsCount {
		return false
	}

	var trans = make([]int, len(m1.Graph))
	for i := range trans {
		trans[i] = -1
	}

	var q []int = []int{m1.StartState}
	trans[m1.StartState] = m2.StartState

	var moves1 = make([][]int, len(m1.Graph))
	for u, es := range m1.Graph {
		moves1[u] = make([]int, m1.WordsCount)
		for i := range moves1[u] {
			moves1[u][i] = -1
		}
		for _, e := range es {
			moves1[u][e.Symbol] = e.V
		}
	}
	var moves2 = make([][]int, len(m2.Graph))
	for u, es := range m2.Graph {
		moves2[u] = make([]int, m2.WordsCount)
		for i := range moves2[u] {
			moves2[u][i] = -1
		}
		for _, e := range es {
			moves2[u][e.Symbol] = e.V
		}
	}

	for ; len(q) > 0; q = q[1:] {
		u := q[0]

		//fmt.Println(u, trans)

		for w := 0; w < m1.WordsCount; w++ {
			//fmt.Printf(" %d %d\n", u, w)
			c1, c2 := moves1[u][w], moves2[trans[u]][w]

			if c1 != -1 {
				if c2 == -1 {
					fmt.Println("E1")
					return false
				}

				if trans[c1] == -1 {
					q = append(q, c1)
					trans[c1] = c2
				} else if trans[c1] != c2 {
					fmt.Println("E2")
					return false
				}
			} else {
				if c2 != -1 {
					fmt.Println("E3")
					return false
				}
			}

			//if c1 != -1 && c2 == -1 || c1 == -1 && c2 != -1 || c1 != -1 && c2 != -1 && trans[c1] != c2 {
			//	fmt.Println(c1, c2, trans[c1])
			//	return false
			//}
		}
	}

	return true
}

func encode(s ...int) int {
	var ans int
	for _, a := range s {
		ans |= int(1 << uint(a))
	}

	return ans - 1
}

func encodeOne(a int) int {
	return (1 << uint(a)) - 1
}

func decode(a int) []int {
	var ans []int

	a += 1
	for i := 0; i < 32; i++ {
		if a&(1<<uint(i)) > 0 {
			ans = append(ans, i)
		}
	}

	return ans
}

func main() {
	var (
		m1, m2   StateMachine
		r, r2, w *os.File
		err      error
	)

	flag.Parse()

	r, err = os.Open(*inFileName)
	if err != nil {
		log.Error(err)
		return
	}

	r2, err = os.Open(*inFileName2)
	if err != nil {
		log.Error(err)
		return
	}

	w, err = os.Create(*outFileName)
	if err != nil {
		log.Error(err)
		return
	}

	m1.Read(r)
	m2.Read(r2)

	m1 = m1.Determine()
	//m1.Write(w)
	//return

	m1 = m1.RemoveUnreachable()
	//m1.Write(w)
	//return
	m1 = m1.Minimaize()
	m1.OrganizeEdges()
	//m1.Write(w)
	//return

	m2 = m2.Determine()
	m2 = m2.RemoveUnreachable()
	m2 = m2.Minimaize()
	m2.OrganizeEdges()

	m1.Write(w)
	w.Write([]byte("\n"))
	m2.Write(w)

	fmt.Println(Equivalent(m1, m2))
}
