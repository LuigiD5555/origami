package integration

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"sort"

	"github.com/LuigiD5555/origami/internal/fixedcarrier"
	"github.com/LuigiD5555/origami/internal/memory"
)

type CarrierCounts struct {
	Pages     uint32
	Blocks    uint32
	Documents uint32
	Objects   uint32
}

// FixedCarrierMetadata binds the frozen visual control plane to the current
// canonical memory state. The pixels remain a boot/navigation/control plane;
// StoreRoot commits to the external semantic/exact memory universe.
func FixedCarrierMetadata(system *System, sourceHash string, counts CarrierCounts) (fixedcarrier.Metadata, error) {
	nodes:=system.Nodes()
	sort.Slice(nodes,func(i,j int)bool{return nodes[i].Address<nodes[j].Address})
	body,err:=json.Marshal(nodes);if err!=nil{return fixedcarrier.Metadata{},err}
	root:=sha256.Sum256(body)
	signature:=system.Signature()
	sigBody,err:=compactSignature(signature);if err!=nil{return fixedcarrier.Metadata{},err}
	objects:=counts.Objects
	if objects==0{objects=uint32(len(nodes))}
	return fixedcarrier.Metadata{
		CarrierID:system.CarrierID,
		StoreRoot:hex.EncodeToString(root[:]),
		SourceSHA256:sourceHash,
		PageCount:counts.Pages,
		BlockCount:counts.Blocks,
		DocumentCount:counts.Documents,
		ObjectCount:objects,
		GraphSignature:sigBody,
		Flags:0x0001, // canonical external memory binding present
	},nil
}

func CanonicalStoreRoot(nodes []memory.Node) string {
	copyNodes:=append([]memory.Node(nil),nodes...);sort.Slice(copyNodes,func(i,j int)bool{return copyNodes[i].Address<copyNodes[j].Address});body,_:=json.Marshal(copyNodes);sum:=sha256.Sum256(body);return hex.EncodeToString(sum[:])
}

func compactSignature(signature memory.GraphSignature)([]byte,error){
	// Keep the R2 record bounded. The complete GraphSignature remains in the
	// canonical store; this compact copy is only a routing hint bound by StoreRoot.
	payload:=struct{Schema string `json:"s"`;Carrier string `json:"c"`;Roots []string `json:"r,omitempty"`;Classes []string `json:"k,omitempty"`;Clusters []string `json:"g,omitempty"`;Nodes int `json:"n"`}{signature.Schema,signature.CarrierID,limitStrings(signature.RootTerms,8),limitStrings(signature.SemanticClasses,6),limitStrings(signature.Clusters,6),signature.NodeCount}
	body,err:=json.Marshal(payload);if err!=nil{return nil,err};if len(body)>256{body=body[:256]};return body,nil
}
func limitStrings(values []string,n int)[]string{if len(values)<=n{return append([]string(nil),values...)};return append([]string(nil),values[:n]...)}
