package env

import (
	"bytes"
	"io"
	"os"
	"strings"
)

// 加载配置文件
func Load(filenames ...string) (err error) {
	filenames = filenamesOrDefault(filenames)

	for _, filename := range filenames {
		err = loadFile(filename, false)
		if err != nil {
			return // return early on a spazout
		}
	}
	return
}

// 处理默认值的问题
func filenamesOrDefault(filenames []string) []string {
	if len(filenames) == 0 {
		return []string{".env"}
	}
	return filenames
}

// 整体基本流程
func loadFile(filename string, overload bool) error {

	// 加载文件
	envMap, err := readFile(filename)
	if err != nil {
		return err
	}

	// 将值设置到环境变量中
	currentEnv := map[string]bool{}

	// 获取当前环境变量值
	rawEnv := os.Environ()
	for _, rawEnvLine := range rawEnv {
		key := strings.Split(rawEnvLine, "=")[0]
		currentEnv[key] = true
	}

	// 将解析出来的值，设置到环境变量中
	for key, value := range envMap {
		if !currentEnv[key] || overload {
			_ = os.Setenv(key, value)
		}
	}

	return nil
}

func readFile(filename string) (envMap map[string]string, err error) {
	// 打开文件
	file, err := os.Open(filename)
	if err != nil {
		return
	}
	defer file.Close()

	return Parse(file)
}

// Parse reads an env file from io.Reader, returning a map of keys and values.
// 解析文件，返回map string => string
func Parse(r io.Reader) (map[string]string, error) {
	var buf bytes.Buffer

	// 给转换为二进制了
	_, err := io.Copy(&buf, r)
	if err != nil {
		return nil, err
	}

	return UnmarshalBytes(buf.Bytes())
}

// UnmarshalBytes parses env file from byte slice of chars, returning a map of keys and values.
func UnmarshalBytes(src []byte) (map[string]string, error) {
	out := make(map[string]string)
	err := parseBytes(src, out)

	return out, err
}
