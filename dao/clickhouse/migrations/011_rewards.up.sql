CREATE TABLE IF NOT EXISTS rewards (
    blk_lvl UInt64,
    blk_epoch UInt64,
    created_at DateTime,
    rwd_amount UInt64,
    reg_entity_address  FixedString(46)
) ENGINE ReplacingMergeTree()
PARTITION BY toYYYYMMDD(blk_created_at)
ORDER BY (reg_entity_address, blk_lvl);