-- Multi-client login and Xiaoyuanzi family seed data.

ALTER TABLE `users`
  ADD COLUMN `username` VARCHAR(50) DEFAULT NULL COMMENT '用户名' AFTER `phone`,
  ADD UNIQUE KEY `uk_username` (`username`);

UPDATE `users`
SET `username` = `phone`
WHERE `username` IS NULL OR `username` = '';

INSERT INTO `users` (`id`, `phone`, `username`, `nickname`, `is_admin`, `password`)
VALUES
  ('100e8400-e29b-41d4-a716-446655440000', '13800138000', 'mom', '妈妈', 1, 'yuanzi123'),
  ('100e8400-e29b-41d4-a716-446655440001', '13800138001', 'dad', '爸爸', 0, 'yuanzi123'),
  ('100e8400-e29b-41d4-a716-446655440002', '13800138002', 'grandpa', '爷爷', 0, 'yuanzi123'),
  ('100e8400-e29b-41d4-a716-446655440003', '13800138003', 'grandma', '奶奶', 0, 'yuanzi123'),
  ('100e8400-e29b-41d4-a716-446655440004', '13800138004', 'waigong', '外公', 0, 'yuanzi123'),
  ('100e8400-e29b-41d4-a716-446655440005', '13800138005', 'waipo', '外婆', 0, 'yuanzi123')
ON DUPLICATE KEY UPDATE
  `username` = VALUES(`username`),
  `nickname` = VALUES(`nickname`),
  `is_admin` = VALUES(`is_admin`),
  `password` = VALUES(`password`);

INSERT INTO `families` (`id`, `name`, `invite_code`, `created_by`, `is_paid`, `storage_limit`)
VALUES ('200e8400-e29b-41d4-a716-446655440000', '小园子的家', 'ABC12345', '100e8400-e29b-41d4-a716-446655440000', 0, 1073741824)
ON DUPLICATE KEY UPDATE `name` = VALUES(`name`), `created_by` = VALUES(`created_by`);

INSERT INTO `family_members` (`id`, `family_id`, `user_id`, `role`, `elder_mode`, `notifications`)
VALUES
  ('300e8400-e29b-41d4-a716-446655440000', '200e8400-e29b-41d4-a716-446655440000', '100e8400-e29b-41d4-a716-446655440000', 'admin', 0, JSON_OBJECT('feed', true, 'sleep', true)),
  ('300e8400-e29b-41d4-a716-446655440001', '200e8400-e29b-41d4-a716-446655440000', '100e8400-e29b-41d4-a716-446655440001', 'member', 0, JSON_OBJECT('feed', true, 'sleep', true)),
  ('300e8400-e29b-41d4-a716-446655440002', '200e8400-e29b-41d4-a716-446655440000', '100e8400-e29b-41d4-a716-446655440002', 'elder', 1, JSON_OBJECT('feed', true, 'sleep', true)),
  ('300e8400-e29b-41d4-a716-446655440003', '200e8400-e29b-41d4-a716-446655440000', '100e8400-e29b-41d4-a716-446655440003', 'elder', 1, JSON_OBJECT('feed', true, 'sleep', true)),
  ('300e8400-e29b-41d4-a716-446655440004', '200e8400-e29b-41d4-a716-446655440000', '100e8400-e29b-41d4-a716-446655440004', 'elder', 1, JSON_OBJECT('feed', true, 'sleep', true)),
  ('300e8400-e29b-41d4-a716-446655440005', '200e8400-e29b-41d4-a716-446655440000', '100e8400-e29b-41d4-a716-446655440005', 'elder', 1, JSON_OBJECT('feed', true, 'sleep', true))
ON DUPLICATE KEY UPDATE `role` = VALUES(`role`), `elder_mode` = VALUES(`elder_mode`);

INSERT INTO `babies` (`id`, `family_id`, `name`, `birthday`, `gender`, `birth_weight`, `birth_height`)
VALUES ('400e8400-e29b-41d4-a716-446655440000', '200e8400-e29b-41d4-a716-446655440000', '小园子', '2024-01-01', 2, 3.20, 50.0)
ON DUPLICATE KEY UPDATE `name` = VALUES(`name`);

INSERT INTO `records` (`id`, `baby_id`, `family_id`, `type`, `started_at`, `ended_at`, `content`, `note`, `created_by`)
VALUES
  ('500e8400-e29b-41d4-a716-446655440005', '400e8400-e29b-41d4-a716-446655440000', '200e8400-e29b-41d4-a716-446655440000', 'excretion', NOW(3) - INTERVAL 30 MINUTE, NULL, JSON_OBJECT('type','poop','color','黄色','consistency','糊状','amount','normal'), '排便正常', '100e8400-e29b-41d4-a716-446655440001'),
  ('500e8400-e29b-41d4-a716-446655440006', '400e8400-e29b-41d4-a716-446655440000', '200e8400-e29b-41d4-a716-446655440000', 'feeding', NOW(3) - INTERVAL 1 DAY, NULL, JSON_OBJECT('type','left_breast','side','left','duration',18,'unit','minute'), '左侧母乳', '100e8400-e29b-41d4-a716-446655440000'),
  ('500e8400-e29b-41d4-a716-446655440007', '400e8400-e29b-41d4-a716-446655440000', '200e8400-e29b-41d4-a716-446655440000', 'feeding', NOW(3) - INTERVAL 2 DAY, NULL, JSON_OBJECT('type','right_breast','side','right','duration',16,'unit','minute'), '右侧母乳', '100e8400-e29b-41d4-a716-446655440000')
ON DUPLICATE KEY UPDATE `note` = VALUES(`note`);

INSERT INTO `ai_chat_records` (`id`, `user_id`, `baby_id`, `question`, `answer`, `tokens_used`, `model`, `created_at`)
VALUES
  ('700e8400-e29b-41d4-a716-446655440001', '100e8400-e29b-41d4-a716-446655440000', '400e8400-e29b-41d4-a716-446655440000', '近一周小园子睡眠和喝奶趋势怎么样？', '近一周睡眠整体平稳，奶量记录充足，排泄没有异常信号。', 96, 'seed-model', NOW(3) - INTERVAL 1 DAY),
  ('700e8400-e29b-41d4-a716-446655440002', '100e8400-e29b-41d4-a716-446655440001', '400e8400-e29b-41d4-a716-446655440000', '今天排泄需要注意什么？', '目前记录看排泄正常，继续观察颜色、次数和精神状态。', 72, 'seed-model', NOW(3) - INTERVAL 2 HOUR)
ON DUPLICATE KEY UPDATE `answer` = VALUES(`answer`);
