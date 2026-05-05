SET NAMES utf8mb4;
USE `music_box`;

-- ==========================================
-- 1. 插入演示 / 测试数据
-- ==========================================

-- 1.1 插入 4 个测试用户
INSERT INTO `users` (`id`, `username`, `password_hash`, `avatar_url`, `role`) VALUES
(2, 'testuser', '$2a$10$emZD53.StcoycaaBcIXGg.CcPKQMUl0d5FCHCVznV8xMq.ikHPNJe', 'uploads/avatars/2.jpg', 'User'),
(3, 'alice', '$2a$10$emZD53.StcoycaaBcIXGg.CcPKQMUl0d5FCHCVznV8xMq.ikHPNJe', 'uploads/avatars/3.jpg', 'User'),
(4, 'bob', '$2a$10$emZD53.StcoycaaBcIXGg.CcPKQMUl0d5FCHCVznV8xMq.ikHPNJe', 'uploads/avatars/4.jpg', 'User'),
(5, 'charlie', '$2a$10$emZD53.StcoycaaBcIXGg.CcPKQMUl0d5FCHCVznV8xMq.ikHPNJe', 'uploads/avatars/5.jpg', 'User');

-- 1.2 插入 30 首歌曲
INSERT INTO `songs` (`id`, `title`, `artist`, `album`, `duration`, `file_url`, `cover_url`, `lyric_url`) VALUES
(1, 'Just Say Hello', 'melo-D', 'Just Say Hello', 215, 'uploads/music/Just Say Hello-melo-D.mp3', 'uploads/covers/Just Say Hello-melo-D.jpg', 'uploads/lyrics/Just Say Hello-melo-D.lrc'),
(2, 'Panoramic View', 'Awolnation', 'Megalithic Symphony', 240, 'uploads/music/Panoramic View-Awolnation.mp3', 'uploads/covers/Panoramic View-Awolnation.jpg', 'uploads/lyrics/Panoramic View-Awolnation.lrc'),
(3, 'Let go', 'Nancy Kwai', 'Let go', 198, 'uploads/music/Let go-Nancy Kwai.flac', 'uploads/covers/Let go-Nancy Kwai.jpg', 'uploads/lyrics/Let go-Nancy Kwai.lrc'),
(4, 'Will you answer me', 'Nancy Kwai', 'Will you answer me', 205, 'uploads/music/Will you answer me-Nancy Kwai.flac', 'uploads/covers/Will you answer me-Nancy Kwai.jpg', 'uploads/lyrics/Will you answer me-Nancy Kwai.lrc'),
(5, '一路向北', '周杰伦', '11月的萧邦', 295, 'uploads/music/一路向北-周杰伦.mp3', 'uploads/covers/一路向北-周杰伦.jpg', 'uploads/lyrics/一路向北-周杰伦.lrc'),
(6, '交换余生', '林俊杰', '幸存者 • 如你', 276, 'uploads/music/交换余生-林俊杰.mp3', 'uploads/covers/交换余生-林俊杰.png', 'uploads/lyrics/交换余生-林俊杰.lrc'),
(7, '倒数', '邓紫棋', '另一个童话', 229, 'uploads/music/倒数-邓紫棋.mp3', 'uploads/covers/倒数-邓紫棋.png', 'uploads/lyrics/倒数-邓紫棋.lrc'),
(8, '光年之外', 'G.E.M.邓紫棋', 'Passengers OST', 235, 'uploads/music/光年之外-G.E.M.邓紫棋.mp3', 'uploads/covers/光年之外-G.E.M.邓紫棋.jpg', 'uploads/lyrics/光年之外-G.E.M.邓紫棋.lrc'),
(9, '刚刚好', '薛之谦', '初学者', 251, 'uploads/music/刚刚好-薛之谦.mp3', 'uploads/covers/刚刚好-薛之谦.jpg', 'uploads/lyrics/刚刚好-薛之谦.lrc'),
(10, '动物世界', '薛之谦', '渡 The Crossing', 230, 'uploads/music/动物世界-薛之谦.mp3', 'uploads/covers/动物世界-薛之谦.jpg', 'uploads/lyrics/动物世界-薛之谦.lrc'),
(11, '告白气球', '周杰伦', '周杰伦的床边故事', 215, 'uploads/music/告白气球-周杰伦.mp3', 'uploads/covers/告白气球-周杰伦.png', 'uploads/lyrics/告白气球-周杰伦.lrc'),
(12, '唯一', 'G.E.M.邓紫棋', '摩天动物园', 236, 'uploads/music/唯一-G.E.M.邓紫棋.mp3', 'uploads/covers/唯一-G.E.M.邓紫棋.jpg', 'uploads/lyrics/唯一-G.E.M.邓紫棋.lrc'),
(13, '天外来物', '薛之谦', '天外来物', 257, 'uploads/music/天外来物-薛之谦.mp3', 'uploads/covers/天外来物-薛之谦.jpg', 'uploads/lyrics/天外来物-薛之谦.lrc'),
(14, '富士山下', '陈奕迅', 'What''s Going On...?', 259, 'uploads/music/富士山下-陈奕迅.mp3', 'uploads/covers/富士山下-陈奕迅.png', 'uploads/lyrics/富士山下-陈奕迅.lrc'),
(15, '小半', '陈粒', '小梦大半', 297, 'uploads/music/小半-陈粒.mp3', 'uploads/covers/小半-陈粒.jpg', 'uploads/lyrics/小半-陈粒.lrc'),
(16, '怪咖', '薛之谦', '怪咖', 241, 'uploads/music/怪咖-薛之谦.mp3', 'uploads/covers/怪咖-薛之谦.jpg', 'uploads/lyrics/怪咖-薛之谦.lrc'),
(17, '我好像在哪见过你', '薛之谦', '初学者', 279, 'uploads/music/我好像在哪见过你-薛之谦.mp3', 'uploads/covers/我好像在哪见过你-薛之谦.jpg', 'uploads/lyrics/我好像在哪见过你-薛之谦.lrc'),
(18, '打上花火', '米津玄師', 'BOOTLEG', 289, 'uploads/music/打上花火-米津玄師.mp3', 'uploads/covers/打上花火-米津玄師.jpg', 'uploads/lyrics/打上花火-米津玄師.lrc'),
(19, '搁浅', '周杰伦', '七里香', 239, 'uploads/music/搁浅-周杰伦.mp3', 'uploads/covers/搁浅-周杰伦.jpg', 'uploads/lyrics/搁浅-周杰伦.lrc'),
(20, '无名的人', '毛不易', '雄狮少年 电影原声带', 276, 'uploads/music/无名的人-毛不易.mp3', 'uploads/covers/无名的人-毛不易.jpg', 'uploads/lyrics/无名的人-毛不易.lrc'),
(21, '泡沫', '邓紫棋', 'Xposed', 258, 'uploads/music/泡沫-邓紫棋.mp3', 'uploads/covers/泡沫-邓紫棋.png', 'uploads/lyrics/泡沫-邓紫棋.lrc'),
(22, '消愁(Live)', '毛不易', '明日之子', 178, 'uploads/music/消愁(Live)-毛不易.mp3', 'uploads/covers/消愁(Live)-毛不易.jpg', 'uploads/lyrics/消愁(Live)-毛不易.lrc'),
(23, '演员', '薛之谦', '初学者', 261, 'uploads/music/演员-薛之谦.mp3', 'uploads/covers/演员-薛之谦.jpg', 'uploads/lyrics/演员-薛之谦.lrc'),
(24, '烟花易冷', '周杰伦', '跨时代', 313, 'uploads/music/烟花易冷-周杰伦.mp3', 'uploads/covers/烟花易冷-周杰伦.jpg', 'uploads/lyrics/烟花易冷-周杰伦.lrc'),
(25, '稻香', '周杰伦', '魔杰座', 223, 'uploads/music/稻香-周杰伦.mp3', 'uploads/covers/稻香-周杰伦.jpg', 'uploads/lyrics/稻香-周杰伦.lrc'),
(26, '红尘客栈', '周杰伦', '12新作', 274, 'uploads/music/红尘客栈-周杰伦.mp3', 'uploads/covers/红尘客栈-周杰伦.jpg', 'uploads/lyrics/红尘客栈-周杰伦.lrc'),
(27, '红色高跟鞋', '蔡健雅', '若你碰到他', 206, 'uploads/music/红色高跟鞋-蔡健雅.mp3', 'uploads/covers/红色高跟鞋-蔡健雅.jpg', 'uploads/lyrics/红色高跟鞋-蔡健雅.lrc'),
(28, '起风了', '买辣椒也用券', '起风了', 308, 'uploads/music/起风了-买辣椒也用券.mp3', 'uploads/covers/起风了-买辣椒也用券.jpg', 'uploads/lyrics/起风了-买辣椒也用券.lrc'),
(29, '跳楼机', 'LBI利比 (时伯尘)', '跳楼机', 186, 'uploads/music/跳楼机-LBI利比 (时伯尘).mp3', 'uploads/covers/跳楼机-LBI利比 (时伯尘).png', 'uploads/lyrics/跳楼机-LBI利比 (时伯尘).lrc'),
(30, '青花瓷', '周杰伦', '我很忙', 239, 'uploads/music/青花瓷-周杰伦.mp3', 'uploads/covers/青花瓷-周杰伦.jpg', 'uploads/lyrics/青花瓷-周杰伦.lrc');

-- 1.3 插入 30 个歌单
INSERT INTO `playlists` (`id`, `user_id`, `title`, `description`, `cover_url`, `is_public`, `play_count`) VALUES
(1, 1, '薛氏情歌精选', '薛之谦经典曲目，深夜 emo 必备。', 'uploads/covers/演员-薛之谦.jpg', 1, 1280),
(2, 1, '周杰伦的青春', '那些年我们一起听过的杰伦。', 'uploads/covers/稻香-周杰伦.jpg', 1, 3450),
(3, 1, '华语流行热歌', '近期最火的华语流行榜单。', 'uploads/covers/起风了-买辣椒也用券.jpg', 1, 890),
(4, 1, '深夜emo专属', '网抑云时间到了。', 'uploads/covers/消愁(Live)-毛不易.jpg', 1, 2100),
(5, 1, '开车防困指南', '节奏感拉满，拒绝疲劳驾驶。', 'uploads/covers/Panoramic View-Awolnation.jpg', 1, 560),
(6, 1, '经典老歌回顾', '老歌越听越有味道。', 'uploads/covers/富士山下-陈奕迅.png', 1, 430),
(7, 2, '邓紫棋高音赏析', '铁肺天后的极致震撼。', 'uploads/covers/光年之外-G.E.M.邓紫棋.jpg', 1, 1500),
(8, 2, '米津玄师特辑', '八爷的日系神曲合集。', 'uploads/covers/打上花火-米津玄師.jpg', 1, 2300),
(9, 2, '欧美抖腿向', '听了就停不下来的欧美节奏。', 'uploads/covers/Just Say Hello-melo-D.jpg', 1, 780),
(10, 2, '沉浸式学习BGM', '适合看书写码时听的轻柔人声。', 'uploads/covers/Let go-Nancy Kwai.jpg', 1, 3200),
(11, 2, '健身房燃脂', 'BPM120+，跑起来！', 'uploads/covers/倒数-邓紫棋.png', 1, 940),
(12, 2, '华语金曲100首', 'KTV必点，首首大合唱。', 'uploads/covers/红色高跟鞋-蔡健雅.jpg', 1, 4100),
(13, 3, '独立女声精选', '那些独特而慵懒的女声。', 'uploads/covers/小半-陈粒.jpg', 1, 670),
(14, 3, '民谣与诗', '一把吉他，一个故事。', 'uploads/covers/无名的人-毛不易.jpg', 1, 1100),
(15, 3, '周末慵懒时光', '周末宅家必备背景音。', 'uploads/covers/唯一-G.E.M.邓紫棋.jpg', 1, 850),
(16, 3, '午后咖啡馆', '充满咖啡香气的音乐。', 'uploads/covers/告白气球-周杰伦.png', 1, 420),
(17, 3, 'KTV必点神曲', '麦霸养成计划。', 'uploads/covers/青花瓷-周杰伦.jpg', 1, 5200),
(18, 3, '治愈系轻音乐', '放松身心，洗涤灵魂。', 'uploads/covers/Will you answer me-Nancy Kwai.jpg', 1, 990),
(19, 4, '失恋阵线联盟', '哭出来就好了。', 'uploads/covers/怪咖-薛之谦.jpg', 1, 1800),
(20, 4, '怀旧金曲大赏', '时代的眼泪。', 'uploads/covers/烟花易冷-周杰伦.jpg', 1, 2150),
(21, 4, '抖音爆款热歌', '你一定在刷短视频时听过。', 'uploads/covers/跳楼机-LBI利比 (时伯尘).png', 1, 6600),
(22, 4, '单曲循环预警', '好听到出不去的歌。', 'uploads/covers/交换余生-林俊杰.png', 1, 1340),
(23, 4, '华语男声实力派', '行走的CD们。', 'uploads/covers/动物世界-薛之谦.jpg', 1, 880),
(24, 4, '雨天听歌指南', '下雨天和音乐更配。', 'uploads/covers/红尘客栈-周杰伦.jpg', 1, 750),
(25, 5, '纯粹人声清唱', '感受嗓音的魅力。', 'uploads/covers/刚刚好-薛之谦.jpg', 1, 310),
(26, 5, '随机音乐日记', '每天一首不重样。', 'uploads/covers/泡沫-邓紫棋.png', 1, 460),
(27, 5, '华语R&B启蒙', '回到那个神仙打架的年代。', 'uploads/covers/稻香-周杰伦.jpg', 1, 1250),
(28, 5, '少数派宝藏', '评论不过999的好歌。', 'uploads/covers/我好像在哪见过你-薛之谦.jpg', 1, 890),
(29, 5, '听海的落日', '海风吹过的声音。', 'uploads/covers/搁浅-周杰伦.jpg', 1, 1050),
(30, 5, '动感节奏控', '跟着节奏摇摆。', 'uploads/covers/一路向北-周杰伦.jpg', 1, 2200);

-- 1.4 填充歌单与歌曲的关联关系
INSERT INTO `playlist_songs` (`playlist_id`, `song_id`, `sort_order`) VALUES
(1, 9, 1), (1, 10, 2), (1, 13, 3), (1, 16, 4), (1, 17, 5), (1, 23, 6),
(2, 5, 1), (2, 11, 2), (2, 19, 3), (2, 24, 4), (2, 25, 5), (2, 26, 6), (2, 30, 7),
(3, 28, 1), (3, 6, 2), (3, 7, 3), (3, 20, 4),
(4, 22, 1), (4, 14, 2), (4, 15, 3), (4, 21, 4), (4, 19, 5),
(5, 2, 1), (5, 1, 2), (5, 18, 3), (5, 29, 4),
(6, 14, 1), (6, 27, 2), (6, 24, 3), (6, 30, 4),
(7, 7, 1), (7, 8, 2), (7, 12, 3), (7, 21, 4),
(8, 18, 1),
(9, 1, 1), (9, 2, 2), (9, 3, 3), (9, 4, 4),
(10, 3, 1), (10, 4, 2), (10, 15, 3), (10, 28, 4),
(11, 7, 1), (11, 2, 2), (11, 29, 3), (11, 8, 4),
(12, 5, 1), (12, 11, 2), (12, 27, 3), (12, 14, 4), (12, 23, 5),
(13, 15, 1), (13, 27, 2), (13, 28, 3),
(14, 20, 1), (14, 22, 2), (14, 15, 3),
(15, 3, 1), (15, 4, 2), (15, 11, 3), (15, 25, 4),
(16, 27, 1), (16, 11, 2), (16, 15, 3),
(17, 23, 1), (17, 5, 2), (17, 7, 3), (17, 21, 4), (17, 8, 5),
(18, 4, 1), (18, 3, 2), (18, 25, 3),
(19, 19, 1), (19, 21, 2), (19, 23, 3), (19, 16, 4), (19, 14, 5),
(20, 24, 1), (20, 30, 2), (20, 26, 3),
(21, 29, 1), (21, 28, 2), (21, 18, 3), (21, 8, 4),
(22, 6, 1), (22, 13, 2), (22, 20, 3), (22, 5, 4),
(23, 6, 1), (23, 14, 2), (23, 20, 3), (23, 13, 4),
(24, 24, 1), (24, 17, 2), (24, 22, 3),
(25, 9, 1), (25, 12, 2), (25, 22, 3),
(26, 1, 1), (26, 7, 2), (26, 15, 3), (26, 29, 4),
(27, 25, 1), (27, 11, 2), (27, 5, 3),
(28, 1, 1), (28, 3, 2), (28, 4, 3),
(29, 19, 1), (29, 26, 2), (29, 8, 3),
(30, 2, 1), (30, 29, 2), (30, 7, 3);

-- 1.5 补充点赞数据
INSERT INTO `song_likes` (`song_id`, `user_id`) VALUES
(5, 1), (11, 1), (14, 1), (25, 1), (30, 1),
(7, 2), (8, 2), (21, 2), (18, 2),
(15, 3), (27, 3), (28, 3),
(19, 4), (23, 4), (6, 4),
(1, 5), (2, 5), (29, 5);

-- ==========================================
-- 2. 仪表盘演示数据
-- ==========================================

-- 将用户的注册时间打散到近 7 天内
UPDATE `users` SET `created_at` = DATE_SUB(CURDATE(), INTERVAL id DAY);

-- 将歌曲的上传时间打散到近 7 天内，制造折线图波动
UPDATE `songs` SET `upload_at` = DATE_SUB(CURDATE(), INTERVAL (id % 7) DAY);

-- 制造歌单审核状态分布
UPDATE `playlists` SET `status` = 0 WHERE `id` IN (4, 11, 19, 26, 28);
UPDATE `playlists` SET `status` = 2, `reject_reason` = '封面包含低俗违规信息' WHERE `id` IN (8, 22);

-- 插入今日播放
INSERT INTO `play_histories` (`user_id`, `song_id`, `duration`, `source`, `played_at`) VALUES
(1, 5, 295, '每日推荐', NOW()), (2, 7, 229, '搜索', NOW()), (3, 15, 297, '歌单', NOW()),
(4, 19, 239, '我喜欢的音乐', NOW()), (5, 1, 215, '每日推荐', NOW()), (1, 23, 120, '歌单', NOW()),
(2, 8, 235, '歌单', NOW()), (3, 28, 308, '每日推荐', NOW()), (1, 11, 215, '搜索', NOW()),
(4, 23, 261, '歌手主页', NOW()), (5, 2, 240, '我喜欢的音乐', NOW()), (1, 9, 251, '歌单', NOW()),
(2, 12, 236, '搜索', NOW()), (3, 4, 205, '每日推荐', NOW()), (4, 6, 276, '排行榜', NOW()),
(5, 29, 186, '排行榜', NOW()), (1, 10, 230, '歌单', NOW()), (2, 21, 258, '歌手主页', NOW()),
(3, 27, 206, '每日推荐', NOW()), (4, 14, 259, '我喜欢的音乐', NOW()), (5, 30, 239, '搜索', NOW()),
(1, 13, 257, '歌单', NOW()), (2, 18, 289, '我喜欢的音乐', NOW()), (3, 25, 223, '每日推荐', NOW());

-- 插入近 1~6 天的历史播放记录
INSERT INTO `play_histories` (`user_id`, `song_id`, `duration`, `source`, `played_at`) VALUES
(1, 1, 100, '搜索', DATE_SUB(NOW(), INTERVAL 1 DAY)), (2, 2, 120, '歌单', DATE_SUB(NOW(), INTERVAL 1 DAY)),
(3, 3, 200, '每日推荐', DATE_SUB(NOW(), INTERVAL 1 DAY)), (4, 4, 150, '我喜欢的音乐', DATE_SUB(NOW(), INTERVAL 1 DAY)),
(5, 5, 210, '排行榜', DATE_SUB(NOW(), INTERVAL 1 DAY)), (1, 6, 180, '搜索', DATE_SUB(NOW(), INTERVAL 1 DAY)),
(2, 7, 220, '每日推荐', DATE_SUB(NOW(), INTERVAL 1 DAY)), (3, 8, 190, '歌单', DATE_SUB(NOW(), INTERVAL 1 DAY)),
(4, 9, 240, '歌手主页', DATE_SUB(NOW(), INTERVAL 1 DAY)), (5, 10, 260, '每日推荐', DATE_SUB(NOW(), INTERVAL 1 DAY)),
(1, 11, 215, '我喜欢的音乐', DATE_SUB(NOW(), INTERVAL 1 DAY)), (2, 12, 230, '搜索', DATE_SUB(NOW(), INTERVAL 1 DAY)),
(3, 13, 250, '歌单', DATE_SUB(NOW(), INTERVAL 2 DAY)), (4, 14, 260, '排行榜', DATE_SUB(NOW(), INTERVAL 2 DAY)),
(5, 15, 290, '我喜欢的音乐', DATE_SUB(NOW(), INTERVAL 2 DAY)), (1, 16, 240, '每日推荐', DATE_SUB(NOW(), INTERVAL 2 DAY)),
(2, 17, 270, '搜索', DATE_SUB(NOW(), INTERVAL 2 DAY)), (3, 18, 280, '歌单', DATE_SUB(NOW(), INTERVAL 2 DAY)),
(4, 19, 230, '歌手主页', DATE_SUB(NOW(), INTERVAL 2 DAY)), (5, 20, 270, '我喜欢的音乐', DATE_SUB(NOW(), INTERVAL 2 DAY)),
(1, 21, 250, '每日推荐', DATE_SUB(NOW(), INTERVAL 2 DAY)), (2, 22, 170, '搜索', DATE_SUB(NOW(), INTERVAL 2 DAY)),
(3, 23, 260, '歌单', DATE_SUB(NOW(), INTERVAL 2 DAY)), (4, 24, 310, '排行榜', DATE_SUB(NOW(), INTERVAL 2 DAY)),
(5, 25, 220, '每日推荐', DATE_SUB(NOW(), INTERVAL 2 DAY)), (1, 26, 270, '我喜欢的音乐', DATE_SUB(NOW(), INTERVAL 2 DAY)),
(2, 27, 200, '歌手主页', DATE_SUB(NOW(), INTERVAL 2 DAY)), (3, 28, 300, '搜索', DATE_SUB(NOW(), INTERVAL 2 DAY)),
(4, 29, 180, '歌单', DATE_SUB(NOW(), INTERVAL 2 DAY)), (5, 30, 230, '每日推荐', DATE_SUB(NOW(), INTERVAL 2 DAY)),
(1, 5, 295, '搜索', DATE_SUB(NOW(), INTERVAL 3 DAY)), (2, 11, 215, '每日推荐', DATE_SUB(NOW(), INTERVAL 3 DAY)),
(3, 15, 297, '歌单', DATE_SUB(NOW(), INTERVAL 4 DAY)), (4, 23, 120, '我喜欢的音乐', DATE_SUB(NOW(), INTERVAL 4 DAY)),
(5, 7, 229, '排行榜', DATE_SUB(NOW(), INTERVAL 5 DAY)), (1, 8, 235, '歌手主页', DATE_SUB(NOW(), INTERVAL 5 DAY)),
(2, 28, 308, '每日推荐', DATE_SUB(NOW(), INTERVAL 6 DAY)), (3, 19, 239, '搜索', DATE_SUB(NOW(), INTERVAL 6 DAY));
