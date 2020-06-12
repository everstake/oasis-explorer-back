CREATE VIEW IF NOT EXISTS entity_depositors_view AS
select *, add.input - remove.output balance from (
  select tx_escrow_account, tx_sender, sum(tx_escrow_amount) input
  from transactions
  where tx_type = 'addescrow'
  group by  tx_escrow_account, tx_sender) add
   ANY
       LEFT JOIN (
         select tx_escrow_account, tx_sender, sum(tx_escrow_reclaim_amount) output
         from transactions
         where tx_type = 'reclaimescrow'
         group by tx_escrow_account, tx_sender) remove USING tx_escrow_account, tx_sender;

CREATE VIEW IF NOT EXISTS entity_active_depositors_counter_view AS
  SELECT tx_escrow_account reg_entity_id, count() depositors_num from entity_depositors_view
  where balance > 0
  group by tx_escrow_account;

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
