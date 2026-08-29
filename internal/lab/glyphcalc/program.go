package glyphcalc

import (
	"bytes"
	"fmt"
)

const (
	TransportProgram byte = 0xA0
	TransportRaw     byte = 0xA1

	OpHalt byte = 0x00
	OpLit  byte = 0x01
	OpRun  byte = 0x02
	OpRef  byte = 0x03
	OpSeq  byte = 0x04
)

type Stats struct {
	Instructions int `json:"instructions"`
	LitOps       int `json:"lit_ops"`
	RunOps       int `json:"run_ops"`
	RefOps       int `json:"ref_ops"`
	SeqOps       int `json:"seq_ops"`
}

type Transport struct {
	Mode        string `json:"mode"`
	Stream      []byte `json:"-"`
	Program     []byte `json:"-"`
	ProgramStat Stats  `json:"program_stats"`
}

func CompileGreedy(src []byte) ([]byte, Stats) {
	var p []byte
	var st Stats
	for i := 0; i < len(src); {
		runLen := detectRun(src, i)
		seqLen, delta := detectSeq(src, i)
		refDist, refLen := detectRef(src, i)
		best := "lit"
		bestLen := 1
		if runLen >= 4 && runLen > bestLen { best, bestLen = "run", runLen }
		if seqLen >= 5 && seqLen > bestLen { best, bestLen = "seq", seqLen }
		if refLen >= 4 && refLen > bestLen { best, bestLen = "ref", refLen }
		switch best {
		case "run":
			p = append(p, OpRun, src[i], byte(runLen)); st.RunOps++; i += runLen
		case "seq":
			p = append(p, OpSeq, src[i], delta, byte(seqLen)); st.SeqOps++; i += seqLen
		case "ref":
			p = append(p, OpRef, byte(refDist), byte(refLen)); st.RefOps++; i += refLen
		default:
			p = append(p, OpLit, src[i]); st.LitOps++; i++
		}
		st.Instructions++
	}
	p = append(p, OpHalt)
	st.Instructions++
	return p, st
}

func CompileBest(src []byte) (Transport, error) {
	if len(src) > 65535 { return Transport{}, fmt.Errorf("source too large for glyph-calculus v0: %d > 65535", len(src)) }
	prog, st := CompileGreedy(src)
	programStream := append([]byte{TransportProgram}, prog...)
	rawStream := make([]byte, 3+len(src))
	rawStream[0] = TransportRaw
	rawStream[1] = byte(len(src) >> 8)
	rawStream[2] = byte(len(src))
	copy(rawStream[3:], src)
	if len(programStream) < len(rawStream) { return Transport{Mode:"program", Stream:programStream, Program:prog, ProgramStat:st}, nil }
	return Transport{Mode:"raw", Stream:rawStream, Program:prog, ProgramStat:st}, nil
}

func ExecuteTransport(stream []byte, maxOutput int, maxOps int) ([]byte, Stats, error) {
	if len(stream)==0 { return nil, Stats{}, fmt.Errorf("empty transport") }
	switch stream[0] {
	case TransportRaw:
		if len(stream)<3 { return nil,Stats{},fmt.Errorf("short raw header") }
		n:=int(stream[1])<<8|int(stream[2])
		if n>maxOutput { return nil,Stats{},fmt.Errorf("raw output %d exceeds budget %d",n,maxOutput) }
		if len(stream)!=n+3 { return nil,Stats{},fmt.Errorf("raw length mismatch: header=%d payload=%d",n,len(stream)-3) }
		return append([]byte(nil),stream[3:]...),Stats{},nil
	case TransportProgram:
		return ExecuteProgram(stream[1:],maxOutput,maxOps)
	default:
		return nil,Stats{},fmt.Errorf("unknown transport header 0x%02x",stream[0])
	}
}

func ExecuteProgram(p []byte, maxOutput int, maxOps int) ([]byte, Stats, error) {
	var out []byte
	var st Stats
	for pc:=0;; {
		if st.Instructions>=maxOps { return nil,st,fmt.Errorf("operation budget exceeded: %d",maxOps) }
		if pc>=len(p) { return nil,st,fmt.Errorf("program ended without HALT") }
		op:=p[pc];pc++;st.Instructions++
		switch op {
		case OpHalt:return out,st,nil
		case OpLit:
			if pc>=len(p){return nil,st,fmt.Errorf("LIT missing operand")};out=append(out,p[pc]);pc++;st.LitOps++
		case OpRun:
			if pc+1>=len(p){return nil,st,fmt.Errorf("RUN missing operands")};v,n:=p[pc],int(p[pc+1]);pc+=2;if n==0{return nil,st,fmt.Errorf("RUN count 0")};if len(out)+n>maxOutput{return nil,st,fmt.Errorf("RUN exceeds output budget")};out=append(out,bytes.Repeat([]byte{v},n)...);st.RunOps++
		case OpRef:
			if pc+1>=len(p){return nil,st,fmt.Errorf("REF missing operands")};dist,n:=int(p[pc]),int(p[pc+1]);pc+=2;if dist==0||n==0||dist>len(out)||n>dist{return nil,st,fmt.Errorf("invalid REF dist=%d len=%d out=%d",dist,n,len(out))};if len(out)+n>maxOutput{return nil,st,fmt.Errorf("REF exceeds output budget")};start:=len(out)-dist;out=append(out,out[start:start+n]...);st.RefOps++
		case OpSeq:
			if pc+2>=len(p){return nil,st,fmt.Errorf("SEQ missing operands")};start,delta,n:=p[pc],p[pc+1],int(p[pc+2]);pc+=3;if n==0||len(out)+n>maxOutput{return nil,st,fmt.Errorf("invalid SEQ count/budget")};v:=start;for j:=0;j<n;j++{out=append(out,v);v+=delta};st.SeqOps++
		default:return nil,st,fmt.Errorf("unknown opcode 0x%02x at pc=%d",op,pc-1)
		}
	}
}

func detectRun(src []byte,i int)int{if i>=len(src){return 0};n:=1;for i+n<len(src)&&src[i+n]==src[i]&&n<255{n++};return n}
func detectSeq(src []byte,i int)(int,byte){if i+2>=len(src){return 0,0};delta:=src[i+1]-src[i];n:=2;v:=src[i+1];for i+n<len(src)&&n<255{v+=delta;if src[i+n]!=v{break};n++};return n,delta}
func detectRef(src []byte,i int)(bestDist,bestLen int){maxDist:=i;if maxDist>255{maxDist=255};for dist:=1;dist<=maxDist;dist++{start:=i-dist;maxLen:=dist;if rem:=len(src)-i;maxLen>rem{maxLen=rem};if maxLen>255{maxLen=255};n:=0;for n<maxLen&&src[start+n]==src[i+n]{n++};if n>bestLen{bestDist,bestLen=dist,n}};return bestDist,bestLen}
