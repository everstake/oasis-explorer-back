CREATE TABLE IF NOT EXISTS transactions (
    tx_blk_lvl UInt64,
    tx_time DateTime,
    tx_hash FixedString(64),
    tx_amount Decimal(36, 18),
    tx_escrow_amount Decimal(36, 18),
    tx_escrow_reclaim_amount Decimal(36, 18),
    tx_escrow_account FixedString(64),
    tx_type  FixedString(64),
    tx_sender FixedString(64),
    tx_receiver FixedString(64),
    tx_nonce UInt64,
    tx_fee UInt64,
    tx_gas_limit UInt64,
    tx_gas_price UInt64
) ENGINE ReplacingMergeTree()
PARTITION BY toYYYYMMDD(tx_time)
ORDER BY (tx_hash);