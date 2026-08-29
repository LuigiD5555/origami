package glyphcalc

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type NativeResponse struct { Status string `json:"status"`; BytesHex string `json:"bytes_hex"`; GlyphCount int `json:"glyph_count"` }
type NativeCheck struct { Status string `json:"status"`; PerceptionExact bool `json:"perception_exact"`; ExactReconstruction bool `json:"exact_reconstruction"`; ObservedGlyphCount int `json:"observed_glyph_count"`; ExpectedGlyphCount int `json:"expected_glyph_count"`; ObservedTransportBytes int `json:"observed_transport_bytes"` }
func CheckNative(runDir,responsePath string)(NativeCheck,error){b,err:=os.ReadFile(responsePath);if err!=nil{return NativeCheck{},err};var r NativeResponse;if err:=json.Unmarshal(b,&r);err!=nil{return NativeCheck{},err};if r.Status!="KNOWN"{return NativeCheck{Status:r.Status},nil};obs,err:=hex.DecodeString(r.BytesHex);if err!=nil{return NativeCheck{},fmt.Errorf("bytes_hex: %w",err)};exp,err:=os.ReadFile(filepath.Join(runDir,"transport.bin"));if err!=nil{return NativeCheck{},err};src,err:=os.ReadFile(filepath.Join(runDir,"source.bin"));if err!=nil{return NativeCheck{},err};perception:=string(obs)==string(exp)&&r.GlyphCount==len(exp);check:=NativeCheck{Status:r.Status,PerceptionExact:perception,ObservedGlyphCount:r.GlyphCount,ExpectedGlyphCount:len(exp),ObservedTransportBytes:len(obs)};if perception{recon,_,err:=ExecuteTransport(obs,max(len(src)*2,1024),max(len(obs)*2,128));if err!=nil{return check,err};check.ExactReconstruction=string(recon)==string(src)};return check,nil}
