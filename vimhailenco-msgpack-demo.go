package main

import (
	"fmt"
	// 另外一个msgpack的包
	"development/msgpack-demo/utils"
	"github.com/vmihailenco/msgpack"
	"strconv"
	"time"
)

//这个demo是可以正确解码的
//这个msgpack支持tag

func main() {

	redis_init()
	defer redis_destory()
	redis_unpck()
}

func redis_init() {
	redisHost := "192.168.10.150"
	redisPort := strconv.Itoa(10002)
	redisServer := redisHost + ":" + redisPort
	//初始化redis
	utils.RedisInit(redisServer, "")
}

func redis_destory() {
	utils.RedisDestroy()
}

// 视频创意组合
type VideoCreativeComb struct {
	CreativeID    string `msgpack:"creative_id"`
	SizeID        string `msgpack:"size_id"`
	AdvertiserID  string `msgpack:"advertiser_id"`
	CombinationID string `msgpack:"combination_id"`
	FileType      string `msgpack:"file_type"`
	Height        int32  `msgpack:"height"`
	Width         int32  `msgpack:"width"`
	VideoBitRate  int32  `msgpack:"video_bit_rate"`
	AudioBitRate  int32  `msgpack:"audio_bit_rate"`
	// Sar                float32 `msgpack:"sar_pixel"`
	VideoCombinationID string `msgpack:"dynamic_combination_id"`
	Version            string `msgpack:"version"`
	AudioEnc           string `msgpack:"audio_encoder"` // 添加字段
	elements           []string
}

func redis_unpck() error {
	// 从redis获取任务 creative 信息
	key := "test_creative_video_combination_list"
	var msgComb string
	for {
		msgGet, err := utils.RedisLPop(key)
		if err != nil {
			fmt.Println(err.Error())
			time.Sleep(1000 * time.Millisecond)
			continue
		} else {
			msgComb = msgGet
			break
		}
	}

	fmt.Println(msgComb)
	// 对任务信息进行msgpack解码
	// var jsonMsgComb interface{}
	// var jsonMsgComb map[interface{}]interface{}
	var jsonMsgComb VideoCreativeComb
	err := msgpack.Unmarshal([]byte(msgComb), &jsonMsgComb)
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println("decode ok")
		fmt.Println(jsonMsgComb)
		fmt.Println(jsonMsgComb.AdvertiserID)
		fmt.Println(jsonMsgComb.Height)
	}
	return nil
}
