package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/owenliang/blockchain/server/blockchain"
)

// 上传Block的HTTP应答
type UploadBlockResult struct {
	Error string `json:"error"`
}

// 区块链
var ChainMutex sync.Mutex
var Blockchain = []*blockchain.Block{}
var BlockFile *os.File // 持久化用

// 下载最新区块
func HandleLatest(w http.ResponseWriter, req *http.Request) {
	ChainMutex.Lock()
	block := Blockchain[len(Blockchain)-1]
	ChainMutex.Unlock()

	var blockJson []byte
	var err error
	if blockJson, err = json.Marshal(block); err != nil {
		w.WriteHeader(500)
		return
	}
	w.Write(blockJson)
}

// 校验区块
func ValidateBlock(prevBlock *blockchain.Block, newBlock *blockchain.Block) (err error) {
	// Index越界
	if newBlock.Index > prevBlock.Index+1 || newBlock.Index < 0 {
		err = fmt.Errorf("区块Index越界，请检查索引=%v", newBlock.Index)
		return
	}
	// 区块已被提交
	if newBlock.Index != prevBlock.Index+1 {
		err = fmt.Errorf("该区块已被别人提交, 索引=%v", newBlock.Index)
		return
	}
	// 时间字段格式检测
	if _, err = time.Parse("2006-01-02 15:04:05", newBlock.Timestamp); err != nil {
		err = fmt.Errorf("该区块的Timestamp错误，请检查一下格式：%v", err.Error())
		return
	}
	// Difficulty检查
	if (newBlock.Index%2 == 0 && newBlock.Difficulty != 5) || (newBlock.Index%2 == 1 && newBlock.Difficulty != 6) {
		err = fmt.Errorf("该区块的Difficulty错误，偶数Index应该传2，奇数Index应该传3")
		return
	}
	// PrevHash检查
	if newBlock.PrevHash != prevBlock.Hash {
		err = fmt.Errorf("prevHash错误，你传的是%v，实际上是%v", newBlock.PrevHash, prevBlock.Hash)
		return
	}
	// Hash检查
	beforeHash := fmt.Sprintf("%v%v%v%v%v%v", newBlock.Index, newBlock.Timestamp, newBlock.Content, newBlock.PrevHash, newBlock.Difficulty, newBlock.Nonce)
	hash := blockchain.DoubleSha256([]byte(beforeHash))
	if hash != newBlock.Hash {
		err = fmt.Errorf("hash错误, 你传的是%v，实际上是%v", newBlock.Hash, hash)
		return
	}
	// Hash前导0检查
	s := "00000"
	if newBlock.Difficulty == 6 {
		s = "000000"
	}
	if !strings.HasPrefix(hash, s) {
		err = fmt.Errorf("hash不满足Difficulty，你传的hash是%v，需要前导0个数为%v", hash, newBlock.Difficulty)
		return
	}
	return
}

// 上传区块
func HandleUpload(w http.ResponseWriter, req *http.Request) {
	var err error

	req.ParseForm()
	blockJson := req.PostForm.Get("block")

	result := &UploadBlockResult{}
	newBlock := &blockchain.Block{}

	// 解析Block
	if err = json.Unmarshal([]byte(blockJson), newBlock); err != nil {
		result.Error = "解析Block JSON失败：" + err.Error()
		goto RET
	}

	// 校验Block
	ChainMutex.Lock()
	if err = ValidateBlock(Blockchain[len(Blockchain)-1], newBlock); err != nil {
		result.Error = "校验Block失败：" + err.Error()
		ChainMutex.Unlock()
		goto RET
	}
	// 添加Block到链
	SaveBlock(newBlock)
	ChainMutex.Unlock()
RET:
	var resultJson []byte
	if resultJson, err = json.Marshal(result); err != nil {
		w.WriteHeader(500)
	} else {
		w.Write(resultJson)
	}
}

// 持久化Block
func SaveBlock(block *blockchain.Block) {
	Blockchain = append(Blockchain, block)
	row, _ := json.Marshal(block)
	BlockFile.WriteString(string(row) + "\n")
	BlockFile.Sync() // 强刷到磁盘上
}

// 初始化区块链
func InitBlockChain() (err error) {
	// 创建或者追加
	if BlockFile, err = os.OpenFile("./blockchain.txt", os.O_RDWR|os.O_CREATE, 0666); err != nil {
		fmt.Println(err)
		return
	}

	// 指向文件末尾，看一下文件是否为空
	offset, _ := BlockFile.Seek(0, os.SEEK_END)
	if offset == 0 {
		// 创世区块（创世块的Hash字段任意生成，不需要参与校验）
		SaveBlock(&blockchain.Block{
			Index:      0,
			Timestamp:  "",
			Content:    "",
			Hash:       blockchain.DoubleSha256([]byte("hello blockchain")),
			PrevHash:   "",
			Difficulty: 0,
			Nonce:      "",
		})
	} else { // 已经有数据，恢复区块链
		BlockFile.Seek(0, os.SEEK_SET)
		scanner := bufio.NewScanner(BlockFile)
		for scanner.Scan() {
			row := scanner.Text()
			block := &blockchain.Block{}
			if err = json.Unmarshal([]byte(row), &block); err != nil {
				return
			}
			Blockchain = append(Blockchain, block)
		}
		err = scanner.Err()
	}
	return
}

func main() {
	var err error
	if err = InitBlockChain(); err != nil {
		return
	}
	http.HandleFunc("/blockchain/latest", HandleLatest)
	http.HandleFunc("/blockchain/upload", HandleUpload)
	http.ListenAndServe(":6543", nil)
}
