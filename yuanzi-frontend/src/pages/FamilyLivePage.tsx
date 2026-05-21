import React, { useCallback, useEffect, useMemo, useState } from 'react';
import { api, type ListResponse } from '@/services/api';
import type { Baby, Family, FamilyMember, Photo, Record as BabyRecord, User } from '@/types/models';
import { useAuthStore } from '@/stores/useAuthStore';

type TabKey = 'home' | 'record' | 'stats' | 'ai' | 'family' | 'photos';
type RecordKind = 'feeding' | 'sleep' | 'diaper' | 'excretion' | 'temperature';

interface DailyStats {
  feeding: { count: number; total_amount: number };
  sleep: { count: number; total_hours: number };
  diaper: { count: number };
}

interface SummaryStats {
  range: string;
  dates: string[];
  daily_avg_sleep_hours: number[];
  daytime_single_sleep_hours: number[];
  daily_avg_milk_amount: number[];
  summary: {
    avg_daily_sleep_hours: number;
    avg_daytime_single_sleep_hours: number;
    avg_daily_milk_amount: number;
  };
}

interface AIChatItem {
  id: string;
  question: string;
  answer: string;
  created_at: string;
}

interface PhotoComment {
  id: string;
  photo_id: string;
  nickname: string;
  content: string;
  created_at: string;
}

const tabs: Array<{ key: TabKey; label: string }> = [
  { key: 'home', label: '首页' },
  { key: 'record', label: '记录' },
  { key: 'stats', label: '统计' },
  { key: 'ai', label: 'AI' },
  { key: 'family', label: '家庭' },
  { key: 'photos', label: '照片' },
];

export const FamilyLivePage: React.FC = () => {
  const { user } = useAuthStore();
  const [activeTab, setActiveTab] = useState<TabKey>('home');
  const [recordInitialType, setRecordInitialType] = useState<RecordKind>('feeding');
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');
  const [babies, setBabies] = useState<Baby[]>([]);
  const [currentBabyId, setCurrentBabyId] = useState('');
  const [dailyStats, setDailyStats] = useState<DailyStats | null>(null);
  const [summaryStats, setSummaryStats] = useState<SummaryStats | null>(null);
  const [records, setRecords] = useState<BabyRecord[]>([]);
  const [photos, setPhotos] = useState<Photo[]>([]);
  const [family, setFamily] = useState<Family | null>(null);
  const [members, setMembers] = useState<FamilyMember[]>([]);
  const [aiHistory, setAIHistory] = useState<AIChatItem[]>([]);
  const [commentsByPhoto, setCommentsByPhoto] = useState<Record<string, PhotoComment[]>>({});
  const [photoUploadMessage, setPhotoUploadMessage] = useState('');

  const currentBaby = useMemo(
    () => babies.find((baby) => baby.id === currentBabyId) ?? babies[0],
    [babies, currentBabyId]
  );

  const refresh = useCallback(async () => {
    setLoading(true);
    setError('');
    try {
      const babyList = (await api.baby.getList()) as Baby[];
      setBabies(babyList);
      const baby = babyList.find((item) => item.id === currentBabyId) ?? babyList[0];
      if (!baby) {
        setLoading(false);
        return;
      }
      setCurrentBabyId(baby.id);

      const [daily, summary, recordList, photoList, familyDetail, memberList, chats] = await Promise.all([
        api.record.getDailyStats(baby.id) as Promise<DailyStats>,
        api.record.getSummaryStats(baby.id, { range: 'week' }) as Promise<SummaryStats>,
        api.record.getList(baby.id, { page_size: 10 }) as Promise<ListResponse<BabyRecord>>,
        api.photo.getList(baby.id, { page_size: 30 }) as Promise<ListResponse<Photo>>,
        api.family.getDetail(baby.family_id) as Promise<Family>,
        api.family.getMembers(baby.family_id) as Promise<FamilyMember[]>,
        api.ai.getHistory({ baby_id: baby.id, page_size: 10 }) as Promise<ListResponse<AIChatItem>>,
      ]);

      setDailyStats(daily);
      setSummaryStats(summary);
      setRecords(recordList.list);
      setPhotos(photoList.list);
      setFamily(familyDetail);
      setMembers(memberList);
      setAIHistory(chats.list);
    } catch (err) {
      setError(err instanceof Error ? err.message : '加载家庭数据失败');
    } finally {
      setLoading(false);
    }
  }, [currentBabyId]);

  useEffect(() => {
    void refresh();
  }, [refresh]);

  const displayUser = (user as User | null)?.nickname || '家庭成员';

  return (
    <div className="min-h-screen bg-[#f8fafc] text-[#172033]">
      <div className="mx-auto flex min-h-screen max-w-5xl flex-col">
        <header className="sticky top-0 z-20 border-b border-slate-200 bg-white px-4 py-4">
          <div className="flex flex-wrap items-center justify-between gap-3">
            <div>
              <div className="text-sm text-slate-500">BabyGarden</div>
              <h1 className="text-2xl font-semibold">小园子家庭</h1>
            </div>
            <div className="flex items-center gap-3">
              <select
                value={currentBaby?.id ?? ''}
                onChange={(event) => setCurrentBabyId(event.target.value)}
                className="rounded-md border border-slate-300 bg-white px-3 py-2"
              >
                {babies.map((baby) => (
                  <option key={baby.id} value={baby.id}>
                    {baby.name}
                  </option>
                ))}
              </select>
              <button onClick={() => void refresh()} className="rounded-md bg-slate-900 px-4 py-2 text-white">
                刷新
              </button>
            </div>
          </div>
          <nav className="mt-4 flex gap-2 overflow-x-auto">
            {tabs.map((tab) => (
              <button
                key={tab.key}
                onClick={() => setActiveTab(tab.key)}
                className={`rounded-md px-3 py-2 text-sm font-medium ${
                  activeTab === tab.key ? 'bg-[#ef7868] text-white' : 'bg-slate-100 text-slate-600'
                }`}
              >
                {tab.label}
              </button>
            ))}
          </nav>
        </header>

        <main className="flex-1 px-4 py-6">
          {loading && <Panel>正在读取后端真实数据...</Panel>}
          {error && <Panel tone="danger">{error}</Panel>}
          {!loading && !currentBaby && <Panel>当前账号还没有宝宝数据，请先创建或加入家庭。</Panel>}
          {!loading && currentBaby && (
            <>
              {activeTab === 'home' && (
                <HomeSection
                  userName={displayUser}
                  baby={currentBaby}
                  dailyStats={dailyStats}
                  records={records}
                  photos={photos}
                  onQuickRecord={(type) => {
                    setRecordInitialType(type);
                    setActiveTab('record');
                  }}
                />
              )}
              {activeTab === 'record' && <RecordSection babyId={currentBaby.id} initialType={recordInitialType} />}
              {activeTab === 'stats' && <StatsSection babyId={currentBaby.id} initialSummary={summaryStats} />}
              {activeTab === 'ai' && <AISection babyId={currentBaby.id} history={aiHistory} onUpdated={() => void refresh()} />}
              {activeTab === 'family' && (
                <FamilySection family={family} members={members} onUpdated={() => void refresh()} />
              )}
              {activeTab === 'photos' && (
                <PhotoSection
                  photos={photos}
                  babyId={currentBaby.id}
                  commentsByPhoto={commentsByPhoto}
                  uploadMessage={photoUploadMessage}
                  onUploadMessage={setPhotoUploadMessage}
                  onLoadComments={async (photoId) => {
                    const response = (await api.photo.getComments(photoId)) as ListResponse<PhotoComment>;
                    setCommentsByPhoto((prev) => ({ ...prev, [photoId]: response.list }));
                  }}
                  onChanged={() => void refresh()}
                />
              )}
            </>
          )}
        </main>
      </div>
    </div>
  );
};

const Panel: React.FC<{ children: React.ReactNode; tone?: 'default' | 'danger' }> = ({ children, tone = 'default' }) => (
  <div className={`rounded-lg border p-4 ${tone === 'danger' ? 'border-red-200 bg-red-50 text-red-700' : 'border-slate-200 bg-white'}`}>
    {children}
  </div>
);

const HomeSection: React.FC<{
  userName: string;
  baby: Baby;
  dailyStats: DailyStats | null;
  records: BabyRecord[];
  photos: Photo[];
  onQuickRecord: (type: RecordKind) => void;
}> = ({ userName, baby, dailyStats, records, photos, onQuickRecord }) => (
  <div className="space-y-6">
    <section className="rounded-lg bg-white p-5 shadow-sm">
      <div className="text-sm text-slate-500">你好，{userName}</div>
      <h2 className="mt-1 text-3xl font-semibold">{baby.name}</h2>
      <div className="mt-4 grid grid-cols-3 gap-3">
        <Metric label="今日喝奶" value={`${dailyStats?.feeding.total_amount ?? 0} ml`} />
        <Metric label="今日睡眠" value={`${(dailyStats?.sleep.total_hours ?? 0).toFixed(1)} h`} />
        <Metric label="换尿布" value={`${dailyStats?.diaper.count ?? 0} 次`} />
      </div>
    </section>

    <section className="grid gap-3 md:grid-cols-5">
      {(['feeding', 'diaper', 'excretion', 'temperature', 'sleep'] as RecordKind[]).map((type) => (
        <button key={type} onClick={() => onQuickRecord(type)} className="rounded-lg bg-[#ef7868] px-4 py-3 font-medium text-white">
          {recordLabel(type)}
        </button>
      ))}
    </section>

    <section className="grid gap-4 lg:grid-cols-2">
      <Panel>
        <h3 className="mb-3 font-semibold">最近记录</h3>
        <RecordList records={records} />
      </Panel>
      <Panel>
        <h3 className="mb-3 font-semibold">照片墙动态</h3>
        <div className="grid grid-cols-3 gap-2">
          {photos.slice(0, 6).map((photo) => (
            <img key={photo.id} src={photo.thumb_url || photo.url} alt="" className="aspect-square rounded-md object-cover" />
          ))}
        </div>
      </Panel>
    </section>
  </div>
);

const Metric: React.FC<{ label: string; value: string }> = ({ label, value }) => (
  <div className="rounded-lg bg-slate-50 p-4">
    <div className="text-sm text-slate-500">{label}</div>
    <div className="mt-1 text-xl font-semibold">{value}</div>
  </div>
);

const RecordSection: React.FC<{ babyId: string; initialType: RecordKind }> = ({ babyId, initialType }) => {
  const [type, setType] = useState<RecordKind>('feeding');
  const [startedAt, setStartedAt] = useState(toDateTimeLocal(new Date()));
  const [endedAt, setEndedAt] = useState(toDateTimeLocal(new Date(Date.now() + 60 * 60000)));
  const [amount, setAmount] = useState(120);
  const [duration, setDuration] = useState(60);
  const [temperature, setTemperature] = useState(36.8);
  const [note, setNote] = useState('');
  const [records, setRecords] = useState<BabyRecord[]>([]);
  const [editing, setEditing] = useState<BabyRecord | null>(null);
  const [saving, setSaving] = useState(false);

  const loadRecords = useCallback(async () => {
    const response = (await api.record.getList(babyId, { page_size: 20, type })) as ListResponse<BabyRecord>;
    setRecords(response.list);
  }, [babyId, type]);

  useEffect(() => {
    setType(initialType);
  }, [initialType]);

  useEffect(() => {
    void loadRecords();
  }, [loadRecords]);

  const resetEditing = () => {
    setEditing(null);
    setNote('');
  };

  const save = async () => {
    setSaving(true);
    const started = new Date(startedAt);
    const body: Record<string, unknown> = {
      baby_id: babyId,
      type,
      started_at: started.toISOString(),
      note,
      content: buildRecordContent(type, amount, temperature),
    };
    if (type === 'sleep') {
      body.ended_at = new Date(endedAt).toISOString();
    }
    try {
      if (editing) {
        await api.record.update(editing.id, {
          started_at: body.started_at,
          ended_at: body.ended_at,
          content: body.content,
          note,
        });
      } else {
        await api.record.create(body);
      }
      resetEditing();
      await loadRecords();
    } finally {
      setSaving(false);
    }
  };

  const startEdit = (record: BabyRecord) => {
    setEditing(record);
    setType(record.type as RecordKind);
    setStartedAt(toDateTimeLocal(new Date(record.started_at)));
    if (record.ended_at) setEndedAt(toDateTimeLocal(new Date(record.ended_at)));
    setNote(record.note ?? '');
    if (record.type === 'feeding') setAmount(Number(record.content?.amount ?? amount));
    if (record.type === 'temperature') setTemperature(Number(record.content?.value ?? temperature));
  };

  const deleteRecord = async (record: BabyRecord) => {
    await api.record.delete(record.id);
    await loadRecords();
  };

  return (
    <div className="grid gap-4 lg:grid-cols-2">
      <Panel>
        <div className="grid gap-4 md:grid-cols-2">
          <div>
            <label className="text-sm font-medium">记录类型</label>
            <select disabled={Boolean(editing)} value={type} onChange={(event) => setType(event.target.value as RecordKind)} className="mt-2 w-full rounded-md border px-3 py-2">
              {(['feeding', 'diaper', 'excretion', 'temperature', 'sleep'] as RecordKind[]).map((item) => (
                <option key={item} value={item}>{recordLabel(item)}</option>
              ))}
            </select>
          </div>
          <div>
            <label className="text-sm font-medium">开始时间</label>
            <input type="datetime-local" value={startedAt} onChange={(event) => setStartedAt(event.target.value)} className="mt-2 w-full rounded-md border px-3 py-2" />
          </div>
          {type === 'feeding' && <NumberInput label="奶量 ml" value={amount} onChange={setAmount} />}
          {type === 'sleep' && (
            <>
              <NumberInput label="睡眠分钟" value={duration} onChange={(value) => {
                setDuration(value);
                setEndedAt(toDateTimeLocal(new Date(new Date(startedAt).getTime() + value * 60000)));
              }} />
              <div>
                <label className="text-sm font-medium">结束时间</label>
                <input type="datetime-local" value={endedAt} onChange={(event) => setEndedAt(event.target.value)} className="mt-2 w-full rounded-md border px-3 py-2" />
              </div>
            </>
          )}
          {type === 'temperature' && <NumberInput label="体温 °C" value={temperature} onChange={setTemperature} step="0.1" />}
        </div>
        <textarea value={note} onChange={(event) => setNote(event.target.value)} placeholder="备注" className="mt-4 min-h-24 w-full rounded-md border px-3 py-2" />
        <div className="mt-4 flex gap-2">
          <button disabled={saving} onClick={() => void save()} className="rounded-md bg-slate-900 px-4 py-2 text-white disabled:opacity-50">
            {saving ? '保存中...' : editing ? '保存修改' : '保存记录'}
          </button>
          {editing && <button onClick={resetEditing} className="rounded-md bg-slate-100 px-4 py-2">取消编辑</button>}
        </div>
      </Panel>
      <Panel>
        <h3 className="mb-3 font-semibold">{recordLabel(type)}记录</h3>
        <div className="space-y-2">
          {records.map((record) => (
            <div key={record.id} className="rounded-md bg-slate-50 p-3">
              <div className="flex items-center justify-between gap-3">
                <div>
                  <div className="font-medium">{recordLabel(record.type as RecordKind)}</div>
                  <div className="text-sm text-slate-500">{new Date(record.started_at).toLocaleString()}</div>
                </div>
                <div className="flex gap-2">
                  <button onClick={() => startEdit(record)} className="rounded bg-white px-3 py-1 text-sm">修改</button>
                  <button onClick={() => void deleteRecord(record)} className="rounded bg-red-50 px-3 py-1 text-sm text-red-600">删除</button>
                </div>
              </div>
              {record.note && <div className="mt-2 text-sm text-slate-600">{record.note}</div>}
            </div>
          ))}
        </div>
      </Panel>
    </div>
  );
};

const NumberInput: React.FC<{ label: string; value: number; onChange: (value: number) => void; step?: string }> = ({ label, value, onChange, step = '1' }) => (
  <div>
    <label className="text-sm font-medium">{label}</label>
    <input type="number" step={step} value={value} onChange={(event) => onChange(Number(event.target.value))} className="mt-2 w-full rounded-md border px-3 py-2" />
  </div>
);

const StatsSection: React.FC<{ babyId: string; initialSummary: SummaryStats | null }> = ({ babyId, initialSummary }) => {
  const [range, setRange] = useState('week');
  const [summary, setSummary] = useState<SummaryStats | null>(initialSummary);
  const [startDate, setStartDate] = useState(() => new Date(Date.now() - 6 * 86400000).toISOString().slice(0, 10));
  const [endDate, setEndDate] = useState(() => new Date().toISOString().slice(0, 10));

  useEffect(() => {
    setSummary(initialSummary);
  }, [initialSummary]);

  const load = async (nextRange: string) => {
    setRange(nextRange);
    const params: Record<string, string> = nextRange === 'custom'
      ? { range: nextRange, start_date: startDate, end_date: endDate }
      : { range: nextRange };
    setSummary((await api.record.getSummaryStats(babyId, params)) as SummaryStats);
  };

  const max = Math.max(1, ...(summary?.daily_avg_sleep_hours ?? []), ...(summary?.daily_avg_milk_amount ?? []));

  return (
    <div className="space-y-4">
      <div className="flex gap-2">
        {['day', 'week', 'month', 'custom'].map((item) => (
          <button key={item} onClick={() => void load(item)} className={`rounded-md px-4 py-2 ${range === item ? 'bg-[#ef7868] text-white' : 'bg-white'}`}>
            {item === 'day' ? '日' : item === 'week' ? '周' : item === 'month' ? '月' : '自定义'}
          </button>
        ))}
      </div>
      {range === 'custom' && (
        <div className="flex flex-wrap gap-2">
          <input type="date" value={startDate} onChange={(event) => setStartDate(event.target.value)} className="rounded-md border px-3 py-2" />
          <input type="date" value={endDate} onChange={(event) => setEndDate(event.target.value)} className="rounded-md border px-3 py-2" />
          <button onClick={() => void load('custom')} className="rounded-md bg-slate-900 px-4 py-2 text-white">更新区间</button>
        </div>
      )}
      <Panel>
        <div className="grid gap-3 md:grid-cols-3">
          <Metric label="每日平均睡眠" value={`${summary?.summary.avg_daily_sleep_hours ?? 0} h`} />
          <Metric label="白天单次睡眠" value={`${summary?.summary.avg_daytime_single_sleep_hours ?? 0} h`} />
          <Metric label="日均喝奶量" value={`${summary?.summary.avg_daily_milk_amount ?? 0} ml`} />
        </div>
        <div className="mt-6 grid gap-2">
          {summary?.dates.map((date, index) => (
            <div key={date} className="grid grid-cols-[92px_1fr] items-center gap-3 text-sm">
              <span className="text-slate-500">{date.slice(5)}</span>
              <div className="space-y-1">
                <Bar label="睡眠" value={summary.daily_avg_sleep_hours[index]} max={max} suffix="h" />
                <Bar label="喝奶" value={summary.daily_avg_milk_amount[index]} max={max} suffix="ml" />
              </div>
            </div>
          ))}
        </div>
      </Panel>
    </div>
  );
};

const Bar: React.FC<{ label: string; value: number; max: number; suffix: string }> = ({ label, value, max, suffix }) => (
  <div className="flex items-center gap-2">
    <span className="w-10 text-xs text-slate-500">{label}</span>
    <div className="h-5 flex-1 rounded bg-slate-100">
      <div className="h-5 rounded bg-[#ef7868]" style={{ width: `${Math.min(100, (value / max) * 100)}%` }} />
    </div>
    <span className="w-16 text-right text-xs">{value}{suffix}</span>
  </div>
);

const AISection: React.FC<{ babyId: string; history: AIChatItem[]; onUpdated: () => void }> = ({ babyId, history, onUpdated }) => {
  const [question, setQuestion] = useState('');
  const [answer, setAnswer] = useState('');
  const [asking, setAsking] = useState(false);
  const [trendLoading, setTrendLoading] = useState(false);

  const ask = async () => {
    if (!question.trim()) return;
    setAsking(true);
    setAnswer('');
    try {
      await api.ai.chatStream(question, {
        baby_id: babyId,
        onDelta: (delta) => setAnswer((prev) => prev + delta),
      });
      setQuestion('');
      onUpdated();
    } finally {
      setAsking(false);
    }
  };

  const analyzeTrend = async () => {
    setTrendLoading(true);
    setAnswer('');
    try {
      const stats = await api.record.getSummaryStats(babyId, { range: 'week' }) as SummaryStats;
      const prompt = `请分析小园子近一周统计趋势，必须包含睡眠、奶量、排泄。统计摘要：${JSON.stringify(stats.summary)}，日期：${stats.dates.join(',')}`;
      await api.ai.chatStream(prompt, {
        baby_id: babyId,
        onDelta: (delta) => setAnswer((prev) => prev + delta),
      });
      onUpdated();
    } finally {
      setTrendLoading(false);
    }
  };

  return (
    <div className="grid gap-4 lg:grid-cols-2">
      <Panel>
        <h3 className="mb-3 font-semibold">AI 育儿问答</h3>
        <textarea value={question} onChange={(event) => setQuestion(event.target.value)} className="min-h-28 w-full rounded-md border px-3 py-2" placeholder="输入你的问题" />
        <div className="mt-3 flex flex-wrap gap-2">
          <button disabled={asking || trendLoading} onClick={() => void ask()} className="rounded-md bg-slate-900 px-4 py-2 text-white disabled:opacity-50">
            {asking ? '流式回答中...' : '提问'}
          </button>
          <button disabled={asking || trendLoading} onClick={() => void analyzeTrend()} className="rounded-md bg-[#ef7868] px-4 py-2 text-white disabled:opacity-50">
            {trendLoading ? '分析中...' : '分析近一周趋势'}
          </button>
        </div>
        {answer && <div className="mt-4 rounded-md bg-slate-50 p-3">{answer}</div>}
      </Panel>
      <Panel>
        <h3 className="mb-3 font-semibold">历史会话</h3>
        <div className="space-y-3">
          {history.map((item) => (
            <div key={item.id} className="rounded-md bg-slate-50 p-3">
              <div className="font-medium">{item.question}</div>
              <div className="mt-1 text-sm text-slate-600">{item.answer}</div>
            </div>
          ))}
        </div>
      </Panel>
    </div>
  );
};

const FamilySection: React.FC<{ family: Family | null; members: FamilyMember[]; onUpdated: () => void }> = ({ family, members, onUpdated }) => {
  const [phone, setPhone] = useState('');
  const [inviteCode, setInviteCode] = useState('');

  return (
    <div className="space-y-4">
      <Panel>
        <h3 className="font-semibold">{family?.name ?? '家庭'}</h3>
        <div className="mt-2 text-sm text-slate-500">邀请码：{family?.invite_code ?? '-'}</div>
      </Panel>
      <Panel>
        <h3 className="mb-3 font-semibold">成员管理</h3>
        <div className="grid gap-2">
          {members.map((member) => (
            <div key={member.user_id} className="flex items-center justify-between rounded-md bg-slate-50 p-3">
              <span>{member.nickname || member.user_id}</span>
              <span className="text-sm text-slate-500">{member.role}</span>
            </div>
          ))}
        </div>
        <div className="mt-4 flex gap-2">
          <input value={phone} onChange={(event) => setPhone(event.target.value)} placeholder="手机号邀请成员" className="flex-1 rounded-md border px-3 py-2" />
          <button onClick={() => family && void api.family.invite(family.id, phone).then(onUpdated)} className="rounded-md bg-slate-900 px-4 py-2 text-white">邀请</button>
        </div>
        <div className="mt-3 flex gap-2">
          <input value={inviteCode} onChange={(event) => setInviteCode(event.target.value)} placeholder="邀请码加入家庭" className="flex-1 rounded-md border px-3 py-2" />
          <button onClick={() => void api.family.join(inviteCode).then(onUpdated)} className="rounded-md bg-slate-100 px-4 py-2">加入</button>
          <button onClick={() => family && void api.family.leave(family.id).then(onUpdated)} className="rounded-md bg-red-50 px-4 py-2 text-red-600">离开</button>
        </div>
      </Panel>
    </div>
  );
};

const PhotoSection: React.FC<{
  photos: Photo[];
  babyId: string;
  commentsByPhoto: Record<string, PhotoComment[]>;
  uploadMessage: string;
  onUploadMessage: (message: string) => void;
  onLoadComments: (photoId: string) => Promise<void>;
  onChanged: () => void;
}> = ({ photos, babyId, commentsByPhoto, uploadMessage, onUploadMessage, onLoadComments, onChanged }) => {
  const [commentText, setCommentText] = useState<Record<string, string>>({});
  const [uploading, setUploading] = useState(false);

  const uploadPhoto = async (file: File) => {
    if (!babyId) {
      onUploadMessage('请先选择宝宝后再上传照片');
      return;
    }
    setUploading(true);
    onUploadMessage('正在上传...');
    try {
      const result = await api.photo.getUploadUrl({
        baby_id: babyId,
        filename: file.name,
        content_type: file.type || 'image/jpeg',
        size: file.size,
      }) as { upload_url: string; photo_id: string; upload_headers?: Record<string, string> };
      await fetch(result.upload_url, {
        method: 'PUT',
        headers: result.upload_headers,
        body: file,
      });
      await api.photo.confirmUpload(result.photo_id, file.size);
      onUploadMessage('上传成功');
      onChanged();
    } catch (err) {
      onUploadMessage(err instanceof Error ? err.message : '上传失败');
    } finally {
      setUploading(false);
    }
  };

  return (
    <div className="space-y-4">
      <Panel>
        <label className="inline-flex cursor-pointer items-center rounded-md bg-slate-900 px-4 py-2 text-white">
          {uploading ? '上传中...' : '上传照片'}
          <input
            type="file"
            accept="image/*"
            className="hidden"
            disabled={uploading}
            onChange={(event) => {
              const file = event.target.files?.[0];
              if (file) void uploadPhoto(file);
            }}
          />
        </label>
        {uploadMessage && <span className="ml-3 text-sm text-slate-500">{uploadMessage}</span>}
      </Panel>
      <div className="grid gap-4 md:grid-cols-2 xl:grid-cols-3">
        {photos.map((photo) => (
          <Panel key={photo.id}>
            <img src={photo.url} alt="" className="aspect-square w-full rounded-md object-cover" />
            <div className="mt-3 flex items-center justify-between text-sm text-slate-600">
              <span>{photo.like_count} 赞 · {photo.comment_count} 评论</span>
              <button onClick={() => void (photo.liked_by_me ? api.photo.unlike(photo.id) : api.photo.like(photo.id)).then(onChanged)} className="rounded bg-slate-100 px-3 py-1">
                {photo.liked_by_me ? '取消赞' : '点赞'}
              </button>
            </div>
            <button onClick={() => void onLoadComments(photo.id)} className="mt-2 text-sm text-[#ef7868]">查看评论</button>
            <div className="mt-2 space-y-2">
              {(commentsByPhoto[photo.id] ?? []).map((comment) => (
                <div key={comment.id} className="rounded bg-slate-50 p-2 text-sm">
                  <b>{comment.nickname || '成员'}：</b>{comment.content}
                </div>
              ))}
            </div>
            <div className="mt-3 flex gap-2">
              <input value={commentText[photo.id] ?? ''} onChange={(event) => setCommentText((prev) => ({ ...prev, [photo.id]: event.target.value }))} placeholder="写评论" className="min-w-0 flex-1 rounded-md border px-3 py-2" />
              <button
                onClick={() => void api.photo.comment(photo.id, commentText[photo.id] ?? '').then(async () => {
                  setCommentText((prev) => ({ ...prev, [photo.id]: '' }));
                  await onLoadComments(photo.id);
                  onChanged();
                })}
                className="rounded-md bg-slate-900 px-3 py-2 text-white"
              >
                发送
              </button>
            </div>
          </Panel>
        ))}
      </div>
    </div>
  );
};

const RecordList: React.FC<{ records: BabyRecord[] }> = ({ records }) => (
  <div className="space-y-2">
    {records.map((record) => (
      <div key={record.id} className="rounded-md bg-slate-50 p-3">
        <div className="font-medium">{recordLabel(record.type as RecordKind)}</div>
        <div className="text-sm text-slate-500">{new Date(record.started_at).toLocaleString()}</div>
      </div>
    ))}
  </div>
);

function buildRecordContent(type: RecordKind, amount: number, temperature: number): Record<string, unknown> {
  switch (type) {
    case 'feeding':
      return { type: 'formula', amount, unit: 'ml' };
    case 'diaper':
      return { type: 'wet' };
    case 'excretion':
      return { type: 'poop', amount: 'normal' };
    case 'temperature':
      return { value: temperature, unit: 'celsius', position: 'armpit' };
    case 'sleep':
      return { quality: 'good', location: 'crib' };
    default:
      return {};
  }
}

function recordLabel(type: RecordKind): string {
  const labels: Record<RecordKind, string> = {
    feeding: '喂养',
    sleep: '睡眠',
    diaper: '换尿布',
    excretion: '排泄',
    temperature: '测温',
  };
  return labels[type];
}

function toDateTimeLocal(date: Date): string {
  const offset = date.getTimezoneOffset();
  const local = new Date(date.getTime() - offset * 60000);
  return local.toISOString().slice(0, 16);
}

export default FamilyLivePage;
