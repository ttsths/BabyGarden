import { useEffect, useMemo, useState, type ReactNode } from 'react';
import { NavLink, Route, Routes, useNavigate, useParams, useSearchParams } from 'react-router-dom';
import { api } from '../services/api';
import { useAuthStore } from '../stores/useAuthStore';
import '../styles/open-design.css';

type NavItem = { to: string; label: string; short: string; description: string };
type RecordKind = 'feeding' | 'sleep' | 'diaper' | 'temperature';
type FeedingMethod = 'formula' | 'left' | 'right';
type StatsMode = '日' | '周' | '月' | '自定义';

type DailyStats = {
  feeding: { count: number; total_amount: number; average_amount?: number };
  sleep: { count: number; total_hours: number; average_duration_hours?: number; daytime_single_hours?: number };
  diaper: { count: number };
  temperature?: { count: number; latest: number };
};

type WeeklyStats = {
  dates: string[];
  feeding: number[];
  sleep: number[];
  diaper: number[];
  daily_average_sleep_hours?: number[];
  daytime_single_sleep_hours?: number[];
  daily_average_milk_amount?: number[];
  temperature_latest?: number[];
};

type YuanziRecord = {
  id: string;
  baby_id?: string;
  type: RecordKind;
  started_at: string;
  ended_at?: string;
  content: Record<string, unknown>;
  note?: string;
  duration_hours?: number;
};

type YuanziPhoto = {
  id: string;
  url?: string;
  thumb_url?: string;
  taken_at?: string;
  description?: string;
  likes_count?: number;
  comments_count?: number;
  liked_by_me?: boolean;
  recent_comments?: PhotoComment[];
};

type PhotoComment = {
  id: string;
  photo_id: string;
  user_id: string;
  nickname: string;
  avatar_url?: string;
  content: string;
  created_at: string;
};

type FamilyMember = {
  user_id: string;
  nickname: string;
  avatar_url?: string;
  role: string;
  elder_mode: boolean;
};

const navItems: NavItem[] = [
  { to: '/', label: '首页', short: 'H', description: '今日概览' },
  { to: '/record', label: '记录', short: 'R', description: '3 秒完成' },
  { to: '/stats', label: '统计', short: 'S', description: '趋势判断' },
  { to: '/ai', label: 'AI 小助手', short: 'A', description: '带记录上下文' },
  { to: '/photos', label: '照片墙', short: 'P', description: '家庭互动' },
  { to: '/family', label: '家庭共享', short: 'F', description: '权限同步' },
  { to: '/settings', label: '设置', short: 'M', description: '提醒与隐私' },
];

const fallbackStats: DailyStats = {
  feeding: { count: 5, total_amount: 620 },
  sleep: { count: 2, total_hours: 12 },
  diaper: { count: 3 },
};

const fallbackWeekly: WeeklyStats = {
  dates: ['05-14', '05-15', '05-16', '05-17', '05-18', '05-19', '05-20'],
  feeding: [4, 5, 4, 6, 5, 6, 5],
  sleep: [11.5, 12, 10.8, 12.6, 12.2, 13, 12],
  diaper: [3, 4, 2, 4, 3, 5, 3],
};

const fallbackRecords: YuanziRecord[] = [
  { id: 'mock-feed-1', type: 'feeding', started_at: '2026-05-20T10:30:00+08:00', content: { type: 'formula', amount: 120, unit: 'ml' }, note: '奶粉喂养，状态安稳' },
  { id: 'mock-diaper-1', type: 'diaper', started_at: '2026-05-20T09:15:00+08:00', content: { type: 'mixed', color: '黄色', consistency: '糊状' }, note: '记录为正常' },
  { id: 'mock-sleep-1', type: 'sleep', started_at: '2026-05-20T07:00:00+08:00', ended_at: '2026-05-20T09:00:00+08:00', content: { quality: 'stable', location: 'crib' }, note: '夜间连续睡眠 6 小时' },
  { id: 'mock-temp-1', type: 'temperature', started_at: '2026-05-20T06:40:00+08:00', content: { value: 36.5, unit: 'C', position: 'armpit' }, note: '体温正常' },
];

const fallbackPhotos: YuanziPhoto[] = [
  { id: 'photo-1', description: '第一次翻身', taken_at: '2026-05-20 10:40' },
  { id: 'photo-2', description: '阳台晒太阳', taken_at: '2026-05-20 11:20' },
  { id: 'photo-3', description: '睡醒笑了', taken_at: '2026-05-20 14:10' },
  { id: 'photo-4', description: '奶后拍嗝', taken_at: '2026-05-19 19:00' },
];

const fallbackMembers: FamilyMember[] = [
  { user_id: 'mom', nickname: '妈妈', role: 'admin', elder_mode: false },
  { user_id: 'dad', nickname: '爸爸', role: 'member', elder_mode: false },
  { user_id: 'grandma', nickname: '外婆', role: 'elder', elder_mode: true },
];

const today = new Date().toISOString().slice(0, 10);
const demoBabyId = import.meta.env.VITE_DEMO_BABY_ID || localStorage.getItem('yuanzi-demo-baby-id') || '';
const demoFamilyId = import.meta.env.VITE_DEMO_FAMILY_ID || localStorage.getItem('yuanzi-demo-family-id') || '';

export function OpenDesignApp() {
  return (
    <Routes>
      <Route path="/login" element={<LoginScreen />} />
      <Route path="/baby-profile" element={<BabyProfileScreen />} />
      <Route path="/*" element={<ProductShell />} />
    </Routes>
  );
}

function useYuanziContext() {
  const token = useAuthStore((state) => state.token);
  const isAuthenticated = useAuthStore((state) => state.isAuthenticated);
  const [babyId, setBabyId] = useState(demoBabyId);
  const [familyId, setFamilyId] = useState(demoFamilyId);
  const [babyName, setBabyName] = useState('小园子');
  const [source, setSource] = useState(demoBabyId ? 'seed' : 'mock');

  useEffect(() => {
    let active = true;
    async function load() {
      if (!token && !isAuthenticated) return;
      try {
        const response = await api.baby.getList();
        const babies = unwrap<Array<Record<string, unknown>>>(response, []);
        const current = babies[0];
        if (!active || !current) return;
        const nextBabyId = String(current.id || '');
        const nextFamilyId = String(current.family_id || current.familyId || '');
        setBabyId(nextBabyId);
        setFamilyId(nextFamilyId);
        setBabyName(String(current.name || '小园子'));
        if (nextBabyId) localStorage.setItem('yuanzi-demo-baby-id', nextBabyId);
        if (nextFamilyId) localStorage.setItem('yuanzi-demo-family-id', nextFamilyId);
        setSource('api');
      } catch {
        setSource(demoBabyId ? 'seed' : 'mock');
      }
    }
    void load();
    return () => { active = false; };
  }, [token, isAuthenticated]);

  return { babyId, familyId, babyName, source };
}

function ProductShell() {
  return (
    <main className="od-shell">
      <div className="od-app">
        <Sidebar />
        <section className="od-main">
          <Routes>
            <Route index element={<DashboardScreen />} />
            <Route path="record" element={<RecordScreen />} />
            <Route path="records" element={<RecordsScreen />} />
            <Route path="record-detail/:id" element={<RecordDetailScreen />} />
            <Route path="stats" element={<StatsScreen />} />
            <Route path="ai" element={<AiScreen />} />
            <Route path="photos" element={<PhotosScreen />} />
            <Route path="family" element={<FamilyScreen />} />
            <Route path="settings" element={<SettingsScreen />} />
            <Route path="elder" element={<ElderScreen />} />
          </Routes>
        </section>
      </div>
      <MobileNav />
    </main>
  );
}

function DashboardScreen() {
  const navigate = useNavigate();
  const { babyId, babyName, source: contextSource } = useYuanziContext();
  const [stats, setStats] = useState<DailyStats>(fallbackStats);
  const [records, setRecords] = useState<YuanziRecord[]>(fallbackRecords);
  const [photos, setPhotos] = useState<YuanziPhoto[]>(fallbackPhotos);
  const [source, setSource] = useState(contextSource);

  useEffect(() => {
    let active = true;
    async function load() {
      if (!babyId) return;
      try {
        const [daily, recordList, photoList] = await Promise.all([
          api.stats.daily(babyId, today),
          api.record.getList(babyId, { page: 1, page_size: 4 }),
          api.photo.getList(babyId, { page: 1, page_size: 3 }),
        ]);
        if (!active) return;
        setStats(unwrap<DailyStats>(daily, fallbackStats));
        setRecords(unwrapList<YuanziRecord>(recordList, fallbackRecords));
        setPhotos(unwrapList<YuanziPhoto>(photoList, fallbackPhotos));
        setSource('api');
      } catch {
        setSource(contextSource);
      }
    }
    void load();
    return () => { active = false; };
  }, [babyId, contextSource]);

  const progress = Math.min(100, Math.round(((stats.feeding.count / 6) + Math.min(stats.sleep.total_hours / 12, 1) + (stats.diaper.count / 4)) / 3 * 100));

  return (
    <>
      <ScreenTitle
        eyebrow={`TODAY · ${today}`}
        title={`照看今天，也看见${babyName}的成长节奏`}
        desc={`把喂养、睡眠、排泄、体温和家庭照片放到同一个温柔但高效的工作台里。数据源：${source === 'api' ? '后端 API' : '本地示例'}`}
        action={<div className="od-actions"><button className="od-btn" onClick={() => navigate('/family')}>家庭成员</button><button className="od-btn primary" onClick={() => navigate('/record')}>快速记录</button></div>}
      />
      <div className="od-dashboard">
        <div className="od-stack">
          <section className="od-panel od-hero">
            <div className="od-hero-head">
              <h2>今日记录完整度良好，喂奶和睡眠节奏稳定。</h2>
              <span>距下次喂奶约 42 分钟</span>
            </div>
            <div className="od-progress">
              <div className="od-ring" style={{ background: `conic-gradient(#ff9a8b ${progress}%, #ffe4dc 0)` }}><strong>{progress}%</strong></div>
              <div className="od-task-list">
                <TaskRow label="喂奶" value={`${stats.feeding.count}/6`} progress={Math.min(100, stats.feeding.count / 6 * 100)} />
                <TaskRow label="睡眠" value={`${round(stats.sleep.total_hours)}h`} progress={Math.min(100, stats.sleep.total_hours / 12 * 100)} />
                <TaskRow label="排泄" value={`${stats.diaper.count}/4`} progress={Math.min(100, stats.diaper.count / 4 * 100)} />
                <TaskRow label="体温" value="1/2" progress={50} />
              </div>
            </div>
          </section>

          <SectionHead title="今日概览" note="get_api_v1_stats_daily" />
          <section className="od-stat-grid">
            <StatCard label="喂奶" value={`${stats.feeding.count} 次`} note={`总量 ${stats.feeding.total_amount}ml`} icon="N" />
            <StatCard label="睡眠" value={`${round(stats.sleep.total_hours)}h`} note={`${stats.sleep.count} 段睡眠`} icon="S" />
            <StatCard label="排泄" value={`${stats.diaper.count} 次`} note="颜色与形态正常" icon="D" />
            <StatCard label="体温" value="36.5" note="最近测量 06:40" icon="T" />
          </section>

          <SectionHead title="快速记录" note="常用动作置前" />
          <section className="od-quick-grid">
            {[
              ['喂奶', 'feeding', 'N'],
              ['睡眠', 'sleep', 'S'],
              ['换尿布', 'diaper', 'D'],
              ['体温', 'temperature', 'T'],
            ].map(([label, type, icon]) => (
              <button key={type} className="od-quick" onClick={() => navigate(`/record?type=${type}`)}>
                <span>{icon}</span><strong>{label}</strong>
              </button>
            ))}
          </section>
        </div>

        <aside className="od-stack">
          <section className="od-panel od-ai-card">
            <h2>AI 育儿助手</h2>
            <p>结合今天的记录，优先回答与睡眠、喂养、体温相关的问题。</p>
            <form className="od-ask-box"><input aria-label="向 AI 提问" placeholder="例如：宝宝午睡少怎么办？" /><button type="button" onClick={() => navigate('/ai')}>提问</button></form>
          </section>
          <SectionHead title="最近记录" note={<button className="od-link" onClick={() => navigate('/records')}>查看全部</button>} />
          <Timeline records={records} compact />
          <SectionHead title="家庭照片" note="get_api_v1_photo" />
          <div className="od-photo-strip">{photos.slice(0, 3).map((photo) => <PhotoTile key={photo.id} photo={photo} />)}</div>
        </aside>
      </div>
    </>
  );
}

function RecordScreen() {
  const [searchParams] = useSearchParams();
  const { babyId } = useYuanziContext();
  const [selectedType, setSelectedType] = useState<RecordKind>(coerceRecordKind(searchParams.get('type')));
  const [feedingMethod, setFeedingMethod] = useState<FeedingMethod>('formula');
  const [startedAt, setStartedAt] = useState(toLocalInputValue(new Date()));
  const [amount, setAmount] = useState(120);
  const [temperature, setTemperature] = useState(36.5);
  const [note, setNote] = useState('喝奶后状态安稳，拍嗝后无明显吐奶。');
  const [picker, setPicker] = useState<'time' | 'amount' | 'temperature' | null>(null);
  const [lastRecord, setLastRecord] = useState<YuanziRecord>(fallbackRecords[0]);
  const [status, setStatus] = useState('LIVE PREVIEW');

  useEffect(() => {
    setSelectedType(coerceRecordKind(searchParams.get('type')));
  }, [searchParams]);

  async function saveRecord() {
    const payload = buildRecordPayload(selectedType, startedAt, feedingMethod, amount, temperature, note);
    try {
      if (!babyId) throw new Error('missing baby id');
      const response = await api.record.create({ baby_id: babyId, ...payload });
      const saved = unwrap<YuanziRecord>(response, { id: `local-${Date.now()}`, baby_id: babyId, ...payload });
      setLastRecord(saved);
      setStatus('已保存到后端');
    } catch {
      setLastRecord({ id: `local-${Date.now()}`, ...payload });
      setStatus('后端不可用，已保留本地预览');
    }
  }

  const previewPayload = buildRecordPayload(selectedType, startedAt, feedingMethod, amount, temperature, note);
  const preview = formatRecord(lastRecord.type === selectedType ? lastRecord : previewPayload);

  return (
    <>
      <ScreenTitle eyebrow="RECORD · ALL-IN-ONE" title="一页完成核心记录" desc="入口按高频排序，表单只暴露当前类型需要的字段；保存后即时写入时间轴并同步家庭空间。" action={<button className="od-btn primary" onClick={saveRecord}>保存记录</button>} />
      <div className="od-content-grid">
        <section className="od-panel od-pad od-stack">
          <div className="od-choice-grid">
            <Choice active={selectedType === 'feeding'} title="喂奶" desc="母乳左右侧 / 配方奶 ml" onClick={() => setSelectedType('feeding')} />
            <Choice active={selectedType === 'sleep'} title="睡眠" desc="开始计时 / 结束入睡" onClick={() => setSelectedType('sleep')} />
            <Choice active={selectedType === 'diaper'} title="排泄" desc="大便 / 小便 / 混合" onClick={() => setSelectedType('diaper')} />
            <Choice active={selectedType === 'temperature'} title="测温" desc="体温 / 异常备注" onClick={() => setSelectedType('temperature')} />
          </div>
          {selectedType === 'feeding' && (
            <div className="od-segmented">
              {[
                ['formula', '配方奶'],
                ['left', '左侧母乳'],
                ['right', '右侧母乳'],
              ].map(([value, label]) => <button key={value} className={feedingMethod === value ? 'active' : ''} onClick={() => setFeedingMethod(value as FeedingMethod)}>{label}</button>)}
            </div>
          )}
          <div className="od-form-grid">
            <label>记录时间<input aria-label="记录时间" readOnly value={formatLocalTime(startedAt)} onClick={() => setPicker('time')} /></label>
            {selectedType === 'feeding' && <label>奶量 ml<input aria-label="奶量 ml" readOnly value={amount} onClick={() => setPicker('amount')} /></label>}
            {selectedType === 'temperature' && <label>体温 °C<input aria-label="体温 °C" readOnly value={temperature} onClick={() => setPicker('temperature')} /></label>}
            <label>宝宝状态<textarea value={note} onChange={(event) => setNote(event.target.value)} /></label>
          </div>
        </section>
        <aside className="od-record-preview">
          <div><div className="od-eyebrow">{status}</div><h2>{preview.title}</h2><p>{preview.desc}</p></div>
          <div><Metric label="下次提醒" value={selectedType === 'feeding' ? '14:30' : '已更新'} note="按设置自动计算" /><Metric label="家庭同步" value="3 人可见" note="SSE 准实时" /><Metric label="AI 上下文" value="已纳入" note="问答时引用" /></div>
        </aside>
      </div>
      {picker === 'time' && <PickerModal title="选择记录时间" onClose={() => setPicker(null)}><DateTimePicker value={startedAt} onChange={setStartedAt} /></PickerModal>}
      {picker === 'amount' && <PickerModal title="选择奶量" onClose={() => setPicker(null)}><NumberPicker min={30} max={260} step={10} value={amount} unit="ml" onChange={setAmount} /></PickerModal>}
      {picker === 'temperature' && <PickerModal title="选择体温" onClose={() => setPicker(null)}><NumberPicker min={350} max={390} step={1} value={Math.round(temperature * 10)} unit="°C" formatter={(value) => (value / 10).toFixed(1)} onChange={(value) => setTemperature(value / 10)} /></PickerModal>}
    </>
  );
}

function StatsScreen() {
  const navigate = useNavigate();
  const { babyId } = useYuanziContext();
  const [mode, setMode] = useState<StatsMode>('日');
  const [daily, setDaily] = useState<DailyStats>(fallbackStats);
  const [weekly, setWeekly] = useState<WeeklyStats>(fallbackWeekly);
  const [records, setRecords] = useState<YuanziRecord[]>(fallbackRecords);
  const [rangeStart, setRangeStart] = useState(today);
  const [rangeEnd, setRangeEnd] = useState(today);

  useEffect(() => {
    let active = true;
    async function load() {
      if (!babyId) return;
      try {
        const statsRequest = mode === '月'
          ? api.stats.monthly(babyId, today)
          : mode === '自定义'
            ? api.stats.range(babyId, rangeStart, rangeEnd)
            : api.stats.weekly(babyId, today);
        const [dailyRes, rangeRes, recordRes] = await Promise.all([
          api.stats.daily(babyId, today),
          statsRequest,
          api.record.getList(babyId, { page: 1, page_size: 10, date: today }),
        ]);
        if (!active) return;
        setDaily(unwrap<DailyStats>(dailyRes, fallbackStats));
        setWeekly(unwrap<WeeklyStats>(rangeRes, fallbackWeekly));
        setRecords(unwrapList<YuanziRecord>(recordRes, fallbackRecords));
      } catch {
        // 保留 mock，页面仍可交互。
      }
    }
    void load();
    return () => { active = false; };
  }, [babyId, mode, rangeStart, rangeEnd]);

  const chartValues = mode === '日'
    ? [daily.feeding.count, daily.sleep.total_hours, daily.diaper.count, 1]
    : mode === '周'
      ? weekly.daily_average_milk_amount || weekly.feeding
      : mode === '月'
        ? weekly.daily_average_sleep_hours || weekly.sleep
        : weekly.daytime_single_sleep_hours || weekly.sleep;

  return (
    <>
      <ScreenTitle eyebrow="TIMELINE · STATS" title="从记录流看到趋势" desc="保留时间轴的可追溯性，同时展示平均睡眠、白天单次睡眠和日均喝奶量。" action={<Segmented values={['日', '周', '月', '自定义']} value={mode} onChange={(value) => setMode(value as StatsMode)} />} />
      <div className="od-content-grid">
        <section className="od-panel od-pad"><h2>今日时间轴</h2><Timeline records={records} onRecordClick={(record) => navigate(`/record-detail/${record.id}`)} /></section>
        <aside className="od-stack">
          {mode === '自定义' && <section className="od-panel od-pad"><div className="od-form-grid two"><label>开始日期<input type="date" value={rangeStart} onChange={(event) => setRangeStart(event.target.value)} /></label><label>结束日期<input type="date" value={rangeEnd} onChange={(event) => setRangeEnd(event.target.value)} /></label></div></section>}
          <section className="od-panel od-pad"><h2>{mode}统计趋势</h2><Chart values={chartValues} /></section>
          <section className="od-stat-grid two"><StatCard label="日均喝奶" value={`${round(daily.feeding.average_amount || daily.feeding.total_amount)}ml`} note={`今日 ${daily.feeding.count} 次，总量 ${daily.feeding.total_amount}ml`} icon="N" /><StatCard label="平均睡眠" value={`${round(daily.sleep.average_duration_hours || daily.sleep.total_hours)}h`} note={`白天单次最长 ${round(daily.sleep.daytime_single_hours || 0)}h`} icon="S" /></section>
        </aside>
      </div>
    </>
  );
}

function PhotosScreen() {
  const { babyId } = useYuanziContext();
  const [photoList, setPhotoList] = useState<YuanziPhoto[]>(fallbackPhotos);
  const [selectedPhoto, setSelectedPhoto] = useState<YuanziPhoto>(fallbackPhotos[0]);
  const [uploadOpen, setUploadOpen] = useState(false);
  const [uploadStatus, setUploadStatus] = useState('');
  const [comment, setComment] = useState('');

  useEffect(() => {
    async function load() {
      if (!babyId) return;
      try {
        const response = await api.photo.getList(babyId, { page: 1, page_size: 20 });
        const list = unwrapList<YuanziPhoto>(response, fallbackPhotos);
        setPhotoList(list);
        setSelectedPhoto(list[0] || fallbackPhotos[0]);
      } catch {
        // 使用本地示例。
      }
    }
    void load();
  }, [babyId]);

  async function requestUpload(file: File) {
    setUploadStatus('正在向后端申请上传地址...');
    try {
      if (!babyId) throw new Error('missing baby id');
      const response = await api.photo.getUploadUrl({ baby_id: babyId, filename: file.name, content_type: file.type || 'image/jpeg', size: file.size });
      const data = unwrap<Record<string, string | number>>(response, {});
      setUploadStatus(`已获取上传地址：${String(data.upload_url || 'upload_url 返回为空')}`);
    } catch {
      setUploadStatus('后端上传地址暂不可用，已停留在上传组件。');
    }
  }

  async function toggleLike() {
    try {
      const response = selectedPhoto.liked_by_me ? await api.photo.unlike(selectedPhoto.id) : await api.photo.like(selectedPhoto.id);
      const summary = unwrap<Partial<YuanziPhoto>>(response, {});
      const updated = { ...selectedPhoto, ...summary };
      setSelectedPhoto(updated);
      setPhotoList((items) => items.map((item) => item.id === updated.id ? updated : item));
    } catch {
      setUploadStatus('点赞接口暂不可用。');
    }
  }

  async function sendComment() {
    if (!comment.trim()) return;
    try {
      const response = await api.photo.comment(selectedPhoto.id, comment.trim());
      const created = unwrap<PhotoComment>(response, {
        id: `local-${Date.now()}`,
        photo_id: selectedPhoto.id,
        user_id: 'local',
        nickname: '我',
        content: comment.trim(),
        created_at: new Date().toISOString(),
      });
      const updated = {
        ...selectedPhoto,
        comments_count: (selectedPhoto.comments_count || 0) + 1,
        recent_comments: [...(selectedPhoto.recent_comments || []), created],
      };
      setSelectedPhoto(updated);
      setPhotoList((items) => items.map((item) => item.id === updated.id ? updated : item));
      setComment('');
    } catch {
      setUploadStatus('评论接口暂不可用。');
    }
  }

  return (
    <>
      <ScreenTitle eyebrow="PHOTO WALL · FAMILY MEMORY" title="家庭照片和成长记录在一起" desc="照片墙保留上传状态、家庭互动和关联记录，避免照片只是散落在相册里。" action={<button className="od-btn primary" onClick={() => setUploadOpen(true)}>上传照片</button>} />
      <div className="od-content-grid">
        <section className="od-panel od-pad"><div className="od-photo-grid">{photoList.map((photo) => <PhotoTile key={photo.id} photo={photo} active={selectedPhoto.id === photo.id} onClick={() => setSelectedPhoto(photo)} />)}</div></section>
        <aside className="od-stack">
          <section className="od-panel od-pad"><h2>照片详情</h2><Metric label="当前照片" value={selectedPhoto.description || selectedPhoto.id} note={selectedPhoto.taken_at || '家庭空间'} /><Metric label="点赞" value={`${selectedPhoto.likes_count || 0}`} note={selectedPhoto.liked_by_me ? '我已点赞' : '点击下方点赞'} /><Metric label="评论" value={`${selectedPhoto.comments_count || 0}`} note="家庭成员互动" /><button className="od-btn full" onClick={toggleLike}>{selectedPhoto.liked_by_me ? '取消点赞' : '点赞'}</button></section>
          <section className="od-invite-card"><h2>家庭互动</h2>{(selectedPhoto.recent_comments || []).map((item) => <p key={item.id}><b>{item.nickname || '家人'}：</b>{item.content}</p>)}<div className="od-ask-box light"><input value={comment} onChange={(event) => setComment(event.target.value)} placeholder="在这张照片下评论" /><button type="button" onClick={sendComment}>发送</button></div></section>
        </aside>
      </div>
      {uploadOpen && <PickerModal title="上传照片" onClose={() => setUploadOpen(false)}><div className="od-upload-box"><input type="file" accept="image/*" onChange={(event) => { const file = event.target.files?.[0]; if (file) void requestUpload(file); }} /><p>{uploadStatus || '选择照片后会调用 post_api_v1_photo_upload_url 获取直传地址。'}</p></div></PickerModal>}
    </>
  );
}

function FamilyScreen() {
  const { familyId } = useYuanziContext();
  const [members, setMembers] = useState<FamilyMember[]>(fallbackMembers);
  const [inviteOpen, setInviteOpen] = useState(false);
  const [selectedMember, setSelectedMember] = useState<FamilyMember | null>(null);
  const [invitePhone, setInvitePhone] = useState('');
  const [inviteRole, setInviteRole] = useState<'member' | 'elder'>('member');
  const [inviteResult, setInviteResult] = useState('');
  const [joinCode, setJoinCode] = useState('');

  useEffect(() => {
    async function load() {
      if (!familyId) return;
      try {
        const response = await api.family.getMembers(familyId);
        setMembers(unwrap<FamilyMember[]>(response, fallbackMembers));
      } catch {
        // 使用本地示例。
      }
    }
    void load();
  }, [familyId]);

  async function inviteMember() {
    try {
      if (!familyId) throw new Error('missing family id');
      const response = await api.family.invite(familyId, invitePhone, inviteRole);
      const data = unwrap<Record<string, string | number>>(response, {});
      setInviteResult(`邀请已发送，邀请码 ${String(data.invite_code || '已生成')}`);
    } catch {
      setInviteResult('后端邀请接口暂不可用，已保留输入。');
    }
  }

  async function joinFamily() {
    try {
      const response = await api.family.join(joinCode, inviteRole);
      const data = unwrap<Record<string, string>>(response, {});
      setInviteResult(`已加入家庭：${data.name || data.id || joinCode}`);
    } catch {
      setInviteResult('加入失败，请检查邀请码或登录状态。');
    }
  }

  async function leaveFamily() {
    try {
      if (!familyId) throw new Error('missing family id');
      await api.family.leave(familyId);
      setInviteResult('已离开家庭');
      setMembers([]);
    } catch {
      setInviteResult('离开家庭失败，最后一个管理员不能离开。');
    }
  }

  return (
    <>
      <ScreenTitle eyebrow="FAMILY SHARE · SSE SYNC" title="每个人只看到自己该做的事" desc="家庭邀请、成员权限和同步状态集中呈现，降低新手父母解释成本。" action={<button className="od-btn primary" onClick={() => setInviteOpen(true)}>邀请家人</button>} />
      <div className="od-content-grid">
        <section className="od-panel od-pad od-stack">{members.map((member) => <MemberRow key={member.user_id} member={member} onClick={() => setSelectedMember(member)} />)}</section>
        <aside className="od-stack"><section className="od-invite-card"><h2>邀请链接</h2><p>链接 24 小时有效。加入前需确认儿童隐私声明和照片可见范围。</p><button className="od-btn primary" onClick={() => setInviteOpen(true)}>按手机号邀请</button></section><section className="od-panel od-pad"><h2>加入 / 离开</h2><div className="od-ask-box light"><input value={joinCode} onChange={(event) => setJoinCode(event.target.value.toUpperCase())} placeholder="输入邀请码" /><button type="button" onClick={joinFamily}>加入</button></div><button className="od-btn full" onClick={leaveFamily}>离开家庭</button><p className="od-muted">{inviteResult}</p></section><section className="od-panel od-pad"><h2>同步状态</h2><Metric label="记录" value="准实时" note="SSE connected" /><Metric label="照片" value="排队上传" note="2 张待同步" /></section></aside>
      </div>
      {inviteOpen && <PickerModal title="邀请家人" onClose={() => setInviteOpen(false)}><div className="od-form-grid"><label>手机号<input value={invitePhone} onChange={(event) => setInvitePhone(event.target.value)} placeholder="请输入 11 位手机号" /></label><label>角色<select value={inviteRole} onChange={(event) => setInviteRole(event.target.value as 'member' | 'elder')}><option value="member">照护成员</option><option value="elder">祖辈模式</option></select></label><button className="od-btn primary full" onClick={inviteMember}>发送邀请</button><p className="od-muted">{inviteResult}</p></div></PickerModal>}
      {selectedMember && <PickerModal title="家庭成员详情" onClose={() => setSelectedMember(null)}><Metric label="成员" value={selectedMember.nickname || selectedMember.user_id} note={selectedMember.user_id} /><Metric label="角色" value={roleText(selectedMember.role)} note={selectedMember.elder_mode ? '已开启祖辈模式' : '标准模式'} /></PickerModal>}
    </>
  );
}

function SettingsScreen() {
  const [mode, setMode] = useState('标准模式');
  const [phoneOpen, setPhoneOpen] = useState(false);
  const [picker, setPicker] = useState<'feed' | 'practice' | null>(null);
  const [feedInterval, setFeedInterval] = useState(4);
  const [practiceHour, setPracticeHour] = useState(20);
  const [phone, setPhone] = useState('138****2026');

  return (
    <>
      <ScreenTitle eyebrow="SETTINGS · PRIVACY" title="把高风险设置放在清晰分组里" desc="设置页覆盖提醒、显示模式、家庭权限和隐私声明；儿童数据相关动作必须说明影响范围。" action={<button className="od-btn primary">保存设置</button>} />
      <div className="od-content-grid equal">
        <section className="od-panel od-pad"><h2>照护提醒</h2><Metric label="喂奶间隔" value={`${feedInterval} 小时`} note="点击更换" onClick={() => setPicker('feed')} /><Metric label="睡眠结束" value="开启" note="计时器提醒" /><Metric label="成长练习" value={`${practiceHour}:00`} note="点击选择小时" onClick={() => setPicker('practice')} /></section>
        <section className="od-panel od-pad od-stack"><h2>显示与模式</h2>{['夜间暗色模式', '祖辈极简模式', '标准模式'].map((item) => <Choice key={item} active={mode === item} title={item} desc={item === '祖辈极简模式' ? '大字体、高对比、一键记录' : '点击切换显示偏好'} onClick={() => setMode(item)} />)}</section>
        <section className="od-panel od-pad"><h2>隐私与数据</h2><p className="od-muted">展示儿童隐私保护声明、照片上传范围、家庭成员权限、AI 问答上下文使用说明。MVP 不提供数据导出入口。</p></section>
        <section className="od-panel od-pad"><h2>账号安全</h2><Metric label="手机号" value={phone} note="点击更换" onClick={() => setPhoneOpen(true)} /><Metric label="登录设备" value="2 台" note="iPhone / Web" /></section>
      </div>
      {picker === 'feed' && <PickerModal title="喂奶间隔" onClose={() => setPicker(null)}><NumberPicker min={2} max={6} step={1} value={feedInterval} unit="小时" onChange={setFeedInterval} /></PickerModal>}
      {picker === 'practice' && <PickerModal title="成长练习提醒" onClose={() => setPicker(null)}><NumberPicker min={6} max={23} step={1} value={practiceHour} unit="点" onChange={setPracticeHour} /></PickerModal>}
      {phoneOpen && <PickerModal title="更换手机号" onClose={() => setPhoneOpen(false)}><div className="od-form-grid"><label>新手机号<input placeholder="请输入新手机号" onChange={(event) => setPhone(maskPhone(event.target.value))} /></label><button className="od-btn primary full" onClick={() => setPhoneOpen(false)}>确认更换</button></div></PickerModal>}
    </>
  );
}

function RecordsScreen() {
  const navigate = useNavigate();
  const { babyId } = useYuanziContext();
  const [records, setRecords] = useState<YuanziRecord[]>(fallbackRecords);
  useEffect(() => {
    async function load() {
      if (!babyId) return;
      try {
        const response = await api.record.getList(babyId, { page: 1, page_size: 30 });
        setRecords(unwrapList<YuanziRecord>(response, fallbackRecords));
      } catch {
        // 使用本地示例。
      }
    }
    void load();
  }, [babyId]);
  return <><ScreenTitle eyebrow="RECORD DETAIL · LIST" title="记录明细" desc="从首页最近记录进入，查看全部记录并继续打开单条详情。" /><section className="od-panel od-pad"><Timeline records={records} onRecordClick={(record) => navigate(`/record-detail/${record.id}`)} /></section></>;
}

function RecordDetailScreen() {
  const { id } = useParams();
  const [record, setRecord] = useState<YuanziRecord>(fallbackRecords.find((item) => item.id === id) || fallbackRecords[0]);
  useEffect(() => {
    async function load() {
      if (!id || id.startsWith('mock')) return;
      try {
        const response = await api.record.getDetail(id);
        setRecord((current) => unwrap<YuanziRecord>(response, current));
      } catch {
        // 保留当前展示。
      }
    }
    void load();
  }, [id]);
  const display = formatRecord(record);
  return <><ScreenTitle eyebrow="RECORD DETAIL" title={display.title} desc={display.desc} action={<NavLink className="od-btn" to="/records">返回明细</NavLink>} /><section className="od-panel od-pad"><Metric label="记录ID" value={record.id} note={record.type} /><Metric label="时间" value={formatTime(record.started_at)} note={record.ended_at ? formatTime(record.ended_at) : '无结束时间'} /><pre className="od-json">{JSON.stringify(record.content, null, 2)}</pre></section></>;
}

function AiScreen() {
  const { babyId } = useYuanziContext();
  const [question, setQuestion] = useState('宝宝今天午睡少 25 分钟，要不要提前哄睡？');
  const [answer, setAnswer] = useState('先观察精神状态和傍晚困倦信号。若精神好，可以保持原入睡节奏；若明显烦躁，可把睡前流程提前 20 分钟。');
  const [history, setHistory] = useState<Array<{ id: string; question: string; answer: string; created_at: string }>>([]);
  const [status, setStatus] = useState('历史会话加载中');

  useEffect(() => {
    async function load() {
      try {
        const response = await api.ai.history({ baby_id: babyId, page: 1, page_size: 8 });
        setHistory(unwrapList(response, []));
        setStatus('历史会话');
      } catch {
        setStatus('历史会话暂不可用');
      }
    }
    void load();
  }, [babyId]);

  async function sendQuestion() {
    if (!question.trim()) return;
    setStatus('AI 正在思考');
    try {
      const response = await api.ai.chat(question.trim(), {
        baby_id: babyId || undefined,
        history: history.slice(0, 4).flatMap((item) => [
          { role: 'user', content: item.question },
          { role: 'assistant', content: item.answer },
        ]),
      });
      const data = unwrap<{ answer: string }>(response, { answer: 'AI暂时无法回答，请稍后再试。' });
      setAnswer(data.answer);
      setHistory((items) => [{ id: `local-${Date.now()}`, question, answer: data.answer, created_at: new Date().toISOString() }, ...items]);
      setQuestion('');
      setStatus('已保存到历史会话');
    } catch {
      setStatus('AI 服务暂不可用');
    }
  }

  return <><ScreenTitle eyebrow="AI CHAT · VOICE INPUT" title="抱娃时也能开口问" desc="AI 回答引用近期记录，但避免医疗诊断式结论；历史会话由后端保存。" action={<button className="od-btn primary">按住说话</button>} /><div className="od-content-grid"><section className="od-panel od-pad od-stack"><article className="od-event"><span className="od-dot brand" /><div><strong>{history[0]?.question || '宝宝今天午睡少 25 分钟，要不要提前哄睡？'}</strong><p>{status}</p></div><time>{history[0] ? formatClock(history[0].created_at) : '10:42'}</time></article><article className="od-advice"><h2>Yuanzi 建议</h2><p>{answer}</p></article><form className="od-ask-box light"><input value={question} onChange={(event) => setQuestion(event.target.value)} placeholder="继续追问：今晚睡前怎么安排？" /><button type="button" onClick={sendQuestion}>发送</button></form></section><aside className="od-voice"><h2>历史 AI 会话</h2><div className="od-stack">{history.map((item) => <button className="od-choice" key={item.id} onClick={() => setAnswer(item.answer)}><strong>{item.question}</strong><small>{formatTime(item.created_at)}</small></button>)}</div><div className="od-wave">{[14, 36, 24, 58, 30, 68, 42, 76, 44, 64, 28, 48, 22, 36, 18, 30, 12].map((height, index) => <span key={index} style={{ height: `${18 + height}px` }} />)}</div><button>按住说话</button></aside></div></>;
}

function LoginScreen() {
  const navigate = useNavigate();
  const loginStore = useAuthStore((state) => state.login);
  const [phone, setPhone] = useState('13800138000');
  const [code, setCode] = useState('');
  const [status, setStatus] = useState('');

  async function sendCode() {
    try {
      await api.auth.sendCode(phone.replace(/\D/g, ''));
      setStatus('验证码已发送');
    } catch {
      setStatus('验证码发送失败，请检查后端短信/Redis 配置');
    }
  }

  async function login() {
    try {
      const response = await api.auth.login(phone.replace(/\D/g, ''), code);
      const data = unwrap<{ access_token: string; refresh_token: string }>(response, { access_token: '', refresh_token: '' });
      await loginStore(data.access_token, data.refresh_token);
      navigate('/');
    } catch {
      setStatus('登录失败，请确认验证码或后端状态');
    }
  }

  return <main className="od-shell"><div className="od-app"><Sidebar /><section className="od-main"><ScreenTitle eyebrow="ACCOUNT · MVP P0" title="用最少步骤创建家庭空间" desc="手机号登录后会按当前用户拉取家庭、宝宝、记录和照片。" action={<NavLink className="od-btn primary" to="/baby-profile">宝宝档案</NavLink>} /><div className="od-content-grid"><section className="od-panel od-pad"><h2>手机号登录</h2><div className="od-form-grid"><label>手机号<input value={phone} onChange={(event) => setPhone(event.target.value)} inputMode="tel" /></label><label>验证码<div className="od-inline-field"><input value={code} onChange={(event) => setCode(event.target.value)} placeholder="6 位验证码" /><button className="od-btn" type="button" onClick={sendCode}>获取验证码</button></div></label><button className="od-btn primary full" type="button" onClick={login}>登录 / 注册</button><p className="od-muted">{status}</p></div></section><aside className="od-panel od-pad"><h2>隐私确认</h2><p className="od-muted">Yuanzi 会保存宝宝昵称、月龄、喂养睡眠记录和家庭照片。创建家庭空间前，需要明确监护人授权。</p><Choice active title="我已阅读并同意" desc="儿童隐私保护声明、家庭成员共享范围、照片上传说明" /></aside></div></section></div></main>;
}

function BabyProfileScreen() {
  return <main className="od-shell"><div className="od-app"><Sidebar /><section className="od-main"><ScreenTitle eyebrow="BABY PROFILE · FAMILY INIT" title="先把提醒和家庭空间准备好" desc="宝宝档案只收集影响记录体验的必要信息，并自动生成默认提醒。" action={<NavLink className="od-btn primary" to="/">进入首页</NavLink>} /><div className="od-content-grid"><section className="od-panel od-pad"><div className="od-form-grid"><label>宝宝昵称<input defaultValue="圆子" /></label><label>生日<input defaultValue="2025-11-21" /></label><label>喂养方式<select defaultValue="mixed"><option value="mixed">混合喂养</option><option value="formula">配方奶</option><option value="breast">母乳</option></select></label></div></section><aside className="od-record-preview"><div><div className="od-eyebrow">DEFAULTS</div><h2>默认提醒已生成</h2><p>喂奶 4 小时、睡眠计时、每日成长练习会在首页持续提示。</p></div><Metric label="家庭空间" value="已初始化" note="等待邀请" /></aside></div></section></div></main>;
}

function ElderScreen() {
  return <><ScreenTitle eyebrow="ELDER MODE · LARGE TOUCH" title="祖辈也能放心记录" desc="大字体、高对比、减少干扰，只保留一键记录和最近提醒。" action={<NavLink className="od-btn" to="/settings">返回标准设置</NavLink>} /><section className="od-panel od-pad od-elder">{['喂奶', '睡觉', '换尿布'].map((label) => <button key={label}>{label}</button>)}</section></>;
}

function Sidebar() {
  return <aside className="od-side"><NavLink className="od-brand" to="/"><span className="od-logo">圆</span><span>Yuanzi</span></NavLink><section className="od-baby-card"><div><span>圆子</span><b>98%</b></div><p>3 个月 15 天 · 家庭空间 4 人</p></section><nav className="od-nav-list">{navItems.map((item) => <NavLink key={item.to} to={item.to} end={item.to === '/'} className={({ isActive }) => `od-nav-item ${isActive ? 'active' : ''}`}><span>{item.short}</span><div><b>{item.label}</b><small>{item.description}</small></div></NavLink>)}</nav><p className="od-side-note">今天还有 2 项推荐记录未完成。睡前建议补充一次体温和睡眠结束时间。</p></aside>;
}

function MobileNav() {
  return <nav className="od-mobile-nav">{navItems.slice(0, 5).map((item) => <NavLink key={item.to} to={item.to} end={item.to === '/'} className={({ isActive }) => (isActive ? 'active' : '')}><span>{item.short}</span><b>{item.label}</b></NavLink>)}</nav>;
}

function ScreenTitle({ eyebrow, title, desc, action }: { eyebrow: string; title: string; desc: string; action?: ReactNode }) {
  return <header className="od-screen-title"><div><div className="od-eyebrow">{eyebrow}</div><h1>{title}</h1><p>{desc}</p></div>{action}</header>;
}

function SectionHead({ title, note }: { title: string; note: ReactNode }) {
  return <div className="od-section-head"><h2>{title}</h2><small>{note}</small></div>;
}

function TaskRow({ label, value, progress }: { label: string; value: string; progress: number }) {
  return <div className="od-task-row"><span>{label}</span><div><i style={{ width: `${progress}%` }} /></div><b>{value}</b></div>;
}

function StatCard({ label, value, note, icon }: { label: string; value: string; note: string; icon: string }) {
  return <article className="od-stat-card"><div><span>{label}</span><i>{icon}</i></div><strong>{value}</strong><p>{note}</p></article>;
}

function Choice({ title, desc, active = false, onClick }: { title: string; desc: string; active?: boolean; onClick?: () => void }) {
  return <button className={`od-choice ${active ? 'active' : ''}`} onClick={onClick}><strong>{title}</strong><small>{desc}</small></button>;
}

function Segmented({ values, value, onChange }: { values: string[]; value: string; onChange: (value: string) => void }) {
  return <div className="od-segmented">{values.map((item) => <button key={item} className={value === item ? 'active' : ''} onClick={() => onChange(item)}>{item}</button>)}</div>;
}

function Metric({ label, value, note, onClick }: { label: string; value: string; note: string; onClick?: () => void }) {
  return <button className={`od-metric ${onClick ? 'clickable' : ''}`} onClick={onClick}><span>{label}</span><b>{value}</b><small>{note}</small></button>;
}

function Timeline({ records, compact = false, onRecordClick }: { records: YuanziRecord[]; compact?: boolean; onRecordClick?: (record: YuanziRecord) => void }) {
  return <section className={compact ? 'od-panel od-pad od-timeline' : 'od-timeline'}>{records.map((record) => { const item = formatRecord(record); return <article className={`od-event ${onRecordClick ? 'clickable' : ''}`} key={record.id} onClick={() => onRecordClick?.(record)}><span className={`od-dot ${item.tone}`} /><div><strong>{item.title}</strong><p>{item.desc}</p></div><time>{formatClock(record.started_at)}</time></article>; })}</section>;
}

function PhotoTile({ photo, active = false, onClick }: { photo: YuanziPhoto; active?: boolean; onClick?: () => void }) {
  const imageUrl = photo.thumb_url || photo.url;
  return <button className={`od-photo ${active ? 'active' : ''}`} onClick={onClick} style={imageUrl ? { backgroundImage: `url(${imageUrl})` } : undefined}><span>{photo.description || photo.id}</span></button>;
}

function MemberRow({ member, onClick }: { member: FamilyMember; onClick: () => void }) {
  return <button className="od-member clickable" onClick={onClick}><span style={member.avatar_url ? { backgroundImage: `url(${member.avatar_url})` } : undefined} /><div><strong>{member.nickname || member.user_id}</strong><p>{roleText(member.role)}{member.elder_mode ? ' · 祖辈模式' : ''}</p></div><small>查看详情</small></button>;
}

function PickerModal({ title, children, onClose }: { title: string; children: ReactNode; onClose: () => void }) {
  return <div className="od-modal-backdrop" role="dialog" aria-modal="true"><section className="od-modal"><header><h2>{title}</h2><button onClick={onClose}>关闭</button></header>{children}</section></div>;
}

function DateTimePicker({ value, onChange }: { value: string; onChange: (value: string) => void }) {
  return <input className="od-native-picker" type="datetime-local" value={value} onChange={(event) => onChange(event.target.value)} />;
}

function NumberPicker({ min, max, step, value, unit, formatter, onChange }: { min: number; max: number; step: number; value: number; unit: string; formatter?: (value: number) => string; onChange: (value: number) => void }) {
  const values = useMemo(() => {
    const result: number[] = [];
    for (let item = min; item <= max; item += step) result.push(item);
    return result;
  }, [min, max, step]);
  return <div className="od-number-wheel">{values.map((item) => <button key={item} className={item === value ? 'active' : ''} onClick={() => onChange(item)}>{formatter ? formatter(item) : item} {unit}</button>)}</div>;
}

function Chart({ values }: { values: number[] }) {
  const max = Math.max(...values, 1);
  return <div className="od-chart">{values.map((item, index) => <span key={`${item}-${index}`} style={{ height: `${Math.max(16, item / max * 92)}%` }} />)}</div>;
}

function unwrap<T>(response: unknown, fallback: T): T {
  if (isRecord(response) && 'data' in response) return response.data as T;
  return fallback;
}

function unwrapList<T>(response: unknown, fallback: T[]): T[] {
  const data = unwrap<{ list?: T[] } | T[]>(response, fallback);
  return Array.isArray(data) ? data : data.list || fallback;
}

function isRecord(value: unknown): value is Record<string, unknown> {
  return typeof value === 'object' && value !== null;
}

function buildRecordPayload(type: RecordKind, localTime: string, feedingMethod: FeedingMethod, amount: number, temperature: number, note: string) {
  const startedAt = new Date(localTime).toISOString();
  if (type === 'feeding') return { type, started_at: startedAt, content: { type: feedingMethod === 'formula' ? 'formula' : 'breast', side: feedingMethod === 'formula' ? undefined : feedingMethod, amount, unit: 'ml' }, note };
  if (type === 'sleep') return { type, started_at: startedAt, ended_at: new Date(new Date(localTime).getTime() + 90 * 60 * 1000).toISOString(), content: { quality: 'stable', location: 'crib' }, note };
  if (type === 'diaper') return { type, started_at: startedAt, content: { type: 'mixed', color: '黄色', consistency: '糊状' }, note };
  return { type, started_at: startedAt, content: { value: temperature, unit: 'C', position: 'armpit' }, note: `体温 ${temperature}°C。${note}` };
}

function formatRecord(record: Pick<YuanziRecord, 'type' | 'started_at' | 'content' | 'note'>) {
  if (record.type === 'feeding') return { title: `喂奶 ${Number(record.content.amount || 0) || 120}ml`, desc: record.note || '奶粉喂养，状态安稳', tone: 'brand' };
  if (record.type === 'sleep') return { title: '睡眠记录', desc: record.note || '夜间连续睡眠', tone: 'lavender' };
  if (record.type === 'diaper') return { title: '换尿布', desc: record.note || '排泄记录正常', tone: 'mint' };
  return { title: `体温 ${record.content.value || record.content.temperature || 36.5}°C`, desc: record.note || '正常范围，已进入 AI 上下文', tone: 'info' };
}

function coerceRecordKind(value: string | null): RecordKind {
  return value === 'sleep' || value === 'diaper' || value === 'temperature' ? value : 'feeding';
}

function toLocalInputValue(date: Date) {
  const offset = date.getTimezoneOffset();
  return new Date(date.getTime() - offset * 60000).toISOString().slice(0, 16);
}

function formatLocalTime(value: string) {
  return value.replace('T', ' ');
}

function formatTime(value: string) {
  return new Date(value).toLocaleString('zh-CN', { hour12: false });
}

function formatClock(value: string) {
  return new Date(value).toLocaleTimeString('zh-CN', { hour: '2-digit', minute: '2-digit', hour12: false });
}

function roleText(role: string) {
  if (role === 'admin') return '管理员';
  if (role === 'elder') return '祖辈';
  return '照护成员';
}

function round(value: number) {
  return Math.round(value * 10) / 10;
}

function maskPhone(value: string) {
  const digits = value.replace(/\D/g, '').slice(0, 11);
  if (digits.length < 7) return digits;
  return `${digits.slice(0, 3)}****${digits.slice(7)}`;
}
