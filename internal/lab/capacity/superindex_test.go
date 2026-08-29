package capacity

import "testing"

func chain(n int)[]Node{ out:=make([]Node,n); for i:=0;i<n;i++ { a:=addr(i); links:=map[string][]string{}; if i+1<n{links["depends"]=[]string{addr(i+1)}}; out[i]=Node{Address:a,Value:"V",Links:links} }; return out }
func addr(i int)string{return fmtInt(i)}
func fmtInt(i int)string{if i==0{return "N0"};const d="0123456789";b:=[]byte{};for i>0{b=append(b,d[i%10]);i/=10};for l,r:=0,len(b)-1;l<r;l,r=l+1,r-1{b[l],b[r]=b[r],b[l]};return "N"+string(b)}

func TestSuperIndexDoesNotGlobalScan(t *testing.T){
	i:=Build(chain(100000))
	r,err:=i.Follow("N10","depends",2,4000);if err!=nil{t.Fatal(err)}
	if r.Metrics.Touched!=3||r.Metrics.Unfolded!=3||r.Metrics.Exposed!=3{t.Fatalf("not selective: %+v",r.Metrics)}
}

func TestWorkingWindowIsEnforced(t *testing.T){
	nodes:=chain(1000);for i:=range nodes{nodes[i].Value="0123456789012345678901234567890123456789"}
	i:=Build(nodes)
	if _,err:=i.Follow("N0","depends",999,256);err==nil{t.Fatal("expected working-window failure")}
}

func TestUnknownDoesNotFallbackToScan(t *testing.T){
	i:=Build(chain(100000))
	r,err:=i.Follow("MISSING","depends",5,4000);if err!=nil{t.Fatal(err)}
	if len(r.Nodes)!=0||r.Metrics.Touched!=1||r.Metrics.Exposed!=0{t.Fatalf("unknown leaked into scan: %+v",r)}
}
