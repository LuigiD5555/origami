package geometry

import (
	"image"
	"image/color"
)

var (
	White = color.RGBA{255,255,255,255}
	Black = color.RGBA{0,0,0,255}
	Gray1 = color.RGBA{224,224,224,255}
	Gray2 = color.RGBA{128,128,128,255}
)

type Canvas struct { Img *image.RGBA }
func New(width,height int,bg color.RGBA)*Canvas{ img:=image.NewRGBA(image.Rect(0,0,width,height)); c:=&Canvas{Img:img}; c.FillRect(0,0,width,height,bg); return c }
func(c *Canvas)set(x,y int,col color.RGBA){ if image.Pt(x,y).In(c.Img.Bounds()){ c.Img.SetRGBA(x,y,col) } }
func(c *Canvas)FillRect(x,y,w,h int,col color.RGBA){ if w<=0||h<=0{return}; r:=image.Rect(x,y,x+w,y+h).Intersect(c.Img.Bounds()); for yy:=r.Min.Y; yy<r.Max.Y; yy++{ for xx:=r.Min.X; xx<r.Max.X; xx++{ c.Img.SetRGBA(xx,yy,col) } } }
func(c *Canvas)StrokeRect(x,y,w,h,t int,col color.RGBA){ if t<=0||w<=0||h<=0{return}; c.FillRect(x,y,w,t,col); c.FillRect(x,y+h-t,w,t,col); c.FillRect(x,y,t,h,col); c.FillRect(x+w-t,y,t,h,col) }
func(c *Canvas)FillCircle(cx,cy,radius int,col color.RGBA){ if radius<=0{return}; rr:=radius*radius; for y:=cy-radius;y<=cy+radius;y++{dy:=y-cy; for x:=cx-radius;x<=cx+radius;x++{dx:=x-cx; if dx*dx+dy*dy<=rr{c.set(x,y,col)}}}}
func(c *Canvas)FillDiamond(cx,cy,radius int,col color.RGBA){ if radius<=0{return}; for dy:=-radius;dy<=radius;dy++{span:=radius-abs(dy); for dx:=-span;dx<=span;dx++{c.set(cx+dx,cy+dy,col)}}}
func(c *Canvas)FillTriangle(cx,cy,radius int,col color.RGBA){ if radius<=0{return}; top:=cy-radius; bottom:=cy+radius; height:=bottom-top; if height==0{c.set(cx,cy,col);return}; for y:=top;y<=bottom;y++{rel:=y-top;span:=(radius*rel)/height;for x:=cx-span;x<=cx+span;x++{c.set(x,y,col)}}}
func(c *Canvas)Line(x0,y0,x1,y1,thickness int,col color.RGBA){ if thickness<=0{return};dx:=abs(x1-x0);sx:=-1;if x0<x1{sx=1};dy:=-abs(y1-y0);sy:=-1;if y0<y1{sy=1};err:=dx+dy;r:=thickness/2;for{c.FillRect(x0-r,y0-r,thickness,thickness,col);if x0==x1&&y0==y1{break};e2:=2*err;if e2>=dy{err+=dy;x0+=sx};if e2<=dx{err+=dx;y0+=sy}}}
func abs(v int)int{if v<0{return -v};return v}
