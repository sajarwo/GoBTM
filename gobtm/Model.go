package gobtm

import(
	"fmt"
	"time"
	"os"
	"bufio"
	"strings"
	"strconv"
)

type BTM struct{
	K		int						// number of topics
	W 		float64					// vocabulary size
	Iter 	int						// maximum number of iteration of Gibbs Sampling
	Alpha	float64					// hyperparameters of p(z)
	Beta 	float64					// hyperparameters of p(w|z)
	Pw_b	[]float64				// word distribution TF (term-freq) / dictionary term
	Nb_z	[]float64				// n(b|z), size K*1 (theta)
	Nw_z	[][]float64				// n(w,z), size K*W (phi)
	BS 		map[string]*Biterm		// biterm distribution
}


func Run(){

	dir:= "set path to your work-dir " // ex. C://output/
	
	Btm := BTM{}
	Btm.Alpha = 1.0
	Btm.Beta = 0.01
	Btm.Iter = 5
	Btm.K = 20
	Btm.W = 28634

	Btm.BS = make(map[string]*Biterm)
	Btm.Pw_b = make([]float64,int(Btm.W))
	Btm.Nb_z = make([]float64,Btm.K)
	Btm.Nw_z = make([][]float64,Btm.K)
	for i:=0;i<Btm.K;i++{
		Btm.Nw_z[i]=make([]float64,int(Btm.W))
	}	
	
	t := time.Now()
	Btm.Load_docs(dir,"doc_wids.txt")
	Btm.Model_init()
	for i:=0;i<Btm.Iter;i++{
		fmt.Println("Iter = ", i)
		Btm.Update_biterm()
	}
	
	Btm.Save_Nb_z(dir,"go.pz")
	Btm.Save_Pw_z(dir,"go.pw_z")
	fmt.Println(time.Since(t))
}

func (Btm *BTM)Load_docs(dir, src string){
	fmt.Println("Loading docs...")
	f, _ := os.Open(dir+src)
	rl := bufio.NewReader(f)
	defer f.Close();
	line ,_, err := rl.ReadLine()
	for err == nil {
		Btm.GenBiterm(string(line),15)
		line ,_, err = rl.ReadLine()
	}
	Normalize(&Btm.Pw_b,0.0)
}

func (Btm *BTM)Model_init(){
	for _,v:=range Btm.BS{
		k:=Uni_sample(Btm.K)
		Btm.Assign_biterm_topic(k,v)
	}
}


func (Btm *BTM)Update_biterm(){
	for _,v:=range Btm.BS{
		Btm.Reset_biterm_topic(v)
		pz:=Btm.Compute_pz_b(v)
		k:=Mult_sample(Btm.K, pz)
		Btm.Assign_biterm_topic(k,v)
	}
}

func (Btm *BTM)Reset_biterm_topic(bi *Biterm){
	k := bi.Get_z()
	w1 := bi.Get_wi();
	w2 := bi.Get_wj();
	Btm.Nb_z[k] -= 1;
	Btm.Nw_z[k][w1] -= 1;
	Btm.Nw_z[k][w2] -= 1;
	bi.Set_z(-1)
}

func (Btm *BTM)Compute_pz_b(bi *Biterm)[]float64{
	pz:=make([]float64,Btm.K)
	w1 := bi.Get_wi();
	w2 := bi.Get_wj();
	var pw1k, pw2k, pk float64

	for k:=0;k<Btm.K;k++{
		if k == 0{
			pw1k = Btm.Pw_b[w1]
			pw2k = Btm.Pw_b[w2]
		}else{
			pw1k = (Btm.Nw_z[k][w1] + Btm.Beta) / (2 * Btm.Nb_z[k] + Btm.W * Btm.Beta);
			pw2k = (Btm.Nw_z[k][w2] + Btm.Beta) / (2 * Btm.Nb_z[k] + 1 + Btm.W * Btm.Beta);
		}
		pk = (Btm.Nb_z[k] + Btm.Alpha) / (float64(len(Btm.BS) + Btm.K) * Btm.Alpha);
		pz[k] = pk * pw1k * pw2k;
	}
	return pz
}

func (Btm *BTM)Assign_biterm_topic(k int, bi *Biterm){
	bi.Set_z(k);
	w1 := bi.Get_wi();
	w2 := bi.Get_wj();
	Btm.Nb_z[k] += 1;
	Btm.Nw_z[k][w1] += 1;
	Btm.Nw_z[k][w2] += 1;
}


func (Btm *BTM)GenBiterm(document string, win int){
	doc:=strings.Fields(document)
	i:=0
	for ;i<len(doc)-1;i++{
		w0,_:=strconv.Atoi(doc[i])
		for j:=i+1;j <Min(i+win, len(doc));j++{
			w1,_:=strconv.Atoi(doc[j])
			bi:=Biterm{}
			bi.Create(w0, w1)
			Btm.BS[bi.ToKey()]=&bi
		}
		Btm.Pw_b[w0]+=1
	}
	w0,_:=strconv.Atoi(doc[i])
	Btm.Pw_b[w0]+=1
}


func (Btm *BTM)Save_Nb_z(dir, fname string){
	fmt.Println("\nNb_z =\n", Btm.Nb_z)
	f, err := os.Create(dir+fname)
	defer f.Close()
	if(err == nil) {
		Normalize(&Btm.Nb_z, Btm.Alpha)
		for _,v := range Btm.Nb_z{
			a:=fmt.Sprintf("%.6f ",v)
			f.WriteString(a)
		}
	}
	fmt.Println("\nNb_z norm =\n",Btm.Nb_z)
}

func (Btm *BTM)Save_Pw_z(dir, fname string){
	fmt.Println("\nPw_z = ")
	f, err := os.Create(dir+fname)
	defer f.Close()
	if(err == nil) {
		Normalize(&Btm.Nb_z, Btm.Alpha)
		for k:=0; k< Btm.K;k++{
			for w:=0;w<int(Btm.W);w++{
				pwz:= (Btm.Nw_z[k][w] + Btm.Beta) / (Btm.Nb_z[k] * 2 + Btm.W * Btm.Beta);
				a:=fmt.Sprintf("%.6f ",pwz)
				f.WriteString(a)
				if w<50{
					fmt.Print(a)
				}
			}
			fmt.Println("\n")
		}
	}
	fmt.Println(dir+fname)
}



