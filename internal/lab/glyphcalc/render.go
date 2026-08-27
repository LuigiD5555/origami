package glyphcalc

import (
	"image"
	"image/color"
	"image/draw"
	"math"

	"github.com/LuigiD5555/origami/internal/lab/renderer"
)

type RenderMeta struct {
	Columns int `json:"columns"`
	Rows int `json:"rows"`
	CellPixels int `json:"cell_pixels"`
	Margin int `json:"margin"`
	GlyphCount int `json:"glyph_count"`
	Width int `json:"width"`
	Height int `json:"height"`
	PixelSHA256 string `json:"pixel_sha256"`
	PNGSHA256 string `json:"png_sha256"`
	PNGBytes int `json:"png_bytes"`
}

func RenderGlyphGrid(gs []Glyph, columns, cell, margin int)([]byte,RenderMeta,error){if columns<=0{columns=16};if cell<36{cell=36};if margin<4{margin=4};rows:=(len(gs)+columns-1)/columns;w:=margin*2+columns*cell;h:=margin*2+rows*cell;img:=image.NewRGBA(image.Rect(0,0,w,h));draw.Draw(img,img.Bounds(),&image.Uniform{C:color.White},image.Point{},draw.Src);for i,g:=range gs{cx:=margin+(i%columns)*cell+cell/2;cy:=margin+(i/columns)*cell+cell/2;drawGlyph(img,cx,cy,cell,g)};pngb,err:=renderer.EncodePNG(img);if err!=nil{return nil,RenderMeta{},err};meta:=RenderMeta{Columns:columns,Rows:rows,CellPixels:cell,Margin:margin,GlyphCount:len(gs),Width:w,Height:h,PixelSHA256:renderer.PixelHash(img),PNGSHA256:renderer.PNGHash(pngb),PNGBytes:len(pngb)};return pngb,meta,nil}

var black=color.RGBA{0,0,0,255}
var white=color.RGBA{255,255,255,255}
func drawGlyph(img *image.RGBA,cx,cy,cell int,g Glyph){baseR:=cell/4;for f:=0;f<int(g.Frames);f++{o:=baseR+5+f*3;drawRectOutline(img,cx-o,cy-o,cx+o,cy+o,black)};r:=baseR;switch g.Shape{case 0:drawCircleOutline(img,cx,cy,r,black);case 1:drawRectOutline(img,cx-r,cy-r,cx+r,cy+r,black);case 2:drawPolyOutline(img,[][2]int{{cx,cy-r},{cx+r,cy},{cx,cy+r},{cx-r,cy}},black);case 3:drawPolyOutline(img,[][2]int{{cx,cy-r},{cx+r,cy+r},{cx-r,cy+r}},black)};for h:=0;h<int(g.Holes);h++{off:=(h*2-int(g.Holes-1))*5;fillCircle(img,cx+off,cy,3,white);drawCircleOutline(img,cx+off,cy,3,black)};bx,by:=cx,cy;d:=r+8;switch g.Direction{case 0:by-=d;case 1:bx+=d;case 2:by+=d;case 3:bx-=d};fillCircle(img,bx,by,3,black)}
func set(img *image.RGBA,x,y int,c color.RGBA){if image.Pt(x,y).In(img.Bounds()){img.SetRGBA(x,y,c)}}
func drawRectOutline(img *image.RGBA,x0,y0,x1,y1 int,c color.RGBA){for x:=x0;x<=x1;x++{set(img,x,y0,c);set(img,x,y1,c)};for y:=y0;y<=y1;y++{set(img,x0,y,c);set(img,x1,y,c)}}
func fillCircle(img *image.RGBA,cx,cy,r int,c color.RGBA){rr:=r*r;for y:=-r;y<=r;y++{for x:=-r;x<=r;x++{if x*x+y*y<=rr{set(img,cx+x,cy+y,c)}}}}
func drawCircleOutline(img *image.RGBA,cx,cy,r int,c color.RGBA){steps:=int(2*math.Pi*float64(r)*2);if steps<24{steps=24};px,py:=cx+r,cy;for i:=1;i<=steps;i++{a:=2*math.Pi*float64(i)/float64(steps);x:=cx+int(math.Round(float64(r)*math.Cos(a)));y:=cy+int(math.Round(float64(r)*math.Sin(a)));drawLine(img,px,py,x,y,c);px,py=x,y}}
func drawPolyOutline(img *image.RGBA,p [][2]int,c color.RGBA){for i:=range p{j:=(i+1)%len(p);drawLine(img,p[i][0],p[i][1],p[j][0],p[j][1],c)}}
func drawLine(img *image.RGBA,x0,y0,x1,y1 int,c color.RGBA){dx:=abs(x1-x0);sx:=-1;if x0<x1{sx=1};dy:=-abs(y1-y0);sy:=-1;if y0<y1{sy=1};err:=dx+dy;for{set(img,x0,y0,c);if x0==x1&&y0==y1{break};e2:=2*err;if e2>=dy{err+=dy;x0+=sx};if e2<=dx{err+=dx;y0+=sy}}}
func abs(x int)int{if x<0{return -x};return x}
