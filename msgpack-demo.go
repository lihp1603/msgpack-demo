package main

import (
	"bytes"
	"development/msgpack-demo/utils"
	"fmt"
	msgpack "github.com/msgpack/msgpack-go"
	"github.com/ugorji/go/codec"
	"reflect"
	"strconv"
	"time"
)

//这个官方的msgpack不支持tag操作
//虽然能解码，但很难提取出来，需要自己进行部分复杂的代码才行

func main() {

	// msgpck_demo()

	// testMsgPack()

	redis_init()
	defer redis_destory()
	// redis_unpck()
	// redis_unpck2()
	redis_unpck3()
}

func msgpck_demo() {
	packStr := "{\"null\":\"stdClass\",\"creative_id\":\"7660870\",\"size_id\":\"499\",\"advertiser_id\":\"8\",\"combination_id\":\"1398460,1398459,0,0,0,0,0,0,0,0,0,1398459,1\",\"file_type\":\"mp4\",\"height\":1080,\"width\":1920,\"video_bit_rate\":4000,\"audio_bit_rate\":128,\"sar_pixel\":0.0,\"dynamic_combination_id\":\"1398459\",\"version\":\"v2\"}"
	fmt.Println(packStr)

	//进行压缩打包
	b := &bytes.Buffer{}
	n, err := msgpack.PackBytes(b, []byte(packStr))
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println(n)
	fmt.Println(b.String())

	//从redis取出来的打印信息
	// bb := `(stdClass«creative_id§7668538§size_id£499­advertiser_id¡8®combination_id׫1488881,1488882,0,0,0,0,0,0,0,0,0,1488881,0©file_type£mp4¦height˄8¥width̀®video_bit_rate͠®audio_bit_ratè©sar_pixelɿ𵣹namic_combination_id§1488881§version¢v2`
	// b := &bytes.Buffer{}
	// n, _ := b.WriteString(bb)
	// fmt.Println(n)
	// // 解码解包
	unpackValue, n, err := msgpack.Unpack(b)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println(n)
	fmt.Println(unpackValue)
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
	CreativeID         string  `msgpack:"creative_id"`
	SizeID             string  `msgpack:"size_id"`
	AdvertiserID       string  `msgpack:"advertiser_id"`
	CombinationID      string  `msgpack:"combination_id"`
	FileType           string  `msgpack:"file_type"`
	Height             int32   `msgpack:"height"`
	Width              int32   `msgpack:"width"`
	VideoBitRate       int32   `msgpack:"video_bit_rate"`
	AudioBitRate       int32   `msgpack:"audio_bit_rate"`
	Sar                float32 `msgpack:"sar_pixel"`
	VideoCombinationID string  `msgpack:"dynamic_combination_id"`
	Version            string  `msgpack:"version"`
	AudioEnc           string  `msgpack:"audio_encoder"` // 添加字段
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
	b := bytes.NewBuffer([]byte(msgComb))

	// // b := bytes.Buffer{}
	// // n, _ := b.WriteString(msgComb)
	// // fmt.Println(n)

	// unpackValue, n, err := msgpack.Unpack(b)
	// if err != nil {
	// 	fmt.Println(err.Error())
	// 	return err
	// }
	// fmt.Println(n)
	// fmt.Println(unpackValue)

	//添加一个codec
	var jsonMsgComb interface{}
	// var jsonMsgComb map[interface{}]interface{}
	// var jsonMsgComb VideoCreativeComb
	mh := codec.MsgpackHandle{}
	mh.StructToArray = true
	dec := codec.NewDecoder(b, &mh)
	// dec := codec.NewDecoderBytes(b, &mh)
	err := dec.Decode(&jsonMsgComb)
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println("decode ok")
		fmt.Println(jsonMsgComb)

		t := reflect.TypeOf(jsonMsgComb)
		fmt.Println(t.String())
		v := reflect.ValueOf(jsonMsgComb)
		fmt.Println(v)
		fmt.Println(t.Kind())

	}

	// utils.LogTraceI(jsonMsgComb)
	return nil
}

func redis_unpck2() error {
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
	b := bytes.NewBuffer([]byte(msgComb))

	//添加一个codec
	var jsonMsgComb map[interface{}]interface{}
	// var jsonMsgComb VideoCreativeComb
	mh := codec.MsgpackHandle{}
	mh.StructToArray = true
	dec := codec.NewDecoder(b, &mh)
	// dec := codec.NewDecoderBytes(b, &mh)
	err := dec.Decode(&jsonMsgComb)
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println("decode ok")
		fmt.Println(jsonMsgComb)

		for key, value := range jsonMsgComb {
			fmt.Println("..............")
			kt := reflect.TypeOf(key)
			fmt.Println(kt.Kind())
			kv := reflect.ValueOf(key)
			fmt.Println(kv)

			vt := reflect.TypeOf(value)
			fmt.Println(vt.Kind())
			vv := reflect.ValueOf(value)
			fmt.Println(vv)
		}

	}

	return nil
}

func redis_unpck3() error {
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
	b := bytes.NewBuffer([]byte(msgComb))

	//添加一个codec
	var jsonMsgComb VideoCreativeComb
	mh := codec.MsgpackHandle{}
	mh.StructToArray = true
	dec := codec.NewDecoder(b, &mh)
	// dec := codec.NewDecoderBytes(b, &mh)
	err := dec.Decode(&jsonMsgComb)
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println("decode ok")
		fmt.Println(jsonMsgComb)
	}
	return nil
}

type UserStruct struct {
	Age  string
	Name string
	ID   int32
}

var (
	b  []byte
	mh codec.MsgpackHandle
)

func testMsgPack() {
	// user := UserStruct{9, "abcd"}
	user := UserStruct{Name: "abcd"}
	//关键调用
	mh.StructToArray = true

	enc := codec.NewEncoderBytes(&b, &mh)
	err := enc.Encode(user)
	if err == nil {
		fmt.Println("data:", b)
	} else {
		fmt.Println("err:", err)
	}

	dec := codec.NewDecoderBytes(b, &mh)
	var new_user UserStruct
	err = dec.Decode(&new_user)
	if err == nil {
		fmt.Println("new_user:", new_user)
		fmt.Println(new_user.Name)
	} else {
		fmt.Println("err:", err)
	}
}
