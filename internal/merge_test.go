package internal

import (
	"fmt"
	"testing"
	"time"
)

func TestMergeScore(t *testing.T) {

	// info:id 838296345096884224, pubTime 1688977034000, timeScore 0.604627, mergeScore 0.290584, rw 0.540000
	// info:id 836493340328464386, pubTime 1688547377000, timeScore 0.547089, mergeScore 0.407363, rw 0.510000
	// info:id 836493318882988032, pubTime 1688547174000, timeScore 0.547089, mergeScore 0.253625, rw 0.510000
	pubTime := int64(1687836021000)
	weight := float32(0.510000)
	d := int((time.Now().UnixMilli() - pubTime) / 86400000)
	timeScore := expDecay(0.8, 0.02, float32(d))
	mergeScore := (0.44*timeScore + 0.5*timeScore) * weight
	fmt.Printf("pubTime %d, timeScore %f, mergeScore %f, rw %f\n", pubTime, timeScore, mergeScore, weight)
}
