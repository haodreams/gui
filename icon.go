/*
 * @Author: wangjun haodreams@163.com
 * @Date: 2024-07-27 08:19:01
 * @LastEditors: wangjun haodreams@163.com
 * @LastEditTime: 2024-08-03 01:12:10
 * @FilePath: \gui\icon.go
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
package gui

import (
	"gioui.org/widget"
	"golang.org/x/exp/shiny/materialdesign/icons"
)

var InfoIcon *widget.Icon = func() *widget.Icon {
	icon, _ := widget.NewIcon(icons.ActionInfo)
	return icon
}()

var WarnIcon *widget.Icon = func() *widget.Icon {
	icon, _ := widget.NewIcon(icons.AlertWarning)
	return icon
}()

var ErrorIcon *widget.Icon = func() *widget.Icon {
	icon, _ := widget.NewIcon(icons.AlertError)
	return icon
}()

var DeleteIcon *widget.Icon = func() *widget.Icon {
	icon, _ := widget.NewIcon(icons.ActionDelete)
	return icon
}()

var ExpandIcon *widget.Icon = func() *widget.Icon {
	icon, _ := widget.NewIcon(icons.NavigationExpandMore)
	return icon
}()

var SaveIcon *widget.Icon = func() *widget.Icon {
	icon, _ := widget.NewIcon(icons.ContentSave)
	return icon
}()

var MenuIcon *widget.Icon = func() *widget.Icon {
	icon, _ := widget.NewIcon(icons.NavigationMenu)
	return icon
}()

var CopyIcon *widget.Icon = func() *widget.Icon {
	icon, _ := widget.NewIcon(icons.ContentContentCopy)
	return icon
}()

var SearchIcon *widget.Icon = func() *widget.Icon {
	icon, _ := widget.NewIcon(icons.ActionSearch)
	return icon
}()

var HomeIcon *widget.Icon = func() *widget.Icon {
	icon, _ := widget.NewIcon(icons.ActionHome)
	return icon
}()

var HeartIcon *widget.Icon = func() *widget.Icon {
	icon, _ := widget.NewIcon(icons.ActionFavorite)
	return icon
}()

var PlusIcon *widget.Icon = func() *widget.Icon {
	icon, _ := widget.NewIcon(icons.ContentAdd)
	return icon
}()

var EditIcon *widget.Icon = func() *widget.Icon {
	icon, _ := widget.NewIcon(icons.ContentCreate)
	return icon
}()

var VisibilityIcon *widget.Icon = func() *widget.Icon {
	icon, _ := widget.NewIcon(icons.ActionVisibility)
	return icon
}()

var CloseIcon *widget.Icon = func() *widget.Icon {
	icon, _ := widget.NewIcon(icons.NavigationClose)
	return icon
}()
var ArrowDropDownIcon *widget.Icon = func() *widget.Icon {
	icon, _ := widget.NewIcon(icons.NavigationArrowDropDown)
	return icon
}()

var NaviLeftIcon *widget.Icon = func() *widget.Icon {
	icon, _ := widget.NewIcon(icons.NavigationChevronLeft)
	return icon
}()

var NaviRightIcon *widget.Icon = func() *widget.Icon {
	icon, _ := widget.NewIcon(icons.NavigationChevronRight)
	return icon
}()

var FileFolderIcon *widget.Icon = func() *widget.Icon {
	icon, _ := widget.NewIcon(icons.FileFolder)
	return icon
}()
var UploadIcon *widget.Icon = func() *widget.Icon {
	icon, _ := widget.NewIcon(icons.FileFileUpload)
	return icon
}()

var DownloadIcon *widget.Icon = func() *widget.Icon {
	icon, _ := widget.NewIcon(icons.FileFileDownload)
	return icon
}()
var RefreshIcon *widget.Icon = func() *widget.Icon {
	icon, _ := widget.NewIcon(icons.NavigationRefresh)
	return icon
}()

var CleanIcon *widget.Icon = func() *widget.Icon {
	icon, _ := widget.NewIcon(icons.EditorFormatColorText)
	return icon
}()
var FileFolderOpen *widget.Icon = func() *widget.Icon {
	icon, _ := widget.NewIcon(icons.FileFolderOpen)
	return icon
}()
