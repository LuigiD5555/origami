package observation

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/LuigiD5555/origami/internal/lab/artifact"
	"github.com/LuigiD5555/origami/internal/lab/nativebundle"
)

const SchemaV1 = "ohf.observation.v1"
type Status string
const(Known Status="KNOWN";Ambiguous Status="AMBIGUOUS";Unknown Status="UNKNOWN";Invalid Status="INVALID")
type Envelope struct{Schema string `json:"schema"`;BundleID string `json:"bundle_id"`;Model string `json:"model"`;Trial int `json:"trial"`;Status Status `json:"status"`;Observations map[string]any `json:"observations"`;Notes string `json:"notes,omitempty"`}
func(e Envelope)Validate()error{if e.Schema!=SchemaV1{return fmt.Errorf("observation schema: expected %q",SchemaV1)};if e.BundleID==""||e.Model==""{return fmt.Errorf("bundle_id and model required")};if e.Trial<=0{return fmt.Errorf("trial must be >0")};switch e.Status{case Known,Ambiguous,Unknown,Invalid:default:return fmt.Errorf("invalid epistemic status %q",e.Status)};return nil}
func LoadStrict(path string)(Envelope,[]byte,error){b,err:=os.ReadFile(path);if err!=nil{return Envelope{},nil,err};dec:=json.NewDecoder(bytes.NewReader(b));dec.DisallowUnknownFields();dec.UseNumber();var e Envelope;if err:=dec.Decode(&e);err!=nil{return Envelope{},nil,fmt.Errorf("decode observation: %w",err)};var x any;if err:=dec.Decode(&x);err!=io.EOF{return Envelope{},nil,fmt.Errorf("observation trailing JSON")};if err:=e.Validate();err!=nil{return Envelope{},nil,err};canon,err:=json.Marshal(e);if err!=nil{return Envelope{},nil,err};return e,canon,nil}
func Import(root,runID,path string)(string,string,error){e,canon,err:=LoadStrict(path);if err!=nil{return "","",err};if _,err:=nativebundle.Verify(root,runID,e.BundleID);err!=nil{return "","",fmt.Errorf("observation bundle isolation/identity: %w",err)};id:=hash(canon)[:16];model:=sanitize(e.Model);out:=filepath.Join(root,"runs",runID,"native","observations",model,id+".json");if err:=artifact.WriteExact(out,append(canon,'\n'));err!=nil{return "","",err};if _,err:=artifact.Write(filepath.Join(root,"runs",runID));err!=nil{return "","",err};return id,out,nil}
func hash(b []byte)string{s:=sha256.Sum256(b);return hex.EncodeToString(s[:])}
func sanitize(s string)string{s=strings.TrimSpace(s);var b strings.Builder;for _,r:=range s{if(r>='a'&&r<='z')||(r>='A'&&r<='Z')||(r>='0'&&r<='9')||r=='-'||r=='_'{b.WriteRune(r)}else{b.WriteByte('_')}};if b.Len()==0{return "model"};return b.String()}
func Find(root,runID,id string)(Envelope,string,error){base:=filepath.Join(root,"runs",runID,"native","observations");var found string;err:=filepath.WalkDir(base,func(path string,d os.DirEntry,err error)error{if err!=nil{return err};if d.IsDir(){return nil};if strings.TrimSuffix(filepath.Base(path),filepath.Ext(path))==id{found=path};return nil});if err!=nil{return Envelope{},"",err};if found==""{return Envelope{},"",fmt.Errorf("observation %s not found",id)};e,_,err:=LoadStrict(found);return e,found,err}
