/*
 * @Author: wangjun haodreams@163.com
 * @Date: 2024-07-23 11:14:35
 * @LastEditors: wangjun haodreams@163.com
 * @LastEditTime: 2025-05-24 00:42:50
 * @FilePath: \gui\image.go
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
package gui

import (
	"image/png"
	"os"
	"testing"

	"golang.org/x/exp/shiny/materialdesign/icons"
)

func TestIco2Img(t *testing.T) {
	bg := NRGB("00ff00")
	img, err := IconToImage(icons.ActionInfo, 64,  bg)
	if err != nil {
		t.Fatal(err)
	}

	// 可选：保存图像到文件（需导入 "os" 和 "image/png" 等包）
	file, _ := os.Create("info_icon.png")
	defer file.Close()
	png.Encode(file, img)
}
