DELETE FROM crawler_sources WHERE source_type IN (
  'html_ibps', 'html_sbi', 'html_gds', 'html_nta', 'html_uppsc', 'html_upsssc',
  'html_bpsc', 'html_mppsc', 'html_mpsc', 'html_rpsc', 'html_rsmssb', 'html_dsssb',
  'html_hssc', 'html_kpsc', 'html_tnpsc', 'html_wbpsc', 'html_army', 'html_afcat',
  'html_navy', 'html_coastguard', 'html_isro', 'html_drdo', 'html_barc', 'html_ongc',
  'html_ntpc', 'html_bhel', 'html_iocl', 'api_greenhouse', 'api_lever', 'api_ashby',
  'api_smartrecruiters', 'api_workable', 'api_recruitee', 'api_remoteok',
  'rss_weworkremotely', 'api_hackernews'
);
