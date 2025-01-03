package mytrade

import (
	"strconv"
	"strings"
	"time"
)

var (
	machineID     int64 // 机器 id 占10位, 十进制范围是 [ 0, 1023 ]
	sn            int64 // 序列号占 12 位,十进制范围是 [ 0, 4095 ]
	lastTimeStamp int64 // 上次的时间戳(毫秒级), 1秒=1000毫秒, 1毫秒=1000微秒,1微秒=1000纳秒
)

func SetMachineId(mid int64) {
	// 把机器 id 左移 12 位,让出 12 位空间给序列号使用
	machineID = mid << 12
}

// 雪花算法
func GetSnowflakeId() int64 {
	lastTimeStamp = time.Now().UnixNano() / 1000000
	SetMachineId(888)
	curTimeStamp := time.Now().UnixNano() / 1000000
	// 同一毫秒
	if curTimeStamp == lastTimeStamp {
		sn++
		// 序列号占 12 位,十进制范围是 [ 0, 4095 ]
		if sn > 4095 {
			time.Sleep(time.Millisecond)
			curTimeStamp = time.Now().UnixNano() / 1000000
			lastTimeStamp = curTimeStamp
			sn = 0
		}
		// 取 64 位的二进制数 0000000000 0000000000 0000000000 0001111111111 1111111111 1111111111  1 ( 这里共 41 个 1 )和时间戳进行并操作
		// 并结果( 右数 )第 42 位必然是 0,  低 41 位也就是时间戳的低 41 位
		rightBinValue := curTimeStamp & 0x1FFFFFFFFFF
		// 机器 id 占用10位空间,序列号占用12位空间,所以左移 22 位; 经过上面的并操作,左移后的第 1 位,必然是 0
		rightBinValue <<= 22
		id := rightBinValue | machineID | sn
		return id
	}
	if curTimeStamp > lastTimeStamp {
		sn = 0
		lastTimeStamp = curTimeStamp
		// 取 64 位的二进制数 0000000000 0000000000 0000000000 0001111111111 1111111111 1111111111  1 ( 这里共 41 个 1 )和时间戳进行并操作
		// 并结果( 右数 )第 42 位必然是 0,  低 41 位也就是时间戳的低 41 位
		rightBinValue := curTimeStamp & 0x1FFFFFFFFFF
		// 机器 id 占用10位空间,序列号占用12位空间,所以左移 22 位; 经过上面的并操作,左移后的第 1 位,必然是 0
		rightBinValue <<= 22
		id := rightBinValue | machineID | sn
		return id
	}
	if curTimeStamp < lastTimeStamp {
		return 0
	}
	return 0
}

// 获取雪花算法id
func GetInstanceId(strategyName string) string {
	return strategyName + strconv.FormatInt(GetSnowflakeId(), 10)
}

// 计算小数点后有效位数的函数
func countDecimalPlaces(str string) int {
	// 去除尾部的零
	str = strings.TrimRight(str, "0")
	// 分割字符串以获取小数部分
	parts := strings.Split(str, ".")
	// 如果有小数部分，则返回其长度
	if len(parts) == 2 {
		return len(parts[1])
	}
	// 如果没有小数部分，则返回0
	return 0
}

// 计算指定精度下最小值  如0 返回1 如1 返回0.1
func getSizeFromPrecision(precision int) string {
	if precision == 0 {
		return "1"
	}
	return "0." + strings.Repeat("0", precision-1) + "1"
}
