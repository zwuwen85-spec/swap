// 商品类型
export const GoodsType = {
  SALE: 1,      // 售卖
  EXCHANGE: 2,  // 交换
  BOTH: 3       // 均可
}

// 商品状态
export const GoodsStatus = {
  OFFLINE: 0,   // 下架
  ON_SALE: 1,   // 在售
  SOLD: 2,      // 已售
  EXCHANGED: 3  // 已交换
}

// 商品成色
export const GoodsCondition = {
  NEW: 1,        // 全新
  LIKE_NEW: 2,   // 九成新
  GOOD: 3,       // 八成新
  FAIR: 4        // 七成新
}

// 交换状态
export const ExchangeStatus = {
  PENDING: 0,    // 待处理
  ACCEPTED: 1,   // 已接受
  REJECTED: 2,   // 已拒绝
  CANCELLED: 3,  // 已取消
  COMPLETED: 4   // 已完成
}

// 用户状态
export const UserStatus = {
  DISABLED: 0,   // 禁用
  NORMAL: 1,     // 正常
  FROZEN: 2      // 冻结
}
