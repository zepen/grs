package util

import (
	"encoding/base64"
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"testing"
	"time"
)

func TestGetNowTimeFmt(t *testing.T) {
	timeNow := time.Now()
	timeString := timeNow.Format("2006-01-02 15:04:05") //2015-06-15 08:52:32
	fmt.Println(timeString)
}

func TestRandOneLabel(t *testing.T) {
	userIntentSet := make(map[string]struct{})
	newIntentList := strings.Split("4,3,10", ",")
	for _, v := range newIntentList {
		userIntentSet[v] = struct{}{}
	}
	// 随机生成用户兴趣
	for {
		rand.Seed(time.Now().UnixNano())
		m := rand.Intn(16) + 1
		labelStr := strconv.FormatInt(int64(m), 10)
		if _, ok := userIntentSet[labelStr]; !ok {
			fmt.Println(labelStr)
			break
		}
	}
}

func TestRandGroupLabel(t *testing.T) {
	userIntentSet := make(map[string]struct{})
	// 随机生成用户兴趣
	groupSize := 3
	for {
		rand.Seed(time.Now().UnixNano())
		m := rand.Intn(16) + 1
		labelStr := strconv.FormatInt(int64(m), 10)
		if _, ok := userIntentSet[labelStr]; !ok {
			userIntentSet[labelStr] = struct{}{}
			if len(userIntentSet) == groupSize {
				labelListStr := ""
				for k, _ := range userIntentSet {
					labelListStr += "," + k
				}
				fmt.Println(labelListStr[1:])
				break
			}
		}
	}
}

func TestSplitIntentList(t *testing.T) {
	intentList1 := "12,11,8"
	//intentList2 := "12:1,11:1,8:1"

	intentListSet := make(map[string]struct{})
	intentListSplit := strings.Split(intentList1, ",")

	for _, v := range intentListSplit {
		sv := strings.Split(v, ":")
		if len(sv) > 1 {
			intentListSet[sv[0]] = struct{}{}
		} else {
			intentListSet[v] = struct{}{}
		}
	}
	fmt.Println(intentListSet)
}

func TestRandFloat32CostTime(t *testing.T) {
	m := make(map[string]map[string]string, 10000)
	for i := 0; i < 10000; i++ {
		m[string(rune(i))] = make(map[string]string)
	}
	start := time.Now()
	for _ = range m {
		rand.Seed(time.Now().UnixNano())
		_ = fmt.Sprintf("%.4f", rand.Float32())
	}
	cost := time.Since(start)
	fmt.Println(cost)
}

func TestJoinStr(t *testing.T) {
	var a []string
	fmt.Println(strings.Join(a, ","))
}

func TestRandomString(t *testing.T) {
	length := 10
	randomBytes := make([]byte, length)
	_, err := rand.Read(randomBytes)
	if err != nil {
		panic(err)
	}
	randomString := base64.URLEncoding.EncodeToString(randomBytes)[:length]
	fmt.Println(randomString)
}
