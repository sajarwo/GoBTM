package gobtm

import(
	//"time"
	"math/rand"
)
func Uni_sample(k int)int{
	//r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return rand.Intn(k)
}
func Mult_sample(K int, p []float64)int{
	for i:=1;i<K;i++{
		p[i]+=p[i-1]
	}
	//r := rand.New(rand.NewSource(time.Now().UnixNano()))
	u:= rand.Float64()

	k:=0
	for ;k<K;k++{
		if p[k] >= u*p[K-1]{
			break
		}
	}
	if k==K{
		k--
	}
	return k
}

func Normalize(p *[]float64, smoother float64){
	s := Sum(*p)
	K := len(*p)
	// avoid numerical problem
	for i:=0; i<K; i++ {
		(*p)[i] = ((*p)[i] + smoother)/(s + float64(K)*smoother);
	}
}

func Sum(p []float64)float64{
	s:=0.0
	for i:=0; i<len(p); i++ {
		s+=p[i]
	}
	return s
}