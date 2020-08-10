CREATE TABLE IF NOT EXISTS transactions (
    blk_lvl UInt64,
    blk_hash FixedString(64),
    tx_time DateTime,
    tx_hash FixedString(64),
    tx_amount UInt64,
    tx_escrow_amount UInt64,
    tx_escrow_reclaim_amount UInt64,
    tx_type  String,
    tx_status UInt8,
    tx_error  String,
    tx_sender FixedString(46),
    tx_receiver FixedString(46),
    tx_nonce UInt64,
    tx_fee UInt64,
    tx_gas_limit UInt64,
    tx_gas_price UInt64
) ENGINE ReplacingMergeTree()
PARTITION BY toYYYYMMDD(tx_time)
ORDER BY (tx_hash);