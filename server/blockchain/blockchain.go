package blockchain

import (
	"crypto/sha256"
	"encoding/hex"
)

// 区块
type Block struct {
	Index      int    `json:"index"`      // 区块链中的第几块（从0开始递增）
	Timestamp  string `json:"timestamp"`  // 时间字符串（格式要求：2021-03-15 03:15:30）
	Content    string `json:"content"`    // 区块内容（用户想保存到链上的内容）
	Hash       string `json:"hash"`       // 区块哈希（Index+Timestamp+Content+PrevHash+Difficulty+Nonce做DoubleSha256）
	PrevHash   string `json:"prev_hash"`  // 前一个区块的哈希（取自前一个区块）
	Difficulty int    `json:"difficulty"` // 满足几个前导0（取自全局变量）
	Nonce      string `json:"nonce"`      // 碰撞随机值（任意值）
}

// 双重SHA256哈希算法：用于生成Block.Hash
func DoubleSha256(raw []byte) (enc string) {
	hash := sha256.Sum256(raw)
	hash = sha256.Sum256(hash[:])
	enc = hex.EncodeToString(hash[:])
	return
}
