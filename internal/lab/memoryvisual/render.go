package memoryvisual

import (
	"image"
	"image/color"
	"image/draw"
	"math"
	"sort"

	"github.com/LuigiD5555/origami/internal/memory"
)

const DefaultBandHeight = 192

// FamilyRole documents which already-registered visual families are used by
// this experimental navigation projection. The projection is not a claim that
// those families are perceptually promoted; Perception Lab remains the gate.
type FamilyRole struct { Family string `json:"family"`; Role string `json:"role"` }

var NavigationFamilyRoles = []FamilyRole{
	{Family:"F01",Role:"node_kind"},
	{Family:"F21",Role:"adjacency"},
	{Family:"F22",Role:"containment"},
	{Family:"F24",Role:"connection"},
	{Family:"F25",Role:"edge_direction"},
	{Family:"F29",Role:"grouping"},
	{Family:"F32",Role:"hierarchy_depth"},
	{Family:"F34",Role:"graph_degree"},
	{Family:"F39",Role:"spatial_zone"},
	{Family:"F40",Role:"macro_region"},
	{Family:"F41",Role:"meso_micro_transition"},
	{Family:"F42",Role:"verification_mark"},
}

type Point struct { Address string `json:"address"`; Scale string `json:"scale"`; Cluster string `json:"cluster,omitempty"`; X int `json:"x"`; Y int `json:"y"`; Radius int `json:"radius"` }
type Plan struct { Width int `json:"width"`; Height int `json:"height"`; Points []Point `json:"points"`; Families []FamilyRole `json:"families"` }

// RenderBand makes a deterministic macro->meso->micro navigation map. It is a
// redundant visual index over the logical memory, never the sole exact source.
func RenderBand(nodes []memory.Node, width, height int) (*image.RGBA, Plan) {
	if width < 300 { width = 300 }
	if height <= 0 { height = DefaultBandHeight }
	img:=image.NewRGBA(image.Rect(0,0,width,height));draw.Draw(img,img.Bounds(),&image.Uniform{C:white},image.Point{},draw.Src)
	drawRect(img,0,0,width-1,height-1,black)
	third:=width/3;drawLine(img,third,0,third,height-1,gray);drawLine(img,2*third,0,2*third,height-1,gray)

	sorted:=append([]memory.Node(nil),nodes...);sort.Slice(sorted,func(i,j int)bool{return sorted[i].Address<sorted[j].Address})
	clusters:=map[string][]memory.Node{};var roots []memory.Node
	for _,n:=range sorted{cluster:=n.ClusterID;if cluster==""{cluster="_"};clusters[cluster]=append(clusters[cluster],n);if n.Parent==""{roots=append(roots,n)}}
	clusterIDs:=make([]string,0,len(clusters));for id:=range clusters{clusterIDs=append(clusterIDs,id)};sort.Strings(clusterIDs)
	positions:=map[string]Point{};plan:=Plan{Width:width,Height:height,Families:append([]FamilyRole(nil),NavigationFamilyRoles...)}

	// MACRO: roots. MESO: cluster anchors. MICRO: addressable nodes.
	for i,n:=range roots{y:=slot(i,len(roots),height);r:=radiusFor(n);p:=Point{Address:n.Address,Scale:"MACRO",Cluster:n.ClusterID,X:third/2,Y:y,Radius:r};positions[n.Address]=p;plan.Points=append(plan.Points,p);drawNode(img,p,n)}
	clusterCenters:=map[string]Point{}
	for i,id:=range clusterIDs{y:=slot(i,len(clusterIDs),height);p:=Point{Address:"cluster:"+id,Scale:"MESO",Cluster:id,X:third+third/2,Y:y,Radius:10};clusterCenters[id]=p;plan.Points=append(plan.Points,p);drawCluster(img,p)}
	microX:=2*third+third/2
	microIndex:=0
	for _,id:=range clusterIDs{for _,n:=range clusters[id]{y:=slot(microIndex,len(sorted),height);microIndex++;r:=radiusFor(n);p:=Point{Address:n.Address,Scale:"MICRO",Cluster:id,X:microX,Y:y,Radius:r};if _,exists:=positions[n.Address];!exists{positions[n.Address]=p};plan.Points=append(plan.Points,p);drawNode(img,p,n);c:=clusterCenters[id];drawLine(img,c.X+c.Radius,c.Y,p.X-r,p.Y,light)}}

	// Explicit parent/edge relations are rendered after nodes so graph structure
	// is visible without changing the canonical logical graph.
	microByAddress:=map[string]Point{};for _,p:=range plan.Points{if p.Scale=="MICRO"{microByAddress[p.Address]=p}}
	for _,n:=range sorted{from,ok:=microByAddress[n.Address];if !ok{continue};targets:=map[string]struct{}{};for _,list:=range n.Links{for _,target:=range list{targets[target]=struct{}{}}};for _,e:=range n.Edges{targets[e.Target]=struct{}{}};for target:=range targets{to,ok:=microByAddress[target];if !ok{continue};drawArrow(img,from.X,from.Y,to.X,to.Y,dark)}}
	return img,plan
}

func radiusFor(n memory.Node) int { r:=5+int(math.Round(memoryImportance(n)*4));if r<5{r=5};if r>10{r=10};return r }
func memoryImportance(n memory.Node) float64 { v:=n.Importance.Global;if n.Importance.Local>v{v=n.Importance.Local};if v<0{return 0};if v>1{return 1};return v }
func slot(i,total,height int) int { if total<=0{return height/2};return 14+(i+1)*(height-28)/(total+1) }

var white=color.RGBA{255,255,255,255};var black=color.RGBA{0,0,0,255};var dark=color.RGBA{55,55,55,255};var gray=color.RGBA{120,120,120,255};var light=color.RGBA{190,190,190,255}

func drawCluster(img *image.RGBA,p Point){drawRect(img,p.X-12,p.Y-8,p.X+12,p.Y+8,dark);drawRect(img,p.X-9,p.Y-5,p.X+9,p.Y+5,light)}
func drawNode(img *image.RGBA,p Point,n memory.Node){switch n.Kind{case "evidence":drawRect(img,p.X-p.Radius,p.Y-p.Radius,p.X+p.Radius,p.Y+p.Radius,black);case "reference":drawDiamond(img,p.X,p.Y,p.Radius,black);case "program","rule":drawTriangle(img,p.X,p.Y,p.Radius,black);default:drawCircle(img,p.X,p.Y,p.Radius,black)};if verified(n){drawCircle(img,p.X,p.Y,p.Radius+3,gray)};degree:=len(n.Edges);for _,targets:=range n.Links{degree+=len(targets)};for i:=0;i<degree&&i<4;i++{set(img,p.X-6+i*4,p.Y+p.Radius+5,dark)}}
func verified(n memory.Node)bool{if n.VerificationRef!=""{return true};for _,f:=range []memory.Fidelity{memory.FidelityEvidence,memory.FidelityExact}{if l,ok:=n.Level(f);ok&&l.Verified{return true}};return false}
func drawArrow(img *image.RGBA,x0,y0,x1,y1 int,c color.RGBA){drawLine(img,x0,y0,x1,y1,c);dx:=x1-x0;dy:=y1-y0;l:=math.Hypot(float64(dx),float64(dy));if l<1{return};ux:=float64(dx)/l;uy:=float64(dy)/l;px:=-uy;py:=ux;ax:=float64(x1)-ux*7;ay:=float64(y1)-uy*7;drawLine(img,x1,y1,int(ax+px*3),int(ay+py*3),c);drawLine(img,x1,y1,int(ax-px*3),int(ay-py*3),c)}
func set(img *image.RGBA,x,y int,c color.RGBA){if image.Pt(x,y).In(img.Bounds()){img.SetRGBA(x,y,c)}}
func drawRect(img *image.RGBA,x0,y0,x1,y1 int,c color.RGBA){for x:=x0;x<=x1;x++{set(img,x,y0,c);set(img,x,y1,c)};for y:=y0;y<=y1;y++{set(img,x0,y,c);set(img,x1,y,c)}}
func drawDiamond(img *image.RGBA,cx,cy,r int,c color.RGBA){drawLine(img,cx,cy-r,cx+r,cy,c);drawLine(img,cx+r,cy,cx,cy+r,c);drawLine(img,cx,cy+r,cx-r,cy,c);drawLine(img,cx-r,cy,cx,cy-r,c)}
func drawTriangle(img *image.RGBA,cx,cy,r int,c color.RGBA){drawLine(img,cx,cy-r,cx+r,cy+r,c);drawLine(img,cx+r,cy+r,cx-r,cy+r,c);drawLine(img,cx-r,cy+r,cx,cy-r,c)}
func drawCircle(img *image.RGBA,cx,cy,r int,c color.RGBA){steps:=int(2*math.Pi*float64(r)*2);if steps<24{steps=24};px,py:=cx+r,cy;for i:=1;i<=steps;i++{a:=2*math.Pi*float64(i)/float64(steps);x:=cx+int(math.Round(float64(r)*math.Cos(a)));y:=cy+int(math.Round(float64(r)*math.Sin(a)));drawLine(img,px,py,x,y,c);px,py=x,y}}
func drawLine(img *image.RGBA,x0,y0,x1,y1 int,c color.RGBA){dx:=abs(x1-x0);sx:=-1;if x0<x1{sx=1};dy:=-abs(y1-y0);sy:=-1;if y0<y1{sy=1};err:=dx+dy;for{set(img,x0,y0,c);if x0==x1&&y0==y1{break};e2:=2*err;if e2>=dy{err+=dy;x0+=sx};if e2<=dx{err+=dx;y0+=sy}}}
func abs(x int)int{if x<0{return -x};return x}
