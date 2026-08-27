package pipeline

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/LuigiD5555/origami/internal/lab/artifact"
	"github.com/LuigiD5555/origami/internal/lab/groundtruth"
	"github.com/LuigiD5555/origami/internal/lab/prepare"
	"github.com/LuigiD5555/origami/internal/lab/renderer"
	"github.com/LuigiD5555/origami/internal/lab/renderer/labgeometry"
	"github.com/LuigiD5555/origami/internal/lab/runid"
	"github.com/LuigiD5555/origami/internal/lab/seed"
	"github.com/LuigiD5555/origami/internal/lab/spec"
	"github.com/LuigiD5555/origami/internal/lab/visualverify"
)

const MaxCarrierBytes=500000
type Event struct{Schema string `json:"schema"`;Sequence int `json:"sequence"`;State string `json:"state"`;Status string `json:"status"`;Detail string `json:"detail,omitempty"`}
type Journal struct{events []Event}
func(j *Journal)Add(state,status,detail string){j.events=append(j.events,Event{Schema:"ohf.pipeline-event.v1",Sequence:len(j.events)+1,State:state,Status:status,Detail:detail})}
func(j *Journal)Write(path string)error{if err:=os.MkdirAll(filepath.Dir(path),0o755);err!=nil{return err};f,err:=os.Create(path+".tmp");if err!=nil{return err};w:=bufio.NewWriter(f);for _,e:=range j.events{b,err:=json.Marshal(e);if err!=nil{_ = f.Close();return err};if _,err:=w.Write(append(b,'\n'));err!=nil{_ = f.Close();return err}};if err:=w.Flush();err!=nil{_ = f.Close();return err};if err:=f.Close();err!=nil{return err};return os.Rename(path+".tmp",path)}
type Result struct{prepare.Result;PixelSHA256 string;PNGSHA256 string;PNGBytes int}
func Run(root string,s spec.ExperimentSpec)(Result,error){var journal Journal;_,specHash,err:=runid.ExperimentID(s);if err!=nil{return Result{},err};journal.Add("LOAD_SPEC","PASS",specHash);if err:=s.Validate();err!=nil{journal.Add("VALIDATE_SPEC","FAIL",err.Error());return Result{},err};journal.Add("VALIDATE_SPEC","PASS",s.Schema);layoutSeed:=seed.Derive(s.Seed.Value,"layout");journal.Add("RESOLVE_SEEDS","PASS",fmt.Sprintf("layout=%d",layoutSeed));truth,material,err:=groundtruth.Build(s);if err!=nil{journal.Add("BUILD_GROUND_TRUTH","FAIL",err.Error());return Result{},err};journal.Add("BUILD_GROUND_TRUTH","PASS",truth.Source.SHA256);foundation,err:=prepare.Run(root,s);if err!=nil{journal.Add("MATERIALIZE_FOUNDATION","FAIL",err.Error());return Result{},err};journal.Add("MATERIALIZE_FOUNDATION","PASS",foundation.RunID);r,err:=resolveRenderer(s);if err!=nil{journal.Add("RESOLVE_RENDERER","FAIL",err.Error());_ = journal.Write(filepath.Join(foundation.RunDir,"events.jsonl"));return Result{},err};journal.Add("RESOLVE_RENDERER","PASS",r.ID()+"@"+r.Version());output,err:=r.Render(renderer.Request{Width:s.Renderer.Width,Height:s.Renderer.Height,Source:material.Bytes,SourceSHA:truth.Source.SHA256,LayoutSeed:layoutSeed});if err!=nil{journal.Add("RENDER","FAIL",err.Error());_ = journal.Write(filepath.Join(foundation.RunDir,"events.jsonl"));return Result{},err};journal.Add("RENDER","PASS",output.Result.PixelSHA256);if err:=artifact.WriteExact(filepath.Join(foundation.RunDir,output.Result.ImagePath),output.PNG);err!=nil{return Result{},err};journal.Add("SAVE","PASS",output.Result.PNGSHA256);if err:=writeJSON(filepath.Join(foundation.RunDir,"render.json"),output.Result);err!=nil{return Result{},err};reopened,err:=renderer.DecodePNG(output.PNG);if err!=nil{return Result{},err};reopenedPixel:=renderer.PixelHash(reopened);journal.Add("REOPEN","PASS",reopenedPixel);if reopenedPixel!=output.Result.PixelSHA256{return Result{},fmt.Errorf("reopened pixel hash mismatch")};if output.Result.PNGBytes>MaxCarrierBytes{return Result{},fmt.Errorf("carrier size %d exceeds hard limit %d",output.Result.PNGBytes,MaxCarrierBytes)};journal.Add("VERIFY_VISUAL","PASS",fmt.Sprintf("png_bytes=%d",output.Result.PNGBytes));if err:=journal.Write(filepath.Join(foundation.RunDir,"events.jsonl"));err!=nil{return Result{},err};verdictPath:=filepath.Join(foundation.RunDir,"verdict.json");vb,err:=os.ReadFile(verdictPath);if err!=nil{return Result{},err};var verdict prepare.FoundationVerdict;if err:=json.Unmarshal(vb,&verdict);err!=nil{return Result{},err};verdict.Schema="ohf.foundation-verdict.v3";verdict.Phase="F1.0-F1.10";verdict.Gates=append(verdict.Gates,prepare.Gate{ID:"G08_RENDERER_CONTRACT_PASS",Status:"PASS",Detail:r.ID()+"@"+r.Version()},prepare.Gate{ID:"G09_LAB_GEOMETRY_RENDERER_PASS",Status:"PASS",Detail:output.Result.PixelSHA256},prepare.Gate{ID:"G10_SAVE_REOPEN_PASS",Status:"PASS",Detail:reopenedPixel},prepare.Gate{ID:"G10B_PNG_FILE_HASH_PASS",Status:"PASS",Detail:output.Result.PNGSHA256},prepare.Gate{ID:"G10C_CARRIER_SIZE_PASS",Status:"PASS",Detail:fmt.Sprintf("%d<=%d",output.Result.PNGBytes,MaxCarrierBytes)},prepare.Gate{ID:"G10D_EVENT_JOURNAL_PASS",Status:"PASS",Detail:fmt.Sprintf("events=%d",len(journal.events))});if verdict.Status!="BLOCKED_BASELINE"{verdict.Status="READY_FOR_EVALUATOR_FOUNDATION"};if err:=writeJSON(verdictPath,verdict);err!=nil{return Result{},err};if _,err:=artifact.Write(foundation.RunDir);err!=nil{return Result{},err};if err:=artifact.Verify(foundation.RunDir);err!=nil{return Result{},err};if _,err:=visualverify.Verify(foundation.RunDir);err!=nil{return Result{},err};return Result{Result:foundation,PixelSHA256:output.Result.PixelSHA256,PNGSHA256:output.Result.PNGSHA256,PNGBytes:output.Result.PNGBytes},nil}
func resolveRenderer(s spec.ExperimentSpec)(renderer.Renderer,error){if s.Renderer.ID==labgeometry.RendererID&&s.Renderer.Version==labgeometry.VersionV1{return labgeometry.Renderer{},nil};return nil,fmt.Errorf("renderer %s@%s unavailable in F1 visual foundation",s.Renderer.ID,s.Renderer.Version)}
func writeJSON(path string,v any)error{b,err:=json.MarshalIndent(v,"","  ");if err!=nil{return err};return artifact.WriteExact(path,append(b,'\n'))}
