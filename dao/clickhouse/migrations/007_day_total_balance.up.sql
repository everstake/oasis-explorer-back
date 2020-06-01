CREATE VIEW IF NOT EXISTS day_total_balance_view AS
select start_of_period,
       sum(toUInt64(acb_general_balance)) general_balance,
       sum(toUInt64(acb_escrow_balance_active)) escrow_balance_active,
       sum(toUInt64(acb_escrow_debonding_active)) escrow_debonding_active
from (select acb_account,start_of_period, max(blk_lvl) blk_lvl from account_balance group by acb_account, toStartOfDay(blk_time) as start_of_period) s
       ANY
       LEFT JOIN account_balance USING acb_account, blk_lvl
       GROUP BY start_of_period;
