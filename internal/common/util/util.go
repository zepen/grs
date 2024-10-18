package util

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"
	"unsafe"
)

func GetNowTimeFmt() string {
	timeNow := time.Now()
	timeString := timeNow.Format("2006-01-02 15:04:05") //2015-06-15 08:52:32
	return timeString
}

func GetNowHourAndMin() string {
	timeNow := time.Now()
	hour := timeNow.Hour()
	hourStr := ""
	if hour < 10 {
		hourStr = fmt.Sprintf("0%v", hour)
	} else {
		hourStr = fmt.Sprintf("%v", hour)
	}
	min := timeNow.Minute()
	minStr := ""
	if min < 10 {
		minStr = fmt.Sprintf("0%v", min)
	} else {
		minStr = fmt.Sprintf("%v", min)
	}
	return hourStr + minStr
}

func GenToday() string {
	timeNow := time.Now()
	timeString := timeNow.Format("20060102")
	return timeString
}

func GenYestday() string {
	timeNow := time.Now()
	yesTime := timeNow.AddDate(0, 0, -1)
	timeString := yesTime.Format("20060102")
	return timeString
}

func TimeStrToMills(stringTime string) (int64, error) {
	if stringTime == "" {
		return 0, nil
	}
	loc, _ := time.LoadLocation("Local")
	theTime, err := time.ParseInLocation("2006-01-02 15:04:05", stringTime, loc)
	if err != nil {
		return 0, err
	}
	return theTime.Unix(), nil
}

func String2Byte(s string) []byte {
	x := (*[2]uintptr)(unsafe.Pointer(&s))
	h := [3]uintptr{x[0], x[1], x[1]}
	return *(*[]byte)(unsafe.Pointer(&h))
}

func ToString(i interface{}) string {
	bb, _ := json.Marshal(i)
	return string(bb)
}

func JSON2Map(str string) map[string]interface{} {
	var tempMap = make(map[string]interface{})
	err := json.Unmarshal([]byte(str), &tempMap)
	if err != nil {
		panic(err)
	}
	return tempMap
}

func RandOneLabel(intentList string, n int) string {
	userIntentSet := make(map[string]struct{})
	newIntentList := strings.Split(intentList, ",")
	for _, v := range newIntentList {
		userIntentSet[v] = struct{}{}
	}
	// 随机生成用户兴趣
	for {
		rand.Seed(time.Now().UnixNano())
		m := rand.Intn(n) + 1
		labelStr := strconv.FormatInt(int64(m), 10)
		if _, ok := userIntentSet[labelStr]; !ok {
			return labelStr
		}
	}
}

func RandGroupLabel(n int, groupSize int) string {
	userIntentSet := make(map[string]struct{})
	// 随机生成用户兴趣
	for {
		rand.Seed(time.Now().UnixNano())
		m := rand.Intn(n) + 1
		labelStr := strconv.FormatInt(int64(m), 10)
		if _, ok := userIntentSet[labelStr]; !ok {
			userIntentSet[labelStr] = struct{}{}
			if len(userIntentSet) == groupSize {
				labelListStr := ""
				for k, _ := range userIntentSet {
					labelListStr += "," + k
				}
				return labelListStr[1:]
			}
		}
	}
}

func RandIF() bool {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(2) == 0
}

func RandFloat32() float32 {
	rand.Seed(time.Now().UnixNano())
	return rand.Float32()
}

func RandFloat32String(a float32) string {
	return fmt.Sprintf("%.4f", RandFloat32()*a)
}

func RandomString(length int) string {
	randomBytes := make([]byte, length)
	_, err := rand.Read(randomBytes)
	if err != nil {
		panic(err)
	}
	return base64.URLEncoding.EncodeToString(randomBytes)[:length]
}
