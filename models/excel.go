package models

import (
	"github.com/astaxie/beego"
	"io/ioutil"
	"github.com/astaxie/beego/logs"
	"os"
	"bufio"
	"io"
	"fmt"
)

type GetDataParams struct {
	name string
	date string
}

type TableDataObj struct {
	structure []Structure
	tableData []TableData
}

type Structure struct {
	index int32
	key   string
	name  string
}

type TableData struct {
	s string
}

func init() {
	beego.LoadAppConfig("ini", "conf/excel.conf")
}

// 读取指定文件夹下的文件夹并以字符串数组的形式返回
func GetDateByFolder() ([]string, error) {
	var dateFolder []string
	files, err := ioutil.ReadDir(beego.AppConfig.String("file"))
	if err != nil {
		logs.Error("error", err)
		return dateFolder, err
	}
	for _, file := range files {
		if file.IsDir() {
			dateFolder = append(dateFolder, file.Name())
		}
	}
	return dateFolder, nil
}

// @Title 根据表名和日期读取结构文件
func ReadStructureFile(fileName, date string) ([]Structure, error) {
	var structureArr []Structure
	filePath := beego.AppConfig.String("file") + "/" + date + "/" + fileName + "_" + date + ".ddl"
	fileData, err := os.Open(filePath)
	defer fileData.Close()
	if err != nil {
		return structureArr, err
	}
	buf := bufio.NewReader(fileData)
	for {
		line, err := buf.ReadBytes('\n')
		if err != nil {
			if err == io.EOF {
				return structureArr, err
			}
			return structureArr, err
		}
		fmt.Println(line)
	}
	return structureArr, nil
}

func GetTableData(params GetDataParams) (TableDataObj, error) {
	var tableDataObj TableDataObj
	structure, err := ReadStructureFile(params.name, params.date)
	tableDataObj.structure = structure
	if err != nil {
		logs.Error(err)
		tableDataObj.tableData = nil
		tableDataObj.structure = nil
		return tableDataObj, err
	}
	return tableDataObj, nil
}
