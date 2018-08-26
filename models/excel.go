package models

import (
	"bufio"
	"dataAnalysis/utils"
	"errors"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"io"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"encoding/json"
	"github.com/360EntSecGroup-Skylar/excelize"
)

type GetDataParams struct {
	Name string `json:"name"`
	Date string `json:"date"`
}

type GetExcelParams struct {
	Name string      `json:"name"`
	Date string      `json:"date"`
	Key  []Structure `json:"key"`
}

type TableDataObj struct {
	Structure []Structure              `json:"structure"`
	TableData []map[string]interface{} `json:"tableData"`
}

type Structure struct {
	Index int    `json:"index"`
	Key   string `json:"key"`
	Name  string `json:"name"`
}

func init() {
	beego.LoadAppConfig("ini", "conf/excel.conf")
}

// 获取缓存路径
func getTableDataCachePath(params GetDataParams) string {
	filePath := "static/config/" + params.Name + "_" + params.Date + ".json"
	return filePath
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

// 格式化表格数据
func FormatTableLineData(str string, structure []Structure) map[string]interface{} {
	var tableData = make(map[string]interface{})
	tableArr := strings.Split(str, "|$|")
	tableArr[len(tableArr)-1] = strings.Replace(tableArr[len(tableArr)-1], "\n", "", -1)
	for index, item := range structure {
		tableData[string(item.Key)] = tableArr[index]
	}
	return tableData
}

// 装载结构和数据
func GetTableData(params GetDataParams) (TableDataObj, error) {
	var tableDataObj TableDataObj
	cachePath := getTableDataCachePath(params)
	// 是否有缓存文件
	hasCache, err := PathExists(cachePath)
	if err == nil && hasCache {
		jsonData, err := ioutil.ReadFile(cachePath)
		if err == nil {
			err = json.Unmarshal(jsonData, &tableDataObj)
			if err == nil {
				logs.Info("读取缓存成功")
				return tableDataObj, nil
			}
		}
	}
	structure, err := ReadStructureFile(params.Name, params.Date)
	tableDataObj.Structure = structure
	// 结构数据读取失败
	if err != nil {
		logs.Error(err)
		tableDataObj.TableData = nil
		tableDataObj.Structure = nil
		return tableDataObj, err
	}
	tableData, err := readTableDataFile(params, structure)
	if err != nil {
		tableDataObj.TableData = nil
		tableDataObj.Structure = structure
		return tableDataObj, err
	}
	tableDataObj.TableData = tableData
	saveTableDataToJson(params, tableDataObj)
	return tableDataObj, nil
}

// 保存数据到json
func saveTableDataToJson(params GetDataParams, obj TableDataObj) {
	filePath := getTableDataCachePath(params)
	// map 转 json
	objData, err := json.Marshal(obj)
	if err != nil {
		logs.Info("保存json失败")
		logs.Error(err)
	}
	// 存储 json
	jsonFile, err := os.Create(filePath)
	if err != nil {
		logs.Error(err)
	}
	defer jsonFile.Close()
	jsonFile.Write(objData)
	logs.Info("保存json成功：" + filePath)
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
func readTableDataFile(params GetDataParams, structure []Structure) ([]map[string]interface{}, error) {
	var tableData []map[string]interface{}
	isGz, filePath, err := isNonSplitFile(params)
	if err != nil {
		return tableData, err
	}
	// gzip file
	if isGz {
		gzSuccess, err := utils.Gzip(filePath, "")
		if err != nil {
			logs.Error(err)
			logs.Error("解压失败")
			return tableData, err
		}
		if gzSuccess {
			logs.Error("解压成功")
		}
		filePath = strings.Replace(filePath, ".gz", "", -1)
	}
	logs.Info("是否是压缩文件：%v，文件路径：%s", isGz, filePath)
	// 读取文件
	fileData, err := os.Open(filePath)
	defer fileData.Close()
	if err != nil {
		return tableData, err
	}
	buf := bufio.NewReader(fileData)
	lineNum := 0
	for {
		// 读取每一行数据
		line, err := buf.ReadBytes('\n')
		// 只提取100条数据
		if lineNum > 100 {
			break
		}
		if err != nil || err == io.EOF {
			break
		}
		// 转换编码之后格式化并追加到结构体数据数组中
		tableData = append(tableData, FormatTableLineData(utils.ConvertToString(string(line), "GBK", "UTF-8"), structure))
		lineNum++
	}
	return tableData, nil
}

// 创建excel
func CareatTableExcel(params GetExcelParams) (string, error) {
	var jsonObj TableDataObj
	var tableObj map[string]utils.TableNameObj
	jsonPath := getTableDataCachePath(GetDataParams{params.Name, params.Date})
	jsonData, err := ioutil.ReadFile(jsonPath)
	tableData, tErr := ioutil.ReadFile("static/config/tableName.json")
	if err != nil || tErr != nil {
		return "", err
	}
	json.Unmarshal(jsonData, &jsonObj)
	json.Unmarshal(tableData, &tableObj)
	tableNameArr := strings.Split(params.Name, "_")
	tableNameStr := strings.Join(tableNameArr[:len(tableNameArr)-1], "_")
	// 表名
	tableNameText := tableObj[tableNameStr].Text
	excelPath := "static/data/" + tableNameText + "_" + params.Date + ".xlsx"
	xlsx := excelize.NewFile()
	sheetName := "Sheet1"
	// 创建一个工作表
	index := xlsx.NewSheet(sheetName)
	// 设置单元格的值
	utils.ModifyExcelCellByAxis(xlsx, sheetName, utils.ChangIndexToAxis(0, 0), "附件")
	// 单元格名称
	utils.ModifyExcelCellByAxis(xlsx, sheetName, utils.ChangIndexToAxis(1, 0), tableNameText)
	for yIndex, key := range params.Key {
		// 表头
		utils.ModifyExcelCellByAxis(xlsx, sheetName, utils.ChangIndexToAxis(2, yIndex), key.Name)
		for xIndex, item := range jsonObj.TableData {
			utils.ModifyExcelCellByAxis(xlsx, sheetName, utils.ChangIndexToAxis(xIndex+3, yIndex), item[key.Key])
		}
	}
	/*        单元格样式      */
	// 合并表头
	xlsx.MergeCell(sheetName, utils.ChangIndexToAxis(1, 0), utils.ChangIndexToAxis(1, len(params.Key)))
	style, err := xlsx.NewStyle(`{"type"":1,"horizontal":"center","vertical":"center"}`)
	xlsx.SetCellStyle(sheetName, utils.ChangIndexToAxis(1, 0), utils.ChangIndexToAxis(1, len(params.Key)), style)
	// 设置工作簿的默认工作表
	xlsx.SetActiveSheet(index)
	// 根据指定路径保存文件
	err = xlsx.SaveAs(excelPath)
	if err != nil {
		return "", err
	}
	return "/" + excelPath, nil
}
