package main

import (
	"fmt"
	"io/ioutil"
	"path"

	"excel/db"

	"github.com/xuri/excelize/v2"
)

var gyomumap = map[string]string{
	"発注": "02",
	"会計": "05",
}

func main() {
	db.CheckTableDef()

	verMap, _ := getSpecialFolderVer(`xxxx`)
	for k, v := range verMap {
		fmt.Println(k, v)
	}

	gamenInfos := readList(`D:\xxxxxx.xlsx`)
	for _, gamenInfo := range gamenInfos {
		if gamenInfo.ApiID != "" {
			path := fmt.Sprintf(`D:/xxxx/%s.%s/xxxx((%s)_(%s)).xlsx`,
				gyomumap[gamenInfo.Kino],
				gamenInfo.Kino,
				gamenInfo.ApiID,
				gamenInfo.GamenName)

			ver, _ := readFdVer(path)
			fmt.Printf("%s\t%s\n", gamenInfo.GamenID, ver)
		} else {
			path := fmt.Sprintf(`D:/xxxx/%s.%s/xxxxxx(%s_%s).xlsx`,
				gyomumap[gamenInfo.Kino],
				gamenInfo.Kino,
				gamenInfo.GamenID,
				gamenInfo.GamenName)

			ver, _ := readFdVer(path)
			fmt.Printf("%s\t%s\n", gamenInfo.GamenID, ver)
		}
	}
}

type GamenInfo struct {
	ApiID     string
	Kino      string
	GamenID   string
	GamenName string
}

// 横展開対応&仕変対応スケジュールから仕様変更一覧シートのリストを取得する
func readList(path string) []GamenInfo {
	f, err := excelize.OpenFile(path)
	if err != nil {
		fmt.Printf("open %s, error %v \n", path, err)
	}
	defer func() {
		if err := f.Close(); err != nil {
			fmt.Printf("close %s, error %v \n", path, err)
		}
	}()

	rows, _ := f.GetRows("仕様変更一覧")
	var gamenInfos []GamenInfo
	for _, row := range rows[3:] {
		if len(row) < 6 {
			continue
		}

		gamenInfo := GamenInfo{
			ApiID:     row[1],
			Kino:      row[3],
			GamenID:   row[4],
			GamenName: row[5],
		}
		gamenInfos = append(gamenInfos, gamenInfo)
	}
	return gamenInfos
}

// バージョン番号を取得する
func readFdVer(path string) (ver string, his string) {
	f, err := excelize.OpenFile(path)
	if err != nil {
		fmt.Printf("open %s, error %v \n", path, err)
	}
	defer func() {
		if err := f.Close(); err != nil {
			fmt.Printf("close %s, error %v \n", path, err)
		}
	}()

	rows, _ := f.GetRows("改版履歴")
	for _, row := range rows[1:] {

		if len(row) < 1 {
			continue
		}
		ver = row[0]
		his = row[5]
	}

	return ver, his
}

func getSpecialFolderVer(basePath string) (map[string]string, error) {
	retMap := make(map[string]string, 10)
	fis, err := ioutil.ReadDir(basePath)
	if err != nil {
		fmt.Printf("read directory failed. pathname=%v, err=%v", basePath, err)
		return retMap, err
	}

	for _, fi := range fis {
		if !fi.IsDir() && path.Ext(fi.Name()) == ".xlsx" {
			fullpath := fmt.Sprintf(`%s\%s`, basePath, fi.Name())
			_, his := readFdVer(fullpath)
			//retMap[fi.Name()] = fmt.Sprintf("%s_%s", ver, his)
			retMap[fi.Name()] = fmt.Sprintf("%s", his)
		}
	}
	return retMap, nil
}
