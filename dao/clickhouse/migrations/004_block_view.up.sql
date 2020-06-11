CREATE VIEW IF NOT EXISTS block_row_view AS
select * from blocks ANY LEFT JOIN
              (select blk_lvl, sum(tx_gas_price), sum(tx_fee), count()
from transactions group by blk_lvl) as s USING blk_lvl ORDER BY blk_lvl DESC;

CREATE VIEW IF NOT EXISTS blocks_sig_count AS select blk_lvl, count() sig_count from block_signatures group by blk_lvl;