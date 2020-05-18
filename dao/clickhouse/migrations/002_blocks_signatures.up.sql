CREATE TABLE IF NOT EXISTS block_signatures (
    sig_blk_lvl UInt64,
    sig_timestamp DateTime,
    sig_block_id_flag UInt64,
    sig_validator_address FixedString(40),
    sig_blk_signature FixedString(128)
) ENGINE ReplacingMergeTree()
PARTITION BY sig_blk_lvl
ORDER BY (sig_blk_lvl, sig_validator_address);