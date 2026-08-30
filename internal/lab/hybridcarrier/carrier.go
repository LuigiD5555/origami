package hybridcarrier

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"image"
	"image/draw"
	"image/png"
	"sort"

	"github.com/LuigiD5555/origami/internal/lab/glyphcalc"
	"github.com/LuigiD5555/origami/internal/lab/memoryvisual"
	"github.com/LuigiD5555/origami/internal/lab/renderer"
	"github.com/LuigiD5555/origami/internal/memory"
	"github.com/LuigiD5555/origami/internal/receiver"
)

const (
	SchemaR0           = "origami.hybrid-carrier.r0"
	MagicR0            = "ORIGAMI-HYBRID-R0\n"
	DefaultColumns     = 24
	DefaultCellPixels  = 56
	DefaultMargin      = 8
	DefaultMaxPNGBytes = 512000
)

type IndexEntry struct { Key string `json:"key"`; Address string `json:"address"` }
type MemoryEntry = memory.Node

type MemoryProfile struct {
	Schema               string                 `json:"schema"`
	AddressABI            string                 `json:"address_abi"`
	CarrierID             string                 `json:"carrier_id"`
	DefaultWorkingWindow int                    `json:"default_working_window"`
	Signature             memory.GraphSignature `json:"graph_signature"`
}

type Verification struct { MemorySHA256 string `json:"memory_sha256"` }

type Envelope struct {
	Schema        string         `json:"schema"`
	Receiver      receiver.Spec  `json:"receiver"`
	Index         []IndexEntry   `json:"index"`
	Memory        []MemoryEntry  `json:"memory"`
	MemoryProfile *MemoryProfile `json:"memory_profile,omitempty"`
	Verification  Verification   `json:"verification"`
}

type Rendered struct {
	PNG              []byte               `json:"-"`
	PNGBytes         int                  `json:"png_bytes"`
	PNGSHA256        string               `json:"png_sha256"`
	PayloadBytes     int                  `json:"payload_bytes"`
	PayloadSHA256    string               `json:"payload_sha256"`
	GlyphCount       int                  `json:"glyph_count"`
	MemoryBandHeight int                  `json:"memory_band_height,omitempty"`
	VisualPlan       *memoryvisual.Plan   `json:"visual_plan,omitempty"`
	RenderMeta       glyphcalc.RenderMeta `json:"render_meta"`
}

func BuildEnvelope(spec receiver.Spec, index []IndexEntry, mem []MemoryEntry) (Envelope, error) {
	if err := receiver.Validate(spec); err != nil { return Envelope{}, fmt.Errorf("receiver spec: %w", err) }
	idx:=append([]IndexEntry(nil),index...);nodes:=append([]MemoryEntry(nil),mem...)
	sort.Slice(idx,func(i,j int)bool{if idx[i].Key==idx[j].Key{return idx[i].Address<idx[j].Address};return idx[i].Key<idx[j].Key});sort.Slice(nodes,func(i,j int)bool{return nodes[i].Address<nodes[j].Address})
	if len(nodes)==0{return Envelope{},fmt.Errorf("memory cannot be empty")};seen:=map[string]bool{}
	for _,entry:=range nodes{if err:=memory.ValidateNode(entry);err!=nil{return Envelope{},err};if seen[entry.Address]{return Envelope{},fmt.Errorf("duplicate memory address %q",entry.Address)};seen[entry.Address]=true}
	for _,entry:=range idx{if entry.Key==""||entry.Address==""{return Envelope{},fmt.Errorf("index entries require key and address")};if !seen[entry.Address]{return Envelope{},fmt.Errorf("index key %q points to missing address %q",entry.Key,entry.Address)}}
	memoryBytes,err:=json.Marshal(nodes);if err!=nil{return Envelope{},err};return Envelope{Schema:SchemaR0,Receiver:spec,Index:idx,Memory:nodes,Verification:Verification{MemorySHA256:hash(memoryBytes)}},nil
}

func BuildVirtualMemoryEnvelope(spec receiver.Spec, carrierID string, index []IndexEntry, mem []MemoryEntry) (Envelope, error) {
	envelope,err:=BuildEnvelope(spec,index,mem);if err!=nil{return Envelope{},err};if carrierID==""{carrierID="carrier"};store,err:=memory.BuildStore(carrierID,envelope.Memory);if err!=nil{return Envelope{},err};window:=spec.Boot.WorkingWindow;if window<=0{window=4000};envelope.MemoryProfile=&MemoryProfile{Schema:memory.SchemaR0,AddressABI:memory.AddressABIR0,CarrierID:carrierID,DefaultWorkingWindow:window,Signature:store.Signature()};return envelope,nil
}

func EncodePayload(envelope Envelope) ([]byte,error){if envelope.Schema!=SchemaR0{return nil,fmt.Errorf("carrier schema must be %q",SchemaR0)};if err:=receiver.Validate(envelope.Receiver);err!=nil{return nil,err};body,err:=json.Marshal(envelope);if err!=nil{return nil,err};payload:=make([]byte,0,len(MagicR0)+len(body));payload=append(payload,[]byte(MagicR0)...);payload=append(payload,body...);return payload,nil}
func DecodePayload(payload []byte)(Envelope,error){if !bytes.HasPrefix(payload,[]byte(MagicR0)){return Envelope{},fmt.Errorf("missing Hybrid carrier magic")};var envelope Envelope;if err:=json.Unmarshal(payload[len(MagicR0):],&envelope);err!=nil{return Envelope{},err};if envelope.Schema!=SchemaR0{return Envelope{},fmt.Errorf("unexpected carrier schema %q",envelope.Schema)};if err:=receiver.Validate(envelope.Receiver);err!=nil{return Envelope{},err};memoryBytes,err:=json.Marshal(envelope.Memory);if err!=nil{return Envelope{},err};if got:=hash(memoryBytes);got!=envelope.Verification.MemorySHA256{return Envelope{},fmt.Errorf("memory verification mismatch: expected=%s actual=%s",envelope.Verification.MemorySHA256,got)};return envelope,nil}

func Render(envelope Envelope, columns, cell, margin, maxPNGBytes int)(Rendered,error){if columns<=0{columns=DefaultColumns};if cell<=0{cell=DefaultCellPixels};if margin<=0{margin=DefaultMargin};if cell<DefaultCellPixels{return Rendered{},fmt.Errorf("Hybrid carrier requires cell >= %d for independent deterministic decoding",DefaultCellPixels)};if maxPNGBytes<=0{maxPNGBytes=DefaultMaxPNGBytes};payload,err:=EncodePayload(envelope);if err!=nil{return Rendered{},err};glyphs:=glyphcalc.GlyphsFromBytes(payload);pngBytes,meta,err:=glyphcalc.RenderGlyphGrid(glyphs,columns,cell,margin);if err!=nil{return Rendered{},err};if len(pngBytes)>maxPNGBytes{return Rendered{},fmt.Errorf("carrier exceeds PNG target: bytes=%d max=%d",len(pngBytes),maxPNGBytes)};return Rendered{PNG:pngBytes,PNGBytes:len(pngBytes),PNGSHA256:hash(pngBytes),PayloadBytes:len(payload),PayloadSHA256:hash(payload),GlyphCount:len(glyphs),RenderMeta:meta},nil}

func RenderMemory(envelope Envelope, columns, cell, margin, maxPNGBytes, bandHeight int)(Rendered,error){
	if maxPNGBytes<=0{maxPNGBytes=DefaultMaxPNGBytes};if bandHeight<=0{bandHeight=memoryvisual.DefaultBandHeight};base,err:=Render(envelope,columns,cell,margin,maxPNGBytes);if err!=nil{return Rendered{},err};baseImage,err:=png.Decode(bytes.NewReader(base.PNG));if err!=nil{return Rendered{},err};bounds:=baseImage.Bounds();band,plan:=memoryvisual.RenderBand(envelope.Memory,bounds.Dx(),bandHeight);canvas:=image.NewRGBA(image.Rect(0,0,bounds.Dx(),bandHeight+bounds.Dy()));draw.Draw(canvas,image.Rect(0,0,bounds.Dx(),bandHeight),band,band.Bounds().Min,draw.Src);draw.Draw(canvas,image.Rect(0,bandHeight,bounds.Dx(),bandHeight+bounds.Dy()),baseImage,bounds.Min,draw.Src);pngBytes,err:=renderer.EncodePNG(canvas);if err!=nil{return Rendered{},err};if len(pngBytes)>maxPNGBytes{return Rendered{},fmt.Errorf("memory carrier exceeds PNG target: bytes=%d max=%d",len(pngBytes),maxPNGBytes)};meta:=base.RenderMeta;meta.Width=canvas.Bounds().Dx();meta.Height=canvas.Bounds().Dy();meta.PixelSHA256=renderer.PixelHash(canvas);meta.PNGSHA256=renderer.PNGHash(pngBytes);meta.PNGBytes=len(pngBytes);base.PNG=pngBytes;base.PNGBytes=len(pngBytes);base.PNGSHA256=hash(pngBytes);base.MemoryBandHeight=bandHeight;base.VisualPlan=&plan;base.RenderMeta=meta;return base,nil
}

func DecodePNG(pngBytes []byte, columns, cell, margin int)(Envelope,error){if columns<=0{columns=DefaultColumns};if cell<=0{cell=DefaultCellPixels};if margin<=0{margin=DefaultMargin};glyphs,err:=glyphcalc.DecodeGlyphGridPNG(pngBytes,columns,cell,margin);if err!=nil{return Envelope{},err};payload,err:=glyphcalc.BytesFromGlyphs(glyphs);if err!=nil{return Envelope{},err};return DecodePayload(payload)}
func DecodeMemoryPNG(pngBytes []byte, columns, cell, margin, bandHeight int)(Envelope,error){if bandHeight<=0{return DecodePNG(pngBytes,columns,cell,margin)};img,err:=png.Decode(bytes.NewReader(pngBytes));if err!=nil{return Envelope{},err};b:=img.Bounds();if b.Dy()<=bandHeight{return Envelope{},fmt.Errorf("carrier is smaller than declared memory band")};exact:=image.NewRGBA(image.Rect(0,0,b.Dx(),b.Dy()-bandHeight));draw.Draw(exact,exact.Bounds(),img,image.Pt(b.Min.X,b.Min.Y+bandHeight),draw.Src);exactPNG,err:=renderer.EncodePNG(exact);if err!=nil{return Envelope{},err};return DecodePNG(exactPNG,columns,cell,margin)}
func hash(b []byte)string{s:=sha256.Sum256(b);return hex.EncodeToString(s[:])}
