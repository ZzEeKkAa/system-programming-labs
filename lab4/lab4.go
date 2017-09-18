package main

import (
	"fmt"
	"io"
	"math/rand"
	"sync"
	"time"
)

type Matrix [][]int

func (mx *Matrix) Read(r io.Reader) {
	var n, m int
	fmt.Fscanf(r, "%d %d", &n, &m)

	*mx = make([][]int, n)

	for i := range *mx {
		a := make([]int, m)

		for j := range a {
			fmt.Fscanf(r, "%d", &a[j])
		}

		(*mx)[i] = a
	}
}

func (mx *Matrix) Generate(n, m int) {
	*mx = make([][]int, n)

	for i := range *mx {
		a := make([]int, m)

		for j := range a {
			a[j] = rand.Intn(20) - 10
		}

		(*mx)[i] = a
	}
}

func (mx *Matrix) New(n, m int) {
	*mx = make([][]int, n)

	for i := range *mx {
		(*mx)[i] = make([]int, m)
	}
}

func (mx *Matrix) N() int {
	return len(*mx)
}

func (mx *Matrix) M() int {
	if len(*mx) == 0 {
		return 0
	}

	return len((*mx)[0])
}

func (mx Matrix) String() string {
	var s string
	for _, r := range mx {
		s += fmt.Sprintln(r)
	}
	return s
}

func main() {
	var A, B Matrix
	A.Generate(900, 600)
	B.Generate(600, 1200)

	dt := time.Now()
	C := multiplyOne(A, B)
	fmt.Printf("%v\t%v\t%v\n", time.Now().Sub(dt), C.N(), C.M())

	dt = time.Now()
	C = multiplyMulti(A, B)
	fmt.Printf("%v\t%v\t%v\n", time.Now().Sub(dt), C.N(), C.M())

	dt = time.Now()
	C = multiplyWorkers(A, B, 2)
	fmt.Printf("%v\t%v\t%v\n", time.Now().Sub(dt), C.N(), C.M())

	dt = time.Now()
	C = multiplyWorkers(A, B, 3)
	fmt.Printf("%v\t%v\t%v\n", time.Now().Sub(dt), C.N(), C.M())

	dt = time.Now()
	C = multiplyWorkers(A, B, 4)
	fmt.Printf("%v\t%v\t%v\n", time.Now().Sub(dt), C.N(), C.M())

	dt = time.Now()
	C = multiplyWorkers(A, B, 5)
	fmt.Printf("%v\t%v\t%v\n", time.Now().Sub(dt), C.N(), C.M())

	dt = time.Now()
	C = multiplyWorkers(A, B, 6)
	fmt.Printf("%v\t%v\t%v\n", time.Now().Sub(dt), C.N(), C.M())

	dt = time.Now()
	C = multiplyWorkers(A, B, 7)
	fmt.Printf("%v\t%v\t%v\n", time.Now().Sub(dt), C.N(), C.M())

	dt = time.Now()
	C = multiplyWorkers(A, B, 8)
	fmt.Printf("%v\t%v\t%v\n", time.Now().Sub(dt), C.N(), C.M())

	dt = time.Now()
	C = multiplyWorkers(A, B, 9)
	fmt.Printf("%v\t%v\t%v\n", time.Now().Sub(dt), C.N(), C.M())

	dt = time.Now()
	C = multiplyWorkers(A, B, 10)
	fmt.Printf("%v\t%v\t%v\n", time.Now().Sub(dt), C.N(), C.M())

	//fmt.Println(A)
	//fmt.Println(B)
	//fmt.Println(C)
}

func multiplyOne(A, B Matrix) Matrix {
	var (
		C       Matrix
		n, k, m = A.N(), A.M(), B.M()
	)
	if k != B.N() {
		return nil
	}
	C.New(n, m)

	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			for l := 0; l < k; l++ {
				C[i][j] += A[i][l] * B[l][j]
			}
		}
	}

	return C
}

func multiplyMulti(A, B Matrix) Matrix {
	var (
		wg      sync.WaitGroup
		C       Matrix
		n, k, m = A.N(), A.M(), B.M()
	)
	if k != B.N() {
		return nil
	}
	C.New(n, m)

	wg.Add(n * m)
	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			go func(i, j int) {
				for l := 0; l < k; l++ {
					C[i][j] += A[i][l] * B[l][j]
				}
				wg.Done()
			}(i, j)
		}
	}
	wg.Wait()

	return C
}

func multiplyWorkers(A, B Matrix, workersCount int) Matrix {
	var (
		wg      sync.WaitGroup
		C       Matrix
		n, k, m = A.N(), A.M(), B.M()
	)
	if k != B.N() {
		return nil
	}
	C.New(n, m)

	wg.Add(workersCount)

	for w, cnt := 0, n/workersCount; w < workersCount; w++ {
		var s, e = w * cnt, (w + 1) * cnt
		if w == workersCount {
			e = n
		}
		go func(s, e int) {
			for i := s; i < e; i++ {
				for j := 0; j < m; j++ {
					for l := 0; l < k; l++ {
						C[i][j] += A[i][l] * B[l][j]
					}
				}
			}
			wg.Done()
		}(s, e)
	}
	wg.Wait()

	return C
}
