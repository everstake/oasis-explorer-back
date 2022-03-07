CREATE MATERIALIZED VIEW  IF NOT EXISTS entity_depositors_view
(
 tx_receiver FixedString(46),
 tx_sender FixedString(46),
 escrow_since DateTime CODEC(DoubleDelta, ZSTD(13)),
 input UInt64 CODEC(DoubleDelta, LZ4),
 output UInt64 CODEC(DoubleDelta, LZ4),
 balance UInt64 CODEC(DoubleDelta, LZ4) 
) 
ENGINE = MergeTree
PARTITION BY toYYYYMM(escrow_since)
ORDER BY (tx_receiver)
SETTINGS index_granularity = 8192 POPULATE

AS
(select *, add.input - remove.output balance from (
        select tx_receiver, tx_sender, min(tx_time) escrow_since, sum(tx_escrow_amount) input
        from transactions
        where tx_type = 'addescrow' and tx_status = 1
        group by  tx_receiver, tx_sender) add
        ANY
        LEFT JOIN (
    select tx_receiver, tx_sender, sum(tx_escrow_reclaim_amount) output
    from transactions
    where tx_type = 'reclaimescrow' and tx_status = 1
    group by tx_receiver, tx_sender) remove USING tx_receiver, tx_sender
);

     
CREATE VIEW IF NOT EXISTS entity_active_depositors_counter_view AS
  SELECT tx_receiver reg_entity_address, count() depositors_num from entity_depositors_view
  where balance > 0
  group by tx_receiver;


CREATE MATERIALIZED VIEW IF NOT EXISTS entity_nodes_view -- need to replace body from newer script from 010_validators.up.sql
(
    reg_entity_id FixedString(44),
    reg_entity_address FixedString(46),
    reg_id FixedString(44),
    reg_address FixedString(46),
    reg_consensus_address FixedString(40),
    created_time DateTime CODEC(DoubleDelta, ZSTD(13)),
    blk_lvl UInt64 CODEC(DoubleDelta),
    reg_expiration UInt32 ,
    last_block_time DateTime CODEC(DoubleDelta, ZSTD(13)),
    blocks UInt32 ,
    last_signature_time DateTime CODEC(DoubleDelta, ZSTD(13)),
    signed_blocks UInt32,
    signatures UInt32 
)
ENGINE = ReplacingMergeTree
PARTITION BY toYYYYMM(created_time)
ORDER BY (reg_entity_address)
SETTINGS index_granularity = 8192 
POPULATE
AS
(
  SELECT
    reg_entity_id, reg_entity_address ,reg_id, reg_address, reg_consensus_address,
     created_time, blk_lvl, reg_expiration, last_block_time, signed_blocks AS blocks, last_signature_time, signed_blocks,
     signatures    
FROM
(
    SELECT *
    FROM
    (
        SELECT
            reg_entity_id,
            reg_id,
            reg_consensus_address,
            min(tx_time) AS created_time,
            max(blk_lvl) AS blk_lvl,
            max(reg_expiration) AS reg_expiration,
            reg_entity_address,
            reg_address
        FROM register_node_transactions
        GROUP BY
            reg_entity_id,
            reg_id,
            reg_consensus_address,
            reg_entity_address,
            reg_address
    ) AS nodes
    ANY LEFT JOIN
    (
        SELECT
            blk_proposer_address AS reg_consensus_address,
            max(blk_created_at) AS last_block_time,
            count() AS signed_blocks
        FROM blocks
        GROUP BY blk_proposer_address
    ) AS blk USING (reg_consensus_address)
) AS prep
ANY LEFT JOIN
(
    SELECT
        sig_validator_address AS reg_consensus_address,
        max(sig_timestamp) AS last_signature_time,
        count() AS signatures
    FROM block_signatures
    GROUP BY sig_validator_address
) AS blk USING (reg_consensus_address)
)