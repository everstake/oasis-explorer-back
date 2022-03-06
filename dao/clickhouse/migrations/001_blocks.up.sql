CREATE TABLE IF NOT EXISTS blocks (
    blk_lvl UInt64,
    blk_created_at DateTime,
    blk_hash FixedString(64),
    blk_proposer_address FixedString(40),
    blk_validator_hash FixedString(64),
    blk_epoch UInt64
) ENGINE MergeTree()
PARTITION BY toYYYYMMDD(blk_created_at)
ORDER BY (blk_lvl);