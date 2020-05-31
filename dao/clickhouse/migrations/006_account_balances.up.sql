CREATE TABLE IF NOT EXISTS account_balance (
    blk_lvl UInt64,
    blk_time DateTime,
    acb_account FixedString(44),
    acb_nonce UInt64,
    acb_general_balance String,
    acb_escrow_balance_active String,
    acb_escrow_balance_share String,
    acb_escrow_debonding_active String,
    acb_escrow_debonding_share String
) ENGINE ReplacingMergeTree()
PARTITION BY toYYYYMMDD(blk_time)
ORDER BY (acb_account,blk_lvl);
