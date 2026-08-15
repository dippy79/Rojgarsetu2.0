-- Insert complete registry of Central Govt, State PSCs, Defense, PSUs, Open ATS & Remote APIs
INSERT INTO crawler_sources (name, category, source_type, base_url, robots_txt_url) VALUES
-- Central Govt & Banking
('IBPS Portal', 'GOVT_JOB', 'html_ibps', 'https://www.ibps.in', 'https://www.ibps.in/robots.txt'),
('SBI Careers', 'GOVT_JOB', 'html_sbi', 'https://sbi.co.in/web/careers', NULL),
('India Post GDS', 'GOVT_JOB', 'html_gds', 'https://indiapostgdsonline.gov.in', NULL),
('NTA Testing Agency', 'GOVT_JOB', 'html_nta', 'https://nta.ac.in', NULL),

-- State Public Service Commissions & Subordinate Boards
('UPPSC (Uttar Pradesh)', 'GOVT_JOB', 'html_uppsc', 'https://uppsc.up.nic.in', NULL),
('UPSSSC (UP Subordinate)', 'GOVT_JOB', 'html_upsssc', 'https://upsssc.gov.in', NULL),
('BPSC (Bihar)', 'GOVT_JOB', 'html_bpsc', 'https://bpsc.bih.nic.in', NULL),
('MPPSC (Madhya Pradesh)', 'GOVT_JOB', 'html_mppsc', 'https://mppsc.mp.gov.in', NULL),
('MPSC (Maharashtra)', 'GOVT_JOB', 'html_mpsc', 'https://mpsc.gov.in', NULL),
('RPSC (Rajasthan)', 'GOVT_JOB', 'html_rpsc', 'https://rpsc.rajasthan.gov.in', NULL),
('RSMSSB (Rajasthan Staff)', 'GOVT_JOB', 'html_rsmssb', 'https://rsmssb.rajasthan.gov.in', NULL),
('DSSSB (Delhi)', 'GOVT_JOB', 'html_dsssb', 'https://dsssb.delhi.gov.in', NULL),
('HSSC (Haryana)', 'GOVT_JOB', 'html_hssc', 'https://hssc.gov.in', NULL),
('KPSC (Karnataka)', 'GOVT_JOB', 'html_kpsc', 'https://kpsc.kar.nic.in', NULL),
('TNPSC (Tamil Nadu)', 'GOVT_JOB', 'html_tnpsc', 'https://tnpsc.gov.in', NULL),
('WBPSC (West Bengal)', 'GOVT_JOB', 'html_wbpsc', 'https://psc.wb.gov.in', NULL),

-- Defense & Security Forces
('Join Indian Army', 'GOVT_JOB', 'html_army', 'https://joinindianarmy.nic.in', NULL),
('Indian Air Force (AFCAT)', 'GOVT_JOB', 'html_afcat', 'https://afcat.cdac.in', NULL),
('Indian Navy', 'GOVT_JOB', 'html_navy', 'https://joinindiannavy.gov.in', NULL),
('Indian Coast Guard', 'GOVT_JOB', 'html_coastguard', 'https://joinindiancoastguard.cdac.in', NULL),

-- Scientific Research & PSUs
('ISRO Careers', 'GOVT_JOB', 'html_isro', 'https://www.isro.gov.in/Careers.html', NULL),
('DRDO RAC', 'GOVT_JOB', 'html_drdo', 'https://drdo.gov.in/careers', NULL),
('BARC Careers', 'GOVT_JOB', 'html_barc', 'https://barc.gov.in/careers', NULL),
('ONGC India', 'GOVT_JOB', 'html_ongc', 'https://ongcindia.com/careers', NULL),
('NTPC Careers', 'GOVT_JOB', 'html_ntpc', 'https://careers.ntpc.co.in', NULL),
('BHEL Trainee', 'GOVT_JOB', 'html_bhel', 'https://bhel.com/careers', NULL),
('IOCL Apprentice', 'GOVT_JOB', 'html_iocl', 'https://iocl.com/apprenticeships', NULL),

-- Open ATS Public REST APIs (Direct Corporate Hiring)
('Greenhouse ATS', 'PRIVATE_JOB', 'api_greenhouse', 'https://boards-api.greenhouse.io', NULL),
('Lever ATS', 'PRIVATE_JOB', 'api_lever', 'https://api.lever.co', NULL),
('Ashby ATS', 'PRIVATE_JOB', 'api_ashby', 'https://api.ashbyhq.com', NULL),
('SmartRecruiters ATS', 'PRIVATE_JOB', 'api_smartrecruiters', 'https://api.smartrecruiters.com', NULL),
('Workable ATS', 'PRIVATE_JOB', 'api_workable', 'https://apply.workable.com', NULL),
('Recruitee ATS', 'PRIVATE_JOB', 'api_recruitee', 'https://recruitee.com', NULL),

-- Remote & Developer Open Feeds
('RemoteOK API', 'PRIVATE_JOB', 'api_remoteok', 'https://remoteok.com/api', NULL),
('WeWorkRemotely RSS', 'PRIVATE_JOB', 'rss_weworkremotely', 'https://weworkremotely.com/remote-jobs.rss', NULL),
('HackerNews Jobs API', 'PRIVATE_JOB', 'api_hackernews', 'https://hacker-news.firebaseio.com/v0', NULL)
ON CONFLICT DO NOTHING;
