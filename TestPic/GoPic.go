package TestPic


import (
"fmt"
"image"
"image/color"
"image/draw"
"image/png"
"net/http"
)


var (
blue  color.Color = color.RGBA{0, 0, 255, 255}
picwidth int = 640
picheight int = 480
)


// 大家可以查看这个网址看看这个image包的使用方法 http://golang.org/doc/articles/image_draw.html
func main() {
http.HandleFunc("/", TTT)
http.ListenAndServe(":999", nil)
}


func TTT(rw http.ResponseWriter, req *http.Request) {

fmt.Printf("somebody come in")
//创建一个图像

m := image.NewRGBA(image.Rect(0, 0, picwidth , picheight)) //*NRGBA (image.Image interface)
// 填充蓝色,并把其写入到m
draw.Draw(m, m.Bounds(), &image.Uniform{blue}, image.ZP, draw.Src)
//以png编码格式,并将m写入到 rw里面去
png.Encode(rw, m) //Encode writes the Image m to w in PNG format.



}
