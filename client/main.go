package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	"github.com/owenliang/blockchain/client/blockchain"
)

// 上传Block的HTTP应答
type UploadBlockResult struct {
	Error string `json:"error"`
}

// ----------《作业A》：解析json -> *Block
// 注意：根据服务端返回的实际JSON，给Block结构体字段配置正确的tag
func Json2Block(blockJson []byte) (block *blockchain.Block, err error) {
	err = errors.New("[请实现Json2Block函数] 你需要解析服务端返回的JSON并保存到Block对象，当前JSON为：" + string(blockJson))
	return
}

// ----------《作业B》：序列化*Block -> json
func Block2Json(block *blockchain.Block) (blockJson []byte, err error) {
	err = fmt.Errorf("[请实现Block2Json函数] 你需要将Block对象序列化为JSON，当前Block对象为%v", block)
	return
}

// 下载最新Block
func FetchLatestBlock(address string) (block *blockchain.Block, err error) {
	// 拼接HTTP接口地址
	api := fmt.Sprintf("http://%s/blockchain/latest", address)

	// 请求HTTP接口
	var resp *http.Response
	if resp, err = http.Get(api); err != nil {
		return
	}
	defer resp.Body.Close()

	// 读取返回的JSON
	var blockJson []byte
	if blockJson, err = ioutil.ReadAll(resp.Body); err != nil {
		return
	}

	// 解析JSON为*Block对象
	if block, err = Json2Block(blockJson); err != nil {
		return
	}
	fmt.Printf("下载到最新区块：%v，解析Block为：%v\n", string(blockJson), *block)
	return
}

// 《作业B》：实现挖矿函数，传入前一个Block，循环计算直到挖出下一个Block
func MiningBlock(prevBlock *blockchain.Block, content string) (newBlock *blockchain.Block) {
	// 服务端会校验提交Block的如下信息：
	// 1，Index：上一个Block的Index+1
	// 2，Timestamp：注意格式
	// 3，Hash：必须正确，计算方法：Index+Timestamp+Content+PrevHash+Difficulty+Nonce做DoubleSha256
	// 4，PrevHash：必须是Index-1区块的Hash
	// 5，Difficulty：如果该Index是偶数，那么Difficulty设置5（意味着你的Hash需要有5个前导0），否则Difficulty设置为6（意味着你的Hash需要有6个前导0），你的Hash值必须符合这个要求
	// 6，Nonce：随意内容，你需要调整它来用来碰撞到符合要求的Hash值
	return
}

// 上传新Block
func UploadNewBlock(address string, block *blockchain.Block) (err error) {
	api := fmt.Sprintf("http://%s/blockchain/upload", address)

	// 序列化JSON
	var blockJson []byte
	if blockJson, err = Block2Json(block); err != nil {
		return
	}

	// 填写POST表单
	params := url.Values{}
	params.Set("block", string(blockJson))

	// POST上传服务端
	var resp *http.Response
	if resp, err = http.PostForm(api, params); err != nil {
		return
	}
	defer resp.Body.Close()

	// 读取返回的JSON
	var resultJson []byte
	if resultJson, err = ioutil.ReadAll(resp.Body); err != nil {
		return
	}

	// 解析应答
	uploadResult := &UploadBlockResult{}
	if err = json.Unmarshal(resultJson, uploadResult); err != nil {
		return
	}

	// 判断提交结果
	if uploadResult.Error != "" {
		err = errors.New(uploadResult.Error)
	}
	return
}

// 命令行参数
var (
	address string
	content string
)

func main() {
	// 解析命令行参数
	flag.StringVar(&address, "address", "", "填写矿池地址")
	flag.StringVar(&content, "content", "", "写入到区块的内容（相当于交作业，需要包含你的名字作为标识）")
	flag.Parse()

	// 开始挖矿
	fmt.Printf("已加入矿池:%v, 启动挖矿...\n", address)
	for {
		var err error
		var prevBlock *blockchain.Block
		var newBlock *blockchain.Block

		// 下载最新Block
		if prevBlock, err = FetchLatestBlock(address); err != nil {
			fmt.Println(err)
			goto FAIL
		}

		// 碰撞下一个Block
		fmt.Println("正在挖下一个Block...请关注CPU变化!")
		newBlock = MiningBlock(prevBlock, content)
		if newBlock == nil {
			fmt.Println("[请实现MiningBlock函数]你挖到的Block为nil，是不是忘记实现MiningBlock函数了？")
			goto FAIL
		}

		// 提交区块到服务端
		if err = UploadNewBlock(address, newBlock); err != nil {
			fmt.Println("上传区块失败：", err)
			goto FAIL
		}
		fmt.Printf("区块已提交, 索引=%v 哈希=%v 内容=%v\n", newBlock.Index, newBlock.Hash, newBlock.Content)
		return
	FAIL:
		// 间隔一下重试
		time.Sleep(1 * time.Second)
	}
}
