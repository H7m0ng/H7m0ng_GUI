package building

import (
	"bytes"
	"fmt"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

func GetMacheId() string {
	// 调用lscpu命令
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

func CheckMacheID(uuid string) string {
	now := time.Now()
	currentDate := now.Format("20060102")
	currentDateNum, _ := strconv.Atoi(currentDate)
	if uuid == "51FA45CC-2C7B-11B2-A85C-E245BDC0124E" {
		if currentDateNum < 20240803 {
			return "success"
		}
		println("已经过期啦")
		return "late"
	}
	println("机器码不正确！")
	return "fail"
}
