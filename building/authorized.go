package building

import (
	"bytes"
	"fmt"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

// GetMacheId 通过命令执行获取机器码
func GetMacheId() string {
	cmd := exec.Command("wmic", "csproduct", "get", "UUID")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		fmt.Println("Failed to execute command:", err)
		return "Error"
	}

	lines := strings.Split(out.String(), "\r\n") // 注意：Windows命令输出通常使用\r\n作为行分隔符

	// 检查是否有足够的行
	if len(lines) < 2 {
		fmt.Println("Insufficient output from command")
		return "Error"
	}

	// 跳过第一行（标题），获取第二行（UUID）
	uuid := strings.TrimSpace(lines[1])
	return uuid
}

// CheckMacheID 硬编码进行认证操作，方法原始单简单有效
func CheckMacheID(uuid string) string {
	now := time.Now()
	currentDate := now.Format("20060102")
	currentDateNum, _ := strconv.Atoi(currentDate)
	//  wmic csproduct get UUID    51FA45CC-2C7B-11B2-A85C-E245BDC0124E
	//  lsq  3F08FB00-B72A-11ED-8C91-088FC3DF0CAE
	if uuid == "51FA45CC-2C7B-11B2-A85C-E245BDC0124E" {
		if currentDateNum < 20250101 {
			return "success"
		}
		println("已经过期啦")
		return "late"
	}
	println("机器码不正确！")
	return "fail"
}
