package help

import (
	"strconv"

	"github.com/google/uuid"
	"github.com/sony/sonyflake"
)

// Uuid 生成新的 UUID v4 字符串。
// 返回格式为 "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx" 的字符串。
//
// Deprecated: 数据库主键请使用 Uuid7()，UUIDv7 因其时间有序特性
// 具有更好的索引性能。
func Uuid() string {
	return uuid.New().String()
}

// Uuid7 生成新的 UUID v7 字符串。
// UUIDv7 按时间排序，推荐用作数据库主键。
// 相比 UUIDv4 的优势：
//   - 时间有序：新 ID 排在旧 ID 之后
//   - 索引友好：顺序插入减少 B-tree 页分裂
//   - 时间戳可提取：可从 ID 中推导出创建时间
//
// 返回格式为 "xxxxxxxx-xxxx-7xxx-xxxx-xxxxxxxxxxxx" 的字符串。
func Uuid7() string {
	id, err := uuid.NewV7()
	if err != nil {
		// 若 v7 生成失败则回退到 v4（极为罕见）
		return uuid.New().String()
	}
	return id.String()
}

// MustUuid7 生成新的 UUID v7 字符串，出错时 panic。
// 适用于 UUID 生成失败即视为致命错误的场景。
func MustUuid7() string {
	return uuid.Must(uuid.NewV7()).String()
}

// Uuid7Time 从 UUID v7 字符串中提取时间戳。
// 成功时返回 Unix 时间戳（毫秒）和 true，
// 若 UUID 不是有效的 v7 则返回 0 和 false。
func Uuid7Time(s string) (int64, bool) {
	id, err := uuid.Parse(s)
	if err != nil {
		return 0, false
	}
	if id.Version() != 7 {
		return 0, false
	}
	// UUIDv7 将 Unix 毫秒时间戳存储在前 48 位
	ts := int64(id[0])<<40 | int64(id[1])<<32 | int64(id[2])<<24 |
		int64(id[3])<<16 | int64(id[4])<<8 | int64(id[5])
	return ts, true
}

// SF 是用于生成分布式唯一 ID 的全局 Sonyflake 实例。
// Sonyflake 生成 63 位唯一 ID，设计灵感来自 Twitter 的 Snowflake。
//
// 默认配置：
//   - MachineID：私有 IP 地址的低 16 位
//   - StartTime：2014-09-01 00:00:00 UTC（Sonyflake 默认值）
//
// 在容器化/分布式环境中，多个实例可能共享同一 IP（如 Kubernetes pod），
// 应在应用启动时用自定义配置的实例替换 SF：
//
//	func init() {
//	    help.SF = sonyflake.NewSonyflake(sonyflake.Settings{
//	        MachineID: func() (uint16, error) {
//	            // Use pod name hash, environment variable, or other unique identifier
//	            id, _ := strconv.Atoi(os.Getenv("MACHINE_ID"))
//	            return uint16(id), nil
//	        },
//	    })
//	}
//
// 注意：若初始化失败 SF 可能为 nil（罕见，通常仅当机器没有有效的私有 IP 时发生）。
// SID() 和 SIDWithError() 会妥善处理这种情况。
var SF = sonyflake.NewSonyflake(sonyflake.Settings{})

// SID 生成新的 Sonyflake ID 字符串。
// 若 ID 生成失败（如 SF 为 nil 或时钟溢出）则返回空字符串。
//
// 生产环境中若需显式处理错误，请改用 SIDWithError()。
func SID() string {
	if SF == nil {
		return ""
	}
	id, err := SF.NextID()
	if err != nil {
		return ""
	}
	return strconv.FormatUint(id, 10)
}

// SIDWithError 生成新的 Sonyflake ID 字符串。
// 返回 ID 以及生成过程中发生的错误。
//
// 可能的错误：
//   - ErrSonyflakeNil：SF 为 nil（初始化失败）
//   - 时钟溢出：时间超出 Sonyflake 自 StartTime 起 174 年的上限
//   - 时钟回拨：系统时间被调整
func SIDWithError() (string, error) {
	if SF == nil {
		return "", ErrSonyflakeNil
	}
	id, err := SF.NextID()
	if err != nil {
		return "", err
	}
	return strconv.FormatUint(id, 10), nil
}

// ErrSonyflakeNil 在 SF 为 nil（Sonyflake 初始化失败）时返回。
var ErrSonyflakeNil = &sonyflakeError{"sonyflake: not initialized, SF is nil"}

type sonyflakeError struct {
	msg string
}

func (e *sonyflakeError) Error() string {
	return e.msg
}
