package glyphcalc

import "fmt"

type Glyph struct { Index int `json:"index"`; Byte uint8 `json:"byte"`; Shape uint8 `json:"shape"`; Holes uint8 `json:"holes"`; Direction uint8 `json:"direction"`; Frames uint8 `json:"frames"` }
func EncodeByte(index int,b byte)Glyph{return Glyph{Index:index,Byte:b,Shape:b&0x03,Holes:(b>>2)&0x03,Direction:(b>>4)&0x03,Frames:(b>>6)&0x03}}
func DecodeGlyph(g Glyph)(byte,error){if g.Shape>3||g.Holes>3||g.Direction>3||g.Frames>3{return 0,fmt.Errorf("invalid glyph states: shape=%d holes=%d direction=%d frames=%d",g.Shape,g.Holes,g.Direction,g.Frames)};b:=g.Shape+4*g.Holes+16*g.Direction+64*g.Frames;return byte(b),nil}
func GlyphsFromBytes(b []byte)[]Glyph{out:=make([]Glyph,len(b));for i,v:=range b{out[i]=EncodeByte(i,v)};return out}
func BytesFromGlyphs(gs []Glyph)([]byte,error){out:=make([]byte,len(gs));for i,g:=range gs{v,err:=DecodeGlyph(g);if err!=nil{return nil,fmt.Errorf("glyph %d: %w",i,err)};out[i]=v};return out,nil}
