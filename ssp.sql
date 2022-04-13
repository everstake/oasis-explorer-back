select *
from (
         select *, CASE WHEN last_signature_time >= (select now() - INTERVAL 3 HOUR) THEN 1 ELSE 0 END is_active
from (
    SELECT *
    FROM (select p.*, b.blk_lvl max_day_block, b.blk_count day_blocks
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
    max(last_signature_time)       last_signature_time,
    max(last_block_time)           last_block_time,
    sum(blocks)                    blocks,
    sum(signed_blocks)             signed_blocks,
    sum(signatures)                signatures
    from (select *
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
    LEFT JOIN validator_block_signatures_count_view USING reg_consensus_address
    order by last_signature_time asc ) alle
    GROUP BY reg_entity_address) g
    any
    left join (
    select reg_consensus_address, day_signatures, day_signed_blocks
    from validator_day_stats_view
    where day = toStartOfDay(now())) sigs USING reg_consensus_address) p
    ANY
    LEFT JOIN day_max_block_lvl_view b USING day) lele -- created MVIEW to don't make fullscan of txs
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
    LEFT JOIN (SELECT tx_receiver AS reg_entity_address,
    count()     AS depositors_num
    FROM (select *, add.input - remove.output balance from (
    select tx_receiver, tx_sender, min(tx_time) escrow_since, sum(tx_escrow_amount) input
    from transactions
    where tx_type = 'addescrow' and tx_status = 1
    group by  tx_receiver, tx_sender) add
    ANY
    LEFT JOIN (
    select tx_receiver, tx_sender, sum(tx_escrow_reclaim_amount) output
    from transactions
    where tx_type = 'reclaimescrow' and tx_status = 1
    group by tx_receiver, tx_sender) remove USING tx_receiver, tx_sender)
    WHERE balance > 0
    GROUP BY tx_receiver) abba -- fixed to make aggregate from MVIEW
    USING reg_entity_address -- entity_active_depositors_counter_view  is aggregate upon `entity_depositors_view`, a very huge select with joins.
    ) b USING reg_entity_address
where blocks > 0
   OR signatures > 0
   OR reg_expiration >= (select max(blk_epoch) from blocks)) prep
    ANY
    LEFT JOIN public_validators USING reg_entity_address;


--70616789F5D173B678560DED991D8D9C55C3C666 53B4C392C33E9B69FF9D4E973E203DB980A2AA83 4024D5C2B90DAC79F9AC75A9E5E8E6058BCD3D04 104D0EB71B5067DF92D8380EA4E5341FF9D3734B
--oasis1qr0jwz65c29l044a204e3cllvumdg8cmsgt2k3ql
select reg_consensus_address, day_signatures, day_signed_blocks
from validator_day_stats_view
where day = toStartOfDay(now()) AND reg_consensus_address='104D0EB71B5067DF92D8380EA4E5341FF9D3734B';

select reg_entity_id,
       reg_entity_address,
       reg_id,
       reg_address,
       reg_consensus_address,
       min(tx_time)        created_time,
       max(blk_lvl)        blk_lvl,
       max(reg_expiration) reg_expiration
from register_node_transactions
group by reg_entity_id, reg_entity_address, reg_address, reg_id, reg_consensus_address;





select p.*, b.blk_lvl max_day_block, b.blk_count day_blocks
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
                         max(last_signature_time)       last_signature_time,
                         sum(blocks)                    blocks,
                         sum(signed_blocks)             signed_blocks,
                         sum(signatures)                signatures
                  from (select *
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
                                 LEFT JOIN validator_block_signatures_count_view USING reg_consensus_address
                        order by last_signature_time asc ) alle
                  GROUP BY reg_entity_address) g
             any
                  left join (
             select reg_consensus_address, day_signatures, day_signed_blocks
             from validator_day_stats_view
             where day = toStartOfDay(now())) sigs USING reg_consensus_address) p
    ANY
         LEFT JOIN day_max_block_lvl_view b USING day;

--entity_nodes_view
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
         LEFT JOIN validator_block_signatures_count_view USING reg_consensus_address
order by last_signature_time desc;


select * from validator_blocks_count_view







----------------------------------------------------------- change to this on server entity_depositors_view
select *, add.input - remove.output balance from (
                                                     select tx_receiver, tx_sender, min(tx_time) escrow_since, sum(tx_escrow_amount) input
                                                     from transactions
                                                     where tx_type = 'addescrow' and tx_status = 1
                                                     group by  tx_receiver, tx_sender) add
     ANY
     LEFT JOIN (
    select tx_receiver, tx_sender, sum(tx_escrow_reclaim_amount) output
    from transactions
    where tx_type = 'reclaimescrow' and tx_status = 1
    group by tx_receiver, tx_sender) remove USING tx_receiver, tx_sender;


