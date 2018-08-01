package main

import (
	"fmt"
	// 另外一个msgpack的包
	"development/msgpack-demo/utils"
	// "flag"
	"github.com/vmihailenco/msgpack"
	"strconv"
	"time"
)

func init() {
	//初始化解析用户参数
	// flag.Parse()
}

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
// type VideoCreativeComb struct {
// 	CreativeID         string  `msgpack:"creative_id"`
// 	SizeID             string  `msgpack:"size_id"`
// 	AdvertiserID       string  `msgpack:"advertiser_id"`
// 	CombinationID      string  `msgpack:"combination_id"`
// 	FileType           string  `msgpack:"file_type"`
// 	Height             int32   `msgpack:"height"`
// 	Width              int32   `msgpack:"width"`
// 	VideoBitRate       int32   `msgpack:"video_bit_rate"`
// 	AudioBitRate       int32   `msgpack:"audio_bit_rate"`
// 	Sar                float64 `msgpack:"sar_pixel"`
// 	VideoCombinationID string  `msgpack:"dynamic_combination_id"`
// 	Version            string  `msgpack:"version"`
// 	AudioEnc           string  `msgpack:"audio_encoder"` // 添加字段
// 	elements           []string
// }

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
		fmt.Println("1111 decode ok")
		fmt.Println(jsonMsgComb)
		fmt.Println(jsonMsgComb.AdvertiserID)
		fmt.Println(jsonMsgComb.Height)
		//获取
		err = GetCreativeTemplate(jsonMsgComb)
		if err == nil {
			fmt.Println(" 2222 decode ok")
		}
	}
	return nil
}

func GetCreativeTemplate(videoCreateCombination VideoCreativeComb) error {
	var videoCreativeTemplate VideoCreativeTemp
	// 用creative_id +size_id拉取视频合成模板数据和元素数据
	fieldCreativeSize := videoCreateCombination.CreativeID + "_" + videoCreateCombination.SizeID
	key := "creative_video"
	msgCreativeTemplate, err := utils.RedisHGet(key, fieldCreativeSize)
	if err != nil {
		fmt.Printf("redis HGet failed,err:%s\n", err.Error())
		return err
	}
	//msgpack
	err = msgpack.Unmarshal([]byte(msgCreativeTemplate), &videoCreativeTemplate)
	if err != nil {
		fmt.Printf("msgpack unmarshal failed,err:%s\n", err.Error())
		return err
	}
	fmt.Printf("%v\n", videoCreativeTemplate)

	return nil
}

////////////////////////////////////////////////////////////////////////
// 主要用于从redis获取的msg消息解码以后的格式

// 元素值
type ElemValue struct {
	Value    string  `msgpack:"value"`
	Duration float32 `msgpack:"duration"`
	Erasure  bool    `msgpack:"erasure"`
	Fade     bool    `msgpack:"fade"`
	Volume   int32   `msgpack:"volume"`
	Type     string  //存储类型值
	Offset   float32 //存储偏移值
	ElemId   string  //元素id
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
	Sar                float64 `msgpack:"sar_pixel"`
	VideoCombinationID string  `msgpack:"dynamic_combination_id"`
	Version            string  `msgpack:"version"`
	AudioEnc           string  `msgpack:"audio_encoder"` // 添加字段
	elements           []string
}

type FontInfo struct {
	FontFamily    string `msgpack:"font-family"`
	FontSize      int32  `msgpack:"font-size"`
	TextAlign     string `msgpack:"text-align"`
	VerticalAlign string `msgpack:"vertical-align"`
}

type Source struct {
	ElementID int32  `msgpack:"element_id"`
	Width     int32  `msgpack:"width"`
	Height    int32  `msgpack:"height"`
	URL       string `msgpack:"url"`
}

type AnimateInfo struct {
	Type   string `msgpack:"type"`
	Delay  int32  `msgpack:"delay"`
	Appear struct {
		Effect   string `msgpack:"effect"`
		Duration string `msgpack:"duration"`
	} `msgpack:"appear"`
	Disappear struct {
		Effect   string `msgpack:"effect"`
		Duration string `msgpack:"duration"`
	} `msgpack:"disappear"`
}

type Options struct {
	X           int32         `msgpack:"x"`
	Y           int32         `msgpack:"y"`
	RotateX     int32         `msgpack:"rotateX"`
	RotateY     int32         `msgpack:"rotateY"`
	Opacity     int32         `msgpack:"opacity"`
	Color       string        `msgpack:"color"`
	Font        FontInfo      `msgpack:"font"`
	Src         Source        `msgpack:"src"`
	W           int32         `msgpack:"w"`
	H           int32         `msgpack:"h"`
	Objposition interface{}   `msgpack:"objposition"`
	Start       int32         `msgpack:"start"`
	End         int32         `msgpack:"end"`
	Shadow      []interface{} `msgpack:"shadow,asArray"`
	Video       []interface{} `msgpack:"video,asArray"`
	Animate     AnimateInfo   `msgpack:"animate"`
	Appear      int32         `msgpack:"appear"`
	Disappear   int32         `msgpack:"disappear"`
	Tag         string        `msgpack:"tag"`
	ZAxis       int32         `msgpack:"z"`
}

type ObjsInfo struct {
	Type   string  `msgpack:"type"`
	UID    string  `msgpack:"uid"`
	Option Options `msgpack:"options"`
}

type LayoutInfo struct {
	TotalDuration string `msgpack:"total_duration"`
	KzSystemRatio struct {
		HeightRatio float64 `msgpack:"heightRatio"`
		WidthRatio  float64 `msgpack:"widthRatio"`
	} `msgpack:"kzSystemRatio"`
}

//视频创意模板信息
type VideoCreativeTemp struct {
	AdvertiserID      string        `msgpack:"advertiser_id"`
	SizeID            string        `msgpack:"size_id"`
	CreativeID        string        `msgpack:"creative_id"`
	AdID              string        `msgpack:"ad_id"`
	OriginWidth       string        `msgpack:"width"`
	OriginHeight      string        `msgpack:"height"`
	CampaignID        int32         `msgpack:"campaign_id"`
	TemplateID        string        `msgpack:"template_id"`
	SourceTemplateID  string        `msgpack:"source_template_id"`
	Objs              []ObjsInfo    `msgpack:"objs,asArray"`
	Layout            LayoutInfo    `msgpack:"layout"`
	CompressQuality   string        `msgpack:"compress_quality"`
	CompressSizeLimit string        `msgpack:"compress_size_limit"`
	Background        []interface{} `msgpack:"background,asArray"`
}
