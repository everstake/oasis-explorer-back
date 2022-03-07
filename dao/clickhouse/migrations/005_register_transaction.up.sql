CREATE TABLE IF NOT EXISTS register_node_transactions (
    blk_lvl UInt64,
    tx_time DateTime,
    tx_hash FixedString(64),
    reg_id FixedString(44),
    reg_address FixedString(46),
    reg_entity_id FixedString(44),
    reg_entity_address  FixedString(46),
    reg_expiration UInt32,
    reg_p2p_id FixedString(44),
    reg_consensus_id FixedString(44),
    reg_consensus_address  FixedString(40),
    reg_physical_address String,
    reg_roles UInt32
) ENGINE ReplacingMergeTree()
PARTITION BY toYYYYMMDD(tx_time)
ORDER BY (tx_hash);

CREATE TABLE IF NOT EXISTS register_entity_transactions (
    blk_lvl UInt64,
    tx_time DateTime,
    tx_hash FixedString(64),
    reg_entity_id FixedString(44),
    reg_entity_address FixedString(46),
    reg_nodes Array(FixedString(44)),
    reg_allow_entity_signed_nodes UInt8 --Bool
) ENGINE ReplacingMergeTree()
PARTITION BY toYYYYMMDD(tx_time)
ORDER BY (blk_lvl);