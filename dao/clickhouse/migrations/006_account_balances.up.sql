CREATE TABLE IF NOT EXISTS account_balance (
    blk_lvl UInt64,
    blk_time DateTime,
    acb_account FixedString(46),
    acb_nonce UInt64,
    acb_general_balance UInt64,
    acb_escrow_balance_active UInt64,
    acb_escrow_balance_share UInt64,
    acb_escrow_debonding_active UInt64,
    acb_escrow_debonding_share UInt64
) ENGINE ReplacingMergeTree()
PARTITION BY toYYYYMMDD(blk_time)
ORDER BY (acb_account,blk_lvl);
