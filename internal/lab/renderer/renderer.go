package renderer
import("bytes";"crypto/sha256";"encoding/hex";"fmt";"image";"image/color";"image/png")
const ResultSchemaV1="ohf.render-result.v1"
type Request struct{Width int;Height int;Source []byte;SourceSHA string;LayoutSeed uint64}
type Result struct{Schema string `json:"schema"`;RendererID string `json:"renderer_id"`;Version string `json:"version"`;Width int `json:"width"`;Height int `json:"height"`;PixelSHA256 string `json:"pixel_sha256"`;PNGSHA256 string `json:"png_sha256"`;PNGBytes int `json:"png_bytes"`;SourceSHA256 string `json:"source_sha256"`;LayoutSeed uint64 `json:"layout_seed"`;ImagePath string `json:"image_path"`}
type Output struct{PNG []byte;Image image.Image;Result Result};type Renderer interface{ID()string;Version()string;Render(Request)(Output,error)}
func EncodePNG(img image.Image)([]byte,error){var buf bytes.Buffer;enc:=png.Encoder{CompressionLevel:png.BestCompression};if err:=enc.Encode(&buf,img);err!=nil{return nil,err};return buf.Bytes(),nil};func PNGHash(b []byte)string{return hashBytes(b)}
func PixelHash(img image.Image)string{b:=img.Bounds();h:=sha256.New();fmt.Fprintf(h,"%d:%d:",b.Dx(),b.Dy());var px[4]byte;for y:=b.Min.Y;y<b.Max.Y;y++{for x:=b.Min.X;x<b.Max.X;x++{c:=color.RGBAModel.Convert(img.At(x,y)).(color.RGBA);px[0],px[1],px[2],px[3]=c.R,c.G,c.B,c.A;_,_=h.Write(px[:])}};return hex.EncodeToString(h.Sum(nil))}
func DecodePNG(b []byte)(image.Image,error){img,err:=png.Decode(bytes.NewReader(b));if err!=nil{return nil,fmt.Errorf("decode png: %w",err)};return img,nil};func hashBytes(b []byte)string{sum:=sha256.Sum256(b);return hex.EncodeToString(sum[:])}
