/*
 * @Author: wangjun haodreams@163.com
 * @Date: 2024-07-21 00:00:46
 * @LastEditors: wangjun haodreams@163.com
 * @LastEditTime: 2025-05-22 16:45:36
 * @FilePath: \gui\example\demo1.go
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
package main

import (
	"time"

	"gioui.org/app"
	"gioui.org/layout"
	"gioui.org/unit"
	"gitee.com/haodreams/gui"
)

type BasicFan struct {
	ID          int     `json:"id"`          // id
	CompanyID   int     `json:"companyId"`   // 公司id
	CompanyName string  `json:"companyName"` // 公司名称
	PlantID     int     `json:"plantId"`     // 场站id
	PlantName   string  `json:"plantName"`   // 场站名称
	StagingID   int     `json:"stagingId"`   // 工期id
	StagingName string  `json:"stagingName"` // 工期名称
	CircuitID   int     `json:"circuitId"`   // 集电线id
	CircuitName string  `json:"circuitName"` // 集电线名称
	FanName     string  `json:"fanName"`     // 风机名称
	PowerField  string  `json:"powerField"`  //电量计算的原始点
	FanCode     string  `json:"fanCode"`     // 风机编码
	InnerCode   string  `json:"innerCode"`   // 内部编码
	ModelID     int     `json:"modelId"`     // 型号id
	ModelName   string  `json:"modelName"`   // 型号名称
	Status      int     `json:"status"`      // 1 运行 2 调试 3 未接入
	StartSpeed  float64 `json:"startSpeed"`  // 切入风速(m/s)
	StopSpeed   float64 `json:"stopSpeed"`   // 切出风速(m/s)
	FanCap      float64 `json:"fanCap"`      // 装机容量
	//Host        string  `json:"host"`

	IsParadigm   int    `json:"isParadigm"`   //是否是标杆
	FanLocalType string `json:"fanLocalType"` //fan_local_type 海风陆风
}

func Init() *gui.Window {
	win := gui.NewWindow()
	win.Option(
		app.Title("Counter"),
		app.Size(unit.Dp(800), unit.Dp(600)),
	)
	var fans []*BasicFan

	for i := 0; i < 100*10000; i++ {
		fans = append(fans, &BasicFan{ID: i + 1, CompanyID: 1, CompanyName: "company1", PlantName: "plant1", StagingName: "staging1", CircuitName: "circuit1", FanName: "#1风机"})
	}

	//array := gui.NewTableObject(fans)

	datatable := gui.NewDataTable(fans, nil, nil)

	table := gui.NewTable(win, datatable)

	header := gui.NewContainer(win)

	header.AddWidget(gui.NewLabel(win, "Hello:"))

	header.AddWidget(gui.NewEdit(win, "test").SetWidth(100))
	header.AddWidget(gui.NewEdit(win, "test2"))
	header.Add(gui.NewSpace(8))

	button := gui.NewButton(win, "Count", func() { win.Log(time.Now()) })
	header.Add(gui.NewSpace(8))

	button2 := gui.NewButton(win, "测试", func() {
		//table.GridState.Vertical.Last
		win.Log(table.Vertical.First, table.VScrollbar.ScrollDistance(), table.Vertical.OffsetAbs, table.Vertical.Length)
		table.Vertical.First = 0
		table.Vertical.Offset = (table.Vertical.Length / (table.Size())) * (100 - 1)
	})

	header.AddWidget(button)
	header.Add(gui.NewSpace(8))

	header.Add(layout.Rigid(button2.Layout))

	root := gui.NewContainer(win)
	root.SetAxis(layout.Vertical)

	root.AddWidget(header)

	memu := gui.NewMenu(win)
	memu.AddItem("add", nil, func() {
		row, col := table.GetSelectedCell()
		win.Log("add", row, col)
	})
	table.SetMenu(memu)

	root.Add(layout.Flexed(1, table.Layout))

	win.SetContent(root)
	return win
}

func main() {
	gui.Run(Init)
}
