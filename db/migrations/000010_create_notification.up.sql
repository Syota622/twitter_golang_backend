CREATE TABLE IF NOT EXISTS notifications (
  id SERIAL PRIMARY KEY,
  user_id INTEGER NOT NULL,
  notified_by_id INTEGER NOT NULL,
  type TEXT NOT NULL,
  post_id INTEGER,
  comment_id INTEGER,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- user_id: 通知を受け取るユーザーのID
-- notified_by_id: 通知を作成したユーザーのID
-- type: 通知の種類（例：'like', 'follow', 'comment'）
-- post_id: 関連する投稿のID（'like'や'comment'の場合）
-- comment_id: 関連するコメントのID（'comment'の場合）