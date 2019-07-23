package keygen

import (
	"math"
	"math/rand"
	"time"
)

//注意，生成的f其1系数个数比-1系数多一个
func Random_f(n,d int) []int {
	d = d+d-1
	j := -1
	var i int
	var f = make([]int, n)
	for i=0;i<n;i++ {f[i]=0}

	//根据时间生成不同的随机源，根据随机源生成不同的随机数
	random := rand.New(rand.NewSource(time.Now().Unix()))
	k := 0

	for k<d {
		i = int(math.Floor(float64(n * random.Intn(100) / 100)))
		if f[i] == 0 {
			j = -j
			f[i] = j
			k++
		}
	}


	return f
}
