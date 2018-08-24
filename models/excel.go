package models

import (
	"bufio"
	"dataAnalysis/utils"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"io"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"errors"
)

type GetDataParams struct {
	Name string
	Date string
}

type TableDataObj struct {
	Structure []Structure `json:"structure"`
	TableData []TableData `json:"tableData"`
}

type Structure struct {
	Index int    `json:"index"`
	Key   string `json:"key"`
	Name  string `json:"name"`
}

type TableData struct {
	key map[string]interface{}
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
	filePath := beego.AppConfig.String("file") + date + "/" + fileName + "_" + date + ".ddl"
	logs.Info("数据文件路径：%s", filePath)
	fileData, err := os.Open(filePath)
	defer fileData.Close()
	if err != nil {
		return structureArr, err
	}
	buf := bufio.NewReader(fileData)
	for {
		// 读取每一行数据
		line, err := buf.ReadBytes('\n')
		if err != nil || err == io.EOF {
			break
		}
		// 转换编码之后格式化并追加到结构体数据数组中
		structureArr = append(structureArr, FormatStructureLineData(utils.ConvertToString(string(line), "GBK", "UTF-8")))
	}
	return structureArr, nil
}

// 格式化结构数据
func FormatStructureLineData(str string) Structure {
	structureArr := strings.Split(str, "|")
	index, err := strconv.Atoi(structureArr[0])
	if err != nil {
		index = 0
	}
	structure := Structure{index - 1, strings.ToLower(structureArr[1]), strings.Replace(structureArr[len(structureArr)-1], "\n", "", -1)}
	return structure
}

// 装载结构和数据
func GetTableData(params GetDataParams) (TableDataObj, error) {
	var tableDataObj TableDataObj
	structure, rsErr := ReadStructureFile(params.Name, params.Date)
	tableDataObj.Structure = structure
	tableDataObj.TableData = nil
	// 结构数据读取失败
	if rsErr != nil {
		logs.Error(rsErr)
		tableDataObj.TableData = nil
		tableDataObj.Structure = nil
		return tableDataObj, rsErr
	}
	readTableDataFile(params)
	return tableDataObj, nil
}

// 根据路径判断文件是否存在
func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

// 是否为非拆分文件
func isNonSplitFile(params GetDataParams) (isGz bool, path string, error error) {
	splitFilePath, nonSplitFilePath := getTableDataFileName(params)
	// 是否是拆分文件
	hasSplitFilePath, err := PathExists(splitFilePath)
	if err != nil {
		return false, "", err
	}
	if hasSplitFilePath { // 已经解压的拆分文件存在
		return false, splitFilePath, nil
	} else { // 拆分文件不存在
		hasSplitFileGzPath, err := PathExists(splitFilePath + ".gz") // 是否未解压
		if err != nil {
			return false, "", err
		}
		if hasSplitFileGzPath {
			return true, splitFilePath + ".gz", nil
		}
	}
	// 非拆分文件
	hasNonSplitFilePath, err := PathExists(nonSplitFilePath)
	if err != nil {
		return false, "", err
	}
	if hasNonSplitFilePath { // 已经解压了的非拆分文件
		return false, nonSplitFilePath, nil
	} else { // 非拆分文件不存在
		hasNonSplitFileGzPath, err := PathExists(nonSplitFilePath + ".gz")
		if err != nil {
			return false, "", err
		}
		if hasNonSplitFileGzPath {
			return true, nonSplitFilePath + ".gz", nil
		}
	}
	return false, "", errors.New("unknown error")
}

// 获取数据文件命名
func getTableDataFileName(params GetDataParams) (splitFilePath, nonSplitFilePath string) {
	agencyCode := beego.AppConfig.String("agencyCode")
	baseFilePath := beego.AppConfig.String("file") + params.Date + "/" + params.Name + "_" + params.Date + "_"
	splitFilePath = baseFilePath + agencyCode + ".del"
	nonSplitFilePath = baseFilePath + agencyCode + "0000000.del"
	return splitFilePath, nonSplitFilePath
}

// 读取数据文件
func readTableDataFile(params GetDataParams) {
	isGz, filePath, err := isNonSplitFile(params)
	if err != nil {

	}
	// gzip file
	if isGz {
		gzSuccess, err := utils.Gzip(filePath, "")
		if err != nil {
			logs.Error(err)
			logs.Error("解压失败")
		}
		if gzSuccess {
			logs.Error("解压成功")
		}
		filePath = strings.Replace(filePath,".gz", "", -1)
	}
	logs.Info("是否是压缩文件：%v，文件路径：%s", isGz, filePath)
}
