package blockchain

import (
	"crypto/sha256"
	"encoding/hex"
)

// 区块
type Block struct {
	Index      int    // 区块链中的第几块（从0开始递增）
	Timestamp  string // 时间字符串（格式要求：2021-03-15 03:15:30）
	Content    string // 区块内容（用户想保存到链上的内容）
	Hash       string // 区块哈希（Index+Timestamp+Content+PrevHash+Difficulty+Nonce做DoubleSha256），提示：+表示连接它们，不是真的加号
	PrevHash   string // 前一个区块的哈希（取自前一个区块的Hash字段）
	Difficulty int    // 满足几个前导0
	Nonce      string // 碰撞随机值（任意值）
}

// 双重SHA256哈希算法：用于生成Block.Hash
func DoubleSha256(raw []byte) (enc string) {
	hash := sha256.Sum256(raw)
	hash = sha256.Sum256(hash[:])
	enc = hex.EncodeToString(hash[:])
	return
}
