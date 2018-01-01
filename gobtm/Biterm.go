package gobtm

import "fmt"

type Biterm struct{
	Wi 	int
	Wj 	int
	Z 	int
}

func (w *Biterm)Create(w1, w2 int){
	w.Wi = Min(w1, w2);
	w.Wj = Max(w1, w2);
	w.Z  = -1
}

func (w *Biterm)Get_wi()int{return w.Wi}
func (w *Biterm)Get_wj()int{return w.Wj}
func (w *Biterm)Get_z()int{return w.Z}
func (w *Biterm)Set_z(z int){w.Z=z}
func (w *Biterm)Reset_z(){ w.Z=-1}
func (w *Biterm)ToKey()string{ return fmt.Sprintf("%d%d",w.Wi,w.Wj)}

func Min(a, b int)int{
	if a>b{
		a=b
	}
	return a
}

func Max(a, b int)int{
	if a<b{
		a=b
	}
	return a
}