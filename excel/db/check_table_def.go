package db

import (
	"fmt"
	"io/ioutil"
	"reflect"
	"regexp"
	"strings"

	_ "github.com/mattn/go-oci8"
	"github.com/xuri/excelize/v2"
	"xorm.io/xorm"
)

type TableInfo struct {
	TblID     string
	TblName   string
	Columns   []ColumnInfo
	PkColumns []string
}

type ColumnInfo struct {
	ColId    string
	ColName  string
	ColType  string
	ColLen   string
	Nullable string
	DefValue string
}

func CheckTableDef() {
	tabs := readTableList()

	for _, tab := range tabs {
		//fmt.Println(tab.TblID, tab.TblName, tab.Columns)
		tblOrcl, err := getOracleDbEngine(tab.TblID)
		if err != nil {
			fmt.Printf("table %s is not exists.\n", tab.TblID)
		}
		if reflect.DeepEqual(tab.Columns, tblOrcl) {
			//fmt.Printf("table %s is same.\n", tab.TblID)
		} else {
			//fmt.Printf("table %s xxx%v -------------- %vx.\n", tab.TblID, tab.Columns, tblOrcl)
			//fmt.Printf("table %s xxx\n", tab.TblID)
			compareColInfo(tab.TblID, tab.Columns, tblOrcl)
		}
	}

}

// D:\メモ.xlsx
func readTableList() []TableInfo {
	var path string = `D:\メモ.xlsx`
	f, err := excelize.OpenFile(path)
	if err != nil {
		fmt.Printf("open %s, error %v \n", path, err)
	}
	defer func() {
		if err := f.Close(); err != nil {
			fmt.Printf("close %s, error %v \n", path, err)
		}
	}()

	allTabDefs, _ := readTabDefFileList()

	rows, _ := f.GetRows("ストコンチームメンテテーブル")
	var tblis []TableInfo
	for _, row := range rows[1:] {
		if len(row) < 3 {
			continue
		}

		tbli := TableInfo{
			TblID:   row[1],
			TblName: row[2],
		}

		// 定義ファイルを取得する
		var tabDefFile string
		for k, v := range allTabDefs {
			//if strings.Contains(v, tbli.TblName) {
			if k == tbli.TblName {
				tabDefFile = v
				readColInfo(&tbli, tabDefFile)
				tblis = append(tblis, tbli)
				break
			}
		}

	}
	return tblis
}

func readColInfo(tbli *TableInfo, tabDefFile string) {
	f, err := excelize.OpenFile(tabDefFile)
	if err != nil {
		fmt.Printf("open %s, error %v \n", tabDefFile, err)
	}
	defer func() {
		if err := f.Close(); err != nil {
			fmt.Printf("close %s, error %v \n", tabDefFile, err)
		}
	}()

	rows, _ := f.Rows("ファイルレイアウト(表)")
	b := false
	for rows.Next() {
		row, _ := rows.Columns()

		if row[0] == "1" {
			b = true
		}

		if b && len(row[0]) > 0 {

			// columns
			col := ColumnInfo{
				ColId:    row[7],
				ColName:  row[6],
				ColType:  row[8],
				ColLen:   row[9],
				Nullable: row[10],
				DefValue: row[11],
			}
			// pk
			if len(row[1]) > 0 {
				tbli.PkColumns = append(tbli.PkColumns, row[7])

				//col.IsPk = true
			}

			tbli.Columns = append(tbli.Columns, col)
			//fmt.Printf("xxxx %v    %v \n", row, col.ColID)
		}
	}
}

// 指定目录下全体表定义文件
func readTabDefFileList() (map[string]string, error) {
	retMap := make(map[string]string, 10)
	path := `D:\xxxxx\31.テーブルレイアウト`
	fis, err := ioutil.ReadDir(path)
	if err != nil {
		fmt.Printf("read directory failed. pathname=%v, err=%v", path, err)
		return retMap, err
	}

	// 正規表現
	reg1 := regexp.MustCompile(`テーブルレイアウト\([A-Z]{2}_(.*)\) *.xlsm`)

	for _, fi := range fis {
		params := reg1.FindStringSubmatch(fi.Name())
		fullpath := fmt.Sprintf(`%s\%s`, path, fi.Name())
		if len(params) > 1 {
			retMap[params[1]] = fullpath
		} else {
			continue
		}
	}
	return retMap, nil
}

func getOracleDbEngine(tblId string) ([]ColumnInfo, error) {
	engine, err := xorm.NewEngine("oci8", "xxx/xxx@192.168.11.11:1521/pdb")
	if err != nil {
		fmt.Println(err)
	}

	// tbs, _ := engine.DBMetas()
	// for i, tb := range tbs {
	// 	fmt.Println("index:", i, "tbName", tb.Name, tb.Comment)

	// 	for _, col := range tb.Columns() {
	// 		// .Name, col.Comment, col.Default, col.FieldName, col.Nullable, col.SQLType.Name, col.Length
	// 		fmt.Printf("%v\n", col)
	// 	}
	// }
	var colInfo []ColumnInfo
	err = engine.SQL(sql, tblId).Find(&colInfo)
	//fmt.Printf("%v", colInfo)
	return colInfo, err
}

var sql string = `
SELECT
     T1.COLUMN_NAME COL_ID
    , T2.COMMENTS COL_NAME
    , T1.DATA_TYPE COL_TYPE
    , CASE 
        WHEN T1.DATA_TYPE = 'NUMBER' 
            THEN T1.DATA_PRECISION || ',' || T1.DATA_SCALE
        WHEN T1.DATA_TYPE = 'DATE' 
            THEN '' 
        ELSE TO_CHAR(T1.DATA_LENGTH)
        END AS COL_LEN
    , CASE WHEN T1.NULLABLE = 'N' THEN 'NOT NULL' ELSE '' END Nullable
    , T1.DATA_DEFAULT DEF_VALUE
FROM
    USER_TAB_COLS T1 
    LEFT JOIN USER_COL_COMMENTS T2 
        ON T1.TABLE_NAME = T2.TABLE_NAME 
        AND T1.COLUMN_NAME = T2.COLUMN_NAME 
WHERE
    T1.TABLE_NAME = ? 
ORDER BY
    T1.COLUMN_ID`

func compareColInfo(tblId string, a []ColumnInfo, b []ColumnInfo) {
	fmt.Printf("-------------------テーブル名　%s------------------\n", tblId)
	for i, infoA := range a {

		if i >= len(b) {
			fmt.Printf("カラム個数不一致\n")
			return
		}

		if infoA.ColId != b[i].ColId {
			fmt.Printf("カラム(%-40s)のID不一致、テーブル定義：%s -- オラクル：%s \n", infoA.ColId, infoA.ColId, b[i].ColId)
		}
		if infoA.ColName != b[i].ColName {
			fmt.Printf("カラム(%-40s)の論理名不一致、テーブル定義：%s--オラクル：%s \n", infoA.ColId, infoA.ColName, b[i].ColName)
		}
		if infoA.ColType != b[i].ColType {
			fmt.Printf("カラム(%-40s)の属性不一致、テーブル定義：%s--オラクル：%s \n", infoA.ColId, infoA.ColType, b[i].ColType)
		}
		if infoA.ColLen != b[i].ColLen {
			fmt.Printf("カラム(%-40s)の桁数不一致、テーブル定義：%s--オラクル：%s \n", infoA.ColId, infoA.ColLen, b[i].ColLen)
		}

		if infoA.Nullable != b[i].Nullable {
			fmt.Printf("カラム(%-40s)のの列制約不一致、テーブル定義：%s--オラクル：%s \n", infoA.ColId, infoA.Nullable, b[i].Nullable)
		}
		if infoA.DefValue != strings.TrimSpace(b[i].DefValue) {
			fmt.Printf("カラム(%-40s)のデフォルト値不一致、テーブル定義：%s--オラクル：%s \n", infoA.ColId, infoA.DefValue, b[i].DefValue)
		}
	}
	fmt.Println("")
}
