package main

import (
	"math/rand"
)

// Style 样式
type Style struct {
	backgroundColor uint32
	bgTransparent   bool
	borderColor     uint32
	borderWeight    int
}

// StyleSheetManager 样式管理器
type StyleSheetManager struct {
	styleSheet map[string]*Style
}

// NewStyleSheetManager 构造函数
func NewStyleSheetManager() (styleSheet *StyleSheetManager) {
	styleSheet = &StyleSheetManager{make(map[string]*Style)}
	hoverborder := &Style{0x00000000, true, 0xffff0000, 2}
	styleSheet.AddStyle("hoverborder", hoverborder)
	return styleSheet
}

// AddStyle 添加样式
func (t *StyleSheetManager) AddStyle(key string, style *Style) {
	t.styleSheet[key] = style
}

// RemoveStyle 删除样式
func (t *StyleSheetManager) RemoveStyle(key string) {
	delete(t.styleSheet, key)
}

// GetStyle 获取样式
func (t *StyleSheetManager) GetStyle(key string) *Style {
	_, ok := t.styleSheet[key]
	if key == "" || !ok {
		return t.GetRandStyle()
	}
	return t.styleSheet[key]
}

// GetRandStyle 随机样式 用于测试
func (t *StyleSheetManager) GetRandStyle() *Style {
	return &Style{uint32(rand.Intn(0xffffffff)), false, 0xffffffff, 0}
}
