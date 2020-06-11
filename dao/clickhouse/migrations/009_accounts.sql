CREATE VIEW IF NOT EXISTS oasis.account_balance_view AS
select *
from (select acb_account, min(blk_time) created_on , max(blk_lvl) blk_lvl
      from oasis.account_balance
      group by acb_account
       ) s
       ANY
       LEFT JOIN oasis.account_balance USING acb_account, blk_lvl;

--Last registered unique entity
CREATE VIEW IF NOT EXISTS oasis.entity_register_view AS
select *
from (
      select reg_entity_id, max(blk_lvl) blk_lvl
      from oasis.register_entity_transactions
      group by reg_entity_id) m
 ANY
       LEFT JOIN oasis.register_entity_transactions USING reg_entity_id, blk_lvl;

--Last registered unique node
CREATE VIEW IF NOT EXISTS oasis.node_registry_view AS
select *
from (
      select reg_id, max(blk_lvl) blk_lvl
      from oasis.register_node_transactions
      group by reg_id) m
 ANY
       LEFT JOIN oasis.register_node_transactions USING reg_id, blk_lvl;


---
CREATE MATERIALIZED VIEW IF NOT EXISTS account_balance_merge_mv
ENGINE = AggregatingMergeTree() PARTITION BY toYYYYMM(created_at) ORDER BY (acb_account) AS
SELECT
    acb_account,
    min(blk_time) created_at,
    anyLastState(acb_general_balance) acb_general_balance,
    anyLastState(acb_escrow_balance_active) acb_escrow_balance_active,
    anyLastState(acb_escrow_balance_share) acb_escrow_balance_share
FROM account_balance
GROUP BY acb_account;

CREATE VIEW IF NOT EXISTS account_last_balance_view AS
SELECT
    acb_account,
    min(created_at) created_at,
    anyLastMerge(acb_general_balance) acb_general_balance,
    anyLastMerge(acb_escrow_balance_active) acb_escrow_balance_active,
    anyLastMerge(acb_escrow_balance_share) acb_escrow_balance_share
FROM account_balance_merge_mv
GROUP BY acb_account;

--
CREATE VIEW IF NOT EXISTS account_list_view AS
select acb_account, created_at, acb_general_balance, acb_escrow_balance_active, tx_escrow_account delegate, entity.blk_lvl entity, prp.blk_lvl node from (
select *
from (select *
      from account_last_balance_view
             ANY
             LEFT JOIN (SELECT tx_sender acb_account, * from entity_depositors_view where balance > 0) al
                       USING acb_account
       ) s
       ANY
       LEFT JOIN (select reg_id acb_account, CASE WHEN max(blk_lvl) = 0 THEN 1 ELSE max(blk_lvl) END as blk_lvl from register_node_transactions group by acb_account) node USING acb_account) prp
ANY
LEFT JOIN (select reg_entity_id acb_account, CASE WHEN max(blk_lvl) = 0 THEN 1 ELSE max(blk_lvl) END as blk_lvl from register_entity_transactions group by acb_account) entity USING acb_account;