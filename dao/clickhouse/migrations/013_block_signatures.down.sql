DROP TABLE block_signatures_count_view;
DROP TABLE block_signatures_count_mv;

CREATE VIEW IF NOT EXISTS blocks_sig_count AS select blk_lvl, count() sig_count from block_signatures group by blk_lvl;