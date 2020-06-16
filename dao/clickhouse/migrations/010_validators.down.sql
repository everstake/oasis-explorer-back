DROP TABLE validator_blocks_count_mv;
DROP TABLE validator_blocks_count_view;
DROP TABLE validator_block_signatures_count_mv;
DROP TABLE validator_block_signatures_count_view;

--Rollback entity_nodes_view
DROP TABLE entity_nodes_view;
CREATE VIEW IF NOT EXISTS entity_nodes_view AS
select *
from (
       select *
       from (
              --Group all register txs by entity and node
              select reg_entity_id, reg_id, reg_consensus_address, min(tx_time) created_time, max(blk_lvl) blk_lvl
              from register_node_transactions
              group by reg_entity_id, reg_id, reg_consensus_address
              ) nodes
              ANY
              LEFT JOIN (
                --Block proposed count
                select blk_proposer_address reg_consensus_address, max(blk_created_at) last_block_time, count() c_blocks
                         from blocks
                         group by blk_proposer_address) blk USING reg_consensus_address
       ) prep
       ANY
       LEFT JOIN (
         --Blocks signatures count
         select sig_validator_address reg_consensus_address, max(sig_timestamp) last_signature_time, count() c_block_signatures
                  from block_signatures
                  group by sig_validator_address) blk USING reg_consensus_address;

DROP TABLE public_validators;
DROP TABLE validators_list_view;
