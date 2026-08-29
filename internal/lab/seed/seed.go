package seed

import("crypto/sha256";"encoding/binary";"fmt")
const Algorithm="splitmix64-v1"
func Derive(master uint64,namespace string)uint64{h:=sha256.New();var b[8]byte;binary.LittleEndian.PutUint64(b[:],master);h.Write([]byte(Algorithm));h.Write([]byte{0});h.Write(b[:]);h.Write([]byte{0});h.Write([]byte(namespace));sum:=h.Sum(nil);return binary.LittleEndian.Uint64(sum[:8])}
type SplitMix64 struct{state uint64};func New(v uint64)*SplitMix64{return &SplitMix64{state:v}};func(s *SplitMix64)Uint64()uint64{s.state+=0x9e3779b97f4a7c15;z:=s.state;z=(z^(z>>30))*0xbf58476d1ce4e5b9;z=(z^(z>>27))*0x94d049bb133111eb;return z^(z>>31)}
func Namespaces(master uint64,names []string)map[string]uint64{out:=make(map[string]uint64,len(names));for _,n:=range names{out[n]=Derive(master,n)};return out};func Explain(master uint64,namespace string)string{return fmt.Sprintf("%d/%s -> %d",master,namespace,Derive(master,namespace))}
