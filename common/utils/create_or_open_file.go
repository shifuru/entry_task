package utils

import (
	"os"
	"strings"
)

func CreateOrOpenFile(uri string) (*os.File, error) {
	//判断目录
	path := uri[0:strings.LastIndex(uri, "/")]
	_, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			//创建目录
			err = os.Mkdir(path, os.ModePerm)
			if nil != err {
				return nil, err
			}
		}
		return nil, err
	}

	file, err := os.OpenFile(uri, os.O_WRONLY|os.O_APPEND, 0666)
	if nil != err {
		if os.IsExist(err) {
			return nil, err
		}
		//创建文件
		file, err = os.Create(uri)
		if nil != err {
			return nil, err
		}
	}

	return file, nil
}
