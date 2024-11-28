package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"strings"

	"github.com/sergi/go-diff/diffmatchpatch"
)

func main() {
	// 读取文件内容
	oldFilePath := "file_old.txt"
	newFilePath := "file_new.txt"

	oldContent, err := ioutil.ReadFile(oldFilePath)
	if err != nil {
		log.Fatalf("无法读取旧文件: %v", err)
	}

	newContent, err := ioutil.ReadFile(newFilePath)
	if err != nil {
		log.Fatalf("无法读取新文件: %v", err)
	}

	oldText := string(oldContent)
	newText := string(newContent)

	// 创建 diff-match-patch 对象
	dmp := diffmatchpatch.New()

	// 计算差异
	diffs := dmp.DiffMain(oldText, newText, false)

	// 用于存储差异的左右两栏内容
	var oldLines []string
	var newLines []string

	// 临时变量，用于构建 HTML 行
	var oldLineBuilder strings.Builder
	var newLineBuilder strings.Builder

	// 处理差异
	for _, diff := range diffs {
		switch diff.Type {
		case diffmatchpatch.DiffDelete:
			// 将删除的文本标记为红色
			oldLineBuilder.WriteString(fmt.Sprintf("<span style=\"color:red\">%s</span>", diff.Text))
		case diffmatchpatch.DiffInsert:
			// 将插入的文本以绿色显示
			newLineBuilder.WriteString(fmt.Sprintf("<span style=\"color:green\">%s</span>", diff.Text))
		case diffmatchpatch.DiffEqual:
			// 添加相等的文本到当前行
			lines := splitLines(diff.Text)
			for _, line := range lines {
				if oldLineBuilder.Len() > 0 || newLineBuilder.Len() > 0 {
					oldLines = append(oldLines, oldLineBuilder.String())
					newLines = append(newLines, newLineBuilder.String())
					oldLineBuilder.Reset()
					newLineBuilder.Reset()
				}
				oldLines = append(oldLines, line)
				newLines = append(newLines, line)
			}
		}
	}

	// 处理最后的行
	if oldLineBuilder.Len() > 0 || newLineBuilder.Len() > 0 {
		oldLines = append(oldLines, oldLineBuilder.String())
		newLines = append(newLines, newLineBuilder.String())
	}

	// 打印左右两栏
	fmt.Println("Old File\t\tNew File")
	fmt.Println("-------------------------")

	maxLines := max(len(oldLines), len(newLines))
	for i := 0; i < maxLines; i++ {
		oldLine := ""
		newLine := ""

		if i < len(oldLines) {
			oldLine = oldLines[i]
		}
		if i < len(newLines) {
			newLine = newLines[i]
		}

		// 打印格式化的行
		fmt.Printf("%-80s\t%s\n", oldLine, newLine)
	}
}

// 辅助函数：按行分割文本
func splitLines(text string) []string {
	lines := strings.Split(text, "\n")
	// 如果文本末尾有空行，移除这些空行
	if len(lines) > 0 && lines[len(lines)-1] == "" {
		lines = lines[:len(lines)-1]
	}
	return lines
}

// 辅助函数：获取两个整数中的较大者
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
