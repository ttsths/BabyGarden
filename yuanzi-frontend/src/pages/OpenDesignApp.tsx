import { useEffect, useMemo, useState, type ReactNode } from 'react';
import { NavLink, Route, Routes, useNavigate, useParams, useSearchParams } from 'react-router-dom';
import { api } from '../services/api';
import { useAuthStore } from '../stores/useAuthStore';
import { useThemeStore } from '../stores/useThemeStore';
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
  daily_avg_sleep_hours?: number[];
  daytime_single_sleep_hours?: number[];
  daily_avg_milk_amount?: number[];
  daily_average_sleep_hours?: number[];
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
  like_count?: number;
  comment_count?: number;
  liked_by_me?: boolean;
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

type BabyOption = {
  id: string;
  family_id?: string;
  familyId?: string;
  name?: string;
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

const emptyStats: DailyStats = {
  feeding: { count: 0, total_amount: 0 },
  sleep: { count: 0, total_hours: 0 },
  diaper: { count: 0 },
};

const emptyWeekly: WeeklyStats = {
  dates: [],
  feeding: [],
  sleep: [],
  diaper: [],
};

const today = new Date().toISOString().slice(0, 10);

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
  const [babyId, setBabyId] = useState('');
  const [familyId, setFamilyId] = useState('');
  const [babyName, setBabyName] = useState('小园子');
  const [status, setStatus] = useState<'loading' | 'ready' | 'empty' | 'error'>('loading');

  useEffect(() => {
    let active = true;
    async function load() {
      if (!token && !isAuthenticated) {
        setStatus('empty');
        return;
      }
      setStatus('loading');
      try {
        const response = await api.baby.getList();
        const babies = unwrap<Array<BabyOption>>(response, []);
        const current = babies[0];
        if (!active) return;
        if (!current) {
          setBabyId('');
          setFamilyId('');
          setStatus('empty');
          return;
        }
        const nextBabyId = current.id;
        const nextFamilyId = current.family_id || current.familyId || '';
        setBabyId(nextBabyId);
        setFamilyId(nextFamilyId);
        setBabyName(current.name || '小园子');
        setStatus('ready');
      } catch {
        if (active) setStatus('error');
      }
    }
    void load();
    return () => { active = false; };
  }, [token, isAuthenticated]);

  return { babyId, familyId, babyName, status };
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
  const { babyId, babyName, status: contextStatus } = useYuanziContext();
  const [stats, setStats] = useState<DailyStats>(emptyStats);
  const [records, setRecords] = useState<YuanziRecord[]>([]);
  const [photos, setPhotos] = useState<YuanziPhoto[]>([]);
  const [status, setStatus] = useState<'loading' | 'ready' | 'empty' | 'error'>('loading');

  useEffect(() => {
    let active = true;
    async function load() {
      if (!babyId) {
        setStatus(contextStatus === 'error' ? 'error' : 'empty');
        return;
      }
      setStatus('loading');
      try {
        const [daily, recordList, photoList] = await Promise.all([
          api.stats.daily(babyId, today),
          api.record.getList(babyId, { page: 1, page_size: 4 }),
          api.photo.getList(babyId, { page: 1, page_size: 3 }),
        ]);
        if (!active) return;
        setStats(unwrap<DailyStats>(daily, emptyStats));
        setRecords(unwrapList<YuanziRecord>(recordList));
        setPhotos(unwrapList<YuanziPhoto>(photoList));
        setStatus('ready');
      } catch {
        if (active) setStatus('error');
      }
    }
    void load();
    return () => { active = false; };
  }, [babyId, contextStatus]);

  const progress = Math.min(100, Math.round(((stats.feeding.count / 6) + Math.min(stats.sleep.total_hours / 12, 1) + (stats.diaper.count / 4)) / 3 * 100));
  const latestTemperature = stats.temperature?.latest;

  return (
    <>
      <ScreenTitle
        eyebrow={`TODAY · ${today}`}
        title={`照看今天，也看见${babyName}的成长节奏`}
        desc={`把喂养、睡眠、排泄、体温和家庭照片放到同一个温柔但高效的工作台里。数据源：后端 API${status === 'loading' ? '，正在加载' : ''}`}
        action={<div className="od-actions"><button className="od-btn" onClick={() => navigate('/family')}>家庭成员</button><button className="od-btn primary" onClick={() => navigate('/record')}>快速记录</button></div>}
      />
      {status === 'error' && <InlineNotice tone="danger" text="PC 首页数据加载失败，请检查登录状态或后端接口。" />}
      {status === 'empty' && <InlineNotice text="当前账号还没有宝宝数据，请先创建或加入家庭。" />}
      <div className="od-dashboard">
        <div className="od-stack">
          <section className="od-panel od-hero">
            <div className="od-hero-head">
              <h2>{records.length > 0 ? '后端记录已同步，今日照护节奏可追溯。' : '今天还没有后端记录，先从一次快速记录开始。'}</h2>
              <span>{status === 'ready' ? '后端 API 已同步' : '等待后端数据'}</span>
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
            <StatCard label="排泄" value={`${stats.diaper.count} 次`} note="来自后端今日概览" icon="D" />
            <StatCard label="体温" value={latestTemperature ? `${latestTemperature}` : '--'} note={latestTemperature ? '后端最近测量' : '暂无后端测温'} icon="T" />
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
          <Timeline records={records} compact emptyText="暂无后端记录" />
          <SectionHead title="家庭照片" note="get_api_v1_photo" />
          <div className="od-photo-strip">{photos.length > 0 ? photos.slice(0, 3).map((photo) => <PhotoTile key={photo.id} photo={photo} />) : <EmptyState text="暂无后端照片" />}</div>
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
  const [lastRecord, setLastRecord] = useState<YuanziRecord | null>(null);
  const [status, setStatus] = useState('LIVE PREVIEW');

  useEffect(() => {
    setSelectedType(coerceRecordKind(searchParams.get('type')));
  }, [searchParams]);

  useEffect(() => {
    let active = true;
    async function loadLatest() {
      if (!babyId) return;
      try {
        const response = await api.record.getList(babyId, { page: 1, page_size: 1 });
        const [latest] = unwrapList<YuanziRecord>(response);
        if (active) {
          setLastRecord(latest ?? null);
          setStatus(latest ? 'LIVE PREVIEW · 后端最近记录' : 'LIVE PREVIEW · 暂无后端记录');
        }
      } catch {
        if (active) setStatus('LIVE PREVIEW · 后端读取失败');
      }
    }
    void loadLatest();
    return () => { active = false; };
  }, [babyId]);

  async function saveRecord() {
    const payload = buildRecordPayload(selectedType, startedAt, feedingMethod, amount, temperature, note);
    try {
      if (!babyId) throw new Error('missing baby id');
      const response = await api.record.create({ baby_id: babyId, ...payload });
      const saved = unwrap<YuanziRecord>(response, { id: '', baby_id: babyId, ...payload });
      setLastRecord(saved);
      setStatus('已保存到后端');
    } catch {
      setStatus('保存失败，未写入后端');
    }
  }

  const previewPayload = buildRecordPayload(selectedType, startedAt, feedingMethod, amount, temperature, note);
  const preview = formatRecord(lastRecord && lastRecord.type === selectedType ? lastRecord : previewPayload);

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
          <div><Metric label="保存状态" value={status.includes('失败') ? '失败' : '待写入'} note="以接口返回为准" /><Metric label="宝宝ID" value={babyId || '--'} note="后端记录归属" /><Metric label="记录类型" value={recordLabel(selectedType)} note="当前表单" /></div>
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
  const [daily, setDaily] = useState<DailyStats>(emptyStats);
  const [weekly, setWeekly] = useState<WeeklyStats>(emptyWeekly);
  const [records, setRecords] = useState<YuanziRecord[]>([]);
  const [status, setStatus] = useState<'loading' | 'ready' | 'empty' | 'error'>('loading');
  const [rangeStart, setRangeStart] = useState(today);
  const [rangeEnd, setRangeEnd] = useState(today);

  useEffect(() => {
    let active = true;
    async function load() {
      if (!babyId) {
        setStatus('empty');
        return;
      }
      setStatus('loading');
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
        setDaily(unwrap<DailyStats>(dailyRes, emptyStats));
        setWeekly(unwrap<WeeklyStats>(rangeRes, emptyWeekly));
        setRecords(unwrapList<YuanziRecord>(recordRes));
        setStatus('ready');
      } catch {
        if (active) setStatus('error');
      }
    }
    void load();
    return () => { active = false; };
  }, [babyId, mode, rangeStart, rangeEnd]);

  const chartValues = mode === '日'
    ? [daily.feeding.count, daily.sleep.total_hours, daily.diaper.count, 1]
    : mode === '周'
      ? weekly.daily_avg_milk_amount || weekly.daily_average_milk_amount || weekly.feeding
      : mode === '月'
        ? weekly.daily_avg_sleep_hours || weekly.daily_average_sleep_hours || weekly.sleep
        : weekly.daytime_single_sleep_hours || weekly.sleep;

  return (
    <>
      <ScreenTitle eyebrow="TIMELINE · STATS" title="从记录流看到趋势" desc="保留时间轴的可追溯性，同时展示平均睡眠、白天单次睡眠和日均喝奶量。" action={<Segmented values={['日', '周', '月', '自定义']} value={mode} onChange={(value) => setMode(value as StatsMode)} />} />
      {status === 'loading' && <InlineNotice text="正在从后端加载统计数据..." />}
      {status === 'error' && <InlineNotice tone="danger" text="统计接口加载失败。" />}
      <div className="od-content-grid">
        <section className="od-panel od-pad"><h2>今日时间轴</h2><Timeline records={records} emptyText="暂无后端时间轴记录" onRecordClick={(record) => navigate(`/record-detail/${record.id}`)} /></section>
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
  const [photoList, setPhotoList] = useState<YuanziPhoto[]>([]);
  const [selectedPhoto, setSelectedPhoto] = useState<YuanziPhoto | null>(null);
  const [comments, setComments] = useState<PhotoComment[]>([]);
  const [status, setStatus] = useState<'loading' | 'ready' | 'empty' | 'error'>('loading');
  const [uploadOpen, setUploadOpen] = useState(false);
  const [uploadStatus, setUploadStatus] = useState('');
  const [comment, setComment] = useState('');

  useEffect(() => {
    async function load() {
      if (!babyId) {
        setStatus('empty');
        return;
      }
      setStatus('loading');
      try {
        const response = await api.photo.getList(babyId, { page: 1, page_size: 20 });
        const list = unwrapList<YuanziPhoto>(response);
        setPhotoList(list);
        setSelectedPhoto(list[0] ?? null);
        setStatus(list.length > 0 ? 'ready' : 'empty');
      } catch {
        setStatus('error');
      }
    }
    void load();
  }, [babyId]);

  useEffect(() => {
    let active = true;
    async function loadComments() {
      if (!selectedPhoto) {
        setComments([]);
        return;
      }
      try {
        const response = await api.photo.getComments(selectedPhoto.id);
        if (active) setComments(unwrapList<PhotoComment>(response));
      } catch {
        if (active) setComments([]);
      }
    }
    void loadComments();
    return () => { active = false; };
  }, [selectedPhoto]);

  async function requestUpload(file: File) {
    setUploadStatus('正在向后端申请上传地址...');
    try {
      if (!babyId) throw new Error('missing baby id');
      const response = await api.photo.getUploadUrl({ baby_id: babyId, filename: file.name, content_type: file.type || 'image/jpeg', size: file.size });
      const data = unwrap<Record<string, string | number>>(response, {});
      const uploadUrl = String(data.upload_url || '');
      const photoId = String(data.photo_id || '');
      if (!uploadUrl || !photoId) throw new Error('missing upload data');
      await fetch(uploadUrl, { method: 'PUT', body: file });
      await api.photo.confirmUpload(photoId, file.size);
      const nextList = unwrapList<YuanziPhoto>(await api.photo.getList(babyId, { page: 1, page_size: 20 }));
      setPhotoList(nextList);
      setSelectedPhoto(nextList[0] ?? null);
      setUploadStatus('上传成功，已从后端刷新照片墙。');
    } catch {
      setUploadStatus('上传失败，未写入后端。');
    }
  }

  async function toggleLike() {
    if (!selectedPhoto) return;
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
    if (!selectedPhoto || !comment.trim()) return;
    try {
      await api.photo.comment(selectedPhoto.id, comment.trim());
      const nextComments = unwrapList<PhotoComment>(await api.photo.getComments(selectedPhoto.id));
      const updated = {
        ...selectedPhoto,
        comment_count: nextComments.length || (selectedPhoto.comment_count || 0) + 1,
      };
      setSelectedPhoto(updated);
      setPhotoList((items) => items.map((item) => item.id === updated.id ? updated : item));
      setComments(nextComments);
      setComment('');
    } catch {
      setUploadStatus('评论接口暂不可用。');
    }
  }

  return (
    <>
      <ScreenTitle eyebrow="PHOTO WALL · FAMILY MEMORY" title="家庭照片和成长记录在一起" desc="照片墙保留上传状态、家庭互动和关联记录，避免照片只是散落在相册里。" action={<button className="od-btn primary" onClick={() => setUploadOpen(true)}>上传照片</button>} />
      {status === 'loading' && <InlineNotice text="正在从后端加载照片墙..." />}
      {status === 'error' && <InlineNotice tone="danger" text="照片接口加载失败。" />}
      <div className="od-content-grid">
        <section className="od-panel od-pad"><div className="od-photo-grid">{photoList.length > 0 ? photoList.map((photo) => <PhotoTile key={photo.id} photo={photo} active={selectedPhoto?.id === photo.id} onClick={() => setSelectedPhoto(photo)} />) : <EmptyState text="暂无后端照片" />}</div></section>
        <aside className="od-stack">
          <section className="od-panel od-pad"><h2>照片详情</h2>{selectedPhoto ? <><Metric label="当前照片" value={selectedPhoto.description || selectedPhoto.id} note={selectedPhoto.taken_at || '家庭空间'} /><Metric label="点赞" value={`${selectedPhoto.like_count || 0}`} note={selectedPhoto.liked_by_me ? '我已点赞' : '点击下方点赞'} /><Metric label="评论" value={`${selectedPhoto.comment_count || comments.length}`} note="家庭成员互动" /><button className="od-btn full" onClick={toggleLike}>{selectedPhoto.liked_by_me ? '取消点赞' : '点赞'}</button></> : <EmptyState text="请选择一张后端照片" />}</section>
          <section className="od-invite-card"><h2>家庭互动</h2>{comments.map((item) => <p key={item.id}><b>{item.nickname || '家人'}：</b>{item.content}</p>)}{comments.length === 0 && <p>暂无后端评论</p>}<div className="od-ask-box light"><input value={comment} onChange={(event) => setComment(event.target.value)} placeholder="在这张照片下评论" /><button type="button" onClick={sendComment} disabled={!selectedPhoto}>发送</button></div></section>
        </aside>
      </div>
      {uploadOpen && <PickerModal title="上传照片" onClose={() => setUploadOpen(false)}><div className="od-upload-box"><input type="file" accept="image/*" onChange={(event) => { const file = event.target.files?.[0]; if (file) void requestUpload(file); }} /><p>{uploadStatus || '选择照片后会调用 post_api_v1_photo_upload_url 获取直传地址。'}</p></div></PickerModal>}
    </>
  );
}

function FamilyScreen() {
  const { familyId } = useYuanziContext();
  const [members, setMembers] = useState<FamilyMember[]>([]);
  const [status, setStatus] = useState<'loading' | 'ready' | 'empty' | 'error'>('loading');
  const [inviteOpen, setInviteOpen] = useState(false);
  const [selectedMember, setSelectedMember] = useState<FamilyMember | null>(null);
  const [invitePhone, setInvitePhone] = useState('');
  const [inviteRole, setInviteRole] = useState<'member' | 'elder'>('member');
  const [inviteResult, setInviteResult] = useState('');
  const [joinCode, setJoinCode] = useState('');

  useEffect(() => {
    async function load() {
      if (!familyId) {
        setStatus('empty');
        return;
      }
      setStatus('loading');
      try {
        const response = await api.family.getMembers(familyId);
        const list = unwrap<FamilyMember[]>(response, []);
        setMembers(list);
        setStatus(list.length > 0 ? 'ready' : 'empty');
      } catch {
        setStatus('error');
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
      {status === 'loading' && <InlineNotice text="正在从后端加载家庭成员..." />}
      {status === 'error' && <InlineNotice tone="danger" text="家庭成员接口加载失败。" />}
      <div className="od-content-grid">
        <section className="od-panel od-pad od-stack">{members.length > 0 ? members.map((member) => <MemberRow key={member.user_id} member={member} onClick={() => setSelectedMember(member)} />) : <EmptyState text="暂无后端家庭成员" />}</section>
        <aside className="od-stack"><section className="od-invite-card"><h2>邀请链接</h2><p>链接 24 小时有效。加入前需确认儿童隐私声明和照片可见范围。</p><button className="od-btn primary" onClick={() => setInviteOpen(true)}>按手机号邀请</button></section><section className="od-panel od-pad"><h2>加入 / 离开</h2><div className="od-ask-box light"><input value={joinCode} onChange={(event) => setJoinCode(event.target.value.toUpperCase())} placeholder="输入邀请码" /><button type="button" onClick={joinFamily}>加入</button></div><button className="od-btn full" onClick={leaveFamily}>离开家庭</button><p className="od-muted">{inviteResult}</p></section><section className="od-panel od-pad"><h2>同步状态</h2><Metric label="记录" value="准实时" note="SSE connected" /><Metric label="照片" value="排队上传" note="2 张待同步" /></section></aside>
      </div>
      {inviteOpen && <PickerModal title="邀请家人" onClose={() => setInviteOpen(false)}><div className="od-form-grid"><label>手机号<input value={invitePhone} onChange={(event) => setInvitePhone(event.target.value)} placeholder="请输入 11 位手机号" /></label><label>角色<select value={inviteRole} onChange={(event) => setInviteRole(event.target.value as 'member' | 'elder')}><option value="member">照护成员</option><option value="elder">祖辈模式</option></select></label><button className="od-btn primary full" onClick={inviteMember}>发送邀请</button><p className="od-muted">{inviteResult}</p></div></PickerModal>}
      {selectedMember && <PickerModal title="家庭成员详情" onClose={() => setSelectedMember(null)}><Metric label="成员" value={selectedMember.nickname || selectedMember.user_id} note={selectedMember.user_id} /><Metric label="角色" value={roleText(selectedMember.role)} note={selectedMember.elder_mode ? '已开启祖辈模式' : '标准模式'} /></PickerModal>}
    </>
  );
}

function SettingsScreen() {
  const user = useAuthStore((state) => state.user);
  const { isDarkMode, toggleDarkMode } = useThemeStore();
  const [mode, setMode] = useState('标准模式');
  const [phoneOpen, setPhoneOpen] = useState(false);
  const [picker, setPicker] = useState<'feed' | 'practice' | null>(null);
  const [feedInterval, setFeedInterval] = useState(4);
  const [practiceHour, setPracticeHour] = useState(20);
  const [phoneOverride, setPhoneOverride] = useState('');
  const displayedPhone = phoneOverride || maskPhone(user?.phone || '');

  return (
    <>
      <ScreenTitle eyebrow="SETTINGS · PRIVACY" title="把高风险设置放在清晰分组里" desc="设置页覆盖提醒、显示模式、家庭权限和隐私声明；儿童数据相关动作必须说明影响范围。" action={<button className="od-btn primary">保存设置</button>} />
      <div className="od-content-grid equal">
        <section className="od-panel od-pad"><h2>照护提醒</h2><Metric label="喂奶间隔" value={`${feedInterval} 小时`} note="点击更换" onClick={() => setPicker('feed')} /><Metric label="睡眠结束" value="开启" note="计时器提醒" /><Metric label="成长练习" value={`${practiceHour}:00`} note="点击选择小时" onClick={() => setPicker('practice')} /></section>
        <section className="od-panel od-pad od-stack"><h2>显示与模式</h2><Choice active={isDarkMode} title="夜间暗色模式" desc="一键切换全局暗色模式" onClick={() => { toggleDarkMode(); setMode('夜间暗色模式'); }} />{['祖辈极简模式', '标准模式'].map((item) => <Choice key={item} active={mode === item} title={item} desc={item === '祖辈极简模式' ? '大字体、高对比、一键记录' : '点击切换显示偏好'} onClick={() => setMode(item)} />)}</section>
        <section className="od-panel od-pad"><h2>隐私与数据</h2><p className="od-muted">展示儿童隐私保护声明、照片上传范围、家庭成员权限、AI 问答上下文使用说明。MVP 不提供数据导出入口。</p></section>
        <section className="od-panel od-pad"><h2>账号安全</h2><Metric label="手机号" value={displayedPhone || '--'} note="点击更换" onClick={() => setPhoneOpen(true)} /><Metric label="登录设备" value="--" note="等待后端设备接口" /></section>
      </div>
      {picker === 'feed' && <PickerModal title="喂奶间隔" onClose={() => setPicker(null)}><NumberPicker min={2} max={6} step={1} value={feedInterval} unit="小时" onChange={setFeedInterval} /></PickerModal>}
      {picker === 'practice' && <PickerModal title="成长练习提醒" onClose={() => setPicker(null)}><NumberPicker min={6} max={23} step={1} value={practiceHour} unit="点" onChange={setPracticeHour} /></PickerModal>}
      {phoneOpen && <PickerModal title="更换手机号" onClose={() => setPhoneOpen(false)}><div className="od-form-grid"><label>新手机号<input placeholder="请输入新手机号" onChange={(event) => setPhoneOverride(maskPhone(event.target.value))} /></label><button className="od-btn primary full" onClick={() => setPhoneOpen(false)}>确认更换</button></div></PickerModal>}
    </>
  );
}

function RecordsScreen() {
  const navigate = useNavigate();
  const { babyId } = useYuanziContext();
  const [records, setRecords] = useState<YuanziRecord[]>([]);
  const [status, setStatus] = useState<'loading' | 'ready' | 'empty' | 'error'>('loading');
  useEffect(() => {
    async function load() {
      if (!babyId) {
        setStatus('empty');
        return;
      }
      setStatus('loading');
      try {
        const response = await api.record.getList(babyId, { page: 1, page_size: 30 });
        const list = unwrapList<YuanziRecord>(response);
        setRecords(list);
        setStatus(list.length > 0 ? 'ready' : 'empty');
      } catch {
        setStatus('error');
      }
    }
    void load();
  }, [babyId]);
  return <><ScreenTitle eyebrow="RECORD DETAIL · LIST" title="记录明细" desc="从首页最近记录进入，查看全部记录并继续打开单条详情。" />{status === 'loading' && <InlineNotice text="正在从后端加载记录明细..." />}{status === 'error' && <InlineNotice tone="danger" text="记录列表接口加载失败。" />}<section className="od-panel od-pad"><Timeline records={records} emptyText="暂无后端记录" onRecordClick={(record) => navigate(`/record-detail/${record.id}`)} /></section></>;
}

function RecordDetailScreen() {
  const { id } = useParams();
  const [record, setRecord] = useState<YuanziRecord | null>(null);
  const [status, setStatus] = useState<'loading' | 'ready' | 'error'>('loading');
  useEffect(() => {
    async function load() {
      if (!id) return;
      setStatus('loading');
      try {
        const response = await api.record.getDetail(id);
        setRecord(unwrap<YuanziRecord | null>(response, null));
        setStatus('ready');
      } catch {
        setStatus('error');
      }
    }
    void load();
  }, [id]);
  if (!record) {
    return <><ScreenTitle eyebrow="RECORD DETAIL" title="记录详情" desc="从后端读取单条记录。" action={<NavLink className="od-btn" to="/records">返回明细</NavLink>} />{status === 'error' ? <InlineNotice tone="danger" text="记录详情接口加载失败。" /> : <InlineNotice text="正在从后端加载记录详情..." />}</>;
  }
  const display = formatRecord(record);
  return <><ScreenTitle eyebrow="RECORD DETAIL" title={display.title} desc={display.desc} action={<NavLink className="od-btn" to="/records">返回明细</NavLink>} /><section className="od-panel od-pad"><Metric label="记录ID" value={record.id} note={record.type} /><Metric label="时间" value={formatTime(record.started_at)} note={record.ended_at ? formatTime(record.ended_at) : '无结束时间'} /><pre className="od-json">{JSON.stringify(record.content, null, 2)}</pre></section></>;
}

function AiScreen() {
  const { babyId } = useYuanziContext();
  const [question, setQuestion] = useState('');
  const [answer, setAnswer] = useState('');
  const [history, setHistory] = useState<Array<{ id: string; question: string; answer: string; created_at: string }>>([]);
  const [status, setStatus] = useState('历史会话加载中');

  useEffect(() => {
    async function load() {
      try {
        const response = await api.ai.history({ baby_id: babyId, page: 1, page_size: 8 });
        setHistory(unwrapList(response));
        setStatus('历史会话');
      } catch {
        setStatus('历史会话暂不可用');
      }
    }
    void load();
  }, [babyId]);

  async function sendQuestion() {
    if (!question.trim()) return;
    setStatus('AI 正在流式回答');
    setAnswer('');
    try {
      await api.ai.chatStream(question.trim(), {
        baby_id: babyId || undefined,
        history: history.slice(0, 4).flatMap((item) => [
          { role: 'user', content: item.question },
          { role: 'assistant', content: item.answer },
        ]),
        onDelta: (delta) => {
          setAnswer((prev) => prev + delta);
        },
      });
      const latestHistory = await api.ai.history({ baby_id: babyId, page: 1, page_size: 8 });
      setHistory(unwrapList(latestHistory));
      setQuestion('');
      setStatus('已保存到历史会话');
    } catch {
      setStatus('AI 服务暂不可用');
    }
  }

  return <><ScreenTitle eyebrow="AI CHAT · STREAM" title="用文字和 AI 一起看趋势" desc="AI 回答引用近期记录，但避免医疗诊断式结论；历史会话由后端按用户保存。" action={<button className="od-btn primary" onClick={sendQuestion}>发送文字问题</button>} /><div className="od-content-grid"><section className="od-panel od-pad od-stack"><article className="od-event"><span className="od-dot brand" /><div><strong>{history[0]?.question || '暂无后端历史会话'}</strong><p>{status}</p></div><time>{history[0] ? formatClock(history[0].created_at) : '--:--'}</time></article><article className="od-advice"><h2>Yuanzi 建议</h2><p>{answer || '发送问题后，这里会展示后端 AI 流式响应。'}</p></article><form className="od-ask-box light"><input value={question} onChange={(event) => setQuestion(event.target.value)} placeholder="继续追问：今晚睡前怎么安排？" /><button type="button" onClick={sendQuestion}>发送</button></form></section><aside className="od-voice"><h2>历史 AI 会话</h2><div className="od-stack">{history.length > 0 ? history.map((item) => <button className="od-choice" key={item.id} onClick={() => setAnswer(item.answer)}><strong>{item.question}</strong><small>{formatTime(item.created_at)}</small></button>) : <EmptyState text="暂无后端 AI 会话" />}</div></aside></div></>;
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
  const { babyId } = useYuanziContext();
  const [baby, setBaby] = useState<BabyOption | null>(null);
  const [status, setStatus] = useState<'loading' | 'ready' | 'empty' | 'error'>('loading');

  useEffect(() => {
    let active = true;
    async function load() {
      if (!babyId) {
        setStatus('empty');
        return;
      }
      try {
        const response = await api.baby.getDetail(babyId);
        if (active) {
          setBaby(unwrap<BabyOption | null>(response, null));
          setStatus('ready');
        }
      } catch {
        if (active) setStatus('error');
      }
    }
    void load();
    return () => { active = false; };
  }, [babyId]);

  return <main className="od-shell"><div className="od-app"><Sidebar /><section className="od-main"><ScreenTitle eyebrow="BABY PROFILE · API" title="宝宝档案来自后端" desc="PC 端宝宝档案只展示后端接口返回的资料。" action={<NavLink className="od-btn primary" to="/">进入首页</NavLink>} />{status === 'loading' && <InlineNotice text="正在从后端加载宝宝档案..." />}{status === 'error' && <InlineNotice tone="danger" text="宝宝档案接口加载失败。" />}<div className="od-content-grid"><section className="od-panel od-pad">{baby ? <div className="od-form-grid"><label>宝宝昵称<input value={baby.name || ''} readOnly /></label><label>宝宝ID<input value={baby.id} readOnly /></label><label>家庭ID<input value={baby.family_id || baby.familyId || ''} readOnly /></label></div> : <EmptyState text="暂无后端宝宝档案" />}</section><aside className="od-record-preview"><div><div className="od-eyebrow">BACKEND</div><h2>{status === 'ready' ? '后端档案已同步' : '等待后端数据'}</h2><p>创建和修改宝宝档案请以后台接口为准。</p></div></aside></div></section></div></main>;
}

function ElderScreen() {
  return <><ScreenTitle eyebrow="ELDER MODE · LARGE TOUCH" title="祖辈也能放心记录" desc="大字体、高对比、减少干扰，只保留一键记录和最近提醒。" action={<NavLink className="od-btn" to="/settings">返回标准设置</NavLink>} /><section className="od-panel od-pad od-elder">{['喂奶', '睡觉', '换尿布'].map((label) => <button key={label}>{label}</button>)}</section></>;
}

function Sidebar() {
  const { babyName, status } = useYuanziContext();
  const user = useAuthStore((state) => state.user);
  return <aside className="od-side"><NavLink className="od-brand" to="/"><span className="od-logo">圆</span><span>Yuanzi</span></NavLink><section className="od-baby-card"><div><span>{babyName}</span><b>{status === 'ready' ? 'API' : '--'}</b></div><p>{user?.nickname || user?.phone || '家庭成员'} · {status === 'ready' ? '后端家庭空间已连接' : '等待后端家庭数据'}</p></section><nav className="od-nav-list">{navItems.map((item) => <NavLink key={item.to} to={item.to} end={item.to === '/'} className={({ isActive }) => `od-nav-item ${isActive ? 'active' : ''}`}><span>{item.short}</span><div><b>{item.label}</b><small>{item.description}</small></div></NavLink>)}</nav><p className="od-side-note">PC 端页面只展示后端接口返回的数据；接口失败时会显示空状态或错误提示。</p></aside>;
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

function InlineNotice({ text, tone = 'default' }: { text: string; tone?: 'default' | 'danger' }) {
  return <div className={`od-inline-notice ${tone === 'danger' ? 'danger' : ''}`}>{text}</div>;
}

function EmptyState({ text }: { text: string }) {
  return <div className="od-empty-state">{text}</div>;
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

function Timeline({ records, compact = false, emptyText = '暂无记录', onRecordClick }: { records: YuanziRecord[]; compact?: boolean; emptyText?: string; onRecordClick?: (record: YuanziRecord) => void }) {
  return <section className={compact ? 'od-panel od-pad od-timeline' : 'od-timeline'}>{records.length > 0 ? records.map((record) => { const item = formatRecord(record); return <article className={`od-event ${onRecordClick ? 'clickable' : ''}`} key={record.id} onClick={() => onRecordClick?.(record)}><span className={`od-dot ${item.tone}`} /><div><strong>{item.title}</strong><p>{item.desc}</p></div><time>{formatClock(record.started_at)}</time></article>; }) : <EmptyState text={emptyText} />}</section>;
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
  if (response === undefined || response === null) return fallback;
  if (isRecord(response) && 'data' in response) return response.data as T;
  return response as T;
}

function unwrapList<T>(response: unknown): T[] {
  if (Array.isArray(response)) return response as T[];
  if (isRecord(response) && Array.isArray(response.list)) return response.list as T[];
  if (isRecord(response) && 'data' in response) return unwrapList<T>(response.data);
  return [];
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

function recordLabel(type: RecordKind) {
  if (type === 'feeding') return '喂奶';
  if (type === 'sleep') return '睡眠';
  if (type === 'diaper') return '排泄';
  return '测温';
}

function round(value: number) {
  return Math.round(value * 10) / 10;
}

function maskPhone(value: string) {
  const digits = value.replace(/\D/g, '').slice(0, 11);
  if (digits.length < 7) return digits;
  return `${digits.slice(0, 3)}****${digits.slice(7)}`;
}
