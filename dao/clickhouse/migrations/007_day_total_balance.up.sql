CREATE VIEW account_day_balance_view AS
select acb_account,
       toStartOfDay(blk_time) start_of_period,
       anyLast(acb_general_balance) acb_general_balance,
       anyLast(acb_escrow_balance_active)   escrow_balance_active,
       anyLast(acb_escrow_debonding_active) escrow_debonding_active
from account_balance
group by acb_account, start_of_period

CREATE VIEW IF NOT EXISTS day_total_balance_view AS
select start_of_period,
       sum(acb_general_balance)         general_balance,
       sum(escrow_balance_active)   escrow_balance_active,
       sum(escrow_debonding_active) escrow_debonding_active
from account_day_balance_view
GROUP BY start_of_period
ORDER BY start_of_period desc;
