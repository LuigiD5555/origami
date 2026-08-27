package glyphcalc

import "fmt"

func Generate(pattern string, length int, seed uint64) ([]byte, error) {
	if length <= 0 { return nil, fmt.Errorf("length must be > 0") }
	out := make([]byte,0,length)
	switch pattern {
	case "runs": vals:=[]byte{0x11,0x44,0xAA,0xF0}; for len(out)<length{v:=vals[(len(out)/64)%len(vals)];for i:=0;i<64&&len(out)<length;i++{out=append(out,v)}}
	case "sequence": for len(out)<length{for i:=0;i<128&&len(out)<length;i++{out=append(out,byte(i*3+7))}}
	case "motif": motif:=[]byte("ORIGAMI-HYPERFOLD|"); for len(out)<length{need:=length-len(out);if need>len(motif){need=len(motif)};out=append(out,motif[:need]...)}
	case "mixed": for len(out)<length{phase:=(len(out)/128)%4;switch phase{case 0:for i:=0;i<128&&len(out)<length;i++{out=append(out,0xCC)};case 1:for i:=0;i<128&&len(out)<length;i++{out=append(out,byte(5+i*2))};case 2:m:=[]byte("OHF-CALC|");for i:=0;i<128&&len(out)<length;i++{out=append(out,m[i%len(m)])};case 3:base:=len(out)-128;if base<0{base=0};chunk:=append([]byte(nil),out[base:]...);if len(chunk)==0{chunk=[]byte{0}};for i:=0;i<128&&len(out)<length;i++{out=append(out,chunk[i%len(chunk)])}}}
	case "random": x:=seed;for len(out)<length{x+=0x9e3779b97f4a7c15;z:=x;z=(z^(z>>30))*0xbf58476d1ce4e5b9;z=(z^(z>>27))*0x94d049bb133111eb;z^=z>>31;for j:=0;j<8&&len(out)<length;j++{out=append(out,byte(z>>(8*j)))}}
	default:return nil,fmt.Errorf("unknown pattern %q (runs|sequence|motif|mixed|random)",pattern)
	}
	return out,nil
}
