CREATE MATERIALIZED VIEW IF NOT EXISTS validator_blocks_count_mv
  ENGINE = SummingMergeTree()
    PARTITION BY toYYYYMM(day)
    ORDER BY (reg_consensus_address, day) POPULATE AS
select blk_proposer_address reg_consensus_address,
       toStartOfDay(blk_created_at) day,
       max(blk_created_at) last_block_time,
       count()                      blocks
from blocks
group by reg_consensus_address, day;

CREATE VIEW IF NOT EXISTS validator_blocks_day_count_view AS
select reg_consensus_address,
       day,
       max(last_block_time) last_block_time,
       sum(blocks)             blocks
from validator_blocks_count_mv
group by reg_consensus_address, day;

CREATE VIEW IF NOT EXISTS validator_blocks_count_view AS
SELECT reg_consensus_address,
       max(last_block_time) last_block_time,
       sum(blocks) blocks
FROM validator_blocks_count_mv
GROUP BY reg_consensus_address;

CREATE MATERIALIZED VIEW IF NOT EXISTS validator_block_signatures_count_mv
  ENGINE = SummingMergeTree()
    PARTITION BY toYYYYMM(day)
    ORDER BY (reg_consensus_address, day) POPULATE AS
select sig_validator_address reg_consensus_address,
       toStartOfDay(sig_timestamp) day,
       max(sig_timestamp) last_signature_time,
       count(distinct (blk_lvl))   blocks,
       count()                     signatures
from block_signatures
group by reg_consensus_address, day;

CREATE VIEW IF NOT EXISTS validator_block_signatures_day_count_view AS
select reg_consensus_address,
       day,
       max(last_signature_time) last_signature_time,
       sum(blocks)   day_signed_blocks,
       sum(signatures)                     day_signatures
from validator_block_signatures_count_mv
group by reg_consensus_address, day;

CREATE VIEW IF NOT EXISTS validator_block_signatures_count_view AS
SELECT reg_consensus_address,
       max(last_signature_time) last_signature_time,
       sum(blocks)     signed_blocks,
       sum(signatures) signatures
FROM validator_block_signatures_count_mv
GROUP BY reg_consensus_address;

DROP TABLE entity_nodes_view;
CREATE MATERIALIZED VIEW IF NOT EXISTS entity_nodes_view AS --OMG. seems like here is rewrited query from old 008_entity_nodes_vew.sql
select *
from (
       select *
       from (
              --Group all register txs by entity and node
              select reg_entity_id,
                     reg_entity_address,
                     reg_id,
                     reg_address,
                     reg_consensus_address,
                     min(tx_time)        created_time,
                     max(blk_lvl)        blk_lvl,
                     max(reg_expiration) reg_expiration
              from register_node_transactions
              group by reg_entity_id, reg_entity_address, reg_address, reg_id, reg_consensus_address
              ) nodes
              ANY
              LEFT JOIN validator_blocks_count_view USING reg_consensus_address
       ) prep
       ANY
       LEFT JOIN validator_block_signatures_count_view USING reg_consensus_address;

CREATE TABLE IF NOT EXISTS public_validators
(
  partition UInt64 DEFAULT 1,
  reg_entity_id FixedString(44),
  reg_entity_address FixedString(46),
  pvl_name      String,
  pvl_info   String
) ENGINE ReplacingMergeTree()
    PARTITION BY partition
    ORDER BY (reg_entity_id);

CREATE VIEW IF NOT EXISTS day_max_block_lvl_view AS
select toStartOfDay(blk_created_at) day, count() blk_count, max(blk_lvl) blk_lvl
from blocks
group by day;

CREATE VIEW IF NOT EXISTS validator_day_stats_view AS
select *
      from validator_block_signatures_day_count_view ANY
             LEFT JOIN validator_blocks_day_count_view USING reg_consensus_address, day;

CREATE VIEW IF NOT EXISTS validator_entity_view AS
select p.*, day_max_block_lvl_view.blk_lvl max_day_block, day_max_block_lvl_view.blk_count day_blocks
from (
       select *
       from (
              select reg_entity_address,
                     anyLast(reg_consensus_address) reg_consensus_address,
                     anyLast(reg_address)           node_address,
                     max(blk_lvl)                   blk_lvl,
                     toStartOfDay(now())            day,
                     min(created_time)              created_time,
                     max(reg_expiration)            reg_expiration,
                     max(last_block_time)           last_block_time,
                     sum(blocks)                    blocks,
                     sum(signed_blocks)             signed_blocks,
                     sum(signatures)                signatures
              from entity_nodes_view
              GROUP BY reg_entity_address) g
              any
              left join (
         select reg_consensus_address, day_signatures, day_signed_blocks
         from validator_day_stats_view
         where day = toStartOfDay(now())) sigs USING reg_consensus_address) p
       ANY
       LEFT JOIN day_max_block_lvl_view b USING day;


CREATE VIEW IF NOT EXISTS validators_list_view AS
select *
from (
       select *, CASE WHEN reg_expiration >= (select max(blk_epoch) from blocks) THEN 1 ELSE 0 END is_active
       from (
              SELECT *
              FROM validator_entity_view -- created MVIEW to don't make fullscan of txs
                    ANY
                    LEFT JOIN
                            ( select reg_entity_address,
                            min(blk_lvl) start_blk_lvl
                            from register_node_transactions -- 420 entries
                            group by reg_entity_address ) val_lvl USING reg_entity_address) validator
              ANY
              LEFT JOIN (SELECT acb_account reg_entity_address, acb_escrow_balance_active, 
              acb_general_balance, acb_escrow_balance_share, acb_escrow_debonding_active,
               acb_delegations_balance , acb_debonding_delegations_balance ,
                acb_self_delegation_balance, acb_commission_schedule ,
                 depositors_num
                         from account_last_balance_view ANY -- this is an aggregation from mview, haven't any optimizations yet
                                LEFT JOIN entity_active_depositors_counter_view -- fixed to make aggregate from MVIEW
                                USING reg_entity_address -- entity_active_depositors_counter_view  is aggregate upon `entity_depositors_view`, a very huge select with joins. 
         ) b USING reg_entity_address
       where blocks > 0
          OR signatures > 0
          OR reg_expiration >= (select max(blk_epoch) from blocks)) prep
       ANY
       LEFT JOIN public_validators USING reg_entity_address;