package futil

import (
	"bufio"
	"errors"
	"fmt"
	"go.uber.org/zap"
	"gocore/util"
	"io"
	"io/ioutil"
	"os"
	"strings"
)

// 文件是否已存在
func Exist(filename string) bool {
	_, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return true
}

//读取所有行
func Readlines(filepath string) ([]string, error) {
	content, err := ioutil.ReadFile(filepath)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	allline := strings.Split(string(content), "\n")
	return allline, nil
}

func Readlines_buf(filepath string) ([]string, error) {
	// 使用bufio
	fd, err := os.Open(filepath)
	if err != nil {
		util.LogErr(err)
		return nil, err
	}
	defer fd.Close()

	allline := make([]string, 0)
	buff := bufio.NewReader(fd)

	for {
		data, _, eof := buff.ReadLine()

		if eof == io.EOF {
			break
		}

		allline = append(allline, string(data))
	}

	return allline, nil
}

func createDir(path string) {

}

//是否是目录
func IsDir(path string) bool {
	fileInfo, err := os.Stat(path)
	if err != nil {
		util.LogErr(err)
		return false
	}
	return fileInfo.IsDir()
}

func IsFile(path string) bool {
	fileInfo, err := os.Stat(path)
	if err != nil {
		util.LogErr(err)
		return false
	}
	return !fileInfo.IsDir()
}

func Write(data []byte, writer io.Writer) error {
	n, err := writer.Write(data)
	if err != nil {
		util.LogErr(err)
		return err
	}

	if n <= 0 || len(data) != n {
		errStr := fmt.Sprintf("write error, want write %d, but writed length = %d", len(data), n)
		util.LogError(errStr)
		return errors.New(errStr)
	}

	return nil
}

func WriteFile(data []byte, path string) error {
	err := ioutil.WriteFile(path, data, 0666)
	if err != nil {
		util.LogErr(err)
		return err
	}
	return nil
}

func WriteUrl(data []byte, url string) error {
	if isGcpOssUrl(url) {
		return writeGcpUrl(data, url)
	}

	return errors.New("url is not found")
}

func writeGcpUrl(data []byte, url string) error {
	strs := strings.Split(url, ":")
	err := WriteGcpOss(data, strs[0], strs[1], nil)
	if err != nil {
		util.LogErr(err)
		return err
	}
	return nil
}

//>>>>>>>>>>>>>>>>>>  Read  <<<<<<<<<<<<<<<<<<<<<<
func ReadFile(path string) ([]byte, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		util.LogErr(err)
		return nil, err
	}

	if len(data) <= 0 {
		util.LogError("read futil eof", zap.Int("read byte []", len(data)), zap.String("filepath", path))
		return nil, errors.New("read futil eof")
	}

	return data, nil
}

func Read(reader io.Reader) ([]byte, error) {
	data := make([]byte, 0)
	n, err := reader.Read(data)
	if err != nil {
		util.LogErr(err)
		return nil, err
	}

	if n <= 0 {
		errStr := fmt.Sprintf("read EOF or error, read byte length = %d", n)
		util.LogError(errStr)
		return nil, errors.New(errStr)
	}

	return data, nil
}

func ReadUrl(url string) ([]byte, error) {
	if isGcpOssUrl(url) {
		return readGcpOssUrl(url)
	}
	return nil, errors.New("url not found")
}

func readGcpOssUrl(url string) ([]byte, error) {
	strs := strings.Split(url, ":")
	data, err := ReadGcpOss(strs[0], strs[1])
	if err != nil {
		util.LogErr(err)
		return nil, err
	}
	return data, nil
}

//>>>>>>>>>>>>>>>>>>  Read  <<<<<<<<<<<<<<<<<<<<<<
func downloadGcpOss(bucket string, object string, path string) error {
	data, err := ReadGcpOss(bucket, object)
	if err != nil {
		util.LogErr(err)
		return err
	}

	err = WriteFile(data, path)
	if err != nil {
		util.LogErr(err)
		return err
	}
	return nil
}

func downloadGcpOssName(bucket, object string) {

}

func uploadGcpOss(path, bucket, object string) error {
	data, err := ReadFile(path)
	if err != nil {
		util.LogErr(err)
		return err
	}

	err = WriteGcpOss(data, bucket, object, nil)
	if err != nil {
		util.LogErr(err)
		return err
	}
	return nil
}

func uploadGcpOssName(path, bucket string) error {
	data, err := ReadFile(path)
	if err != nil {
		util.LogErr(err)
		return err
	}

	strs := strings.Split(path, "/")
	lastStr := strs[len(strs)-1]
	subfixs := strings.Split(lastStr, ".")
	subfix := subfixs[len(subfixs)-1]

	f := util.Sha256Hash(data)
	dd := util.Base64Encoding(f)
	filename := string(dd)

	filename = filename + "." + subfix
	object := GetDir(subfix, filename) + filename

	err = WriteGcpOss(data, bucket, object, nil)
	if err != nil {
		util.LogErr(err)
		return err
	}
	return nil
}

func isGcpOssUrl(url string) bool {
	return true
}
